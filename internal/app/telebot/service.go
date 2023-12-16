package telebot

import (
	botEnt "PriceWatcher/internal/entities/telebot"
	"PriceWatcher/internal/interfaces/telebot"
	"context"
	"fmt"
	"sync"
)

func Start(wg *sync.WaitGroup, ctx context.Context, bot telebot.Bot) error {
	commands := createCommands()
	if err := bot.Start(commands...); err != nil {
		wg.Done()

		return fmt.Errorf("can not start the bot: %v", err)
	}

	if err := bot.RegisterCommands(commands); err != nil {
		wg.Done()

		return fmt.Errorf("can not register commands in the bot: %v", err)
	}

	go func(wg *sync.WaitGroup) {
		<-ctx.Done()
		bot.Stop()
		wg.Done()
	}(wg)

	return nil
}

func createCommands() []botEnt.Command {
	botComms := make([]botEnt.Command, len(commands))

	for i, command := range commands {
		botCommand := botEnt.Command{
			Name:        command.name,
			Description: command.description,
			Action:      command.action,
		}

		botComms[i] = botCommand
	}

	return botComms
}
