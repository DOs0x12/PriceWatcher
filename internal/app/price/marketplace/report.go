package marketplace

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func whenToSendRep(now time.Time) (time.Duration, error) {
	targetMin := 30
	return custTime.WhenToSendRep(now, targetMin)
}
