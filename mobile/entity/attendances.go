package entity

import "time"

var (
	TYPE_CHECKIN  string = "checkin"
	TYPE_CHECKOUT string = "checkout"

	TYPE_CHECKIN_ID  int = 1
	TYPE_CHECKOUT_ID int = 2
)

type AttendanceRequest struct {
	Latitude  string  `json:"latitude" validate:"required,latitude"`
	Longitude string  `json:"longitude" validate:"required,longitude"`
	Type      *string `json:"type" validate:"required,oneof=checkin checkout"`
	LeaveID   *int64  `json:"leave_id" `
	Email     string
	CustID    string
}

func ConvTypeStringToCode(typeAtt string) int {
	if typeAtt == TYPE_CHECKIN {
		return TYPE_CHECKIN_ID
	} else {
		return TYPE_CHECKOUT_ID
	}
}

type AttendanceResponse struct {
	CurrentTime time.Time  `json:"current_time"`
	CheckIn     *time.Time `json:"checkin_at"`
	CheckOut    *time.Time `json:"checkout_at"`
}
type AttendanceGetRequest struct {
	Email  string
	CustID string
}

type AttendanceCheckRequest struct {
	Date          int64  `query:"date" validate:"required"`
	EmpID         int64  `query:"emp_id" validate:"required"`
	DistributorID *int64 `query:"distributor_id"` // Optional: NULL for principal, NOT NULL for distributor
	CustID        string
}

type AttendanceCheckResponse struct {
	Success     bool                        `json:"success"`
	Message     string                      `json:"message"`
	Description string                      `json:"description"`
	Data        AttendanceCheckResponseData `json:"data"`
}

type AttendanceCheckResponseData struct {
	Plan int `json:"plan"`

	EmpID         *int64  `json:"emp_id,omitempty"`
	EmpCode       *string `json:"emp_code,omitempty"`
	EmpName       *string `json:"emp_name,omitempty"`
	OprType       *string `json:"opr_type,omitempty"`        // "O" for Taking Order
	OprTypeCanvas *string `json:"opr_type_canvas,omitempty"` // "C" for Canvas
	WhID          *int64  `json:"wh_id,omitempty"`
	WhCode        *string `json:"wh_code,omitempty"`
	WhNameCanvas  *string `json:"wh_name_canvas,omitempty"`
	Stock         *int    `json:"stock,omitempty"`
}
