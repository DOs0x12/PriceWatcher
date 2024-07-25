package main

import (
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
	appCtx, appCancel := context.WithCancel(context.Background())

	logrus.Infoln("Start the application")

	interruption.WatchForInterruption(appCancel)

	servCount := 2
	wg := &sync.WaitGroup{}
	wg.Add(servCount)

	configer := NewConfiger()
	conf, err := configer.GetConfig()
	if err != nil {
		logrus.Error("Cannot get the config: %w", err)
	}

	priceRegEx := `([0-9]).*([0-9])*,([0-9])*`
	priceTag := "div"
	priceExtractor := bankDom.NewPriceExtractor(priceRegEx, priceTag)
	bankService := bankApp.NewService(bankInfra.BankRequester{}, priceExtractor, conf)

	subService := bankInfra.SubscribingService{}
	subscribers, err := subService.GetSubscribers()
	if err != nil {
		logrus.Errorf("Cannot get subscribers: %v", err)

		return
	}

	commands := createCommands(subscribers)
	bot, err := botInfra.NewTelebot(wg, configer, commands)
	if err != nil {
		logrus.Errorf("A bot error occurs: %v", err)

		return
	}

	err = bot.Start(appCtx)
	if err != nil {
		logrus.Errorf("Cannot start a bot: %v", err)

		return
	}

	bankService.WatchPrice(appCtx, wg, bot, subscribers)

	wg.Wait()

	err = subService.SaveSubscribers(subscribers)
	if err != nil {
		logrus.Errorf("Cannot save the data of subscribers: %v", err)
	}

	logrus.Infoln("The application is done")
}

func NewConfiger() config.Configer {
	configPath := "gold-price-watcher-data/config.yml"

	return config.NewConfiger(configPath)
}

func createCommands(subscribers *subEnt.Subscribers) []botEnt.Command {
	subCom := appPrice.SubscribingComm{Subscribers: subscribers}

	return appBotComm.CreateCommands(subCom)
}
