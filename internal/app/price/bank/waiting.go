package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func getWaitTime() time.Duration {
	base := 40
	variation := 20
	return custTime.GetWaitTimeInMinutes(base, variation)
}
