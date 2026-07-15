package router

import (
	"github.com/gin-gonic/gin"
	"scyllax-pjp/controller"
)

func RegisterVisitDayRoutes(rg *gin.RouterGroup, ctrl *controller.VisitDayMapController) {
	vd := rg.Group("/visit-day-maps")
	vd.POST("", ctrl.SaveWeekly)
	vd.GET("", ctrl.GetAllVisitDayMap)
}
