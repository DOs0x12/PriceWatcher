package domain

import "io"

type Processor struct{}

func (p Processor) Process(body io.ReadCloser) float32 {
	return 5.55
}
