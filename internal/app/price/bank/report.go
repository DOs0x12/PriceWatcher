package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func waitPerStart(now time.Time) (time.Duration, error) {
	targetMin := 60
	return custTime.WaitPerStart(now, targetMin)
}
