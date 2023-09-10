package file

type WriteReader interface {
	Write(prices map[string]float64) error
	Read() (float32, error)
}
