package bot

import (
	"PriceWatcher/internal/entities/bot"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DOs0x12/TeleBot/client/v2/broker"
	"github.com/google/uuid"
)

type Broker struct {
	mu          *sync.RWMutex
	commands    []bot.Command
	address     string
	kafkaBroker broker.KafkaBroker
}

func NewBroker(ctx context.Context, address, serviceName string, commands []bot.Command) (Broker, error) {
	kafkaBroker, err := broker.NewKafkaBroker(ctx, address, serviceName)
	if err != nil {
		return Broker{}, fmt.Errorf("failed to create a broker: %w", err)
	}

	return Broker{
		mu:          &sync.RWMutex{},
		kafkaBroker: *kafkaBroker,
		address:     address,
		commands:    commands,
	}, nil
}

func (b Broker) Start(ctx context.Context, serviceName string) (<-chan bot.Message, error) {
	for _, comm := range b.commands {
		commData := broker.BrokerCommandData{Name: comm.Name, Description: comm.Description}
		err := b.kafkaBroker.RegisterCommand(ctx, commData, serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to register the command %v: %w", comm.Name, err)
		}
	}

	kafkaBrChan := b.kafkaBroker.StartGetData(ctx)
	msgChan := make(chan bot.Message)

	go pipelineData(ctx, kafkaBrChan, msgChan)

	return msgChan, nil
}

func pipelineData(ctx context.Context,
	brokerDataChan <-chan broker.BrokerData,
	msgChan chan<- bot.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case brokerData := <-brokerDataChan:
			msg := bot.Message{
				ChatID:  brokerData.ChatID,
				Value:   brokerData.Value,
				MsgUuid: brokerData.MessageUuid,
			}

			msgChan <- msg
		}
	}

}

func (b Broker) Stop() {
	b.kafkaBroker.Stop()
}

func (b Broker) SendMessage(ctx context.Context, msg string, chatID int64) error {
	botData := broker.BrokerData{ChatID: chatID, Value: msg}
	maxRetries := 10
	cnt := 0
	var err error

	for cnt < maxRetries {
		if err = b.kafkaBroker.SendData(ctx, botData); err != nil {
			time.Sleep(5 * time.Second)
			cnt++

			continue
		}

		return nil
	}

	return fmt.Errorf("an error occurs at sending a message to the broker: %v", err)
}

func (b Broker) CommitMessage(ctx context.Context, msgUuid uuid.UUID) error {
	err := b.kafkaBroker.Commit(ctx, msgUuid)
	if err != nil {
		return fmt.Errorf("an error occurs at commiting a message in the broker: %v", err)
	}

	return nil
}
