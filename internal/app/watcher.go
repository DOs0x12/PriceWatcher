package app

import (
	"PriceWatcher/internal/app/service"
	custTime "PriceWatcher/internal/app/time"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func Watch(done <-chan struct{}, serv service.PriceWatcherService, clock custTime.Clock) {
	dur := getWaitTimeWithLogs(serv)
	t := time.NewTicker(dur)
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Infoln("Shutting down the application")
			return
		case <-t.C:
			serveWithLogs(serv)
			waitToSendRepWithLogs(serv, clock.Now())
			sendReportWithLogs(serv)
			dur = getWaitTimeWithLogs(serv)
			t.Reset(dur)
		}
	}
}

func serveWithLogs(serv service.PriceWatcherService) {
	if err := serv.Serve(); err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)
	}

	logrus.Info("The price is processed")
}

func sendReportWithLogs(serv service.PriceWatcherService) {
	err := serv.SendReport()
	if err != nil {
		logrus.Errorf("cannot send the report: %v", err)
	}
}

func waitToSendRepWithLogs(serv service.PriceWatcherService, now time.Time) {
	err := serv.WaitToSendRep(now)
	if err != nil {
		msg := fmt.Sprintf("An error occurs while waiting when to send a report: %v", err)
		panic(msg)
	}
}

func getWaitTimeWithLogs(serv service.PriceWatcherService) time.Duration {
	dur := serv.GetWaitTime()
	logrus.Infof("Waiting %v", dur)

	return dur
}
