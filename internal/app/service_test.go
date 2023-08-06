package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/entities"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
)

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

func (s testSender) Send(price float32) error { return errors.New("HOHO") }

func TestServe(t *testing.T) {
	serv := GoldPriceService{
		req:    testRequester{},
		sender: testSender{},
		ext:    domain.PriceExtractor{},
		val:    domain.MessageHourVal{}}
	patches := gomonkey.ApplyFunc(time.Time.Hour, func(time.Time) int {
		return 12
	})
	defer patches.Reset()
	var want error = nil
	got := serv.serve()

	if want != got {
		t.Errorf("got error: %v", got)
	}

}

// type testProcessor struct{}

// func (p testProcessor) Process(page *string) float32 {
// 	return 55.55
// }

// func TestHandleGoldPrice(t *testing.T) {
// 	HandleGoldPrice(testRequester{}, testProcessor{}, testSender{})
// }
