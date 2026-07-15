package live_monitoring

import (
	livemonitoringservice "scyllax-pjp/service/live_monitoring"

	"github.com/gin-gonic/gin"
)

// LiveMonitoringController defines the interface for live monitoring endpoints
type LiveMonitoringController interface {
	GetPrincipalMonitoring(ctx *gin.Context)
	GetDistributorMonitoring(ctx *gin.Context)
	GetMonitoringDetail(ctx *gin.Context)
	GetUpdateLocations(ctx *gin.Context)
}

type liveMonitoringController struct {
	service livemonitoringservice.LiveMonitoringService
}

// NewLiveMonitoringController creates a new instance of LiveMonitoringController
func NewLiveMonitoringController(service livemonitoringservice.LiveMonitoringService) LiveMonitoringController {
	return &liveMonitoringController{
		service: service,
	}
}
