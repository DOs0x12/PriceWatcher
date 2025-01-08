package time

import (
	"math/rand"
	"time"
)

func GetWaitDurWithRandomComp(now time.Time, callHours []int) time.Duration {
	variation := 1800
	randDur := randomDurSec(variation)
	callTime := getCallTime(now, callHours)
	waitDur := callTime.Sub(now)
	processingTime := 3 * time.Minute
	randComp := randDur + processingTime

	if waitDur < randComp {
		return 0 * time.Second
	}

	return waitDur - randComp
}

func randomDurSec(variationInSec int) time.Duration {
	randComp := rand.Intn(variationInSec)
	return time.Duration(randComp) * time.Second
}

func getCallTime(now time.Time, callHours []int) time.Time {
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
