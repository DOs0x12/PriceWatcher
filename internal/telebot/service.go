package telebot

import (
	"PriceWatcher/internal/config"
	"PriceWatcher/internal/entities/subscribing"
	botEnt "PriceWatcher/internal/entities/telebot"
	botCom "PriceWatcher/internal/telebot/command"
	"PriceWatcher/internal/telebot/command/price"
	"context"
	"fmt"
	"sync"
)

func Start(ctx context.Context,
	wg *sync.WaitGroup,
	bot Telebot,
	configer config.Configer,
	restart chan<- interface{},
	subscribers *subscribing.Subscribers) error {
	defer wg.Done()

	commands := createCommands(subscribers)
	if err := bot.Start(commands, restart); err != nil {
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

func createCommands(subscribers *subscribing.Subscribers) []botEnt.Command {
	pCom := price.NewPriceCommand()
	subCom := price.SubscribingComm{Subscribers: subscribers}
	commands := botCom.CreateCommands(pCom, subCom)

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
