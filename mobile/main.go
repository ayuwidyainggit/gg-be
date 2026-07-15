package main

import (
	"log"
	"mobile/adapter"
	"mobile/controller"
	"mobile/pkg/config"
	"mobile/pkg/config/env"
	"mobile/pkg/validation"
	"mobile/repository"
	"mobile/service"

	"mobile/pkg/middleware"
	"mobile/pkg/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var Redis config.RedisConfig
	// load config env
	envCfg := env.NewCfgEnv()
	validatorPkg := validation.NewValidator()
	// emailPkg := smtp.NE

	postgreDB := config.PostgreSQLConnection(envCfg)

	postgreDB2, err := config.ConnToDb(envCfg)
	if err != nil {
		log.Println("err:", err.Error())
		panic(err)
	}

	redisClient, errRedisCon := Redis.SetConfig().RedisInstance()
	if errRedisCon != nil {
		panic(errRedisCon)
	}

	// Setup Repository
	transactionDB := repository.NewDbtransactionRepo(postgreDB)
	mCustomerRepository := repository.NewMCustomerRepository(postgreDB)
	mEmployeeRepository := repository.NewMEmployeeRepository(postgreDB)
	promotionRepository := repository.NewPromotionRepo(postgreDB)
	mSalesmanRepository := repository.NewMSalesmanRepository(postgreDB)
	mProductDistRepository := repository.NewMProductDistRepository(postgreDB)
	productRepository := repository.NewProductRepository(postgreDB)
	orderRepository := repository.NewOrderRepo(postgreDB)
	collectionRepository := repository.NewCollectionRepo(postgreDB)
	userRepository := repository.NewUserRepository(postgreDB)
	passwordResetRepository := repository.NewPasswordResetRequestRepository(postgreDB)
	attendanceRepository := repository.NewAttendanceRepository(postgreDB)
	leaveRequestRepository := repository.NewLeaveRequestRepository(postgreDB)
	invoicesRepository := repository.NewInvoicesRepository(postgreDB)
	takingOrderRepository := repository.NewMTakingOrderRepository(postgreDB)
	visitRepository := repository.NewVisitsRepository(postgreDB)
	returnRepository := repository.NewReturnRepo(postgreDB)
	pickupReasonRepository := repository.NewPickupReasonRepository(postgreDB2)
	stockRepository := repository.NewStockRepository(postgreDB)
	orderCanvasRepository := repository.NewOrderCanvasRepo(postgreDB)
	orderHistoryRepository := repository.NewOrderHistoryRepo(postgreDB)
	salesRepository := repository.NewSalesRepository(postgreDB)
	activitiesRepository := repository.NewActivitiesRepository(postgreDB)
	// eventsRepository := repository.NewEventsRepository(postgreDB)
	// announcementsRepository := repository.NewAnnouncementsRepository(postgreDB)
	// leaderboardsRepository := repository.NewLeaderboardsRepository(postgreDB)
	// returnReasonsRepository := repository.NewReturnsRepository(postgreDB)

	_ = repository.NewCache(redisClient)
	// mConfigRepository := repository.NewMConfigRepository(postgreDB)
	mOutletRepository := repository.NewMOutletRepository(postgreDB2)
	pjpDistributorRepository := repository.NewPjpDistributorRepository(postgreDB)
	pjpPrincipalRepository := repository.NewPjpPrincipalRepository(postgreDB)
	outletListRepository := repository.NewOutletListRepository(postgreDB2)
	mEmployee2Repository := repository.NewEmployeeRepository(postgreDB2)
	mEmpGroupRepository := repository.NewEmpGroupRepository(postgreDB2)
	discountRepository := repository.NewDiscountRepo(postgreDB)

	validateOrderRepository := repository.NewValidateOrderRepo(postgreDB)
	expenseRepository := repository.NewExpenseRepository(postgreDB)
	printDailyReportRepository := repository.NewPrintDailyReportRepository(postgreDB)
	weekRepository := repository.NewWeekRepository(postgreDB)
	workingDayCalendarRepository := repository.NewWorkingDayCalendarRepository(postgreDB)
	distributorRepository := repository.NewDistributorRepository(postgreDB)
	mWarahouseRepository := repository.NewMWarehouseRepository(postgreDB)
	paymentRepository := repository.NewPaymentRepository(postgreDB)
	outletBankRepository := repository.NewOutletBankRepository(postgreDB)
	surveyRepository := repository.NewSurveyRepository(postgreDB)
	bankRepository := repository.NewBankRepository(postgreDB)
	regionRepository := repository.NewRegionRepository(postgreDB)
	areaRepository := repository.NewAreaRepository(postgreDB)
	locationRepo := repository.NewUserLocationRepository(postgreDB)

	// setup adapter - use local storage if USE_LOCAL_STORAGE is true, otherwise use OBS
	var obsAdapter adapter.ObsAdapter
	useLocalStorage := envCfg.Get("USE_LOCAL_STORAGE")
	if useLocalStorage == "true" || useLocalStorage == "1" {
		log.Println("Using local storage for file uploads")
		localBasePath := envCfg.Get("LOCAL_STORAGE_PATH")
		if localBasePath == "" {
			localBasePath = "./uploads"
		}
		localBaseURL := envCfg.Get("LOCAL_STORAGE_URL")
		if localBaseURL == "" {
			localBaseURL = "http://localhost:9008/uploads"
		}
		localPublicPath := envCfg.Get("LOCAL_STORAGE_PUBLIC_PATH")
		if localPublicPath == "" {
			localPublicPath = "./public/uploads"
		}

		localAdapter, err := adapter.InitLocalStorageAdapter(localBasePath, localBaseURL, localPublicPath)
		if err != nil {
			log.Printf("Failed to initialize local storage adapter: %v", err)
			panic(err)
		}
		obsAdapter = localAdapter
	} else {
		log.Println("Using OBS for file uploads")
		obsAK := envCfg.Get("OBS_HUAWEI_AK")
		obsSK := envCfg.Get("OBS_HUAWEI_SK")
		obsEndpoint := envCfg.Get("OBS_HUAWEI_ENDPOINT")
		obsBucket := envCfg.Get("OBS_HUAWEI_BUCKET")

		log.Printf("OBS Config - AK: %s, SK: %s, Endpoint: %s, Bucket: %s",
			func() string {
				if obsAK != "" {
					return "***"
				} else {
					return "EMPTY"
				}
			}(),
			func() string {
				if obsSK != "" {
					return "***"
				} else {
					return "EMPTY"
				}
			}(),
			obsEndpoint,
			obsBucket)

		obsAdapterImpl, err := adapter.InitObsAdapter(obsAK, obsSK, obsEndpoint, obsBucket)
		if err != nil {
			log.Printf("Failed to initialize OBS adapter: %v", err)
			panic(err)
		}
		obsAdapter = obsAdapterImpl
	}

	// Setup Service
	userService := service.NewUserService(
		envCfg,
		mCustomerRepository,
		mEmployeeRepository,
		mSalesmanRepository,
		userRepository,
		locationRepo,
		passwordResetRepository,
		transactionDB)
	salesService := service.NewSalesService(envCfg, transactionDB, salesRepository)
	activitiesService := service.NewActivitiesService(envCfg, transactionDB, activitiesRepository)
	atttendanceService := service.NewAttendanceService(envCfg, transactionDB, attendanceRepository, mEmployeeRepository)
	leaveService := service.NewLeaveService(leaveRequestRepository, attendanceRepository, obsAdapter)
	eventsService := service.NewEventsService(envCfg, transactionDB)
	announcementsService := service.NewAnnouncementsService(envCfg, transactionDB)
	leaderboardsService := service.NewLeaderboardsService(envCfg, transactionDB)
	visitsService := service.NewVisitsService(envCfg, transactionDB, mEmployeeRepository, visitRepository, obsAdapter, invoicesRepository, pjpDistributorRepository, pjpPrincipalRepository)
	invoiceService := service.NewInvoicesService(envCfg, invoicesRepository, paymentRepository, outletBankRepository, transactionDB)
	returnReasonsService := service.NewReturnService(envCfg, returnRepository, transactionDB)
	promotionService := service.NewPromotionService(promotionRepository, envCfg, transactionDB)

	productService := service.NewProductService(
		envCfg,
		productRepository,
		mProductDistRepository,
		transactionDB)

	orderService := service.NewOrderService(
		envCfg,
		orderRepository,
		discountRepository,
		pjpPrincipalRepository,
		pjpDistributorRepository,
		transactionDB)

	collectionService := service.NewCollectionService(
		envCfg,
		collectionRepository,
		transactionDB,
		orderRepository)

	// orderService := service.NewOrderService(orderRepository, transactionDB)
	orderHistoryService := service.NewOrderHistoryService(envCfg, orderHistoryRepository, transactionDB)

	filesService := service.NewFilesService(envCfg, obsAdapter)
	takingOrderService := service.NewTakingOrderService(envCfg, takingOrderRepository)
	pickupReasonService := service.NewPickupReasonService(pickupReasonRepository)
	// configService := service.NewMConfigService(mConfigRepository, transactionDB)
	mOutletService := service.NewMOutletService(mOutletRepository, pjpDistributorRepository, pjpPrincipalRepository, invoicesRepository)
	outletListService := service.NewOutletListService(outletListRepository)
	mEmployeeService := service.NewEmployeeService(mEmployee2Repository)
	mEmployeeGroupService := service.NewEmpGroupService(mEmpGroupRepository)
	discountService := service.NewDiscountService(discountRepository, transactionDB)

	validateOrderService := service.NewValidateOrderService(validateOrderRepository, transactionDB)
	stockService := service.NewStockService(stockRepository, transactionDB)
	orderCanvasService := service.NewOrderCanvasService(
		envCfg,
		orderCanvasRepository,
		discountRepository,
		transactionDB)

	expenseService := service.NewExpenseService(
		envCfg,
		expenseRepository,
		transactionDB,
		obsAdapter,
		mOutletRepository,
	)

	pjpService := service.NewPjpService(mSalesmanRepository, mWarahouseRepository, pjpDistributorRepository, pjpPrincipalRepository, transactionDB)

	printDailyReportService := service.NewPrintDailyReportService(
		envCfg,
		transactionDB,
		printDailyReportRepository,
		userRepository)
	weekService := service.NewWeekService(weekRepository, pjpDistributorRepository, pjpPrincipalRepository)
	workingDayCalendarService := service.NewWorkingDayCalendarService(workingDayCalendarRepository)
	distributorService := service.NewDistributorService(distributorRepository)
	pjpPrincipalService := service.NewPjpPrincipalService(pjpPrincipalRepository, mCustomerRepository, transactionDB)
	surveyService := service.NewSurveyService(surveyRepository, distributorRepository, mOutletRepository)
	collectionPayService := service.NewCollectionPayService(paymentRepository, bankRepository, collectionRepository, transactionDB)
	regionService := service.NewRegionService(regionRepository)
	areaService := service.NewAreaService(areaRepository)
	xtraCallService := service.NewExtraCallService(transactionDB, pjpDistributorRepository, pjpPrincipalRepository, mOutletRepository)

	// Setup Controller
	userController := controller.NewUserController(userService, validatorPkg)
	salesController := controller.NewSalesController(salesService, validatorPkg)
	activitiesController := controller.NewActivitiesController(activitiesService, validatorPkg)
	attendanceController := controller.NewAttendanceController(atttendanceService, validatorPkg)
	leaveController := controller.NewLeaveController(leaveService, validatorPkg)
	eventsController := controller.NewEventsController(eventsService, validatorPkg)
	announcementsController := controller.NewAnnouncementsController(announcementsService, validatorPkg)
	leaderboardsController := controller.NewLeaderboardsController(leaderboardsService, validatorPkg)
	visitsController := controller.NewVisitsController(visitsService, validatorPkg)
	productController := controller.NewProductController(productService, validatorPkg)
	orderController := controller.NewOrderController(orderService, discountService, validatorPkg)
	collectionController := controller.NewCollectionController(collectionService, validatorPkg)
	invoiceController := controller.NewInvoicesController(invoiceService, validatorPkg)
	returnReasonsController := controller.NewReturnController(returnReasonsService, validatorPkg)
	filesController := controller.NewFilesController(filesService, validatorPkg)
	promotionController := controller.NewPromotionController(promotionService, validatorPkg)
	takingOrderController := controller.NewTakingOrderController(takingOrderService, validatorPkg)
	pickupReasonController := controller.NewPickupReasonController(pickupReasonService, validatorPkg)
	// configController := controller.NewMConfigController(configService, validatorPkg)
	mOutletController := controller.NewMOutletController(mOutletService, validatorPkg)
	outletListController := controller.NewOutletListController(outletListService, validatorPkg)
	mEmployeeController := controller.NewEmployeeController(mEmployeeService, validatorPkg)
	mEmployeeGroupController := controller.NewEmpGroupController(mEmployeeGroupService, validatorPkg)
	discountController := controller.NewDiscountController(discountService, validatorPkg)

	validateOrderController := controller.NewValidateOrderController(validateOrderService, validatorPkg)
	stockController := controller.NewStockController(stockService, validatorPkg)
	orderCanvasController := controller.NewOrderCanvasController(orderCanvasService, discountService, validatorPkg)
	orderHistoryController := controller.NewOrderHistoryController(orderHistoryService, discountService, validatorPkg)
	expenseController := controller.NewExpenseController(expenseService, validatorPkg)
	printDailyReportController := controller.NewPrintDailyReportController(printDailyReportService, validatorPkg)
	weekController := controller.NewWeekController(weekService, validatorPkg)
	workingDayCalendarController := controller.NewWorkingDayCalendarController(workingDayCalendarService, validatorPkg)
	pjpController := controller.NewPjpController(pjpService, validatorPkg)
	distributorController := controller.NewDistributorController(distributorService, validatorPkg)
	pjpDistributorController := controller.NewPjpDistributorController(pjpService, validatorPkg)
	pjpPrincipalController := controller.NewPjpPrincipalController(pjpPrincipalService, validatorPkg)
	surveyController := controller.NewSurveyController(surveyService, validatorPkg)
	collectionPayController := controller.NewCollectionPayController(collectionPayService, validatorPkg)
	regionController := controller.NewRegionController(regionService, validatorPkg)
	areaController := controller.NewAreaController(areaService, validatorPkg)
	xtraCallController := controller.NewExtraCallController(xtraCallService, validatorPkg)

	// define Fiber Framework config
	fiberCfg := config.NewFiberConfig(envCfg)
	app := fiber.New(fiberCfg)

	// middleware
	middleware.AppMiddleware(app) // Register Fiber's middleware for app.

	// Serve static files for local storage if enabled
	if useLocalStorage == "true" || useLocalStorage == "1" {
		localPublicPath := envCfg.Get("LOCAL_STORAGE_PUBLIC_PATH")
		if localPublicPath == "" {
			localPublicPath = "./public/uploads"
		}
		app.Static("/uploads", localPublicPath)
		log.Printf("Serving static files from: %s at /uploads", localPublicPath)
	}

	// This route path for test service is running "/"
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("It works")
	})

	// route.PublicRoutes(app)
	userController.Route(app)
	salesController.Route(app)
	activitiesController.Route(app)
	attendanceController.Route(app)
	leaveController.Route(app)
	eventsController.Route(app)
	announcementsController.Route(app)
	leaderboardsController.Route(app)
	visitsController.Route(app)
	productController.Route(app)
	orderController.Route(app)
	collectionController.Route(app)
	invoiceController.Route(app)
	returnReasonsController.Route(app)
	filesController.Route(app)
	promotionController.Route(app)
	takingOrderController.Route(app)
	pickupReasonController.Route(app)
	// configController.Route(app)
	// itemController.Route(app)
	mOutletController.Route(app)
	outletListController.Route(app)
	mEmployeeController.Route(app)
	mEmployeeGroupController.Route(app)
	discountController.Route(app)

	validateOrderController.Route(app)
	stockController.Route(app)
	orderCanvasController.Route(app)
	orderHistoryController.Route(app)
	expenseController.Route(app)
	printDailyReportController.Route(app)
	weekController.Route(app)
	workingDayCalendarController.Route(app)
	pjpController.Route(app)
	distributorController.Route(app)
	pjpDistributorController.Route(app)
	pjpPrincipalController.Route(app)
	surveyController.Route(app)
	collectionPayController.Route(app)
	regionController.Route(app)
	areaController.Route(app)
	xtraCallController.Route(app)
	// start fiber server

	server.FiberServerWithGracefulShutdown(app)
}
