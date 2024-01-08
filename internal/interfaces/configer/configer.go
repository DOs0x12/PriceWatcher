package configer

import "PriceWatcher/internal/entities/config"

type Configer interface {
	GetConfig() (config.Config, error)
	GetMarketplaceConfig(name string) (config.ServiceConf, error)
	AddItemToWatch(name, address, priceType string) error
	RemoveItemFromWatching(name, priceType string) error
}
