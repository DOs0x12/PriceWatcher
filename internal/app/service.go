package app

import (
	"GoldRateGetter/internal/domain"
	"GoldRateGetter/internal/interfaces"
)

func HandleGoldRate(req interfaces.Requester, pr domain.Processor, sender interfaces.Sender) {
	response := req.RequestPage()
	rate := pr.Process(response.Body)
	sender.Send(rate)
}
