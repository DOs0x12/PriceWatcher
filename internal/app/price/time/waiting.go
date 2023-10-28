package time

import (
	"math/rand"
	"time"
)

func GetWaitDurWithRandomComp(now time.Time, callTime time.Time, variation int) time.Duration {
	waitDur := callTime.Sub(now)

	if waitDur < 0 {
		return time.Duration(0)
	}

	randComp := rand.Intn(variation)
	randDur := time.Duration(randComp)
	return waitDur - randDur
}
