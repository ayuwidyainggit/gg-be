package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/rabbitmq"
	"sales/pkg/responsebuild"
	"sales/pkg/structs"
	"sales/pkg/validation"
	"sales/service"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
)

// secondarySalesExportBody is the ONLY struct that should be used for BodyParser on the
// POST /secondary-sales export endpoint. Never bind request body directly into
// entity.SecondarySalesReportQueryFilter — its CustID/ParentCustID fields are auth-only
// and must always be sourced from JWT locals, not from the request body.
type rawSecondarySalesExportBody struct {
	RequestedCustIDRaw json.RawMessage `json:"cust_id"`
	From               *int64          `json:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To                 *int64          `json:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Sort               string          `json:"sort"`
	Page               int             `json:"page"`
	Limit              int             `json:"limit"`
	DistributorIDs     []int64         `json:"distributor_ids"`
	SalesmanIDs        []int64         `json:"salesman_ids"`
	OutletIDs          []int64         `json:"outlet_ids"`
	ProIDs             []int64         `json:"pro_ids"`
}

type rawActivityReportExportBody struct {
	RequestedCustIDRaw json.RawMessage `json:"cust_id"`
	DistributorCodeRaw json.RawMessage `json:"distributor_code"`
	SalesmanIDs        []int           `json:"salesman_ids"`
	FromDate           string          `json:"from" validate:"required"`
	ToDate             string          `json:"to" validate:"required"`
	Sort               string          `json:"sort"`
	Page               int             `json:"page"`
	Limit              int             `json:"limit"`
}

type activityReportExportBody struct {
	RequestedCustID  string                    `json:"-" validate:"omitempty,alphanum,max=20"`
	RequestedCustIDs entity.StringListOrScalar `json:"cust_id,omitempty"`
	DistributorCodes []string                  `json:"distributor_code,omitempty"`
	SalesmanIDs      []int                     `json:"salesman_ids"`
	FromDate         string                    `json:"from" validate:"required"`
	ToDate           string                    `json:"to" validate:"required"`
	Sort             string                    `json:"sort"`
	Page             int                       `json:"page"`
	Limit            int                       `json:"limit"`
}

type secondarySalesExportBody struct {
	RequestedCustID  string                    `json:"-" validate:"omitempty,alphanum,max=20"`
	RequestedCustIDs entity.StringListOrScalar `json:"cust_id,omitempty"`
	From             *int64                    `json:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To               *int64                    `json:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Sort             string                    `json:"sort"`
	Page             int                       `json:"page"`
	Limit            int                       `json:"limit"`
	DistributorIDs   []int64                   `json:"distributor_ids"`
	SalesmanIDs      []int64                   `json:"salesman_ids"`
	OutletIDs        []int64                   `json:"outlet_ids"`
	ProIDs           []int64                   `json:"pro_ids"`
}

type ReportController struct {
	ReportService service.ReportService
	validator     *validation.Validate
}

func NewReportController(reportService service.ReportService, validator *validation.Validate) *ReportController {
	return &ReportController{
		ReportService: reportService,
		validator:     validator,
	}
}
func (controller *ReportController) Route(app *fiber.App) {
	go rabbitmq.Subscribe(constant.RMQ_SECONDARY_SALES_EXPORT, controller.processSecondarySalesExportMessage)
	go rabbitmq.Subscribe(constant.RMQ_SALESMAN_ACTIVITY_REPORT_SALES_EXPORT, controller.SalesmanActivityReportExportMessage)

	// qParamId := ":ro_no"
	reportRouteV1 := app.Group("/v1/reports", middleware.JWTProtected())
	reportRouteV1.Post("/secondary-sales", controller.SecondarySales)
	reportRouteV1.Get("", controller.List)
	reportRouteV1.Get("/secondary-sales/sum-date", controller.SecondaryReportSalesSumMonth)
	reportRouteV1.Get("/secondary-sales/group", controller.SecondaryReportSalesGroup)
	reportRouteV1.Get("/secondary-sales/trend-sales", controller.SecondaryReportSalesTrendSales)

	reportRouteV1.Post("/activity-report-sales", controller.ActivityReportSales)
	reportRouteV1.Get("/activity-report-sales", controller.ActivityReportSalesList)
	reportRouteV1.Get("/activity-report-sales/sum-date", controller.SalesmanActivitySumMonth)
	reportRouteV1.Get("/activity-report-sales/trend-sales", controller.SalesmanActivityTrendSales)
	reportRouteV1.Get("/activity-report-sales/group", controller.SalesmanActivityReportSalesGroup)
	reportRouteV1.Get("/activity-report-sales/geotag", controller.SalesmanActivityGeotag)
	reportRouteV1.Get("/activity-report-sales/salesman-list", controller.SalesmanActivitySalesmanList)

	reportRouteExtract := app.Group("/v1/extract")
	reportRouteExtract.Post("/secondary-sales", controller.SecondarySalesDashboardExtract)

}

