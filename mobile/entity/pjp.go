package entity

type PjpSalesmanQueryParams struct {
	EmpId int64 `query:"emp_id" validate:"required"`
}

type PjpSalesmanWarehouseData struct {
	OprType string `json:"opr_type"`
	WhId    int64  `json:"wh_id"`
	WhCode  string `json:"wh_code"`
	WhName  string `json:"wh_name"`
}

type PjpSalesmanResponse struct {
	CustId          string         `json:"cust_id"`
	CustName        string         `json:"cust_name"`
	DistributorId   *int           `json:"distributor_id"` // null if principal, not null if distributor
	DistributorCode string         `json:"distributor_code"`
	DistributorName string         `json:"distributor_name"`
	EmpId           int64          `json:"emp_id"`
	SalesName       string         `json:"sales_name"`
	SalesTeamId     int64          `json:"sales_team_id"`
	SalesTeamCode   string         `json:"sales_team_code"`
	SalesTeamName   string         `json:"sales_team_name"`
	Data            map[string]any `json:"data"`
	PJPCode         *int64         `json:"pjp_code"`
	PJPStatus       *string        `json:"pjp_status"`
}
