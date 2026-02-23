package utils

import (
	"time"
)

func 	Now() time.Time {
	return time.Now() // pakai local WIB aja
}

func DayRange(dateStr string) (time.Time, time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	var t time.Time
	var err error

	if dateStr == "" {
		t = time.Now().In(loc)
	} else {
		t, err = time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	return start, end, nil
}

