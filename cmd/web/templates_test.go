package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2024, 10, 17, 3, 41, 59, 0, time.UTC)
	hd := humanDate(tm)

	expect := "17 Oct 2024 at 03:41"
	if hd != expect {
		t.Errorf("Got %q, expect %q\n", hd, expect)
	}
}
