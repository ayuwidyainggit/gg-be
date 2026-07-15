package service

import (
	"context"
	"errors"
	"math"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/exception"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/repository"
	"scyllax-pjp/repository/pjp"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type RouteServiceImpl struct {
	RouteRepository             repository.RouteRepository
	RouteOutletRepository       repository.RouteOutletRepository
	OutletVisitRepo             repository.OutletVisitRepo
	RoutePopDailyRepository     repository.RoutePopDailyRepository
	RoutePopPermanentRepository repository.RoutePopPermanentRepository
	PjpRepository               pjp.PjpRepository
	Validate                    *validator.Validate
}

func NewRouteServiceImpl(
	routeRepository repository.RouteRepository,
	routeOutletRepository repository.RouteOutletRepository,
	routePopDailyRepository repository.RoutePopDailyRepository,
	routePopPermanentRepository repository.RoutePopPermanentRepository,
	outletVisitRepo repository.OutletVisitRepo,
	pjpRepository pjp.PjpRepository,
	validate *validator.Validate,
) RouteService {
	return &RouteServiceImpl{
		RouteRepository:             routeRepository,
		RouteOutletRepository:       routeOutletRepository,
		RoutePopDailyRepository:     routePopDailyRepository,
		RoutePopPermanentRepository: routePopPermanentRepository,
		OutletVisitRepo:             outletVisitRepo,
		PjpRepository:               pjpRepository,
		Validate:                    validate,
	}
}

func (service *RouteServiceImpl) Create(ctx context.Context, request request.CreateRouteRequest, currentCustomerId string) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	code := utils.GenerateCode(4)
	helper.ErrorPanic(err)

	dataset := model.Route{
		RouteCode: code,
		RouteName: request.RouteName,
		CustID:    currentCustomerId,
	}

	service.RouteRepository.Insert(ctx, dataset)
}

