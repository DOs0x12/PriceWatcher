package bot

import "context"

type Worker interface {
	Start(ctx context.Context) error
	Stop()
	SendMessage(msg string, chatID int64) error
}
