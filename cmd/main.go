package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	startInterruptionWatch(cancel)
	serv := newService()

	logrus.Infoln("Start the application")

	app.WatchGoldPrice(serv, ctx.Done())

	logrus.Infoln("The application is done")
}

func newService() *app.GoldPriceService {
	req := requester.Requester{}
	ext := domain.PriceExtractor{}
	sen := sender.Sender{}

	return app.NewGoldPriceService(req, ext, sen)
}

func startInterruptionWatch(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logrus.Infof("The application has been stopped")
		cancel()
	}()
}
