package pjp

import (
	"scyllax-pjp/service/pjp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PjpController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetPjpWithRoute(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	ListPjpApprove(ctx *gin.Context)
	GetById(ctx *gin.Context)
}

type pjpController struct {
	pjpService pjp.PjpService
}

func NewPjpController(service pjp.PjpService) PjpController {
	return &pjpController{
		pjpService: service,
	}
}

// extractFilters parses all supported filters from query params
func extractFilters(ctx *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})

	if v := ctx.Query("pjp_code"); v != "" {
		filters["pjp_code"] = strings.Split(v, "%2C")
	}
	if v := ctx.Query("team_salesman"); v != "" {
		filters["team_salesman"] = strings.Split(v, "%2C")
	}

	filters["operation_type"] = ctx.Query("operation_type")
	filters["salesman_name"] = ctx.Query("salesman_name")
	filters["salesman_code"] = ctx.Query("salesman_code")
	filters["status"] = ctx.Query("is_active")
	filters["q"] = ctx.Query("q")

	return filters
}

// parseQueryInt tries to parse an int query, returns default if invalid
func parseQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	if val, err := strconv.Atoi(ctx.Query(key)); err == nil && val > 0 {
		return val
	}
	return defaultVal
}
