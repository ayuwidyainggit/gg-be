package thirdparty

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get outlet for picklist.
// @Description		Return list of outlet.
// @Param		    salesman_code query string false "SalesmanCode"
// @Produce		    application/json
// @Tags			outlet
// @Success         200 {object} response.Response{}
// @Router			/outlets-picklist/salesman [get]
// @Security        Bearer
func (controller *thirdPartyController) GetOutletPicklistBySalesCodes(ctx *gin.Context) {
	log.Info().Msg("get outlet picklist")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		webResponse := response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "empty token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		ctx.JSON(http.StatusUnauthorized, webResponse)
		return
	}

	var dataFilter model.OutletBySalesman

	dataFilter.SalesmanCode = strings.Split(ctx.Query("salesman_code"), ",")

	headers := map[string]string{
		"Authorization": authHeader,
		"Accept":        "application/json",
	}

	responses, pagination, err := controller.master.GetOutletPicklistBySalesCodes(c, dataFilter, headers, currentCustomerId)
	helper.ErrorPanic(err)

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
