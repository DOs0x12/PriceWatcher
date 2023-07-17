package entities

import "io"

type Response struct {
	Body io.ReadCloser
}
