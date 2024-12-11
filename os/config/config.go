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

var Config *ConfigStr
var once sync.Once

type ConfigStr struct {
	Viper    *viper.Viper
	filePath []string
	UserName string
	Password string
}

// New<初始化配置>
func New(path ...string) *ConfigStr {
	once.Do(func() {
		Config = &ConfigStr{Viper: viper.New(), filePath: path}

		if len(path) > 0 {
			Config.setupConfigFile(path)
		}
	})

	return Config
}

// AddConfigFile<加载一个新的配置并且合并>
func AddConfigFile(path string) error {
	newViper := viper.New()
	newViper.SetConfigType(filepath.Ext(path)[1:])
	newViper.SetConfigFile(path)
	if err := newViper.ReadInConfig(); err != nil {
		return err
	}
	newViper.WatchConfig()
	//合并配置
	Config.Viper.MergeConfigMap(newViper.AllSettings())

	// 当文件发生修改时触发回调
	newViper.OnConfigChange(func(e fsnotify.Event) {
		if err := newViper.ReadInConfig(); err == nil {
			// 合并最新的配置到主配置
			Config.Viper.MergeConfigMap(newViper.AllSettings())
		} else {
			alog.Error("Viper ReadInConfig error", zap.Error(err))
		}
	})
	return nil
}

func (c *ConfigStr) setupConfigFile(path []string) {
	if len(path) == 1 {
		c.Viper.SetConfigFile(path[0])
		c.Viper.WatchConfig()
	}
	if len(path) >= 2 {
		c.Viper.SetConfigType(filepath.Ext(path[len(path)-1])[1:])
	}

	if len(path) == 3 {
		if err := c.Viper.AddRemoteProvider(path[0], path[1], path[2]); err != nil {
			panic(err)
		}
	}
	if len(path) == 4 {
		if err := c.Viper.AddSecureRemoteProvider(path[0], path[1], path[2], path[3]); err != nil {
			panic(err)
		}
	}
}

// Etcd3 ETCD3 configuration link
func (c *ConfigStr) Etcd3(hosts, path []string, username, pwd string) (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: hosts,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
		},
		Username: username,
		Password: pwd,
		Logger:   alog.Write,
	})

	if err != nil {
		return
	}

	if err = c.loadEtcdFormToViper(path, cli); err != nil {
		return err
	}
	go c.watchEtcd3(path, cli)
	return nil
}

// loadEtcdFormToViper
func (c *ConfigStr) loadEtcdFormToViper(path []string, client *clientv3.Client) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < len(path); i++ {
		newViper := viper.New()
		types := strings.Split(path[i], ".")
		if len(types) == 2 {
			newViper.SetConfigType(types[1])
		} else {
			newViper.SetConfigType("toml")
		}

		resp, err := client.Get(ctx, path[i], clientv3.WithPrefix())
		if err != nil {
			return err
		}

		err = newViper.ReadConfig(bytes.NewReader(resp.Kvs[0].Value))
		if err != nil {
			return err
		}

		filenameWithExt := filepath.Base(path[i])
		filename := strings.TrimSuffix(filenameWithExt, filepath.Ext(filenameWithExt))
		for key, value := range newViper.AllSettings() {
			if i == 0 {
				c.Viper.Set(fmt.Sprintf("%s", key), value)
			}
			c.Viper.Set(fmt.Sprintf("%s.%s", filename, key), value)
		}

	}

	return nil
}

// watchEtcd3 监听etcd
func (c *ConfigStr) watchEtcd3(path []string, cli *clientv3.Client) {
	for i := 0; i < len(path); i++ {
		newViper := viper.New()
		go func(p []string) {

			watcher := clientv3.NewWatcher(cli)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			for {
				watchChan := watcher.Watch(ctx, p[i])
				for resp := range watchChan {
					for _, event := range resp.Events {
						switch event.Type {
						case clientv3.EventTypePut:
							types := strings.Split(p[i], ".")
							if len(types) == 2 {
								newViper.SetConfigType(types[1])
							} else {
								newViper.SetConfigType("toml")
							}
							err := newViper.ReadConfig(bytes.NewReader(event.Kv.Value))
							if err != nil {
								alog.Error("Viper ReadConfig error", zap.Error(err))
							}

							filenameWithExt := filepath.Base(p[i])
							filename := strings.TrimSuffix(filenameWithExt, filepath.Ext(filenameWithExt))
							for key, value := range newViper.AllSettings() {
								if i == 0 {
									c.Viper.Set(fmt.Sprintf("%s", key), value)
								}
								c.Viper.Set(fmt.Sprintf("%s.%s", filename, key), value)
							}
						}
					}
				}
			}

		}(path)
	}
	select {}
}

