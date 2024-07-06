package command

import "PriceWatcher/internal/telebot/command/price"

type Command struct {
	Name        string
	Description string
	Action      func(interface{}) string
}

func CreateCommands(pCom price.CurrentPriceComm, subCom price.SubscribingComm) []Command {
	return []Command{
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
