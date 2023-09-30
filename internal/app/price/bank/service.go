package bank

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Service struct {
	val      message.HourValidator
	analyser analyser.Analyser
	wr       file.WriteReader
	conf     configer.Configer
	cl       clock.Clock
}

func (s Service) ServePrice() (string, error) {
	curHour := s.cl.Now().Hour()
	conf, err := s.conf.GetConfig()
	if err != nil {
		return "", fmt.Errorf("on getting the config an error occurs: %w", err)
	}
	logrus.Infof("Check time for processing a price. The time value: %v", curHour)
	if !s.val.Validate(curHour, conf.SendingHours) {
		logrus.Info("It is not appropriate time for getting a price")

		return nil
	}

	return ""
}
