package file

type WriteReader interface {
	Write(data float32)
	Read() float32
}
