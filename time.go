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
	// LastMinute is the translation ID for the "1 minute ago" text.
	LastMinute
	// Minutes is the translation ID for the "%d minutes ago" text.
	Minutes
	// LastHour is the translation ID for the "1 hour ago" text.
	LastHour
	// Hours is the translation ID for the "%d hours ago" text.
	Hours
	// Yesterday is the translation ID for the "yesterday" text.
	Yesterday
	// Days is the translation ID for the "%d days ago" text.
	Days
	// Weeks is the translation ID for the "`%d weeks ago" text.
	Weeks
	// Months is the translation ID for the "%d months ago" text.
	Months
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
		NotYet:     `noch nicht`,
		JustNow:    `im Moment`,
		LastMinute: `vor 1 Minute`,
		Minutes:    `vor %d Minuten`,
		LastHour:   `vor 1 Stunde`,
		Hours:      `vor %d Stunden`,
		Yesterday:  `gestern`,
		Days:       `vor %d Tagen`,
		Weeks:      `vor %d Wochen`,
		Months:     `vor %d Monaten`,
		Years:      `vor %d Jahren`,
	},
	"en": {
		NotYet:     `not yet`,
		JustNow:    `just now`,
		LastMinute: `1 minute ago`,
		Minutes:    `%d minutes ago`,
		LastHour:   `1 hour ago`,
		Hours:      `%d hours ago`,
		Yesterday:  `yesterday`,
		Days:       `%d days ago`,
		Weeks:      `%d weeks ago`,
		Months:     `%d months ago`,
		Years:      `%d years ago`,
	},
	"es": {
		NotYet:     `aún no`,
		JustNow:    `al instante`,
		LastMinute: `hace 1 minuto`,
		Minutes:    `hace %d minutos`,
		LastHour:   `hace 1 hora`,
		Hours:      `hace %d horas`,
		Yesterday:  `ayer`,
		Days:       `hace %d días`,
		Weeks:      `hace %d semanas`,
		Months:     `hace %d meses`,
		Years:      `hace %d años`,
	},
	"fr": {
		NotYet:     `pas encore`,
		JustNow:    `à l'instant`,
		LastMinute: `il y a 1 minute`,
		Minutes:    `il y a %d minutes`,
		LastHour:   `il y a 1 heure`,
		Hours:      `il y a %d heures`,
		Yesterday:  `hier`,
		Days:       `il y a %d jours`,
		Weeks:      `il y a %d semaines`,
		Months:     `il y a %d mois`,
		Years:      `il y a %d ans`,
	},
	"it": {
		NotYet:     `non ancora`,
		JustNow:    `al momento`,
		LastMinute: `1 minuto fa`,
		Minutes:    `%d minuti fa`,
		LastHour:   `1 ora fa`,
		Hours:      `%d ore fa`,
		Yesterday:  `ieri`,
		Days:       `da %d giorni`,
		Weeks:      `da %d settimane`,
		Months:     `da %d mesi`,
		Years:      `da %d anni`,
	},
	"nl": {
		NotYet:     `nog niet`,
		JustNow:    `dit moment`,
		LastMinute: `1 minuut geleden`,
		Minutes:    `%d minuten geleden`,
		LastHour:   `1 uur geleden`,
		Hours:      `%d uren geleden`,
		Yesterday:  `gisteren`,
		Days:       `%d dagen geleden`,
		Weeks:      `%d weken geleden`,
		Months:     `%d maanden geleden`,
		Years:      `%d jaar geleden.`,
	},
	"pl": {
		NotYet:     `jeszcze nie`,
		JustNow:    `w tej chwili`,
		LastMinute: `1 minutę temu`,
		Minutes:    `%d minuty temu`,
		LastHour:   `1 godzinę temu`,
		Hours:      `%d godziny temu`,
		Yesterday:  `wczoraj`,
		Days:       `%d dni temu`,
		Weeks:      `%d tygodnie temu`,
		Months:     `%d miesiące temu`,
		Years:      `%d lata temu`,
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
		return tr(NotYet, lang)
	}
	diff := time.Since(t)
	// Duration in seconds
	s := diff.Seconds()
	// Duration in days
	d := int(s / 86400)
	switch {
	case s < 60:
		return tr(JustNow, lang)
	case s < 120:
		return tr(LastMinute, lang)
	case s < 3600:
		return fmt.Sprintf(tr(Minutes, lang), int(diff.Minutes()))
	case s < 7200:
		return tr(LastHour, lang)
	case s < 86400:
		return fmt.Sprintf(tr(Hours, lang), int(diff.Hours()))
	case d == 1:
		return tr(Yesterday, lang)
	case d < 7:
		return fmt.Sprintf(tr(Days, lang), d)
	case d < 31:
		return fmt.Sprintf(tr(Weeks, lang), int(math.Ceil(float64(d)/7)))
	case d < 365:
		return fmt.Sprintf(tr(Months, lang), int(math.Ceil(float64(d)/30)))
	default:
		return fmt.Sprintf(tr(Years, lang), int(math.Ceil(float64(d)/365)))
	}
}

func tr(id TrID, lang string) string {
	ltr, ok := i18n[lang]
	if !ok {
		// Uses the english language as fail over.
		ltr = i18n["en"]
	}
	return ltr[id]
}
