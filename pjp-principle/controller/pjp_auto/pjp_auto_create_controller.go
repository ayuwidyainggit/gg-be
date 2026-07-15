package pjpauto

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

// CreateTags		godoc
// @Summary			Create pjp auto.
// @Description		Create pjp auto data in Db.
// @Param			pjp body request.CreatePjpAuto true "Create pjp auto summary"
// @Produce			application/json
// @Tags			pjp auto
// @Success			200 {object} response.Response{}
// @Router			/web/pjp/auto [post]
// @Security        Bearer
func (controller *pjpAutoController) CreatePjpAuto(ctx *gin.Context) {
	log.Info().Msg("create optimized route")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	request := request.CreatePjpAuto{}

	err := ctx.ShouldBindJSON(&request)
	helper.ErrorPanic(err)

	controller.pjpAutoService.Create(c, request, currentCustomerId)
	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}
