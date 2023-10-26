package price

import "time"

func IsAppropriateDelay(del time.Duration) bool {
	thresholdInMin := 10

	return del > time.Duration(thresholdInMin)
}
