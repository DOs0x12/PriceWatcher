package app

import (
	"PriceWatcher/internal/app/service"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func Watch(done <-chan struct{}, serv service.PriceWatcherService) {
	dur := getWaitTimeWithLogs(serv, time.Now())

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Infoln("Shutting down the application")
			return
		case <-callChan:
			msg, sub := serveWithLogs(serv)

			now := time.Now()
			dur = perStartWithLogs(serv, now)

			time.Sleep(dur)

			if msg != "" {
				sendReportWithLogs(serv, msg, sub)
			}

			now = time.Now()
			dur = getWaitTimeWithLogs(serv, now)

			t.Reset(dur)
		}
	}
}

func serveWithLogs(serv service.PriceWatcherService) (string, string) {
	msg, sub, err := serv.Serve()
	if err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)

		return "", ""
	}

	logrus.Info("The price is processed")

	return msg, sub
}

func sendReportWithLogs(serv service.PriceWatcherService, msg, sub string) {
	err := serv.SendReport(msg, sub)
	if err != nil {
		logrus.Errorf("cannot send the report: %v", err)
	}

	logrus.Info("A report is sended")
}

func perStartWithLogs(serv service.PriceWatcherService, now time.Time) time.Duration {
	dur, err := serv.PerStartDur(now)
	if err != nil {
		msg := fmt.Sprintf("An error occurs while waiting the next period start: %v", err)
		panic(msg)
	}

	logrus.Infof("Waiting the start of the next period %v", dur)

	return dur
}

func getWaitTimeWithLogs(serv service.PriceWatcherService, now time.Time) time.Duration {
	dur := serv.GetWaitTime(now)
	logrus.Infof("Waiting %v", dur)

	return dur
}
