package router

import (
	"scyllax-pjp/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRouteMappingRoutes(rg *gin.RouterGroup, route *controller.RouteMappingController, visitDay *controller.VisitDayMapController) {
	routeRouter := rg.Group("/route-mappings")

	routeRouter.GET("", route.FindAll)
	routeRouter.POST("", route.Create)
	routeRouter.POST("/outlets", route.SaveOutlet)
	routeRouter.POST("/pjp", route.SavePjp)
	routeRouter.PATCH("/outlets/update", route.DeleteOutlet)
	routeRouter.PATCH("/outlets-additional/update", route.DeleteOutletAdditional)
	routeRouter.PATCH("/pjp", route.UpdatePjp)
	routeRouter.GET("/:routeCode/:pjpCode", route.FindByRouteCode)
	routeRouter.GET("/additional/:routeCode", visitDay.FindByParentRoute)
	routeRouter.DELETE("/:routeId", route.Delete)
	routeRouter.GET("/pjp/:pjpCode/:routeCode", route.FindRouteByPjpCode)
	routeRouter.PATCH("/:routeId", route.Update)
	routeRouter.POST("/save/route", route.SaveRouteConfirmation)
	routeRouter.PATCH("/remove/route", route.RemoveRouteInPjp)
	routeRouter.PATCH("/new-route", route.NewRoutePropose)
	routeRouter.GET("/:routeCode", route.FindByRouteOutlet)
	routeRouter.POST("duplicate", route.RouteDuplicate)
}
