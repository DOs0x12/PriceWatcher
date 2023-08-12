package requester

import (
	"GoldPriceGetter/internal/entities/page"
)

type Requester interface {
	RequestPage() (page.Response, error)
}