func (service *RouteServiceImpl) SaveOutlet(ctx context.Context, request request.SaveOutletRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	route, err := service.RouteRepository.FindByRouteCode(ctx, request.RouteCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	for _, outlet := range request.Outlets {
		route := model.Destination{
			RouteCode:          route.RouteCode,
			RouteName:          route.RouteName,
			DestinationID:      outlet.DestinationID,
			DestinationCode:    outlet.DestinationCode,
			DestinationName:    outlet.DestinationName,
			Longitude:          outlet.Longitude,
			Latitude:           outlet.Latitude,
			AvgSalesWeek:       outlet.AvgSalesWeek,
			DestinationStatus:  strconv.Itoa(outlet.DestinationStatus),
			DestinationAddress: outlet.DestinationAddress,
			OldRouteCode:       route.RouteCode,
			OldRouteName:       route.RouteName,
		}
		service.RouteOutletRepository.Create(ctx, route)
	}
}

// TODO remove one route to one pjp
/*
	make condition to check wheteher the code in db there is a route outlet where pjp_code and pjp_id is null
	if null just update the pjp_code and pjp_id if not create new route outlet and assign pjp_id and pjp_code
	1. if route outlet but
*/
func (service *RouteServiceImpl) SavePjp(ctx context.Context, request request.SavePjpRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, routeCode := range request.RouteCode {
		route, err := service.RouteRepository.FindByRouteCode(ctx, routeCode)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		// Parse PjpCode to integer
		resPjpCode, err := strconv.Atoi(request.PjpCode)
		if err != nil {
			panic(err)
		}

		routes := model.Destination{
			RouteCode:  route.RouteCode,
			RouteName:  route.RouteName,
			PjpID:      &request.PjpID,
			PjpCode:    &resPjpCode,
			OldPjpID:   &request.PjpID,
			OldPjpCode: &resPjpCode,
			CustID:     route.CustID,
		}
		service.RouteRepository.Update(ctx, route)
		service.RouteOutletRepository.Save(ctx, routes)

		//fmt.Println("service", routes)
	}

	// for _, routeCode := range request.RouteCode {
	// 	route, err := service.RouteRepository.FindByRouteCode(ctx, routeCode)
	// 	if err != nil {
	// 		panic(exception.NewNotFoundError(err.Error()))
	// 	}

	// 	fmt.Println(route)
	// 	outlets, err := service.RouteRepository.FindAllByRouteCode(ctx, routeCode)
	// 	if err != nil {
	// 		panic(exception.NewNotFoundError(err.Error()))
	// 	}

	// 	//Parse PjpCode to integer
	// 	resPjpCode, err := strconv.Atoi(request.PjpCode)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	//As route can  be assign to many pjp it should be false
	// 	// route.IsAssign = false
	// 	// routes := model.Destination{
	// 	// 	RouteCode:  route.RouteCode,
	// 	// 	RouteName:  route.RouteName,
	// 	// 	PjpID:      &request.PjpID,
	// 	// 	PjpCode:    &resPjpCode,
	// 	// 	OldPjpID:   &request.PjpID,
	// 	// 	OldPjpCode: &resPjpCode,
	// 	// 	CustID:     route.CustID,
	// 	// }
	// 	// service.RouteRepository.Update(ctx, route)
	// 	// service.RouteOutletRepository.Save(ctx, routes)

	// 	for _, outlet := range outlets {
	// 		existingRouteOutlet, err := service.RouteOutletRepository.FindByRouteCodeAndDestinationIDAndPjpNull(ctx, route.RouteCode, outlet.DestinationID)
	// 		if err != nil {
	// 			panic(exception.NewNotFoundError(err.Error()))
	// 		}
	// 		if existingRouteOutlet != nil {
	// 			// update exist data
	// 			existingRouteOutlet.PjpID = &request.PjpID
	// 			existingRouteOutlet.PjpCode = &resPjpCode
	// 			existingRouteOutlet.OldPjpID = &request.PjpID
	// 			existingRouteOutlet.OldPjpCode = &resPjpCode
	// 			existingRouteOutlet.RouteCode = route.RouteCode
	// 			existingRouteOutlet.RouteName = route.RouteName
	// 			existingRouteOutlet.CustID = route.CustID
	// 			service.RouteOutletRepository.Save(ctx, *existingRouteOutlet)
	// 		} else {
	// 			// Create new route outlet
	// 			routeOutlet := model.Destination{
	// 				RouteCode:     route.RouteCode,
	// 				RouteName:     route.RouteName,
	// 				DestinationID:      outlet.DestinationID,
	// 				PjpID:         &request.PjpID,
	// 				PjpCode:       &resPjpCode,
	// 				DestinationCode:    outlet.DestinationCode,
	// 				OldPjpID:      &request.PjpID,
	// 				OldPjpCode:    &resPjpCode,
	// 				DestinationName:    outlet.DestinationName,
	// 				Longitude:     outlet.Longitude,
	// 				Latitude:      outlet.Latitude,
	// 				DestinationStatus:  outlet.DestinationStatus,
	// 				DestinationAddress: outlet.DestinationAddress,
	// 				OldRouteCode:  route.RouteCode,
	// 				OldRouteName:  route.RouteName,
	// 				CustID:        route.CustID,
	// 			}
	// 			service.RouteOutletRepository.Create(ctx, routeOutlet)
	// 		}
	// 		//fmt.Println("service", routes)
	// 	}
	// }
}

func (service *RouteServiceImpl) UpdatePjp(ctx context.Context, request request.UpdatePjpInRouteRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, route := range request.RouteCode {
		routeOutlet, err := service.RouteOutletRepository.FindByRouteCodeAndPjpCode(ctx, route, request.PjpCode)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		for _, data := range routeOutlet {
			data.PjpID = nil
			data.PjpCode = nil

			service.RouteOutletRepository.Save(ctx, data)
		}
	}
}

func (service *RouteServiceImpl) DeleteOutlet(ctx context.Context, request request.DeleteOutletRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, DestinationCode := range request.DestinationCode {
		route := model.Destination{
			RouteCode:       request.RouteCode,
			DestinationCode: DestinationCode,
		}
		err = service.RouteOutletRepository.DeleteByDestinationCode(ctx, route)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}
	}
}

