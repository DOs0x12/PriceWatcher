package marketplace

import (
	"PriceWatcher/internal/app/price/shared"
	"time"
)

func waitNextStart(now time.Time) (time.Duration, error) {
	period := 30
	return shared.WaitNextStart(now, period)
}
