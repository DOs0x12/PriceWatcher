package shared

import (
	"math/rand"
	"time"
)

func GetWaitTimeInMinutes(base, variation int) time.Duration {
	return time.Duration(base+rand.Intn(variation)) * time.Minute
}
