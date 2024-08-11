package subscriber

import "PriceWatcher/internal/entities/bank/subscribing"

type Worker interface {
	GetSubscribers() (*subscribing.Subscribers, error)
	SaveSubscribers(subs *subscribing.Subscribers) error
}
