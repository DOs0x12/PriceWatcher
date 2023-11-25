package service

import (
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/entities/config"
	interSend "PriceWatcher/internal/interfaces/sender"
)

func NewWatcherService(sender interSend.Sender, config config.ServiceConf) (PriceWatcherService, error) {
	priceService, err := price.NewPriceService(config)
	if err != nil {
		return PriceWatcherService{}, err
	}

	service := PriceWatcherService{
		priceService: priceService,
		sender:       sender,
		conf:         config,
	}

	return service, nil
}
