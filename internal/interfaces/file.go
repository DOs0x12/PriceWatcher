package interfaces

import "PriceWatcher/internal/entities/subscribing"

type Subscribers interface {
	GetSubscribers() (*subscribing.Subscribers, error)
	SaveSubscribers(subs *subscribing.Subscribers) error
}
