package entity

type ExtraCallOutlet struct {
	PJPID         int    `json:"pjp_id"`
	PJPCode       string `json:"pjp_code"`
	Date          string `json:"date"`
	RouteCode     string `json:"route_code"`
	OutletIDs     []int  `json:"outlet_id" validate:"dive,required"`
	CustID        string `json:"-"`
	ParentCustID  string `json:"-"`
	EmpID         int64  `json:"-"`
	IsDistributor bool   `json:"-"`
}

type CreateExtraCallRequest struct {
	DestinationIDs  []int  `json:"destination_id" validate:"required"`
	DestinationType string `json:"destination_type" validate:"required"`
	Start           int64  `json:"start"`
	Date            string `json:"date" validate:"required"`
	RouteCode       string `json:"route_code"`
	PJPCode         string `json:"pjp_code"`
	PJPID           int64  `json:"pjp_id"`
	EmpID           int
	CustID          string
	IsDistributor   bool
}

type CreateExtraCallResponse struct {
}
