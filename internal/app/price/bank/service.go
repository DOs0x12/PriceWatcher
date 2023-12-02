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
	conf config.ServiceConf
}

func NewService(
	req requester.Requester,
	ext extractor.Extractor,
	conf config.ServiceConf) Service {
	return Service{
		req:  req,
		ext:  ext,
		conf: conf,
	}
}

var bankUrl = "https://investzoloto.ru/gold-sber-oms/"

func (s Service) ServePrice() (message, subject string, err error) {
	serveName := s.GetName()

	logrus.Infof("%v: start processing a price", serveName)

	response, err := s.req.RequestPage(bankUrl)
	if err != nil {
		return "", "", fmt.Errorf("%v: cannot get a page with the current price: %w", serveName, err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("%v: cannot extract the price from the body: %w", serveName, err)
	}

	msg := fmt.Sprintf("Курс золота. Продажа: %.2fр", price)
	sub := "Че по золоту?"

	return msg, sub, nil
}

func (s Service) GetWaitTime(now time.Time) time.Duration {
	variation := 1800
	randDur := priceTime.RandomSec(variation)
	callTime := getCallTime(now, s.conf.SendingHours)

	return getWaitDurWithRandomComp(now, callTime, randDur)
}

func (Service) PerStartDur(now time.Time) time.Duration {
	return priceTime.PerStartDur(now, priceTime.Hour)
}

func (s Service) GetName() string {
	return s.conf.PriceType
}
