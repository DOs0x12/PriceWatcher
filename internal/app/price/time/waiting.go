package time

import (
	"math/rand"
	"time"
)

func GetWaitDurWithRandomComp(now time.Time, callTime time.Time, variation int) time.Duration {
	waitDur := callTime.Sub(now)
	randComp := rand.Intn(variation)
	randMin := time.Duration(randComp)
	return waitDur - randMin
}
