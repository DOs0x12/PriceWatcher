package analyser

type Analyser interface {
	AnalysePrice(price float32) (changed, up bool, amount float32)
}

type MarketplaceAnalyser struct {
	currentMinPrice float32 `default:"0.0"`
}

func (a MarketplaceAnalyser) AnalysePrice(price float32) (changed, up bool, amount float32) {
	if price > a.currentMinPrice {
		return true, true, price - a.currentMinPrice
	}

	if price < a.currentMinPrice {
		a.currentMinPrice = price

		return true, false, a.currentMinPrice - price
	}

	return false, false, 0.0
}
