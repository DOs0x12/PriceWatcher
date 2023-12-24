package price

import (
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
	items, err := c.wr.Read()
	if err != nil {
		logrus.Errorf("Can not get the current prices: %v", err)

		return ""
	}

	return formMessage(items)
}

func formMessage(items map[string]float64) string {
	builder := strings.Builder{}

	for name, val := range items {
		builder.WriteString(name + ": ")
		builder.WriteString(fmt.Sprintf("%.2f\n", val))
	}

	return builder.String()
}
