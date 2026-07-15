package pjpenhance

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

// GetPjpById		godoc
// @Summary			Get pjp by ID.
// @Description		Get detail of a specific pjp by ID.
// @Param			id path string true "PJP ID"
// @Produce			application/json
// @Tags			pjp-enhance
// @Success			200 {object} response.Response{data=response.PjpEnhanceResponse}
// @Failure			404 {object} response.Response{}
// @Router			/web/pjp-enhance/{id} [get]
// @Security        Bearer
func (controller *pjpEnhanceController) GetPjpById(ctx *gin.Context) {
	log.Info().Msg("get pjp enhnace by id")

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

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	pjpDetail := controller.pjpEnhanceService.GetById(c, idInt, currentCustomerId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Success",
		Data:    pjpDetail,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
