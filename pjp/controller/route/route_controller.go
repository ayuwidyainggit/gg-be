package route

import (
	"scyllax-pjp/service/route"

	"github.com/gin-gonic/gin"
)

type RouteController interface {
	UpdateStatus(ctx *gin.Context)
	FindAllApproval(ctx *gin.Context)
	FindAllApprovalEnhance(ctx *gin.Context)
	UpdateStatusEnhance(ctx *gin.Context)
}
type routeController struct {
	routeService route.RouteService
}

func NewRouteController(routeService route.RouteService) RouteController {
	return &routeController{
		routeService: routeService,
	}
}
