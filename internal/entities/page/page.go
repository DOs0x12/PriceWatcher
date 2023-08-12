package page

import "io"

type Response struct {
	Body io.ReadCloser
}