// AddRemoteProvider 添加远程连接
func (c *ConfigStr) AddRemoteProvider(provider, endpoint, path string) error {
	c.Viper.SetConfigType("toml")
	err := c.Viper.AddRemoteProvider(provider, endpoint, path)

	if err != nil {
		return err
	}
	return c.Viper.ReadRemoteConfig()
}

// AddSecureRemoteProvider 添加远程连接
func (c *ConfigStr) AddSecureRemoteProvider(provider, endpoint, path, secretKeyRing string) error {
	err := c.Viper.AddSecureRemoteProvider(provider, endpoint, path, secretKeyRing)
	if err != nil {
		return err
	}
	return c.Viper.ReadRemoteConfig()
}

// AddPath 增加配置文件路径
func (c *ConfigStr) AddPath(in string) *ConfigStr {
	c.Viper.AddConfigPath(in)
	return c
}

// SetFile 设置目录
func (c *ConfigStr) SetFile(in string) *ConfigStr {
	c.Viper.SetConfigFile(in)
	return c
}

// SetType 设置文件类型
func (c *ConfigStr) SetType(in string) *ConfigStr {
	c.Viper.SetConfigType(in)
	return c
}

// Register 注册读取配置
func (c *ConfigStr) Register() (err error) {
	if len(c.filePath) == 1 {
		return c.Viper.ReadInConfig()
	}
	if err = c.Viper.ReadRemoteConfig(); err != nil {
		return
	}

	if err = c.Viper.WatchRemoteConfigOnChannel(); err != nil {
		return
	}

	return
}

// SetKey 设置值
func SetKey(key string, value any) {
	Config.Viper.Set(key, value)
}

// Get 获取
func Get(key string) any {
	if Config == nil {
		return nil
	}
	return Config.Viper.Get(key)
}

// GetByte 获取
func GetByte(key string) []byte {
	if Config == nil {
		return []byte{}
	}
	return conv.Bytes(Config.Viper.GetString(key))
}

// GetString 获取
func GetString(key string) string {
	if Config == nil {
		return ""
	}
	return Config.Viper.GetString(key)
}

// GetBool 获取
func GetBool(key string) bool {
	if Config == nil {
		return false
	}
	return Config.Viper.GetBool(key)
}

// GetInt 获取
func GetInt(key string) int {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetInt(key)
}

// GetInt64 获取
func GetInt64(key string) int64 {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetInt64(key)
}

// GetStringMap 获取
func GetStringMap(key string) map[string]any {
	if Config == nil {
		return map[string]any{}
	}
	return Config.Viper.GetStringMap(key)
}

// GetStringMapStringSlice 获取
func GetStringMapStringSlice(key string) map[string][]string {
	if Config == nil {
		return map[string][]string{}
	}
	return Config.Viper.GetStringMapStringSlice(key)
}

// GetDuration 获取
func GetDuration(key string) time.Duration {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetDuration(key)
}

// GetFloat64 获取
func GetFloat64(key string) float64 {
	if Config == nil {
		return 0
	}
	return Config.Viper.GetFloat64(key)
}

// GetTime 获取
func GetTime(key string) time.Time {
	if Config == nil {
		return time.Time{}
	}
	return Config.Viper.GetTime(key)
}

// GetStringSlice 获取
func GetStringSlice(key string) []string {
	if Config == nil {
		return []string{}
	}
	return Config.Viper.GetStringSlice(key)
}

// GetIntSlice 获取
func GetIntSlice(key string) []int {
	if Config == nil {
		return []int{}
	}
	return Config.Viper.GetIntSlice(key)
}

// GetMaps 获取
func GetMaps(key string) []map[string]any {
	if Config == nil {
		return []map[string]any{}
	}
	return conv.Maps(Config.Viper.Get(key))
}
