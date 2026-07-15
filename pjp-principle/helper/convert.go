package helper

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertStringToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	return strconv.Atoi(str)
}

func ParseDateFilter(dateStr, layout string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	endDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return endDate, nil
}

func ParseDate(layout, date string) (time.Time, error) {
	return time.Parse(layout, date)
}

// Format pjp response
func FormatPjpCode(pjpCode int) string {
	return fmt.Sprintf("%04d", pjpCode)
}
