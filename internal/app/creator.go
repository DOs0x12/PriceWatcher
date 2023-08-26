package app

import (
	"PriceWatcher/internal/domain/hour"
	"PriceWatcher/internal/domain/page"
	"PriceWatcher/internal/infrastructure/requester/bank"
	"PriceWatcher/internal/infrastructure/requester/marketplace"
	"PriceWatcher/internal/interfaces/configer"
	interReq "PriceWatcher/internal/interfaces/requester"
	interSend "PriceWatcher/internal/interfaces/sender"
	"fmt"
	"strings"
)

func NewPriceService(
	sender interSend.Sender,
	ext page.Extractor,
	val hour.HourValidator,
	conf configer.Configer) (*PriceService, error) {

	config, err := conf.GetConfig()
	if err != nil {
		return nil, err
	}

	req, err := createRequester(strings.ToLower(config.PriceType))
	if err != nil {
		return nil, err
	}

	crt := PriceService{
		req:    req,
		sender: sender,
		ext:    ext,
		val:    val,
		conf:   conf,
	}

	return &crt, err
}

func createRequester(priceType string) (interReq.Requester, error) {
	switch priceType {
	case "bank":
		return bank.BankRequester{}, nil
	case "marketplace":
		return marketplace.MarketplaceRequester{}, nil
	default:
		return nil, fmt.Errorf("have the unknown price type: %v", priceType)
	}
}
