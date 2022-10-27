package mail

import (
	"crypto/tls"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	UserPW string `json:"userpw"`
}

type Mail struct {
	To         string // 接收邮箱
	Subject    string // 主题
	Body       string // 内容
	AttachPath string // 附件路径
}

// 发送单封邮件
func SendEmail(smtp *Smtp, mail *Mail) error {
	// 收件人不能为空
	if mail.To == "" {
		return fmt.Errorf("%#v can not empty", mail.To)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", smtp.User)
	m.SetHeader("To", mail.To)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	if mail.AttachPath != "" {
		m.Attach(mail.AttachPath)
	}

	e := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.UserPW)
	e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := e.DialAndSend(m); err != nil {
		// 失败暂停 1s 重发
		time.Sleep(1 * time.Second)
		if err := e.DialAndSend(m); err != nil {
			return err
		}
	}
	return nil
}

// 发送单封邮件给多人
func SendEmailMP(smtp *Smtp, mail *Mail, mailList []string) []error {
	var errList []error
	for _, to := range mailList {
		mail.To = to
		if err := SendEmail(smtp, mail); err != nil {
			errList = append(errList, err)
		}

		// 间隔 0.5 秒
		time.Sleep(500 * time.Millisecond)
	}
	return errList
}
