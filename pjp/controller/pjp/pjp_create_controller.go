package pjp

import (
	"context"
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
// @Param			pjp body request.PjpRequest true "Create pjp"
// @Produce			application/json
// @Tags			pjp
// @Success			201 {object} response.Response{}
// @Router			/web/pjp [post]
// @Security        Bearer
func (controller *pjpController) Create(ctx *gin.Context) {
	log.Info().Msg("Create PJP called")

	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "customer_id not found",
		})
		return
	}

	var req request.PjpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	controller.pjpService.Create(c, req, customerID)

	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created successfully",
	}

	utils.ResponseInterceptor(c, &webResponse)
	ctx.JSON(http.StatusCreated, webResponse)
}
