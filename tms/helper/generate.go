package helper

import (
	"fmt"
	"math/rand"
	"time"
)

// func MappingAutoShipmentNo(vehicleID string, lastNumber string) string {
// 	currentYear := time.Now().Year()
// 	year := fmt.Sprintf("%02d", currentYear%100)

// 	rand.New(rand.NewSource(time.Now().UnixNano()))
// 	autoGenerate := fmt.Sprintf("%05d", rand.Intn(10000))

// 	return "1" + vehicleID + year + autoGenerate + lastNumber
// }

// func ManualGenerateShipmentNo(vehicleID int) string {
// 	currentYear := time.Now().Year()
// 	year := fmt.Sprintf("%02d", currentYear%100)

// 	rand.New(rand.NewSource(time.Now().UnixNano()))
// 	autoGenerate := fmt.Sprintf("%06d", rand.Intn(10000))

// 	return "2" + fmt.Sprintf("%d", vehicleID) + year + autoGenerate
// }

func MappingAutoShipmentNo(vehicleID string, lastNumber string) string {
	currentYear := time.Now().Year()
	year := fmt.Sprintf("%02d", currentYear%100)
	month := fmt.Sprintf("%02d", int(time.Now().Month()))
	day := fmt.Sprintf("%02d", time.Now().Day())

	rand.Seed(time.Now().UnixNano())
	autoGenerate := fmt.Sprintf("%03d", rand.Intn(1000))
	seq4digit := "1" + autoGenerate

	return "DO" + year + month + day + seq4digit + lastNumber
}

func ManualGenerateShipmentNo(vehicleID int) string {
	currentYear := time.Now().Year()
	year := fmt.Sprintf("%02d", currentYear%100)
	month := fmt.Sprintf("%02d", int(time.Now().Month()))
	day := fmt.Sprintf("%02d", time.Now().Day())

	rand.Seed(time.Now().UnixNano())
	autoGenerate := fmt.Sprintf("%03d", rand.Intn(1000))
	seq4digit := "2" + autoGenerate

	return "DO" + year + month + day + seq4digit
}

func GenerateShipmentNo() string {
	currentYear := time.Now().Year()
	year := fmt.Sprintf("%02d", currentYear)
	month := fmt.Sprintf("%02d", int(time.Now().Month()))
	day := fmt.Sprintf("%02d", time.Now().Day())

	rand.Seed(time.Now().UnixNano())
	autoGenerate := fmt.Sprintf("%04d", rand.Intn(1000))

	return "DO" + year + month + day + autoGenerate
}
