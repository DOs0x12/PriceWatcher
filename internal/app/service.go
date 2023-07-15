package app

import (
	"GoldRateGetter/internal/interfaces"
)

func HandleGoldRate(req interfaces.Requester, pr interfaces.Processor, sender interfaces.Sender) {
	response := req.RequestPage()
	rate := pr.Process(response.Page)
	sender.Send(rate)
}
