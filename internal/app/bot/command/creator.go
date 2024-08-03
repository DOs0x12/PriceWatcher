package command

import (
	subComm "PriceWatcher/internal/app/bot/command/subscribing"
	unsubComm "PriceWatcher/internal/app/bot/command/unsubscribing"
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
	subCom := subComm.SubscribingComm{Subscribers: subs}
	return telebot.Command{
		Name:        "/subscribe",
		Description: "Subscribe to messages of the current price ",
		Action:      subCom.SubscribeUser,
	}
}

func CreateUnsubCommand(subs *subscribing.Subscribers) telebot.Command {
	unsubCom := unsubComm.UnsubscribingComm{Subscribers: subs}
	return telebot.Command{
		Name:        "/subscribe",
		Description: "Subscribe to messages of the current price ",
		Action:      unsubCom.UnsubscribeUser,
	}
}
