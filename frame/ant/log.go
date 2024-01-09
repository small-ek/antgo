package ant

import (
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
)

// Log Get Log content
func Log() *zap.Logger {
	return alog.Write
}
