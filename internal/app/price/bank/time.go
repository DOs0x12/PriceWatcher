package bank

import (
	"PriceWatcher/internal/app/price/shared"
	"time"
)

func getWaitTime() time.Duration {
	base := 50
	variation := 10
	return shared.GetWaitTimeInMinutes(base, variation)
}
