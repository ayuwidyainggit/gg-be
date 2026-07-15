package route

import (
	"scyllax-pjp/model"
	"testing"
	"time"
)

func TestShouldCreateDailyRouteSkipsPastWeeks(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, loc)

	if shouldCreateDailyRoute(model.RoutePopPermanent{Date: time.Date(2026, 6, 28, 0, 0, 0, 0, loc)}, now) {
		t.Fatal("expected past week route to be skipped")
	}

	if !shouldCreateDailyRoute(model.RoutePopPermanent{Date: time.Date(2026, 6, 29, 0, 0, 0, 0, loc)}, now) {
		t.Fatal("expected current week route to be created")
	}
}
