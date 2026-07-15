package pjpenhance

import (
	"scyllax-pjp/service"
	pjpenhance "scyllax-pjp/service/pjp_enhance"

	"github.com/gin-gonic/gin"
)

type PjpEnhanceController interface {
	Create(ctx *gin.Context)
	GetPjpById(ctx *gin.Context)
	UpdatePjpById(ctx *gin.Context)
	UpdateStatusPjpById(ctx *gin.Context)
	UpdateStatusPjpByEmpId(ctx *gin.Context)
}

type pjpEnhanceController struct {
	pjpEnhanceService pjpenhance.PjpEnhanceService
	routePopService   service.RoutePopService
}

func NewPjpEnhanceController(pjpEnhanceService pjpenhance.PjpEnhanceService, popService service.RoutePopService) PjpEnhanceController {
	return &pjpEnhanceController{
		pjpEnhanceService: pjpEnhanceService,
		routePopService:   popService,
	}
}
