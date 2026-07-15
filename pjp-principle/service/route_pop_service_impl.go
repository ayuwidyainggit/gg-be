package service

import (
	"context"
	"fmt"
	"log"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/exception"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/repository"
	"scyllax-pjp/repository/pjp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type RoutePopServiceImpl struct {
	RoutePopPermanentRepository repository.RoutePopPermanentRepository
	RoutePopDailyRepository     repository.RoutePopDailyRepository
	PjpRepository               pjp.PjpRepository
	RouteRepository             repository.RouteRepository
	RouteOutletRepository       repository.RouteOutletRepository
	OutletVisitListRepo         repository.OutletVisitRepo
	Validate                    *validator.Validate
}

func NewRoutePopServiceImpl(routePopPermanentRepository repository.RoutePopPermanentRepository, routePopDailyRepository repository.RoutePopDailyRepository, pjpRepository pjp.PjpRepository, routeRepository repository.RouteRepository, routeOutletRepository repository.RouteOutletRepository, OutletVisitListRepo repository.OutletVisitRepo, validate *validator.Validate) RoutePopService {
	return &RoutePopServiceImpl{
		RoutePopPermanentRepository: routePopPermanentRepository,
		RoutePopDailyRepository:     routePopDailyRepository,
		PjpRepository:               pjpRepository,
		RouteRepository:             routeRepository,
		RouteOutletRepository:       routeOutletRepository,
		OutletVisitListRepo:         OutletVisitListRepo,
		Validate:                    validate,
	}
}

func (service *RoutePopServiceImpl) SaveWeekly(ctx context.Context, request request.SaveWeeklyRequest, currentCustomerId string) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, value := range request.Data {
		for _, week := range value.Weeks {
			var dateTime time.Time
			if week.Date != "-" {
				// First try to parse with RFC3339 format
				dateTime, err = time.Parse(time.RFC3339, week.Date)
				if err != nil {
					// If RFC3339 fails, try the original format
					dateTime, err = time.Parse("2006-01-02", week.Date)
					helper.ErrorPanic(err)
				}
			}

			// Parse PjpCode to integer
			resPjpCode, err := strconv.Atoi(value.PjpCode)
			if err != nil {
				panic(err)
			}

			if week.Day == "no visit" || week.Day == "-" || week.Date == "-" || week.Day == "" {
				service.RoutePopPermanentRepository.DeleteByParams(ctx, value.RouteCode, value.PjpID, resPjpCode, week.Year, week.Week, currentCustomerId)
				service.RoutePopDailyRepository.DeleteByParams(ctx, value.RouteCode, value.PjpID, resPjpCode, week.Year, week.Week, currentCustomerId)
			} else {
				routePopPermanent := model.RoutePopPermanent{
					RouteCode: &value.RouteCode,
					Week:      week.Week,
					Day:       week.Day,
					Date:      dateTime,
					PjpCode:   &resPjpCode,
					PjpID:     &value.PjpID,
					Year:      week.Year,
					CustID:    currentCustomerId,
				}

				routePopDaily := model.RoutePopDaily{
					RouteCode: &value.RouteCode,
					Week:      week.Week,
					Day:       week.Day,
					Date:      dateTime,
					PjpCode:   &resPjpCode,
					PjpID:     &value.PjpID,
					Year:      week.Year,
					CustID:    currentCustomerId,
					Status:    "permanent",
				}
				service.RoutePopPermanentRepository.UpdateOrCreate(ctx, routePopPermanent)
				service.RoutePopDailyRepository.UpdateOrCreate(ctx, routePopDaily)
			}
		}
	}
}

