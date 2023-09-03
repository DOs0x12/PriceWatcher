package sender

import (
	"PriceWatcher/internal/entities/config"
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type Sender struct{}

func (s Sender) Send(message, subject string, conf config.Email) error {
	msg := message
	sub := subject
	m := gomail.NewMessage()
	configureMsg(m, sub, msg, conf)
	d := gomail.NewDialer(conf.SmtpHost, conf.SmtpPort, conf.From, conf.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("an error occurs when sending an email: %w", err)
	}

	return nil
}

func configureMsg(m *gomail.Message, sub, msg string, conf config.Email) {
	m.SetHeader("From", conf.From)
	m.SetHeader("To", conf.To)
	m.SetHeader("Subject", sub)
	m.SetBody("text/plain", msg)
}
