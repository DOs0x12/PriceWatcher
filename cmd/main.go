package main

import (
	"PriceWatcher/internal/app/interruption"
	"PriceWatcher/internal/app/telebot"
	"PriceWatcher/internal/app/watcher"
	"PriceWatcher/internal/infrastructure/configer"
	"PriceWatcher/internal/infrastructure/file"
	"PriceWatcher/internal/infrastructure/sender"
	infraTelebot "PriceWatcher/internal/infrastructure/telebot"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func main() {
	watcherCtx, watcherCancel := newContext()
	botCtx, botCancel := newContext()

	logrus.Infoln("Start the application")

	interruption.WatchForInterruption(watcherCancel, botCancel)

	servCount := 2
	wg := &sync.WaitGroup{}
	wg.Add(servCount)

	configer := GetConfiger()

	wr := file.NewWR()

	startBot(botCtx, wg, configer, wr)
	startWatching(watcherCtx, wg, configer, wr)

	wg.Wait()

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context, wg *sync.WaitGroup, configer configer.Configer, wr file.WriteReader) {
	sen := sender.Sender{}

	watcher.ServeWatchers(ctx, wg, configer, sen, wr)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(ctx context.Context, wg *sync.WaitGroup, configer configer.Configer, wr file.WriteReader) {
	bot, err := infraTelebot.NewTelebot(configer)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}

	err = telebot.Start(ctx, wg, bot, wr)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}

func GetConfiger() configer.Configer {
	configPath := "config.yml"

	return configer.NewConfiger(configPath)
}
