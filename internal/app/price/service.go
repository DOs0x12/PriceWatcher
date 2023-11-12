package price

import (
	"time"
)

type PriceService interface {
	ServePrice() (message, subject string, err error)
	GetWaitTime(now time.Time) time.Duration
	PerStartDur(now time.Time) time.Duration
}
