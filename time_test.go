// Copyright (c) 2017-2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package elapsed_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rvflash/elapsed"
)

func TestAddTranslation(t *testing.T) {
	var dt = []struct {
		lang string
		tr   elapsed.Terms
		err  error
	}{
		{lang: "", err: elapsed.ErrISOCode},
		{lang: "fr", err: elapsed.ErrExists},
		{lang: "ru", tr: elapsed.Terms{elapsed.Yesterday: "euh"}, err: elapsed.ErrIncomplete},
		{lang: "en-gb", tr: elapsed.Terms{
			elapsed.NotYet:     `not yet`,
			elapsed.JustNow:    `just now`,
			elapsed.LastMinute: `1 minute ago`,
			elapsed.Minutes:    `%d minutes ago`,
			elapsed.LastHour:   `1 hour ago`,
			elapsed.Hours:      `%d hours ago`,
			elapsed.Yesterday:  `yesterday`,
			elapsed.Days:       `%d days ago`,
			elapsed.Weeks:      `%d weeks ago`,
			elapsed.Months:     `%d months ago`,
			elapsed.Years:      `%d years ago`,
		}},
	}
	for i, tt := range dt {
		if err := elapsed.AddTranslation(tt.lang, tt.tr); err != tt.err {
			t.Errorf("%d. error mismatch: exp=%q got=%q", i, tt.err, err)
		}
	}
}

func TestLocalTime(t *testing.T) {
	var dt = []struct {
		in  time.Time
		out string
	}{
		{time.Time{}, "not yet"},
		{time.Now().Add(time.Hour), "not yet"},
		{time.Now(), "just now"},
		{time.Now().Add(-time.Minute), "1 minute ago"},
		{time.Now().Add(-time.Minute * 40), "40 minutes ago"},
		{time.Now().Add(-time.Hour), "1 hour ago"},
		{time.Now().Add(-time.Hour * 3), "3 hours ago"},
		{time.Now().Add(-time.Hour * 32), "yesterday"},
		{time.Now().Add(-time.Hour * 24 * 3), "3 days ago"},
		{time.Now().Add(-time.Hour * 24 * 14), "2 weeks ago"},
		{time.Now().Add(-time.Hour * 24 * 60), "2 months ago"},
		{time.Now().Add(-time.Hour * 24 * 365 * 3), "3 years ago"},
	}
	for i, tt := range dt {
		// Requests an unknown language.
		if out := elapsed.LocalTime(tt.in, "ru"); out != tt.out {
			t.Errorf("%d. content mismatch for %v: exp=%q got=%q", i, tt.in, tt.out, out)
		}
	}
}

func ExampleTime() {
	t := time.Now().Add(-time.Hour)
	fmt.Println(elapsed.Time(t))

	t = time.Now().Add(-time.Hour * 24 * 3)
	fmt.Println(elapsed.Time(t))

	t, _ = time.Parse("2006-02-01", "2049-08-19")
	fmt.Println(elapsed.Time(t))

	t = time.Now().Add(-time.Hour * 24 * 3)
	fmt.Println(elapsed.LocalTime(t, "fr"))
	// Output: 1 hour ago
	// 3 days ago
	// not yet
	// il y a 3 jours
}
