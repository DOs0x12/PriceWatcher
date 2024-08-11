package command

import (
	"PriceWatcher/internal/entities/bank/subscribing"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
)

type subscribingComm struct {
	mu          *sync.Mutex
	Subscribers *subscribing.Subscribers
}

func newSubCommand(mu *sync.Mutex, subscribers *subscribing.Subscribers) subscribingComm {
	return subscribingComm{mu: mu, Subscribers: subscribers}
}

func (c subscribingComm) subscribeUser(input interface{}) string {
	upd := input.(tgbotapi.Update)
	c.mu.Lock()
	defer c.mu.Unlock()

	idIndex := slices.Index(c.Subscribers.ChatIDs, upd.Message.Chat.ID)
	if idIndex != -1 {
		return "The user is already subscribed!"
	}

	c.Subscribers.ChatIDs = append(c.Subscribers.ChatIDs, upd.Message.Chat.ID)

	return "The user is subscribed for current gold price notifications!"
}
