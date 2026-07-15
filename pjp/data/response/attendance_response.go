package response

type AttendanceCheckResponse struct {
	Success     bool                 `json:"success"`
	Message     string               `json:"message"`
	Description string               `json:"description"`
	Data        *AttendanceCheckData `json:"data,omitempty"`
	RequestID   string               `json:"request_id"`
}

// AttendanceCheckData contains salesman info for distributor responses
// For Principal responses, only Plan field is used
type AttendanceCheckData struct {
	EmpID         int         `json:"emp_id,omitempty"`
	EmpCode       string      `json:"emp_code,omitempty"`
	EmpName       string      `json:"emp_name,omitempty"`
	OprType       string      `json:"opr_type,omitempty"`        // O for Taking Order
	OprTypeCanvas string      `json:"opr_type_canvas,omitempty"` // C for Canvas
	WhID          interface{} `json:"wh_id,omitempty"`           // int or empty string for TO
	WhCode        string      `json:"wh_code,omitempty"`
	WhNameCanvas  string      `json:"wh_name_canvas,omitempty"`
	Stock         interface{} `json:"stock,omitempty"` // int or empty string for TO
	Plan          int         `json:"plan"`
}
