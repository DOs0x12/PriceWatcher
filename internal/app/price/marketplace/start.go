package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func perStartDur(now time.Time) time.Duration {
	return custTime.PerStartDur(now, custTime.HalfHour)
}
