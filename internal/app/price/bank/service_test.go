package bank

import (
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/page"
	"io"
	"testing"
	"time"
)

var testNow time.Time

type testClock struct{}

func (testClock) Now() time.Time                         { return testNow }
func (testClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

var (
	reqCall bool
	extCall bool
	valCall bool
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

type valWithCall struct{}

func (valWithCall) Validate(hour int, sendHours []int) bool {
	valCall = true

	return true
}

func testServePriceCalls(t *testing.T) {
	serv := NewService(reqWithCall{}, extWithCall{}, valWithCall{}, testClock{})

	config := config.Config{}
	serv.ServePrice(config)

	if !reqCall {
		t.Error("The method for requesting a page is not called in the app layer")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called in the app layer")
	}
	if !valCall {
		t.Error("The method for validating the price is not called in the app layer")
	}
}

func TestServePrice(t *testing.T) {
	testServePriceCalls(t)
}
