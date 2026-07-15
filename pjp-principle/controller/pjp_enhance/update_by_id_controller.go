package pjpenhance

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

// UpdatePjpById godoc
// @Summary			Update pjp by ID.
// @Description		Update detail of a specific PJP by ID.
// @Param			id path string true "PJP ID"
// @Param			pjp body request.CreatePjpEnhanceRequest true "Update PJP payload"
// @Produce			application/json
// @Tags			pjp-enhance
// @Success			200 {object} response.Response{data=model.Pjp}
// @Failure			400 {object} response.Response{}
// @Failure			404 {object} response.Response{}
// @Router			/web/pjp-enhance/{id} [put]
// @Security        Bearer
func (controller *pjpEnhanceController) UpdatePjpById(ctx *gin.Context) {
	log.Info().Msg("update pjp enhance by id")

	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "PJP ID is required",
			Data:    nil,
		})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid PJP ID",
			Data:    nil,
		})
		return
	}

	var request request.CreatePjpEnhanceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	controller.pjpEnhanceService.UpdatePjp(c, idInt, request, currentCustomerId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "PJP updated successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
