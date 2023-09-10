package file

type WriteReader interface {
	Write(prices map[string]float64) error
	Read() (map[string]float64, error)
}
