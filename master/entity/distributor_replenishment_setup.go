package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ReplenishmentCreditLimitAction int

func (c *ReplenishmentCreditLimitAction) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "restricted":
			*c = ReplenishmentCreditLimitAction(1)
		case "unrestricted":
			*c = ReplenishmentCreditLimitAction(2)
		default:
			return fmt.Errorf("credit_limit_action must be 1, 2, Restricted, or Unrestricted")
		}
		return nil
	}
	var n int
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	if n != 1 && n != 2 {
		return fmt.Errorf("credit_limit_action must be 1 or 2")
	}
	*c = ReplenishmentCreditLimitAction(n)
	return nil
}

// DistributorReplenishmentSetupQueryFilter query params for list API
type DistributorReplenishmentSetupQueryFilter struct {
	CustId         string `query:"-"`
	ParentCustID   string `query:"-"`
	Page           int    `query:"page"`
	Limit          int    `query:"limit"`
	Q              string `query:"q"`
	Sort           string `query:"sort"`
	DistributorIDs []int  `query:"distributor_id"`
	SupplierIDs    []int  `query:"supplier_id"`
}

// DistributorReplenishmentSetupListItem response row
type DistributorReplenishmentSetupListItem struct {
	ID                 int     `json:"id"`
	SupID              int     `json:"sup_id"`
	SupCode            string  `json:"sup_code"`
	SupName            string  `json:"sup_name"`
	DistributorID      int     `json:"distributor_id"`
	DistributorCode    string  `json:"distributor_code"`
	DistributorName    string  `json:"distributor_name"`
	DistributorType    string  `json:"distributor_type"`
	WhLimitAction      *string `json:"wh_limit_action,omitempty"`
	WhCapacity         *int    `json:"wh_capacity,omitempty"`
	WhVolume           *int    `json:"wh_volume,omitempty"`
	CreditLimitAction  string  `json:"credit_limit_action"`
	PlafonCredit       *int    `json:"plafon_credit,omitempty"`
	LeadTimeDays       int     `json:"lead_time_days"`
	IsApprovalRequired bool    `json:"is_approval_required"`
	CreatedBy          int64   `json:"created_by"`
	CreatedByName      string  `json:"created_by_name"`
	CreatedAt          string  `json:"created_at"`
	UpdatedBy          int64   `json:"updated_by"`
	UpdatedByName      string  `json:"updated_by_name"`
	UpdatedAt          string  `json:"updated_at"`
}

// DistributorReplenishmentSetupPaging paging block with request_id (per API spec)
type DistributorReplenishmentSetupPaging struct {
	TotalRecord int    `json:"total_record"`
	PageCurrent int    `json:"page_current"`
	PageLimit   int    `json:"page_limit"`
	PageTotal   int    `json:"page_total"`
	RequestID   string `json:"request_id"`
}

// DistributorReplenishmentSetupDetailParams path params for detail API
type DistributorReplenishmentSetupDetailParams struct {
	ID int `params:"id" validate:"required"`
}

type DistributorReplenishmentSetupPicParams struct {
	UserID int `params:"user_id" validate:"required"`
}

// DistributorReplenishmentSetupDetailResponse detail payload (single object + approvals)
type DistributorReplenishmentSetupDetailResponse struct {
	ID                 int     `json:"id"`
	SupID              int     `json:"sup_id"`
	SupCode            string  `json:"sup_code"`
	SupName            string  `json:"sup_name"`
	DistributorID      int     `json:"distributor_id"`
	DistributorCode    string  `json:"distributor_code"`
	DistributorName    string  `json:"distributor_name"`
	DistributorType    string  `json:"distributor_type"`
	WhLimitAction      *string `json:"wh_limit_action,omitempty"`
	WhCapacity         *int    `json:"wh_capacity,omitempty"`
	WhVolume           *int    `json:"wh_volume,omitempty"`
	CreditLimitAction  string  `json:"credit_limit_action"`
	PlafonCredit       *int    `json:"plafon_credit,omitempty"`
	LeadTimeDays       int     `json:"lead_time_days"`
	IsApprovalRequired bool    `json:"is_approval_required"`
	CreatedBy          int64   `json:"created_by"`
	CreatedByName      string  `json:"created_by_name"`
	CreatedAt          string  `json:"created_at"`
	UpdatedBy          int64   `json:"updated_by"`
	UpdatedByName      string  `json:"updated_by_name"`
	UpdatedAt          string  `json:"updated_at"`
	ApprovalData       []DistributorReplenishmentApprovalItem `json:"approval_data"`
}

// DistributorReplenishmentApprovalItem row in approval_data[]
type DistributorReplenishmentApprovalItem struct {
	ID                       int    `json:"id"`
	DistReplenishmentSetupID int    `json:"dist_replenishment_setup_id"`
	Level                    int    `json:"level"`
	Sequence                 int    `json:"sequence"`
	BusinessUnit             int    `json:"business_unit"`
	BusinessUnitName         string `json:"business_unit_name"`
	Pic                      int    `json:"pic"`
	PicName                  string `json:"pic_name"`
	IsActive                 bool   `json:"is_active"`
}

type DistributorReplenishmentSetupApprovalPayload struct {
	Level        int   `json:"level" validate:"required"`
	Sequence     int   `json:"sequence" validate:"required"`
	BusinessUnit int   `json:"business_unit" validate:"required"`
	Pic          int   `json:"pic" validate:"required"`
	IsActive     *bool `json:"is_active"`
}

type DistributorReplenishmentSetupCreatePayload struct {
	SupID              int                                            `json:"sup_id" validate:"required"`
	DistributorID      int                                            `json:"distributor_id" validate:"required"`
	DistributorType    string                                         `json:"distributor_type" validate:"required"`
	WhLimitAction      *string                                        `json:"wh_limit_action" validate:"omitempty,oneof=Restricted Unrestricted"`
	WhCapacity         *int                                           `json:"wh_capacity"`
	WhVolume           *int                                           `json:"wh_volume"`
	CreditLimitAction  ReplenishmentCreditLimitAction                 `json:"credit_limit_action" validate:"required"`
	PlafonCredit       *int                                           `json:"plafon_credit"`
	LeadTimeDays       int                                            `json:"lead_time_days" validate:"required"`
	IsApprovalRequired *bool                                          `json:"is_approval_required"`
	ApprovalData       []DistributorReplenishmentSetupApprovalPayload `json:"approval_data"`
}

type DistributorReplenishmentSupplierQueryFilter struct {
	CustId        string `query:"-"`
	ParentCustID  string `query:"-"`
	Pic           int    `query:"pic" validate:"required"`
	DistributorID *int   `query:"distributor_id"`
	Q             string `query:"q"`
	Page          int    `query:"page"`
	Limit         int    `query:"limit"`
	Sort          string `query:"sort"`
}

type DistributorReplenishmentSupplierItem struct {
	ID      int    `json:"id"`
	SupID   int    `json:"sup_id"`
	SupCode string `json:"sup_code"`
	SupName string `json:"sup_name"`
}

type DistributorReplenishmentDistributorQueryFilter struct {
	CustId       string `query:"-"`
	ParentCustID string `query:"-"`
	Pic          int    `query:"pic" validate:"required"`
	Q            string `query:"q"`
	Page         int    `query:"page" validate:"required"`
	Limit        int    `query:"limit" validate:"required"`
	Sort         string `query:"sort" validate:"required"`
}

type DistributorReplenishmentDistributorItem struct {
	ID              int    `json:"id"`
	DistributorID   int    `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
}
