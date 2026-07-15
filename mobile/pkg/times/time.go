package times

import (
	"fmt"
	"os"
	"time"
)

const DateTimeLayout = "2006-01-02 15:04:05.0"

// FormatDateWithZeroTime validates the input date and formats it
// Input format example: "2006-01-02"
func FormatDateWithZeroTime(dateStr string) (string, error) {
	// Parse input date (validation happens here)
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date: %w", err)
	}

	// Force time to midnight
	result := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	return result.Format(DateTimeLayout), nil
}

func GetCurrentTime() (time.Time, error) {
	locationString := os.Getenv("TIMEZONE")
	if locationString == "" {
		locationString = "Asia/Jakarta"
	}
	loc, err := time.LoadLocation(locationString)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now().In(loc)
	return now, nil
}
