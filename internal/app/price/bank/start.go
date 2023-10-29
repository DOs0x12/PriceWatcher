package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func waitNextStart(now time.Time) (time.Duration, error) {
	targetMin := 60
	return custTime.WaitNextStart(now, targetMin)
}