func (service *RouteServiceImpl) DeleteOutletAdditional(ctx context.Context, request request.DeleteOutletAdditionalRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, DestinationCode := range request.DestinationCode {
		route := model.DestinationAdditional{
			RouteCode:       request.RouteCode,
			DestinationCode: DestinationCode,
		}
		err = service.RouteOutletRepository.DeleteByDestinationCodeAdditional(ctx, route)
		if err != nil {
			helper.ErrorPanic(err)
		}
		err = service.OutletVisitRepo.Delete(ctx, request.Date, request.Week, DestinationCode, request.RouteCode)
		if err != nil {
			helper.ErrorPanic(err)
		}
	}

}

func (service *RouteServiceImpl) FindAll(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingResponse, response.Meta, error) {
	result := service.RouteOutletRepository.FindAll(ctx, page, limit, filters, currentCustomerId)

	var routes []response.ApprovalRouteMappingResponse

	for _, value := range result {
		routeResponse := response.ApprovalRouteMappingResponse{
			ID:           value.ID,
			RouteCode:    value.OldRouteCode,
			RouteName:    value.OldRouteName,
			NewRouteName: value.RouteName,
			NewRouteCode: value.RouteCode,
			Status:       value.Status,
			Date:         value.CreatedAt.Format("2006-01-02 15:04"),
		}

		if value.VerifiedDate != nil {
			routeResponse.VerifiedDate = time.Now().Format("2006-01-02 15:04")
		}

		if value.PjpID != nil || value.PjpCode != nil {
			routeResponse.PjpID = value.PjpID
			routeResponse.PjpCode = value.PjpCode
		}

		if value.Pjp != nil && value.Pjp.SalesmanCode != "" {
			routeResponse.SalesmanCode = &value.Pjp.SalesmanCode
		}

		if value.Pjp != nil && value.Pjp.SalesmanName != "" {
			routeResponse.SalesmanName = &value.Pjp.SalesmanName
		}

		if value.PjpOld != nil && value.PjpOld.SalesmanCode != "" {
			routeResponse.NewSalesmanCode = &value.PjpOld.SalesmanCode
		}

		if value.PjpOld != nil && value.PjpOld.SalesmanName != "" {
			routeResponse.NewSalesmanName = &value.PjpOld.SalesmanName
		}

		if value.DestinationCode != "" {
			outlet := response.DestinationResponse{
				DestinationID:      value.DestinationID,
				DestinationCode:    value.DestinationCode,
				DestinationName:    value.DestinationName,
				Longitude:          value.Longitude,
				Latitude:           value.Latitude,
				DestinationStatus:  value.DestinationStatus,
				DestinationAddress: value.DestinationAddress,
				AvgSalesWeek:       value.AvgSalesWeek,
				Status:             value.Status,
			}
			routeResponse.Outlets = &outlet
		} else {
			routeResponse.Outlets = nil
		}

		routes = append(routes, routeResponse)
	}

	totalData, _ := service.RouteOutletRepository.Count(ctx, currentCustomerId)

	pagination := &response.Meta{
		TotalData: int(totalData),
		Page:      page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(totalData) / float64(limit))),
	}

	return routes, *pagination, nil
}
func (service *RouteServiceImpl) FindAllEnhance(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingEnhanceResponse, response.Meta, error) {
	result := service.RouteOutletRepository.FindAllEnhance(ctx, page, limit, filters, currentCustomerId)

	var routes []response.ApprovalRouteMappingEnhanceResponse

	for _, value := range result {
		var verifiedDate string
		if len(value.RouteOutlets) > 0 && value.RouteOutlets[0].VerifiedDate != nil {
			verifiedDate = value.RouteOutlets[0].VerifiedDate.Format("2006-01-02 15:04")
		} else {
			verifiedDate = ""
		}

		routeResponse := response.ApprovalRouteMappingEnhanceResponse{
			PjpID:        value.ID,
			PjpCode:      helper.FormatPjpCode(value.PjpCode),
			Status:       value.ApprovalStatus,
			SalesmanName: value.SalesmanName,
			SalesmanCode: value.SalesmanCode,
			Date:         value.CreatedAt.Format("2006-01-02 15:04"),
			VerifiedDate: verifiedDate,
		}

		routes = append(routes, routeResponse)
	}

	totalData, _ := service.RouteOutletRepository.CountAllEnhance(ctx, currentCustomerId)

	pagination := &response.Meta{
		TotalData: int(totalData),
		Page:      page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(totalData) / float64(limit))),
	}

	return routes, *pagination, nil
}

