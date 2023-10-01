package app

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"
	interReq "PriceWatcher/internal/interfaces/requester"
	interSend "PriceWatcher/internal/interfaces/sender"
	"fmt"

	"github.com/sirupsen/logrus"
)

type PriceService struct {
	req          interReq.Requester
	sender       interSend.Sender
	ext          extractor.Extractor
	val          message.HourValidator
	analyser     analyser.Analyser
	wr           file.WriteReader
	conf         configer.Configer
	priceService price.PriceService
}

func (s *PriceService) serve(clock clock.Clock) error {
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

	logrus.Info("The price is processed")

	return nil
}
