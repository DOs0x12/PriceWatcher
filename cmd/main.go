package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	startInterruptionWatch()
	done := make(chan interface{})
	serv := newService()

	logrus.Infoln("Start the application")

	app.WatchGoldPrice(serv, done)

	logrus.Infoln("The application is done")
}

func newService() *app.GoldPriceService {
	req := requester.Requester{}
	ext := domain.PriceExtractor{}
	sen := sender.Sender{}

	return app.NewGoldPriceService(req, ext, sen)
}

func startInterruptionWatch() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logrus.Infof("The application has been stopped")
			os.Exit(0)
		}
	}()
}
