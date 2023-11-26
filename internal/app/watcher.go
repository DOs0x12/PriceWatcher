package app

import (
	"PriceWatcher/internal/app/service"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func watch(ctx context.Context, serv service.PriceWatcherService) {
	dur := getWaitTimeWithLogs(serv, time.Now())

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	callCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logrus.Infoln("Shutting down the application")
			return
		case <-callChan:
			go servePriceWithTiming(callCtx, serv, t)
		}
	}
}

func servePriceWithTiming(ctx context.Context, serv service.PriceWatcherService, timer *time.Timer) {
	msg, sub := serveWithLogs(serv)

	now := time.Now()
	dur := perStartWithLogs(serv, now)

	select {
	case <-ctx.Done():
		logrus.Infoln("Interrupting waiting the next period")
		return
	case <-time.After(dur):
	}

	if msg != "" {
		sendReportWithLogs(serv, msg, sub)
	}

	now = time.Now()
	dur = getWaitTimeWithLogs(serv, now)

	timer.Reset(dur)
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
	dur := serv.PerStartDur(now)
	logrus.Infof("Waiting the start of the next period %v", dur)

	return dur
}

func getWaitTimeWithLogs(serv service.PriceWatcherService, now time.Time) time.Duration {
	dur := serv.GetWaitTime(now)
	logrus.Infof("Waiting %v", dur)

	return dur
}
