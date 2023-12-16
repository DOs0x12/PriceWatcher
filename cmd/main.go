package main

import (
	"PriceWatcher/internal/app/interruption"
	"PriceWatcher/internal/app/telebot"
	"PriceWatcher/internal/app/watcher"
	"PriceWatcher/internal/infrastructure/configer"
	"PriceWatcher/internal/infrastructure/sender"
	infraTelebot "PriceWatcher/internal/infrastructure/telebot"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	watcherCtx, watcherCancel := newContext()
	botCtx, botCancel := newContext()

	defer botCancel()

	logrus.Infoln("Start the application")

	interruption.WatchForInterruption(watcherCancel)
	startBot(botCtx)
	startWatching(watcherCtx)

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context) {
	sen := sender.Sender{}

	configPath := "config.yml"
	conf := configer.NewConfiger(configPath)

	watcher.ServeWatchers(ctx, conf, sen)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(ctx context.Context) {
	bot, err := infraTelebot.NewTelebot("6892592660:AAEf69s7JICdEKVTCboSGBeRC43HELUcfiY")
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}

	_, err = telebot.Start(ctx, bot)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}
