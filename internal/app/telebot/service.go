package telebot

import (
	entTelebot "PriceWatcher/internal/entities/telebot"
	"PriceWatcher/internal/interfaces/telebot"
	"context"
	"fmt"
)

func Start(ctx context.Context, bot telebot.Bot) (<-chan string, error) {
	commands := createCommands()
	if err := bot.Start(commands...); err != nil {
		return nil, fmt.Errorf("can not start the bot: %v", err)
	}

	if err := bot.RegisterCommands(commands); err != nil {
		return nil, fmt.Errorf("can not register commands in the bot: %v", err)
	}

	ch := make(chan string)

	go func() {
		<-ctx.Done()
		bot.Stop()
		ch <- "bot"
	}()

	return ch, nil
}

func createCommands() []entTelebot.Command {
	botComms := make([]entTelebot.Command, len(commands))

	for i, command := range commands {
		botCommand := entTelebot.Command{
			Name:        command.name,
			Description: command.description,
			Action:      command.action,
		}

		botComms[i] = botCommand
	}

	return botComms
}
