package telebot

import (
	"fmt"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (t Telebot) SendCurrentPrice(msg string, chatID int64) error {
	tgMsg := tgbot.NewMessage(chatID, msg)

	maxRetries := 10
	cnt := 0

	for cnt < maxRetries {
		if _, err := t.bot.Send(tgMsg); err != nil {
			logrus.Errorf("Cannot send a message: %v", err)

			time.Sleep(5 * time.Second)
			cnt++

			continue
		}

		return nil
	}

	return fmt.Errorf("cannot send the message")
}
