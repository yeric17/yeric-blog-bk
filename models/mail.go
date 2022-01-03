package models

import (
	"fmt"
	"yeric-blog/config"

	gomail "gopkg.in/mail.v2"
)

func SendMail(subject string, body string, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Mail.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	fmt.Printf("Sending mail to %s\n", to)

	d := gomail.NewDialer(config.Mail.Host, config.Mail.Port, config.Mail.Username, config.Mail.Password)

	return d.DialAndSend(m)
}
