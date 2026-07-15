package str

import (
	"fmt"
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

func UnixTimestampToUtcTime(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}

func DateStrToRfc3339String(dateStr string) (result string, err error) {
	layoutFormat := "2006-01-02"

	dateObj, err := time.Parse(layoutFormat, dateStr)
	if err != nil {
		return result, err
	}

	return dateObj.Format(time.RFC3339), nil
}

func ArrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func DateTimeToISO9075String(dateTime time.Time) (result string, err error) {
	layoutFormat := "2006-01-02 15:04:05"
	return dateTime.Format(layoutFormat), nil
}
