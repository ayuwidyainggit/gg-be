package thirdparty

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get outlet.
// @Description		Return list of outlet.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    sort query string false "Sort"
// @Param		    is_active query string false "IsActive"
// @Produce		    application/json
// @Tags			destination
// @Success         200 {object} response.Response{}
// @Router			/distributors [get]
// @Security        Bearer
func (controller *thirdPartyController) GetDistributor(ctx *gin.Context) {
	log.Info().Msg("get outlet")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter model.DistributorQueryFilter

	dataFilter.Limit = ctx.Query("limit")
	dataFilter.Page = ctx.Query("page")
	dataFilter.Sort = ctx.Query("sort")
	dataFilter.Query = ctx.Query("q")

	responses, pagination := controller.master.GetDistributor(c, dataFilter, currentCustomerId)

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
