package app

import (
	"fmt"
	"testing"
	"time"
)

func TestGetWaitTime(t *testing.T) {
	nT := time.Now()
	testMin := 45
	testSec := 45
	testNow :=
		time.Date(nT.Year(), nT.Month(), nT.Day(), nT.Hour(), testMin, testSec, nT.Nanosecond(), nT.Location())

	waitMin := 60 - testMin
	waitSec := 60 - testSec

	durStr := fmt.Sprintf("%vm%vs", waitMin, waitSec)
	want, err := time.ParseDuration(durStr)
	if err != nil {
		t.Errorf("An error occurs while parsing duration in the test: %v", err)
	}

	got, err := getWaitTime(testNow)
	if err != nil {
		t.Errorf("The getWaitTime method retuns an error: %v", err)
	}

	if want != got {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
