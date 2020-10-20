package logs

import (
	"log"
	"runtime/debug"
)

//Error Print error stack
func Error(err string) {
	log.Println(err)
	debug.PrintStack()
}
