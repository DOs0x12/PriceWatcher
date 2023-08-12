package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/entities"
	"fmt"
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

func (r testRequester) RequestPage() (entities.Response, error) {
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
	return entities.Response{Body: rc}, nil
}

type testSender struct{}

func (s testSender) Send(price float32) error { return nil }

func serveWithTrueValue(t *testing.T) {
	serv := NewGoldPriceService(
		testRequester{},
		testSender{},
		domain.PriceExtractor{},
		domain.MessageHourVal{})

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
)

type reqWithCall struct{}

func (r reqWithCall) RequestPage() (entities.Response, error) {
	s := "test"
	reader := strings.NewReader(s)
	rc := testReadCloser{reader}

	reqCall = true

	return entities.Response{Body: rc}, nil
}

type extWithCall struct{}

func (extWithCall) ExtractPrice(body io.ReadCloser) (float32, error) {
	extCall = true

	return 0.00, nil
}

type valWithCall struct{}

func (valWithCall) Validate(hour int) bool {
	valCall = true

	return true
}

type sendWithCall struct{}

func (sendWithCall) Send(price float32) error {
	sendCall = true

	return nil
}

func serveWithCall(t *testing.T) {
	serv := NewGoldPriceService(
		reqWithCall{},
		sendWithCall{},
		extWithCall{},
		valWithCall{})

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

func TestGetTuneTime(t *testing.T) {
	nT := time.Now()
	testMin := 45
	testSec := 45
	testNow :=
		time.Date(nT.Year(), nT.Month(), nT.Day(), nT.Hour(), testMin, testSec, nT.Nanosecond(), nT.Location())

	waitMin := 60 - testMin
	waitSec := 60 - testSec

	durStr := fmt.Sprintf("%vm%vs", waitMin, waitSec)
	want, err := time.ParseDuration(durStr)
	if err != nil {
		t.Errorf("An error occurs while parsing duration in the test: %v", err)
	}

	got, err := getWaitTime(testNow)
	if err != nil {
		t.Errorf("The getWaitTime method retuns an error: %v", err)
	}

	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
