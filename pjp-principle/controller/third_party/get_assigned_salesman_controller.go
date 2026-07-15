package thirdparty

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get list salesman by PJP.
// @Description		Return list of salesman.
// @Produce		    application/json
// @Tags			sales
// @Success         200 {object} response.Response{}
// @Router			/list-salesman [get]
// @Security	   Bearer
func (controller *thirdPartyController) GetAssignedSalesman(ctx *gin.Context) {
	log.Info().Msg("get salesman")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	headers := make(map[string]string)
	headers["Accept"] = ctx.GetHeader("application/json")
	headers["Authorization"] = ctx.GetHeader("Authorization")

	responses := controller.master.GetAssignedSalesman(c, headers, currentCustomerId)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, responses)
}
