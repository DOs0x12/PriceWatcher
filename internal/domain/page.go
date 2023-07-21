package domain

import "io"

type Processor interface {
	Process(body io.ReadCloser) float32
}

type PageProcessor struct{}

func (p PageProcessor) Process(body io.ReadCloser) float32 {
	return 5.55
}
