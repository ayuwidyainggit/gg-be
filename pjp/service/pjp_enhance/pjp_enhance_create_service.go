package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"time"

	"gorm.io/gorm"
)

func (service *pjpEnhanceService) Create(ctx context.Context, request request.CreatePjpEnhanceRequest, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	pjp := buildPjpModel(request, currentCustomerId)
	savedPjp := service.pjpRepository.Create(ctx, tx, pjp)

	savedRoutes := service.createRoutes(ctx, tx, savedPjp, request.Routes, currentCustomerId, nil)
	routePopPermanents := service.createVisitHistory(ctx, tx, savedPjp, savedRoutes, request.VisitDay, currentCustomerId)

	if routePopPermanents != nil {
		service.routePopRepository.CreateBulk(ctx, tx, routePopPermanents)
	}
}

func (service *pjpEnhanceService) createRoutes(
	ctx context.Context,
	tx *gorm.DB,
	savedPjp model.Pjp,
	routesReq []request.RoutesCreatePjp,
	currentCustomerId string,
	existingRoutes []model.Route,
) []model.Route {
	var savedRoutes []model.Route

	for indexRoute, routeReq := range routesReq {
		route := buildRoute(routeReq, savedPjp.ID, currentCustomerId, existingRoutes, indexRoute+1)
		route.Sequence = indexRoute + 1 // set sequence monday is 1, tuesday is 2, and so on
		savedRoute := service.routeRepository.Create(ctx, tx, route)
		savedRoutes = append(savedRoutes, savedRoute)

		outlets := buildRouteOutlets(routeReq.Destination, savedPjp, savedRoute, currentCustomerId)
		service.routeOutletRepository.CreateBulk(ctx, tx, outlets)
	}

	return savedRoutes
}

func (service *pjpEnhanceService) createVisitHistory(ctx context.Context, tx *gorm.DB, savedPjp model.Pjp, savedRoutes []model.Route, visitDays []request.VisitDayCreatePjp, currentCustomerId string) []model.RoutePopPermanent {
	var routePopPermanents []model.RoutePopPermanent

	for _, visit := range visitDays {
		var dateTime time.Time
		var startWeek time.Time
		var err error

		if visit.Date != "-" {
			dateTime, err = parseVisitDate(visit.Date)
			helper.ErrorPanic(err)

			startWeek, err = parseVisitDate(visit.Date)
			helper.ErrorPanic(err)
		}

		routeCode := getRouteCodeByName(savedRoutes, visit.Visit.RouteName)

		var outletHistories []model.RouteOutletHistory
		if len(visit.Visit.Destination) > 0 {
			for _, destination := range visit.Visit.Destination {
				outletHistories = append(outletHistories, model.RouteOutletHistory{
					PjpID:           &savedPjp.ID,
					PjpCode:         &savedPjp.PjpCode,
					CustID:          currentCustomerId,
					RouteCode:       routeCode,
					RouteName:       visit.Visit.RouteName,
					OutletID:        destination.ID,
					OutletCode:      destination.Code,
					OutletName:      destination.Name,
					Longitude:       destination.Longitude,
					Latitude:        destination.Latitude,
					AvgSalesWeek:    destination.AvgSalesWeek,
					OutletStatus:    destination.Status,
					OutletAddress:   destination.Address,
					Week:            visit.Week,
					Year:            visit.Year,
					IndexDay:        visit.IndexDay,
					IsInCurrentYear: visit.IsInCurrentYear,
					StartWeek:       &startWeek,
					Date:            dateTime,
					IsAdditional:    destination.IsAdditional,
				})
			}
			service.routeOutletHistoryRepository.CreateBulk(ctx, tx, outletHistories)

			routePopPermanents = append(routePopPermanents, model.RoutePopPermanent{
				RouteCode:            &routeCode,
				Week:                 visit.Week,
				Day:                  visit.Day,
				Date:                 dateTime,
				WorkingDayCalendarID: visit.WorkingDayCalendarID,
				PjpCode:              &savedPjp.PjpCode,
				PjpID:                &savedPjp.ID,
				Year:                 visit.Year,
				CustID:               currentCustomerId,
			})
		}
	}

	return routePopPermanents
}
