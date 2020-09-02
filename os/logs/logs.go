package logs

import (
	"log"
	"runtime/debug"
)

//打印错误堆栈
func Error(err string) {
	log.Println(err)
	debug.PrintStack()
}
