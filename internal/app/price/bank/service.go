package bank

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/interfaces/file"

	"github.com/sirupsen/logrus"
)

type Service struct {
	val      message.HourValidator
	analyser analyser.Analyser
	wr       file.WriteReader
	cl       clock.Clock
}

func (s Service) ServePrice(conf config.Config) (string, error) {
	curHour := s.cl.Now().Hour()

	logrus.Infof("Check time for processing a price. The time value: %v", curHour)

	if !s.val.Validate(curHour, conf.SendingHours) {
		logrus.Info("It is not appropriate time for getting a price")

		return "", nil
	}

	logrus.Info("Start processing a price")

	return "", nil
}
