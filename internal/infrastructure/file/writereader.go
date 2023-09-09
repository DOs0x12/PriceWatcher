package file

import (
	"fmt"
	"os"
	"strconv"
)

var fileName = "last_price.txt"

type WriteReader struct{}

func Write(price float32) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create a file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := file.WriteString(strconv.FormatFloat(float64(price), 'f', 2, 32)); err != nil {
		return fmt.Errorf("cannot write a price to the file %v: %v", fileName, err)
	}

	return nil
}

func Read() (float32, error) {
	return 0.0, nil
}