// TODO will evaluate total outlet as the code changes
func (service *RouteServiceImpl) FindAllRoute(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RouteResponse {
	result := service.RouteRepository.FindAll(ctx, filters, currentCustomerId)

	var responses []response.RouteResponse
	for _, value := range result {
		totalOutlet := len(value.Destinations) // TODO

		data := response.RouteResponse{
			ID:          value.ID,
			RouteName:   value.RouteName,
			RouteCode:   value.RouteCode,
			TotalOutlet: totalOutlet,
		}

		responses = append(responses, data)
	}

	return responses
}

func (service *RouteServiceImpl) FindByRouteOutlet(ctx context.Context, routeCode, pjpCode int) []response.RouteOutletsResponse {

	data := service.RouteOutletRepository.FindByRouteCodes(ctx, routeCode, pjpCode)

	outletMap := make(map[int][]response.DestinationResponse)
	routeNames := make(map[int]string)

	for _, outlet := range data {
		outletMap[outlet.RouteCode] = append(outletMap[outlet.RouteCode], response.DestinationResponse{
			DestinationID:      outlet.DestinationID,
			DestinationCode:    outlet.DestinationCode,
			DestinationName:    outlet.DestinationName,
			Longitude:          outlet.Longitude,
			Latitude:           outlet.Latitude,
			DestinationStatus:  outlet.DestinationStatus,
			DestinationAddress: outlet.DestinationAddress,
			AvgSalesWeek:       outlet.AvgSalesWeek,
			Status:             outlet.Status,
		})

		if _, ok := routeNames[outlet.RouteCode]; !ok {
			routeNames[outlet.RouteCode] = outlet.RouteName
		}
	}

	var result []response.RouteOutletsResponse

	for routeCode, outlets := range outletMap {
		routeName := routeNames[routeCode]
		result = append(result, response.RouteOutletsResponse{
			RouteCode: routeCode,
			RouteName: routeName,
			Outlets:   outlets,
		})
	}

	//fmt.Println("cel", result)
	return result
}

func (service *RouteServiceImpl) Delete(ctx context.Context, routeId int) {
	err := service.RouteRepository.Delete(ctx, routeId)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *RouteServiceImpl) FindRouteByPjpCode(ctx context.Context, pjpCode, routeCode int) []response.RouteResponse {
	result := service.RouteRepository.FindByPjpCode(ctx, pjpCode, routeCode)

	seenRouteNames := make(map[string]bool)
	totalOutlet := len(result)
	var data []response.RouteResponse
	for _, row := range result {
		if !seenRouteNames[row.RouteName] {
			var res response.RouteResponse
			helper.Automapper(row, &res)
			res.TotalOutlet = totalOutlet

			var outlets []response.DestinationResponse
			for _, value := range result {
				if value.RouteCode == row.RouteCode && value.DestinationCode != "" {
					outlet := response.DestinationResponse{
						DestinationID:      value.DestinationID,
						DestinationCode:    value.DestinationCode,
						DestinationName:    value.DestinationName,
						Longitude:          value.Longitude,
						Latitude:           value.Latitude,
						DestinationStatus:  value.DestinationStatus,
						DestinationAddress: value.DestinationAddress,
						AvgSalesWeek:       value.AvgSalesWeek,
						Status:             value.Status,
					}
					outlets = append(outlets, outlet)
				}
			}
			res.Destinations = &outlets
			data = append(data, res)
			seenRouteNames[row.RouteName] = true
		}
	}

	return data
}

func (service *RouteServiceImpl) FindDailyRouteByPjpCode(ctx context.Context, pjpCode, routeCode int, date string) []response.RouteDailyResponse {
	result := service.RouteRepository.FindByPjpCodeRouteAdditional(ctx, pjpCode, routeCode, date)

	seenRouteNames := make(map[string]bool)
	totalOutlet := len(result)
	var data []response.RouteDailyResponse
	for _, row := range result {
		if !seenRouteNames[row.RouteName] {
			var res response.RouteDailyResponse
			helper.Automapper(row, &res)
			res.TotalOutlet = totalOutlet

			var outlets []response.OutletDailyResponse
			for _, value := range result {
				if value.RouteCode == row.RouteCode && value.DestinationCode != "" {
					DestinationStatus, err := strconv.Atoi(value.DestinationStatus)
					if err != nil {
						// Handle the error, maybe set a default value or log the error
						DestinationStatus = 0 // Default value or handle as needed
					}
					// routePop := service.RoutePopDailyRepository.FindStatusRoutePop(ctx, pjpCode, routeCode)
					outlet := response.OutletDailyResponse{
						DestinationID:      value.DestinationID,
						DestinationCode:    value.DestinationCode,
						DestinationName:    value.DestinationName,
						Longitude:          value.Longitude,
						Latitude:           value.Latitude,
						DestinationStatus:  DestinationStatus,
						DestinationAddress: value.DestinationAddress,
						Status:             value.RoutePopStatus,
						AvgSalesWeek:       value.AvgSalesWeek,
						PjpCode:            pjpCode,
						RouteCode:          routeCode,
						PjpID:              value.PjpID,
					}
					outlets = append(outlets, outlet)
				}
			}
			res.Destinations = &outlets
			data = append(data, res)
			seenRouteNames[row.RouteName] = true
		}
	}

	return data
}

func (service *RouteServiceImpl) UpdateRoute(ctx context.Context, request request.UpdateRoutesRequest, currentCustomerId string) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	data, err := service.RouteRepository.FindById(ctx, request.ID, currentCustomerId)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	data.RouteName = request.RouteName

	service.RouteRepository.Update(ctx, data)
	service.RouteOutletRepository.Update(ctx, data.RouteCode, data.RouteName)
}

