package bot

import (
	"PriceWatcher/internal/entities/bot"
	"context"
	"fmt"
	"strings"
	"time"

	extBroker "github.com/Guise322/TeleBot/client/pkg/broker"
)

type Broker struct {
	commands  []bot.Command
	extBroker extBroker.Broker
	address   string
}

func NewBroker(commands []bot.Command, address string) Broker {
	extBroker := extBroker.NewBroker(address)
	return Broker{commands: commands, extBroker: extBroker, address: address}
}

func (b Broker) Start(ctx context.Context) (<-chan bot.Message, error) {
	dataChan := make(chan bot.Message)
	for _, comm := range b.commands {
		commData := extBroker.CommandData{Name: comm.Name, Description: comm.Description}

		if err := b.extBroker.RegisterCommand(ctx, commData); err != nil {
			return nil, fmt.Errorf("cannot start the broker: %v", err)
		}

		topicName := strings.Trim(comm.Name, "/")

		brokerDataChan := b.extBroker.StartGetData(ctx, topicName, b.address)
		go pipelineData(ctx, brokerDataChan, dataChan, comm.Name)
	}

	return dataChan, nil
}

func pipelineData(ctx context.Context,
	brokerDataChan <-chan extBroker.BotData,
	msgChan chan<- bot.Message,
	command string) {
	for {
		select {
		case <-ctx.Done():
			return
		case brokerData := <-brokerDataChan:
			msg := bot.Message{ChatID: brokerData.ChatID, Command: command, Value: brokerData.Value}
			msgChan <- msg
		}
	}

}

func (b Broker) Stop() error {
	if err := b.extBroker.Stop(); err != nil {
		return fmt.Errorf("an error occurs at stopping the broker worker: %v", err)
	}

	return nil
}

func (b Broker) SendMessage(ctx context.Context, msg string, chatID int64) error {
	botData := extBroker.BotData{ChatID: chatID, Value: msg}
	maxRetries := 10
	cnt := 0
	var err error

	for cnt < maxRetries {
		if err = b.extBroker.SendData(ctx, botData); err != nil {
			time.Sleep(5 * time.Second)
			cnt++

			continue
		}

		return nil
	}

	return fmt.Errorf("an error occurs at sending a message to the broker: %v", err)
}
