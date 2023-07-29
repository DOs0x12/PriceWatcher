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

var (
	req requester.Requester
	ext domain.PriceExtractor
	sen sender.Sender
)

func main() {
	startInterruptionWatch()
	setGlobalVals()
	serv := app.NewGoldPriceService(req, ext, sen)

	logrus.Infoln("Start the application")

	app.WatchGoldPrice(serv)

	logrus.Infoln("The application is done")
}

func setGlobalVals() {
	req = requester.Requester{}
	ext = domain.PriceExtractor{}
	sen = sender.Sender{}
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
