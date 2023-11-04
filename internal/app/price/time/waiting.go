package time

import (
	"time"
)

func GetWaitDurWithRandomComp(now time.Time, callTime time.Time, randDur time.Duration) time.Duration {
	waitDur := callTime.Sub(now)

	if waitDur < 0 {
		var zeroDur time.Duration
		return zeroDur
	}

	if waitDur < randDur {
		return waitDur
	}

	return waitDur - randDur
}
