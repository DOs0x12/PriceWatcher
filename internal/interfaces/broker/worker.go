package broker

import (
	"PriceWatcher/internal/entities/bot"
	"context"
)

type Worker interface {
	Start(ctx context.Context) (chan<- bot.Message, error)
	Stop() error
	SendMessage(ctx context.Context, msg string, chatID int64) error
}
