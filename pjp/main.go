package main

import (
	"log"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/controller"
	attendanceCtrl "scyllax-pjp/controller/attendance"
	"scyllax-pjp/controller/daily"
	"scyllax-pjp/controller/geotaging"
	livemonitoringctrl "scyllax-pjp/controller/live_monitoring"
	pjpCtrl "scyllax-pjp/controller/pjp"
	pjpauto "scyllax-pjp/controller/pjp_auto"
	pjpenhance "scyllax-pjp/controller/pjp_enhance"
	"scyllax-pjp/controller/route"
	thirdparty "scyllax-pjp/controller/third_party"
	"scyllax-pjp/controller/visit"
	"scyllax-pjp/docs"
	"scyllax-pjp/helper"
	"scyllax-pjp/repository"
	arrivalReportRepo "scyllax-pjp/repository/arrival_report"
	attendanceRepo "scyllax-pjp/repository/attendance"
	livemonitoringrepo "scyllax-pjp/repository/live_monitoring"
	outletVisitRepo "scyllax-pjp/repository/outlet_visit"
	"scyllax-pjp/repository/outlet_visit_principle"
	pjpRepo "scyllax-pjp/repository/pjp"
	"scyllax-pjp/repository/pjp_principle"
	routeRepo "scyllax-pjp/repository/route"
	routeOutletRepo "scyllax-pjp/repository/route_outlet"
	routeOutletHistoryRepo "scyllax-pjp/repository/route_outlet_history"
	routePopDailyRepo "scyllax-pjp/repository/route_pop_daily"
	routePopPermRepo "scyllax-pjp/repository/route_pop_permanent"
	"scyllax-pjp/router"
	"scyllax-pjp/service"
	attendanceServ "scyllax-pjp/service/attendance"
	geotagingServ "scyllax-pjp/service/geotaging"
	livemonitoringserv "scyllax-pjp/service/live_monitoring"
	pjpServ "scyllax-pjp/service/pjp"
	pjpAutoServ "scyllax-pjp/service/pjp_auto"
	pjpEnhanceServ "scyllax-pjp/service/pjp_enhance"
	routeServ "scyllax-pjp/service/route"
	routepop "scyllax-pjp/service/route_pop"
	thirdPartyServ "scyllax-pjp/service/third_party"
	visitServ "scyllax-pjp/service/visit"
	"scyllax-pjp/utils"
	"time"
)

// @title 	ScyllaX-PJP API
// @version	1.0
// @description ScyllaX-PJP API in Go using Gin framework

// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	//Database
	db := config.ConnectionDB(&loadConfig)

	// environment swagger
	if loadConfig.Environment != "dev" {
		// docs.SwaggerInfo.Host = "103.28.219.73:5001"
		docs.SwaggerInfo.Host = loadConfig.SwaggerHost
		docs.SwaggerInfo.BasePath = "/scylla-pjp/api/v1"
		// docs.SwaggerInfo.BasePath = loadConfig.SwaggerUrl
	} else {
		docs.SwaggerInfo.Host = "localhost:8888"
		docs.SwaggerInfo.BasePath = "/api/v1"
	}

	//Validation
	validate := utils.InitializeValidator(db)

	//Init Repository
	attendanceRepository := attendanceRepo.NewAttendanceRepository()
	pjpRepository := pjpRepo.NewPjpRepository()
	pjpPrincipleRepo := pjp_principle.NewPjpPrincipleRepository()
	routeOutletRepositoryNew := routeOutletRepo.NewRouteOutletRepository()
	routeOutletHistoryRepository := routeOutletHistoryRepo.NewRouteOutletHistoryRepository()
	routeRepositoryNew := routeRepo.NewRouteRepository()
	routePopPermanentRepositoryNew := routePopPermRepo.NewRoutePopPermanentRepository()
	routePopDailyRepositoryNew := routePopDailyRepo.NewRoutePopDailyRepositoryImpl()
	liveMonitoringRepository := livemonitoringrepo.NewLiveMonitoringRepository()

	routeRepository := repository.NewRouteRepositoryImpl(db)
	routeOutletRepository := repository.NewRouteOutletRepositoryImpl(db)
	routePopPermanentRepository := repository.NewRoutePopPermanentRepositoryImpl(db)
	routePopDailyRepository := repository.NewRoutePopDailyRepositoryImpl(db)
	// outletRepository := repository.NewOutletRepository(db)
	outletVisitList := repository.NewOutletVisitRepoImpl(db)
	outletVisitRepo := outletVisitRepo.NewOutletVisitRepository()
	outletVisitPrincipleRepo := outlet_visit_principle.NewOutletVisitPrincipleRepository()
	outletCrRepository := repository.NewOutletCrRepository(db)
	arrivalReportRepository := arrivalReportRepo.NewArrivalReportRepository()
	customerRepository := repository.NewCustomerRepositoryImpl(db)

	//Init Service
	attendanceService := attendanceServ.NewAttendanceService(attendanceRepository, db)
	pjpService := pjpServ.NewPjpService(pjpRepository, routeOutletRepositoryNew, routeRepositoryNew, validate, db)
	pjpEnhanceService := pjpEnhanceServ.NewPjpEnhanceService(pjpRepository, routeOutletRepositoryNew, routeOutletHistoryRepository, routeRepositoryNew, routePopPermanentRepositoryNew, validate, db)
	pjpAutoService := pjpAutoServ.NewPjpAutoService(pjpRepository, routeRepository, routeOutletRepository, validate, db)
	routeService := service.NewRouteServiceImpl(
		routeRepository,
		routeOutletRepository,
		routePopDailyRepository,
		routePopPermanentRepository,
		outletVisitList,
		pjpRepository,
		validate,
	)

	routeServiceNew := routeServ.NewRouteService(validate, pjpRepository, routeOutletRepositoryNew, routePopPermanentRepositoryNew, routePopDailyRepositoryNew, db)

	routePopDailyService := service.NewRoutePopServiceImpl(routePopPermanentRepository, routePopDailyRepository, pjpRepository, routeRepository, routeOutletRepository, outletVisitList, validate)
	routePopService := routepop.NewRoutePopService(routeOutletHistoryRepository, routePopDailyRepository, routeRepository, routeOutletRepository, validate, db)
	thirdPartyService := thirdPartyServ.NewThirdPartyService(pjpRepository, routeOutletRepositoryNew, routeOutletHistoryRepository, db)
	visitService := service.NewVisitServiceImpl(routeOutletRepository, outletVisitList, outletVisitPrincipleRepo, pjpRepository, outletCrRepository, arrivalReportRepository, customerRepository, validate, db)
	visitServiceNew := visitServ.NewVisitService(pjpRepository, pjpPrincipleRepo, outletVisitRepo, outletVisitPrincipleRepo, validate, db)
	geotagingService := geotagingServ.NewGeotagingService(validate, db, loadConfig)
	liveMonitoringService := livemonitoringserv.NewLiveMonitoringService(liveMonitoringRepository, validate, db)

	//Init controller
	attendanceController := attendanceCtrl.NewAttendanceController(attendanceService)
	pjpController := pjpCtrl.NewPjpController(pjpService)

	pjpEnhanceController := pjpenhance.NewPjpEnhanceController(pjpEnhanceService, routePopDailyService)
	pjpAutoController := pjpauto.NewPjpAutoController(pjpAutoService)
	routeMappingController := controller.NewRouteMappingController(routeService)
	routeController := route.NewRouteController(routeServiceNew)
	routePopDailyController := controller.NewVisitDayMapController(routePopDailyService)
	dailyRouteMapController := controller.NewDailyRouteMapController(routePopDailyService)
	dailyRouteMapControllerNew := daily.NewDailyRouteMapController(routePopService)
	thirdPartyController := thirdparty.NewThirdPartyController(thirdPartyService)
	visitController := controller.NewVisitController(visitService)
	visitControllerNew := visit.NewVisitController(visitServiceNew)
	geotagingController := geotaging.NewGeotagingController(geotagingService)
	liveMonitoringController := livemonitoringctrl.NewLiveMonitoringController(liveMonitoringService)

	//scheduler
	utils.CheckWeeklyJob(routePopPermanentRepository, routePopDailyRepository)

	//Router
	routes := router.NewRouter(
		pjpController,
		pjpEnhanceController,
		pjpAutoController,
		routeMappingController,
		routeController,
		routePopDailyController,
		thirdPartyController,
		dailyRouteMapController,
		dailyRouteMapControllerNew,
		visitController,
		visitControllerNew,
		geotagingController,
		attendanceController,
		liveMonitoringController,
	)

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	error := server.ListenAndServe()
	helper.ErrorPanic(error)
}
