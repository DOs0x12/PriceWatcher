package internal

import (
	"PriceWatcher/internal/app/bank"
	"PriceWatcher/internal/entities/subscribing"
	"PriceWatcher/internal/interfaces"
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func ServeMetalPrice(ctx context.Context,
	wg *sync.WaitGroup,
	bankService bank.Service,
	bot interfaces.Bot,
	subscribers *subscribing.Subscribers) {
	defer wg.Done()

	dur := getWaitTimeWithLogs(bankService, time.Now())

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-callChan:
			go servePriceWithTiming(ctx, bankService, t, bot, subscribers)
		}
	}
}

func servePriceWithTiming(
	ctx context.Context,
	serv bank.Service,
	timer *time.Timer,
	bot interfaces.Bot,
	subscribers *subscribing.Subscribers) {
	msg, _ := serveWithLogs(serv)
	if msg != "" {
		for _, chatID := range subscribers.ChatIDs {
			bot.SendMessage(msg, chatID)
		}
	}

	now := time.Now()
	dur := perStartWithLogs(serv, now)

	select {
	case <-ctx.Done():
		logrus.Infoln("Interrupting waiting the next period")

		return
	case <-time.After(dur):
	}

	now = time.Now()
	dur = getWaitTimeWithLogs(serv, now)

	timer.Reset(dur)
}

func serveWithLogs(serv bank.Service) (string, string) {
	msg, sub, err := serv.ServePrice()
	if err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)

		return "", ""
	}

	logrus.Info("The price is processed")

	return msg, sub
}

func perStartWithLogs(serv bank.Service, now time.Time) time.Duration {
	dur := serv.PerStartDur(now)
	logrus.Infof("Waiting the start of the next period %v", dur)

	return dur
}

func getWaitTimeWithLogs(serv bank.Service, now time.Time) time.Duration {
	dur := serv.GetWaitTime(now)
	logrus.Infof("Waiting %v", dur)

	return dur
}
