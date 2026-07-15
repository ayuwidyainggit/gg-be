package router

import (
	"scyllax-pjp/controller/pjp"
	pjpauto "scyllax-pjp/controller/pjp_auto"
	pjpenhance "scyllax-pjp/controller/pjp_enhance"

	"github.com/gin-gonic/gin"
)

func RegisterPjpRoutes(rg *gin.RouterGroup, pjp pjp.PjpController, enhance pjpenhance.PjpEnhanceController, auto pjpauto.PjpAutoController) {
	pjpRouter := rg.Group("/pjp")

	pjpRouter.POST("", pjp.Create)
	pjpRouter.GET("", pjp.GetAll)
	pjpRouter.GET("/list", pjp.GetPjpWithRoute)
	pjpRouter.PATCH("/:pjpId", pjp.Update)
	pjpRouter.DELETE("/:pjpId", pjp.Delete)
	pjpRouter.GET("/visit-list", pjp.ListPjpApprove)

	pjpRouter.GET("/:pjpId", pjp.GetById)
	// pjpRouter.PATCH("/update/:empId", enhance.UpdateStatusPjpByEmpId)

	pjpAuto := pjpRouter.Group("/auto")
	pjpAuto.POST("", auto.CreatePjpAuto)

	enhanceRouter := rg.Group("/pjp-enhance")
	enhanceRouter.POST("", enhance.Create)
	enhanceRouter.GET("/:id", enhance.GetPjpById)
	enhanceRouter.PUT("/:id", enhance.UpdatePjpById)
	enhanceRouter.PUT("/:id/status", enhance.UpdateStatusPjpById)
}
