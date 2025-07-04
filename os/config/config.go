package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config is the global configuration instance.
// 全局配置实例
var Config *ConfigStr
var once sync.Once

// ConfigStr 封装viper配置结构体
// ConfigStr wraps the viper configuration instance.
type ConfigStr struct {
	Viper    *viper.Viper
	filePath []string
	UserName string
	Password string
}

// New 初始化配置实例（单例）
// New creates a singleton configuration instance.
// 参数 path 可传入一个或多个配置文件路径，具体逻辑参见 setupConfigFile。
// Parameter path: one or more configuration file paths.
func New(path ...string) *ConfigStr {
	once.Do(func() {
		Config = &ConfigStr{
			Viper:    viper.New(),
			filePath: path,
		}
		if len(path) > 0 {
			Config.setupConfigFile(path)
		}
	})
	return Config
}

// AddConfigFile 加载一个新的配置文件并合并到全局配置中
// AddConfigFile loads a new configuration file and merges it into the global configuration.
func AddConfigFile(path string) error {
	newViper := viper.New()
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	newViper.SetConfigType(ext)
	newViper.SetConfigFile(path)
	if err := newViper.ReadInConfig(); err != nil {
		return err
	}
	newViper.WatchConfig()

	// 合并配置
	Config.Viper.MergeConfigMap(newViper.AllSettings())

	// 文件变化时自动重新加载并合并配置
	newViper.OnConfigChange(func(e fsnotify.Event) {
		if err := newViper.ReadInConfig(); err == nil {
			Config.Viper.MergeConfigMap(newViper.AllSettings())
		} else {
			alog.Error(context.Background(), "Viper ReadInConfig error", zap.Error(err))
		}
	})
	return nil
}

// setupConfigFile 根据传入的路径参数初始化配置文件或远程配置
// setupConfigFile initializes local or remote configuration based on the provided paths.
func (c *ConfigStr) setupConfigFile(path []string) {
	switch len(path) {
	case 1:
		// 单个配置文件
		c.Viper.SetConfigFile(path[0])
		c.Viper.WatchConfig()
	case 2:
		// 配置类型由第二个参数的文件扩展名决定
		c.Viper.SetConfigType(strings.TrimPrefix(filepath.Ext(path[len(path)-1]), "."))
	case 3:
		// 远程配置：无安全认证
		if err := c.Viper.AddRemoteProvider(path[0], path[1], path[2]); err != nil {
			panic(err)
		}
	case 4:
		// 远程配置：带安全认证
		if err := c.Viper.AddSecureRemoteProvider(path[0], path[1], path[2], path[3]); err != nil {
			panic(err)
		}
	}
}

// Etcd3 连接ETCD3并加载配置
// Etcd3 connects to ETCD3, loads configuration into viper, and starts watching for changes.
// Parameters:
//   - hosts: etcd 服务器地址列表
//   - paths: 配置在 etcd 中的键列表（可包含文件名和类型信息，如 "config.toml"）
//   - username: etcd 用户名
//   - pwd: etcd 密码
func (c *ConfigStr) Etcd3(hosts, paths []string, username, pwd string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   hosts,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		Username:    username,
		Password:    pwd,
		Logger:      alog.Write,
	})
	if err != nil {
		return err
	}

	if err = c.loadEtcdToViper(paths, cli); err != nil {
		return err
	}
	go c.watchEtcd3(paths, cli)
	return nil
}

