package telebot

import (
	"PriceWatcher/internal/entities/telebot"
	"PriceWatcher/internal/interfaces/configer"
	"fmt"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Telebot struct {
	bot *tgbot.BotAPI
}

func NewTelebot(configer configer.Configer) (Telebot, error) {
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

func (t Telebot) Start(commands []telebot.Command, commandsWithInput []telebot.CommandWithInput) error {
	updConfig := tgbot.NewUpdate(0)
	updCh := t.bot.GetUpdatesChan(updConfig)
	go t.watchUpdates(updCh, commands, commandsWithInput)

	return nil
}

func (t Telebot) RegisterCommands(commands []telebot.Command, commandsWithInput []telebot.CommandWithInput) error {
	if err := t.configureCommands(commands, commandsWithInput); err != nil {
		return fmt.Errorf("getting an error at registering commands: %v", err)
	}

	return nil
}

func (t Telebot) Stop() {
	t.bot.StopReceivingUpdates()
}

var addItemComm = false

func (t Telebot) watchUpdates(updCh tgbot.UpdatesChannel, commands []telebot.Command, commandsWithInput []telebot.CommandWithInput) {
	for upd := range updCh {
		if upd.Message == nil {
			continue
		}

		if !upd.Message.IsCommand() {
			if addItemComm {
				//TODO: add the item to the config for watching
				addItemComm = false
			}

			continue
		}

		for _, command := range commands {
			if upd.Message.Text == command.Name {
				msg := tgbot.NewMessage(upd.Message.Chat.ID, command.Action())

				maxRetries := 10
				cnt := 0

				for cnt < maxRetries {
					if _, err := t.bot.Send(msg); err != nil {
						logrus.Errorf("Can not send a message: %v", err)

						time.Sleep(5 * time.Second)
						cnt++

						continue
					}

					break
				}
			}
		}

		for _, CommandWithInput := range commandsWithInput {
			if upd.Message.Text == CommandWithInput.Name {
				//msg := tgbot.NewMessage(upd.Message.Chat.ID, CommandWithInput.Action())
				text := "Чтобы добавить товар для отслеживания пришлите даннные в формате: <наименование> <ссылка на товар>"
				msg := tgbot.NewMessage(upd.Message.Chat.ID, text)
				maxRetries := 10
				cnt := 0

				for cnt < maxRetries {
					if _, err := t.bot.Send(msg); err != nil {
						logrus.Errorf("Can not send a message: %v", err)

						time.Sleep(5 * time.Second)
						cnt++

						continue
					}

					break
				}

				addItemComm = true
			}
		}
	}
}
