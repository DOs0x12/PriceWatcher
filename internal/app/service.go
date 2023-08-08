package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/interfaces"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type GoldPriceService struct {
	req    interfaces.Requester
	sender interfaces.Sender
	ext    domain.Extractor
	val    domain.HourValidator
}

func NewGoldPriceService(
	req interfaces.Requester,
	sender interfaces.Sender,
	ext domain.Extractor,
	val domain.HourValidator) *GoldPriceService {

	serv := GoldPriceService{
		req:    req,
		sender: sender,
		ext:    ext,
		val:    val,
	}

	return &serv
}

var nowHour = time.Now().Hour

func (s *GoldPriceService) serve() error {
	logrus.Infoln("Check time for processing a gold price")

	curHour := nowHour()

	if !s.val.Validate(curHour) {
		return nil
	}

	logrus.Infoln("Start processing a gold price")

	response, err := s.req.RequestPage()
	if err != nil {
		return fmt.Errorf("cannot get a page with the current price of gold: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return fmt.Errorf("cannot extract the gold price from the body: %w", err)
	}

	err = s.sender.Send(price)
	if err != nil {
		return fmt.Errorf("cannot send the gold price: %w", err)
	}

	logrus.Info("The gold price is processed")

	return nil
}

func (s *GoldPriceService) Watch(done <-chan struct{}, cancel context.CancelFunc) {
	watchForInterruption(cancel)

	t := time.NewTicker(1 * time.Hour)

	for {
		select {
		case <-done:
			logrus.Infoln("Shut down the application")
			t.Stop()
			return
		case <-t.C:
			err := s.serve()
			if err != nil {
				logrus.Errorf("The error occurs while serving a gold price: %v", err)
			}
		}
	}
}

func watchForInterruption(cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
}
