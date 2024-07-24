package command

import (
	"PriceWatcher/internal/app/bot/command/price"
	"PriceWatcher/internal/entities/telebot"
)

func CreateCommands(pCom price.CurrentPriceComm, subCom price.SubscribingComm) []telebot.Command {
	return []telebot.Command{
		{
			Name:        "/hello",
			Description: "Say hello to the bot",
			Action: func(interface{}) string {
				return "Hello there!"
			},
		}, {
			Name:        "/prices",
			Description: "Get the currents prices",
			Action:      pCom.GetCurrentPrices,
		}, {
			Name:        "/subscribe",
			Description: "Subscribe to messages of the current gold price ",
			Action:      subCom.SubscribeUser,
		},
	}
}
