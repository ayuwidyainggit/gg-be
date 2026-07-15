package str

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

func PhoneConvertToAbbv(phone string) string {
	isMatch, _ := regexp.MatchString(`^[0{1}]`, phone)
	if isMatch {
		re := regexp.MustCompile(`^[0{1}]`)
		s := re.ReplaceAllString(phone, `+62`)

		return s
	}

	return phone
}

func Replacer(source string, replacer *strings.Replacer) string {
	return replacer.Replace(source)
}

func DateStrToRfc3339String(dateStr string) (result string, err error) {
	layoutFormat := "2006-01-02"

	dateObj, err := time.Parse(layoutFormat, dateStr)
	if err != nil {
		return result, err
	}

	return dateObj.Format(time.RFC3339), nil
}

func ConvertStringTimeToTimeObject(timeStr string) (*time.Time, error) {
	// Format sesuai dengan layout time Go
	layout := time.RFC3339
	// Parse string ke time.Time
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return nil, errors.New("Error parsing date")
	}

	return &parsedTime, nil
}

func UnixTimestampToUtcTime(unix int64) time.Time {
	if unix > 1e11 {
		unix = unix / 1000
	}
	return time.Unix(unix, 0).UTC()
}

func UnixTimestampToUtcDate(unix int64) time.Time {
	result := UnixTimestampToUtcTime(unix)
	return time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, 0, time.UTC)
}

func UnixTimestampToAsiaJakartaTime(unix int64) time.Time {
	if unix > 1e11 {
		unix = unix / 1000
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Unix(unix, 0).In(loc)
}

func UnixTimestampToAsiaJakartaTimePointer(unix int64) string {
	if unix <= 0 {
		return ""
	}

	// kalau ternyata unix dalam millisecond
	if unix > 1e11 {
		unix = unix / 1000
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")

	result := time.Unix(unix, 0).In(loc)

	// format sesuai kebutuhan, misal "2006-01-02 15:04:05"
	return result.Format("2006-01-02 15:04:05")
}

func ConvertStringDateToTimeObject(dateStr string) (time.Time, error) {
	layout := "2006-01-02"

	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %w", err)
	}

	return parsedTime, nil
}

func FormatTimeToDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

func DiffDuration(ts1, ts2 int64) string {
	if ts1 == 0 || ts2 == 0 {
		return ""
	}

	t1 := time.UnixMilli(ts1)
	t2 := time.UnixMilli(ts2)
	diff := t2.Sub(t1)

	// hitung menit bulat
	minutes := int(math.Round(diff.Minutes()))

	return fmt.Sprintf("%d", minutes)
}

// ExtractDateParts mengembalikan day, month, dan year dari string tanggal "YYYY-MM-DD"
func ExtractDateParts(dateStr string) (day, month, year int, err error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("format tanggal tidak valid: %w", err)
	}

	year = t.Year()
	month = int(t.Month())
	day = t.Day()

	return day, month, year, nil
}
