package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func whenToSendRep(now time.Time) (time.Duration, error) {
	targetMin := 60
	return custTime.WhenToSendRep(now, targetMin)
}
