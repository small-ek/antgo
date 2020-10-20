package email

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

//Email parameter structure
type New struct {
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

var Config New

//SetFrom Set Send email
func (this *New) SetFrom(from string) *New {
	this.From = from
	return this
}

//SetTo Set To
func (this *New) SetTo(to []string) *New {
	this.To = to
	return this
}

//SetTitle Set Title
func (this *New) SetTitle(title string) *New {
	this.Title = title
	return this
}

//SetText Set Text
func (this *New) SetText(text string) *New {
	this.Text = text
	return this
}

//SetHtml Set Html
func (this *New) SetHtml(html string) *New {
	this.Html = html
	return this
}

//SetPassword Set Password
func (this *New) SetPassword(password string) *New {
	this.Password = password
	return this
}

//SetAddress Set Address
func (this *New) SetAddress(address string) *New {
	this.Address = address
	return this
}

//SetHost Set Host
func (this *New) SetHost(host string) *New {
	this.Host = host
	return this
}

//SetFilePath Set Email attachment path
func (this *New) SetFilePath(file_path []string) *New {
	this.FilePath = file_path
	return this
}

//Send Email
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
