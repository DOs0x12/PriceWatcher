package requester

import (
	"GoldRateGetter/internal/entities"
	"io"
	"net/http"
)

// TODO: add a config and get the url value
const goldRateUrl = "https://investzoloto.ru/gold-sber-oms/"

type Requester struct{}

func (r Requester) RequestPage() entities.Response {
	resp, err := http.Get(goldRateUrl)
	if err != nil {
		//TODO: appropriate error handling and logging
		return entities.Response{Page: nil}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		//TODO: appropriate error handling and logging
		return entities.Response{Page: nil}
	}

	sb := string(data)

	return entities.Response{Page: &sb}
}
