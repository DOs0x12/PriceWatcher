package main

import (
	"PriceWatcher/internal/app"
	"PriceWatcher/internal/app/interruption"
	"PriceWatcher/internal/infrastructure/configer"
	"PriceWatcher/internal/infrastructure/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := newContext()

	logrus.Infoln("Start the application")

	interruption.WatchForInterruption(cancel)
	startWatching(ctx)

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context) {
	sen := sender.Sender{}

	configPath := "config.yml"
	conf := configer.NewConfiger(configPath)

	app.StartWatchers(ctx, conf, sen)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}
