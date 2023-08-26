package app

import (
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price"
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
	val message.HourValidator,
	conf configer.Configer) (*PriceService, error) {

	config, err := conf.GetConfig()
	if err != nil {
		return nil, err
	}

	priceType := strings.ToLower(config.PriceType)

	req, err := createRequester(priceType)
	if err != nil {
		return nil, err
	}

	ext := createPriceExtractor(priceType)

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

func createPriceExtractor(priceType string) price.Extractor {
	pageReg, priceReg, tag := getSearchData(priceType)

	return price.New(pageReg, priceReg, tag)
}

func getSearchData(priceType string) (pageReg, priceReg, tag string) {
	switch priceType {
	case "bank":
		pageReg = `(^ покупка: [0-9]{4,5}\.[0-9][0-9])`
		priceReg = `([0-9]{4,5}\.[0-9][0-9])`
		tag = "td"
		return pageReg, priceReg, tag
	case "marketplace":
		pageReg = `([0-9])*(&nbsp;)*([0-9])*(&nbsp;)[₽];`
		priceReg = `([0-9]{4,5}\.[0-9][0-9])`
		tag = "ins"
		return pageReg, priceReg, tag
	default:
		return "", "", ""
	}
}
