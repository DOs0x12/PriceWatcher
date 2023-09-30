package app

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/app/price/bank"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/entities/config"
	infraFile "PriceWatcher/internal/infrastructure/file"
	bankReq "PriceWatcher/internal/infrastructure/requester/bank"
	"PriceWatcher/internal/infrastructure/requester/marketplace"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"
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

	req, err := createRequester(config)
	if err != nil {
		return nil, err
	}

	priceType := strings.ToLower(config.PriceType)
	marketplaceType := strings.ToLower(config.Marketplace)

	ext := createPriceExtractor(priceType, marketplaceType)
	wr := createWriteReader(priceType)

	analyser, err := createAnalyser(config.PriceType)
	if err != nil {
		return nil, err
	}

	bankType := "bank"
	var priceService price.PriceService

	if strings.ToLower(config.PriceType) == bankType {
		priceService = bank.NewService(req, ext, val, clock.RealClock{})
	}

	crt := PriceService{
		req:          req,
		sender:       sender,
		ext:          ext,
		val:          val,
		analyser:     analyser,
		wr:           wr,
		conf:         conf,
		priceService: priceService,
	}

	return &crt, err
}

func createRequester(conf config.Config) (interReq.Requester, error) {
	priceType := strings.ToLower(conf.PriceType)

	switch priceType {
	case "bank":
		return bankReq.BankRequester{}, nil
	case "marketplace":
		return marketplace.MarketplaceRequester{}, nil
	default:
		return nil, fmt.Errorf("have the unknown price type: %v", conf.PriceType)
	}
}

func createPriceExtractor(priceType, marketplaceType string) extractor.Extractor {
	pageReg, priceReg, tag := getSearchData(priceType, marketplaceType)

	return extractor.New(pageReg, priceReg, tag)
}

func getSearchData(priceType, marketplaceType string) (pageReg, priceReg, tag string) {
	switch priceType {
	case "bank":
		pageReg = `(^ покупка: [0-9]{4,5}\.[0-9][0-9])`
		tag = "td"
		return pageReg, priceReg, tag
	case "marketplace":
		if marketplaceType == "wb" {
			pageReg = "([0-9])*(\u00a0)*([0-9])*(\u00a0)[₽]"
			tag = "ins"

			return pageReg, priceReg, tag
		}

		pageReg = "([0-9])*(\u2009)*([0-9])*(\u2009)[₽]"
		tag = "span"

		return pageReg, priceReg, tag
	default:
		return "", "", ""
	}
}

func createAnalyser(priceType string) (analyser.Analyser, error) {
	pType := strings.ToLower(priceType)

	switch pType {
	case "bank":
		return nil, nil
	case "marketplace":
		return analyser.MarketplaceAnalyser{}, nil
	default:
		return nil, fmt.Errorf("have the unknown price type: %v", pType)
	}
}

func createWriteReader(priceType string) file.WriteReader {
	pType := strings.ToLower(priceType)

	switch pType {
	case "bank":
		return nil
	case "marketplace":
		return infraFile.WriteReader{}
	default:
		return nil
	}
}
