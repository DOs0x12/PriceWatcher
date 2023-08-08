package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/entities"
	"io"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestServe(t *testing.T) {
	logrus.Info("Start to test the func serve with true value")
	withTrueValue(t)
	logrus.Info("Start to test the func serve for checking that the all methods are called")
	withCall(t)
}

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

func withTrueValue(t *testing.T) {
	serv := NewGoldPriceService(
		testRequester{},
		testSender{},
		domain.PriceExtractor{},
		domain.MessageHourVal{})

	workHour := 12
	nowHour = func() int { return workHour }

	got := serv.serve()

	if got != nil {
		t.Errorf("got error: %v", got.Error())
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

func withCall(t *testing.T) {
	serv := NewGoldPriceService(
		reqWithCall{},
		sendWithCall{},
		extWithCall{},
		valWithCall{})

	serv.serve()

	if !reqCall {
		t.Errorf("the method for requesting a page is not called in the app layer")
	}
	if !extCall {
		t.Errorf("the method for extracting a price is not called in the app layer")
	}
	if !valCall {
		t.Errorf("the method for validating the price is not called in the app layer")
	}
	if !sendCall {
		t.Errorf("the method for sending the price is not called in the app layer")
	}
}
