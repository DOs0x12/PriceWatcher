package telebot

import "PriceWatcher/internal/entities/telebot"

type Bot interface {
	Start(commands []telebot.Command, commandsWithInput []telebot.CommandWithInput, restart chan<- interface{}) error
	RegisterCommands(commands []telebot.Command, commandsWithInput []telebot.CommandWithInput) error
	Stop()
}
