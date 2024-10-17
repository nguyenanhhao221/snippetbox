package main

import (
	"testing"
	"time"

	"snippetbox.haonguyen.tech/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name   string
		tm     time.Time
		expect string
	}{
		{name: "UTC", tm: time.Date(2024, 10, 17, 3, 41, 59, 0, time.UTC), expect: "17 Oct 2024 at 03:41"},
		{name: "Empty", tm: time.Time{}, expect: ""},
		{name: "CET", tm: time.Date(2024, 10, 17, 3, 41, 59, 0, time.FixedZone("CET", 1*60*60)), expect: "17 Oct 2024 at 02:41"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hd := humanDate(tc.tm)
			assert.Equal(t, hd, tc.expect)
		})
	}
}
