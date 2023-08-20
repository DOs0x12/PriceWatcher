package requester

import (
	"PriceWatcher/internal/entities/page"
)

type Requester interface {
	RequestPage(url string) (page.Response, error)
}