// loadEtcdToViper 从 etcd 中加载配置并合并到 viper 中
// loadEtcdToViper loads configuration from etcd keys and merges them into viper.
func (c *ConfigStr) loadEtcdToViper(paths []string, client *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for idx, pathKey := range paths {
		// 初始化一个新的 viper 实例读取远程配置
		newViper := viper.New()
		// 根据路径后缀设置配置类型，默认 toml
		parts := strings.Split(pathKey, ".")
		if len(parts) == 2 {
			newViper.SetConfigType(parts[1])
		} else {
			newViper.SetConfigType("toml")
		}

		resp, err := client.Get(ctx, pathKey, clientv3.WithPrefix())
		if err != nil {
			return err
		}
		// 如果没有获取到配置内容，则跳过该 key
		if len(resp.Kvs) == 0 {
			alog.Warn(context.Background(), fmt.Sprintf("No configuration found for key: %s", pathKey))
			continue
		}

		if err = newViper.ReadConfig(bytes.NewReader(resp.Kvs[0].Value)); err != nil {
			return err
		}

		filename := strings.TrimSuffix(filepath.Base(pathKey), filepath.Ext(pathKey))
		// 根据索引决定合并到全局配置的方式
		for key, value := range newViper.AllSettings() {
			if idx == 0 {
				c.Viper.Set(key, value)
			}
			c.Viper.Set(fmt.Sprintf("%s.%s", filename, key), value)
		}
	}
	return nil
}

// watchEtcd3 监听 etcd 配置变化，并更新到 viper 中
// watchEtcd3 watches the etcd configuration keys and updates the viper configuration on changes.
func (c *ConfigStr) watchEtcd3(paths []string, cli *clientv3.Client) {
	// 为每个 etcd 配置 key 启动一个独立的 watcher
	for idx, pathKey := range paths {
		// 捕获循环变量
		index := idx
		key := pathKey

		go func(index int, key string) {
			newViper := viper.New()
			watcher := clientv3.NewWatcher(cli)
			// 这里可以根据需求设置更长的 context 以及退出条件
			ctx := context.Background()
			defer watcher.Close()
			for watchChan := range watcher.Watch(ctx, key) {
				for _, event := range watchChan.Events {
					if event.Type == clientv3.EventTypePut {
						// 根据 key 后缀判断配置类型
						parts := strings.Split(key, ".")
						if len(parts) == 2 {
							newViper.SetConfigType(parts[1])
						} else {
							newViper.SetConfigType("toml")
						}
						if err := newViper.ReadConfig(bytes.NewReader(event.Kv.Value)); err != nil {
							alog.Error(context.Background(), "Viper ReadConfig error", zap.Error(err))
							continue
						}
						filename := strings.TrimSuffix(filepath.Base(key), filepath.Ext(key))
						for k, v := range newViper.AllSettings() {
							if index == 0 {
								c.Viper.Set(k, v)
							}
							c.Viper.Set(fmt.Sprintf("%s.%s", filename, k), v)
						}
					}
				}
			}
		}(index, key)
	}
	// 阻塞防止退出，如果需要退出逻辑，请自行控制 context 取消
	select {}
}

// AddRemoteProvider 添加远程配置提供者
// AddRemoteProvider adds a remote configuration provider (without security).
func (c *ConfigStr) AddRemoteProvider(provider, endpoint, path string) error {
	c.Viper.SetConfigType("toml")
	if err := c.Viper.AddRemoteProvider(provider, endpoint, path); err != nil {
		return err
	}
	return c.Viper.ReadRemoteConfig()
}

// AddSecureRemoteProvider 添加带安全认证的远程配置提供者
// AddSecureRemoteProvider adds a remote configuration provider with security.
func (c *ConfigStr) AddSecureRemoteProvider(provider, endpoint, path, secretKeyRing string) error {
	if err := c.Viper.AddSecureRemoteProvider(provider, endpoint, path, secretKeyRing); err != nil {
		return err
	}
	return c.Viper.ReadRemoteConfig()
}

// AddPath 增加配置文件搜索路径
// AddPath adds a configuration search path.
func (c *ConfigStr) AddPath(in string) *ConfigStr {
	c.Viper.AddConfigPath(in)
	return c
}

// SetFile 设置配置文件路径
// SetFile sets the configuration file path.
func (c *ConfigStr) SetFile(in string) *ConfigStr {
	c.Viper.SetConfigFile(in)
	return c
}

