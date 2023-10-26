package app

import (
	"PriceWatcher/internal/app/service"
	custTime "PriceWatcher/internal/app/time"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func Watch(done <-chan struct{}, serv service.PriceWatcherService, clock custTime.Clock) {
	serveWithLogs(serv)
	sendReportWithLogs(serv)

	delay := waitNextStartWithLogs(serv, clock.Now())

	if serv.CanCall(delay) {
		serveWithLogs(serv)
		sendReportWithLogs(serv)
	}

	dur := serv.GetWaitTime()

	t := time.NewTicker(dur)
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Info("Shutting down the application")
			return
		case <-t.C:
			serveWithLogs(serv)
			_ = waitNextStartWithLogs(serv, clock.Now())
			sendReportWithLogs(serv)
			dur = serv.GetWaitTime()
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

func waitNextStartWithLogs(serv service.PriceWatcherService, now time.Time) time.Duration {
	delay, err := serv.WaitNextStart(now)
	if err != nil {
		msg := fmt.Sprintf("An error occurs while waiting when the next hour begins: %v", err)
		panic(msg)
	}

	return delay
}
