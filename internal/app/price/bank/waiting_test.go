package bank

import (
	"testing"
	"time"
)

var testWNow time.Time

type testWClock struct{}

func (testWClock) Now() time.Time                         { return testWNow }
func (testWClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

func TestGetCallTime(t *testing.T) {
	now := time.Now()
	tW := testWClock{}
	testHour := 5
	testWNow = time.Date(now.Year(), now.Month(), now.Day(), testHour, 1, 0, 0, now.Location())
	callHours := []int{15, 16}

	got := getCallTime(tW.Now(), callHours)

	want := time.Date(now.Year(), now.Month(), now.Day(), callHours[0], 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testHour = 15
	testWNow = time.Date(now.Year(), now.Month(), now.Day(), testHour, 0, 0, 0, now.Location())

	got = getCallTime(tW.Now(), callHours)

	want = time.Date(now.Year(), now.Month(), now.Day(), callHours[1], 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testHour = 16
	testWNow = time.Date(now.Year(), now.Month(), now.Day(), testHour, 0, 0, 0, now.Location())

	got = getCallTime(tW.Now(), callHours)

	want = time.Date(now.Year(), now.Month(), now.Day()+1, callHours[0], 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
