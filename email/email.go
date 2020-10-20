package email

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

//New Email parameter structure
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
func (this *New) Send() {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = this.From
	// 设置接收方的邮箱
	e.To = this.To
	//设置主题
	e.Subject = this.Title
	//设置文件发送的内容
	if this.Text != "" {
		e.Text = []byte(this.Text)
	}
	//设置文件发送的html
	if this.Html != "" {
		e.HTML = []byte(this.Html)
	}
	//附件
	if len(this.FilePath) > 0 {
		for i := 0; i < len(this.FilePath); i++ {
			e.AttachFile(this.FilePath[i])
		}
	}
	//设置服务器相关的配置
	if this.Address == "" {
		this.Address = "smtp.qq.com:25"
	}
	//发送地址
	if this.Host == "" {
		this.Host = "smtp.qq.com"
	}
	//设置抄送如果抄送多人逗号隔开
	if len(this.Cc) > 0 {
		e.Cc = this.Cc
	}
	//设置秘密抄送
	if len(this.Bcc) > 0 {
		e.Bcc = this.Bcc
	}

	err := e.Send(this.Address, smtp.PlainAuth("", this.From, this.Password, this.Host))
	if err != nil {
		log.Fatal(err)
	}
}
