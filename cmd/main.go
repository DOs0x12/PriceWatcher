package main

import (
	"GoldRateGetter/internal/app"
	"GoldRateGetter/internal/domain"
	"GoldRateGetter/internal/infrastructure/requester"
	"GoldRateGetter/internal/infrastructure/sender"
)

func main() {
	req := requester.Requester{}
	ext := domain.RateExtractor{}
	sen := sender.Sender{}

	app.HandleGoldRate(req, ext, sen)
}
