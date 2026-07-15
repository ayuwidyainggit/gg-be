package pjp

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

// FindByIdTags 		godoc
// @Summary				Get Single pjp by id.
// @Param				pjpId path string true "update pjp by id"
// @Description			Return the pjp whoes pjpId value mathes id.
// @Produce				application/json
// @Tags				pjp
// @Success				200 {object} response.PjpResponse{}
// @Router				/web/pjp/{pjpId} [get]
// @Security            Bearer
func (controller *pjpController) GetById(ctx *gin.Context) {
	log.Info().Msg("find by id pjp")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	ParamID := ctx.Param("pjpId")
	id, err := strconv.Atoi(ParamID)
	helper.ErrorPanic(err)

	data := controller.pjpService.GetById(c, id, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
