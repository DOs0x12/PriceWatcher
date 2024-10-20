package bot

import (
	"PriceWatcher/internal/entities/bot"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Guise322/TeleBot/client/pkg/service"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	commands []bot.Command
	w        *kafka.Writer
}

func NewBroker(commands []bot.Command, w *kafka.Writer) Broker {
	return Broker{commands: commands, w: w}
}

func (b Broker) Start(ctx context.Context) (chan<- bot.Message, error) {
	dataChan := make(chan bot.Message)
	for _, comm := range b.commands {
		commData := service.CommandData{Name: comm.Name, Description: comm.Description}

		if err := service.RegisterCommand(ctx, b.w, commData); err != nil {
			return nil, fmt.Errorf("cannot start the broker: %v", err)
		}

		topicName := strings.Trim("/", comm.Name)

		brokerDataChan := service.StartGetData(ctx, topicName, b.w.Addr.String())
		go pipelineData(ctx, brokerDataChan, dataChan, comm.Name)
	}

	return dataChan, nil
}

func pipelineData(ctx context.Context,
	brokerDataChan <-chan service.BotData,
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
	if err := b.w.Close(); err != nil {
		return fmt.Errorf("an error occurs at stopping the broker worker: %v", err)
	}

	return nil
}

func (b Broker) SendMessage(ctx context.Context, msg string, chatID int64) error {
	botData := service.BotData{ChatID: chatID, Value: msg}
	maxRetries := 10
	cnt := 0
	var err error

	for cnt < maxRetries {
		if err = service.SendData(ctx, b.w, botData); err != nil {
			time.Sleep(5 * time.Second)
			cnt++

			continue
		}

		return nil
	}

	return fmt.Errorf("an error occurs at sending a message to the broker: %v", err)
}
