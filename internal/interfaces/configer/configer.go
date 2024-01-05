package configer

import "PriceWatcher/internal/entities/config"

type Configer interface {
	GetConfig() (config.Config, error)
	AddItemToWatch(address, name, priceType string) error
	RemoveItemFromWatching(address, name, priceType string) error
}
