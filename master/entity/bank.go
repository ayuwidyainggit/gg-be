package entity

import (
	"time"
)

type BankResponse struct {
	BankId        int        `json:"bank_id"`
	BankCode      string     `json:"bank_code"`
	BankName      string     `json:"bank_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type BankLookupResponse struct {
	BankId   int    `json:"bank_id"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

type CreateBankBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	BankCode  string `json:"bank_code" validate:"required,max=3,numeric"`
	BankName  string `json:"bank_name" validate:"required,max=25,alphanumericSpace"`
	IsActive  bool   `json:"is_active"`
}

type DetailBankParams struct {
	BankId int `params:"bank_id" validate:"required"`
}

type UpdateBankParams struct {
	BankId int `params:"bank_id" validate:"required"`
}

type DeleteBankParams struct {
	BankId int `params:"bank_id" validate:"required"`
}

type UpdateBankRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	BankCode  string `json:"bank_code,omitempty" validate:"required,max=3,numeric"`
	BankName  string `json:"bank_name,omitempty" validate:"max=25,alphanumericSpace,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

type QueryFilterOutletBank struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	OutletID     *int   `query:"outlet_id"`
	BankID       *int   `query:"bank_id"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type OutletBankList struct {
	BankId      int64  `db:"bank_id" json:"bank_id"`
	AccountNo   string `json:"account_no"`
	AccountName string `json:"account_name"`
}
