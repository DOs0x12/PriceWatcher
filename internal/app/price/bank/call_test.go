package bank

import (
	"testing"
	"time"
)

func TestGetCallTime(t *testing.T) {
	now := time.Now()
	testHour := 5
	testNow := time.Date(now.Year(), now.Month(), now.Day(), testHour, 1, 0, 0, now.Location())
	callHours := []int{15, 16}

	got := getCallTime(testNow, callHours)

	want := time.Date(now.Year(), now.Month(), now.Day(), callHours[0], 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testHour = 15
	testNow = time.Date(now.Year(), now.Month(), now.Day(), testHour, 0, 0, 0, now.Location())

	got = getCallTime(testNow, callHours)

	want = time.Date(now.Year(), now.Month(), now.Day(), callHours[1], 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testHour = 16
	testNow = time.Date(now.Year(), now.Month(), now.Day(), testHour, 0, 0, 0, now.Location())

	got = getCallTime(testNow, callHours)

	want = time.Date(now.Year(), now.Month(), now.Day()+1, callHours[0], 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
