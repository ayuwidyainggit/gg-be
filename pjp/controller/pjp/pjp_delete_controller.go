package pjp

import (
	"context"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// DeleteTags		godoc
// @Summary			Delete pjp.
// @Param		    pjpId path string true "delete pjp by id"
// @Description		Remove pjp data by id.
// @Produce			application/json
// @Tags			pjp
// @Success			200 {object} response.Response{}
// @Router			/web/pjp/{pjpId} [delete]
// @Security        Bearer
func (controller *pjpController) Delete(ctx *gin.Context) {
	log.Info().Msg("delete pjp")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	ParamID := ctx.Param("pjpId")
	id, err := strconv.Atoi(ParamID)
	helper.ErrorPanic(err)

	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "customer_id not found",
		})
		return
	}

	controller.pjpService.Delete(c, id, customerID)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Deleted successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
