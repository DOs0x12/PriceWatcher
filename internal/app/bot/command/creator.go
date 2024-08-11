package command

import (
	subComms "PriceWatcher/internal/app/bot/command/subscribing"
	"PriceWatcher/internal/entities/bank/subscribing"
	"PriceWatcher/internal/entities/bot"
	"sync"
)

func CreateHelloCommand() bot.Command {
	return bot.Command{
		Name:        "/hello",
		Description: "Say hello to the bot",
		Action: func(interface{}) string {
			return "Hello there!"
		},
	}
}

func CreateSubCommand(mu *sync.Mutex, subs *subscribing.Subscribers) bot.Command {
	subCom := subComms.NewSubCommand(mu, subs)
	return bot.Command{
		Name:        "/start",
		Description: "Start getting messages of the current gold price ",
		Action:      subCom.SubscribeUser,
	}
}

func CreateUnsubCommand(mu *sync.Mutex, subs *subscribing.Subscribers) bot.Command {
	unsubCom := subComms.NewUnsubCommand(mu, subs)
	return bot.Command{
		Name:        "/stop",
		Description: "Stop getting notifications about the current gold price ",
		Action:      unsubCom.UnsubscribeUser,
	}
}
