package daily

import (
	routepop "scyllax-pjp/service/route_pop"

	"github.com/gin-gonic/gin"
)

type DailyRouteMapController interface {
	SaveDailyRouteMap(ctx *gin.Context)
}

type dailyRouteMapController struct {
	routePopService routepop.RoutePopService
}

func NewDailyRouteMapController(routePopService routepop.RoutePopService) *dailyRouteMapController {
	return &dailyRouteMapController{
		routePopService: routePopService,
	}
}
