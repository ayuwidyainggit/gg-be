package request

type StartVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  int64  `json:"current_time" validate:"required"`
	// RouteCode    string `validate:"required" json:"route_code"`
}

type OutletVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
}

type FinishVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  int64  `json:"current_time" validate:"required"`
	// Id           int64  `validate:"required" json:"id"`
}

type SkipVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  int64  `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
	SkipReason   string `validate:"required" json:"skip_reason"`
}

type ResumeVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  *int64 `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
}

type ArriveVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  *int64 `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
}

type LeaveVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  *int64 `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
}

type OnholdVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  *int64 `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
}
