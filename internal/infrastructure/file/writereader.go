package file

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

var fileName = "last_price.txt"

type WriteReader struct{}

func (WriteReader) Write(price float32) error {
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

func (WriteReader) Read() (float32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0.0, nil
		}

		return 0.0, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return 0.0, fmt.Errorf("cannot read a price from the file %v: %v", fileName, err)
		}
		if n == 0 {
			break
		}
	}

	bits := binary.LittleEndian.Uint32(buf)
	pr := math.Float32frombits(bits)

	return pr, nil
}
