package bank

import (
	"PriceWatcher/internal/entities/page"
	"fmt"
	"net/http"
)

type BankRequester struct {
	Url string
}

func (r BankRequester) RequestPage() (page.Response, error) {
	resp, err := http.Get(r.Url)
	if err != nil {
		return page.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", err)
	}

	return page.Response{Body: resp.Body}, nil
}
