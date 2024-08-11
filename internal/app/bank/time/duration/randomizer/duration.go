package randomizer

import (
	"math/rand"
	"time"
)

func RandomSec(variationInSec int) time.Duration {
	randComp := rand.Intn(variationInSec)
	return time.Duration(randComp) * time.Second
}
