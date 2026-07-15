package pjpenhance

import (
	"context"
	"errors"
	"fmt"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"sort"
	"time"
)

func (service *pjpEnhanceService) GetById(ctx context.Context, id int, currentCustomerId, parentCustomerId string) *response.PjpEnhanceResponse {
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)
	pjp := service.pjpRepository.GetPjpById(ctx, tx, id, currentCustomerId)

	routes := service.routeRepository.FindByPjpID(ctx, tx, pjp.ID, currentCustomerId)
	if len(routes) == 0 {
		helper.ErrorPanic(errors.New("PJP was created using the old menu"))
	}

	routeOutlets := service.routeOutletRepository.FindAllOutletsByPjpIds(ctx, tx, []int{pjp.ID}, currentCustomerId)

	routeOutletsHistory := service.routeOutletHistoryRepository.FindByPjpId(ctx, tx, pjp.ID, currentCustomerId)
	routePopPermanents := service.routePopRepository.FindByPjpID(ctx, tx, pjp.ID, currentCustomerId, parentCustomerId)

	pjpResponse := mapToPjpResponse(pjp)
	routesResponse := mapToRoutesResponse(routes, routeOutlets)
	visitDay := mapToVisitDays(routeOutletsHistory, routePopPermanents)

	if visitDay == nil {
		visitDay = []response.VisitDay{}
	}

	if routesResponse == nil {
		routesResponse = []response.Routes{}
	}

	return &response.PjpEnhanceResponse{
		Pjp:      pjpResponse,
		Routes:   routesResponse,
		VisitDay: visitDay, // jika tetap kosong, akan di-encode jadi [] bukan null
	}

}

func mapToPjpResponse(pjp model.Pjp) response.Pjp {
	var pjpResponse response.Pjp
	pjpResponse.PjpCode = helper.FormatPjpCode(pjp.PjpCode)
	helper.Automapper(&pjp, &pjpResponse)
	return pjpResponse
}

func mapToRoutesResponse(routes []model.Route, routeOutlets []model.RouteOutlet) []response.Routes {
	outletsByRoute := make(map[int][]model.RouteOutlet)
	for _, outlet := range routeOutlets {
		outletsByRoute[outlet.RouteCode] = append(outletsByRoute[outlet.RouteCode], outlet)
	}

	var routesResponse []response.Routes
	for _, route := range routes {
		var routeResponse response.Routes
		helper.Automapper(&route, &routeResponse)

		routeResponse.Outlets = []response.Outlets{}

		if outlets, exists := outletsByRoute[route.RouteCode]; exists {
			var outletsResponse []response.Outlets
			helper.Automapper(&outlets, &outletsResponse)
			for i := range outletsResponse {
				outletsResponse[i].Type = "outlet"
				outletsResponse[i].RouteName = route.RouteName
			}
			routeResponse.Outlets = outletsResponse
		}

		routesResponse = append(routesResponse, routeResponse)
	}
	return routesResponse
}

func mapToVisitDays(routeOutletsHistory []model.RouteOutletHistory, routePopPemanents []model.RoutePopPermanent) []response.VisitDay {
	var visitDays []response.VisitDay
	usedHistory := make(map[string]bool)
	for _, perm := range routePopPemanents {
		if perm.RouteCode == nil {
			continue
		}
		var histories []response.OutletHistory
		for _, history := range routeOutletsHistory {
			if historyKey(history) == permanentKey(perm) {
				var h response.OutletHistory
				helper.Automapper(&history, &h)
				histories = append(histories, h)
			}
		}
		if len(histories) == 0 {
			continue
		}
		usedHistory[permanentKey(perm)] = true
		visitDays = append(visitDays, response.VisitDay{
			ID:                   perm.ID,
			Day:                  perm.Day,
			IndexDay:             histories[0].IndexDay,
			Week:                 perm.Week,
			WorkingDayCalendarID: perm.WorkingDayCalendarID,
			StartWeek:            formatDatePointer(histories[0].StartWeek),
			Year:                 perm.Year,
			Date:                 perm.Date.Format("2006-01-02"),
			IsInCurrentYear:      histories[0].IsInCurrentYear,
			Visit: response.RoutesHistory{
				RouteCode: histories[0].RouteCode,
				RouteName: histories[0].RouteName,
				CustID:    perm.CustID,
				Outlets:   histories,
			},
		})
	}

	historiesByDay := make(map[string][]model.RouteOutletHistory)
	for _, history := range routeOutletsHistory {
		key := historyKey(history)
		if usedHistory[key] {
			continue
		}
		historiesByDay[key] = append(historiesByDay[key], history)
	}

	for _, histories := range historiesByDay {
		visitDays = append(visitDays, visitDayFromHistory(histories))
	}

	sort.Slice(visitDays, func(i, j int) bool {
		if visitDays[i].Date == visitDays[j].Date {
			return visitDays[i].Visit.RouteCode < visitDays[j].Visit.RouteCode
		}
		return visitDays[i].Date < visitDays[j].Date
	})

	return visitDays
}

func visitDayFromHistory(histories []model.RouteOutletHistory) response.VisitDay {
	history := histories[0]
	var outletHistories []response.OutletHistory
	helper.Automapper(&histories, &outletHistories)

	day := history.Date.Format("Mon")
	if history.Date.IsZero() {
		day = ""
	}

	return response.VisitDay{
		Day:             day,
		IndexDay:        history.IndexDay,
		Week:            history.Week,
		StartWeek:       formatDatePointer(history.StartWeek),
		Year:            history.Year,
		Date:            history.Date.Format("2006-01-02"),
		IsInCurrentYear: history.IsInCurrentYear,
		Visit: response.RoutesHistory{
			RouteCode: history.RouteCode,
			RouteName: history.RouteName,
			CustID:    history.CustID,
			Outlets:   outletHistories,
		},
	}
}

func historyKey(history model.RouteOutletHistory) string {
	return visitDayKey(history.RouteCode, history.Week, history.Year, history.Date)
}

func permanentKey(permanent model.RoutePopPermanent) string {
	return visitDayKey(*permanent.RouteCode, permanent.Week, permanent.Year, permanent.Date)
}

func visitDayKey(routeCode int, week int, year int, date time.Time) string {
	return fmt.Sprintf("%s|%d|%d|%d", date.Format("2006-01-02"), routeCode, week, year)
}

func formatDatePointer(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format("2006-01-02")
}
