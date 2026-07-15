package visit

import (
	"scyllax-pjp/service/visit"

	"github.com/gin-gonic/gin"
)

type VisitController interface {
	StartVisit(ctx *gin.Context)
}
type visitController struct {
	visitService visit.VisitService
}

func NewVisitController(visitService visit.VisitService) VisitController {
	return &visitController{
		visitService: visitService,
	}
}
