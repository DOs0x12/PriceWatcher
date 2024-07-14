package bank

import (
	priceTime "PriceWatcher/internal/app/bank/time"
	domBank "PriceWatcher/internal/domain/bank"
	entConfig "PriceWatcher/internal/entities/config"
	infraBank "PriceWatcher/internal/infrastructure/bank"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	req  infraBank.BankRequester
	ext  domBank.PriceExtractor
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

func (s Service) ServePrice() (message, subject string, err error) {
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

func (s Service) GetWaitTime(now time.Time) time.Duration {
	variation := 1800
	randDur := priceTime.RandomSec(variation)
	callTime := priceTime.GetCallTime(now, s.conf.SendingHours)

	return getWaitDurWithRandomComp(now, callTime, randDur)
}

func (Service) PerStartDur(now time.Time) time.Duration {
	return priceTime.PerStartDur(now)
}
