package requester

import (
	"GoldPriceGetter/internal/entities"
	"fmt"
	"net/http"
)

// TODO: add a config and get the url value
const goldPriceUrl = "https://investzoloto.ru/gold-sber-oms/"

type Requester struct{}

func (r Requester) RequestPage() (entities.Response, error) {
	resp, err := http.Get(goldPriceUrl)
	if err != nil {
		requestError := fmt.Sprintf("Cannot get the data from the address: %v", goldPriceUrl)

		return entities.Response{Body: nil}, fmt.Errorf(requestError, err)
	}

	return entities.Response{Body: resp.Body}, nil
}
