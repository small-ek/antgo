package config

import (
	"github.com/spf13/viper"
	"sync"
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
		if len(filePath) > 0 && filePath[0] != "" {
			Config.Viper.AddConfigPath(filePath[0])
		}
		if len(filePath) > 1 && filePath[1] != "" {
			Config.Viper.SetConfigFile(filePath[1])
		}
	})
	return Config
}

// AddRemoteProvider 添加远程连接
func (c *ConfigStr) AddRemoteProvider(provider, endpoint, path string) error {
	return c.Viper.AddRemoteProvider(provider, endpoint, path)
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
	return c.Viper.ReadInConfig()
}

// SetKey 设置值
func SetKey(key string, value any) {
	Config.Viper.Set(key, value)
}

// Get 获取
func Get(key string) any {
	return Config.Viper.Get(key)
}
