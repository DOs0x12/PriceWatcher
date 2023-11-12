package bank

import (
	priceTime "PriceWatcher/internal/app/price/time"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/interfaces/requester"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	req  requester.Requester
	ext  extractor.Extractor
	conf config.Config
}

func NewService(
	req requester.Requester,
	ext extractor.Extractor,
	conf config.Config) Service {
	return Service{
		req:  req,
		ext:  ext,
		conf: conf,
	}
}

var bankUrl = "https://investzoloto.ru/gold-sber-oms/"

func (s Service) ServePrice() (message, subject string, err error) {
	logrus.Info("Start processing a price")

	response, err := s.req.RequestPage(bankUrl)
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
	variation := 20
	rand := priceTime.Randomizer{}
	randDur := rand.RandomMin(variation)
	callTime := getCallTime(now, s.conf.SendingHours)

	return priceTime.GetWaitDurWithRandomComp(now, callTime, randDur)
}

func (Service) PerStartDur(now time.Time) time.Duration {
	return perStartDur(now)
}
