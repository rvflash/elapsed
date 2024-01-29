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
		Days:      `hace %d días`,
		Week:      `hace %d semana`,
		Weeks:     `hace %d semanas`,
		Month:     `hace %d mes`,
		Months:    `hace %d meses`,
		Year:      `hace %d año`,
		Years:     `hace %d años`,
	},
	"pt": {
		NotYet:    `ainda não`,
		JustNow:   `agora mesmo`,
		Minute:    `há %d minuto`,
		Minutes:   `há %d minutos`,
		Hour:      `há %d hora`,
		Hours:     `há %d horas`,
		Yesterday: `ontem`,
		Days:      `há %d dias`,
		Week:      `há %d semana`,
		Weeks:     `há %d semanas`,
		Month:     `há %d mês`,
		Months:    `há %d meses`,
		Year:      `há %d ano`,
		Years:     `há %d anos`,
	},
	"ru": {
		NotYet:    `еще нет`,
		JustNow:   `сейчас`,
		Minute:    `%d минуту назад`,
		Minutes:   `%d минут назад`,
		Hour:      `%d час назад`,
		Hours:     `%d часов назад`,
		Yesterday: `вчера`,
		Days:      `%d дней назад`,
		Week:      `%d неделю назад`,
		Weeks:     `%d недели назад`,
		Month:     `%d месяц назад`,
		Months:    `%d месяца назад`,
		Year:      `%d год назад`,
		Years:     `%d года назад`,
	},
	"fr": {
		NotYet:    `pas encore`,
		JustNow:   `à l'instant`,
		Minute:    `il y a %d minute`,
		Minutes:   `il y a %d minutes`,
		Hour:      `il y a %d heure`,
		Hours:     `il y a %d heures`,
		Yesterday: `hier`,
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
		Days:      `%d dni temu`,
		Week:      `%d tydzień temu`,
		Weeks:     `%d tygodnie temu`,
		Month:     `%d miesiąc temu`,
		Months:    `%d miesiące temu`,
		Year:      `%d rok temu`,
		Years:     `%d lata temu`,
	},
	"zh": {
		NotYet:    `未到`,
		JustNow:   `刚刚`,
		Minute:    `%d 分钟前`,
		Minutes:   `%d 分钟前`,
		Hour:      `%d 小时前`,
		Hours:     `%d 小时前`,
		Yesterday: `昨天`,
		Days:      `%d 天前`,
		Week:      `%d 周前`,
		Weeks:     `%d 周前`,
		Month:     `%d 个月前`,
		Months:    `%d 个月前`,
		Year:      `%d 年前`,
		Years:     `%d 年前`,
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
	case s < 3600:
		min := int(diff.Minutes())
		return fmt.Sprintf(tr(changeIfSing(Minutes, min), lang), min)
	case s < 86400:
		hours := int(diff.Hours())
		return fmt.Sprintf(tr(changeIfSing(Hours, hours), lang), hours)
	case d == 1:
		return tr(Yesterday, lang)
	case d < 7:
		return fmt.Sprintf(tr(changeIfSing(Days, d), lang), d)
	case d < 31:
		nbWeek := int(math.Ceil(float64(d) / 7))
		if nbWeek < 4 {
			return fmt.Sprintf(tr(changeIfSing(Weeks, nbWeek), lang), nbWeek)
		}
		fallthrough
	case d < 365:
		nbMonth := int(math.Ceil(float64(d) / 30))
		if nbMonth < 12 {
			return fmt.Sprintf(tr(changeIfSing(Months, nbMonth), lang), nbMonth)
		}
		fallthrough
	default:
		nbYear := int(math.Ceil(float64(d) / 365))
		return fmt.Sprintf(tr(changeIfSing(Years, nbYear), lang), nbYear)
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

func changeIfSing(id TrID, nb int) TrID {
	if nb != 1 {
		return id
	}
	switch id {
	case Minutes:
		return Minute
	case Hours:
		return Hour
	case Months:
		return Month
	case Weeks:
		return Week
	case Years:
		return Year
	default:
		return id
	}
}
