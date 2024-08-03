package main

import (
	bankApp "PriceWatcher/internal/app/bank"
	appBotComm "PriceWatcher/internal/app/bot/command"
	"PriceWatcher/internal/app/interruption"
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

const configPath = "price-watcher-data/config.yml"
const subscribersFilePath = "price-watcher-data/subscribers.yml"

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())

	logrus.Infoln("Start the application")

	interruption.WatchForInterruption(appCancel)

	servCount := 2
	wg := &sync.WaitGroup{}
	wg.Add(servCount)

	configer := config.NewConfiger(configPath)
	conf, err := configer.GetConfig()
	if err != nil {
		logrus.Errorf("Cannot get the config: %v", err)

		return
	}

	priceRegEx := `([0-9]).*([0-9])*,([0-9])*`
	priceTag := "div"
	priceExtractor := bankDom.NewPriceExtractor(priceRegEx, priceTag)
	bankService := bankApp.NewService(bankInfra.BankRequester{}, priceExtractor, conf)

	subService := bankInfra.SubscribingService{}
	subscribers, err := subService.GetSubscribers(subscribersFilePath)
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

	err = subService.SaveSubscribers(subscribers, subscribersFilePath)
	if err != nil {
		logrus.Errorf("Cannot save the data of subscribers: %v", err)
	}

	logrus.Infoln("The application is done")
}

func createCommands(subscribers *subEnt.Subscribers) []botEnt.Command {
	return []botEnt.Command{
		appBotComm.CreateHelloCommand(),
		appBotComm.CreateSubCommand(subscribers),
	}
}
