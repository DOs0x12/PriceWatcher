package price

import (
	"PriceWatcher/internal/app/watcher/price/bank"
	"PriceWatcher/internal/app/watcher/price/marketplace"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/price"
	"reflect"
	"strings"
	"testing"
)

type FakeWriteReader struct{}

func (FakeWriteReader) WritePrices(prices map[string]price.ItemPrice) error {
	return nil
}
func (FakeWriteReader) ReadPrices() (map[string]price.ItemPrice, error) {
	return make(map[string]price.ItemPrice), nil
}

func TestNewPriceService(t *testing.T) {
	createBankService(t)
	createWBService(t)
	createOzonService(t)
	getCreationError(t)
}

func createBankService(t *testing.T) {
	priceType := "bank"
	config := config.ServiceConf{PriceType: priceType}

	serv, err := NewPriceService(config, FakeWriteReader{})
	if err != nil {
		t.Errorf("The method retuns an error: %v", err)
	}

	want := "bank.Service"
	bankService, ok := serv.(bank.Service)

	if !ok {
		t.Error("Cannot convert priceService to bank.Service")
	}

	got := reflect.TypeOf(bankService).String()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func createMarketplaceService(marketplaceType string) (PriceService, error) {
	priceType := "marketplace"
	config := config.ServiceConf{PriceType: priceType}

	return NewPriceService(config, FakeWriteReader{})
}

func createWBService(t *testing.T) {
	marketplaceType := "wb"

	serv, err := createMarketplaceService(marketplaceType)
	if err != nil {
		t.Errorf("The method of creating a WB service retuns an error: %v", err)
	}

	checkMarketplaceService(serv, t)
}

func createOzonService(t *testing.T) {
	marketplaceType := "ozon"

	serv, err := createMarketplaceService(marketplaceType)
	if err != nil {
		t.Errorf("The method of creating a WB service retuns an error: %v", err)
	}

	checkMarketplaceService(serv, t)
}

func checkMarketplaceService(serv PriceService, t *testing.T) {
	want := "marketplace.Service"
	marketplaceService, ok := serv.(marketplace.Service)

	if !ok {
		t.Error("Cannot convert priceService to marketplace.Service")
	}

	got := reflect.TypeOf(marketplaceService).String()
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func getCreationError(t *testing.T) {
	priceType := "test"
	marketplaceType := "test"
	config := config.ServiceConf{PriceType: priceType, Marketplace: marketplaceType}

	_, err := NewPriceService(config, FakeWriteReader{})
	errTemplt := "a price service is not created from the price type"
	if err != nil && !strings.Contains(err.Error(), errTemplt) {
		t.Errorf("Got not wanted error: %v, wanted error template: %v", err, errTemplt)
	}
}
