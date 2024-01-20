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
	restart := make(chan interface{})

	startBot(botCtx, wg, configer, wr, restart)
	startWatching(watcherCtx, wg, configer, wr, restart)

	wg.Wait()

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context,
	wg *sync.WaitGroup,
	configer configer.Configer,
	wr file.WriteReader,
	restart <-chan interface{}) {
	sen := sender.Sender{}

	watcher.ServeWatchers(ctx, wg, configer, sen, wr, restart)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(ctx context.Context,
	wg *sync.WaitGroup,
	configer configer.Configer,
	wr file.WriteReader,
	restart chan<- interface{}) {
	bot, err := infraTelebot.NewTelebot(configer)
	if err != nil {
		logrus.Errorf("bot: %v", err)
		wg.Done()

		return
	}

	err = telebot.Start(ctx, wg, bot, wr, configer, restart)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}

func GetConfiger() configer.Configer {
	configPath := "config.yml"

	return configer.NewConfiger(configPath)
}
