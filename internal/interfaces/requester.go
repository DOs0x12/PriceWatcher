package interfaces

import "GoldRateGetter/internal/entities"

type Requester interface {
	RequestPage() entities.Response
}
