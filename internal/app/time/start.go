package time

import (
	"fmt"
	"time"
)

func WaitHourStart(now time.Time) error {
	waitTime, err := getWaitTime(now)
	if err != nil {
		return err
	}

	time.Sleep(waitTime)

	return nil
}

func getWaitTime(now time.Time) (time.Duration, error) {
	waitMin := 60 - now.Minute()
	waitSec := 60 - now.Second()

	durStr := fmt.Sprintf("%vm%vs", waitMin, waitSec)

	return time.ParseDuration(durStr)
}
