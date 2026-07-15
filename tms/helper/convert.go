package helper

import (
	"errors"
	"time"
)

func ConvertQtyToTotal(qty int) (carton, box, pcs int) {
	carton = qty / 1000
	remaining := qty % 1000
	box = remaining / 100
	pcs = remaining % 100
	return carton, box, pcs
}

func ConvertTotalToQty(carton, box, pcs int) int {
	total := (carton * 1000) + (box * 100) + pcs
	return total
}

func CalculateQty(initial, deduct int64) int64 {
	if deduct != 0 {
		return initial - deduct
	}
	return initial
}

func ParseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"

	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	return parsedDate, nil
}
