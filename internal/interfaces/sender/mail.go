package sender

import "GoldPriceGetter/internal/entities/config"

type Sender interface {
	Send(price float32, conf config.Email) error
}
