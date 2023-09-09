package file

type WriteReader interface {
	Write(price float32) error
	Read() (float32, error)
}
