package pjp

import (
	"context"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags godoc
// @Summary			Get all pjp with route for dropdown.
// @Description		Return list of with route for dropdown.
// @Param		    q query string false "Q"
// @Produce		    application/json
// @Tags			pjp
// @Success         200 {object} response.Pagination{}
// @Router			/web/pjp/list [get]
// @Security        Bearer
func (controller *pjpController) GetPjpWithRoute(ctx *gin.Context) {
	log.Info().Msg("Get PJP with route (dropdown list)")

	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "customer_id not found",
		})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	query := ctx.Query("q")
	data := controller.pjpService.GetPjpWithRoute(reqCtx, query, customerID)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}

	utils.ResponseInterceptor(reqCtx, &webResponse)
	ctx.JSON(http.StatusOK, webResponse)
}
