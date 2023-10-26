package time

import "time"

func CanCall(del time.Duration) bool {
	threshold := 10
	thresholdInMin := time.Duration(threshold) * time.Minute
	return del > thresholdInMin
}