func (service *RoutePopServiceImpl) SaveDelegateRoute(ctx context.Context, request request.SaveDelegateRequest, currentCustomerId string) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	dateTime, err := time.Parse("2006-01-02", request.Date) // yyyy-mm-dd
	helper.ErrorPanic(err)

	// Parse PjpCode to integer
	resPjpCode, err := strconv.Atoi(request.PjpCode)
	if err != nil {
		panic(err)
	}

	routePopDailyPermanent := model.RoutePopDaily{
		RouteCode: &request.RouteCode,
		Week:      request.Week + 1,
		Day:       request.Day,
		Date:      dateTime,
		PjpCode:   &resPjpCode,
		PjpID:     &request.PjpID,
		Year:      request.Year,
		Status:    "additional",
		CustID:    currentCustomerId,
	}
	service.RoutePopDailyRepository.Insert(ctx, routePopDailyPermanent)
}

func (service *RoutePopServiceImpl) CopyAllPermanentToDaily(ctx context.Context, request request.CopyAllRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	routes, err := service.RoutePopPermanentRepository.FindByPjpCodes(ctx, request.PjpCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	for _, route := range routes {
		dataset := model.RoutePopDaily{
			RouteCode: route.RouteCode,
			Week:      route.Week + 1,
			Day:       route.Day,
			Date:      route.Date,
			PjpCode:   route.PjpCode,
			PjpID:     route.PjpID,
			Year:      route.Year,
			CustID:    route.CustID,
		}
		service.RoutePopDailyRepository.Insert(ctx, dataset)
	}
}

func (service *RoutePopServiceImpl) CopyPartialPermanentToDaily(ctx context.Context, request request.CopyPartialRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, routeData := range request.Data {
		pjp, err := service.RoutePopPermanentRepository.FindByPjpCode(ctx, routeData.PjpCode)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		for _, route := range routeData.Routes {
			dataset := model.RoutePopDaily{
				RouteCode: &route.RouteCode,
				Week:      pjp.Week + 1,
				Day:       pjp.Day,
				Date:      pjp.Date,
				PjpCode:   pjp.PjpCode,
				PjpID:     pjp.PjpID,
				Year:      pjp.Year,
				Status:    "additional_route",
				CustID:    pjp.CustID,
			}
			service.RoutePopDailyRepository.Insert(ctx, dataset)
		}
	}
}

func (service *RoutePopServiceImpl) CopyToSpecificDaily(ctx context.Context, request request.RoutesMapping, currentCustomerId string) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, route := range request.Routes {
		parsedTime, err := time.Parse("2006-01-02", request.Date) //yyyy-mm-dd
		if err != nil {
			helper.ErrorPanic(err)
		}
		dataset := model.RoutePopDaily{
			RouteCode:   &route.RouteCode,
			Week:        request.Week + 1,
			Day:         request.Day,
			Date:        parsedTime,
			PjpCode:     &request.PjpCode,
			PjpID:       &request.PjpID,
			Year:        request.Year,
			Status:      "additional_route_outlet",
			ParentRoute: &request.RouteCode,
			CustID:      currentCustomerId,
		}
		service.RoutePopDailyRepository.Insert(ctx, dataset)
	}
}

func (service *RoutePopServiceImpl) CopyRouteDailyToDaily(ctx context.Context, request request.RoutesMapping, currentCustomerId string) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	for _, route := range request.Routes {
		parsedTime, err := time.Parse("2006-01-02", request.Date) //yyyy-mm-dd
		if err != nil {
			helper.ErrorPanic(err)
		}
		dataset := model.RoutePopDaily{
			RouteCode:   &route.RouteCode,
			Week:        request.Week + 1,
			Day:         request.Day,
			Date:        parsedTime,
			PjpCode:     &request.PjpCode,
			PjpID:       &request.PjpID,
			Year:        request.Year,
			Status:      "additional_route_outlet",
			ParentRoute: &request.RouteCode,
			CustID:      currentCustomerId,
		}
		service.RoutePopDailyRepository.Insert(ctx, dataset)
	}
}

