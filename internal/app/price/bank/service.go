package bank

import "PriceWatcher/internal/interfaces/configer"

type Service struct{}

func (Service) ServePrice(conf configer.Configer) string {
	return ""
}
