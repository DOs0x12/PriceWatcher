package interfaces

import "io"

type Processor interface {
	Process(body io.ReadCloser) float32
}
