package time

import (
	"fmt"
	"testing"
	"time"
)

type testTimeParam struct {
	nt         NearestTime
	testMin    int
	testSecond int
	wantMin    int
	wantSec    int
}

func TestGetWaitTime(t *testing.T) {
	par := testTimeParam{
		nt:         Hour,
		testMin:    55,
		testSecond: 31,
		wantMin:    4,
		wantSec:    29,
	}
	testWhenToSendRep(t, par)

	par.nt = HalfHour
	par.testMin = 21
	par.testSecond = 0
	par.wantMin = 9
	par.wantSec = 0
	testWhenToSendRep(t, par)

	par.testMin = 31
	par.testSecond = 5
	par.wantMin = 28
	par.wantSec = 55
	testWhenToSendRep(t, par)
}

func testWhenToSendRep(t *testing.T, par testTimeParam) {
	nT := time.Now()

	testNow :=
		time.Date(nT.Year(), nT.Month(), nT.Day(), nT.Hour(), par.testMin, par.testSecond, 0, nT.Location())

	durStr := fmt.Sprintf("%vm%vs", par.wantMin, par.wantSec)

	want, err := time.ParseDuration(durStr)
	if err != nil {
		t.Errorf("An error occurs while parsing duration in the test: %v", err)
	}

	got := PerStartDur(testNow, par.nt)

	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
