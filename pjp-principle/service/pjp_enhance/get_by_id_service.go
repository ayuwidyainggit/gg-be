package pjpenhance

import (
	"context"
	"errors"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"time"
)

func (service *pjpEnhanceService) GetById(ctx context.Context, id int, currentCustomerId string) *response.PjpEnhanceResponse {
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

	routeOutlets := service.routeOutletRepository.FindAllOutletsByPjpId(ctx, tx, pjp.ID, currentCustomerId)

	routeOutletsHistory := service.routeOutletHistoryRepository.FindByPjpId(ctx, tx, pjp.ID, currentCustomerId)
	routePopPermanents := service.routePopRepository.FindByPjpID(ctx, tx, pjp.ID, currentCustomerId)

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

func mapToRoutesResponse(routes []model.Route, routeOutlets []model.Destination) []response.Routes {
	outletsByRoute := make(map[int][]model.Destination)
	for _, outlet := range routeOutlets {
		outletsByRoute[outlet.RouteCode] = append(outletsByRoute[outlet.RouteCode], outlet)
	}

	var routesResponse []response.Routes
	for _, route := range routes {
		var routeResponse response.Routes
		helper.Automapper(&route, &routeResponse)

		if outlets, exists := outletsByRoute[route.RouteCode]; exists {
			for _, outlet := range outlets {
				routeName := route.RouteName
				if outlet.DestinationType == "distributor" {
					routeResponse.Outlets = append(routeResponse.Outlets, response.Distributors{
						ID:                 outlet.ID,
						RouteCode:          outlet.RouteCode,
						RouteName:          routeName,
						DestinationID:      outlet.DestinationID,
						DestinationCode:    outlet.DestinationCode,
						DestinationName:    outlet.DestinationName,
						DestinationStatus:  outlet.DestinationStatus,
						DestinationAddress: outlet.DestinationAddress,
						Longitude:          outlet.Longitude,
						Latitude:           outlet.Latitude,
						AvgSalesWeek:       outlet.AvgSalesWeek,
						PjpID:              outlet.PjpID,
						PjpCode:            outlet.PjpCode,
						Status:             outlet.Status,
						CustID:             outlet.CustID,
						Type:               "distributor",
					})
				} else {
					routeResponse.Outlets = append(routeResponse.Outlets, response.Outlets{
						ID:                 outlet.ID,
						RouteCode:          outlet.RouteCode,
						RouteName:          routeName,
						DestinationID:      outlet.DestinationID,
						DestinationCode:    outlet.DestinationCode,
						DestinationName:    outlet.DestinationName,
						DestinationStatus:  outlet.DestinationStatus,
						DestinationAddress: outlet.DestinationAddress,
						Longitude:          outlet.Longitude,
						Latitude:           outlet.Latitude,
						AvgSalesWeek:       outlet.AvgSalesWeek,
						PjpID:              outlet.PjpID,
						PjpCode:            outlet.PjpCode,
						Status:             outlet.Status,
						CustID:             outlet.CustID,
						Type:               "outlet",
					})
				}
			}
		}
		routesResponse = append(routesResponse, routeResponse)
	}

	return routesResponse
}

func mapToVisitDays(routeOutletsHistory []model.DestinationHistory, routePopPemanents []model.RoutePopPermanent) []response.VisitDay {
	var visitDays []response.VisitDay

	for _, perm := range routePopPemanents {
		if perm.RouteCode == nil {
			continue
		}
		var histories []response.Destination
		for _, history := range routeOutletsHistory {
			if history.RouteCode == *perm.RouteCode && isSameCalendarDate(history.Date, perm.Date) {
				if history.DestinationType == "distributor" {
					histories = append(histories, response.DistributorsHistory{
						ID:                 history.ID,
						RouteCode:          history.RouteCode,
						RouteName:          history.RouteName,
						DestinationID:      history.DestinationID,
						DestinationCode:    history.DestinationCode,
						DestinationName:    history.DestinationName,
						DestinationStatus:  history.DestinationStatus,
						DestinationAddress: history.DestinationAddress,
						Longitude:          history.Longitude,
						Latitude:           history.Latitude,
						AvgSalesWeek:       history.AvgSalesWeek,
						PjpID:              history.PjpID,
						PjpCode:            history.PjpCode,
						CustID:             history.CustID,
						DestinationType:    "distributor",
						IndexDay:           history.IndexDay,
						StartWeek:          history.StartWeek,
					})
				} else {
					histories = append(histories, response.OutletsHistory{
						ID:                 history.ID,
						RouteCode:          history.RouteCode,
						RouteName:          history.RouteName,
						DestinationID:      history.DestinationID,
						DestinationCode:    history.DestinationCode,
						DestinationName:    history.DestinationName,
						DestinationStatus:  history.DestinationStatus,
						DestinationAddress: history.DestinationAddress,
						Longitude:          history.Longitude,
						Latitude:           history.Latitude,
						AvgSalesWeek:       history.AvgSalesWeek,
						PjpID:              history.PjpID,
						PjpCode:            history.PjpCode,
						CustID:             history.CustID,
						DestinationType:    "outlet",
						IndexDay:           history.IndexDay,
						StartWeek:          history.StartWeek,
					})
				}
			}
		}
		if len(histories) == 0 {
			continue
		}

		visitDays = append(visitDays, response.VisitDay{
			ID:              perm.ID,
			Day:             perm.Day,
			IndexDay:        histories[0].GetIndexDay(),
			Week:            perm.Week,
			StartWeek:       histories[0].GetStartWeek().Format("2006-01-02"),
			Year:            perm.Year,
			Date:            perm.Date.Format("2006-01-02"),
			IsInCurrentYear: true,
			Visit: response.RoutesHistory{
				RouteCode: histories[0].GetRouteCode(),
				RouteName: histories[0].GetRouteName(),
				CustID:    perm.CustID,
				Outlets:   histories,
			},
		})
	}
	return visitDays
}

func isSameCalendarDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
