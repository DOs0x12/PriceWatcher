package bank

import (
	bankTime "PriceWatcher/internal/app/bank/time"
	domBank "PriceWatcher/internal/domain/bank"
	entConfig "PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/subscribing"
	infraBank "PriceWatcher/internal/infrastructure/bank"
	interfBank "PriceWatcher/internal/interfaces/bank"
	"PriceWatcher/internal/interfaces/broker"
	"context"
	"sync"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	req  interfBank.Requester
	ext  interfBank.Extractor
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
	broker broker.Worker,
	subscribers *subscribing.Subscribers) {
	defer wg.Done()

	dur := s.getWaitDurWithLogs(time.Now())

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-callChan:
			go s.servePriceWithTiming(ctx, t, broker, subscribers)
		}
	}
}

func (s Service) servePriceWithTiming(
	ctx context.Context,
	timer *time.Timer,
	broker broker.Worker,
	subscribers *subscribing.Subscribers) {
	msg, err := s.getMessageWithPrice()
	if err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)

		return
	}

	logrus.Info("The price is processed")

	if msg == "" {
		s.resetTimer(timer)
		logrus.Infof("A message of the processed price is empty")

		return
	}

	now := time.Now()
	durForMessage := bankTime.DurToSendMessage(now, s.conf.SendingHours)
	logrus.Infof("Waiting the time to send a message: %v", durForMessage)

	select {
	case <-ctx.Done():
		logrus.Infoln("Interrupting waiting the time when to send a message")

		return
	case <-time.After(durForMessage):
	}

	for _, chatID := range subscribers.ChatIDs {
		broker.SendMessage(ctx, msg, chatID)
	}

	s.resetTimer(timer)
}

func (s Service) getWaitDurWithLogs(now time.Time) time.Duration {
	dur := bankTime.GetWaitDurWithRandomComp(now, s.conf.SendingHours)
	logrus.Infof("Waiting %v", dur)

	return dur
}

func (s Service) getMessageWithPrice() (message string, err error) {
	logrus.Infof("Start processing a price")

	response, err := s.req.RequestPage()
	if err != nil {
		return "", fmt.Errorf("cannot get a page with the current price: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return "", fmt.Errorf("cannot extract the price from the body: %w", err)
	}

	msg := fmt.Sprintf("Курс золота. Продажа: %.2fр", price)

	return msg, nil
}

func (s Service) resetTimer(timer *time.Timer) {
	now := time.Now()
	dur := s.getWaitDurWithLogs(now)
	timer.Reset(dur)
}
