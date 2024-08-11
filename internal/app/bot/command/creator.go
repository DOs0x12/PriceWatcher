package command

import (
	subComms "PriceWatcher/internal/app/bot/command/subscribing"
	"PriceWatcher/internal/entities/subscribing"
	"PriceWatcher/internal/entities/telebot"
	"sync"
)

func CreateHelloCommand() telebot.Command {
	return telebot.Command{
		Name:        "/hello",
		Description: "Say hello to the bot",
		Action: func(interface{}) string {
			return "Hello there!"
		},
	}
}

func CreateSubCommand(mu *sync.Mutex, subs *subscribing.Subscribers) telebot.Command {
	subCom := subComms.NewSubCommand(mu, subs)
	return telebot.Command{
		Name:        "/start",
		Description: "Start getting messages of the current gold price ",
		Action:      subCom.SubscribeUser,
	}
}

func CreateUnsubCommand(mu *sync.Mutex, subs *subscribing.Subscribers) telebot.Command {
	unsubCom := subComms.NewUnsubCommand(mu, subs)
	return telebot.Command{
		Name:        "/stop",
		Description: "Stop getting notifications about the current gold price ",
		Action:      unsubCom.UnsubscribeUser,
	}
}
