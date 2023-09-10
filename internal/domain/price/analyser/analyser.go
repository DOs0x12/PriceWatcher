package analyser

type Analyser interface {
	AnalysePrice(price, initialPrice float32) (changed, up bool, amount float32)
}

type MarketplaceAnalyser struct{}

func (a MarketplaceAnalyser) AnalysePrice(price, initialPrice float32) (changed, up bool, amount float32) {
	if price > initialPrice {
		return true, true, price - initialPrice
	}

	if price < initialPrice {
		return true, false, initialPrice - price
	}

	return false, false, 0.0
}
