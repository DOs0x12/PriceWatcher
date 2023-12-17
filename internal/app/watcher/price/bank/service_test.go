package bank

import (
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/page"
	"io"
	"testing"
)

var (
	reqCall bool
	extCall bool
)

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

func testCalls(t *testing.T) {
	serv := NewService(reqWithCall{}, extWithCall{}, config.ServiceConf{})
	serv.ServePrice()

	if !reqCall {
		t.Error("The method for requesting a page is not called")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called")
	}
}

func TestServePrice(t *testing.T) {
	testCalls(t)
}
