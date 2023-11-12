package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func perStartDur(now time.Time) (time.Duration, error) {
	targetMin := 60
	return custTime.PerStartDur(now, targetMin)
}
