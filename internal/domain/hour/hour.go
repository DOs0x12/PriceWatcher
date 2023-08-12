package hour

type HourValidator interface {
	Validate(hour int, sendHours []int) bool
}

type MessageHourVal struct{}

func (MessageHourVal) Validate(hour int, sendHours []int) bool {
	for _, h := range sendHours {
		if h == hour {
			return true
		}
	}

	return false
}
