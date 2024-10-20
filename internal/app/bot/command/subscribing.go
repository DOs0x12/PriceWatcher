package command

import (
	"PriceWatcher/internal/entities/bot"
	"PriceWatcher/internal/entities/subscribing"

	"sync"

	"golang.org/x/exp/slices"
)

type subscribingComm struct {
	mu          *sync.Mutex
	Subscribers *subscribing.Subscribers
}

func newSubCommand(mu *sync.Mutex, subscribers *subscribing.Subscribers) subscribingComm {
	return subscribingComm{mu: mu, Subscribers: subscribers}
}

func (c subscribingComm) subscribeUser(msg bot.Message) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	idIndex := slices.Index(c.Subscribers.ChatIDs, msg.ChatID)
	if idIndex != -1 {
		return "The user is already subscribed!"
	}

	c.Subscribers.ChatIDs = append(c.Subscribers.ChatIDs, msg.ChatID)

	return "The user is subscribed for current gold price notifications!"
}
