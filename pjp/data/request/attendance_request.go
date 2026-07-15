package request

type AttendanceCheckRequest struct {
	Date          int64 `form:"date" validate:"required"` // epoch timestamp
	EmpID         int   `form:"emp_id" validate:"required"`
	DistributorID *int  `form:"distributor_id"` // nullable
}
