package interfaces

import "context"

type Bot interface {
	Start(ctx context.Context) error
	Stop()
	SendMessage(msg string, chatID int64) error
}
