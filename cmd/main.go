package main

import (
	internalApp "PriceWatcher/internal/app"
	bankApp "PriceWatcher/internal/app/bank"
	"PriceWatcher/internal/app/bank/interruption"
	appBotComm "PriceWatcher/internal/app/bot/command"
	appPrice "PriceWatcher/internal/app/bot/command/price"
	bankDom "PriceWatcher/internal/domain/bank"
	subEnt "PriceWatcher/internal/entities/subscribing"
	botEnt "PriceWatcher/internal/entities/telebot"
	bankInfra "PriceWatcher/internal/infrastructure/bank"
	botInfra "PriceWatcher/internal/infrastructure/bot"
	"PriceWatcher/internal/infrastructure/config"
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

	configer := NewConfiger()
	conf, err := configer.GetConfig()
	if err != nil {
		logrus.Error("%w", err)
	}

	bankService := bankApp.NewService(bankInfra.BankRequester{},
		bankDom.NewPriceExtractor(`([0-9]).*([0-9])*,([0-9])*`, "div"), conf)

	subService := bankInfra.SubscribingService{}
	subscribers, err := subService.GetSubscribers()
	if err != nil {
		logrus.Errorf("cannot get subscribers: %v", err)
		wg.Done()

		return
	}

	commands := createCommands(subscribers)
	bot, err := botInfra.NewTelebot(configer, commands)
	if err != nil {
		logrus.Errorf("bot: %v", err)
		wg.Done()

		return
	}

	startBot(botCtx, bot)
	startWatching(watcherCtx, wg, bankService, bot, subscribers)

	wg.Wait()

	err = subService.SaveSubscribers(subscribers)
	if err != nil {
		logrus.Errorf("cannot save the data of subscribers: %v", err)
	}

	logrus.Infoln("The application is done")
}

func startWatching(ctx context.Context,
	wg *sync.WaitGroup,
	bankService bankApp.Service,
	bot botInfra.Telebot,
	subscribers *subEnt.Subscribers) {
	internalApp.ServeMetalPrice(ctx, wg, bankService, bot, subscribers)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func startBot(ctx context.Context, bot botInfra.Telebot) {

	err := bot.Start(ctx)
	if err != nil {
		logrus.Errorf("bot: %v", err)

		return
	}
}

func NewConfiger() config.Configer {
	configPath := "gold-price-watcher-data/config.yml"

	return config.NewConfiger(configPath)
}

func createCommands(subscribers *subEnt.Subscribers) []botEnt.Command {
	subCom := appPrice.SubscribingComm{Subscribers: subscribers}

	return appBotComm.CreateCommands(subCom)
}
