package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func waitNextStart(now time.Time) (time.Duration, error) {
	period := 60
	return custTime.WaitNextStart(now, period)
}
