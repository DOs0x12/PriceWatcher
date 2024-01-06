package configer

import "PriceWatcher/internal/entities/config"

type Configer interface {
	GetConfig() (config.Config, error)
	GetMarketplaceConfig(name string) (config.ServiceConf, error)
	AddItemToWatch(address, name, priceType string) error
	RemoveItemFromWatching(address, name, priceType string) error
}
