package aemail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

type Mailer struct {
	from     string
	password string
	server   string
	host     string
	port     int
	useTLS   bool
}

type Email struct {
	To      []string
	Subject string
	Text    string
	HTML    string
	Files   []string
}

// NewMailer 创建邮件客户端（推荐方式）
// 自动识别常见邮箱服务的 SMTP 配置
func NewMailer(from, password string) *Mailer {
	m := &Mailer{
		from:     from,
		password: password,
		useTLS:   true, // 默认启用 TLS
	}

	domain := strings.Split(from, "@")[1]
	switch {
	case strings.Contains(domain, "qq.com"):
		m.server = "smtp.qq.com:465"
		m.host = "smtp.qq.com"
		m.port = 465
	case strings.Contains(domain, "gmail.com"):
		m.server = "smtp.gmail.com:587"
		m.host = "smtp.gmail.com"
		m.port = 587
	case strings.Contains(domain, "163.com"):
		m.server = "smtp.163.com:465"
		m.host = "smtp.163.com"
		m.port = 465
	default:
		m.server = "smtp." + domain + ":465"
		m.host = "smtp." + domain
		m.port = 465
	}

	return m
}

// WithCustomSMTP 自定义 SMTP 配置
func (m *Mailer) WithCustomSMTP(host string, port int, useTLS bool) *Mailer {
	m.server = fmt.Sprintf("%s:%d", host, port)
	m.host = host
	m.port = port
	m.useTLS = useTLS
	return m
}

// Send 发送邮件（自动选择 TLS 配置）
func (m *Mailer) Send(e *Email) error {
	em := email.NewEmail()
	em.From = m.from
	em.To = e.To
	em.Subject = e.Subject

	if e.Text != "" {
		em.Text = []byte(e.Text)
	}
	if e.HTML != "" {
		em.HTML = []byte(e.HTML)
	}

	for _, f := range e.Files {
		if _, err := em.AttachFile(f); err != nil {
			return fmt.Errorf("附件添加失败: %w", err)
		}
	}

	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	if m.useTLS {
		return em.SendWithTLS(
			m.server,
			auth,
			&tls.Config{ServerName: m.host},
		)
	}
	return em.Send(m.server, auth)
}

// QuickSend 快速发送文本邮件
func (m *Mailer) QuickSend(to []string, subject, text string) error {
	return m.Send(&Email{
		To:      to,
		Subject: subject,
		Text:    text,
	})
}

//================= 高级用法 =================//

// EmailOption 邮件配置函数类型
type EmailOption func(*Email)

// NewEmail 创建邮件对象
func NewEmail(opts ...EmailOption) *Email {
	e := &Email{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func WithTo(to ...string) EmailOption {
	return func(e *Email) {
		e.To = to
	}
}

func WithSubject(subject string) EmailOption {
	return func(e *Email) {
		e.Subject = subject
	}
}

func WithText(text string) EmailOption {
	return func(e *Email) {
		e.Text = text
	}
}

func WithHTML(html string) EmailOption {
	return func(e *Email) {
		e.HTML = html
	}
}

func WithFiles(files ...string) EmailOption {
	return func(e *Email) {
		e.Files = files
	}
}