func (service *RoutePopServiceImpl) FindAllPermanent(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RoutePopPermanentResponse {
	result := service.RoutePopPermanentRepository.FindAll(ctx, filters, currentCustomerId)
	log.Printf("Data: %+v", result)

	var weekNumber int
	switch v := filters["week"].(type) {
	case int:
		weekNumber = v
	case float64:
		weekNumber = int(v)
	case string:
		parsedWeek, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Error converting `week` filter to int: %v", err)
			return nil
		}
		weekNumber = parsedWeek
	default:
		log.Printf("Invalid type for `week` filter: %T", v)
		return nil
	}

	start, end, err := getWeekRangeFromNumber(weekNumber)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}

	query := fmt.Sprintf("BETWEEN '%s' AND '%s'", start, end)
	log.Println("Query filter:", query)
	outletCounts := service.RoutePopPermanentRepository.CountOutletByRoute(ctx, currentCustomerId, start, end)
	log.Printf("Fetched outlet counts for customer ID: %s: %+v", currentCustomerId, outletCounts)
	mapperByPjp := make(map[int]response.RoutePopPermanentResponse)

	for _, value := range result {
		key := *value.PjpCode
		routePopPermanentResponse, exists := mapperByPjp[key]
		if !exists {
			routePopPermanentResponse = response.RoutePopPermanentResponse{
				ID:      value.ID,
				PjpCode: value.PjpCode,
				PjpID:   value.PjpID,
			}
		}

		if value.Pjp != nil {
			routePopPermanentResponse.SalesmanName = &value.Pjp.SalesmanName
		}

		isDuplicate := false
		for _, route := range routePopPermanentResponse.Routes {
			if route.RouteCode == *value.RouteCode {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			totalOutlet := 0
			routeName := ""
			if count, ok := outletCounts[*value.RouteCode]; ok {
				totalOutlet = count.TotalOutlet
				routeName = count.RouteName
			}
			routePopPermanentResponse.Routes = append(routePopPermanentResponse.Routes, response.RoutesMap{
				RouteCode:   *value.RouteCode,
				RouteName:   routeName,
				Week:        value.Week,
				Date:        value.Date,
				TotalOutlet: totalOutlet,
			})
		}

		mapperByPjp[key] = routePopPermanentResponse
	}

	var routes []response.RoutePopPermanentResponse
	for _, routeResponse := range mapperByPjp {
		routes = append(routes, routeResponse)
	}

	return routes
}

