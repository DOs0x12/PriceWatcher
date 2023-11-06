package service

import (
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/interfaces/configer"
	interSend "PriceWatcher/internal/interfaces/sender"
	"time"
)

type PriceWatcherService struct {
	priceService price.PriceService
	sender       interSend.Sender
	conf         configer.Configer
}

var (
	conf config.Config
	msg  string
	sub  string
)

func (s PriceWatcherService) Serve() error {
	var err error
	msg, sub, err = s.priceService.ServePrice()
	if err != nil {
		return err
	}

	return nil
}

func (s PriceWatcherService) SendReport() error {
	return s.sender.Send(msg, sub, conf.Email)
}

func (s PriceWatcherService) GetWaitTime(now time.Time) time.Duration {
	return s.priceService.GetWaitTime(now)
}

func (s PriceWatcherService) WaitToSendRep(now time.Time) error {
	dur, err := s.priceService.WhenToSendRep(now)
	if err != nil {
		return err
	}

	time.Sleep(dur)

	return nil
}
