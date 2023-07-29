package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"

	"github.com/sirupsen/logrus"
)

var (
	req requester.Requester
	ext domain.PriceExtractor
	sen sender.Sender
)

func main() {
	setGlobalVals()
	serv := app.NewGoldPriceService(req, ext, sen)

	logrus.Infoln("Start the application")

	app.WatchGoldPrice(serv)

	logrus.Infoln("The application has been stopped")
}

func setGlobalVals() {
	req = requester.Requester{}
	ext = domain.PriceExtractor{}
	sen = sender.Sender{}
}
