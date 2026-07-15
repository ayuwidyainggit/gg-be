package router

import (
	"scyllax-pjp/controller/route"

	"github.com/gin-gonic/gin"
)

func RegisterApprovalRoutes(rg *gin.RouterGroup, ctrl route.RouteController) {
	approval := rg.Group("/approval-routes")
	approval.GET("", ctrl.FindAllApproval)
	approval.PATCH("/status", ctrl.UpdateStatus)

	enhance := rg.Group("/approval-routes-enhance")
	enhance.GET("", ctrl.FindAllApprovalEnhance)
	enhance.PATCH("/status", ctrl.UpdateStatusEnhance)
}
