package bank

import (
	priceTime "PriceWatcher/internal/app/price/time"
	"time"
)

func getWaitTime(now time.Time, callHours []int, rand priceTime.Randomizer) time.Duration {
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
	randDur := rand.RandomMin(variation)

	return priceTime.GetWaitDurWithRandomComp(now, callTime, randDur)
}

func getCallTime(now time.Time, hour int, nextDay bool) time.Time {
	day := now.Day()

	if nextDay {
		day++
	}

	return time.Date(now.Year(), now.Month(), day, hour, 0, 0, 0, now.Location())

}
