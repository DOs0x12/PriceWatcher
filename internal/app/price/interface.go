package price

import "PriceWatcher/internal/interfaces/configer"

type PriceService interface {
	ServePrice(conf configer.Configer) string
}
