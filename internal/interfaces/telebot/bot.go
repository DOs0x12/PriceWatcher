package telebot

import "PriceWatcher/internal/entities/telebot"

type Bot interface {
	Start(commands ...telebot.Command) error
	RegisterCommands(commands []telebot.Command) error
	Stop()
}
