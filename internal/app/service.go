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

func (s *GoldPriceService) serve(clock Clock) error {
	curHour := clock.Now().Hour()

	logrus.Infof("Check time for processing a gold price. The time value: %v", curHour)

	if !s.val.Validate(curHour) {
		logrus.Info("It is not appropriate time for getting a price")

		return nil
	}

	logrus.Info("Start processing a gold price")

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

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

func (s *GoldPriceService) Watch(done <-chan struct{}, cancel context.CancelFunc, clock Clock) {
	watchForInterruption(cancel)

	errMes := "An error occurs while serving a gold price: %v"
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

func watchForInterruption(cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
}

func waitHourStart(now time.Time) error {
	waitTime, err := getWaitTime(now)
	if err != nil {
		return err
	}

	time.Sleep(waitTime)

	return nil
}

func getWaitTime(now time.Time) (time.Duration, error) {
	waitMin := 60 - now.Minute()
	waitSec := 60 - now.Second()

	durStr := fmt.Sprintf("%vm%vs", waitMin, waitSec)

	return time.ParseDuration(durStr)
}
