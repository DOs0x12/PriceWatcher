package app

import (
	"time"
)

func WatchGoldPrice(s Service) {
	t := time.NewTicker(1 * time.Hour)
	done := make(chan interface{})

	for {
		select {
		case <-done:
			return
		case <-t.C:
			if isMessageHour(time.Now().Hour()) {
				s.HandlePrice()
			}
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
