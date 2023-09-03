package analyser

type Analyser interface {
	IsPriceChanged(price float32) (changed, up bool, amount float32)
}
