package page

import "PriceWatcher/internal/entities/bank"

type Requester interface {
	RequestPage() (bank.Page, error)
}
