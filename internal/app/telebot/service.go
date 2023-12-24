package telebot

import (
	botCom "PriceWatcher/internal/app/telebot/command"
	"PriceWatcher/internal/app/telebot/command/price"
	botEnt "PriceWatcher/internal/entities/telebot"
	infraFile "PriceWatcher/internal/infrastructure/file"
	"PriceWatcher/internal/interfaces/telebot"
	"context"
	"fmt"
	"sync"
)

func Start(ctx context.Context, wg *sync.WaitGroup, bot telebot.Bot, wr infraFile.WriteReader) error {
	defer wg.Done()

	commands := createCommands(wr)
	if err := bot.Start(commands...); err != nil {
		return fmt.Errorf("can not start the bot: %v", err)
	}

	if err := bot.RegisterCommands(commands); err != nil {
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
