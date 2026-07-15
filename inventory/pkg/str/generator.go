package str

import (
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

func UnixTimestampToUtcTime(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}

func DateTimeStrToRfc3339StringInAsiaJkt(dateStr string) (result string, err error) {
	layoutFormat := "2006-01-02 15:04:05"
	tzJakarta, _ := time.LoadLocation("Asia/Jakarta")
	dateObj, err := time.ParseInLocation(layoutFormat, dateStr, tzJakarta)
	if err != nil {
		return result, err
	}
	dateFormat := dateObj.Format(time.RFC3339)
	return dateFormat, nil
}

// DateStrDdMmYyyyToTime parses DD/MM/YYYY format to time.Time
func DateStrDdMmYyyyToTime(dateStr string) (time.Time, error) {
	layoutFormat := "02/01/2006"
	dateObj, err := time.Parse(layoutFormat, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return dateObj, nil
}
