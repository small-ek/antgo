package ant

import (
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/os/logger"
	"go.uber.org/zap"
)

// Log Get Log content
func Log() *zap.Logger {
	if config.Decode().Get("log.switch").Bool() == false {
		panic("Log profile is not enabled")
	}
	return logger.Write
}
