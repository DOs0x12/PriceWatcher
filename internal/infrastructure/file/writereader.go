package file

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var fileName = "last_price.yaml"

type WriteReader struct{}

func (WriteReader) Write(prices map[string]float64) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create a file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	yVal, err := yaml.Marshal(prices)
	if err != nil {
		return fmt.Errorf("cannot marshall price values to an yaml value: %v", err)
	}

	if _, err := file.Write(yVal); err != nil {
		return fmt.Errorf("cannot write a price to the file %v: %v", fileName, err)
	}

	return nil
}

func (WriteReader) Read() (map[string]float64, error) {
	itemPrices := make(map[string]float64)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read the file %v: %v", fileName, err)
	}

	if err := yaml.Unmarshal(file, &itemPrices); err != nil {
		return nil, fmt.Errorf("cannot unmarshall price values from the file %v: %v", fileName, err)
	}

	return itemPrices, nil
}
