package price

import (
	"time"
)

type PriceService interface {
	ServePrice() (message, subject string, err error)
	GetWaitTime(now time.Time) time.Duration
	WaitPerStart(now time.Time) (time.Duration, error)
}
