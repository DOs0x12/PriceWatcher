package bank

import (
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"

	"github.com/sirupsen/logrus"
)

type Service struct {
	val      message.HourValidator
	analyser analyser.Analyser
	wr       file.WriteReader
}

func (s Service) ServePrice(conf configer.Configer) string {
	logrus.Infof("Check time for processing a price. The time value: %v", curHour)
	if !s.val.Validate(curHour, conf.SendingHours) {
		logrus.Info("It is not appropriate time for getting a price")

		return nil
	}

	return ""
}
