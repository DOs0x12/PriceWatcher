package item

import (
	"PriceWatcher/internal/interfaces/configer"
	"strings"
)

type AddItemComm struct {
	configer configer.Configer
}

func NewAddItemComm(configer configer.Configer) AddItemComm {
	return AddItemComm{configer: configer}
}

const successfullMessage = "The item is added for watching"

func (c AddItemComm) AddItemToWatch(itemValue string) string {
	values := strings.Split(itemValue, " ")
	_ = c.configer.AddItemToWatch(values[0], values[1], values[2])

	return successfullMessage
}
