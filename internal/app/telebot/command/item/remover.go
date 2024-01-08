package item

import (
	"PriceWatcher/internal/interfaces/configer"
	"strings"
)

type RemoveItemComm struct {
	configer configer.Configer
}

func NewRemoveItemComm(configer configer.Configer) RemoveItemComm {
	return RemoveItemComm{configer: configer}
}

const successfullMessageOfRemoving = "The item is removed for watching"

func (c RemoveItemComm) RemoveItemToWatch(itemValue string) string {
	values := strings.Split(itemValue, " ")
	_ = c.configer.RemoveItemFromWatching(values[0], values[1])

	return successfullMessageOfRemoving
}
