package command

import (
	subComms "PriceWatcher/internal/app/bot/command/subscribing"
	"PriceWatcher/internal/entities/subscribing"
	"PriceWatcher/internal/entities/telebot"
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

func CreateSubCommand(subs *subscribing.Subscribers) telebot.Command {
	subCom := subComms.SubscribingComm{Subscribers: subs}
	return telebot.Command{
		Name:        "/subscribe",
		Description: "Subscribe to messages of the current price ",
		Action:      subCom.SubscribeUser,
	}
}

func CreateUnsubCommand(subs *subscribing.Subscribers) telebot.Command {
	unsubCom := subComms.UnsubscribingComm{Subscribers: subs}
	return telebot.Command{
		Name:        "/subscribe",
		Description: "Subscribe to messages of the current price ",
		Action:      unsubCom.UnsubscribeUser,
	}
}
