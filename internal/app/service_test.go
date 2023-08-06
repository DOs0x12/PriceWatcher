package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/entities"
	"io"
	"strings"
	"testing"
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

func (s testSender) Send(price float32) error { return nil }

func TestServe(t *testing.T) {
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
