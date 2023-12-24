package file

import (
	"PriceWatcher/internal/entities/price"
	"errors"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var fileName = "last_price.yaml"

type WriteReader struct {
	mu *sync.Mutex
}

func NewWR() WriteReader {
	return WriteReader{mu: &sync.Mutex{}}
}

type ItemPriceDto struct {
	Address string  `yaml:"address"`
	Price   float64 `yaml:"price"`
}

func (WriteReader) WritePrices(prices map[string]price.ItemPrice) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create a file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	curPrices := castTo(prices)

	yVal, err := yaml.Marshal(curPrices)
	if err != nil {
		return fmt.Errorf("cannot marshall price values to an yaml value: %v", err)
	}

	if _, err := file.Write(yVal); err != nil {
		return fmt.Errorf("cannot write a price to the file %v: %v", fileName, err)
	}

	return nil
}

func castTo(prices map[string]price.ItemPrice) map[string]ItemPriceDto {
	priceDtos := make(map[string]ItemPriceDto)

	for k, v := range prices {
		priceDtos[k] = ItemPriceDto{Address: v.Address, Price: v.Price}
	}

	return priceDtos
}

func (WriteReader) ReadPrices() (map[string]price.ItemPrice, error) {
	itemPrices := make(map[string]price.ItemPrice)
	file, err := os.ReadFile(fileName)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		return itemPrices, nil
	}

	if err != nil {
		return nil, fmt.Errorf("cannot read the file %v: %v", fileName, err)
	}

	itemPriceDtos := make(map[string]ItemPriceDto)

	if err := yaml.Unmarshal(file, &itemPriceDtos); err != nil {
		return nil, fmt.Errorf("cannot unmarshall price values from the file %v: %v", fileName, err)
	}

	return castFrom(itemPriceDtos), nil
}

func castFrom(priceDtos map[string]ItemPriceDto) map[string]price.ItemPrice {
	prices := make(map[string]price.ItemPrice)

	for k, v := range priceDtos {
		prices[k] = price.ItemPrice{Address: v.Address, Price: v.Price}
	}

	return prices
}

func (wr WriteReader) Lock() {
	wr.mu.Lock()
}

func (wr WriteReader) Unlock() {
	wr.mu.Unlock()
}
