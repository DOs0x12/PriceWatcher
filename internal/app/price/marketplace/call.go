package marketplace

import (
	"time"
)

func getCallTime(now time.Time) time.Time {
	curMinutes := now.Minute()
	callPeriod := 30
	var callMinutes int

	if curMinutes <= callPeriod {
		callMinutes = callPeriod
	} else {
		callMinutes = 2 * callPeriod
	}

	return getCallTimeFromMinutes(now, callMinutes)
}

func getCallTimeFromMinutes(now time.Time, minutes int) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minutes, 0, 0, now.Location())
}
