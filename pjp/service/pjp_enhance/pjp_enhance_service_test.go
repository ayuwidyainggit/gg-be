package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/model"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestResolveRouteCode(t *testing.T) {
	existingRoutes := []model.Route{
		{ID: 101, RouteCode: 4321, RouteName: "Route 1", Sequence: 1},
		{ID: 102, RouteCode: 5432, RouteName: "Route 2", Sequence: 2},
	}

	t.Run("uses request route code when provided", func(t *testing.T) {
		routeCode := 9999
		got := resolveRouteCode(request.RoutesCreatePjp{RouteCode: &routeCode}, existingRoutes, 1)
		if got != routeCode {
			t.Fatalf("got %d, want %d", got, routeCode)
		}
	})

	t.Run("uses route_id alias to preserve existing route code", func(t *testing.T) {
		routeID := 102
		got := resolveRouteCode(request.RoutesCreatePjp{RouteID: &routeID}, existingRoutes, 7)
		if got != 5432 {
			t.Fatalf("got %d, want %d", got, 5432)
		}
	})

	t.Run("falls back to sequence when request has no route identifiers", func(t *testing.T) {
		got := resolveRouteCode(request.RoutesCreatePjp{RouteName: "Renamed Route 2"}, existingRoutes, 2)
		if got != 5432 {
			t.Fatalf("got %d, want %d", got, 5432)
		}
	})
}

func TestMapToVisitDaysIncludesWorkingDayCalendarID(t *testing.T) {
	routeCode := 4321
	startWeek := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	workingDayCalendarID := int64(77)

	histories := []model.RouteOutletHistory{
		{
			RouteCode: routeCode,
			RouteName: "Route 1",
			Week:      2,
			Year:      2026,
			Date:      startWeek,
			IndexDay:  1,
			StartWeek: &startWeek,
		},
	}

	permanents := []model.RoutePopPermanent{
		{
			ID:                   1,
			Day:                  "Monday",
			Week:                 2,
			Year:                 2026,
			Date:                 startWeek,
			RouteCode:            &routeCode,
			CustID:               "C26002",
			WorkingDayCalendarID: &workingDayCalendarID,
		},
	}

	visitDays := mapToVisitDays(histories, permanents)
	if len(visitDays) != 1 {
		t.Fatalf("got %d visit days, want 1", len(visitDays))
	}
	if visitDays[0].WorkingDayCalendarID == nil || *visitDays[0].WorkingDayCalendarID != workingDayCalendarID {
		t.Fatalf("got working_day_calendar_id %v, want %d", visitDays[0].WorkingDayCalendarID, workingDayCalendarID)
	}
}

func TestMapToVisitDaysIncludesHistoryWithoutPermanentRow(t *testing.T) {
	startWeek := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)

	visitDays := mapToVisitDays([]model.RouteOutletHistory{
		{
			RouteCode:       4654,
			RouteName:       "Route 1",
			OutletID:        1756,
			OutletCode:      "SJ016",
			OutletName:      "Dicky TK",
			Week:            16,
			Year:            2026,
			Date:            startWeek,
			IndexDay:        1,
			StartWeek:       &startWeek,
			IsInCurrentYear: true,
			CustID:          "C260040002",
		},
	}, nil)

	if len(visitDays) != 1 {
		t.Fatalf("got %d visit days, want 1", len(visitDays))
	}
	if visitDays[0].Date != "2026-04-13" {
		t.Fatalf("got date %s, want 2026-04-13", visitDays[0].Date)
	}
	if len(visitDays[0].Visit.Outlets) != 1 {
		t.Fatalf("got %d outlets, want 1", len(visitDays[0].Visit.Outlets))
	}
	if visitDays[0].Visit.Outlets[0].OutletName != "Dicky TK" {
		t.Fatalf("got outlet %s, want Dicky TK", visitDays[0].Visit.Outlets[0].OutletName)
	}
}

type routeOutletHistoryRepoStub struct{}

func (routeOutletHistoryRepoStub) CreateBulk(context.Context, *gorm.DB, []model.RouteOutletHistory) {}
func (routeOutletHistoryRepoStub) FindByPjpId(context.Context, *gorm.DB, int, string) []model.RouteOutletHistory {
	return nil
}
func (routeOutletHistoryRepoStub) FindByPjpIdToday(context.Context, *gorm.DB, []int, string) []model.RouteOutletHistory {
	return nil
}
func (routeOutletHistoryRepoStub) DeleteByPjpId(context.Context, *gorm.DB, int, string) {}
func (routeOutletHistoryRepoStub) DeleteByVisitDay(context.Context, *gorm.DB, model.RouteOutletHistory) {
}

type routePopPermanentRepoStub struct{}

func (routePopPermanentRepoStub) CreateBulk(context.Context, *gorm.DB, []model.RoutePopPermanent) {}
func (routePopPermanentRepoStub) FindByPjpID(context.Context, *gorm.DB, int, string, string) []model.RoutePopPermanent {
	return nil
}
func (routePopPermanentRepoStub) DeleteByVisitDay(context.Context, *gorm.DB, model.RoutePopPermanent) {
}

func TestCreateVisitHistoryIncludesWorkingDayCalendarID(t *testing.T) {
	workingDayCalendarID := int64(77)
	service := &pjpEnhanceService{
		routeOutletHistoryRepository: routeOutletHistoryRepoStub{},
	}

	savedPjp := model.Pjp{ID: 11, PjpCode: 12}
	savedRoutes := []model.Route{{RouteCode: 4321, RouteName: "Route 1"}}
	visitDays := []request.VisitDayCreatePjp{
		{
			ID:                   1,
			Day:                  "Monday",
			IndexDay:             1,
			Week:                 2,
			WorkingDayCalendarID: &workingDayCalendarID,
			StartWeek:            "2026-01-05",
			Year:                 2026,
			Date:                 "2026-01-05",
			IsInCurrentYear:      true,
			Visit: request.RoutesCreatePjp{
				RouteName: "Route 1",
				Destination: []request.Destination{
					{ID: 1, Code: "OUT1", Name: "Outlet 1", Type: "outlet"},
				},
			},
		},
	}

	routePopPermanents := service.createVisitHistory(context.Background(), nil, savedPjp, savedRoutes, visitDays, "C26002")
	if len(routePopPermanents) != 1 {
		t.Fatalf("got %d route pop permanents, want 1", len(routePopPermanents))
	}
	if routePopPermanents[0].WorkingDayCalendarID == nil || *routePopPermanents[0].WorkingDayCalendarID != workingDayCalendarID {
		t.Fatalf("got working_day_calendar_id %v, want %d", routePopPermanents[0].WorkingDayCalendarID, workingDayCalendarID)
	}
}

func TestFilterEditableVisitDaysSkipsPastWeek(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, loc)

	got := filterEditableVisitDays([]request.VisitDayCreatePjp{
		{Date: "2026-06-28", Week: 26},
		{Date: "2026-06-29", Week: 27},
	}, now)

	if len(got) != 1 {
		t.Fatalf("got %d editable visit days, want 1", len(got))
	}
	if got[0].Date != "2026-06-29" {
		t.Fatalf("got editable date %s, want 2026-06-29", got[0].Date)
	}
}
