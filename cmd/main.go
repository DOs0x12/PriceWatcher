package main

import (
	"PriceWatcher/internal/app"
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/domain/hour"
	"PriceWatcher/internal/domain/page"
	"PriceWatcher/internal/infrastructure/configer"
	"PriceWatcher/internal/infrastructure/requester/bank"
	"PriceWatcher/internal/infrastructure/sender"
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

func newService() *app.PriceService {
	req := bank.Requester{}
	sen := sender.Sender{}
	ext := page.PriceExtractor{}
	val := hour.MessageHourVal{}
	conf := configer.Configer{}

	return app.NewPriceService(req, sen, ext, val, conf)
}
