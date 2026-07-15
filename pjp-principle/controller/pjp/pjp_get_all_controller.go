package pjp

import (
	"context"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags godoc
// @Summary			Get all pjp.
// @Description		Return list of pjp.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    pjp_code query string false "Pjp Code (%2C-separated for multiple values)"
// @Param		    operation_type query string false "OperationType"
// @Param		    team_salesman query string false "TeamSalesman (%2C-separated for multiple values)"
// @Param		    salesman_name query string false "SalesmanName"
// @Param		    salesman_code query string false "SalesmanCode"
// @Param		    q query string false "Q"
// @Param		    is_active query string false "Is Active"
// @Produce		    application/json
// @Tags			pjp
// @Success         200 {object} response.Pagination{}
// @Router			/web/pjp [get]
// @Security        Bearer
func (controller *pjpController) GetAll(ctx *gin.Context) {
	log.Info().Msg("Get all PJP")

	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "customer_id not found",
		})
		return
	}

	// Set request context with timeout
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// Parse pagination query safely
	limit := parseQueryInt(ctx, "limit", 10)
	page := parseQueryInt(ctx, "page", 1)

	// Extract filters from query
	filters := extractFilters(ctx)

	// Service call
	data, pagination, err := controller.pjpService.GetAll(reqCtx, limit, page, filters, customerID)
	if err != nil {
		helper.ErrorPanic(err)
	}

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
		Meta:   &pagination,
	}

	utils.ResponseInterceptor(reqCtx, &webResponse)
	ctx.JSON(http.StatusOK, webResponse)
}
