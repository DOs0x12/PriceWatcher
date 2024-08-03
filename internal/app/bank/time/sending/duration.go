package sending

import (
	"PriceWatcher/internal/app/bank/time/call"
	"time"
)

func DurToSendMessage(now time.Time, callHours []int) time.Duration {
	timeForMessage := call.GetCallTime(now, callHours)

	return timeForMessage.Sub(now)
}
