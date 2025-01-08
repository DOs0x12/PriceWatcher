package bank

import "PriceWatcher/internal/entities/subscribing"

type Worker interface {
	GetSubscribers() (*subscribing.Subscribers, error)
	SaveSubscribers(subs *subscribing.Subscribers) error
}
