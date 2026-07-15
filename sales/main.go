package main

import (
	"sales/adapter"
	"sales/controller"
	"sales/pkg/config"
	"sales/pkg/config/env"
	"sales/pkg/middleware"
	"sales/pkg/server"
	"sales/pkg/validation"
	"sales/repository"
	"sales/service"

	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	envCfg := env.NewCfgEnv()
	validatorPkg := validation.NewValiditor()

	// Validate required OBS config early so missing keys fail clearly at startup.
	if err := env.ValidateRequired(envCfg,
		"OBS_HUAWEI_AK",
		"OBS_HUAWEI_SK",
		"OBS_HUAWEI_ENDPOINT",
		"OBS_HUAWEI_BUCKET",
	); err != nil {
		panic(err)
	}

	postgreDB := config.PostgreSQLConnection(envCfg)

	transactionDB := repository.NewDbtransactionRepo(postgreDB)
	roRepository := repository.NewRoRepo(postgreDB)
	soRepository := repository.NewSoRepo(postgreDB)
	orderRepository := repository.NewOrderRepo(postgreDB)
	invoiceRepository := repository.NewInvoiceRepo(postgreDB)
	returnRepository := repository.NewReturnRepo(postgreDB)
	gamificationRepository := repository.NewGamificationRepo(postgreDB)
	consignmentRepository := repository.NewConsignmentRepo(postgreDB)
	tlsRepository := repository.NewTlsRepo(postgreDB)
	promotionRepository := repository.NewPromotionRepo(postgreDB)
	promotionV2Repository := repository.NewPromotionV2Repo(postgreDB)
	promoTemplateRepository := repository.NewPromoTemplateRepo(postgreDB)
	discountRepository := repository.NewDiscountRepo(postgreDB)
	validateOrderRepository := repository.NewValidateOrderRepo(postgreDB)
	stockRepository := repository.NewStockRepo(postgreDB)
	hierarchyApprovalRepository := repository.NewHierarchyApprovalRepo(postgreDB)
	orderApprovalRequest := repository.NewOrderApprovalRequestRepo(postgreDB)
	orderApprovalRepository := repository.NewOrderApprovalRepo(postgreDB)
	reportRepository := repository.NewReportRepo(postgreDB)
	openAPIRepository := repository.NewOpenAPIRepo(postgreDB)

	// setup adapter
	obsAdapter, err := adapter.InitObsAdapter(envCfg.Get("OBS_HUAWEI_AK"), envCfg.Get("OBS_HUAWEI_SK"), envCfg.Get("OBS_HUAWEI_ENDPOINT"), envCfg.Get("OBS_HUAWEI_BUCKET"))
	if err != nil {
		panic(err)
	}

	roService := service.NewRoService(roRepository, transactionDB)
	orderService := service.NewOrderService(orderRepository, validateOrderRepository, promotionRepository, promotionV2Repository, discountRepository, stockRepository, transactionDB)
	invoiceService := service.NewInvoiceService(invoiceRepository, stockRepository, transactionDB)
	soService := service.NewSoService(soRepository, reportRepository, transactionDB)
	returnService := service.NewReturnService(returnRepository, orderRepository, promotionRepository, discountRepository, transactionDB)
	gamificationService := service.NewGamificationService(gamificationRepository, transactionDB)
	consignmentService := service.NewConsignmentService(consignmentRepository, transactionDB)
	tlsService := service.NewTlsService(tlsRepository, transactionDB)
	filesService := service.NewFilesService(envCfg, obsAdapter)
	promotionService := service.NewPromotionService(promotionRepository, promotionV2Repository, transactionDB)
	promoTemplateService := service.NewPromoTemplateService(promoTemplateRepository, transactionDB)
	discountService := service.NewDiscountService(discountRepository, transactionDB)
	validateOrderService := service.NewValidateOrderService(validateOrderRepository, transactionDB)
	hierarchyApprovalService := service.NewHierarchyApprovalService(hierarchyApprovalRepository, orderRepository, orderApprovalRequest, transactionDB)
	orderApprovalService := service.NewOrderApprovalService(orderApprovalRepository, transactionDB)
	reportService := service.NewReportService(envCfg, reportRepository, transactionDB, obsAdapter)
	openAPIService := service.NewOpenAPIService(openAPIRepository)

	roController := controller.NewRoController(roService, validatorPkg)
	orderController := controller.NewOrderController(orderService, validateOrderService, validatorPkg)
	invoiceController := controller.NewInvoiceController(invoiceService, validatorPkg)
	soController := controller.NewSoController(soService, validatorPkg)
	returnController := controller.NewReturnController(returnService, promotionService, discountService, validatorPkg)
	gamificationController := controller.NewGamificationController(gamificationService, validatorPkg)
	consignmentController := controller.NewConsignmentController(consignmentService, validatorPkg)
	tlsController := controller.NewTlsController(tlsService, validatorPkg)
	filesController := controller.NewFilesController(filesService, validatorPkg)
	promotionController := controller.NewPromotionController(promotionService, validatorPkg)
	openAPIPromotionController := controller.NewOpenAPIPromotionController(promotionController, openAPIService, validatorPkg)
	promoTemplateController := controller.NewPromoTemplateController(promoTemplateService, validatorPkg)
	discountController := controller.NewDiscountController(discountService, validatorPkg)
	validateOrderController := controller.NewValidateOrderController(validateOrderService, validatorPkg)
	hierarchyApprovalController := controller.NewHierarchyApprovalController(hierarchyApprovalService, validatorPkg)
	orderApprovalController := controller.NewOrderApprovalController(orderApprovalService, validatorPkg)
	reportController := controller.NewReportController(reportService, validatorPkg)

	// define Fiber Framework config
	fiberCfg := config.NewFiberConfig(envCfg)
	app := fiber.New(fiberCfg)

	// middleware
	middleware.AppMiddleware(app, envCfg) // Register Fiber's middleware for app.

	// This route path for test service is running "/"
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("It works")
	})
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Master Service Metrics Page"}))

	roController.Route(app)
	orderController.Route(app)
	invoiceController.Route(app)
	soController.Route(app)
	returnController.Route(app)
	gamificationController.Route(app)
	consignmentController.Route(app)
	tlsController.Route(app)
	filesController.Route(app)
	promotionController.Route(app)
	openAPIPromotionController.Route(app)
	promoTemplateController.Route(app)
	discountController.Route(app)
	validateOrderController.Route(app)
	hierarchyApprovalController.Route(app)
	orderApprovalController.Route(app)
	reportController.Route(app)

	scheduler := reportController.Cron()

	_, err = scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(00, 01, 0),
			),
		),
		gocron.NewTask(func() {
			if err := promotionService.CloseExpiredPromotions(); err != nil {
				fmt.Println("error closing expired promotions:", err)
			}
		}),
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("error scheduling CloseExpiredPromotions job: %v", err))
	}

	// start fiber server
	server.FiberServerWithGracefulShutdown(app, scheduler)
}
