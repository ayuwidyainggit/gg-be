package attendance

import (
	"context"
	"net/http"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/service/attendance"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	requestTimeout = 10 * time.Second
)

// AttendanceController defines contract for attendance HTTP handlers
type AttendanceController interface {
	CheckAttendance(ctx *gin.Context)
}

type attendanceController struct {
	service attendance.AttendanceService
}

// NewAttendanceController creates a new attendance controller instance
func NewAttendanceController(service attendance.AttendanceService) AttendanceController {
	return &attendanceController{
		service: service,
	}
}

// CheckAttendance godoc
// @Summary      Check-in validation for salesman
// @Description  Validates if check-in is available based on route plan (PJP) and warehouse stock.
// @Tags         mobile
// @Accept       json
// @Produce      json
// @Param        date            query     int     true   "Unix timestamp (epoch time)"
// @Param        emp_id          query     int     true   "Employee ID"
// @Param        distributor_id  query     int     false  "Distributor ID (NULL for Principal)"
// @Success      200  {object}  response.AttendanceCheckResponse
// @Router       /mobile/v1/attendances/check [get]
func (c *attendanceController) CheckAttendance(ctx *gin.Context) {
	req := request.AttendanceCheckRequest{}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		helper.ErrorPanic(err)
	}

	cxt, cancel := context.WithTimeout(ctx.Request.Context(), requestTimeout)
	defer cancel()

	res := c.service.CheckAttendance(cxt, req)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, res)
}
