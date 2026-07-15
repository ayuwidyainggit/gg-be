package entity

import (
	"time"
)

type SupplierQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	IsActive *int   `query:"is_active"`
	SupType  string `query:"sup_type" validate:"max=1"`
}

type SupplierLookupScope struct {
	CustID          string 
	ParentCustID    string 
	DistributorIDs  []int  
	IncludeParentID bool   
}

type SupplierResponse struct {
	SupplierId      int        `json:"sup_id"`
	SupplierCode    string     `json:"sup_code"`
	SupplierName    string     `json:"sup_name"`
	Address1        string     `json:"address1"`
	Address2        string     `json:"address2"`
	City            string     `json:"city"`
	PhoneNo         string     `json:"phone_no"`
	FaxNo           *string    `json:"fax_no"`
	SupType         string     `json:"sup_type"`
	ContactName     *string    `json:"contact_name"`
	TaxName         *string    `json:"tax_name"`
	TaxNo           *string    `json:"tax_no"`
	PayTerm         *int       `json:"pay_term"`
	IsCreditLimit   bool       `json:"is_credit_limit"`
	CreditLimit     float64    `json:"credit_limit"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedByName   *string    `json:"updated_by_name"`
	UpdatedAt       *time.Time `json:"updated_at"`
	ProvinceId      *string    `json:"province_id"`
	ProvinceName    *string    `json:"province"`
	RegencyId       *string    `json:"regency_id"`
	RegencyName     *string    `json:"regency"`
	SubdustrictId   *string    `json:"sub_district_id"`
	SubdustrictName *string    `json:"sub_district"`
	WardId          *string    `json:"ward_id"`
	WardName        *string    `json:"ward"`
	ZipCode         *string    `json:"zip_code"`
	OtLocId         *int       `json:"ot_loc_id"`
	Latitude        *string    `json:"latitude"`
	Longitude       *string    `json:"longitude"`
	Email           *string    `json:"email"`
	IsWaNo          *bool      `json:"is_wa_no"`
	WaNo            *string    `json:"wa_no"`
	ContactType     *string    `json:"contact_type"`
	CreditLimitType *string    `json:"credit_limit_type"`
	Phone           *string    `json:"phone" `
	FaxNumber       *string    `json:"fax_number" `
	TaxIdentifierNo *string    `json:"tax_identifier_no" `
	Nitku           *string    `json:"nitku" `
	TaxAddress      *string    `json:"tax_address" `
}

type SupplierLookupResponse struct {
	SupplierId   int    `json:"sup_id"`
	SupplierCode string `json:"sup_code"`
	SupplierName string `json:"sup_name"`
	City         string `json:"city"`
	PhoneNo      string `json:"phone_no"`
	SupType      string `json:"sup_type"`
}

type CreateSupplierBody struct {
	CustId          string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy       int64   `json:"created_by" validate:"required"`
	DistributorId   *int64  `json:"distributor_id"`
	SupplierCode    string  `json:"sup_code" validate:"required,max=20,alphanum"`
	SupplierName    string  `json:"sup_name" validate:"required,max=50"`
	Address1        string  `json:"address1" validate:"required,max=100"`
	Address2        string  `json:"address2" validate:"max=100"`
	City            *string `json:"city"`
	PhoneNo         string  `json:"phone_no" validate:"required,max=15"`
	FaxNo           *string `json:"fax_no" validate:"max=15"`
	SupType         string  `json:"sup_type" validate:"required,max=1,oneof='I' 'E' 'C'"`
	ContactName     *string `json:"contact_name" validate:"max=100"`
	TaxName         *string `json:"tax_name,omitempty" validate:"omitempty,max=50"`
	TaxNo           *string `json:"tax_no,omitempty" validate:"omitempty,max=20"`
	PayTerm         *int    `json:"pay_term"`
	IsCreditLimit   bool    `json:"is_credit_limit"`
	CreditLimit     float64 `json:"credit_limit"`
	IsActive        bool    `json:"is_active"`
	UpdatedBy       *int64  `json:"updated_by"`
	ProvinceId      *string `json:"province_id"`
	RegencyId       *string `json:"regency_id"`
	SubDistrictId   *string `json:"sub_district_id"`
	WardId          *string `json:"ward_id"`
	ZipCode         *string `json:"zip_code" validate:"max=6"`
	OtLocId         *int    `json:"ot_loc_id"`
	Latitude        *string `json:"latitude"`
	Longitude       *string `json:"longitude"`
	Email           *string `json:"email" validate:"max=100"`
	IsWaNo          *bool   `json:"is_wa_no"`
	WaNo            *string `json:"wa_no" validate:"max=20"`
	ContactType     *string `json:"contact_type"`
	CreditLimitType *string `json:"credit_limit_type" validate:"max=1,oneof='L' 'U'"`
	Phone           *string `json:"phone" `
	FaxNumber       *string `json:"fax_number" `
	TaxIdentifierNo *string `json:"tax_identifier_no" `
	Nitku           *string `json:"nitku" `
	TaxAddress      *string `json:"tax_address" `
}

type DetailSupplierParams struct {
	SupplierId int `params:"sup_id" validate:"required"`
}

type UpdateSupplierParams struct {
	SupplierId int `params:"sup_id" validate:"required"`
}

type DeleteSupplierParams struct {
	SupplierId int `params:"sup_id" validate:"required"`
}

type UpdateSupplierRequest struct {
	CustId          string  `json:"cust_id" validate:"required,max=10"`
	UpdatedBy       int64   `json:"updated_by" validate:"required"`
	SupplierCode    string  `json:"sup_code,omitempty" validate:"max=20,omitempty,alphanum"`
	SupplierName    string  `json:"sup_name,omitempty" validate:"max=50,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
	Address1        string  `json:"address1" validate:"max=100"`
	Address2        string  `json:"address2" validate:"max=100"`
	City            *string `json:"city"`
	PhoneNo         string  `json:"phone_no" validate:"max=15"`
	FaxNo           string  `json:"fax_no" validate:"max=15"`
	SupType         string  `json:"sup_type" validate:"max=1,oneof='I' 'E' 'C'"`
	ContactName     string  `json:"contact_name" validate:"max=100"`
	TaxName         string  `json:"tax_name" validate:"max=50"`
	TaxNo           string  `json:"tax_no" validate:"max=20"`
	PayTerm         int     `json:"pay_term"`
	IsCreditLimit   *bool   `json:"is_credit_limit"`
	CreditLimit     float64 `json:"credit_limit"`
	ProvinceId      *string `json:"province_id"`
	RegencyId       *string `json:"regency_id"`
	SubDistrictId   *string `json:"sub_district_id"`
	WardId          *string `json:"ward_id"`
	ZipCode         *string `json:"zip_code"`
	OtLocId         *int    `json:"ot_loc_id"`
	Latitude        *string `json:"latitude"`
	Longitude       *string `json:"longitude"`
	Email           *string `json:"email"`
	IsWaNo          *bool   `json:"is_wa_no"`
	WaNo            *string `json:"wa_no"`
	ContactType     *string `json:"contact_type"`
	CreditLimitType *string `json:"credit_limit_type" validate:"max=1,oneof='L' 'U'"`
	Phone           *string `json:"phone" `
	FaxNumber       *string `json:"fax_number" `
	TaxIdentifierNo *string `json:"tax_identifier_no" `
	Nitku           *string `json:"nitku" `
	TaxAddress      *string `json:"tax_address" `
}
