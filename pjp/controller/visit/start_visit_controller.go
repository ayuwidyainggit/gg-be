package visit

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
// @Summary			start visit mobile.
// @Description		start visit.
// @Param			data body request.StartVisitRequest true "start visit"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/visits/start [post]
func (controller *visitController) StartVisit(ctx *gin.Context) {
	log.Info().Msg("start visit")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.StartVisitRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)
	empId, _ := helper.GetCurrentEmpId(ctx)
	Request.EmpID = int(empId)

	controller.visitService.StartVisit(c, Request, Request.CustID)
	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Start Visit successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}
