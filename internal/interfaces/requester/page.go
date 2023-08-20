package requester

import (
	"PriceWatcher/internal/entities/page"
)

type Requester interface {
	RequestPage() (page.Response, error)
}
