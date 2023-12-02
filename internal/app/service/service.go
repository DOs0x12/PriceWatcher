package service

import (
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/entities/config"
	interSend "PriceWatcher/internal/interfaces/sender"
	"time"
)

type PriceWatcherService struct {
	priceService price.PriceService
	sender       interSend.Sender
	conf         config.ServiceConf
}

func (s PriceWatcherService) Serve() (string, string, error) {
	var err error
	msg, sub, err := s.priceService.ServePrice()
	if err != nil {
		return "", "", err
	}

	return msg, sub, nil
}

func (s PriceWatcherService) SendReport(msg, sub string) error {
	return s.sender.Send(msg, sub, s.conf.Email)
}

func (s PriceWatcherService) GetWaitTime(now time.Time) time.Duration {
	return s.priceService.GetWaitTime(now)
}

func (s PriceWatcherService) PerStartDur(now time.Time) time.Duration {
	return s.priceService.PerStartDur(now)
}

func (s PriceWatcherService) GetName() string {
	return s.priceService.GetName()
}
