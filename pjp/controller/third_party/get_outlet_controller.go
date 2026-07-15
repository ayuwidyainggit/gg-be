package thirdparty

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get outlet.
// @Description		Return list of outlet.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    q query string false "Search"
// @Param		    outlet_code query string false "OutletCode"
// @Param		    outlet_id query int false "OutletID"
// @Param		    mode query string false "Mode"
// @Param		    sort query string false "Sort"
// @Param		    is_active query string false "IsActive"
// @Param		    sales_team_id query string false "SalesTeamId"
// @Produce		    application/json
// @Tags			outlet
// @Success         200 {object} response.Response{}
// @Router			/outlets [get]
// @Security        Bearer
func (controller *thirdPartyController) GetOutlet(ctx *gin.Context) {
	log.Info().Msg("get outlet")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter model.DmsQueryFilter

	dataFilter.Limit = ctx.Query("limit")
	dataFilter.Page = ctx.Query("page")
	dataFilter.Query = ctx.Query("q")
	dataFilter.OutletCode = ctx.Query("outlet_code")
	dataFilter.OutletID, _ = strconv.Atoi(ctx.Query("outlet_id"))
	dataFilter.Mode = ctx.Query("mode")
	dataFilter.Sort = ctx.Query("sort")
	dataFilter.IsActive = ctx.Query("is_active")
	dataFilter.SalesTeamID = ctx.Query("sales_team_id")

	responses, pagination := controller.master.GetOutlet(c, dataFilter, currentCustomerId)

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
