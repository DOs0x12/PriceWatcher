package price

import (
	"PriceWatcher/internal/app/price/bank"
	mpService "PriceWatcher/internal/app/price/marketplace"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/entities/config"
	infraFile "PriceWatcher/internal/infrastructure/file"
	bankReq "PriceWatcher/internal/infrastructure/requester/bank"
	"PriceWatcher/internal/infrastructure/requester/marketplace"
	"fmt"
	"strings"
)

func NewPriceService(conf config.Config) (PriceService, error) {
	bankPriceType := "bank"
	marketplacePriceType := "marketplace"

	priceTypeInLowers := strings.ToLower(conf.PriceType)

	if priceTypeInLowers == bankPriceType {
		return createBankPriceService(conf), nil
	}

	if priceTypeInLowers == marketplacePriceType {
		return createMarketplacePriceService(conf), nil
	}

	return nil, fmt.Errorf("a price service is not created from the price type %v", conf.Marketplace)
}

func createBankPriceService(conf config.Config) PriceService {
	req := bankReq.BankRequester{}
	ext := createBankExtractor()

	return bank.NewService(req, ext, conf)
}

func createMarketplacePriceService(conf config.Config) PriceService {
	wr := infraFile.WriteReader{}
	req := marketplace.MarketplaceRequester{}
	marketplaceTypeInLowers := strings.ToLower(conf.Marketplace)
	ext := createMarketplaceExtractor(marketplaceTypeInLowers)
	analyser := analyser.MarketplaceAnalyser{}

	return mpService.NewService(wr, req, ext, analyser, conf)
}

func createMarketplaceExtractor(marketplaceType string) extractor.Extractor {
	var pageReg, tag string

	wbType := "wb"
	ozonType := "ozon"

	marketplaceTypeInLowers := strings.ToLower(marketplaceType)

	//TODO: write an error if no type
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
