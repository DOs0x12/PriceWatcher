package app

import (
	custTime "PriceWatcher/internal/app/time"
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (s PriceWatcherService) Watch(done <-chan struct{}, cancel context.CancelFunc, clock custTime.Clock) {
	WatchForInterruption(cancel)

	serveWithLogs(s)

	config, err := s.conf.GetConfig()
	if err != nil {
		logrus.Errorf("An error occurs while get the config data: %v", err)

		return
	}

	//TODO: rewrite waiting for the all services; if the previous processing occured less than 10 mins ago then no need
	//for the second processing
	if strings.ToLower(config.PriceType) == "bank" {
		err := custTime.WaitHourStart(clock.Now())
		if err != nil {
			logrus.Errorf("An error occurs while waiting when the next hour begins: %v", err)

			return
		}

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
