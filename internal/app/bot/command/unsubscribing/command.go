package unsubscribing

import (
	"PriceWatcher/internal/entities/subscribing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
)

type UnsubscribingComm struct {
	Subscribers *subscribing.Subscribers
}

func (c UnsubscribingComm) UnsubscribeUser(input interface{}) string {
	errMessage := "Error: the user is not subscribed!"

	if len(c.Subscribers.ChatIDs) == 0 {
		return errMessage
	}

	upd := input.(tgbotapi.Update)
	idIndex := slices.Index(c.Subscribers.ChatIDs, upd.Message.Chat.ID)
	if idIndex == -1 {
		return errMessage
	}

	c.Subscribers.ChatIDs = slices.Delete(c.Subscribers.ChatIDs, idIndex, idIndex+1)

	return "The user is subscribed for current metal price notifications!"
}
