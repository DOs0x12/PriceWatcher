package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/interfaces"
)

func HandleGoldPrice(req interfaces.Requester, ext domain.Extractor, sender interfaces.Sender) {
	response := req.RequestPage()
	price := ext.ExtractPrice(response.Body)
	sender.Send(price)
}
