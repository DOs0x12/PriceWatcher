package marketplace

import (
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/page"
	"PriceWatcher/internal/entities/price"
	"io"
	"testing"
)

type wrWithCall struct{}

func (wrWithCall) WritePrices(prices map[string]price.ItemPrice) error {
	rwWriteCall = true

	return nil
}

func (wrWithCall) ReadPrices() (map[string]price.ItemPrice, error) {
	rwReadCall = true

	return map[string]price.ItemPrice{"test": {Address: "address", Price: 0.0}}, nil
}

type reqWithCall struct{}

func (reqWithCall) RequestPage(url string) (page.Response, error) {
	reqCall = true

	return page.Response{}, nil
}

type extWithCall struct{}

func (extWithCall) ExtractPrice(body io.Reader) (float32, error) {
	extCall = true

	return 0.0, nil
}

type analyserWithUpChangedCall struct{}

func (analyserWithUpChangedCall) AnalysePrice(price, initialPrice float32) (changed, up bool, amount float32) {
	analyzerCall = true

	return true, true, 1.0
}

type analyserWithDownChangedCall struct{}

func (analyserWithDownChangedCall) AnalysePrice(price, initialPrice float32) (changed, up bool, amount float32) {
	analyzerCall = true

	return true, false, 1.0
}

type analyserWithNotChangedCall struct{}

func (analyserWithNotChangedCall) AnalysePrice(price, initialPrice float32) (changed, up bool, amount float32) {
	analyzerCall = true

	return false, false, 0.0
}

var (
	rwWriteCall  bool
	rwReadCall   bool
	reqCall      bool
	extCall      bool
	analyzerCall bool
)

func testUpChangedServePriceCalls(t *testing.T) {
	itemName := "test"
	itemValue := "1.0"
	config := config.ServiceConf{Items: map[string]string{itemName: itemValue}, PriceType: "marketplace"}
	serv := NewService(wrWithCall{}, reqWithCall{}, extWithCall{}, analyserWithUpChangedCall{}, config)

	serv.ServePrice()

	if !rwWriteCall {
		t.Error("The method for writing the current prices is not called")
	}
	if !rwReadCall {
		t.Error("The method for writing the current prices is not called")
	}
	if !reqCall {
		t.Error("The method for requesting a page is not called")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called")
	}
	if !analyzerCall {
		t.Error("The method for analyzing the price is not called")
	}
}

func testNotChangedServePriceCalls(t *testing.T) {
	itemName := "test"
	itemValue := "1.0"
	config := config.ServiceConf{Items: map[string]string{itemName: itemValue}, PriceType: "marketplace"}
	serv := NewService(wrWithCall{}, reqWithCall{}, extWithCall{}, analyserWithNotChangedCall{}, config)

	serv.ServePrice()

	if !rwWriteCall {
		t.Error("The method for writing the current prices is not called")
	}
	if !rwReadCall {
		t.Error("The method for writing the current prices is not called")
	}
	if !reqCall {
		t.Error("The method for requesting a page is not called")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called")
	}
	if !analyzerCall {
		t.Error("The method for analyzing the price is not called")
	}
}

func testDownChangedServePriceCalls(t *testing.T) {
	itemName := "test"
	itemValue := "1.0"
	config := config.ServiceConf{Items: map[string]string{itemName: itemValue}, PriceType: "marketplace"}
	serv := NewService(wrWithCall{}, reqWithCall{}, extWithCall{}, analyserWithDownChangedCall{}, config)

	serv.ServePrice()

	if !rwWriteCall {
		t.Error("The method for writing the current prices is not called")
	}
	if !rwReadCall {
		t.Error("The method for writing the current prices is not called")
	}
	if !reqCall {
		t.Error("The method for requesting a page is not called")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called")
	}
	if !analyzerCall {
		t.Error("The method for analyzing the price is not called")
	}
}

func TestServePrice(t *testing.T) {
	testUpChangedServePriceCalls(t)
	testNotChangedServePriceCalls(t)
	testDownChangedServePriceCalls(t)
}
