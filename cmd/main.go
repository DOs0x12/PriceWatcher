package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	serv := newService()

	logrus.Infoln("Start the application")

	serv.Watch(ctx.Done(), cancel)

	logrus.Infoln("The application is done")
}

func newService() *app.GoldPriceService {
	req := requester.Requester{}
	sen := sender.Sender{}
	ext := domain.PriceExtractor{}
	val := domain.MessageHourVal{}

	return app.NewGoldPriceService(req, sen, ext, val)
}
