package app

// import (
// 	"GoldPriceGetter/internal/entities"
// 	"io"
// 	"strings"
// 	"testing"
// )

// type testProcessor struct{}

// func (p testProcessor) Process(page *string) float32 {
// 	return 55.55
// }

// type testRequester struct{}

// func (r testRequester) RequestPage() entities.Response {
// 	reader := io. strings.NewReader("test55.55")
// 	return entities.Response{Body: reader}
// }

// type testSender struct{}

// func (s testSender) Send(price float32) {

// }

// func TestHandleGoldPrice(t *testing.T) {
// 	HandleGoldPrice(testRequester{}, testProcessor{}, testSender{})
// }
