package marketplace

import (
	"PriceWatcher/internal/app/price/shared"
	"time"
)

func getWaitTime() time.Duration {
	base := 20
	variation := 10
	return shared.GetWaitTimeInMinutes(base, variation)
}
