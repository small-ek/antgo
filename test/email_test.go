package test

import (
	"fmt"
	"github.com/small-ek/antgo/aemail"
	"testing"
)

func TestEmail(t *testing.T) {
	mailer := aemail.NewMailer("56494565@qq.com", "fdtshicbbvybbiic")
	err := mailer.QuickSend(
		[]string{"56494565@qq.com"},
		"测试邮件",
		"这是一封测试邮件",
	)
	fmt.Println(err)
}
