package service

import (
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/interfaces/configer"
	interSend "PriceWatcher/internal/interfaces/sender"
	"fmt"
)

func NewWatcherService(sender interSend.Sender, conf configer.Configer) (PriceWatcherService, error) {

	config, err := conf.GetConfig()
	if err != nil {
		return PriceWatcherService{}, fmt.Errorf("can not get the config data: %v", err)
	}

	priceService, err := price.NewPriceService(config.PriceType, config.Marketplace)
	if err != nil {
		return PriceWatcherService{}, err
	}

	service := PriceWatcherService{
		priceService: priceService,
		sender:       sender,
		conf:         conf,
	}

	return service, nil
}
