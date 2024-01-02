package command

import (
	"PriceWatcher/internal/app/telebot/command/price"
)

type Command struct {
	Name        string
	Description string
	Action      func() string
}

func CreateCommands(pCom price.PriceCommand) []Command {
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
