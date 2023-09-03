package app

import (
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/domain/price/extractor"
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

type servTClock struct{}

func (servTClock) Now() time.Time                         { return testNow }
func (servTClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

type servTRequester struct{}

func (r servTRequester) RequestPage() (pEnt.Response, error) {
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

	return pEnt.Response{Body: reader}, nil
}

type servTSender struct{}

func (s servTSender) Send(price float32, config config.Email) error { return nil }

type servTConfiger struct{}

func (servTConfiger) GetConfig() (config.Config, error) {
	return config.Config{}, nil
}

type servTAnalyser struct{}

func (servTAnalyser) AnalysePrice(price float32) (changed, up bool, amount float32) {
	return false, false, 0.0
}

func serveWithTrueValue(t *testing.T) {
	serv := PriceService{
		servTRequester{},
		servTSender{},
		extractor.PriceExtractor{},
		message.MessageHourVal{},
		servTAnalyser{},
		servTConfiger{}}

	workHour := 12
	now := time.Now()
	testNow = time.Date(
		now.Year(), now.Month(), now.Day(), workHour,
		now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	got := serv.serve(servTClock{})

	if got != nil {
		t.Errorf("Got an error in the serve method: %v", got.Error())
	}
}

var (
	reqCall      bool
	extCall      bool
	valCall      bool
	sendCall     bool
	confCall     bool
	analyserCall bool
)

type reqWithCall struct{}

func (r reqWithCall) RequestPage() (pEnt.Response, error) {
	s := "test"
	reader := strings.NewReader(s)

	reqCall = true

	return pEnt.Response{Body: reader}, nil
}

type extWithCall struct{}

func (extWithCall) ExtractPrice(body io.Reader) (float32, error) {
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

type analyserWithCall struct{}

func (analyserWithCall) AnalysePrice(price float32) (changed, up bool, amount float32) {
	analyserCall = true

	return false, false, 0.0
}

func serveWithCall(t *testing.T) {
	serv := PriceService{
		reqWithCall{},
		sendWithCall{},
		extWithCall{},
		valWithCall{},
		analyserWithCall{},
		confWithCall{}}

	workHour := 12
	now := time.Now()
	testNow = time.Date(
		now.Year(), now.Month(), now.Day(), workHour,
		now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	serv.serve(servTClock{})

	if !reqCall {
		t.Error("The method for requesting a page is not called in the app layer")
	}
	if !extCall {
		t.Error("The method for extracting a price is not called in the app layer")
	}
	if !valCall {
		t.Error("The method for validating the price is not called in the app layer")
	}
	if !analyserCall {
		t.Error("The method for analysing the price is not called in the app layer")
	}
	if !sendCall {
		t.Error("The method for sending the price is not called in the app layer")
	}
}
