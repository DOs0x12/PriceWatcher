package price

import (
	"PriceWatcher/internal/entities/config"
	"time"
)

type PriceService interface {
	ServePrice(conf config.Config) (message, subject string, err error)
	GetWaitTime() time.Duration
	WaitNextStart(now time.Time) (time.Duration, error)
}
