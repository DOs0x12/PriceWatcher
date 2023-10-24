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

type valWithCallAndTrue struct{}

func (valWithCallAndTrue) Validate(hour int, sendHours []int) bool {
	valCall = true

	return true
}

type valWithCallAndFalse struct{}

func (valWithCallAndFalse) Validate(hour int, sendHours []int) bool {
	valCall = true

	return false
}

func testCallsWithTimeToCheck(t *testing.T) {
	serv := NewService(reqWithCall{}, extWithCall{}, valWithCallAndTrue{}, testClock{})

	config := config.Config{}
	serv.ServePrice(config)

	if !reqCall {
		t.Error("The method for requesting a page is not called")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called")
	}
	if !valCall {
		t.Error("The method for validating the price is not called")
	}
}

func testCallsWithNoTimeToCheck(t *testing.T) {
	serv := NewService(reqWithCall{}, extWithCall{}, valWithCallAndFalse{}, testClock{})

	config := config.Config{}
	serv.ServePrice(config)

	if reqCall {
		t.Error("The method for requesting a page is called")
	}
	if extCall {
		t.Error("The method for extracting a price is called")
	}
	if !valCall {
		t.Error("The method for validating the price is not called")
	}
}

func setCallsToDefault() {
	reqCall = false
	extCall = false
	valCall = false
}

func TestServePrice(t *testing.T) {
	testCallsWithTimeToCheck(t)
	setCallsToDefault()
	testCallsWithNoTimeToCheck(t)
}
