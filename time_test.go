// Copyright (c) 2017-2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package elapsed

import (
	"fmt"
	"testing"
	"time"
)

func TestAddTranslation(t *testing.T) {
	var dt = []struct {
		lang string
		tr   Terms
		err  error
	}{
		{lang: "", err: ErrISOCode},
		{lang: "fr", err: ErrExists},
		{lang: "ru", tr: Terms{Yesterday: "euh"}, err: ErrIncomplete},
		{lang: "en-gb", tr: Terms{
			NotYet:    `not yet`,
			JustNow:   `just now`,
			Minute:    `1 minute ago`,
			Minutes:   `%d minutes ago`,
			Hour:      `1 hour ago`,
			Hours:     `%d hours ago`,
			Yesterday: `yesterday`,
			Days:      `%d days ago`,
			Week:      `1 weeks ago`,
			Weeks:     `%d weeks ago`,
			Month:     `1 months ago`,
			Months:    `%d months ago`,
			Year:      `1 years ago`,
			Years:     `%d years ago`,
		}},
	}
	for i, tt := range dt {
		if err := AddTranslation(tt.lang, tt.tr); err != tt.err {
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
		{time.Now().Add(-time.Hour * 24 * 6), "6 days ago"},
		{time.Now().Add(-(time.Hour * 24 * 6) - 2*time.Hour), "6 days ago"},
		{time.Now().Add(-time.Hour * 24 * 3), "3 days ago"},
		{time.Now().Add(-time.Hour * 24 * 7), "1 week ago"},
		{time.Now().Add(-time.Hour * 24 * 14), "2 weeks ago"},
		// 4 weeks == 1 month
		{time.Now().Add(-time.Hour * 24 * 28), "1 month ago"},
		{time.Now().Add(-time.Hour * 24 * 60), "2 months ago"},
		// 12 months == 1 year
		{time.Now().Add(-time.Hour * 24 * 360), "1 year ago"},
		{time.Now().Add(-time.Hour * 24 * 365 * 3), "3 years ago"},
	}
	for i, tt := range dt {
		// Requests an unknown language.
		if out := LocalTime(tt.in, "ru"); out != tt.out {
			t.Errorf("%d. content mismatch for %v: exp=%q got=%q", i, tt.in, tt.out, out)
		}
	}
}

func ExampleTime() {
	t := time.Now().Add(-time.Hour)
	fmt.Println(Time(t))

	t = time.Now().Add(-time.Hour * 24 * 3)
	fmt.Println(Time(t))

	t, _ = time.Parse("2006-02-01", "2049-08-19")
	fmt.Println(Time(t))

	t = time.Now().Add(-time.Hour * 24 * 3)
	fmt.Println(LocalTime(t, "fr"))
	// Output: 1 hour ago
	// 3 days ago
	// not yet
	// il y a 3 jours
}