func (service *RoutePopServiceImpl) FindAllDaily(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RoutePopDailyResponse {
	result := service.RoutePopDailyRepository.FindAll(ctx, filters, currentCustomerId)
	mapperByPjp := make(map[int]response.RoutePopDailyResponse)
	routeCodes := make(map[int]bool)

	for _, value := range result {
		key := *value.PjpCode
		routePopDailyResponse, exists := mapperByPjp[key]
		if !exists {
			routePopDailyResponse = response.RoutePopDailyResponse{
				ID:      value.ID,
				PjpCode: value.PjpCode,
				PjpID:   value.PjpID,
				Status:  value.Status,
				Week:    value.Week,
			}
		}

		if value.Pjp != nil {
			routePopDailyResponse.SalesmanName = &value.Pjp.SalesmanName
		}

		var outlets []response.OutletMap
		if value.Route != nil {
			for _, routeOutlet := range value.Route.Destinations {
				outlet := response.OutletMap{
					DestinationID:      routeOutlet.DestinationID,
					DestinationCode:    routeOutlet.DestinationCode,
					DestinationName:    routeOutlet.DestinationName,
					Longitude:          routeOutlet.Longitude,
					Latitude:           routeOutlet.Latitude,
					DestinationStatus:  routeOutlet.DestinationStatus,
					DestinationAddress: routeOutlet.DestinationAddress,
					AvgSalesWeek:       routeOutlet.AvgSalesWeek,
				}
				outlets = append(outlets, outlet)
			}
		}

		if value.RouteCode != nil {
			routeCode := *value.RouteCode
			if _, exists := routeCodes[routeCode]; !exists {
				routeCodes[routeCode] = true
				routeMap := response.RouteMap{
					RouteCode: &routeCode,
					Outlets:   outlets,
				}
				routePopDailyResponse.Routes = append(routePopDailyResponse.Routes, routeMap)
			}
		}

		mapperByPjp[key] = routePopDailyResponse
	}

	var routes []response.RoutePopDailyResponse

	for _, routeResponse := range mapperByPjp {
		routes = append(routes, routeResponse)
	}

	if len(routes) == 0 {
		return nil
	}

	return routes
}

func (service *RoutePopServiceImpl) FindByRouteOutletAdditional(ctx context.Context, code int, currentCustomerId string) response.RouteDetailResponse {
	data, err := service.RoutePopDailyRepository.FindByParentRoute(ctx, code, currentCustomerId)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	var DestinationResponses []response.DestinationResponse
	if len(data) > 0 {
		for _, value := range data {
			for _, routeOutlet := range value.Route.Destinations {
				outlet := response.DestinationResponse{
					DestinationID:      routeOutlet.DestinationID,
					DestinationCode:    routeOutlet.DestinationCode,
					DestinationName:    routeOutlet.DestinationName,
					Longitude:          routeOutlet.Longitude,
					Latitude:           routeOutlet.Latitude,
					DestinationStatus:  routeOutlet.DestinationStatus,
					DestinationAddress: routeOutlet.DestinationAddress,
				}
				DestinationResponses = append(DestinationResponses, outlet)
			}
		}
	}

	response := response.RouteDetailResponse{}
	if len(data) > 0 {
		response.ID = data[0].Route.ID
		response.RouteCode = data[0].Route.RouteCode
		response.RouteName = data[0].Route.RouteName
	}
	response.Outlets = DestinationResponses
	return response

}

func (service *RoutePopServiceImpl) GetAllVisitDayMap(ctx context.Context, dataFilter entity.VisitDayMapQueryFilter, currentCustomerId string) (response []entity.VisitDayMapResponse) {
	result := service.RoutePopPermanentRepository.GetAllVisitDayMap(ctx, dataFilter, currentCustomerId)

	for _, row := range result {
		var res entity.VisitDayMapResponse
		helper.Automapper(row, &res)
		res.WeekNumber = row.Week
		response = append(response, res)
	}

	return response
}

func (service *RoutePopServiceImpl) SaveOutletToRoute(ctx context.Context, request request.AddOutletToRouteRequest, currentCustomerId string, customerCode string) {
	// Validasi input
	if err := service.Validate.Struct(request); err != nil {
		helper.ErrorPanic(err)
	}

	localTime := time.Now()

	// Format waktu menjadi "2006-01-02"
	formattedTime := localTime.Format("2006-01-02")

	// Parse kembali waktu yang diformat menjadi time.Time
	parsedTime, err := time.Parse("2006-01-02", formattedTime)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return
	}

	// Mendapatkan detail waktu
	year, week := localTime.ISOWeek()
	dayString := localTime.Weekday().String()

	for _, value := range request.Data {
		// Cari route berdasarkan RouteCode
		route, err := service.RouteRepository.FindByRouteCode(ctx, value.RouteCode)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		// Konversi PjpCode menjadi integer
		resPjpCode, err := strconv.Atoi(value.PjpCode)
		if err != nil {
			log.Printf("Invalid PjpCode: %v", err)
			panic(err)
		}

		// Siapkan model RoutePopDaily
		routePopDaily := model.RoutePopDaily{
			RouteCode:   &value.RouteCode,
			Week:        week,
			Day:         dayString,
			Date:        parsedTime,
			ParentRoute: &value.RouteCode,
			PjpCode:     &resPjpCode,
			PjpID:       &value.PjpID,
			Year:        year,
			CustID:      currentCustomerId,
			Status:      "unplanned",
		}

		log.Printf("Inserting or updating RoutePopDaily: %+v", routePopDaily)

		// Update atau create RoutePopDaily
		service.RoutePopDailyRepository.UpdateOrCreateDaily(ctx, routePopDaily)

		// Proses setiap outlet
		for _, outlet := range value.Outlets {
			routeOutlet := model.DestinationAdditional{
				RouteCode:          value.RouteCode,
				DestinationID:      outlet.DestinationID,
				DestinationCode:    outlet.DestinationCode,
				DestinationName:    outlet.DestinationName,
				Longitude:          outlet.Longitude,
				Latitude:           outlet.Latitude,
				DestinationStatus:  strconv.Itoa(outlet.DestinationStatus),
				DestinationAddress: outlet.DestinationAddress,
				AvgSalesWeek:       outlet.AvgSalesWeek,
				OldRouteCode:       value.RouteCode,
				Status:             "Approved",
				PjpID:              &value.PjpID,
				PjpCode:            &resPjpCode,
				OldPjpID:           &value.PjpID,
				OldPjpCode:         &resPjpCode,
				RouteName:          route.RouteName,
				CustID:             currentCustomerId,
				OldRouteName:       route.RouteName,
				Date:               parsedTime,
				// IsPlanned:          false,
			}

			log.Printf("Inserting RouteOutletAdditional: %+v", routeOutlet)

			// Create RouteOutletAdditional
			service.RouteOutletRepository.CreateAdditionalRoute(ctx, routeOutlet)

			outletVisitList := model.OutletVisitList{
				Year:            year,
				Week:            week,
				Date:            parsedTime,
				Day:             dayString,
				RouteCode:       intToPtr(value.RouteCode),
				DestinationID:   outlet.DestinationID,
				DestinationCode: outlet.DestinationCode,
				PjpID:           &value.PjpID,
				PjpCode:         &resPjpCode,
				IsPlanned:       false,
			}
			service.OutletVisitListRepo.Create(ctx, outletVisitList)
		}
	}
}

