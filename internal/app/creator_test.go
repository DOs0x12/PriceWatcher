package app

import (
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price"
	"PriceWatcher/internal/entities/config"
	"reflect"
	"testing"
)

type crtTSender struct{}

func (s crtTSender) Send(price float32, config config.Email) error { return nil }

type bankConfiger struct{}

func (bankConfiger) GetConfig() (config.Config, error) {
	return config.Config{PriceType: "bank"}, nil
}

type marketplaceConfiger struct{}

func (marketplaceConfiger) GetConfig() (config.Config, error) {
	return config.Config{PriceType: "marketplace"}, nil
}

func TestNewPriceService(t *testing.T) {
	serv, err := NewPriceService(
		crtTSender{},
		price.PriceExtractor{},
		message.MessageHourVal{},
		bankConfiger{})
	if err != nil {
		t.Errorf("The method retuns an error: %v", err)
	}

	want := "BankRequester"
	got := reflect.TypeOf(serv.req).Name()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	want = "BankExtractor"
	got = reflect.TypeOf(serv.ext).Name()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	serv, err = NewPriceService(
		crtTSender{},
		price.PriceExtractor{},
		message.MessageHourVal{},
		marketplaceConfiger{})
	if err != nil {
		t.Errorf("The method retuns an error: %v", err)
	}

	want = "MarketplaceRequester"
	got = reflect.TypeOf(serv.req).Name()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	want = "MarketplaceExtractor"
	got = reflect.TypeOf(serv.ext).Name()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
