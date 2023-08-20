package requester

import (
	"PriceWatcher/internal/entities/page"
	"fmt"
	"net/http"
)

const goldPriceUrl = "https://investzoloto.ru/gold-sber-oms/"

type Requester struct{}

func (r Requester) RequestPage() (page.Response, error) {
	resp, err := http.Get(goldPriceUrl)
	if err != nil {
		return page.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", err)
	}

	return page.Response{Body: resp.Body}, nil
}
