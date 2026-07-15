package helper

import (
	"log"
	"time"
)

// ParseWIBDateOnly untuk nge-parse tanggal format "2006-01-02"
// startOfDay = true -> jam 00:00:00
// startOfDay = false -> jam 23:59:59
func ParseWIBDateOnly(dateStr string, startOfDay bool) *time.Time {
	if dateStr == "" {
		return nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("failed load WIB timezone:", err)
		return nil
	}

	// Parse tanggal tanpa jam
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		log.Println("failed parse date:", err)
		return nil
	}

	if startOfDay {
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	} else {
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, loc)
	}

	return &t
}
