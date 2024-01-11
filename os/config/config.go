package config

import (
	"github.com/small-ek/antgo/utils/conv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"sync"
	"time"
)

var Config *ConfigStr
var once sync.Once

type ConfigStr struct {
	Viper *viper.Viper
}

// New<初始化配置>
func New(filePath ...string) *ConfigStr {
	once.Do(func() {
		Config = &ConfigStr{Viper: viper.New()}
		if len(filePath) == 1 {
			Config.Viper.SetConfigFile(filePath[0])
			types := strings.Split(filePath[0], ".")

			if len(types) == 2 {
				Config.Viper.SetConfigType(types[1])
			}
		}
		if len(filePath) == 3 {
			Config.Viper.AddRemoteProvider(filePath[0], filePath[1], filePath[2])
			types := strings.Split(filePath[2], ".")

			if len(types) == 2 {
				Config.Viper.SetConfigType(types[1])
			}
		}
	})

	return Config
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
func (c *ConfigStr) Regiter() error {

	return c.Viper.ReadRemoteConfig()
}

// SetKey 设置值
func SetKey(key string, value any) {
	Config.Viper.Set(key, value)
}

// Get 获取
func Get(key string) any {
	return Config.Viper.Get(key)
}

// GetString 获取
func GetString(key string) string {
	return Config.Viper.GetString(key)
}

// GetBool 获取
func GetBool(key string) bool {
	return Config.Viper.GetBool(key)
}

// GetInt 获取
func GetInt(key string) int {
	return Config.Viper.GetInt(key)
}

// GetInt64 获取
func GetInt64(key string) int64 {
	return Config.Viper.GetInt64(key)
}

// GetStringMap 获取
func GetStringMap(key string) map[string]any {
	return Config.Viper.GetStringMap(key)
}

// GetStringMapStringSlice 获取
func GetStringMapStringSlice(key string) map[string][]string {
	return Config.Viper.GetStringMapStringSlice(key)
}

// GetDuration 获取
func GetDuration(key string) time.Duration {
	return Config.Viper.GetDuration(key)
}

// GetFloat64 获取
func GetFloat64(key string) float64 {
	return Config.Viper.GetFloat64(key)
}

// GetTime 获取
func GetTime(key string) time.Time {
	return Config.Viper.GetTime(key)
}

// GetStringSlice 获取
func GetStringSlice(key string) []string {
	return Config.Viper.GetStringSlice(key)
}

// GetIntSlice 获取
func GetIntSlice(key string) []int {
	return Config.Viper.GetIntSlice(key)
}

// GetMaps 获取
func GetMaps(key string) []map[string]any {
	return conv.Maps(Config.Viper.Get(key))
}
