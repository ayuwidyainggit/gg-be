package pjp

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UpdateTags		godoc
// @Summary			Update pjp.
// @Description		Update pjp data.
// @Param			pjpId path string true "update pjp by id"
// @Param			pjp body request.PjpRequest true  "Update pjp"
// @Tags			pjp
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/web/pjp/{pjpId} [patch]
// @Security        Bearer
func (controller *pjpController) Update(ctx *gin.Context) {
	log.Info().Msg("update pjp")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.PjpRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	ParamID := ctx.Param("pjpId")
	id, err := strconv.Atoi(ParamID)
	helper.ErrorPanic(err)
	Request.ID = id

	controller.pjpService.Update(c, Request, currentCustomerId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Updated successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
