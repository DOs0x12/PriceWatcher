package price

import "PriceWatcher/internal/interfaces/configer"

type AddItemComm struct {
	configer configer.Configer
}

func NewAddItemComm(configer configer.Configer) AddItemComm {
	return AddItemComm{configer: configer}
}

const successfullMessage = "The item is added for watching"

func (c AddItemComm) AddItemToWatch(itemValue string) string {
	return successfullMessage
}
