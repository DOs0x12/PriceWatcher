package marketplace

import (
	"time"
)

func getCallTime(now time.Time) time.Time {
	curMinutes := now.Minute()
	callPeriod := 30

	if curMinutes < callPeriod {
		return getCallTimeFromMinutes(now, callPeriod)
	}

	return getCallTimeFromMinutes(now, 2*callPeriod)
}

func getCallTimeFromMinutes(now time.Time, minutes int) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minutes, 0, 0, now.Location())
}
