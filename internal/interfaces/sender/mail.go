package sender

import "PriceWatcher/internal/entities/config"

type Sender interface {
	Send(message, subject string, conf config.Email) error
}
