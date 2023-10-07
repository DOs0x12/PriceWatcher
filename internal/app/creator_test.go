package app

import (
	"PriceWatcher/internal/app/price/bank"
	"PriceWatcher/internal/app/price/marketplace"
	"PriceWatcher/internal/entities/config"
	"reflect"
	"strings"
	"testing"
)

type crtTSender struct{}

func (s crtTSender) Send(message, subject string, conf config.Email) error { return nil }

type bankConfiger struct{}

func (bankConfiger) GetConfig() (config.Config, error) {
	return config.Config{PriceType: "bank"}, nil
}

type marketplaceConfiger struct{}

func (marketplaceConfiger) GetConfig() (config.Config, error) {
	return config.Config{PriceType: "marketplace"}, nil
}

type configerWithUndefinedType struct{}

func (configerWithUndefinedType) GetConfig() (config.Config, error) {
	return config.Config{PriceType: "undefined"}, nil
}

func TestNewPriceService(t *testing.T) {
	createBankService(t)
	createMarketplaceService(t)
	getCreationError(t)
}

func createBankService(t *testing.T) {
	serv, err := NewPriceService(
		crtTSender{},
		bankConfiger{})
	if err != nil {
		t.Errorf("The method retuns an error: %v", err)
	}

	want := "bank.Service"
	bankService, ok := serv.priceService.(bank.Service)

	if !ok {
		t.Error("Cannot convert priceService to bank.Service")
	}

	got := reflect.TypeOf(bankService).String()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func createMarketplaceService(t *testing.T) {
	serv, err := NewPriceService(
		crtTSender{},
		marketplaceConfiger{})
	if err != nil {
		t.Errorf("The method retuns an error: %v", err)
	}

	want := "marketplace.Service"
	marketplaceService, ok := serv.priceService.(marketplace.Service)

	if !ok {
		t.Error("Cannot convert priceService to marketplace.Service")
	}

	got := reflect.TypeOf(marketplaceService).String()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func getCreationError(t *testing.T) {
	errTemplt := "a price service is not created from the price type"
	_, err := NewPriceService(
		crtTSender{},
		configerWithUndefinedType{})
	if err != nil && !strings.Contains(err.Error(), errTemplt) {
		t.Errorf("Got not wanted error: %v, wanted error template: %v", err, errTemplt)
	}
}
