package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func waitNextStart(now time.Time) (time.Duration, error) {
	targetMin := 30
	return custTime.WaitNextStart(now, targetMin)
}
