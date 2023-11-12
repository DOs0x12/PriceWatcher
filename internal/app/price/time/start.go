package time

import (
	"fmt"
	"time"
)

func WaitPerStart(now time.Time, targetMin int) (time.Duration, error) {
	minInHour := 60
	secInMin := 60
	curMin := now.Minute()
	waitMin := targetMin - curMin
	waitSec := secInMin - now.Second()

	if targetMin < curMin {
		waitMin += minInHour
	}

	if waitSec < secInMin {
		waitMin--
	}

	if waitSec == secInMin {
		waitSec = 0
	}

	durStr := fmt.Sprintf("%vm%vs", waitMin, waitSec)

	return time.ParseDuration(durStr)
}
