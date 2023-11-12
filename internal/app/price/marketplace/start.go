package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func waitPerStart(now time.Time) (time.Duration, error) {
	targetMin := 30
	return custTime.WaitPerStart(now, targetMin)
}
