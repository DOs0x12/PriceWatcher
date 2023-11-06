package marketplace

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
	testMinute := 25
	testWNow = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	got := getCallTime(tW.Now())

	wantMinute := 30
	want := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), wantMinute, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testMinute = 30
	testWNow = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	got = getCallTime(tW.Now())

	want = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testMinute = 31
	testWNow = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	got = getCallTime(tW.Now())

	want = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
