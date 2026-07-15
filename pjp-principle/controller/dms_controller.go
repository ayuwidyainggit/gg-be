package controller

import (
	"scyllax-pjp/service"
)

type DmsController struct {
	dmsService service.DmsService
	// pjpService service.PjpService
}

func NewExternalApiController(service service.DmsService) *DmsController {
	return &DmsController{
		dmsService: service,
	}
}

// FindAllTags 		godoc
// @Summary			Get sales team.
// @Description		Return list of sales team.
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/teams [get]
// @Security        Bearer
// func (controller *DmsController) GetSalesTeam(ctx *gin.Context) {
// 	log.Info().Msg("get sales team")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	responses, err := controller.dmsService.GetSalesTeam(c, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get salesman By ID.
// @Description		Return of salesman by selected Id.
// @Param		    empId path string true "get salesman by id"
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/teams/salesman/{empId} [get]
// @Security	   Bearer
// func (controller *DmsController) GetSalesmanByID(ctx *gin.Context) {
// 	log.Info().Msg("get salesman")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		webResponse := response.Response{
// 			Code:    http.StatusUnauthorized,
// 			Status:  "UNAUTHORIZED",
// 			Message: "empty token",
// 		}
// 		utils.ResponseInterceptor(c, &webResponse)
// 		ctx.JSON(http.StatusUnauthorized, webResponse)
// 		return
// 	}

// 	ParamID := ctx.Param("empId")
// 	id, err := strconv.Atoi(ParamID)
// 	helper.ErrorPanic(err)

// 	headers := make(map[string]string)
// 	headers["Accept"] = ctx.GetHeader("application/json")
// 	headers["Authorization"] = ctx.GetHeader("Authorization")

// 	responses, err := controller.dmsService.GetSalesmanByID(c, id, headers, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get sales operation type name based on stored pjp.
// @Description		Return list of sales type name.
// @Param		    pjp_code query string false "Pjp Code (%2C-separated for multiple values)"
// @Param		    team_salesman query string false "TeamSalesman (%2C-separated for multiple values)"
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/teams/operation/type [get]
// @Security        Bearer
// func (controller *DmsController) GetSalesOperationType(ctx *gin.Context) {
// 	log.Info().Msg("get sales team")
// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	filters := make(map[string]interface{})
// 	if ctx.Query("pjp_code") != "" {
// 		filters["pjp_code"] = strings.Split(ctx.Query("pjp_code"), "%2C")
// 	}
// 	if ctx.Query("team_salesman") != "" {
// 		filters["team_salesman"] = strings.Split(ctx.Query("team_salesman"), "%2C")
// 	}

// 	responses := controller.pjpService.GetListOperationtype(c, filters)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get sales type name based on stored pjp.
// @Description		Return list of sales type name.
// @Param		    pjp_code query string false "Pjp Code (%2C-separated for multiple values)"
// @Param		    team_salesman query string false "TeamSalesman (%2C-separated for multiple values)"
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/teams/salesman/type [get]
// func (controller *DmsController) GetSalesTeamType(ctx *gin.Context) {
// 	log.Info().Msg("get sales team")
// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	filters := make(map[string]interface{})
// 	if ctx.Query("pjp_code") != "" {
// 		filters["pjp_code"] = strings.Split(ctx.Query("pjp_code"), "%2C")
// 	}
// 	if ctx.Query("team_salesman") != "" {
// 		filters["team_salesman"] = strings.Split(ctx.Query("team_salesman"), "%2C")
// 	}

// 	responses := controller.pjpService.GetlistSalesmanTeam(c, filters)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get warehouse.
// @Description		Return list of warehouse.
// @Produce		    application/json
// @Tags			warehouse
// @Success         200 {object} response.Response{}
// @Router			/warehouses [get]
// @Security        Bearer
// func (controller *DmsController) GetWarehouse(ctx *gin.Context) {
// 	log.Info().Msg("get warehouse")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	responses, err := controller.dmsService.GetWarehouse(c, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get outlet.
// @Description		Return list of outlet.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    sort query string false "Sort"
// @Param		    is_active query string false "IsActive"
// @Produce		    application/json
// @Tags			outlet
// @Success         200 {object} response.Response{}
// @Router			/outlets [get]
// @Security        Bearer
// func (controller *DmsController) GetOutlet(ctx *gin.Context) {
// 	log.Info().Msg("get outlet")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	var dataFilter model.DmsQueryFilter

// 	dataFilter.Limit = ctx.Query("limit")
// 	dataFilter.Page = ctx.Query("page")
// 	dataFilter.Sort = ctx.Query("sort")
// 	dataFilter.IsActive = ctx.Query("is_active")

