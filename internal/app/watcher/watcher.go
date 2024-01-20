package watcher

import (
	"PriceWatcher/internal/app/watcher/price"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/interfaces/sender"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func watch(ctx context.Context,
	serv price.PriceService,
	sen sender.Sender,
	email config.Email) {
	servName := serv.GetName()
	dur := getWaitTimeWithLogs(serv, time.Now(), servName)

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	defer func() {
		finishJobWithLogs(servName)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-callChan:
			go servePriceWithTiming(ctx, serv, sen, t, servName, email)
		}
	}
}

func servePriceWithTiming(
	ctx context.Context,
	serv price.PriceService,
	sen sender.Sender,
	timer *time.Timer,
	servName string,
	email config.Email) {
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
		sendReportWithLogs(sen, msg, sub, servName, email)
	}

	now = time.Now()
	dur = getWaitTimeWithLogs(serv, now, servName)

	timer.Reset(dur)
}

func serveWithLogs(serv price.PriceService, servName string) (string, string) {
	msg, sub, err := serv.ServePrice()
	if err != nil {
		logrus.Errorf("%v: an error occurs while serving a price: %v", servName, err)

		return "", ""
	}

	logrus.Info(servName + ": the price is processed")

	return msg, sub
}

func sendReportWithLogs(sender sender.Sender, msg, sub, servName string, email config.Email) {
	err := sender.Send(msg, sub, email)
	if err != nil {
		logrus.Errorf("%v: cannot send the report: %v", servName, err)
	}

	logrus.Info(servName + ": a report is sended")
}

func perStartWithLogs(serv price.PriceService, now time.Time, servName string) time.Duration {
	dur := serv.PerStartDur(now)
	logrus.Infof("%v: waiting the start of the next period %v", servName, dur)

	return dur
}

func getWaitTimeWithLogs(serv price.PriceService, now time.Time, servName string) time.Duration {
	dur := serv.GetWaitTime(now)
	logrus.Infof("%v: waiting %v", servName, dur)

	return dur
}

func finishJobWithLogs(servName string) {
	logrus.Infof("%v: the job is done", servName)
}
