// Copyright (c) 2017-2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package elapsed return the elapsed time since a given time in a human readable format.
package elapsed

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// TrID is the ID of a translation.
type TrID int

const (
	// NotYet is the translation ID for the "not yet" text.
	NotYet TrID = iota
	// JustNow is the translation ID for the "just now" text.
	JustNow
	// Minute is the translation ID for the "1 minute ago" text.
	Minute
	// Minutes is the translation ID for the "%d minutes ago" text.
	Minutes
	// Hour is the singular of Hours
	Hour
	// Hours is the translation ID for the "%d hours ago" text.
	Hours
	// Yesterday is the translation ID for the "yesterday" text.
	Yesterday
	// Day is the singular of Days
	Day
	// Days is the translation ID for the "%d days ago" text.
	Days
	// Week is the singular of Weeks
	Week
	// Weeks is the translation ID for the "`%d weeks ago" text.
	Weeks
	// Month is the singular of Months
	Month
	// Months is the translation ID for the "%d months ago" text.
	Months
	// Year is the singular
	Year
	// Years is the translation ID for the "%d years ago" text.
	Years
)

// Lists all translations by identifier.
type Terms map[TrID]string

// Lists all translations by language code.
type Translation map[string]Terms

// i18n is a map of translations by language code.
var i18n = Translation{
	"de": {
		NotYet:    `noch nicht`,
		JustNow:   `im Moment`,
		Minute:    `vor %d Minute`,
		Minutes:   `vor %d Minuten`,
		Hour:      `vor %d Stunde`,
		Hours:     `vor %d Stunden`,
		Yesterday: `gestern`,
		Day:       `vor %d Tag`,
		Days:      `vor %d Tagen`,
		Week:      `vor %d Woche`,
		Weeks:     `vor %d Wochen`,
		Month:     `vor %d Monat`,
		Months:    `vor %d Monaten`,
		Year:      `vor %d Jahr`,
		Years:     `vor %d Jahren`,
	},
	"en": {
		NotYet:    `not yet`,
		JustNow:   `just now`,
		Minute:    `%d minute ago`,
		Minutes:   `%d minutes ago`,
		Hour:      `%d hour ago`,
		Hours:     `%d hours ago`,
		Yesterday: `yesterday`,
		Day:       `%d day ago`,
		Days:      `%d days ago`,
		Week:      `%d week ago`,
		Weeks:     `%d weeks ago`,
		Month:     `%d month ago`,
		Months:    `%d months ago`,
		Year:      `%d year ago`,
		Years:     `%d years ago`,
	},
	"es": {
		NotYet:    `aún no`,
		JustNow:   `al instante`,
		Minute:    `hace %d minuto`,
		Minutes:   `hace %d minutos`,
		Hour:      `hace %d hora`,
		Hours:     `hace %d horas`,
		Yesterday: `ayer`,
		Day:       `hace %d día`,
		Days:      `hace %d días`,
		Week:      `hace %d semana`,
		Weeks:     `hace %d semanas`,
		Month:     `hace %d mes`,
		Months:    `hace %d meses`,
		Year:      `hace %d año`,
		Years:     `hace %d años`,
	},
	"fr": {
		NotYet:    `pas encore`,
		JustNow:   `à l'instant`,
		Minute:    `il y a %d minute`,
		Minutes:   `il y a %d minutes`,
		Hour:      `il y a %d heure`,
		Hours:     `il y a %d heures`,
		Yesterday: `hier`,
		Day:       `il y a %d jour`,
		Days:      `il y a %d jours`,
		Week:      `il y a %d semaine`,
		Weeks:     `il y a %d semaines`,
		Month:     `il y a %d mois`,
		Months:    `il y a %d mois`,
		Year:      `il y a %d an`,
		Years:     `il y a %d ans`,
	},
	"it": {
		NotYet:    `non ancora`,
		JustNow:   `al momento`,
		Minute:    `%d minuto fa`,
		Minutes:   `%d minuti fa`,
		Hour:      `%d ora fa`,
		Hours:     `%d ore fa`,
		Yesterday: `ieri`,
		Day:       `da %d giorno`,
		Days:      `da %d giorni`,
		Week:      `da %d settimana`,
		Weeks:     `da %d settimane`,
		Month:     `da %d mese`,
		Months:    `da %d mesi`,
		Year:      `da %d anno`,
		Years:     `da %d anni`,
	},
	"nl": {
		NotYet:    `nog niet`,
		JustNow:   `dit moment`,
		Minute:    `%d minuut geleden`,
		Minutes:   `%d minuten geleden`,
		Hour:      `%d uur geleden`,
		Hours:     `%d uren geleden`,
		Yesterday: `gisteren`,
		Day:       `%d dag geleden`,
		Days:      `%d dagen geleden`,
		Week:      `%d weke geleden`,
		Weeks:     `%d weken geleden`,
		Month:     `%d maand geleden`,
		Months:    `%d maanden geleden`,
		Year:      `%d jaar geleden.`,
		Years:     `%d jaar geleden.`,
	},
	"pl": {
		NotYet:    `jeszcze nie`,
		JustNow:   `w tej chwili`,
		Minute:    `%d minutę temu`,
		Minutes:   `%d minuty temu`,
		Hour:      `%d godzinę temu`,
		Hours:     `%d godziny temu`,
		Yesterday: `wczoraj`,
		Day:       `%d dzień temu`,
		Days:      `%d dni temu`,
		Week:      `%d tydzień temu`,
		Weeks:     `%d tygodnie temu`,
		Month:     `%d miesiąc temu`,
		Months:    `%d miesiące temu`,
		Year:      `%d rok temu`,
		Years:     `%d lata temu`,
	},
}

