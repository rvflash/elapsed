// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package elapsed

import (
	"fmt"
	"math"
	"time"
)

// Texts to be translated if necessary.
var (
	NotYet     = `not yet`
	JustNow    = `just now`
	LastMinute = `1 minute ago`
	Minutes    = `%d minutes ago`
	LastHour   = `1 hour ago`
	Hours      = `%d hours ago`
	Yesterday  = `Yesterday`
	Days       = `%d days ago`
	Weeks      = `%d weeks ago`
	Months     = `%d months ago`
	Years      = `%d years ago`
)

// Time returns in a human readable format the elapsed time
// since the given datetime.
func Time(t time.Time) string {
	if t.IsZero() || time.Now().Before(t) {
		return NotYet
	}
	diff := time.Since(t)
	// Duration in seconds
	s := diff.Seconds()
	// Duration in days
	d := int(s / 86400)
	switch {
	case s < 60:
		return JustNow
	case s < 120:
		return LastMinute
	case s < 3600:
		return fmt.Sprintf(Minutes, int(diff.Minutes()))
	case s < 7200:
		return LastHour
	case s < 86400:
		return fmt.Sprintf(Hours, int(diff.Hours()))
	case d == 1:
		return Yesterday
	case d < 7:
		return fmt.Sprintf(Days, d)
	case d < 31:
		return fmt.Sprintf(Weeks, int(math.Ceil(float64(d)/7)))
	case d < 365:
		return fmt.Sprintf(Months, int(math.Ceil(float64(d)/30)))
	default:
		return fmt.Sprintf(Years, int(math.Ceil(float64(d)/365)))
	}
}
