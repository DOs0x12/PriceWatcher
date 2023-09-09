package analyser

type Analyser interface {
	AnalysePrice(price float32) (changed, up bool, amount float32)
}

type MarketplaceAnalyser struct {
	CurrentMinPrice float32 `default:"0.0"`
}

func (a MarketplaceAnalyser) AnalysePrice(price float32) (changed, up bool, amount float32) {
	if price > a.CurrentMinPrice {
		return true, true, price - a.CurrentMinPrice
	}

	if price < a.CurrentMinPrice {
		a.CurrentMinPrice = price

		return true, false, a.CurrentMinPrice - price
	}

	return false, false, 0.0
}
