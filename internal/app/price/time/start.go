package time

import (
	"fmt"
	"time"
)

func WaitNextStart(now time.Time, targetMin int) (time.Duration, error) {
	waitTime, err := getWaitTime(now, targetMin)
	if err != nil {
		var zeroDur time.Duration

		return zeroDur, err
	}

	time.Sleep(waitTime)

	return waitTime, nil
}

func getWaitTime(now time.Time, targetMin int) (time.Duration, error) {
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