func (service *RoutePopServiceImpl) CancelOutletToRoute(ctx context.Context, request request.CancelAddOutletToRouteRequest) {
	// Validasi request
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	// Mulai transaction
	tx, err := service.RouteOutletRepository.BeginTx(ctx)
	if err != nil {
		helper.ErrorPanic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// Looping untuk setiap data
	for _, value := range request.Data {
		// Inisialisasi model RouteOutletAdditional
		route := model.DestinationAdditional{
			RouteCode:       value.RouteCode,
			DestinationCode: value.DestinationCode,
			PjpID:           &value.PjpID,
		}

		// Hapus RouteOutletAdditional dengan transaction
		err = service.RouteOutletRepository.MobileCancelAddOutletToRoute(ctx, route, tx)
		if err != nil {
			tx.Rollback() // Rollback jika error
			helper.ErrorPanic(err)
		}

		// Inisialisasi model OutletVisitList
		data := model.OutletVisitList{
			RouteCode:       &value.RouteCode,
			DestinationCode: value.DestinationCode,
			PjpID:           &value.PjpID,
		}

		// Hapus OutletVisitList dengan transaction
		err = service.OutletVisitListRepo.MobileCancelAddOutletToRoute(ctx, data, tx)
		if err != nil {
			tx.Rollback() // Rollback jika error
			helper.ErrorPanic(err)
		}
	}

	// Commit transaction jika semua operasi sukses
	err = tx.Commit().Error
	if err != nil {
		helper.ErrorPanic(err) // Handle error saat commit
	}
}

func intToPtr(i int) *int {
	return &i
}

func getWeekRangeFromNumber(week int) (string, string, error) {
	// Ambil tahun sekarang
	currentYear := time.Now().Year()

	// Tentukan hari pertama tahun ini
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)

	// Cari hari pertama minggu pertama (Senin pertama di tahun ini)
	weekday := startOfYear.Weekday()
	offset := (int(weekday) - 1 + 7) % 7 // Jarak ke Senin
	firstMonday := startOfYear.AddDate(0, 0, -offset)

	// Hitung awal minggu ke-n
	startOfWeek := firstMonday.AddDate(0, 0, (week-1)*7)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	// Format hasil ke string
	layout := "2006-01-02"
	start := startOfWeek.Format(layout)
	end := endOfWeek.Format(layout)

	return start, end, nil
}