// SetType 设置配置文件类型
// SetType sets the configuration file type.
func (c *ConfigStr) SetType(in string) *ConfigStr {
	c.Viper.SetConfigType(in)
	return c
}

// Register 读取配置（本地或远程）
// Register reads the configuration from file or remote provider.
func (c *ConfigStr) Register() error {
	if len(c.filePath) == 1 {
		return c.Viper.ReadInConfig()
	}
	if err := c.Viper.ReadRemoteConfig(); err != nil {
		return err
	}
	return c.Viper.WatchRemoteConfigOnChannel()
}

// 以下是全局封装的辅助函数，便于在项目中直接获取配置值

// SetKey 设置配置键值对
// SetKey sets a configuration key-value pair.
func SetKey(key string, value any) {
	Config.Viper.Set(key, value)
}

// Get 获取配置值
// Get retrieves the configuration value for the given key.
func Get(key string) any {
	if Config == nil {
		return nil
	}
	return Config.Viper.Get(key)
}

// GetByte 获取配置对应的字节数组
// GetByte retrieves the configuration value as a byte slice.
func GetByte(key string) []byte {
	if Config == nil {
		return []byte{}
	}
	return conv.Bytes(Config.Viper.GetString(key))
}

// GetString 获取字符串配置
// GetString retrieves the configuration value as a string.
func GetString(key string) string {
	if Config == nil {
		return ""
	}
	return Config.Viper.GetString(key)
}

// GetBool 获取布尔类型配置
// GetBool retrieves the configuration value as a bool.
func GetBool(key string) bool {
	if Config == nil {
		return false
	}
	return Config.Viper.GetBool(key)
}

// GetInt 获取整型配置
// GetInt retrieves the configuration value as an int.
func GetInt(key string) int {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetInt(key)
}

// GetInt64 获取64位整型配置
// GetInt64 retrieves the configuration value as an int64.
func GetInt64(key string) int64 {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetInt64(key)
}

// GetStringMap 获取字符串映射配置
// GetStringMap retrieves the configuration value as a map[string]any.
func GetStringMap(key string) map[string]any {
	if Config == nil {
		return map[string]any{}
	}
	return Config.Viper.GetStringMap(key)
}

// GetStringMapStringSlice 获取字符串映射切片配置
// GetStringMapStringSlice retrieves the configuration value as a map[string][]string.
func GetStringMapStringSlice(key string) map[string][]string {
	if Config == nil {
		return map[string][]string{}
	}
	return Config.Viper.GetStringMapStringSlice(key)
}

// GetDuration 获取持续时间配置
// GetDuration retrieves the configuration value as a time.Duration.
func GetDuration(key string) time.Duration {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetDuration(key)
}

// GetFloat64 获取浮点型配置
// GetFloat64 retrieves the configuration value as a float64.
func GetFloat64(key string) float64 {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetFloat64(key)
}

// GetTime 获取时间配置
// GetTime retrieves the configuration value as a time.Time.
func GetTime(key string) time.Time {
	if Config == nil {
		return time.Time{}
	}
	return Config.Viper.GetTime(key)
}

// GetStringSlice 获取字符串切片配置
// GetStringSlice retrieves the configuration value as a slice of strings.
func GetStringSlice(key string) []string {
	if Config == nil {
		return []string{}
	}
	return Config.Viper.GetStringSlice(key)
}

// GetIntSlice 获取整型切片配置
// GetIntSlice retrieves the configuration value as a slice of ints.
func GetIntSlice(key string) []int {
	if Config == nil {
		return []int{}
	}
	return Config.Viper.GetIntSlice(key)
}

// GetMaps 获取映射数组配置
// GetMaps retrieves the configuration value as a slice of maps.
func GetMaps(key string) []map[string]any {
	if Config == nil {
		return []map[string]any{}
	}
	return conv.Maps(Config.Viper.Get(key))
}
