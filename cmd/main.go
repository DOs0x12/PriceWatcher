package main

import (
	"PriceWatcher/internal"
	"PriceWatcher/internal/bank"
	"PriceWatcher/internal/common/interruption"
	"PriceWatcher/internal/config"
	"PriceWatcher/internal/entities/subscribing"
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

	bot, err := telebot.NewTelebot(configer)
	if err != nil {
		logrus.Errorf("bot: %v", err)
		wg.Done()

		return
	}

	subscribers := &subscribing.Subscribers{ChatIDs: make([]int64, 0)}

	startBot(botCtx, wg, bot, configer, jobDone, subscribers)
	startWatching(watcherCtx, wg, bankService, jobDone, bot, subscribers)

	wg.Wait()

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context,
	wg *sync.WaitGroup,
	bankService bank.Service,
	jobDone chan<- interface{},
	bot telebot.Telebot,
	subscribers *subscribing.Subscribers) {
	internal.ServeMetalPrice(ctx, wg, bankService, jobDone, bot, subscribers)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(ctx context.Context,
	wg *sync.WaitGroup,
	bot telebot.Telebot,
	configer config.Configer,
	jobDone chan<- interface{},
	subscribers *subscribing.Subscribers) {

	err := telebot.Start(ctx, wg, bot, configer, jobDone, subscribers)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}

func GetConfiger() config.Configer {
	configPath := "config.yml"

	return config.NewConfiger(configPath)
}
