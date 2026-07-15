package main

import (
	"finance/adapter"
	"finance/controller"
	"finance/pkg/config"
	"finance/pkg/config/env"
	"finance/pkg/middleware"
	"finance/pkg/server"
	"finance/pkg/validation"
	"finance/repository"
	"finance/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	// load config env
	envCfg := env.NewCfgEnv()
	validatorPkg := validation.NewValiditor()

	postgreDB := config.PostgreSQLConnection(envCfg)

	transactionDB := repository.NewDbtransactionRepo(postgreDB)
	opexTrRepository := repository.NewOpexRepo(postgreDB)
	mApDiscRepository := repository.NewMApDiscRepo(postgreDB)
	mOpexTrRepository := repository.NewMOpexRepo(postgreDB)
	arPayRepository := repository.NewArPayRepo(postgreDB)
	chequeRepository := repository.NewChequeRepo(postgreDB)
	chequeGiroRepository := repository.NewChequeGiroRepo(postgreDB)
	bankTransferRepository := repository.NewBankTransferRepo(postgreDB)
	depositLookupRepository := repository.NewDepositLookupRepo(postgreDB)
	cashRepository := repository.NewCashRepo(postgreDB)
	arCndnRepository := repository.NewArCndnRepo(postgreDB)
	arSettlementRepository := repository.NewArSettlementRepo(postgreDB)
	arRepository := repository.NewArRepo(postgreDB)
	mChequeRejectRepository := repository.NewMChequeRejectRepo(postgreDB)
	apCndnRepository := repository.NewApCndnRepo(postgreDB)
	apRepository := repository.NewApRepo(postgreDB)
	apPayRepository := repository.NewApPayRepo(postgreDB)
	memoJrRepository := repository.NewMemoJrRepo(postgreDB)
	mCoaRepository := repository.NewMCoaRepo(postgreDB)
	mCoaTypeRepository := repository.NewMCoaTypeRepo(postgreDB)
	mApDistributorDiscountRepository := repository.NewApDistributorDiscountRepo(postgreDB)
	cndnRepository := repository.NewCndnRepo(postgreDB)
	apListRepository := repository.NewApListRepo(postgreDB)
	depositRepository := repository.NewDepositRepo(postgreDB)
	apSupplierInvoiceReturnRepository := repository.NewApSupplierInvoiceReturnRepo(postgreDB)
	cashBankReportRepository := repository.NewCashBankReportRepo(postgreDB)
	chequeGiroClearingRepository := repository.NewChequeGiroClearingRepo(postgreDB)
	apPaymentRepository := repository.NewApPaymentRepo(postgreDB)
	mTaxesRepository := repository.NewMTaxesRepo(postgreDB)
	taxesRepository := repository.NewTaxesRepo(postgreDB)
	vatRepository := repository.NewVatExtractRepo(postgreDB)
	coretaxVatExtractRepository := repository.NewCoreTaxVatExtractRepository(postgreDB)
	expenseRepository := repository.NewExpenseRepo(postgreDB)
	expenseEntryRepository := repository.NewExpenseEntryRepo(postgreDB)
	paymentDepositReportRepository := repository.NewPaymentDepositReportRepo(postgreDB)

	// setup adapter
	obsAdapter, err := adapter.InitObsAdapter(envCfg.Get("OBS_HUAWEI_AK"), envCfg.Get("OBS_HUAWEI_SK"), envCfg.Get("OBS_HUAWEI_ENDPOINT"), envCfg.Get("OBS_HUAWEI_BUCKET"))
	if err != nil {
		panic(err)
	}

	opexService := service.NewOpexService(opexTrRepository, transactionDB)
	apDiscService := service.NewMApDiscService(mApDiscRepository, transactionDB)
	mOpexService := service.NewMOpexService(mOpexTrRepository, transactionDB)
	arPayService := service.NewArPayService(arPayRepository, transactionDB)
	chequeService := service.NewChequeService(chequeRepository, transactionDB)
	chequeGiroService := service.NewChequeGiroService(chequeGiroRepository, transactionDB)
	bankTransferService := service.NewBankTransferService(bankTransferRepository, transactionDB)
	depositLookupService := service.NewDepositLookupService(depositLookupRepository, transactionDB)
	cashService := service.NewCashService(cashRepository, transactionDB)
	arCndnService := service.NewArCndnService(arCndnRepository, transactionDB)
	arSettlementService := service.NewArSettlementService(arSettlementRepository, depositRepository, transactionDB)
	arService := service.NewArService(arRepository, transactionDB)
	mChequeRejectService := service.NewMChequeRejectService(mChequeRejectRepository, transactionDB)
	apCndnService := service.NewApCndnService(apCndnRepository, transactionDB)
	apService := service.NewApService(apRepository, transactionDB)
	apPayService := service.NewApPayService(apPayRepository, transactionDB)
	memoJrService := service.NewMemoJrService(memoJrRepository, transactionDB)
	mCoaService := service.NewMCoaService(mCoaRepository, transactionDB)
	mCoaTypeService := service.NewMCoaTypeService(mCoaTypeRepository, transactionDB)
	filesService := service.NewFilesService(envCfg, obsAdapter)
	mApDistributorDiscountService := service.NewApDistributorDiscountService(mApDistributorDiscountRepository, transactionDB)
	cndnService := service.NewCndnService(cndnRepository, transactionDB)
	apListService := service.NewApListService(apListRepository, transactionDB)
	depositService := service.NewDepositService(depositRepository, transactionDB)
	apSupplierInvoiceReturnService := service.NewApSupplierInvoiceReturnService(apSupplierInvoiceReturnRepository, transactionDB)
	cashBankReportService := service.NewCashBankReportService(cashBankReportRepository, transactionDB)
	chequeGiroClearingService := service.NewChequeGiroClearingService(chequeGiroClearingRepository, transactionDB)
	apPaymentService := service.NewApPaymentService(apPaymentRepository, transactionDB)
	mTaxesService := service.NewMTaxesService(mTaxesRepository, taxesRepository, transactionDB)
	taxesService := service.NewTaxesService(taxesRepository, mTaxesRepository, transactionDB)
	vatExtractService := service.NewVatExtractService(vatRepository, apSupplierInvoiceReturnRepository, transactionDB)
	coretaxExtractVatService := service.NewCoreTaxVatExtractService(coretaxVatExtractRepository, transactionDB)
	expenseService := service.NewExpenseService(expenseRepository, transactionDB)
	expenseEntryService := service.NewExpenseEntryService(expenseEntryRepository, transactionDB)
	paymentDepositReportService := service.NewPaymentDepositReportService(paymentDepositReportRepository, transactionDB)

	opexController := controller.NewOpexController(opexService, validatorPkg)
	mOpexController := controller.NewMOpexController(mOpexService, validatorPkg)
	mApDiscController := controller.NewMApDiscController(apDiscService, validatorPkg)
	arPayController := controller.NewArPayController(arPayService, validatorPkg)
	chequeController := controller.NewChequeController(chequeService, validatorPkg)
	chequeGiroController := controller.NewChequeGiroController(chequeGiroService, validatorPkg)
	bankTransferController := controller.NewBankTransferController(bankTransferService, validatorPkg)
	depositLookupController := controller.NewDepositLookupController(depositLookupService, validatorPkg)
	cashController := controller.NewCashController(cashService, validatorPkg)
	arCndnController := controller.NewArCndnController(arCndnService, validatorPkg)
	arSettlementController := controller.NewArSettlementController(arSettlementService, validatorPkg)
	arController := controller.NewArController(arService, validatorPkg)
	mChequeRejectController := controller.NewMChequeRejectController(mChequeRejectService, validatorPkg)
	apCndnController := controller.NewApCndnController(apCndnService, validatorPkg)
	apController := controller.NewApController(apService, validatorPkg)
	apPayController := controller.NewApPayController(apPayService, validatorPkg)
	memoController := controller.NewMemoJrController(memoJrService, validatorPkg)
	mCoaController := controller.NewMCoaController(mCoaService, validatorPkg)
	mCoaTypeController := controller.NewMCoaTypeController(mCoaTypeService, validatorPkg)
	filesController := controller.NewFilesController(filesService, validatorPkg)
	mApDistribitorDiscountController := controller.NewApDistributorDiscountController(mApDistributorDiscountService, validatorPkg)
	cndnController := controller.NewCndnController(cndnService, validatorPkg)
	apListController := controller.NewApListController(apListService, validatorPkg)
	depositController := controller.NewDepositController(depositService, validatorPkg)
	apSupplierInvoiceReturnController := controller.NewApSupplierInvoiceReturnController(apSupplierInvoiceReturnService, validatorPkg)
	cashBankReportController := controller.NewCachBankReportController(cashBankReportService, validatorPkg)
	chequeGiroClearingController := controller.NewChequeGiroClearingController(chequeGiroClearingService, validatorPkg)
	apPaymentController := controller.NewApPaymentController(apPaymentService, validatorPkg)
	mTaxesController := controller.NewMTaxesController(mTaxesService, validatorPkg)
	taxesController := controller.NewTaxesController(taxesService, validatorPkg)
	vatExtractController := controller.NewVatExtractController(vatExtractService, validatorPkg)
	coretaxVatExtractController := controller.NewCoreTaxVatExtractController(coretaxExtractVatService, validatorPkg)
	expenseController := controller.NewExpenseController(expenseService, validatorPkg)
	expenseEntryController := controller.NewExpenseEntryController(expenseEntryService, validatorPkg)
	paymentDepositReportController := controller.NewPaymentDepositReportController(paymentDepositReportService, validatorPkg)

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

	opexController.Route(app)
	mOpexController.Route(app)
	mApDiscController.Route(app)
	arPayController.Route(app)
	chequeController.Route(app)
	chequeGiroController.Route(app)
	bankTransferController.Route(app)
	depositLookupController.Route(app)
	cashController.Route(app)
	arCndnController.Route(app)
	arSettlementController.Route(app)
	arController.Route(app)
	mChequeRejectController.Route(app)
	apCndnController.Route(app)
	apController.Route(app)
	apPayController.Route(app)
	memoController.Route(app)
	mCoaController.Route(app)
	mCoaTypeController.Route(app)
	filesController.Route(app)
	mApDistribitorDiscountController.Route(app)
	cndnController.Route(app)
	apListController.Route(app)
	depositController.Route(app)
	apSupplierInvoiceReturnController.Route(app)
	cashBankReportController.Route(app)
	chequeGiroClearingController.Route(app)
	apPaymentController.Route(app)
	mTaxesController.Route(app)
	taxesController.Route(app)
	vatExtractController.Route(app)
	coretaxVatExtractController.Route(app)
	expenseController.Route(app)
	expenseEntryController.Route(app)
	paymentDepositReportController.Route(app)

	// start fiber server
	server.FiberServerWithGracefulShutdown(app)
}
