package bot

import "github.com/google/uuid"

type Message struct {
	ChatID       int64
	Command      string
	Value        string
	MsgUuid      uuid.UUID
	ReceiverUuid uuid.UUID
}
