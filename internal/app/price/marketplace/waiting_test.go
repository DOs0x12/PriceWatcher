package marketplace

import (
	"testing"
	"time"
)

func TestGetCallTime(t *testing.T) {
	now := time.Now()
	testMinute := 25
	testNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	got := getCallTime(testNow)

	wantMinute := 30
	want := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), wantMinute, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testMinute = 30
	testNow = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	got = getCallTime(testNow)

	want = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	testMinute = 31
	testNow = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), testMinute, 0, 0, now.Location())

	got = getCallTime(testNow)

	want = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

	if want.Compare(got) != 0 {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
