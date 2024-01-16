package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
func New(filePath ...string) *ConfigStr {
	once.Do(func() {
		Config = &ConfigStr{Viper: viper.New(), filePath: filePath}

		if len(filePath) == 1 {
			Config.Viper.SetConfigFile(filePath[0])
			types := strings.Split(filePath[0], ".")

			if len(types) == 2 {
				Config.Viper.SetConfigType(types[1])
			}
		}

		if len(filePath) == 3 {
			if err := Config.Viper.AddRemoteProvider(filePath[0], filePath[1], filePath[2]); err != nil {
				panic(err)
			}
			types := strings.Split(filePath[2], ".")

			if len(types) == 2 {
				Config.Viper.SetConfigType(types[1])
			}
		}

		//加密链接
		if len(filePath) == 4 {
			if err := Config.Viper.AddSecureRemoteProvider(filePath[0], filePath[1], filePath[2], filePath[3]); err != nil {
				panic(err)
			}
			types := strings.Split(filePath[2], ".")

			if len(types) == 2 {
				Config.Viper.SetConfigType(types[1])
			}
		}
	})

	return Config
}

// Etcd3
func (c *ConfigStr) Etcd3(hosts []string, path, username, pwd string) (err error) {
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
	if path != "" {
		types := strings.Split(path, ".")
		Config.Viper.SetConfigType(types[1])
	}

	if err = c.loadEtcdFormToViper(path, cli); err != nil {
		return err
	}
	go c.watchEtcd3(path, cli)
	//if err = c.Viper.ReadRemoteConfig(); err != nil {
	//	panic(err)
	//}
	return nil
}

// loadEtcdFormToViper
func (c *ConfigStr) loadEtcdFormToViper(prefix string, client *clientv3.Client) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	err = c.Viper.ReadConfig(bytes.NewReader(resp.Kvs[0].Value))
	if err != nil {
		return err
	}

	return nil
}

// watchEtcd3 监听etcd
func (c *ConfigStr) watchEtcd3(path string, cli *clientv3.Client) {
	// 创建一个context
	watcher := clientv3.NewWatcher(cli)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			watchChan := watcher.Watch(ctx, path)

			for resp := range watchChan {
				for _, event := range resp.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						err := c.Viper.ReadConfig(bytes.NewReader(event.Kv.Value))
						if err != nil {
							alog.Error("watchEtcd", zap.Error(err))
						}
						if err = c.Viper.ReadRemoteConfig(); err != nil {
							alog.Error("watchEtcd", zap.Error(err))
						}
					case clientv3.EventTypeDelete:
						fmt.Printf("Key %s deleted.\n", event.Kv.Key)
					}
				}
			}
		}
	}()

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
func (c *ConfigStr) AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	err := c.Viper.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
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

// Regiter 注册读取配置
func (c *ConfigStr) Regiter() (err error) {
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
