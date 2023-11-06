package main

import (
	"PriceWatcher/internal/app"
	"PriceWatcher/internal/app/interruption"
	"PriceWatcher/internal/app/service"
	"PriceWatcher/internal/infrastructure/configer"
	"PriceWatcher/internal/infrastructure/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := newContext()

	serv, err := newService()
	if err != nil {
		logrus.Errorf("Got the error: %v", err)
	}

	logrus.Infoln("Start the application")

	interruption.WatchForInterruption(cancel)
	app.Watch(ctx.Done(), serv)

	logrus.Infoln("The application is done")
}

func newService() (service.PriceWatcherService, error) {
	sen := sender.Sender{}
	conf := configer.Configer{}

	return service.NewWatcherService(sen, conf)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}
