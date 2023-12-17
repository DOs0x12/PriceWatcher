package time

import (
	"time"
)

type NearestTime int

const (
	Hour NearestTime = iota
	HalfHour
)

func PerStartDur(now time.Time, nt NearestTime) time.Duration {
	if nt == Hour {
		starTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

		return starTime.Sub(now)
	}

	halfHour := 30
	starTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), halfHour, 0, 0, now.Location())

	if now.Minute() > halfHour {
		starTime = starTime.Add(time.Duration(halfHour) * time.Minute)
	}

	return starTime.Sub(now)
}
