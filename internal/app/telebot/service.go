package telebot

import (
	botCom "PriceWatcher/internal/app/telebot/command"
	"PriceWatcher/internal/app/telebot/command/item"
	"PriceWatcher/internal/app/telebot/command/price"
	botEnt "PriceWatcher/internal/entities/telebot"
	infraFile "PriceWatcher/internal/infrastructure/file"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/telebot"
	"context"
	"fmt"
	"sync"
)

func Start(ctx context.Context, wg *sync.WaitGroup, bot telebot.Bot, wr infraFile.WriteReader, configer configer.Configer) error {
	defer wg.Done()

	commands := createCommands(wr)
	commandsWithInput := createCommandsWithInput(configer)
	if err := bot.Start(commands, commandsWithInput); err != nil {
		return fmt.Errorf("can not start the bot: %v", err)
	}

	if err := bot.RegisterCommands(commands, commandsWithInput); err != nil {
		return fmt.Errorf("can not register commands in the bot: %v", err)
	}

	go func() {
		<-ctx.Done()
		bot.Stop()
	}()

	return nil
}

func createCommands(wr infraFile.WriteReader) []botEnt.Command {
	pCom := price.NewPriceCommand(wr)
	commands := botCom.CreateCommands(pCom)

	botComms := make([]botEnt.Command, len(commands))

	for i, command := range commands {
		botCommand := botEnt.Command{
			Name:        command.Name,
			Description: command.Description,
			Action:      command.Action,
		}

		botComms[i] = botCommand
	}

	return botComms
}

func createCommandsWithInput(configer configer.Configer) []botEnt.CommandWithInput {
	addCom := item.NewAddItemComm(configer)
	commands := botCom.CreateCommandsWithInput(addCom)

	botComms := make([]botEnt.CommandWithInput, len(commands))

	for i, command := range commands {
		botCommand := botEnt.CommandWithInput{
			Name:        command.Name,
			Description: command.Description,
			Action:      command.Action,
		}

		botComms[i] = botCommand
	}

	return botComms
}
