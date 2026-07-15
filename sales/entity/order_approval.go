package entity

const (
	ORDER_APPROVAL_APPROVED = 1
	ORDER_APPROVAL_REJECTED = 9
)

type OrderApprovalQueryFilter struct {
	CustId       string
	ParentCustId string
	EmpID        int64
	SalesmanId   []int  `query:"salesman_id"`
	OutletID     []int  `query:"outlet_id"`
	RoFrom       *int64 `query:"ro_date_from" validate:"required_with=RoTo,omitempty,gte=1000000000"`
	RoTo         *int64 `query:"ro_date_to" validate:"required_with=RoFrom,omitempty,lte=9999999999,gtefield=RoFrom"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Status       int    `query:"status"`
	Sort         string `query:"sort"`
}

type OrderApprovalListResponse struct {
	OrderApprovalRequestsDetailID int64   `json:"order_approval_request_id"`
	RoNo                          string  `json:"ro_no"`
	RoDate                        string  `json:"ro_date"`
	OutletID                      int64   `json:"outlet_id"`
	OutletCode                    string  `json:"outlet_code"`
	OutletName                    string  `json:"outlet_name"`
	Total                         float64 `json:"total"`
	SalesName                     string  `json:"sales_name"`
	CreditLimit                   float64 `json:"credit_limit"`
	SalesInvLimit                 float64 `json:"sales_inv_limit"`
	CreditLimitType               *int    `json:"credit_limit_type"`
	CreditLimitTypeName           string  `json:"credit_limit_type_name"`
	CreditLimitAction             *int    `json:"credit_limit_action"`
	CreditLimitActionName         string  `json:"credit_limit_action_name"`
	CreditLimitValue              float64 `json:"credit_limit_value"`
	SalesInvLimitType             *int    `json:"sales_inv_limit_type"`
	SalesInvLimitTypeName         string  `json:"sales_inv_limit_type_name"`
	SalesInvLimitAction           *int    `json:"sales_inv_limit_action"`
	SalesInvLimitActionName       string  `json:"sales_inv_limit_action_name"`
	SalesInvLimitValue            int     `json:"sales_inv_limit_value"`
	ObsType                       *int    `json:"obs_type"`
	ObsTypeName                   string  `json:"obs_type_name"`
	ObsLimitAction                *int    `json:"obs_limit_action"`
	ObsLimitActionName            string  `json:"obs_limit_action_name"`
	ObsLimitValue                 int     `json:"obs_limit_value"`
	CustIDOrigin                  string  `json:"cust_id_origin"`
}

type UpdateOrderApprovalDetailParams struct {
	OrderApprovalRequestsDetailID int64 `params:"order_approval_request_id" validate:"required"`
}

type UpdateOrderApprovalDetailBody struct {
	EmpID  int64
	Status int `json:"status" validate:"oneof=1 9"`
}
