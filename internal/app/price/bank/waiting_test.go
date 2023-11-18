package bank

import (
	"fmt"
	"testing"
	"time"
)

func TestGetWaitDurWithRandomComp(t *testing.T) {
	now := time.Now()
	testHour := 5
	testCallHour := testHour + 1
	testNow := time.Date(now.Year(), now.Month(), now.Day(), testHour, 0, 0, 0, now.Location())
	callTime := time.Date(now.Year(), now.Month(), now.Day(), testCallHour, 0, 0, 0, now.Location())
	randDur, _ := time.ParseDuration("0h")

	got := getWaitDurWithRandomComp(testNow, callTime, randDur)

	subHours := testCallHour - testHour
	want, _ := time.ParseDuration(fmt.Sprintf("%vh", subHours))
	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	got = getWaitDurWithRandomComp(testNow, testNow, randDur)

	want, _ = time.ParseDuration(fmt.Sprintf("%vh", 0))
	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	got = getWaitDurWithRandomComp(callTime, testNow, randDur)

	want, _ = time.ParseDuration(fmt.Sprintf("%vh", 0))
	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	randDur, _ = time.ParseDuration("2h")

	got = getWaitDurWithRandomComp(testNow, callTime, randDur)

	want, _ = time.ParseDuration(fmt.Sprintf("%vh", subHours))
	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	randMinutes := 30
	randDur, _ = time.ParseDuration(fmt.Sprintf("%vm", randMinutes))

	got = getWaitDurWithRandomComp(testNow, callTime, randDur)

	want, _ = time.ParseDuration(fmt.Sprintf("%vm", randMinutes))
	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
