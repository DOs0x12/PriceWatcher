package command

import (
	"PriceWatcher/internal/telebot/command/price"
)

type Command struct {
	Name        string
	Description string
	Action      func() string
}

func CreateCommands(pCom price.CurrentPriceComm) []Command {
	return []Command{
		{
			Name:        "/hello",
			Description: "Say hello to the bot",
			Action: func() string {
				return "Hello there!"
			},
		}, {
			Name:        "/prices",
			Description: "Get the currents prices",
			Action:      pCom.GetCurrentPrices,
		},
	}
}
