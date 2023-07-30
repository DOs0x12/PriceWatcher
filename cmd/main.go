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

	watchInterruption(cancel)
	serv := newService()

	logrus.Infoln("Start the application")

	serv.Watch(ctx.Done())

	logrus.Infoln("The application is done")
}

func newService() *app.GoldPriceService {
	req := requester.Requester{}
	sen := sender.Sender{}
	ext := domain.PriceExtractor{}
	val := domain.MessageHourVal{}

	return app.NewGoldPriceService(req, sen, ext, val)
}

func watchInterruption(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logrus.Infof("The application has been stopped")
		cancel()
	}()
}
