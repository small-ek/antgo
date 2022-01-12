package ant

import (
	"github.com/small-ek/antgo/os/logger"
	"go.uber.org/zap"
)

// Log Get Log content
func Log() *zap.Logger {
	return logger.Write
}
