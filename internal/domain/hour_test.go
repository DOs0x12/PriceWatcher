package domain

import "testing"

func TestValidate(t *testing.T) {
	val := MessageHourVal{}

	tHours := [4]int{11, 13, 16, 18}
	want := false
	for i := 0; i < len(tHours); i++ {
		if got := val.Validate(tHours[i]); got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	}

	tHour := 12
	want = true
	if got := val.Validate(tHour); got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

	tHour = 17
	want = true
	if got := val.Validate(tHour); got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
