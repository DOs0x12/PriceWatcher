package time

import (
	"time"
)

func PerStartDur(now time.Time) time.Duration {
	starTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

	return starTime.Sub(now)
}
