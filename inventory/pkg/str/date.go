package str

import "time"

func GetJakartaDate() time.Time {
	tzJakarta, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().Local()

	return time.Date(now.In(tzJakarta).Year(), now.In(tzJakarta).Month(), now.In(tzJakarta).Day(), 0, 0, 0, 0, tzJakarta)
}
