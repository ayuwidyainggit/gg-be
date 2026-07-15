package routepop

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
)

func (service *routePopService) GetAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RoutePopPermanentResponse {
	var result []response.RoutePopPermanentResponse

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	histories := service.destinationHistoryRepo.GetAll(ctx, tx, filters)
	if len(histories) == 0 || histories[0].PjpID == nil {
		return result
	}

	pjp := service.pjpRepo.GetPjpById(ctx, tx, *histories[0].PjpID, currentCustomerId)

	// Grouping histories by RouteCode
	grouped := make(map[int][]model.DestinationHistory)
	for _, h := range histories {
		grouped[h.RouteCode] = append(grouped[h.RouteCode], h)
	}

	var routes []response.RoutesMap

	// Build routes array per group
	for _, group := range grouped {
		first := group[0]
		var totalOutlet, totalDistributor int

		for _, h := range group {
			if h.DestinationType == "outlet" {
				totalOutlet++
			} else if h.DestinationType == "distributor" {
				totalDistributor++
			}
		}

		routes = append(routes, response.RoutesMap{
			RouteCode:        first.RouteCode,
			RouteName:        first.RouteName,
			Week:             first.Week,
			Date:             first.Date,
			TotalOutlet:      totalOutlet,
			TotalDistributor: totalDistributor,
		})
	}

	// Satu response dengan semua routes
	result = append(result, response.RoutePopPermanentResponse{
		PjpID:        &pjp.ID,
		PjpCode:      &pjp.PjpCode,
		SalesmanName: &pjp.SalesmanName,
		Routes:       routes,
	})

	return result
}
