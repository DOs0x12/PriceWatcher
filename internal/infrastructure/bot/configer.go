package telebot

import (
	"PriceWatcher/internal/entities/telebot"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t Telebot) configureCommands(commands []telebot.Command) error {
	commandSet := make([]tgbot.BotCommand, len(commands))

	for i, command := range commands {
		commandSet[i] = tgbot.BotCommand{Command: command.Name, Description: command.Description}
	}

	config := tgbot.NewSetMyCommands(commandSet...)

	_, err := t.bot.Request(config)

	return err
}
