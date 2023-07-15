package interfaces

type Processor interface {
	Process(page *string) float32
}
