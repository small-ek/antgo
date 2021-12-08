package ant

import (
	"github.com/small-ek/antgo/os/config"
)

// GetConfig Get configuration content
func GetConfig(name string) *config.Config {
	return config.Decode().Get(name)
}
