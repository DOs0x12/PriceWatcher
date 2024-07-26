package time

import (
	"time"
)

func DurToSendMessage(now time.Time, callHours []int) time.Duration {
	timeForMessage := GetCallTime(now, callHours)

	return timeForMessage.Sub(now)
}
