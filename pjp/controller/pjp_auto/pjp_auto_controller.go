package pjpauto

import (
	pjpauto "scyllax-pjp/service/pjp_auto"

	"github.com/gin-gonic/gin"
)

type PjpAutoController interface {
	CreatePjpAuto(ctx *gin.Context)
}

type pjpAutoController struct {
	pjpAutoService pjpauto.PjpAutoService
}

func NewPjpAutoController(pjpAutoService pjpauto.PjpAutoService) PjpAutoController {
	return &pjpAutoController{
		pjpAutoService: pjpAutoService,
	}
}
