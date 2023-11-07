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
	conf         config.Config
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

func (s PriceWatcherService) WaitToSendRep(now time.Time) error {
	dur, err := s.priceService.WaitPerStart(now)
	if err != nil {
		return err
	}

	time.Sleep(dur)

	return nil
}
