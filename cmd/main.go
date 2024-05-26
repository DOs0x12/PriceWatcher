package main

import (
	"PriceWatcher/internal"
	"PriceWatcher/internal/bank"
	"PriceWatcher/internal/common/interruption"
	"PriceWatcher/internal/config"
	"PriceWatcher/internal/telebot"
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
	conf, err := configer.GetConfig()
	if err != nil {
		logrus.Error("%w", err)
	}
	jobDone := make(chan interface{})
	bankService := bank.NewService(bank.BankRequester{}, bank.NewPriceExtractor(`([0-9]).*([0-9])*,([0-9])*`, "div"), conf)

	startBot(botCtx, wg, configer, jobDone)
	startWatching(watcherCtx, wg, bankService, jobDone)

	wg.Wait()

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context,
	wg *sync.WaitGroup,
	bankService bank.Service,
	jobDone chan<- interface{}) {
	internal.ServeMetalPrice(ctx, wg, bankService, jobDone)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(ctx context.Context,
	wg *sync.WaitGroup,
	configer config.Configer,
	jobDone chan<- interface{}) {
	bot, err := telebot.NewTelebot(configer)
	if err != nil {
		logrus.Errorf("bot: %v", err)
		wg.Done()

		return
	}

	err = telebot.Start(ctx, wg, bot, configer, jobDone)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}

func GetConfiger() config.Configer {
	configPath := "config.yml"

	return config.NewConfiger(configPath)
}
