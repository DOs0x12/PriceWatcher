package watcher

import (
	"PriceWatcher/internal/app/service"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func watch(ctx context.Context, serv service.PriceWatcherService, finishedJobs chan<- string) {
	servName := serv.GetName()
	dur := getWaitTimeWithLogs(serv, time.Now(), servName)

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	callCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			finishJobWithLogs(servName, finishedJobs)
			return
		case <-callChan:
			go servePriceWithTiming(callCtx, serv, t, servName)
		}
	}
}

func servePriceWithTiming(
	ctx context.Context,
	serv service.PriceWatcherService,
	timer *time.Timer,
	servName string) {
	msg, sub := serveWithLogs(serv, servName)

	now := time.Now()
	dur := perStartWithLogs(serv, now, servName)

	select {
	case <-ctx.Done():
		logrus.Infoln(servName + ": interrupting waiting the next period")
		return
	case <-time.After(dur):
	}

	if msg != "" {
		sendReportWithLogs(serv, msg, sub, servName)
	}

	now = time.Now()
	dur = getWaitTimeWithLogs(serv, now, servName)

	timer.Reset(dur)
}

func serveWithLogs(serv service.PriceWatcherService, servName string) (string, string) {
	msg, sub, err := serv.Serve()
	if err != nil {
		logrus.Errorf("%v: an error occurs while serving a price: %v", servName, err)

		return "", ""
	}

	logrus.Info(servName + ": the price is processed")

	return msg, sub
}

func sendReportWithLogs(serv service.PriceWatcherService, msg, sub, servName string) {
	err := serv.SendReport(msg, sub)
	if err != nil {
		logrus.Errorf("%v: cannot send the report: %v", servName, err)
	}

	logrus.Info(servName + ": a report is sended")
}

func perStartWithLogs(serv service.PriceWatcherService, now time.Time, servName string) time.Duration {
	dur := serv.PerStartDur(now)
	logrus.Infof("%v: waiting the start of the next period %v", servName, dur)

	return dur
}

func getWaitTimeWithLogs(serv service.PriceWatcherService, now time.Time, servName string) time.Duration {
	dur := serv.GetWaitTime(now)
	logrus.Infof("%v: waiting %v", servName, dur)

	return dur
}

func finishJobWithLogs(servName string, finishedJobs chan<- string) {
	logrus.Infoln(servName + ": shutting down the application")
	finishedJobs <- servName
}
