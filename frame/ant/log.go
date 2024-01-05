package ant

import (
	"go.uber.org/zap"
)

// Log Get Log content
func Log() *zap.Logger {
	return alog.Write
}
