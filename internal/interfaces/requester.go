package interfaces

import "GoldPriceGetter/internal/entities"

type Requester interface {
	RequestPage() entities.Response
}
