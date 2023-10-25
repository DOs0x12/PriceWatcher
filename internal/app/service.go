package app

import (
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/interfaces/configer"
	interSend "PriceWatcher/internal/interfaces/sender"
	"fmt"
)

type PriceWatcherService struct {
	sender       interSend.Sender
	conf         configer.Configer
	priceService price.PriceService
}

func (s PriceWatcherService) serve() error {
	conf, err := s.conf.GetConfig()
	if err != nil {
		return fmt.Errorf("on getting the config an error occurs: %w", err)
	}

	msg, sub, err := s.priceService.ServePrice(conf)
	if err != nil {
		return err
	}

	err = s.sender.Send(msg, sub, conf.Email)
	if err != nil {
		return fmt.Errorf("cannot send the price: %w", err)
	}

	return nil
}
