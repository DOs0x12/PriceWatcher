package broker

import (
	"PriceWatcher/internal/entities/bot"
	"context"

	"github.com/google/uuid"
)

type Worker interface {
	Start(ctx context.Context, serviceName string) (<-chan bot.Message, error)
	Stop()
	SendMessage(ctx context.Context, msg string, chatID int64) error
	CommitMessage(ctx context.Context, msgUuid uuid.UUID) error
}
