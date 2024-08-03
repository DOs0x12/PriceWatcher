package requester

import "PriceWatcher/internal/entities/bank"

type Requester interface {
	RequestPage() (bank.Response, error)
}
