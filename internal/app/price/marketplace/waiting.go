package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func getWaitTime() time.Duration {
	base := 20
	variation := 10
	return custTime.GetWaitTimeInMinutes(base, variation)
}
