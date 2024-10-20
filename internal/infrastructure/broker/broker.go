package bot

import (
	telebot "PriceWatcher/internal/entities/bot"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Guise322/TeleBot/client/pkg/service"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	commands []telebot.Command
	w        *kafka.Writer
}

func (b Broker) Start(ctx context.Context) (chan<- telebot.Message, error) {
	dataChan := make(chan telebot.Message)
	for _, comm := range b.commands {
		commData := service.CommandData{Name: comm.Name, Description: comm.Description}

		if err := service.RegisterCommand(ctx, b.w, commData); err != nil {
			return nil, fmt.Errorf("cannot start the bot: %v", err)
		}

		topicName := strings.Trim("/", comm.Name)

		brokerDataChan := service.StartGetData(ctx, topicName, b.w.Addr.String())
		go pipelineData(ctx, brokerDataChan, dataChan)
	}

	return dataChan, nil
}

func pipelineData(ctx context.Context,
	brokerDataChan <-chan service.BotData,
	msgChan chan<- telebot.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case brokerData := <-brokerDataChan:
			msg := telebot.Message{ChatID: brokerData.ChatID, Value: brokerData.Value}
			msgChan <- msg
		}
	}

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

	return err
}
