package main

import (
	bankApp "PriceWatcher/internal/app/bank"
	botApp "PriceWatcher/internal/app/bot"
	appBotComm "PriceWatcher/internal/app/bot/command"
	"PriceWatcher/internal/app/interruption"
	bankDom "PriceWatcher/internal/domain/bank"
	botEnt "PriceWatcher/internal/entities/bot"
	subEnt "PriceWatcher/internal/entities/subscribing"
	infraBank "PriceWatcher/internal/infrastructure/bank"
	brokerInfra "PriceWatcher/internal/infrastructure/broker"
	"PriceWatcher/internal/infrastructure/config"
	infraSub "PriceWatcher/internal/infrastructure/subscribing"
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
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
	bankService := bankApp.NewService(infraBank.BankRequester{}, priceExtractor, conf)

	subService := infraSub.SubscribingService{}
	subscribers, err := subService.GetSubscribers(subscribersFilePath)
	if err != nil {
		logrus.Errorf("Cannot get subscribers: %v", err)

		return
	}

	commands := createCommands(subscribers)
	w := &kafka.Writer{
		Addr:     kafka.TCP(conf.KafkaAddress),
		Balancer: &kafka.LeastBytes{},
	}
	broker := brokerInfra.NewBroker(commands, w)

	err = botApp.Start(appCtx, wg, broker, commands)
	if err != nil {
		logrus.Errorf("Cannot start serving bot messages: %v", err)

		return
	}

	bankService.WatchPrice(appCtx, wg, broker, subscribers)

	wg.Wait()

	err = subService.SaveSubscribers(subscribers, subscribersFilePath)
	if err != nil {
		logrus.Errorf("Cannot save the data of subscribers: %v", err)
	}

	logrus.Infoln("The application is done")
}

func createCommands(subscribers *subEnt.Subscribers) []botEnt.Command {
	mu := &sync.Mutex{}
	return []botEnt.Command{
		appBotComm.CreateHelloCommand(),
		appBotComm.CreateSubCommand(mu, subscribers),
		appBotComm.CreateUnsubCommand(mu, subscribers),
	}
}
