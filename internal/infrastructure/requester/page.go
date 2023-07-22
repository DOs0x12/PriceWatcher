package requester

import (
	"GoldPriceGetter/internal/entities"
	"net/http"
)

// TODO: add a config and get the url value
const goldPriceUrl = "https://investzoloto.ru/gold-sber-oms/"

type Requester struct{}

func (r Requester) RequestPage() entities.Response {
	resp, err := http.Get(goldPriceUrl)
	if err != nil {
		//TODO: appropriate error handling and logging
		return entities.Response{Body: nil}
	}

	return entities.Response{Body: resp.Body}
}