func (service *RouteServiceImpl) SaveRouteConfirmation(ctx context.Context, request request.SaveRouteConfirmationRequest) {
	for _, route := range request.Routes {
		err := service.Validate.Struct(route)
		helper.ErrorPanic(err)

		data, err := service.RouteRepository.FindByRouteCode(ctx, route.RouteCode)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}
		routeOutlet := model.Destination{
			RouteCode:          data.RouteCode,
			RouteName:          data.RouteName,
			DestinationID:      route.DestinationID,
			DestinationCode:    route.DestinationCode,
			DestinationName:    route.DestinationName,
			Longitude:          route.Longitude,
			Latitude:           route.Latitude,
			DestinationStatus:  strconv.Itoa(route.DestinationStatus),
			DestinationAddress: route.DestinationAddress,
		}
		service.RouteOutletRepository.Create(ctx, routeOutlet)
	}
}

func (service *RouteServiceImpl) DeletePjp(ctx context.Context, request request.DeletePjpRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	route, err := service.RouteRepository.FindByRouteCode(ctx, request.RouteCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// Parse PjpCode to integer
	resPjpCode, err := strconv.Atoi(request.PjpCode)
	if err != nil {
		panic(err)
	}

	routes := model.Destination{
		RouteCode: route.RouteCode,
		PjpCode:   &resPjpCode,
		PjpID:     &request.PjpID,
	}

	error := service.RouteOutletRepository.UpdatePjp(ctx, routes)
	helper.ErrorPanic(error)

	err = service.RoutePopDailyRepository.DeleteByRouteCode(ctx, request.RouteCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	err = service.RoutePopPermanentRepository.DeleteByRouteCode(ctx, request.RouteCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *RouteServiceImpl) UpdateNewRoute(ctx context.Context, request request.NewRouteRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, route := range request.NewRoutePropose {

		data, err := service.RouteRepository.FindByRouteCode(ctx, route.OldRouteCode)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		//Parse PjpCode to integer
		resPjpCode, err := strconv.Atoi(route.PjpCode)
		if err != nil {
			panic(err)
		}

		routeOutlet := model.Destination{
			DestinationID:   route.DestinationID,
			DestinationCode: route.DestinationCode,
			RouteCode:       route.RouteCode,
			RouteName:       route.RouteName,
			PjpID:           &route.PjpID,
			PjpCode:         &resPjpCode,
			OldPjpID:        &route.OldPjpID,
			OldPjpCode:      &route.OldPjpCode,
			OldRouteCode:    data.RouteCode,
			OldRouteName:    data.RouteName,
		}
		service.RouteOutletRepository.UpdateNewRoute(ctx, routeOutlet)
	}
}

func (service *RouteServiceImpl) FindByRouteCode(ctx context.Context, routeCode int) []response.RouteResponse {
	result := service.RouteRepository.QueryByRouteCode(ctx, routeCode)

	var data []response.RouteResponse
	routeMap := make(map[int]*response.RouteResponse)

	for _, row := range result {
		if _, exists := routeMap[row.RouteCode]; !exists {
			var res response.RouteResponse
			helper.Automapper(row, &res)
			res.Destinations = nil
			routeMap[row.RouteCode] = &res
		}

		if row.DestinationCode != "" {
			var outlet response.DestinationResponse
			helper.Automapper(row, &outlet)
			if routeMap[row.RouteCode].Destinations == nil {
				routeMap[row.RouteCode].Destinations = &[]response.DestinationResponse{}
			}
			*routeMap[row.RouteCode].Destinations = append(*routeMap[row.RouteCode].Destinations, outlet)
		}
	}

	for _, route := range routeMap {
		data = append(data, *route)
	}

	return data
}

func (service *RouteServiceImpl) DuplicateRoute(ctx context.Context, request request.DuplicateRoute, currentCustomerId string) error {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	routes, err := service.RouteRepository.FindRouteOutletByRouteCode(ctx, request.RouteCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if len(routes.Destinations) == 0 {
		return errors.New("route has no outlets to duplicate")
	}

	newRouteCode := utils.GenerateCode(4)
	helper.ErrorPanic(err)
	newRouteName := "Copy of " + routes.RouteName

	NewRoute := model.Route{
		RouteCode: newRouteCode,
		RouteName: newRouteName,
		CustID:    currentCustomerId,
	}

	service.RouteRepository.Insert(ctx, NewRoute)

	for _, outlet := range routes.Destinations {
		routeOutlet := model.Destination{
			RouteCode:          newRouteCode,
			RouteName:          newRouteName,
			DestinationID:      outlet.DestinationID,
			DestinationCode:    outlet.DestinationCode,
			DestinationName:    outlet.DestinationName,
			Longitude:          outlet.Longitude,
			Latitude:           outlet.Latitude,
			DestinationStatus:  outlet.DestinationStatus,
			DestinationAddress: outlet.DestinationAddress,
		}
		service.RouteOutletRepository.Create(ctx, routeOutlet)
	}

	return nil
}
