package app

import (
	"GoldRateGetter/internal/entities"
	"testing"
)

type testProcessor struct{}

func (p testProcessor) Process(page string) float32 {
	return 55.55
}

type testRequester struct{}

func (r testRequester) RequestPage() entities.Response {
	return entities.Response{Page: "test55.55"}
}

type testSender struct{}

func (s testSender) Send(rate float32) {

}

func TestHandleGoldRate(t *testing.T) {
	HandleGoldRate(testRequester{}, testProcessor{}, testSender{})
}
