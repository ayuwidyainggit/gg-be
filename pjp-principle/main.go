package main

import (
	"log"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/controller"
	"scyllax-pjp/controller/daily"
	pjpCtrl "scyllax-pjp/controller/pjp"
	pjpauto "scyllax-pjp/controller/pjp_auto"
	pjpenhance "scyllax-pjp/controller/pjp_enhance"
	"scyllax-pjp/controller/route"
	thirdparty "scyllax-pjp/controller/third_party"
	"scyllax-pjp/controller/visit"
	"scyllax-pjp/docs"
	"scyllax-pjp/helper"
	"scyllax-pjp/repository"
	routeOutletRepo "scyllax-pjp/repository/destination"
	routeOutletHistoryRepo "scyllax-pjp/repository/destination_history"
	distributordms "scyllax-pjp/repository/distributor_dms"
	outletdms "scyllax-pjp/repository/outlet_dms"
	outletVisitRepo "scyllax-pjp/repository/outlet_visit"
	pjpRepo "scyllax-pjp/repository/pjp"
	routeRepo "scyllax-pjp/repository/route"
	routePopDailyRepo "scyllax-pjp/repository/route_pop_daily"
	routePopPermRepo "scyllax-pjp/repository/route_pop_permanent"
	"scyllax-pjp/router"
	"scyllax-pjp/service"
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

// @title 	ScyllaX-PJP-Principle API
// @version	1.0
// @description ScyllaX-PJP API-Principle in Go using Gin framework

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
	pjpRepository := pjpRepo.NewPjpRepository()
	routeOutletRepositoryNew := routeOutletRepo.NewDestinationRepository()
	routeOutletHistoryRepository := routeOutletHistoryRepo.NewDestinationHistoryRepository()
	routeRepositoryNew := routeRepo.NewRouteRepository()
	routePopPermanentRepositoryNew := routePopPermRepo.NewRoutePopPermanentRepository()
	routePopDailyRepositoryNew := routePopDailyRepo.NewRoutePopDailyRepositoryImpl()

	routeRepository := repository.NewRouteRepositoryImpl(db)
	routeOutletRepository := repository.NewRouteOutletRepositoryImpl(db)
	routePopPermanentRepository := repository.NewRoutePopPermanentRepositoryImpl(db)
	routePopDailyRepository := repository.NewRoutePopDailyRepositoryImpl(db)
	// outletRepository := repository.NewOutletRepository(db)
	outletVisitList := repository.NewOutletVisitRepoImpl(db)
	outletVisitRepo := outletVisitRepo.NewOutletVisitRepository()
	outletDmsRepo := outletdms.NewOutletDmsRepository()
	distributorDmsRepo := distributordms.NewDistributorDmsRepository()

	//Init Service
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
	routePopService := routepop.NewRoutePopService(pjpRepository, routeOutletHistoryRepository, routePopDailyRepository, routeRepository, routeOutletRepository, validate, db)
	thirdPartyService := thirdPartyServ.NewThirdPartyService(pjpRepository, outletDmsRepo, distributorDmsRepo, db)
	visitService := service.NewVisitServiceImpl(routeOutletRepository, outletVisitList, pjpRepository, validate)
	visitServiceNew := visitServ.NewVisitService(pjpRepository, outletVisitRepo, validate, db)

	//Init controller
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
	)

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	error := server.ListenAndServe()
	helper.ErrorPanic(error)
}
