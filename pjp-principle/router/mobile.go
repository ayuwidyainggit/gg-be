package router

import (
	"scyllax-pjp/controller"
	thirdparty "scyllax-pjp/controller/third_party"
	"scyllax-pjp/controller/visit"
	"scyllax-pjp/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterMobileRoutes(rg *gin.RouterGroup, visit *controller.VisitController, visitNew visit.VisitController, dms thirdparty.ThirdPartyController, daily *controller.DailyRouteMapController) {
	m := rg.Group("/mobile")
	m.POST("/visits/start", visitNew.StartVisit)

	m.GET("/visits/outlet/:salesCode/:custId/:date/:routeCode", visit.GetOutletBySalesCode)
	m.POST("/visits/end", visit.EndVisit)
	m.POST("/visits/arrive", visit.ArriveVisit)
	m.POST("/visits/skip", visit.SkipVisit)
	m.POST("/visits/resume", visit.ResumeVisit)
	m.POST("/visits/leave", visit.LeaveVisit)
	m.POST("/visits/onhold", visit.OnholdVisit)
	m.POST("/visits/outlet", visit.OutletVisit)
	m.GET("/visits/summary", visit.SummaryVisit)
	m.GET("/visits/status", visit.SummaryVisitStatus)
	m.GET("/todo/list/:outletVisitId", visit.TravelList)
	m.GET("/visits/outlet/list", visit.GetOutletVisitList)
	m.GET("/salesman/report", visit.GetSalesmanReport)
	// m.GET("/outlets/salesman", middleware.JwtMiddleware, dms.MobileGetOutletsBySalesman)
	m.POST("/add-outlets", middleware.JwtMiddleware, daily.MobileAddOutletToRoute)
	m.POST("/cancel/add-outlets", middleware.JwtMiddleware, daily.MobileCancelAddOutletToRoute)
}
