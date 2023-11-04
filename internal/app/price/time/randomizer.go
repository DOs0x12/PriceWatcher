package time

import (
	"math/rand"
	"time"
)

type Randomizer struct{}

func (Randomizer) RandomMin(variation int) time.Duration {
	randComp := rand.Intn(variation)
	return time.Duration(randComp) * time.Minute
}
