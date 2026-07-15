package helper

import "time"

func GetStartOfISOWeek(year int, week int) time.Time {
	// ISO 8601: week 1 is the week with the first Thursday of the year
	// Gunakan time.Date dengan hari Kamis minggu pertama, lalu geser ke Senin
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)
	isoYearStart := jan4.AddDate(0, 0, -int(jan4.Weekday()-time.Monday))

	// Geser ke minggu ke-n
	startOfWeek := isoYearStart.AddDate(0, 0, (week-1)*7)
	return startOfWeek
}
