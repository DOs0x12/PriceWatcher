package marketplace

import "PriceWatcher/internal/interfaces/configer"

type Service struct{}

func (Service) ServePrice(conf configer.Configer) (message, subject string, err error) {
	return "", "", nil
}
