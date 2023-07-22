package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
)

func main() {
	req := requester.Requester{}
	ext := domain.PriceExtractor{}
	sen := sender.Sender{}

	app.HandleGoldPrice(req, ext, sen)
}
