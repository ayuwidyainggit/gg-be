package routepop

import (
	"context"
	"log"
	"scyllax-pjp/data/request"
	"scyllax-pjp/exception"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strconv"
	"time"
)

func (service *routePopService) SaveDailyRouteMap(ctx context.Context, request request.SaveDailyRouteMap, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	for _, value := range request.Data {

		route, err := service.routeRepository.FindByRouteCode(ctx, value.RouteCode)
		if err != nil {
			log.Printf("Error finding route by route code %d: %v\n", value.RouteCode, err)
			panic(exception.NewNotFoundError(err.Error()))
		}

		// Parse PjpCode to integer
		resPjpCode, err := strconv.Atoi(value.PjpCode)
		if err != nil {
			panic(err)
		}

		var dateParse time.Time

		for _, week := range value.Weeks {
			dateTime, err := time.Parse("2006-01-02", week.Date) // yyyy-mm-dd
			dateParse = dateTime
			helper.ErrorPanic(err)

			// store in route pop daily
			routePopDaily := model.RoutePopDaily{
				RouteCode:   &value.RouteCode,
				Week:        week.Week,
				Day:         week.Day,
				Date:        dateTime,
				ParentRoute: &value.RouteCode,
				PjpCode:     &resPjpCode,
				PjpID:       &value.PjpID,
				Year:        week.Year,
				CustID:      currentCustomerId,
				Status:      "additional",
			}

			log.Printf("Inserting or updating RoutePopDaily: %+v\n", routePopDaily)

			service.routePopDailyRepository.UpdateOrCreateDaily(ctx, routePopDaily)
		}

		var routesHistory []model.RouteOutletHistory

		// Store outlet
		for _, outlet := range value.Outlets {
			routeOutlet := model.RouteOutletAdditional{
				RouteCode:     value.RouteCode,
				OutletID:      outlet.OutletID,
				OutletCode:    outlet.OutletCode,
				OutletName:    outlet.OutletName,
				Longitude:     outlet.Longitude,
				Latitude:      outlet.Latitude,
				OutletStatus:  strconv.Itoa(outlet.OutletStatus),
				OutletAddress: outlet.OutletAddress,
				AvgSalesWeek:  outlet.AvgSalesWeek,
				OldRouteCode:  value.RouteCode,
				Status:        "Approved",
				PjpID:         &value.PjpID,
				PjpCode:       &resPjpCode,
				OldPjpID:      &value.PjpID,
				OldPjpCode:    &resPjpCode,
				RouteName:     route.RouteName,
				CustID:        currentCustomerId,
				OldRouteName:  route.RouteName,
				Date:          dateParse,
			}
			log.Printf("Insert route outlet additional: %v\n", routeOutlet)
			service.routeOutletRepository.CreateAdditionalRoute(ctx, routeOutlet)

			var dayToIndex = map[string]int{
				"Monday":    1,
				"Tuesday":   2,
				"Wednesday": 3,
				"Thursday":  4,
				"Friday":    5,
				"Saturday":  6,
				"Sunday":    7,
			}

			startWeek := helper.GetStartOfISOWeek(value.Weeks[0].Year, value.Weeks[0].Week)

			routesHistory = append(routesHistory, model.RouteOutletHistory{
				PjpID:           &value.PjpID,
				PjpCode:         &resPjpCode,
				CustID:          currentCustomerId,
				RouteCode:       value.RouteCode,
				RouteName:       route.RouteName,
				OutletID:        outlet.OutletID,
				OutletCode:      outlet.OutletCode,
				OutletName:      outlet.OutletName,
				Longitude:       outlet.Longitude,
				Latitude:        outlet.Latitude,
				AvgSalesWeek:    outlet.AvgSalesWeek,
				OutletAddress:   outlet.OutletAddress,
				Week:            value.Weeks[0].Week,
				Year:            value.Weeks[0].Year,
				IndexDay:        dayToIndex[value.Weeks[0].Day],
				IsInCurrentYear: true,
				StartWeek:       &startWeek,
				Date:            dateParse,
				IsAdditional:    true,
			})
		}

		service.routeOutletHistoryRepo.CreateBulk(ctx, tx, routesHistory)
	}
}
