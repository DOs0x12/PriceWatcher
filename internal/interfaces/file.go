package bank

import "PriceWatcher/internal/entities/bank/subscribing"

type Subscribers interface {
	GetSubscribers() (*subscribing.Subscribers, error)
	SaveSubscribers(subs *subscribing.Subscribers) error
}
