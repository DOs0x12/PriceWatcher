package command

import (
	"PriceWatcher/internal/app/telebot/command/item"
	"PriceWatcher/internal/app/telebot/command/price"
)

type Command struct {
	Name        string
	Description string
	Action      func() string
}

type CommandWithInput struct {
	Name        string
	Description string
	Action      func(input string) string
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

func CreateCommandsWithInput(addCom item.AddItemComm, remCom item.RemoveItemComm) []CommandWithInput {
	return []CommandWithInput{
		{
			Name:        "/additem",
			Description: "Add an item for watching its price",
			Action:      addCom.AddItemToWatch,
		}, {
			Name:        "/removeitem",
			Description: "Remove an item from watching its price",
			Action:      remCom.RemoveItemToWatch,
		},
	}
}
