package daytime_test

import (
	"testing"
	"time"

	"github.com/claudemuller/daytime/pkg/daytime"
)

func TestGetDayTime(t *testing.T) {
	tests := []struct {
		name     string
		mockTime daytime.Clockfn
		want     string
	}{
		{
			"Monday 22nd February 1982 17:37:43 PST",
			func() time.Time {
				tz := time.FixedZone("PST", -8*3600)
				return time.Date(1982, 2, 22, 17, 37, 43, 0, tz)
			},
			"Monday, February 22, 1982 17:37:43-PST",
		},
		{
			"Tuesday 29th February 2000 10:11:33 SAST",
			func() time.Time {
				tz := time.FixedZone("SAST", -8*3600)
				return time.Date(2000, 2, 29, 10, 11, 33, 0, tz)
			},
			"Tuesday, February 29, 2000 10:11:33-SAST",
		},
		{
			"Wednesday 1st January 2025 00:00:00 UTC",
			func() time.Time {
				return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			},
			"Wednesday, January 1, 2025 00:00:00-UTC",
		},
	}

	t.Log("Given that we want the current datetime.")
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Test: %s", tc.name)

			got := daytime.Get(tc.mockTime)

			if got != tc.want {
				t.Fatalf("\t\twant = %s, got = %s", tc.want, got)
			}
		})
	}

}
