package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func waitNextStart(now time.Time) (time.Duration, error) {
	period := 30
	return custTime.WaitNextStart(now, period)
}
