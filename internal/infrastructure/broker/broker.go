package bot

import (
	"PriceWatcher/internal/entities/bot"
	"context"
	"fmt"
	"sync"
	"time"

	extBroker "github.com/DOs0x12/TeleBot/client/broker"
	"github.com/google/uuid"
)

type Broker struct {
	mu        *sync.RWMutex
	commands  []bot.Command
	sender    extBroker.Sender
	receivers map[uuid.UUID]extBroker.Receiver
	address   string
}

func NewBroker(commands []bot.Command, address string) Broker {
	sender := extBroker.NewSender(address)
	return Broker{
		mu:        &sync.RWMutex{},
		commands:  commands,
		sender:    sender,
		address:   address,
		receivers: make(map[uuid.UUID]extBroker.Receiver),
	}
}

func (b Broker) Start(ctx context.Context) (<-chan bot.Message, error) {
	dataChan := make(chan bot.Message)
	for _, comm := range b.commands {
		commData := extBroker.CommandData{Name: comm.Name, Description: comm.Description}

		if err := b.sender.RegisterCommand(ctx, commData); err != nil {
			return nil, fmt.Errorf("cannot start the broker: %v", err)
		}

		r := extBroker.NewReceiver(b.address, comm.Name)
		recUuid := uuid.New()
		b.mu.Lock()
		b.receivers[recUuid] = *r
		b.mu.Unlock()
		brokerDataChan := r.StartGetData(ctx)

		go pipelineData(ctx, brokerDataChan, dataChan, comm.Name, recUuid)
	}

	return dataChan, nil
}

func pipelineData(ctx context.Context,
	brokerDataChan <-chan extBroker.BotData,
	msgChan chan<- bot.Message,
	command string,
	recUuid uuid.UUID) {
	for {
		select {
		case <-ctx.Done():
			return
		case brokerData := <-brokerDataChan:
			msg := bot.Message{
				ChatID:       brokerData.ChatID,
				Command:      command,
				Value:        brokerData.Value,
				MsgUuid:      brokerData.MessageUuid,
				ReceiverUuid: recUuid,
			}

			msgChan <- msg
		}
	}

}

func (b Broker) Stop() error {
	if err := b.sender.Stop(); err != nil {
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
		if err = b.sender.SendData(ctx, botData); err != nil {
			time.Sleep(5 * time.Second)
			cnt++

			continue
		}

		return nil
	}

	return fmt.Errorf("an error occurs at sending a message to the broker: %v", err)
}

func (b Broker) CommitMessage(ctx context.Context, recUuid uuid.UUID, msgUuid uuid.UUID) error {
	b.mu.RLock()
	r := b.receivers[recUuid]
	b.mu.RUnlock()
	err := r.Commit(ctx, msgUuid)
	if err != nil {
		return fmt.Errorf("an error occurs at commiting a message in the broker: %v", err)
	}

	return nil
}