func (controller *ReportController) Cron() gocron.Scheduler {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic("error loading Asia/Jakarta timezone")
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		panic("error running cron")
	}

	j, err := s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(00, 01, 0),
			),
		),
		gocron.NewTask(func() {
			dateExtract := time.Now().Add(-24 * time.Hour)
			controller.SecondarySalesDashboardExtractCron(dateExtract)
		},
		),
	)
	if err != nil {
		panic(fmt.Sprintf("error running job %v", err.Error()))
	}
	// each job has a unique id
	fmt.Println(fmt.Sprintf("running job %v", j.ID()))

	// start the scheduler
	s.Start()

	return s
}

func (controller *ReportController) SecondarySalesDashboardExtract(c *fiber.Ctx) error {
	var headerAcceptLang string
	var dataPayload entity.SecondarySalesReportExtractPayload
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.BodyParser(&dataPayload); err != nil {
		log.Error("ReturnController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataPayload, headerAcceptLang)
	if errs != nil {
		log.Error("InvoiceController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	t := time.Date(dataPayload.Year, time.Month(dataPayload.Month), dataPayload.Day, 0, 0, 0, 0, time.Local)

	err := controller.SecondarySalesDashboardExtractCron(t)

	if err != nil {
		log.Error("ReportController, List, ValidateStruct(params), errs:", err)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	return nil
}

func (controller *ReportController) SecondarySalesDashboardExtractCron(date time.Time) error {
	var (
		request entity.SecondarySalesReportDashboardExtractQueryFilter
	)

	request.Date = date
	log.Info("Report - Secondary Sales report extract")
	err := controller.ReportService.ExtractReportSecondary(request)
	if err != nil {
		return err
	}
	return nil
}

func (controller *ReportController) SecondarySales(c *fiber.Ctx) error {
	var rawBody rawSecondarySalesExportBody
	var body secondarySalesExportBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&rawBody); err != nil {
		log.Error("ReportController, SecondarySales, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	body = secondarySalesExportBody{
		From:           rawBody.From,
		To:             rawBody.To,
		Sort:           rawBody.Sort,
		Page:           rawBody.Page,
		Limit:          rawBody.Limit,
		DistributorIDs: rawBody.DistributorIDs,
		SalesmanIDs:    rawBody.SalesmanIDs,
		OutletIDs:      rawBody.OutletIDs,
		ProIDs:         rawBody.ProIDs,
	}
	if len(rawBody.RequestedCustIDRaw) > 0 {
		if err := json.Unmarshal(rawBody.RequestedCustIDRaw, &body.RequestedCustIDs); err != nil {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if len(body.RequestedCustIDs) == 1 {
			body.RequestedCustID = body.RequestedCustIDs[0]
		}
	}

	errs := controller.validator.ValidateStruct(body, headerAcceptLang)
	if errs != nil {
		log.Error("ReportController, SecondarySales, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Build filter from DTO — auth fields always come from JWT locals, never from body.
	requestedCustIDs := []string(body.RequestedCustIDs)
	requestedCustID := strings.TrimSpace(body.RequestedCustID)
	if len(requestedCustIDs) == 0 && requestedCustID != "" {
		requestedCustIDs = []string{requestedCustID}
	}
	if normalized, err := entity.NormalizeStringList(requestedCustIDs); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	} else {
		requestedCustIDs = normalized
	}
	if len(requestedCustIDs) == 1 {
		requestedCustID = requestedCustIDs[0]
	} else if len(requestedCustIDs) == 0 {
		requestedCustID = ""
	}

	request := entity.SecondarySalesReportQueryFilter{
		CustID:           c.Locals("cust_id").(string),
		ParentCustID:     c.Locals("parent_cust_id").(string),
		ExportBy:         c.Locals("user_fullname").(string),
		RequestedCustID:  requestedCustID,
		RequestedCustIDs: requestedCustIDs,
		From:             body.From,
		To:               body.To,
		Sort:             body.Sort,
		Page:             body.Page,
		Limit:            body.Limit,
		DistributorIDs:   body.DistributorIDs,
		SalesmanIDs:      body.SalesmanIDs,
		OutletIDs:        body.OutletIDs,
		ProIDs:           body.ProIDs,
	}

	data, err := controller.ReportService.PublishSecondarySalesReport(request)
	if err != nil {
		log.Error("ReportController, SecondarySales, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: 1,
		PageCurrent: request.Page,
		PageLimit:   request.Limit,
		PageTotal:   1,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) List(c *fiber.Ctx) error {
	var (
		dataFilter entity.ReportQueryFilter
		data       []entity.ReportList
	)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReportController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReportController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.ReportService.List(dataFilter)
	if err != nil {
		log.Error("ReportController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) processSecondarySalesExportMessage(msg amqp.Delivery) {
	log.Infof("Processing message: %s", string(msg.Body))

	// Step 1: Unmarshal the message body into a JSON object
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(msg.Body, &jsonBody); err != nil {
		log.Errorf("Failed to unmarshal JSON body: %v", err)
		msg.Nack(false, false) // Requeue: false, multiple: false
		return
	}

	// Step 2: Map the JSON body to the request struct
	var request entity.SecondarySalesReportQueryFilter
	if err := structs.Automapper(jsonBody, &request); err != nil {
		log.Errorf("Failed to map JSON body to request struct: %v", err)
		msg.Nack(false, false)
		return
	}

	// Step 3: Call the service layer to process the message
	if err := controller.ReportService.SubscribeSecondarySalesReport(request); err != nil {
		log.Errorf("Failed to process message in service: %v", err)
		msg.Nack(false, false)
		return
	}

	// Acknowledge the message after successful processing
	if err := msg.Ack(false); err != nil {
		log.Errorf("Failed to acknowledge message: %v", err)
	}

	log.Infof("Message processed successfully: %s", string(msg.Body))
}

func (controller *ReportController) ActivityReportSales(c *fiber.Ctx) error {
	var (
		rawBody rawActivityReportExportBody
		body    activityReportExportBody
		request entity.ActivityReportQueryFilter
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&rawBody); err != nil {
		log.Error("ReportController, ActivityReportSales, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	body = activityReportExportBody{
		SalesmanIDs: rawBody.SalesmanIDs,
		FromDate:    rawBody.FromDate,
		ToDate:      rawBody.ToDate,
		Sort:        rawBody.Sort,
		Page:        rawBody.Page,
		Limit:       rawBody.Limit,
	}
	if len(rawBody.RequestedCustIDRaw) > 0 {
		if err := json.Unmarshal(rawBody.RequestedCustIDRaw, &body.RequestedCustIDs); err != nil {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if len(body.RequestedCustIDs) == 1 {
			body.RequestedCustID = body.RequestedCustIDs[0]
		}
	}
	if len(rawBody.DistributorCodeRaw) > 0 {
		var distributorCodes entity.StringListOrScalar
		if err := json.Unmarshal(rawBody.DistributorCodeRaw, &distributorCodes); err != nil {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		body.DistributorCodes = distributorCodes
	}

	requestedCustIDs, err := resolveActivityReportCustIDs(c, []string(body.RequestedCustIDs))
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	requestedCustID := ""
	if len(requestedCustIDs) == 1 {
		requestedCustID = requestedCustIDs[0]
	}

	distributorCodes, err := resolveActivityReportDistributorCodes(c, body.DistributorCodes)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.AuthCustID = c.Locals("cust_id").(string)
	request.RequestedCustID = requestedCustID
	request.RequestedCustIDs = requestedCustIDs
	request.CustIDs = requestedCustIDs
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.SalesmanIDs = body.SalesmanIDs
	request.DistributorCodes = distributorCodes
	request.FromDate = body.FromDate
	request.ToDate = body.ToDate
	request.Sort = body.Sort
	request.Page = body.Page
	request.Limit = body.Limit
	request.ExportBy = c.Locals("user_fullname").(string)
	request.IsAdmin = c.Locals("is_admin").(bool)
	request.DistPriceGrpID = c.Locals("dist_price_grp_id").(int)

	data, err := controller.ReportService.PublishActivitySalesReport(request)
	if err != nil {
		log.Error("ReportController, ActivityReportSales, data, err:", err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(constant.MsgActivityReportSalesExportFailed)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: 1,
		PageCurrent: request.Page,
		PageLimit:   request.Limit,
		PageTotal:   1,
	})
	responsePayload.Setmsg(constant.MsgActivityReportSalesExportSuccess)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) ActivityReportSalesList(c *fiber.Ctx) error {
	var (
		request entity.ActivityReportQueryFilterList
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("ReturnController, List, query parser filter:", err.Error())
	}

	requestedCustIDs, err := resolveActivityReportCustIDs(c, nil)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	distributorCodes, err := resolveActivityReportDistributorCodes(c, nil)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.AuthCustID = c.Locals("cust_id").(string)
	request.CustID = request.AuthCustID
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.RequestedCustID = requestedCustIDs[0]
	} else {
		request.RequestedCustID = ""
	}
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.DistributorCodes = distributorCodes
	request.IsAdmin = c.Locals("is_admin").(bool)
	request.DistPriceGrpID = c.Locals("dist_price_grp_id").(int)

	data, total, lastPage, err := controller.ReportService.PublishActivitySalesReportList(request)
	if err != nil {
		log.Error("ReportController, ActivityReportSalesList, data, err:", err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(constant.MsgActivityReportSalesListFailed)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: request.Page,
		PageLimit:   request.Limit,
		PageTotal:   lastPage,
	})
	responsePayload.Setmsg(constant.MsgActivityReportSalesListSuccess)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SecondaryReportSalesSumMonth(c *fiber.Ctx) error {
	var (
		request entity.SecondarySalesReportDashboardSumPayload
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	requestedCustIDs, err := parseSecondarySalesCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(requestedCustIDs) == 0 && strings.TrimSpace(request.CustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(request.CustID)}
	}
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.CustID = requestedCustIDs[0]
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SecondarySalesReportSumReportByMonth(authCustID, parentCustID, request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SecondaryReportSalesTrendSales(c *fiber.Ctx) error {
	var (
		request entity.SecondarySalesReportTrensSalesSumPayload
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("report controller, SecondaryReportSalesTrendSales, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	// GET body is optional. If body is empty, keep zero-value request.CustID and fall back in service.
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&request); err != nil {
			log.Error("report controller, SecondaryReportSalesTrendSales, body parser:", err.Error())
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	requestedCustIDs, err := parseSecondarySalesCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(requestedCustIDs) == 0 && request.CustID != "" {
		requestedCustIDs = []string{request.CustID}
	}
	if normalized, err := entity.NormalizeStringList(requestedCustIDs); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	} else {
		request.CustIDs = normalized
	}
	if len(request.CustIDs) == 1 {
		request.CustID = request.CustIDs[0]
	} else if len(request.CustIDs) == 0 {
		request.CustID = ""
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SecondarySalesReportTrendSales(authCustID, parentCustID, request.Year, request.CustIDs)
	if err != nil {
		log.Error("report controller, SecondaryReportSalesTrendSales, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SecondaryReportSalesGroup(c *fiber.Ctx) error {
	var (
		request entity.SecondarySalesReportDashboardGroupPayload
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	requestedCustIDs, err := parseSecondarySalesCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.CustID = requestedCustIDs[0]
	} else if len(requestedCustIDs) == 0 {
		request.CustID = ""
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SecondarySalesReportGroupSales(authCustID, parentCustID, request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SalesmanActivityReportExportMessage(msg amqp.Delivery) {
	log.Infof("Processing message: %s", string(msg.Body))

	// Step 1: Unmarshal the message body into a JSON object
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(msg.Body, &jsonBody); err != nil {
		log.Errorf("Failed to unmarshal JSON body: %v", err)
		msg.Nack(false, false) // Requeue: false, multiple: false
		return
	}

	// Step 2: Map the JSON body to the request struct
	var request entity.ActivityReportQueryFilter
	if err := structs.Automapper(jsonBody, &request); err != nil {
		log.Errorf("Failed to map JSON body to request struct: %v", err)
		msg.Nack(false, false)
		return
	}

	// Step 3: Call the service layer to process the message
	if err := controller.ReportService.SubscribeActivitySalesReport(request); err != nil {
		log.Errorf("Failed to process message in service: %v", err)
		msg.Nack(false, false)
		return
	}

	// Acknowledge the message after successful processing
	if err := msg.Ack(false); err != nil {
		log.Errorf("Failed to acknowledge message: %v", err)
	}

	log.Infof("Message processed successfully: %s", string(msg.Body))
}

func (controller *ReportController) SalesmanActivitySumMonth(c *fiber.Ctx) error {
	var (
		request entity.SalesmanActivityReportDashboardSumPayload
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("ReturnController, List, query parser filter:", err.Error())
	}

	requestedCustIDs, err := parseActivityReportCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(requestedCustIDs) == 0 && strings.TrimSpace(request.CustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(request.CustID)}
	}
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.CustID = requestedCustIDs[0]
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SalesmanActivityReportSumReportByMonth(authCustID, parentCustID, request)
	if err != nil {
		log.Error("ReportController, SecondarySales, data, err:", err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SalesmanActivityTrendSales(c *fiber.Ctx) error {
	var (
		request entity.ActivityReportTrendSalesPayload
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("ReportController, SalesmanActivityTrendSales, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	requestedCustIDs, err := parseActivityReportCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(requestedCustIDs) == 0 && strings.TrimSpace(request.CustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(request.CustID)}
	}
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.CustID = requestedCustIDs[0]
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SalesmanActivityReportTrendSales(authCustID, parentCustID, request.Year, request.CustIDs)
	if err != nil {
		log.Error("ReportController, SalesmanActivityTrendSales, data, err:", err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SalesmanActivityGeotag(c *fiber.Ctx) error {
	var request entity.ActivityReportGeotagPayload
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("ReportController, SalesmanActivityGeotag, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	requestedCustIDs, err := parseActivityReportCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(requestedCustIDs) == 0 && strings.TrimSpace(request.CustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(request.CustID)}
	}
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.CustID = requestedCustIDs[0]
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SalesmanActivityReportGeotag(authCustID, parentCustID, request)
	if err != nil {
		log.Error("ReportController, SalesmanActivityGeotag, data, err:", err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ReportController) SalesmanActivityReportSalesGroup(c *fiber.Ctx) error {
	var (
		request entity.SalesmanActivityReportDashboardGroupPayload
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("ReturnController, List, query parser filter:", err.Error())
	}

	requestedCustIDs, err := parseActivityReportCustIDQuery(c)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(requestedCustIDs) == 0 && strings.TrimSpace(request.CustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(request.CustID)}
	}
	request.CustIDs = requestedCustIDs
	if len(requestedCustIDs) == 1 {
		request.CustID = requestedCustIDs[0]
	}

	authCustID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	data, err := controller.ReportService.SalesmanActivityReportGroupSales(authCustID, parentCustID, request)
	if err != nil {
		log.Error("ReportController, SecondarySales, data, err:", err.Error())
		if errors.Is(err, service.ErrUnauthorizedCustID) {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func parseSecondarySalesCustIDQuery(c *fiber.Ctx) ([]string, error) {
	return parseActivityReportCustIDQuery(c)
}

func parseActivityReportCustIDQuery(c *fiber.Ctx) ([]string, error) {
	raw := make([]string, 0)
	for _, key := range []string{"cust_id", "cust_id[]"} {
		for _, value := range c.Context().QueryArgs().PeekMulti(key) {
			raw = append(raw, string(value))
		}
	}
	if len(raw) == 0 {
		if value := strings.TrimSpace(c.Query("cust_id")); value != "" {
			raw = append(raw, value)
		}
	}
	return entity.NormalizeStringList(raw)
}

func parseActivityReportDistributorCodeQuery(c *fiber.Ctx) ([]string, error) {
	raw := make([]string, 0)
	for _, key := range []string{"distributor_code", "distributor_code[]"} {
		for _, value := range c.Context().QueryArgs().PeekMulti(key) {
			raw = append(raw, string(value))
		}
	}
	if len(raw) == 0 {
		if value := strings.TrimSpace(c.Query("distributor_code")); value != "" {
			raw = append(raw, value)
		}
	}
	return entity.NormalizeDistributorCodeList(raw)
}

func resolveActivityReportCustIDs(c *fiber.Ctx, bodyCustIDs []string) ([]string, error) {
	queryCustIDs, err := parseActivityReportCustIDQuery(c)
	if err != nil {
		return nil, err
	}

	raw := append([]string{}, queryCustIDs...)
	raw = append(raw, bodyCustIDs...)
	if len(raw) == 0 {
		return nil, nil
	}
	return entity.NormalizeStringList(raw)
}

func resolveActivityReportDistributorCodes(c *fiber.Ctx, bodyCodes []string) ([]string, error) {
	queryCodes, err := parseActivityReportDistributorCodeQuery(c)
	if err != nil {
		return nil, err
	}

	raw := append([]string{}, queryCodes...)
	raw = append(raw, bodyCodes...)
	if len(raw) == 0 {
		return nil, nil
	}
	return entity.NormalizeDistributorCodeList(raw)
}

func (controller *ReportController) SalesmanActivitySalesmanList(c *fiber.Ctx) error {
	var (
		request entity.ActivityReportSalesmanListQueryFilter
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&request); err != nil {
		log.Error("SalesmanActivitySalesmanList, List, query parser filter:", err.Error())
	}

	request.CustID = c.Locals("cust_id").(string)

	data, err := controller.ReportService.SalesmanActivitySalesmanList(request)
	if err != nil {
		log.Error("ReportController, SalesmanActivitySalesmanList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}
