package model

import "time"

const (
	ORDER_APPROVAL_REQUEST_APPROVED = 1
	ORDER_APPROVAL_REQUEST_REJECTED = 9
)

type OrderApprovalRead struct {
	OrderApprovalRequestsDetailID int64      `gorm:"column:order_approval_request_id" json:"order_approval_request_id"`
	RoNo                          string     `gorm:"column:ro_no" json:"ro_no"`
	RoDate                        *time.Time `gorm:"column:ro_date" json:"ro_date"`
	OutletID                      int64      `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode                    string     `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName                    string     `gorm:"column:outlet_name" json:"outlet_name"`
	Total                         float64    `gorm:"column:total" json:"total"`
	SalesName                     string     `gorm:"column:sales_name" json:"sales_name"`
	CreditLimit                   float64    `gorm:"column:credit_limit" json:"credit_limit"`
	SalesInvLimit                 float64    `gorm:"column:sales_inv_limit" json:"sales_inv_limit"`
	CreditLimitType               *int       `gorm:"column:credit_limit_type" json:"credit_limit_type"`
	CreditLimitTypeName           string     `gorm:"column:credit_limit_type_name" json:"credit_limit_type_name"`
	CreditLimitAction             *int       `gorm:"column:credit_limit_action" json:"credit_limit_action"`
	CreditLimitActionName         string     `gorm:"column:credit_limit_action_name" json:"credit_limit_action_name"`
	CreditLimitValue              float64    `gorm:"column:validate_credit_limit_value" json:"validate_credit_limit_value"`
	SalesInvLimitType             *int       `gorm:"column:sales_inv_limit_type" json:"sales_inv_limit_type"`
	SalesInvLimitTypeName         string     `gorm:"column:sales_inv_limit_type_name" json:"sales_inv_limit_type_name"`
	SalesInvLimitAction           *int       `gorm:"column:sales_inv_limit_action" json:"sales_inv_limit_action"`
	SalesInvLimitActionName       string     `gorm:"column:sales_inv_limit_action_name" json:"sales_inv_limit_action_name"`
	SalesInvLimitValue            int        `gorm:"column:validate_overdue_value" json:"validate_overdue_value"`
	ObsType                       *int       `gorm:"column:obs_type" json:"obs_type"`
	ObsTypeName                   string     `gorm:"column:obs_type_name" json:"obs_type_name"`
	ObsLimitAction                *int       `gorm:"column:obs_limit_action" json:"obs_limit_action"`
	ObsLimitActionName            string     `gorm:"column:obs_limit_action_name" json:"obs_limit_action_name"`
	ObsLimitValue                 int        `gorm:"column:validate_outstanding_value" json:"validate_outstanding_value"`
	CustIDOrigin                  string     `gorm:"column:cust_id_origin" json:"cust_id_origin"`
}

type OrderApprovalActiveRead struct {
	OrderApprovalRequestsDetailID int64  `gorm:"column:order_approval_request_id" json:"order_approval_request_id"`
	RoNo                          string `gorm:"column:ro_no" json:"ro_no"`
}
