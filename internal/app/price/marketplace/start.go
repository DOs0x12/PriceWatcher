package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func perStartDur(now time.Time) (time.Duration, error) {
	targetMin := 30
	return custTime.PerStartDur(now, targetMin)
}
