package broker

import (
	"PriceWatcher/internal/entities/bot"
	"context"

	"github.com/google/uuid"
)

type Worker interface {
	Start(ctx context.Context) (<-chan bot.Message, error)
	Stop() error
	SendMessage(ctx context.Context, msg string, chatID int64) error
	CommitMessage(ctx context.Context, recUuid uuid.UUID, msgUuid uuid.UUID) error
}
