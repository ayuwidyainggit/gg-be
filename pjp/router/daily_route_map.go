package router

import (
	"scyllax-pjp/controller"
	"scyllax-pjp/controller/daily"

	"github.com/gin-gonic/gin"
)

func RegisterDailyRouteRoutes(rg *gin.RouterGroup, ctrl *controller.DailyRouteMapController, ctrlNew daily.DailyRouteMapController) {
	dr := rg.Group("/daily-route-maps")
	dr.POST("", ctrlNew.SaveDailyRouteMap)

	dr.GET("", ctrl.FindAllPermanent)
	dr.GET("/daily", ctrl.FindAllDaily)
	dr.POST("/to/daily", ctrl.CopyRouteDailyToDaily)
	dr.POST("/all", ctrl.CopyAllToDaily)
	dr.POST("/partial", ctrl.CopySpecificToDaily)
	dr.POST("/specific", ctrl.CopySpecificToDaily)
}
