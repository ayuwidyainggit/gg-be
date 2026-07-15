package router

import (
	livemonitoring "scyllax-pjp/controller/live_monitoring"

	"github.com/gin-gonic/gin"
)

// RegisterLiveMonitoringRoutes registers all live monitoring routes
func RegisterLiveMonitoringRoutes(router *gin.RouterGroup, controller livemonitoring.LiveMonitoringController) {
	// Live monitoring routes (without /web prefix, directly under /api/v1)
	router.GET("/live-monitoring-principal", controller.GetPrincipalMonitoring)
	router.GET("/live-monitoring-distributor", controller.GetDistributorMonitoring)
	router.GET("/monitoring_locations/details", controller.GetMonitoringDetail)
	router.GET("/monitoring_locations/update-locations", controller.GetUpdateLocations)
}
