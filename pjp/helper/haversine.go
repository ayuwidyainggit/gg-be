package helper

import (
	"math"
)

// earthRadiusMeters is the mean radius of the Earth in meters.
const earthRadiusMeters = 6371000

// CalculateHaversineDistance calculates the great-circle distance between two GPS coordinates
// using the Haversine formula. Returns distance in meters.
//
// Parameters:
//   - lat1, lon1: coordinates of the first point (e.g., outlet position)
//   - lat2, lon2: coordinates of the second point (e.g., arrival/current position)
//
// Returns:
//   - distance in meters as integer (rounded to nearest meter)
func CalculateHaversineDistance(lat1, lon1, lat2, lon2 float64) int {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(math.Round(earthRadiusMeters * c))
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
