package time

import (
	"math/rand"
	"time"
)

func RandomMin(variation int) time.Duration {
	randComp := rand.Intn(variation)
	return time.Duration(randComp) * time.Minute
}
