package bank

import (
	priceTime "PriceWatcher/internal/app/bank/time"
	domBank "PriceWatcher/internal/domain/bank"
	entConfig "PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/subscribing"
	infraBank "PriceWatcher/internal/infrastructure/bank"
	"PriceWatcher/internal/interfaces"
	"context"
	"sync"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	req  interfaces.Requester
	ext  interfaces.Extractor
	conf entConfig.Config
}

func NewService(
	req infraBank.BankRequester,
	ext domBank.PriceExtractor,
	conf entConfig.Config) Service {
	return Service{
		req:  req,
		ext:  ext,
		conf: conf,
	}
}

func (s Service) WatchPrice(ctx context.Context,
	wg *sync.WaitGroup,
	bot interfaces.Bot,
	subscribers *subscribing.Subscribers) {
	defer wg.Done()

	dur := s.getWaitTimeWithLogs(time.Now())

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-callChan:
			go s.servePriceWithTiming(ctx, t, bot, subscribers)
		}
	}
}

func (s Service) servePriceWithTiming(
	ctx context.Context,
	timer *time.Timer,
	bot interfaces.Bot,
	subscribers *subscribing.Subscribers) {
	msg, _ := s.serveWithLogs()
	if msg != "" {
		for _, chatID := range subscribers.ChatIDs {
			bot.SendMessage(msg, chatID)
		}
	}

	now := time.Now()
	dur := priceTime.PerStartDur(now)

	select {
	case <-ctx.Done():
		logrus.Infoln("Interrupting waiting the next period")

		return
	case <-time.After(dur):
	}

	now = time.Now()
	dur = s.getWaitTimeWithLogs(now)

	timer.Reset(dur)
}

func (s Service) serveWithLogs() (string, string) {
	msg, sub, err := s.servePrice()
	if err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)

		return "", ""
	}

	logrus.Info("The price is processed")

	return msg, sub
}

func (s Service) getWaitTimeWithLogs(now time.Time) time.Duration {
	dur := s.getWaitTime(now)
	logrus.Infof("Waiting %v", dur)

	return dur
}

func (s Service) servePrice() (message, subject string, err error) {
	logrus.Infof("Start processing a price")

	response, err := s.req.RequestPage()
	if err != nil {
		return "", "", fmt.Errorf("cannot get a page with the current price: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("cannot extract the price from the body: %w", err)
	}

	msg := fmt.Sprintf("Курс золота. Продажа: %.2fр", price)
	sub := "Че по золоту?"

	return msg, sub, nil
}

func (s Service) getWaitTime(now time.Time) time.Duration {
	variation := 1800
	randDur := priceTime.RandomSec(variation)
	callTime := priceTime.GetCallTime(now, s.conf.SendingHours)

	return getWaitDurWithRandomComp(now, callTime, randDur)
}
