package time

import (
	"time"
)

func DurToSendMessage(now time.Time, callHours []int) time.Duration {
	timeForMessage := getCallTime(now, callHours)

	return timeForMessage.Sub(now)
}
