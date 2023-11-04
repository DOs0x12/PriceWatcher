package bank

import (
	custTime "PriceWatcher/internal/app/price/time"
	"time"
)

func getWaitTime(now time.Time, callHours []int) time.Duration {
	curHour := now.Hour()
	nextHour := -1

	for hour := range callHours {
		if curHour <= hour {
			nextHour = hour

			break
		}
	}

	nextDay := false

	if nextHour == -1 {
		nextHour = callHours[0]
		nextDay = true
	}

	callTime := getCallTime(now, nextHour, nextDay)

	variation := 20

	return custTime.GetWaitDurWithRandomComp(now, callTime, variation)
}

func getCallTime(now time.Time, hour int, nextDay bool) time.Time {
	day := now.Day()

	if nextDay {
		day++
	}

	return time.Date(now.Year(), now.Month(), day, hour, 0, 0, 0, now.Location())

}
