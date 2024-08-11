package command

import (
	"PriceWatcher/internal/entities/subscribing"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
)

type unsubscribingComm struct {
	mu          *sync.Mutex
	Subscribers *subscribing.Subscribers
}

func newUnsubCommand(mu *sync.Mutex, subscribers *subscribing.Subscribers) unsubscribingComm {
	return unsubscribingComm{mu: mu, Subscribers: subscribers}
}

func (c unsubscribingComm) unsubscribeUser(input interface{}) string {
	errMessage := "The user is not subscribed!"

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Subscribers.ChatIDs) == 0 {
		return errMessage
	}

	upd := input.(tgbotapi.Update)
	idIndex := slices.Index(c.Subscribers.ChatIDs, upd.Message.Chat.ID)
	if idIndex == -1 {
		return errMessage
	}

	c.Subscribers.ChatIDs = slices.Delete(c.Subscribers.ChatIDs, idIndex, idIndex+1)

	return "The user is unsubscribed from current gold price notifications!"
}
