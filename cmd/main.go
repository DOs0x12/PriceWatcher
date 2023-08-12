package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/app/clock"
	"GoldPriceGetter/internal/domain/hour"
	"GoldPriceGetter/internal/domain/page"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	serv := newService()

	logrus.Infoln("Start the application")

	serv.Watch(ctx.Done(), cancel, clock.RealClock{})

	logrus.Infoln("The application is done")
}

func newService() *app.GoldPriceService {
	req := requester.Requester{}
	sen := sender.Sender{}
	ext := page.PriceExtractor{}
	val := hour.MessageHourVal{}

	return app.NewGoldPriceService(req, sen, ext, val)
}
