package router

import (
	"scyllax-tms/controller"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(
	shipmentController *controller.ShipmentController,
	dmsController *controller.DmsController,
	visitController *controller.VisitController,
	arriveController *controller.ArriveController,
	unloadController *controller.UnloadController,
	rejectController *controller.RejectController,
	outletController *controller.OutletController,
	productController *controller.ProductController,
	reportController *controller.ReportController,
	pickUpController *controller.PickUpController,
	picklistController *controller.PicklistController,
) *fiber.App {
	service := fiber.New()

	// send-pick third party
	service.Post("/login/send-pick", shipmentController.LoginSendPick)
	service.Post("/generate/send-pick", shipmentController.GenerateSendPick)

	//api third party
	service.Get("/vehicles/dev", dmsController.GetVehicle)
	service.Get("/vehicles", dmsController.GetVehicleByDms)
	service.Get("/reject-reason", dmsController.GetRejectReason)
	service.Get("/invoices", dmsController.GetListInvoice)
	service.Get("/returns", dmsController.GetListReturn)

	//mobile
	mobileRouter := service.Group("/mobile")
	mobileRouter.Post("/leave", visitController.Leave)
	mobileRouter.Post("/skip", visitController.Skip)

	//report
	reportRouter := mobileRouter.Group("/")
	reportRouter.Get("driver/reports", reportController.GetDriverReport)

	//arrive
	arriveRouter := mobileRouter.Group("/")
	arriveRouter.Post("arrive", arriveController.Arrive)

	//unload
	unloadRouter := mobileRouter.Group("/")
	unloadRouter.Post("unload", unloadController.Unload)
	unloadRouter.Post("resume", unloadController.Resume)
	unloadRouter.Post("onhold", unloadController.OnHold)
	unloadRouter.Get("todo/list/:outletId/:shipmentNo", unloadController.TravelList)

	//outlet
	outletRouter := mobileRouter.Group("/")
	outletRouter.Get("outlet/:driverId/:outletId/:shipmentNo", outletController.GetOutletByParams)
	outletRouter.Get("outlets", outletController.GetOutlet)

	//product
	productRouter := mobileRouter.Group("/")
	productRouter.Get("products", productController.GetProduct)

	//reject
	rejectRouter := mobileRouter.Group("/rejects")
	rejectRouter.Get("", rejectController.GetReject)
	rejectRouter.Get("/partial", rejectController.GetRejectPartial)
	rejectRouter.Post("", rejectController.RejectAll)
	rejectRouter.Post("/partial", rejectController.RejectPartial)
	rejectRouter.Post("/cancel", rejectController.RejectCancel)

	//visit
	visitRouter := mobileRouter.Group("/visits")
	visitRouter.Get("/summary/:driverId/:custId", visitController.GetSummaryByDriverIDAndCustID)
	visitRouter.Get("/daily/:shipmentNo/:custId", visitController.GetSummaryDailyByParams)
	visitRouter.Get("/daily-activity/:driverId", visitController.GetDailyActivityByDriverID)
	visitRouter.Post("/start", visitController.Start)
	visitRouter.Post("/end", visitController.End)

	//pickUp
	pickUpRouter := mobileRouter.Group("/pickup")
	pickUpRouter.Post("", pickUpController.PickUpAll)
	pickUpRouter.Post("/partial", pickUpController.PickUpPartial)
	pickUpRouter.Post("/skip", pickUpController.SkipPickUp)
	//web
	shipmentRouter := service.Group("/web/shipments")
	shipmentRouter.Get("", shipmentController.FindAll)
	shipmentRouter.Post("", shipmentController.CreateManual)
	shipmentRouter.Post("/auto", shipmentController.CreateAuto)
	shipmentRouter.Get("/:shipmentNo", shipmentController.FindByShipmentNo)
	shipmentRouter.Get("/invoices/:shipmentNo", shipmentController.FindByOrderNo)
	shipmentRouter.Patch("/submit", shipmentController.Update)
	shipmentRouter.Delete("/:shipmentNo", shipmentController.Delete)
	shipmentRouter.Post("/bulk", shipmentController.DeleteBulk)

	//web shipment report
	shipmentReportRouter := service.Group("/web/shipment-report")
	shipmentReportRouter.Get("/summary", reportController.GetShipmentReportSummary)
	shipmentReportRouter.Get("/detail", reportController.GetShipmentReportDetail)
	shipmentReportRouter.Get("/reject", reportController.GetShipmentReportReject)
	//shipment report dropdown
	shipmentReportRouter.Get("/shipment", reportController.GetShipmentNumberDropdown)
	shipmentReportRouter.Get("/product", reportController.GetProductCodeDropdown)
	shipmentReportRouter.Get("/driver", reportController.GetDriverDropdown)
	shipmentReportRouter.Get("/outlet", reportController.GetOutletDropdown)
	shipmentReportRouter.Get("/reason", reportController.GetReasonDropdown)

	picklistRouter := service.Group("/picklists")
	picklistRouter.Get("/invoices", picklistController.GetPicklistInvoice)
	picklistRouter.Post("", picklistController.CreatePicklist)
	picklistRouter.Put("", picklistController.UpdatePicklist)
	picklistRouter.Delete("", picklistController.DeletePicklist)
	picklistRouter.Get("/:id", picklistController.GetPicklist)
	picklistRouter.Get("", picklistController.GetAllPicklists)

	return service
}