// 	responses, pagination, err := controller.dmsService.GetOutlet(c, dataFilter, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 		Meta:   &pagination,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get list outlet.
// @Description		Return list of outlet.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    sort query string false "Sort"
// @Param		    outlet_type_name query string false "OutletTypeName"
// @Param		    outlet_group_name query string false "OutletGroupName"
// @Param		    is_active query string false "Filter by active status (1 = true, 2 = false, else = all)"
// @Produce		    application/json
// @Tags			outlet
// @Success         200 {object} response.Response{}
// @Router			/outlet/list [get]
// @Security        Bearer
// func (controller *DmsController) GetListOutlet(ctx *gin.Context) {
// 	log.Info().Msg("get list outlet")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	var dataFilter model.OutletQueryFilter

// 	dataFilter.Limit = ctx.Query("limit")
// 	dataFilter.Page = ctx.Query("page")
// 	dataFilter.Sort = ctx.Query("sort")
// 	dataFilter.OutletTypeName = ctx.Query("outlet_type_name")
// 	dataFilter.OutletGroupName = ctx.Query("outlet_group_name")
// 	dataFilter.IsActive = ctx.Query("is_active")

// 	responses, pagination, err := controller.dmsService.GetListOutlet(c, dataFilter, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 		Meta:   &pagination,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get outlet not assign.
// @Description		Return list of outlet not assign.
// @Param		    sort query string false "Sort"
// @Produce		    application/json
// @Tags			outlet
// @Success         200 {object} response.Response{}
// @Router			/outlets/not-assign [get]
// func (controller *DmsController) GetOutletNotAssign(ctx *gin.Context) {
// 	log.Info().Msg("get outlet not assign")
// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	var dataFilter model.DmsQueryFilter

// 	dataFilter.Sort = ctx.Query("sort")

// 	responses := controller.dmsService.GetOutletNotAssign(c, dataFilter)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get outlet.
// @Description		Return list of outlet.
// @Param		    salesman_code query string false "SalesmanCode"
// @Produce		    application/json
// @Tags			outlet
// @Success         200 {object} response.Response{}
// @Router			/outlets/salesman [get]
// @Security        Bearer
// func (controller *DmsController) GetOutletBySalesman(ctx *gin.Context) {
// 	log.Info().Msg("get outlet")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		webResponse := response.Response{
// 			Code:    http.StatusUnauthorized,
// 			Status:  "UNAUTHORIZED",
// 			Message: "empty token",
// 		}
// 		utils.ResponseInterceptor(c, &webResponse)
// 		ctx.JSON(http.StatusUnauthorized, webResponse)
// 		return
// 	}

// 	var dataFilter model.OutletBySalesman

// 	dataFilter.SalesmanCode = strings.Split(ctx.Query("salesman_code"), ",")

// 	headers := map[string]string{
// 		"Authorization": authHeader,
// 		"Accept":        "application/json",
// 	}

// 	responses, pagination, err := controller.dmsService.GetOutletBySalesman(c, dataFilter, headers, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 		Meta:   &pagination,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get outlets by salesman_id.
// @Description		Return list of outlet.
// @Param		    salesman_id query string true "SalesmanId"
// @Param		    search query string false "search"
// @Produce		    application/json
// @Tags			mobile
// @Success         200 {object} response.Response{}
// @Router			/mobile/outlets/salesman [get]
// @Security        Bearer
// func (controller *DmsController) MobileGetOutletsBySalesman(ctx *gin.Context) {
// 	log.Info().Msg("get outlets by salesman_id for mobile")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		webResponse := response.Response{
// 			Code:    http.StatusUnauthorized,
// 			Status:  "UNAUTHORIZED",
// 			Message: "empty token",
// 		}
// 		utils.ResponseInterceptor(c, &webResponse)
// 		ctx.JSON(http.StatusUnauthorized, webResponse)
// 		return
// 	}

// 	var dataFilter model.OutletBySalesmanId

// 	dataFilter.SalesmanId = ctx.Query("salesman_id")
// 	dataFilter.Search = ctx.Query("search")

// 	headers := map[string]string{
// 		"Authorization": authHeader,
// 		"Accept":        "application/json",
// 	}

// 	responses, pagination, err := controller.dmsService.GetOutletBySalesmanId(c, dataFilter, headers, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	webResponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   responses,
// 		Meta:   &pagination,
// 	}
// 	utils.ResponseInterceptor(c, &webResponse)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// FindAllTags 		godoc
// @Summary			Get list salesman by PJP.
// @Description		Return list of salesman.
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/list-salesman [get]
// @Security	   Bearer
// func (controller *DmsController) GetListSalesman(ctx *gin.Context) {
// 	log.Info().Msg("get salesman")
// 	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
// 	if !exists {
// 		helper.ErrorPanic(errors.New("customer_id not found"))
// 	}

// 	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
// 	defer cancel()

// 	responses, err := controller.dmsService.GetListSalesman(c, currentCustomerId)
// 	helper.ErrorPanic(err)

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, responses)
// }
