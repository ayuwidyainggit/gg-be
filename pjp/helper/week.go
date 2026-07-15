package helper

import "time"

func StartOfWeek(t time.Time) time.Time {
	year, month, day := t.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	offset := (int(start.Weekday()) - int(time.Monday) + 7) % 7
	return start.AddDate(0, 0, -offset)
}

func IsBeforeCurrentWeek(date time.Time, now time.Time) bool {
	year, month, day := date.In(now.Location()).Date()
	dateOnly := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	return dateOnly.Before(StartOfWeek(now))
}
