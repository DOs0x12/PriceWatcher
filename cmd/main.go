package main

import (
	"GoldPriceGetter/internal/app"
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/infrastructure/requester"
	"GoldPriceGetter/internal/infrastructure/sender"
)

var (
	req requester.Requester
	ext domain.PriceExtractor
	sen sender.Sender
)

func main() {
	setGlobalVals()
	serv := app.NewGoldPriceService(req, ext, sen)
	app.WatchGoldPrice(serv)
}

func setGlobalVals() {
	req = requester.Requester{}
	ext = domain.PriceExtractor{}
	sen = sender.Sender{}
}
