package router

import (
	thirdparty "scyllax-pjp/controller/third_party"
	"scyllax-pjp/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterThirdPartyRoutes(rg *gin.RouterGroup, master thirdparty.ThirdPartyController) {
	// rg.GET("/teams", middleware.JwtMiddleware, dms.GetSalesTeam)
	rg.GET("/teams/salesman", middleware.JwtMiddleware, master.GetUnassignedSalesman)
	rg.GET("/teams/salesman/:empId", middleware.JwtMiddleware, master.GetSalesmanByID)
	// rg.GET("/teams/operation/type", middleware.JwtMiddleware, dms.GetSalesOperationType)
	// rg.GET("/teams/salesman/type", dms.GetSalesTeamType)
	// rg.GET("/warehouses", middleware.JwtMiddleware, dms.GetWarehouse)
	rg.GET("/outlets", middleware.JwtMiddleware, master.GetOutlet)
	// rg.GET("/outlet/list", middleware.JwtMiddleware, dms.GetListOutlet)
	// rg.GET("/outlets/not-assign", dms.GetOutletNotAssign)
	rg.GET("/outlets/salesman", middleware.JwtMiddleware, master.GetOutletBySalesCodes)
	rg.GET("/outlets-picklist/salesman", middleware.JwtMiddleware, master.GetOutletPicklistBySalesCodes)
	rg.GET("/list-salesman", middleware.JwtMiddleware, master.GetAssignedSalesman)
}
