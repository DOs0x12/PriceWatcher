package app

import (
	"PriceWatcher/internal/domain/hour"
	pDomain "PriceWatcher/internal/domain/page"
	"PriceWatcher/internal/entities/config"
	pEnt "PriceWatcher/internal/entities/page"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestServe(t *testing.T) {
	logrus.Info("Start to test the func serve with true value")
	serveWithTrueValue(t)
	logrus.Info("Start to test the func serve for checking that the all methods are called")
	serveWithCall(t)
}

var testNow time.Time

type testClock struct{}

func (testClock) Now() time.Time                         { return testNow }
func (testClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

type testReadCloser struct {
	Reader io.Reader
}

func (testReadCloser) Close() error                       { return nil }
func (t testReadCloser) Read(p []byte) (n int, err error) { return t.Reader.Read(p) }

type testRequester struct{}

func (r testRequester) RequestPage() (pEnt.Response, error) {
	s := `
		<html>
			<head>
				<title>test</title>
			</head>
			<body>
				<table>
					<td> покупка: 3232.00</td>
				</table>
			</body>
		</html>`
	reader := strings.NewReader(s)
	rc := testReadCloser{reader}
	return pEnt.Response{Body: rc}, nil
}

type testSender struct{}

func (s testSender) Send(price float32, config config.Email) error { return nil }

type testConfiger struct{}

func (testConfiger) GetConfig() (config.Config, error) {
	return config.Config{}, nil
}

func serveWithTrueValue(t *testing.T) {
	serv := NewGoldPriceService(
		testRequester{},
		testSender{},
		pDomain.PriceExtractor{},
		hour.MessageHourVal{},
		testConfiger{})

	workHour := 12
	now := time.Now()
	testNow = time.Date(
		now.Year(), now.Month(), now.Day(), workHour,
		now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	got := serv.serve(testClock{})

	if got != nil {
		t.Errorf("Got an error in the serve method: %v", got.Error())
	}
}

var (
	reqCall  bool
	extCall  bool
	valCall  bool
	sendCall bool
	confCall bool
)

type reqWithCall struct{}

func (r reqWithCall) RequestPage() (pEnt.Response, error) {
	s := "test"
	reader := strings.NewReader(s)
	rc := testReadCloser{reader}

	reqCall = true

	return pEnt.Response{Body: rc}, nil
}

type extWithCall struct{}

func (extWithCall) ExtractPrice(body io.ReadCloser) (float32, error) {
	extCall = true

	return 0.00, nil
}

type valWithCall struct{}

func (valWithCall) Validate(hour int, sendHours []int) bool {
	valCall = true

	return true
}

type sendWithCall struct{}

func (sendWithCall) Send(price float32, conf config.Email) error {
	sendCall = true

	return nil
}

type confWithCall struct{}

func (confWithCall) GetConfig() (config.Config, error) {
	confCall = true

	return config.Config{}, nil
}

func serveWithCall(t *testing.T) {
	serv := NewGoldPriceService(
		reqWithCall{},
		sendWithCall{},
		extWithCall{},
		valWithCall{},
		confWithCall{})

	workHour := 12
	now := time.Now()
	testNow = time.Date(
		now.Year(), now.Month(), now.Day(), workHour,
		now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	serv.serve(testClock{})

	if !reqCall {
		t.Error("The method for requesting a page is not called in the app layer")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called in the app layer")
	}
	if !valCall {
		t.Error("The method for validating the price is not called in the app layer")
	}
	if !sendCall {
		t.Error("The method for sending the price is not called in the app layer")
	}
}
