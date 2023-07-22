package app

import (
	"GoldRateGetter/internal/domain"
	"GoldRateGetter/internal/interfaces"
)

func HandleGoldRate(req interfaces.Requester, ext domain.Extractor, sender interfaces.Sender) {
	response := req.RequestPage()
	rate := ext.ExtractRate(response.Body)
	sender.Send(rate)
}
