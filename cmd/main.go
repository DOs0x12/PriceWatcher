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

	serv := app.NewGoldPriceService(req, ext, sen)

	app.WatchGoldPrice(serv)
}
