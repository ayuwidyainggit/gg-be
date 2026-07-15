package thirdparty

import (
	thirdparty "scyllax-pjp/service/third_party"

	"github.com/gin-gonic/gin"
)

type ThirdPartyController interface {
	GetAssignedSalesman(ctx *gin.Context)
	GetUnassignedSalesman(ctx *gin.Context)
	GetOutlet(ctx *gin.Context)
	GetDistributor(ctx *gin.Context)
	GetSalesmanByID(ctx *gin.Context)
}
type thirdPartyController struct {
	master thirdparty.ThirdPartyService
}

func NewThirdPartyController(master thirdparty.ThirdPartyService) ThirdPartyController {
	return &thirdPartyController{
		master: master,
	}
}
