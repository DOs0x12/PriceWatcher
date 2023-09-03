package app

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/app/interrupt"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/interfaces/configer"
	interReq "PriceWatcher/internal/interfaces/requester"
	interSend "PriceWatcher/internal/interfaces/sender"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type PriceService struct {
	req      interReq.Requester
	sender   interSend.Sender
	ext      extractor.Extractor
	val      message.HourValidator
	analyser analyser.Analyser
	conf     configer.Configer
}

func (s *PriceService) serve(clock clock.Clock) error {
	curHour := clock.Now().Hour()

	conf, err := s.conf.GetConfig()
	if err != nil {
		return fmt.Errorf("on getting the config an error occurs: %w", err)
	}

	logrus.Infof("Check time for processing a price. The time value: %v", curHour)

	if !s.val.Validate(curHour, conf.SendingHours) {
		logrus.Info("It is not appropriate time for getting a price")

		return nil
	}

	logrus.Info("Start processing a price")

	response, err := s.req.RequestPage()
	if err != nil {
		return fmt.Errorf("cannot get a page with the current price: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return fmt.Errorf("cannot extract the price from the body: %w", err)
	}

	if s.analyser != nil {
		changed, _, _ := s.analyser.IsPriceChanged(price)

		if changed {
			//TO DO: send a report
			logrus.Info("The item price has been changed. A report is sended")

			return nil
		}

		logrus.Info("The item price has been not changed")

		return nil
	}

	err = s.sender.Send(price, conf.Email)
	if err != nil {
		return fmt.Errorf("cannot send the price: %w", err)
	}

	logrus.Info("The price is processed")

	return nil
}

func (s *PriceService) Watch(done <-chan struct{}, cancel context.CancelFunc, clock clock.Clock) {
	interrupt.WatchForInterruption(cancel)

	errMes := "An error occurs while serving a price: %v"
	if err := s.serve(clock); err != nil {
		logrus.Errorf(errMes, err)
	}

	err := waitHourStart(clock.Now())
	if err != nil {
		logrus.Errorf("An error occurs while waiting when the next hour begins: %v", err)
	}

	t := time.NewTicker(1 * time.Hour)
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Info("Shut down the application")
			return
		case <-t.C:
			if err := s.serve(clock); err != nil {
				logrus.Errorf(errMes, err)
			}
		}
	}
}
