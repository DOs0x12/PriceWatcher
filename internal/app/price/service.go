package price

import (
	"time"
)

type PriceService interface {
	ServePrice() (message, subject string, err error)
	GetWaitTime() time.Duration
	WaitNextStart(now time.Time) (time.Duration, error)
}
