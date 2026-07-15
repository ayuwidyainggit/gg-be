package routepop

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
)

func (service *routePopService) GetByPjpAndRouteCode(ctx context.Context, pjpCode, routeCode int, date, currentCustomerId string) response.DailyRouteMap {

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	histories := service.destinationHistoryRepo.FindByPjpCodeRouteCodeAndDate(ctx, tx, pjpCode, routeCode, date, currentCustomerId)

	return mapToDestinations(histories)

}

func mapToDestinations(histories []model.DestinationHistory) response.DailyRouteMap {
	var destinations response.DailyRouteMap

	for _, history := range histories {
		if history.DestinationType == "distributor" {
			destinations.Destination = append(destinations.Destination, response.DistributorsHistory{
				ID:                 history.ID,
				RouteCode:          history.RouteCode,
				RouteName:          history.RouteName,
				DestinationID:      history.ID,
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
				Year:               history.Year,
				Week:               history.Week,
				Status:             getStatus(history.IsAdditional),
			})
		} else {
			destinations.Destination = append(destinations.Destination, response.OutletsHistory{
				ID:                 history.ID,
				RouteCode:          history.RouteCode,
				RouteName:          history.RouteName,
				DestinationID:      history.ID,
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
				Year:               history.Year,
				Week:               history.Week,
				Status:             getStatus(history.IsAdditional),
			})
		}
	}

	destinations.CustID = histories[0].CustID
	destinations.RouteCode = histories[0].RouteCode
	destinations.RouteName = histories[0].RouteName
	return destinations
}

func getStatus(isAdditional bool) string {
	if isAdditional {
		return "additional"
	}
	return "permanent"
}
