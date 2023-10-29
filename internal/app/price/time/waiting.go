package time

import (
	"math/rand"
	"time"
)

func GetWaitDurWithRandomComp(now time.Time, callTime time.Time, variation int) time.Duration {
	waitDur := callTime.Sub(now)

	if waitDur < 0 {
		var zeroDur time.Duration
		return zeroDur
	}

	randComp := rand.Intn(variation)
	randDur := time.Duration(randComp)

	if waitDur < randDur {
		return waitDur
	}

	return waitDur - randDur
}
