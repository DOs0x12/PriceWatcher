package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	serv := newService()

	logrus.Infoln("Start the application")

	serv.Watch(ctx.Done(), cancel, realClock{})

	logrus.Infoln("The application is done")
}

func newService() *app.GoldPriceService {
	req := requester.Requester{}
	sen := sender.Sender{}
	ext := domain.PriceExtractor{}
	val := domain.MessageHourVal{}

	return app.NewGoldPriceService(req, sen, ext, val)
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
