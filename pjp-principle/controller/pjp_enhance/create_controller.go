package pjpenhance

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
// @Summary			Create pjp.
// @Description		Save pjp data in Db.
// @Param			pjp body request.CreatePjpEnhanceRequest true "Create pjp"
// @Produce			application/json
// @Tags			pjp-enhance
// @Success			200 {object} response.Response{}
// @Router			/web/pjp-enhance [post]
// @Security        Bearer
func (controller *pjpEnhanceController) Create(ctx *gin.Context) {
	log.Info().Msg("create pjp enhance")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.CreatePjpEnhanceRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.pjpEnhanceService.Create(c, Request, currentCustomerId)
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
