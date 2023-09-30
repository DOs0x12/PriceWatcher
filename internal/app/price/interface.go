package price

import (
	"PriceWatcher/internal/entities/config"
)

type PriceService interface {
	ServePrice(conf config.Config) (message, subject string, err error)
}
