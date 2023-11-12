package time

import (
	"time"
)

func WaitPerStart(now time.Time, period int) (time.Duration, error) {
	curPer := period

	for now.Minute() > curPer {
		curPer += period
	}

	tarTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), curPer, 0, 0, now.Location())

	return tarTime.Sub(now), nil
}