// Common errors
var (
	ErrExists     = errors.New("already exists")
	ErrIncomplete = errors.New("missing translation")
	ErrISOCode    = errors.New("invalid language code")
)

// AddTranslation adds the terms for the given language code.
// It fails to do it if the language code already exists or
// if it misses some translation IDs.
func AddTranslation(lang string, tr Terms) error {
	if lang = strings.TrimSpace(lang); lang == "" {
		return ErrISOCode
	}
	if _, ok := i18n[lang]; ok {
		return ErrExists
	}
	for k := range i18n["en"] {
		if _, ok := tr[k]; !ok {
			return ErrIncomplete
		}
	}
	i18n[lang] = tr

	return nil
}

// Time returns in a human readable format the elapsed time
// since the given datetime in english.
// This methods keeps the interface of the first version of the package.
func Time(t time.Time) string {
	return LocalTime(t, "en")
}

// LocalTime returns in a human readable format the elapsed time
// since the given datetime using the given ISO 639-1 language code.
func LocalTime(t time.Time, lang string) string {
	if t.IsZero() || time.Now().Before(t) {
		return tr(NotYet, lang, -1)
	}
	diff := time.Since(t)
	// Duration in seconds
	s := diff.Seconds()
	// Duration in days
	d := int(s / 86400)
	switch {
	case s < 60:
		return tr(JustNow, lang, -1)
	case s < 3600:
		min := int(diff.Minutes())
		return fmt.Sprintf(tr(Minutes, lang, min), min)
	case s < 86400:
		hours := int(diff.Hours())
		return fmt.Sprintf(tr(Hours, lang, hours), hours)
	case d == 1:
		return tr(Yesterday, lang, -1)
	case d < 7:
		return fmt.Sprintf(tr(Days, lang, d), d)
	case d < 31:
		nbWeek := int(math.Ceil(float64(d) / 7))
		if nbWeek < 4 {
			return fmt.Sprintf(tr(Weeks, lang, nbWeek), nbWeek)
		}
		fallthrough
	case d < 365:
		nbMonth := int(math.Ceil(float64(d) / 30))
		if nbMonth < 12 {
			return fmt.Sprintf(tr(Months, lang, nbMonth), nbMonth)
		}
		fallthrough
	default:
		nbYear := int(math.Ceil(float64(d) / 365))
		return fmt.Sprintf(tr(Years, lang, nbYear), nbYear)
	}
}

func tr(id TrID, lang string, nb int) string {
	ltr, ok := i18n[lang]
	if !ok {
		// Uses the english language as fail over.
		ltr = i18n["en"]
	}
	return ltr[changeIfSing(id, nb)]
}

func changeIfSing(id TrID, nb int) TrID {
	switch id {
	case Minutes:
		if nb == 1 {
			return Minute
		}
		return Minutes
	case Hours:
		if nb == 1 {
			return Hour
		}
		return Hours
	case Months:
		if nb == 1 {
			return Month
		}
		return Months
	case Weeks:
		if nb == 1 {
			return Week
		}
		return Weeks
	case Days:
		if nb == 1 {
			return Day
		}
		return Days
	case Years:
		if nb == 1 {
			return Year
		}
		return Years
	default:
		return id
	}
}
