package interfaces

import "io"

type Extractor interface {
	ExtractPrice(body io.Reader) (float32, error)
}
