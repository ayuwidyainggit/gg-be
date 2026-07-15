package thirdparty

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get salesman By ID.
// @Description		Return of salesman by selected Id.
// @Param		    empId path string true "get salesman by id"
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/teams/salesman/{empId} [get]
// @Security	   Bearer
func (controller *thirdPartyController) GetSalesmanByID(ctx *gin.Context) {
	log.Info().Msg("get salesman")
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

	ParamID := ctx.Param("empId")
	id, err := strconv.Atoi(ParamID)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Accept"] = ctx.GetHeader("application/json")
	headers["Authorization"] = ctx.GetHeader("Authorization")

	responses := controller.master.GetSalesmanByID(c, id, headers, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
