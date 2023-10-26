package time

import (
	"fmt"
	"time"
)

func WaitNextStart(now time.Time, periodInMin int) (time.Duration, error) {
	waitTime, err := getWaitTime(now, periodInMin)
	if err != nil {
		return time.Duration(0), err
	}

	time.Sleep(waitTime)

	return waitTime, nil
}

func getWaitTime(now time.Time, periodInMin int) (time.Duration, error) {
	waitMin := periodInMin - now.Minute()
	waitSec := periodInMin - now.Second()

	durStr := fmt.Sprintf("%vm%vs", waitMin, waitSec)

	return time.ParseDuration(durStr)
}
