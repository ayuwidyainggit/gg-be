package pjpauto

import (
	"context"
	"fmt"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/utils"
	"strconv"
	"time"
)

func (service *pjpAutoService) Create(ctx context.Context, request request.CreatePjpAuto, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	// Delete assigned pjp before create pjp auto
	for _, routeOutlet := range request.Sales {

		resPjpCode, err := strconv.Atoi(routeOutlet.PjpCode)
		if err != nil {
			panic(err)
		}

		outlet := model.RouteOutlet{
			PjpID:   &routeOutlet.ID,
			PjpCode: &resPjpCode,
		}

		error := service.routeOutletRepository.UpdatePjpRouteOutlet(ctx, outlet)
		helper.ErrorPanic(error)

	}

	for i := range request.Day {

		// skip if data null
		if len(request.Route[i]) == 0 {
			continue
		}

		// log.Printf("Processing index: %d, Day: %s\n", i, day)
		// log.Println()
		var routeCode int
		var route model.Route

		// Generate unique route code
		for {
			routeCode = utils.GenerateCode(4)
			now := time.Now()
			// Format Route MMDDhhmmssi
			routeName := fmt.Sprintf("Route%s%06d%d", now.Format("0102"), now.Hour()*10000+now.Minute()*100+now.Second(), i)
			route = model.Route{
				RouteCode: routeCode,
				RouteName: routeName,
				IsAssign:  true,
				CustID:    currentCustomerId,
			}
			// Check if route with this code already exists
			existingRoute, err := service.routeRepository.FindByRouteCode(ctx, routeCode)
			if err == nil && existingRoute.RouteCode == routeCode {
				continue
			}
			break
		}

		// Insert the new route
		service.routeRepository.Insert(ctx, route)
		// log.Printf("Inserting route: %+v\n", route)
		// log.Println()

		// Create Route outlet
		for _, outlet := range request.Route[i] {
			// log.Printf("Processing outlet index: %d, Outlet: %+v\n", j, outlet)
			// log.Println()
			routeOutlet := model.RouteOutlet{
				RouteCode:     route.RouteCode,
				RouteName:     route.RouteName,
				OutletID:      outlet.OutletID,
				OutletCode:    outlet.OutletCode,
				OutletName:    outlet.OutletName,
				Longitude:     outlet.Longitude,
				Latitude:      outlet.Latitude,
				OutletStatus:  strconv.Itoa(outlet.OutletStatus),
				OutletAddress: outlet.OutletAddress,
				OldRouteCode:  route.RouteCode,
				OldRouteName:  route.RouteName,
			}
			// log.Printf("Creating route outlet: %+v\n", routeOutlet)
			// log.Println()
			service.routeOutletRepository.Create(ctx, routeOutlet)
		}

		// Assign route to PJP salesman
		salesman := request.Sales[i]

		//Parse PjpCode to integer
		resPjpCode, err := strconv.Atoi(salesman.PjpCode)
		if err != nil {
			panic(err)
		}

		routeOutlet := model.RouteOutlet{
			RouteCode:  route.RouteCode,
			RouteName:  route.RouteName,
			PjpID:      &salesman.ID,
			PjpCode:    &resPjpCode,
			OldPjpID:   &salesman.ID,
			OldPjpCode: &resPjpCode,
			CustID:     route.CustID,
		}
		// log.Printf("Assigning route to salesman: %+v\n", routeOutlet)
		// log.Println()
		// service.RouteRepository.Update(ctx, route)
		service.routeOutletRepository.Save(ctx, routeOutlet)

		// update pjp status to auto
		service.pjpRepository.Patch(ctx, tx, *routeOutlet.PjpID, "auto", currentCustomerId)
	}

}
