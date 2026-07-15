package controller

import (
	"context"
	"net/http"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/service"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type VisitController struct {
	visitService service.VisitService
}

func NewVisitController(service service.VisitService) *VisitController {
	return &VisitController{
		visitService: service,
	}
}

/*
// FindAllTags 		    godoc
// @Summary				get outlet by sales_code and custId.
// @Param				salesCode path string true "sales_code"
// @Param				custId path string true "cust_id"
// @Param				date path string true "date"
// @Param				routeCode path string true "routeCode"
// @Description			Return list of route outlet by salesCode and custId.
// @Produce				application/json
// @Tags				mobile
// @Success				200 {object} response.Response{}
// @Router				/mobile/visits/outlet/{salesCode}/{custId}/{date}/{routeCode} [get]
*/
func (controller *VisitController) GetOutletBySalesCode(ctx *gin.Context) {
	log.Info().Msg("get outlet by sales_code")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	salesCode := ctx.Param("salesCode")
	custId := ctx.Param("custId")
	date := ctx.Param("date")
	routeCode := ctx.Param("routeCode")

	responses := controller.visitService.GetAllOutletBySalesCode(ctx, salesCode, custId, date, routeCode)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			end visit mobile.
// @Description		end visit.
// @Param			data body request.FinishVisitRequest true "end visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/end [post]
func (controller *VisitController) EndVisit(ctx *gin.Context) {
	log.Info().Msg("end visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.FinishVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.visitService.FinishVisit(c, Request)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "End Visit successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			summary visit mobile.
// @Description		summary visit and example date(yyyy-mm-dd).
// @Param		    salesman_code query string false "salesman_code"
// @Param		    date query string false "date"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/summary [get]
func (controller *VisitController) SummaryVisit(ctx *gin.Context) {
	log.Info().Msg("summary visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.SummaryQueryFilter

	dataFilter.SalesmanCode = ctx.Query("salesman_code")
	dataFilter.Date = ctx.Query("date")
	dataFilter.CustID = ctx.Query("custId")

	data := controller.visitService.SummaryVisit(c, dataFilter)
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// ArriveVisit handles arrive visit request with file URL and optional location update
// @Summary			arrive visit mobile.
// @Description		arrive visit with file URL and optional location update for approval. If is_update_location is true, location change will be saved to outlet_cr for approval.
// @Accept			multipart/form-data
// @Param			date formData string true "Visit date (format: yyyy-mm-dd)"
// @Param			salesman_code formData string true "Salesman code"
// @Param			cust_id formData string true "Customer ID"
// @Param			current_time formData int true "Current timestamp"
// @Param			id formData int true "Outlet visit list ID"
// @Param			outlet_id formData int false "Outlet ID"
// @Param			latitude formData string false "New latitude"
// @Param			longitude formData string false "New longitude"
// @Param			is_update_location formData bool true "Whether to update location (requires approval if true)"
// @Param			file_url formData string true "File URL"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/arrive [post]
func (controller *VisitController) ArriveVisit(ctx *gin.Context) {
	log.Info().Msg("arrive visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// Bind multipart form data
	req := request.ArriveVisitRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		helper.ErrorPanic(err)
	}

	// Validate file_url is provided (required according to docs)
	if req.FileUrl == "" {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "file_url is required",
			Data:    nil,
		}
		utils.ResponseInterceptor(c, &webResponse)
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	// Get cust_id and user_id from JWT context (set by middleware)
	custId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		// Try to get from request if not in context (for backward compatibility)
		custId = req.CustID
	}

	userId, _ := helper.GetCurrentUserId(ctx)

	controller.visitService.ArriveVisit(c, req, custId, userId)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Arrive successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			skip visit mobile.
// @Description		skip visit.
// @Param			data body request.SkipVisitRequest true "skip visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/skip [post]
func (controller *VisitController) SkipVisit(ctx *gin.Context) {
	log.Info().Msg("skip visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.SkipVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.visitService.SkipVisit(c, Request)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Skip successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			resume visit mobile.
// @Description		resume visit.
// @Param			data body request.ResumeVisitRequest true "skip visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/resume [post]
func (controller *VisitController) ResumeVisit(ctx *gin.Context) {
	log.Info().Msg("skip visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.ResumeVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.visitService.ResumeVisit(c, Request)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Resume successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			leave visit mobile.
// @Description		leave visit.
// @Param			data body request.LeaveVisitRequest true "leave visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/leave [post]
func (controller *VisitController) LeaveVisit(ctx *gin.Context) {
	log.Info().Msg("leave visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.LeaveVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.visitService.LeaveVisit(c, Request)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Leave successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Note             godoc
// @Summary			todo list
// @Description		todo list arrive and leave.
// @Param		    outletVisitId path string true "todo list arrive and leave"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/todo/list/{outletVisitId} [get]
func (controller *VisitController) TravelList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	params := ctx.Param("outletVisitId")
	outletVisitId, err := strconv.Atoi(params)
	helper.ErrorPanic(err)
	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "customer_id not found",
		})
		return
	}

	data := controller.visitService.TravelList(c, outletVisitId, customerID)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			onhold visit mobile.
// @Description		onhold visit.
// @Param			data body request.OnholdVisitRequest true "onhold visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/onhold [post]
func (controller *VisitController) OnholdVisit(ctx *gin.Context) {
	log.Info().Msg("onhold visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.OnholdVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.visitService.UnloadVisit(c, Request)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "onhold successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			outlet visit mobile.
// @Description		outlet visit.
// @Param			data body request.OutletVisitRequest true "outlet visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/outlet [post]
func (controller *VisitController) OutletVisit(ctx *gin.Context) {
	log.Info().Msg("outlet visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.OutletVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)
	Request.EmpID, _ = helper.GetCurrentEmpId(ctx)

	controller.visitService.OutletVisit(c, Request)
	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Outlet Visit successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}

// CreateTags		godoc
// @Summary			summary visit Status mobile.
// @Description		summary visit and example date(yyyy-mm-dd).
// @Param		    salesman_code query string false "salesman_code"
// @Param		    date query string false "date"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/status [get]
func (controller *VisitController) SummaryVisitStatus(ctx *gin.Context) {
	log.Info().Msg("summary visit status")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.SummaryQueryFilter

	dataFilter.SalesmanCode = ctx.Query("salesman_code")
	dataFilter.Date = ctx.Query("date")
	dataFilter.CustID = ctx.Query("custId")

	data := controller.visitService.SummaryVisitStatus(c, dataFilter)
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			outlet visit mobile.
// @Description		outlet visit and example date(yyyy-mm-dd).
// @Param		    salesman_code query string true "salesman_code"
// @Param			cust_id query string true "cust_id"
// @Param		    date query string true "date"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/outlet/list [get]
func (controller *VisitController) GetOutletVisitList(ctx *gin.Context) {
	log.Info().Msg("outlet visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.SummaryQueryFilter

	dataFilter.SalesmanCode = ctx.Query("salesman_code")
	dataFilter.CustID = ctx.Query("cust_id")
	dataFilter.Date = ctx.Query("date")
	dataFilter.EmpID, _ = helper.GetCurrentEmpId(ctx)

	data := controller.visitService.GetVisitOutletList(c, dataFilter)
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			salesman report mobile.
// @Description		salesman report and example date(yyyy-mm-dd).
// @Param		    salesman_id query string true "salesman_id"
// @Param		    date query string false "date"
// @Param		    year query string false "year"
// @Param		    month query string false "month"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/salesman/report [get]
func (controller *VisitController) GetSalesmanReport(ctx *gin.Context) {
	log.Info().Msg("outlet visit")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.SalesmanReportQueryFilter

	dataFilter.Date = ctx.Query("date")
	dataFilter.SalesmanId = ctx.Query("salesman_id")
	dataFilter.Month = ctx.Query("month")
	dataFilter.Year = ctx.Query("year")

	data := controller.visitService.GetSalesmanReport(c, dataFilter)
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
