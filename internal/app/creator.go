package app

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/app/price"
	"PriceWatcher/internal/app/price/bank"
	marketplaceService "PriceWatcher/internal/app/price/marketplace"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	infraFile "PriceWatcher/internal/infrastructure/file"
	bankReq "PriceWatcher/internal/infrastructure/requester/bank"
	"PriceWatcher/internal/infrastructure/requester/marketplace"
	"PriceWatcher/internal/interfaces/configer"
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
	var priceService price.PriceService
	bankPriceType := "bank"
	marketplacePriceType := "marketplace"

	if priceType == bankPriceType {
		priceService = createBankPriceService()
	}

	if priceType == marketplacePriceType {
		marketplaceType := strings.ToLower(config.Marketplace)
		priceService = createMarketplacePriceService(marketplaceType)
	}

	if priceService == nil {
		return nil, fmt.Errorf("a price service is not created from the price type %v", priceType)
	}

	service := PriceService{
		sender:       sender,
		conf:         conf,
		priceService: priceService,
	}

	return &service, err
}

func createBankPriceService() price.PriceService {
	req := bankReq.BankRequester{}
	ext := createBankExtractor()
	val := message.MessageHourVal{}

	return bank.NewService(req, ext, val, clock.RealClock{})
}

func createMarketplacePriceService(marketplaceType string) price.PriceService {
	wr := infraFile.WriteReader{}
	req := marketplace.MarketplaceRequester{}
	ext := createMarketplaceExtractor(marketplaceType)
	analyser := analyser.MarketplaceAnalyser{}

	return marketplaceService.NewService(wr, req, ext, analyser)
}

func createMarketplaceExtractor(marketplaceType string) extractor.Extractor {
	var pageReg, tag string

	wbType := "wb"
	ozonType := "ozon"

	marketplaceTypeInLowers := strings.ToLower(marketplaceType)

	if marketplaceTypeInLowers == wbType {
		pageReg = "([0-9])*(\u00a0)*([0-9])*(\u00a0)[₽]"
		tag = "ins"
	}

	if marketplaceTypeInLowers == ozonType {
		pageReg = "([0-9])*(\u2009)*([0-9])*(\u2009)[₽]"
		tag = "span"
	}

	return extractor.New(pageReg, tag)
}

func createBankExtractor() extractor.Extractor {
	pageReg := `(^ покупка: [0-9]{4,5}\.[0-9][0-9])`
	tag := "td"

	return extractor.New(pageReg, tag)
}
