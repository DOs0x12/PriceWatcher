package marketplace

import (
	"time"
)

func getWaitDurWithRandomComp(now time.Time, callTime time.Time, randDur time.Duration) time.Duration {
	waitDur := callTime.Sub(now)

	if waitDur < 0 {
		var zeroDur time.Duration
		return zeroDur
	}

	return waitDur + randDur
}
