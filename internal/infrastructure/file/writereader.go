package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var fileName = "last_price.txt"
var itemPrices = make(map[string]float64)

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

	for k, v := range prices {
		val := strconv.FormatFloat(float64(v), 'f', 2, 32)
		data := fmt.Sprintf("%v %v\n", k, val)

		if _, err := file.WriteString(data); err != nil {
			return fmt.Errorf("cannot write a price to the file %v: %v", fileName, err)
		}
	}

	return nil
}

func (WriteReader) Read() error {
	file, err := os.Open(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = addItemPrice(itemPrices, scanner.Text())
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("cannot read a price from the file %v: %v", fileName, err)
	}

	return nil
}

func addItemPrice(prices map[string]float64, dataLine string) error {
	dataTokenCount := 2
	data := strings.Split(dataLine, " ")
	if len(data) < dataTokenCount {
		return fmt.Errorf("data \"%v\" in the file %v is in not appropriate format", data, fileName)
	}

	price, err := strconv.ParseFloat(data[1], 32)
	if err != nil {
		return err
	}

	prices[data[0]] = price

	return nil
}
