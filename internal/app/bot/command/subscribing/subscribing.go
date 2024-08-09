package subscribing

import (
	"PriceWatcher/internal/entities/subscribing"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SubscribingComm struct {
	mu          *sync.Mutex
	Subscribers *subscribing.Subscribers
}

func NewSubCommand(mu *sync.Mutex, subscribers *subscribing.Subscribers) SubscribingComm {
	return SubscribingComm{mu: mu, Subscribers: subscribers}
}

func (c SubscribingComm) SubscribeUser(input interface{}) string {
	upd := input.(tgbotapi.Update)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Subscribers.ChatIDs = append(c.Subscribers.ChatIDs, upd.Message.Chat.ID)

	return "The user is subscribed for current metal price notifications!"
}
