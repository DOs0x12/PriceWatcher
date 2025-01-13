package telebot

import (
	"PriceWatcher/internal/entities/bot"
	"PriceWatcher/internal/interfaces/broker"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context,
	wg *sync.WaitGroup,
	broker broker.Worker,
	serviceName string,
	commands []bot.Command) error {
	msgChan, err := broker.Start(ctx, serviceName)
	if err != nil {
		return err
	}

	go processeMessages(ctx, broker, msgChan, commands)

	go func() {
		<-ctx.Done()
		broker.Stop()
		wg.Done()
	}()

	return nil
}

func processeMessages(ctx context.Context,
	broker broker.Worker,
	msgChan <-chan bot.Message,
	commands []bot.Command) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			for _, cmd := range commands {
				if cmd.Name == msg.Command {
					result := cmd.Action(msg)
					err := broker.SendMessage(ctx, result, msg.ChatID)
					if err != nil {
						logrus.Error("Can not send a message:", err)
					}
				}
			}

			err := broker.CommitMessage(ctx, msg.MsgUuid)
			if err != nil {
				logrus.Error("Can not commit a message:", err)
			}
		}
	}
}
