package bank

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/interfaces/requester"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Service struct {
	req requester.Requester
	ext extractor.Extractor
	val message.HourValidator
	cl  clock.Clock
}

func NewService(req requester.Requester, ext extractor.Extractor, val message.HourValidator, cl clock.Clock) Service {
	return Service{
		req: req,
		ext: ext,
		val: val,
		cl:  cl,
	}
}

var bankUrl = "https://investzoloto.ru/gold-sber-oms/"

func (s Service) ServePrice(conf config.Config) (message, subject string, err error) {
	curHour := s.cl.Now().Hour()

	logrus.Infof("Check time for processing a price. The time value: %v", curHour)

	if !s.val.Validate(curHour, conf.SendingHours) {
		logrus.Info("It is not appropriate time for getting a price")

		return "", "", nil
	}

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
