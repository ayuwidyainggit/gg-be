package entity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type FlexibleTime struct {
	time.Time
}

func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}

	formats := []string{
		time.RFC3339,                  // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,              // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",        // "2006-01-02T15:04:05Z"
		"2006-01-02T15:04:05.000Z",    // "2006-01-02T15:04:05.000Z"
		"2006-01-02T15:04:05.000000Z", // "2006-01-02T15:04:05.000000Z"
		"2006-01-02",                  // "2006-01-02" (date only)
		"2006-01-02 15:04:05",         // "2006-01-02 15:04:05"
		"2006-01-02 15:04:05.000",     // "2006-01-02 15:04:05.000"
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			ft.Time = t
			return nil
		}
	}

	return &time.ParseError{
		Layout:     "multiple formats",
		Value:      s,
		LayoutElem: "2006-01-02 or RFC3339",
		ValueElem:  s,
		Message:    "unable to parse date",
	}
}

func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time.Format(time.RFC3339))
}

type FlexibleRouteId string

func (frid *FlexibleRouteId) UnmarshalJSON(b []byte) error {
	var num int64
	if err := json.Unmarshal(b, &num); err == nil {
		*frid = FlexibleRouteId(fmt.Sprintf("%d", num))
		return nil
	}

	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*frid = FlexibleRouteId(str)
	return nil
}

func (frid FlexibleRouteId) String() string {
	return string(frid)
}

func (frid FlexibleRouteId) Int64() (int64, error) {
	return strconv.ParseInt(string(frid), 10, 64)
}
