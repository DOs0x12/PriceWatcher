package price

import (
	"PriceWatcher/internal/app/price/bank"
	mpService "PriceWatcher/internal/app/price/marketplace"
	"PriceWatcher/internal/app/time"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	infraFile "PriceWatcher/internal/infrastructure/file"
	bankReq "PriceWatcher/internal/infrastructure/requester/bank"
	"PriceWatcher/internal/infrastructure/requester/marketplace"
	"fmt"
	"strings"
)

func NewPriceService(priceType, marketplaceType string) (PriceService, error) {
	bankPriceType := "bank"
	marketplacePriceType := "marketplace"

	priceTypeInLowers := strings.ToLower(priceType)

	if priceTypeInLowers == bankPriceType {
		return createBankPriceService(), nil
	}

	if priceTypeInLowers == marketplacePriceType {
		marketplaceTypeInLowers := strings.ToLower(marketplaceType)
		return createMarketplacePriceService(marketplaceTypeInLowers), nil
	}

	return nil, fmt.Errorf("a price service is not created from the price type %v", priceType)
}

func createBankPriceService() PriceService {
	req := bankReq.BankRequester{}
	ext := createBankExtractor()
	val := message.MessageHourVal{}

	return bank.NewService(req, ext, val, time.RealClock{})
}

func createMarketplacePriceService(marketplaceType string) PriceService {
	wr := infraFile.WriteReader{}
	req := marketplace.MarketplaceRequester{}
	ext := createMarketplaceExtractor(marketplaceType)
	analyser := analyser.MarketplaceAnalyser{}

	return mpService.NewService(wr, req, ext, analyser)
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
