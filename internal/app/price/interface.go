package price

import (
	"PriceWatcher/internal/entities/config"
)

type PriceService interface {
	ServePrice(conf config.Config) (string, error)
}
