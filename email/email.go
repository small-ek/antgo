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
func (get *New) SetFrom(from string) *New {
	get.From = from
	return get
}

//SetTo Set To
func (get *New) SetTo(to []string) *New {
	get.To = to
	return get
}

//SetTitle Set Title
func (get *New) SetTitle(title string) *New {
	get.Title = title
	return get
}

//SetText Set Text
func (get *New) SetText(text string) *New {
	get.Text = text
	return get
}

//SetHtml Set Html
func (get *New) SetHtml(html string) *New {
	get.Html = html
	return get
}

//SetPassword Set Password
func (get *New) SetPassword(password string) *New {
	get.Password = password
	return get
}

//SetAddress Set Address
func (get *New) SetAddress(address string) *New {
	get.Address = address
	return get
}

//SetHost Set Host
func (get *New) SetHost(host string) *New {
	get.Host = host
	return get
}

//SetFilePath Set Email attachment path
func (get *New) SetFilePath(filePath []string) *New {
	get.FilePath = filePath
	return get
}

//Config ...
var Config New

//Send Email
func (get *New) Send() {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = get.From
	// 设置接收方的邮箱
	e.To = get.To
	//设置主题
	e.Subject = get.Title
	//设置文件发送的内容
	if get.Text != "" {
		e.Text = []byte(get.Text)
	}
	//设置文件发送的html
	if get.Html != "" {
		e.HTML = []byte(get.Html)
	}
	//附件
	if len(get.FilePath) > 0 {
		for i := 0; i < len(get.FilePath); i++ {
			e.AttachFile(get.FilePath[i])
		}
	}
	//设置服务器相关的配置
	if get.Address == "" {
		get.Address = "smtp.qq.com:25"
	}
	//发送地址
	if get.Host == "" {
		get.Host = "smtp.qq.com"
	}
	//设置抄送如果抄送多人逗号隔开
	if len(get.Cc) > 0 {
		e.Cc = get.Cc
	}
	//设置秘密抄送
	if len(get.Bcc) > 0 {
		e.Bcc = get.Bcc
	}

	err := e.Send(get.Address, smtp.PlainAuth("", get.From, get.Password, get.Host))
	if err != nil {
		log.Fatal(err)
	}
}
