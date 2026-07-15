package thirdparty

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get salesman.
// @Description		Return list of salesman.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    sales_team_id query string false "SalesTeamId"
// @Param		    sort query string false "Sort"
// @Param		    is_active query string false "IsActive"
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/teams/salesman [get]
// @Security	   Bearer
func (controller *thirdPartyController) GetUnassignedSalesman(ctx *gin.Context) {
	log.Info().Msg("get salesman list")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter request.SalesmanListQueryFilter

	dataFilter.Limit = ctx.Query("limit")
	dataFilter.Page = ctx.Query("page")
	dataFilter.SalesTeamID = ctx.Query("sales_team_id")
	dataFilter.Sort = ctx.Query("sort")
	dataFilter.IsActive = ctx.Query("is_active")

	headers := make(map[string]string)
	headers["Accept"] = ctx.GetHeader("application/json")
	headers["Authorization"] = ctx.GetHeader("Authorization")

	responses, pagination := controller.master.GetUnassignedSalesman(c, headers, dataFilter, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
		Meta:   &pagination,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
