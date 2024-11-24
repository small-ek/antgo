package aerror

import (
	"github.com/pkg/errors"
	"runtime"
	"strconv"
)

// Wrap 堆栈跟踪
func Wrap(message string, err error) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

// printCallerNameAndLine
func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}
