package live_monitoring

import (
	"math"
	"scyllax-pjp/model"
)

func isValidCurrentCoordinate(longitude, latitude float64) bool {
	if math.IsNaN(longitude) || math.IsNaN(latitude) {
		return false
	}

	if math.IsInf(longitude, 0) || math.IsInf(latitude, 0) {
		return false
	}

	if longitude == 0 || latitude == 0 {
		return false
	}

	if longitude < -180 || longitude > 180 {
		return false
	}

	if latitude < -90 || latitude > 90 {
		return false
	}

	return true
}

func normalizeCurrentCoordinateEpoch(timestamp *int64) int64 {
	if timestamp == nil {
		return 0
	}

	value := *timestamp
	if value > 9999999999 || value < -9999999999 {
		return value / 1000
	}

	return value
}

func hasCurrentCoordinate(candidate model.CurrentCoordinateRow) bool {
	return candidate.Timestamp != nil && isValidCurrentCoordinate(candidate.Longitude, candidate.Latitude)
}

func shouldReplaceCurrentCoordinate(current, candidate model.CurrentCoordinateRow) bool {
	if !hasCurrentCoordinate(candidate) {
		return false
	}

	if !hasCurrentCoordinate(current) {
		return true
	}

	currentTimestamp := normalizeCurrentCoordinateEpoch(current.Timestamp)
	candidateTimestamp := normalizeCurrentCoordinateEpoch(candidate.Timestamp)

	if candidateTimestamp != currentTimestamp {
		return candidateTimestamp > currentTimestamp
	}

	if candidate.SourceRank != current.SourceRank {
		return candidate.SourceRank < current.SourceRank
	}

	return candidate.SourceRecord > current.SourceRecord
}
