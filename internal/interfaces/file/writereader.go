package file

type WriteReader interface {
	WritePrices(prices map[string]float64) error
	ReadPrices() (map[string]float64, error)
}
