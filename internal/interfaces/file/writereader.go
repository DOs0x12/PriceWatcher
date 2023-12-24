package file

import "PriceWatcher/internal/entities/price"

type WriteReader interface {
	WritePrices(prices map[string]price.ItemPrice) error
	ReadPrices() (map[string]price.ItemPrice, error)
	Lock()
	Unlock()
}
