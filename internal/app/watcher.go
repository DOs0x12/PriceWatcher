package app

import (
	"time"

	"github.com/sirupsen/logrus"
)

func WatchGoldPrice(s Service) {
	t := time.NewTicker(1 * time.Hour)
	done := make(chan interface{})

	for {
		select {
		case <-done:
			logrus.Infoln("Done is called")
			return
		case <-t.C:
			handlePrice(s)
		}
	}
}

func isMessageHour(hour int) bool {
	hours := []int{12, 17}

	for _, h := range hours {
		if h == hour {
			return true
		}
	}

	return false
}

func handlePrice(s Service) {
	logrus.Infoln("Check time for handling a gold price")

	if !isMessageHour(time.Now().Hour()) {
		return
	}

	logrus.Infoln("Start handling a gold price")

	err := s.HandlePrice()
	if err != nil {
		logrus.Errorf("Handling the gold price ends with the error: %v", err)
	}

	logrus.Info("The gold price is handled")
}
