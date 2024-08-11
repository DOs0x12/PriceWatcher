package command

import (
	"PriceWatcher/internal/entities/bot"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ConfigureCommands(commands []bot.Command) tgbot.SetMyCommandsConfig {
	commandSet := make([]tgbot.BotCommand, len(commands))

	for i, command := range commands {
		commandSet[i] = tgbot.BotCommand{Command: command.Name, Description: command.Description}
	}

	return tgbot.NewSetMyCommands(commandSet...)
}
