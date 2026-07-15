package router

import (
	"scyllax-pjp/controller"
	"scyllax-pjp/controller/attendance"
	"scyllax-pjp/controller/daily"
	"scyllax-pjp/controller/geotaging"
	livemonitoring "scyllax-pjp/controller/live_monitoring"
	"scyllax-pjp/controller/pjp"
	pjpauto "scyllax-pjp/controller/pjp_auto"
	pjpenhance "scyllax-pjp/controller/pjp_enhance"
	"scyllax-pjp/controller/route"
	thirdparty "scyllax-pjp/controller/third_party"
	"scyllax-pjp/controller/visit"
	"scyllax-pjp/exception"
	"scyllax-pjp/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	pjpController pjp.PjpController,
	pjpEnhanceController pjpenhance.PjpEnhanceController,
	pjpAutoController pjpauto.PjpAutoController,
	routeMappingController *controller.RouteMappingController,
	routeController route.RouteController,
	visitDayMapController *controller.VisitDayMapController,
	thirdPartyController thirdparty.ThirdPartyController,
	dailyRouteMapController *controller.DailyRouteMapController,
	dailyRouteMapControllerNew daily.DailyRouteMapController,
	visitController *controller.VisitController,
	visitControllerNew visit.VisitController,
	geotagingController geotaging.GeotagingController,
	attendanceController attendance.AttendanceController,
	liveMonitoringController livemonitoring.LiveMonitoringController,
) *gin.Engine {
	service := gin.Default()
	//add custom recovery
	service.Use(gin.CustomRecovery(exception.ErrorHandler))
	service.Use(middleware.RequestID())
	service.Use(gin.Logger())
	//add cors
	service.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//add swagger docs
	service.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})

	router := service.Group("/api/v1")
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 200, "message": "OK"})
	})

	// third party
	RegisterThirdPartyRoutes(router, thirdPartyController)

	// mobile
	RegisterMobileRoutes(router, visitController, visitControllerNew, thirdPartyController, dailyRouteMapController, geotagingController, attendanceController)

	//web
	webRouter := router.Group("/web")
	webRouter.Use(middleware.JwtMiddleware)

	// pjp, pjp_enhance, pjp_auto
	RegisterPjpRoutes(webRouter, pjpController, pjpEnhanceController, pjpAutoController)

	// route_mapping
	RegisterRouteMappingRoutes(webRouter, routeMappingController, visitDayMapController)

	// approval
	RegisterApprovalRoutes(webRouter, routeController)

	// visit day
	RegisterVisitDayRoutes(webRouter, visitDayMapController)

	// daily route map
	RegisterDailyRouteRoutes(webRouter, dailyRouteMapController, dailyRouteMapControllerNew)

	// live monitoring (with JWT middleware, directly under /api/v1)
	liveMonitoringRouter := router.Group("")
	liveMonitoringRouter.Use(middleware.JwtMiddleware)
	RegisterLiveMonitoringRoutes(liveMonitoringRouter, liveMonitoringController)

	return service
}
