package time

import (
	"fmt"
	"testing"
	"time"
)

type testTimeParam struct {
	periodMin  int
	testMin    int
	testSecond int
	wantMin    int
	wantSec    int
}

func TestGetWaitTime(t *testing.T) {
	par := testTimeParam{
		periodMin:  45,
		testMin:    55,
		testSecond: 31,
		wantMin:    49,
		wantSec:    29,
	}
	testWhenToSendRep(t, par)

	par.testMin = 46
	par.testSecond = 0
	par.wantMin = 59
	par.wantSec = 0
	testWhenToSendRep(t, par)

	par.testMin = 45
	par.testSecond = 0
	par.wantMin = 0
	par.wantSec = 0
	testWhenToSendRep(t, par)

	par.testMin = 44
	par.testSecond = 0
	par.wantMin = 1
	par.wantSec = 0
	testWhenToSendRep(t, par)

	par.testMin = 44
	par.testSecond = 5
	par.wantMin = 0
	par.wantSec = 55
	testWhenToSendRep(t, par)
}

func testWhenToSendRep(t *testing.T, par testTimeParam) {
	nT := time.Now()

	testNow :=
		time.Date(nT.Year(), nT.Month(), nT.Day(), nT.Hour(), par.testMin, par.testSecond, nT.Nanosecond(), nT.Location())

	durStr := fmt.Sprintf("%vm%vs", par.wantMin, par.wantSec)

	want, err := time.ParseDuration(durStr)
	if err != nil {
		t.Errorf("An error occurs while parsing duration in the test: %v", err)
	}

	got, err := WaitPerStart(testNow, par.periodMin)
	if err != nil {
		t.Errorf("The method retuns an error: %v", err)
	}

	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
