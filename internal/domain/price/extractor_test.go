package price

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

var ext = PriceExtractor{}

func TestExtractPrice(t *testing.T) {
	logrus.Info("Start to test the func ExtractPrice with true value")
	withTrueValue(t)
	logrus.Info("Start to test the func ExtractPrice for errors")
	testErrors(t)
}

type testReadCloser struct {
	Reader io.Reader
}

func (testReadCloser) Close() error                       { return nil }
func (t testReadCloser) Read(p []byte) (n int, err error) { return t.Reader.Read(p) }

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Some error!")
}

func withTrueValue(t *testing.T) {
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
	r := strings.NewReader(s)
	rc := testReadCloser{r}
	var want float32 = 3232.00
	got, err := ext.ExtractPrice(rc)
	if err != nil {
		t.Errorf("Got error: %v, wanted %v", err, want)
	}
	if got != want && err == nil {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func testErrors(t *testing.T) {
	s := `
		<html>
			<head>
				<title>test</title>
			</head>
			<body>
				<table>
					<td> покупка: wrong value</td>
				</table>
			</body>
		</html>`
	r := strings.NewReader(s)
	rc := testReadCloser{r}
	errTemplt := "the document does not have a price value with the tag:"
	want := float32(0.00)
	got, err := ext.ExtractPrice(rc)
	if err != nil && !strings.Contains(err.Error(), errTemplt) {
		t.Errorf("Got not wanted error: %v, wanted error template: %v", err, errTemplt)
	}
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	errTemplt = "cannot parse the body to an HTML document:"
	rc.Reader = errReader{}
	got, err = ext.ExtractPrice(rc)
	if err != nil && !strings.Contains(err.Error(), errTemplt) {
		t.Errorf("Got not wanted error: %v, wanted error template: %v", err, errTemplt)
	}
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
