package bank

import (
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/entities/page"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

var testNow time.Time

type testClock struct{}

func (testClock) Now() time.Time                         { return testNow }
func (testClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

var (
	price   = 3232.00
	pageReg = `(^ покупка: [0-9]{4,5}\.[0-9][0-9])`
	tag     = "td"
)

type testRequester struct{}

func (testRequester) RequestPage(url string) (page.Response, error) {
	s := `
		<html>
			<head>
				<title>test</title>
			</head>
			<body>
				<table>
					<td> покупка: %.2f</td>
				</table>
			</body>
		</html>`
	interS := fmt.Sprintf(string(s), price)
	reader := strings.NewReader(interS)

	return page.Response{Body: reader}, nil
}

func testServePriceE2E(t *testing.T) {
	workHour := 12
	now := time.Now()
	testNow = time.Date(
		now.Year(), now.Month(), now.Day(), workHour,
		now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	serv := NewService(testRequester{}, extractor.New(pageReg, tag), message.MessageHourVal{}, testClock{})

	config := config.Config{SendingHours: []int{workHour}}

	message, subject, err := serv.ServePrice(config)
	if err != nil {
		t.Errorf("Got an error in the serve method: %v", err)
	}

	wantedMes := fmt.Sprintf("Курс золота. Продажа: %.2fр", price)
	wantedSub := "Че по золоту?"

	if message != wantedMes {
		t.Errorf("Got %v, wanted %v", message, wantedMes)
	}

	if subject != wantedSub {
		t.Errorf("Got %v, wanted %v", subject, wantedSub)
	}
}

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

type requesterWithoutPrice struct{}

func (requesterWithoutPrice) RequestPage(url string) (page.Response, error) {
	s := `
		<html>
			<head>
				<title>test</title>
			</head>
			<body>
				<table>
				</table>
			</body>
		</html>`
	reader := strings.NewReader(s)

	return page.Response{Body: reader}, nil
}

func testServePriceError(t *testing.T) {
	workHour := 12
	now := time.Now()
	testNow = time.Date(
		now.Year(), now.Month(), now.Day(), workHour,
		now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	serv := NewService(requesterWithoutPrice{}, extractor.New(pageReg, tag), message.MessageHourVal{}, testClock{})

	config := config.Config{SendingHours: []int{workHour}}

	message, subject, err := serv.ServePrice(config)
	wrappedError := fmt.Errorf("the document does not have a price value with the tag: %v", tag)
	want := fmt.Errorf("cannot extract the price from the body: %w", wrappedError)
	if message == "" && subject == "" && err != nil && want.Error() != err.Error() {
		t.Errorf("Got an error in the serve method: %v", err)
	}
}

func TestServePrice(t *testing.T) {
	testServePriceCalls(t)
	testServePriceE2E(t)
	testServePriceError(t)
}
