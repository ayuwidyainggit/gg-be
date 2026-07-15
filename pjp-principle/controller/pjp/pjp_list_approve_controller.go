package pjp

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get all pjp route with status approve dropdown.
// @Description		Return list of route with status approve dropdown.
// @Param		    q query string false "Q"
// @Produce		    application/json
// @Tags			pjp
// @Success         200 {object} response.Pagination{}
// @Router			/web/pjp/visit-list [get]
// @Security        Bearer
func (controller *pjpController) ListPjpApprove(ctx *gin.Context) {
	log.Info().Msg("list pjp")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	q := ctx.Query("q")
	responses := controller.pjpService.ListPjpApprove(c, q, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
