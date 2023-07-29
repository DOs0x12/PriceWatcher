package sender

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

const (
	from     = "guise322@ya.ru"
	pass     = "nxwamiqmoqdolhds"
	to       = "paulina.urgn@yandex.ru"
	smtpHost = "smtp.yandex.ru"
	smtpPort = 465
)

type Sender struct{}

func (s Sender) Send(price float32) error {
	msg := fmt.Sprintf("Курс золота. Продажа: %.2fр", price)
	sub := "Че по золоту?"
	m := gomail.NewMessage()
	configureMsg(m, sub, msg)
	d := gomail.NewDialer(smtpHost, smtpPort, from, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		sendingErr := "an error occurs when sending an email: %w"

		return fmt.Errorf(sendingErr, err)
	}

	return nil
}

func configureMsg(m *gomail.Message, sub, msg string) {
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", sub)
	m.SetBody("text/plain", msg)
}
