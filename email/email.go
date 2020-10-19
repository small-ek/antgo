package email

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

type Data struct {
	From     string
	To       []string
	Cc       []string
	Bcc      []string
	Title    string
	Text     string
	Html     string
	Password string
	Address  string
	Host     string
	FilePath []string
}

var Config Data

func (this *Data) SetFrom(from string) *Data {
	this.From = from
	return this
}

func (this *Data) SetTo(to []string) *Data {
	this.To = to
	return this
}

func (this *Data) SetTitle(title string) *Data {
	this.Title = title
	return this
}

func (this *Data) SetText(text string) *Data {
	this.Text = text
	return this
}
func (this *Data) SetHtml(html string) *Data {
	this.Html = html
	return this
}
func (this *Data) SetPassword(password string) *Data {
	this.Password = password
	return this
}

func (this *Data) SetAddress(address string) *Data {
	this.Address = address
	return this
}

func (this *Data) SetHost(host string) *Data {
	this.Host = host
	return this
}
func (this *Data) SetFilePath(file_path []string) *Data {
	this.FilePath = file_path
	return this
}
func Send() {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = Config.From
	// 设置接收方的邮箱
	e.To = Config.To
	//设置主题
	e.Subject = Config.Title
	//设置文件发送的内容
	if Config.Text != "" {
		e.Text = []byte(Config.Text)
	}
	//设置文件发送的html
	if Config.Html != "" {
		e.HTML = []byte(Config.Html)
	}
	//附件
	if len(Config.FilePath) > 0 {
		for i := 0; i < len(Config.FilePath); i++ {
			e.AttachFile(Config.FilePath[i])
		}
	}
	//设置服务器相关的配置
	if Config.Address == "" {
		Config.Address = "smtp.qq.com:25"
	}
	//发送地址
	if Config.Host == "" {
		Config.Host = "smtp.qq.com"
	}
	//设置抄送如果抄送多人逗号隔开
	if len(Config.Cc) > 0 {
		e.Cc = Config.Cc
	}
	//设置秘密抄送
	if len(Config.Bcc) > 0 {
		e.Bcc = Config.Bcc
	}

	err := e.Send(Config.Address, smtp.PlainAuth("", Config.From, Config.Password, Config.Host))
	if err != nil {
		log.Fatal(err)
	}
}
