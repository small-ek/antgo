package aemail

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

//Email Email parameter structure
type Email struct {
	From     string   //Send email
	To       []string //Accept mailbox
	Cc       []string //Set cc
	Bcc      []string //Set Bcc
	Title    string   //Email title
	Text     string   //Email Text
	Html     string   //Email Html
	Password string   //Email password
	Address  string   //Send email address
	Host     string   //Send email host
	FilePath []string //Email attachment path
}

//SetFrom Set Send email
func (e *Email) SetFrom(from string) *Email {
	e.From = from
	return e
}

//SetTo Set To
func (e *Email) SetTo(to []string) *Email {
	e.To = to
	return e
}

//SetTitle Set Title
func (e *Email) SetTitle(title string) *Email {
	e.Title = title
	return e
}

//SetText Set Text
func (e *Email) SetText(text string) *Email {
	e.Text = text
	return e
}

//SetHtml Set Html
func (e *Email) SetHtml(html string) *Email {
	e.Html = html
	return e
}

//SetPassword Set Password
func (e *Email) SetPassword(password string) *Email {
	e.Password = password
	return e
}

//SetAddress Set Address
func (e *Email) SetAddress(address string) *Email {
	e.Address = address
	return e
}

//SetHost Set Host
func (e *Email) SetHost(host string) *Email {
	e.Host = host
	return e
}

//SetFilePath Set Email attachment path
func (e *Email) SetFilePath(filePath []string) *Email {
	e.FilePath = filePath
	return e
}

//New 创建
func New(from ...string) *Email {
	if len(from[0]) > 0 {
		return &Email{From: from[0]}
	}
	return &Email{}
}

//Send Email
func (e *Email) Send() {
	emails := email.NewEmail()
	//设置发送方的邮箱
	emails.From = e.From
	// 设置接收方的邮箱
	emails.To = e.To
	//设置主题
	emails.Subject = e.Title
	//设置文件发送的内容
	if e.Text != "" {
		emails.Text = []byte(e.Text)
	}
	//设置文件发送的html
	if e.Html != "" {
		emails.HTML = []byte(e.Html)
	}
	//附件
	if len(e.FilePath) > 0 {
		for i := 0; i < len(e.FilePath); i++ {
			_, err := emails.AttachFile(e.FilePath[i])
			if err != nil {
				panic(err)
			}
		}
	}
	//设置服务器相关的配置
	if e.Address == "" {
		e.Address = "smtp.qq.com:25"
	}
	//发送地址
	if e.Host == "" {
		e.Host = "smtp.qq.com"
	}
	//设置抄送如果抄送多人逗号隔开
	if len(e.Cc) > 0 {
		emails.Cc = e.Cc
	}
	//设置秘密抄送
	if len(e.Bcc) > 0 {
		emails.Bcc = e.Bcc
	}

	err := emails.Send(e.Address, smtp.PlainAuth("", e.From, e.Password, e.Host))
	if err != nil {
		log.Fatal(err)
	}
}
