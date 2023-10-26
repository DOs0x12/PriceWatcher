package app

import (
	"PriceWatcher/internal/app/price"
	custTime "PriceWatcher/internal/app/time"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (s PriceWatcherService) Watch(done <-chan struct{}, cancel context.CancelFunc, clock custTime.Clock) {
	WatchForInterruption(cancel)

	serveWithLogs(s)

	delay, err := s.priceService.WaitNextStart(clock.Now())
	if err != nil {
		logrus.Errorf("An error occurs while waiting when the next hour begins: %v", err)

		return
	}

	if price.IsAppropriateDelay(delay) {
		serveWithLogs(s)
	}

	dur := s.priceService.GetWaitTime()

	t := time.NewTicker(dur)
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Info("Shutting down the application")
			return
		case <-t.C:
			serveWithLogs(s)
			dur = s.priceService.GetWaitTime()
			t.Reset(dur)
		}
	}
}

func serveWithLogs(serv PriceWatcherService) {
	if err := serv.serve(); err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)
	}

	logrus.Info("The price is processed")
}
