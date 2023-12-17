package telebot

import (
	botEnt "PriceWatcher/internal/entities/telebot"
	"PriceWatcher/internal/interfaces/telebot"
	"context"
	"fmt"
	"sync"
)

func Start(ctx context.Context, wg *sync.WaitGroup, bot telebot.Bot) error {
	defer wg.Done()

	commands := createCommands()
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
