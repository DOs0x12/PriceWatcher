package telebot

import (
	"PriceWatcher/internal/config"
	"PriceWatcher/internal/entities/telebot"
	"fmt"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Telebot struct {
	bot *tgbot.BotAPI
}

func NewTelebot(configer config.Configer) (Telebot, error) {
	config, err := configer.GetConfig()
	if err != nil {
		var zero Telebot

		return zero, fmt.Errorf("can not get the config data: %v", err)
	}

	botApi, err := tgbot.NewBotAPI(config.BotKey)
	if err != nil {
		var zero Telebot
		return zero, fmt.Errorf("getting an error at connecting to the bot: %v", err)
	}

	return Telebot{bot: botApi}, nil
}

func (t Telebot) Start(commands []telebot.Command,
	restart chan<- interface{}) error {
	updConfig := tgbot.NewUpdate(0)
	updCh := t.bot.GetUpdatesChan(updConfig)
	go t.watchUpdates(updCh, commands, restart)

	return nil
}

func (t Telebot) RegisterCommands(commands []telebot.Command) error {
	if err := t.configureCommands(commands); err != nil {
		return fmt.Errorf("getting an error at registering commands: %v", err)
	}

	return nil
}

func (t Telebot) Stop() {
	t.bot.StopReceivingUpdates()
}

func (t Telebot) watchUpdates(updCh tgbot.UpdatesChannel,
	commands []telebot.Command,
	restart chan<- interface{}) {
	for upd := range updCh {
		if upd.Message == nil {
			continue
		}

		if !upd.Message.IsCommand() {
			continue
		}

		for _, command := range commands {
			if upd.Message.Text == command.Name {
				msg := tgbot.NewMessage(upd.Message.Chat.ID, command.Action(upd))

				maxRetries := 10
				cnt := 0

				for cnt < maxRetries {
					if _, err := t.bot.Send(msg); err != nil {
						logrus.Errorf("Cannot send a message: %v", err)

						time.Sleep(5 * time.Second)
						cnt++

						continue
					}

					break
				}
			}
		}
	}
}
