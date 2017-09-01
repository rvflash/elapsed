// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package elapsed_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rvflash/elapsed"
)

func TestTime(t *testing.T) {
	var dt = []struct {
		in  time.Time
		out string
	}{
		{time.Time{}, elapsed.NotYet},
		{time.Now().Add(time.Hour), elapsed.NotYet},
		{time.Now(), elapsed.JustNow},
		{time.Now().Add(-time.Minute), elapsed.LastMinute},
		{time.Now().Add(-time.Minute * 40), fmt.Sprintf(elapsed.Minutes, 40)},
		{time.Now().Add(-time.Hour), elapsed.LastHour},
		{time.Now().Add(-time.Hour * 3), fmt.Sprintf(elapsed.Hours, 3)},
		{time.Now().Add(-time.Hour * 32), elapsed.Yesterday},
		{time.Now().Add(-time.Hour * 24 * 3), fmt.Sprintf(elapsed.Days, 3)},
		{time.Now().Add(-time.Hour * 24 * 14), fmt.Sprintf(elapsed.Weeks, 2)},
		{time.Now().Add(-time.Hour * 24 * 60), fmt.Sprintf(elapsed.Months, 2)},
		{time.Now().Add(-time.Hour * 24 * 365 * 3), fmt.Sprintf(elapsed.Years, 3)},
	}
	for i, tt := range dt {
		if out := elapsed.Time(tt.in); out != tt.out {
			t.Errorf("%d. content mismatch for %v:exp=%q got=%q", i, tt.in, tt.out, out)
		}
	}
}

func ExampleTime() {
	t := time.Now().Add(-time.Hour)
	fmt.Println(elapsed.Time(t))
	// Output: 1 hour ago
}
