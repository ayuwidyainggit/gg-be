package main

import (
	"inventory/adapter"
	"inventory/controller"
	"inventory/pkg/config"
	"inventory/pkg/config/env"
	"inventory/pkg/middleware"
	"inventory/pkg/server"
	"inventory/pkg/validation"
	"inventory/repository"
	"inventory/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	// load config env
	envCfg := env.NewCfgEnv()
	validatorPkg := validation.NewValidator()

	postgreDB := config.PostgreSQLConnection(envCfg)

	transactionDB := repository.NewDbtransactionRepo(postgreDB)
	grRepository := repository.NewGrRepo(postgreDB)
	grBranchRepository := repository.NewGrBranchRepo(postgreDB)
	arBranchRepository := repository.NewArBranchRepo(postgreDB)
	bpprRepository := repository.NewBpprRepo(postgreDB)
	whStockRepository := repository.NewWhStockRepo(postgreDB)
	smpIssRepository := repository.NewSmpIssRepo(postgreDB)
	whTrfRepository := repository.NewWhTrfRepo(postgreDB)
	itemStChRepository := repository.NewItemStChRepo(postgreDB)
	whAdjRepository := repository.NewWhAdjRepo(postgreDB)
	whSoRepository := repository.NewWhSoRepo(postgreDB)
	vanSoRepository := repository.NewVanSoRepo(postgreDB)
	gdsRepository := repository.NewGdsRepo(postgreDB)
	vanBsUlRepository := repository.NewVanBsUlRepo(postgreDB)
	vanUlRepository := repository.NewVanUlRepo(postgreDB)
	vanLoRepository := repository.NewVanLoRepo(postgreDB)
	stockRepository := repository.NewStockRepo(postgreDB)
	warehouseStockRepository := repository.NewWarehouseStockRepo(postgreDB)
	SupplierReturnRepository := repository.NewSupplierReturnRepo(postgreDB)
	StockReturnRepository := repository.NewStockReturnRepo(postgreDB)
	StockOpnameRepository := repository.NewStockOpnameRepo(postgreDB)
	OrderBookingRepository := repository.NewOrderBookingRepo(postgreDB)
	StockDisposalRepository := repository.NewStockDisposalRepo(postgreDB)
	ReplenishmentRepository := repository.NewReplenishmentRepo(postgreDB)
	ReportsRepository := repository.NewReportsRepo(postgreDB)

	// setup adapter
	obsAdapter, err := adapter.InitObsAdapter(envCfg.Get("OBS_HUAWEI_AK"), envCfg.Get("OBS_HUAWEI_SK"), envCfg.Get("OBS_HUAWEI_ENDPOINT"), envCfg.Get("OBS_HUAWEI_BUCKET"))
	if err != nil {
		panic(err)
	}

	grService := service.NewGrService(grRepository, warehouseStockRepository, stockRepository, ReplenishmentRepository, transactionDB)
	grBranchService := service.NewGrBranchService(OrderBookingRepository, grBranchRepository, warehouseStockRepository, stockRepository, transactionDB)
	arBranchService := service.NewArBranchService(arBranchRepository, warehouseStockRepository, stockRepository, transactionDB)
	bpprService := service.NewBpprService(bpprRepository, whStockRepository, stockRepository, transactionDB)
	smpIssService := service.NewSmpIssService(smpIssRepository, whStockRepository, stockRepository, transactionDB)
	whTrfService := service.NewWhTrfService(whTrfRepository, transactionDB, stockRepository, envCfg)
	itemStChService := service.NewItemStChService(itemStChRepository, transactionDB)
	whAdjService := service.NewWhAdjService(whAdjRepository, transactionDB, stockRepository)
	whSoService := service.NewWhSoService(whSoRepository, transactionDB)
	vanSoService := service.NewVanSoService(vanSoRepository, transactionDB)
	gdsService := service.NewGdsService(gdsRepository, transactionDB)
	vanBsService := service.NewVanBsUlService(vanBsUlRepository, transactionDB)
	vanUlService := service.NewVanUlService(vanUlRepository, transactionDB)
	vanLoService := service.NewVanLoService(vanLoRepository, transactionDB)
	stockService := service.NewStockService(stockRepository, transactionDB, validatorPkg)
	warehouseStockService := service.NewWarehouseStockService(warehouseStockRepository, transactionDB, validatorPkg)
	filesService := service.NewFilesService(envCfg, obsAdapter)
	SupplierReturnService := service.NewSupplierReturnService(SupplierReturnRepository, transactionDB, stockRepository)
	StockReturnService := service.NewStockReturnService(StockReturnRepository, stockRepository, transactionDB)
	StockOpnameService := service.NewStockOpnameService(StockOpnameRepository, transactionDB, stockRepository, warehouseStockRepository, obsAdapter)
	OrderBookingService := service.NewOrderBookingService(OrderBookingRepository, transactionDB)
	StockDisposalService := service.NewStockDisposalService(StockDisposalRepository, stockRepository, transactionDB, validatorPkg, envCfg)
	ReplenishmentService := service.NewReplenishmentService(ReplenishmentRepository, transactionDB)
	ReportsService := service.NewReportsService(ReportsRepository, obsAdapter)

	grController := controller.NewGrController(grService, validatorPkg)
	grBranchController := controller.NewGrBranchController(grBranchService, validatorPkg)
	arBranchController := controller.NewArBranchController(arBranchService, validatorPkg)
	bpprController := controller.NewBpprController(bpprService, validatorPkg)
	smpIssController := controller.NewSmpIssController(smpIssService, validatorPkg)
	whTrfController := controller.NewWhTrfController(whTrfService, validatorPkg)
	itemStChController := controller.NewItemStChController(itemStChService, validatorPkg)
	whAdjController := controller.NewWhAdjController(whAdjService, validatorPkg)
	whSoAdjController := controller.NewWhSoController(whSoService, validatorPkg)
	vanSoController := controller.NewVanSoController(vanSoService, validatorPkg)
	gdsController := controller.NewGdsController(gdsService, validatorPkg)
	vanBsUlController := controller.NewVanBsUlController(vanBsService, validatorPkg)
	vanUlController := controller.NewVanUlController(vanUlService, validatorPkg)
	vanLoController := controller.NewVanLoController(vanLoService, validatorPkg)
	stockController := controller.NewStockController(stockService, validatorPkg)
	warehouseStockController := controller.NewWarehouseStockController(warehouseStockService, validatorPkg)
	filesController := controller.NewFilesController(filesService, validatorPkg)
	SupplierReturnController := controller.NewSupplierReturnController(SupplierReturnService, validatorPkg)
	stockReturnController := controller.NewStockReturnController(StockReturnService, validatorPkg)
	stockOpnameController := controller.NewStockOpnameController(StockOpnameService, validatorPkg)
	orderBookingController := controller.NewOrderBookingController(OrderBookingService, validatorPkg)
	stockDisposalController := controller.NewStockDisposalController(StockDisposalService, validatorPkg)
	replenishmentController := controller.NewReplenishmentController(ReplenishmentService, validatorPkg)
	reportsController := controller.NewReportsController(ReportsService, validatorPkg)

	// define Fiber Framework config
	fiberCfg := config.NewFiberConfig(envCfg)
	app := fiber.New(fiberCfg)

	// middleware
	middleware.AppMiddleware(app) // Register Fiber's middleware for app.

	// This route path for test service is running "/"
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("It works")
	})
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Master Service Metrics Page"}))

	grController.Route(app)
	grBranchController.Route(app)
	arBranchController.Route(app)
	bpprController.Route(app)
	smpIssController.Route(app)
	whTrfController.Route(app)
	itemStChController.Route(app)
	whAdjController.Route(app)
	whSoAdjController.Route(app)
	vanSoController.Route(app)
	gdsController.Route(app)
	vanBsUlController.Route(app)
	vanUlController.Route(app)
	vanLoController.Route(app)
	stockController.Route(app)
	warehouseStockController.Route(app)
	filesController.Route(app)
	SupplierReturnController.Route(app)
	stockReturnController.Route(app)
	stockOpnameController.Route(app)
	orderBookingController.Route(app)
	stockDisposalController.Route(app)
	sapReplCfg := config.LoadSAPReplenishmentStatusConfig(envCfg)
	controller.RegisterSAPReplenishmentRoutes(app, sapReplCfg, replenishmentController)
	replenishmentController.Route(app)
	reportsController.Route(app)

	// start fiber server
	server.FiberServerWithGracefulShutdown(app)
}
