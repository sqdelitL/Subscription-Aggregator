package util

import (
	"fmt"
	"time"
)

func ParseMMYYYY(date string) (time.Time, error) {
	t, err := time.Parse("01-2006", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format %q, expected MM-YYYY: %w", date, err)
	}
	return t, nil
}

func FormatMMYYYY(t time.Time) string {
	return t.Format("01-2006")
}
