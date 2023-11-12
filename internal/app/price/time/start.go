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
		tarTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

		return tarTime.Sub(now)
	}

	halfHour := 30
	tarTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), halfHour, 0, 0, now.Location())

	if now.Minute() > halfHour {
		tarTime = tarTime.Add(time.Duration(halfHour) * time.Minute)
	}

	return tarTime.Sub(now)

}
