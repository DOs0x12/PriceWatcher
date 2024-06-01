package price

import (
	"PriceWatcher/internal/entities/subscribing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SubscribingComm struct {
	Subscribers subscribing.Subscribers
}

func (c SubscribingComm) SubscribeUser(input interface{}) string {
	upd := input.(tgbotapi.Update)
	c.Subscribers.ChatIDs = append(c.Subscribers.ChatIDs, upd.Message.Chat.ID)

	return "Subscribed!"
}
