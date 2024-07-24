package command

import (
	"PriceWatcher/internal/app/bot/command/price"
	"PriceWatcher/internal/entities/telebot"
)

func CreateCommands(subCom price.SubscribingComm) []telebot.Command {
	return []telebot.Command{
		{
			Name:        "/hello",
			Description: "Say hello to the bot",
			Action: func(interface{}) string {
				return "Hello there!"
			},
		}, {
			Name:        "/subscribe",
			Description: "Subscribe to messages of the current gold price ",
			Action:      subCom.SubscribeUser,
		},
	}
}
