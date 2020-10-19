package test

import (
	"github.com/small-ek/ginp/email"
	"testing"
)

func TestEmail(t *testing.T) {
	email.Config.SetFrom("56494565@qq.com").SetTo([]string{"56494565@qq.com"}).SetTitle("test").SetText("test2223223").SetPassword("fdtshicbbvybbiic").SetFilePath([]string{"test.txt"})
	email.Send()

}
