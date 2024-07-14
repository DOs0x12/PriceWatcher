package time

import (
	"time"
)

func GetCallTime(now time.Time, callHours []int) time.Time {
	curHour := now.Hour()
	nextHour := -1

	for _, hour := range callHours {
		if curHour < hour {
			nextHour = hour

			break
		}
	}

	nextDay := false

	if nextHour == -1 {
		nextHour = callHours[0]
		nextDay = true
	}

	return getCallTimeFromHour(now, nextHour, nextDay)
}

func getCallTimeFromHour(now time.Time, hour int, nextDay bool) time.Time {
	day := now.Day()

	if nextDay {
		day++
	}

	return time.Date(now.Year(), now.Month(), day, hour, 0, 0, 0, now.Location())

}
