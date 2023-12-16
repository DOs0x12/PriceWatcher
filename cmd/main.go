package main

import (
	"PriceWatcher/internal/app/interruption"
	"PriceWatcher/internal/app/telebot"
	"PriceWatcher/internal/app/watcher"
	"PriceWatcher/internal/infrastructure/configer"
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

	startBot(wg, botCtx)
	startWatching(wg, watcherCtx)

	wg.Wait()

	logrus.Infoln("The application is done")
}

func startWatching(wg *sync.WaitGroup, ctx context.Context) {
	sen := sender.Sender{}

	configPath := "config.yml"
	conf := configer.NewConfiger(configPath)

	watcher.ServeWatchers(wg, ctx, conf, sen)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(wg *sync.WaitGroup, ctx context.Context) {
	bot, err := infraTelebot.NewTelebot("6892592660:AAEf69s7JICdEKVTCboSGBeRC43HELUcfiY")
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}

	err = telebot.Start(wg, ctx, bot)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}
