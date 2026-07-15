package live_monitoring

import (
	"scyllax-pjp/model"
	"testing"
)

func TestBuildLiveMonitoringDayRange_UsesJakartaBusinessDate(t *testing.T) {
	startAt, endAt, err := buildLiveMonitoringDayRange("2026-04-07")
	if err != nil {
		t.Fatalf("buildLiveMonitoringDayRange() error = %v", err)
	}

	if startAt != "2026-04-07 00:00:00" {
		t.Fatalf("startAt = %s, want 2026-04-07 00:00:00", startAt)
	}

	if endAt != "2026-04-08 00:00:00" {
		t.Fatalf("endAt = %s, want 2026-04-08 00:00:00", endAt)
	}
}

func TestBuildLiveMonitoringDayRange_ReturnsErrorForInvalidDate(t *testing.T) {
	_, _, err := buildLiveMonitoringDayRange("2026-13-07")
	if err == nil {
		t.Fatal("buildLiveMonitoringDayRange() error = nil, want error")
	}
}

func TestIsValidCurrentCoordinate(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		latitude  float64
		want      bool
	}{
		{name: "valid coordinate", longitude: 106.5979627, latitude: -6.2465789, want: true},
		{name: "zero coordinate", longitude: 0, latitude: 0, want: false},
		{name: "zero latitude", longitude: 106.8, latitude: 0, want: false},
		{name: "zero longitude", longitude: 0, latitude: -6.2, want: false},
		{name: "longitude out of range", longitude: 181, latitude: -6.2, want: false},
		{name: "latitude out of range", longitude: 106.8, latitude: -91, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidCurrentCoordinate(tt.longitude, tt.latitude)
			if got != tt.want {
				t.Fatalf("isValidCurrentCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldReplaceCurrentCoordinate(t *testing.T) {
	checkinAt := int64(1775652000)
	arriveAt := int64(1775652600)
	leaveAt := int64(1775652900)
	checkoutAt := int64(1775653142)
	checkoutAtMillis := int64(1775653142000)

	checkin := model.CurrentCoordinateRow{
		EmpID:        359,
		Longitude:    106.7001,
		Latitude:     -6.2001,
		Timestamp:    &checkinAt,
		Source:       "attendance",
		SourceRank:   2,
		SourceRecord: 10,
	}

	arrive := model.CurrentCoordinateRow{
		EmpID:        359,
		Longitude:    106.7101,
		Latitude:     -6.2101,
		Timestamp:    &arriveAt,
		Source:       "mobile_visit",
		SourceRank:   3,
		SourceRecord: 20,
	}

	leave := model.CurrentCoordinateRow{
		EmpID:        359,
		Longitude:    106.7201,
		Latitude:     -6.2201,
		Timestamp:    &leaveAt,
		Source:       "outlet_visit_list",
		SourceRank:   4,
		SourceRecord: 30,
	}

	checkout := model.CurrentCoordinateRow{
		EmpID:        359,
		Longitude:    106.7301,
		Latitude:     -6.2301,
		Timestamp:    &checkoutAt,
		Source:       "attendance_checkout",
		SourceRank:   1,
		SourceRecord: 40,
	}

	checkoutMillis := model.CurrentCoordinateRow{
		EmpID:        359,
		Longitude:    106.7301,
		Latitude:     -6.2301,
		Timestamp:    &checkoutAtMillis,
		Source:       "attendance_checkout",
		SourceRank:   1,
		SourceRecord: 40,
	}

	invalidCheckout := model.CurrentCoordinateRow{
		EmpID:        359,
		Longitude:    0,
		Latitude:     0,
		Timestamp:    &checkoutAt,
		Source:       "attendance_checkout",
		SourceRank:   1,
		SourceRecord: 41,
	}

	tests := []struct {
		name      string
		current   model.CurrentCoordinateRow
		candidate model.CurrentCoordinateRow
		want      bool
	}{
		{name: "attendance only selected", current: model.CurrentCoordinateRow{}, candidate: checkin, want: true},
		{name: "arrive after attendance selected", current: checkin, candidate: arrive, want: true},
		{name: "leave after arrive selected", current: arrive, candidate: leave, want: true},
		{name: "clock out after leave selected", current: leave, candidate: checkout, want: true},
		{name: "clock out with invalid coordinate ignored", current: leave, candidate: invalidCheckout, want: false},
		{name: "same time prefers better source rank", current: leave, candidate: checkoutMillis, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldReplaceCurrentCoordinate(tt.current, tt.candidate)
			if got != tt.want {
				t.Fatalf("shouldReplaceCurrentCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasCurrentCoordinate(t *testing.T) {
	validAt := int64(1775653142)

	if hasCurrentCoordinate(model.CurrentCoordinateRow{}) {
		t.Fatal("hasCurrentCoordinate() = true, want false for empty row")
	}

	if hasCurrentCoordinate(model.CurrentCoordinateRow{Timestamp: &validAt, Longitude: 0, Latitude: 0}) {
		t.Fatal("hasCurrentCoordinate() = true, want false for zero coordinate")
	}

	if !hasCurrentCoordinate(model.CurrentCoordinateRow{Timestamp: &validAt, Longitude: 106.5979627, Latitude: -6.2465789}) {
		t.Fatal("hasCurrentCoordinate() = false, want true for valid coordinate")
	}
}
