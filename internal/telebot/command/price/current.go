package price

type CurrentPriceComm struct{}

func NewPriceCommand() CurrentPriceComm {
	return CurrentPriceComm{}
}

func (c CurrentPriceComm) GetCurrentPrices() string {
	return ""
}
