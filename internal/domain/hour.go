package domain

type HourValidator interface {
	Validate(hour int) bool
}

type MessageHourVal struct{}

func (MessageHourVal) Validate(hour int) bool {
	hours := []int{12, 17}

	for _, h := range hours {
		if h == hour {
			return true
		}
	}

	return false
}
