package price

import (
	"PriceWatcher/internal/entities/price"
	"PriceWatcher/internal/interfaces/file"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type PriceCommand struct {
	wr file.WriteReader
}

func NewPriceCommand(wr file.WriteReader) PriceCommand {
	return PriceCommand{wr: wr}
}

func (c PriceCommand) GetCurrentPrices() string {
	c.wr.Lock()
	items, err := c.wr.ReadPrices()
	c.wr.Unlock()
	if err != nil {
		logrus.Errorf("Can not get the current prices: %v", err)

		return ""
	}

	return formMessage(items)
}

func formMessage(items map[string]price.ItemPrice) string {
	builder := strings.Builder{}

	for name, val := range items {
		builder.WriteString(name + ": ")
		builder.WriteString(fmt.Sprintf("%.2f\n%v\n\n", val.Price, val.Address))
	}

	return builder.String()
}
