package interfaces

import "PriceWatcher/internal/entities/bank/page"

type Requester interface {
	RequestPage() (page.Response, error)
}
