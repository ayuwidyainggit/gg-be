package model

import (
	"time"
)

type Supplier struct {
	CustId        string  `db:"cust_id" json:"cust_id"`
	DistributorId *int64  `db:"distributor_id" json:"distributor_id"`
	SupplierId    int     `db:"sup_id" json:"sup_id"`
	SupplierCode  string  `db:"sup_code" json:"sup_code"`
	SupplierName  string  `db:"sup_name" json:"sup_name"`
	Address1      string  `db:"address1" json:"address1"`
	Address2      *string `db:"address2" json:"address2"`
	City          string  `db:"city" json:"city"`
	PhoneNo       string  `db:"phone_no" json:"phone_no"`
	FaxNo         *string `db:"fax_no" json:"fax_no"`
	SupType       string  `db:"sup_type" json:"sup_type"`
	ContactName   *string `db:"contact_name" json:"contact_name"`
	TaxName       *string `db:"tax_name" json:"tax_name"`
	TaxNo         *string `db:"tax_no" json:"tax_no"`
	PayTerm       *int    `db:"pay_term" json:"pay_term"`
	IsCreditLimit bool    `db:"is_credit_limit" json:"is_credit_limit"`
	CreditLimit   float64 `db:"credit_limit" json:"credit_limit"`

	ProvinceId      *string `db:"province_id" json:"province_id"`
	RegencyId       *string `db:"regency_id" json:"regency_id"`
	SubDistrictId   *string `db:"sub_district_id" json:"sub_district_id"`
	WardId          *string `db:"ward_id" json:"ward_id"`
	ZipCode         *string `db:"zip_code" json:"zip_code"`
	OtLocId         *int    `db:"ot_loc_id" json:"ot_loc_id"`
	Latitude        *string `db:"latitude" json:"latitude"`
	Longitude       *string `db:"longitude" json:"longitude"`
	Email           *string `db:"email" json:"email"`
	IsWaNo          *bool   `db:"is_wa_no" json:"is_wa_no"`
	WaNo            *string `db:"wa_no" json:"wa_no"`
	ContactType     *string `db:"contact_type" json:"contact_type"`
	CreditLimitType *string `db:"credit_limit_type" json:"credit_limit_type"`
	Phone           *string `db:"phone" json:"phone" `
	FaxNumber       *string `db:"fax_number" json:"fax_number" `
	TaxIdentifierNo *string `db:"tax_identifier_no" json:"tax_identifier_no" `
	Nitku           *string `db:"nitku" json:"nitku" `
	TaxAddress      *string `db:"tax_address" json:"tax_address" `

	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type SupplierRead struct {
	CustId          string  `db:"cust_id" json:"cust_id"`
	DistributorId   *int64  `db:"distributor_id" json:"distributor_id"`
	SupplierId      int     `db:"sup_id" json:"sup_id"`
	SupplierCode    string  `db:"sup_code" json:"sup_code"`
	SupplierName    string  `db:"sup_name" json:"sup_name"`
	Address1        string  `db:"address1" json:"address1"`
	Address2        *string `db:"address2" json:"address2"`
	City            string  `db:"city" json:"city"`
	PhoneNo         string  `db:"phone_no" json:"phone_no"`
	FaxNo           *string `db:"fax_no" json:"fax_no"`
	SupType         string  `db:"sup_type" json:"sup_type"`
	ContactName     *string `db:"contact_name" json:"contact_name"`
	TaxName         *string `db:"tax_name" json:"tax_name"`
	TaxNo           *string `db:"tax_no" json:"tax_no"`
	PayTerm         *int    `db:"pay_term" json:"pay_term"`
	IsCreditLimit   bool    `db:"is_credit_limit" json:"is_credit_limit"`
	CreditLimit     float64 `db:"credit_limit" json:"credit_limit"`
	ProvinceId      *string `db:"province_id" json:"province_id"`
	ProvinceName    *string `db:"province" json:"province"`
	RegencyId       *string `db:"regency_id" json:"regency_id"`
	RegencyName     *string `db:"regency" json:"regency"`
	SubDistrictId   *string `db:"sub_district_id" json:"sub_district_id"`
	SubDistrictName *string `db:"sub_district" json:"sub_district"`
	WardId          *string `db:"ward_id" json:"ward_id"`
	WardName        *string `db:"ward" json:"ward"`
	ZipCode         *string `db:"zip_code" json:"zip_code"`
	OtLocId         *int    `db:"ot_loc_id" json:"ot_loc_id"`
	Latitude        *string `db:"latitude" json:"latitude"`
	Longitude       *string `db:"longitude" json:"longitude"`
	Email           *string `db:"email" json:"email"`
	IsWaNo          *bool   `db:"is_wa_no" json:"is_wa_no"`
	WaNo            *string `db:"wa_no" json:"wa_no"`
	ContactType     *string `db:"contact_type" json:"contact_type"`
	CreditLimitType *string `db:"credit_limit_type" json:"credit_limit_type"`
	Phone           *string `db:"phone" json:"phone" `
	FaxNumber       *string `db:"fax_number" json:"fax_number" `
	TaxIdentifierNo *string `db:"tax_identifier_no" json:"tax_identifier_no" `
	Nitku           *string `db:"nitku" json:"nitku" `
	TaxAddress      *string `db:"tax_address" json:"tax_address" `

	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type SupplierUpdate struct {
	SupplierCode  *string  `json:"sup_code,omitempty" sql:"sup_code"`
	SupplierName  *string  `json:"sup_name,omitempty" sql:"sup_name"`
	Address1      *string  `json:"address1,omitempty" sql:"address1"`
	Address2      *string  `json:"address2,omitempty" sql:"address2"`
	City          *string  `json:"city,omitempty" sql:"city"`
	PhoneNo       *string  `json:"phone_no,omitempty" sql:"phone_no"`
	FaxNo         *string  `json:"fax_no,omitempty" sql:"fax_no"`
	SupType       *string  `json:"sup_type,omitempty" sql:"sup_type"`
	ContactName   *string  `json:"contact_name,omitempty" sql:"contact_name"`
	TaxName       *string  `json:"tax_name,omitempty" sql:"tax_name"`
	TaxNo         *string  `json:"tax_no,omitempty" sql:"tax_no"`
	PayTerm       *int     `json:"pay_term,omitempty" sql:"pay_term"`
	IsCreditLimit *bool    `json:"is_credit_limit" sql:"is_credit_limit"`
	CreditLimit   *float64 `json:"credit_limit,omitempty" sql:"credit_limit"`

	ProvinceId      *string `db:"province_id" json:"province_id"`
	RegencyId       *string `db:"regency_id" json:"regency_id"`
	SubDistrictId   *string `db:"sub_district_id" json:"sub_district_id"`
	WardId          *string `db:"ward_id" json:"ward_id"`
	ZipCode         *string `db:"zip_code" json:"zip_code"`
	OtLocId         *int    `db:"ot_loc_id" json:"ot_loc_id"`
	Latitude        *string `db:"latitude" json:"latitude"`
	Longitude       *string `db:"longitude" json:"longitude"`
	Email           *string `db:"email" json:"email"`
	IsWaNo          *bool   `db:"is_wa_no" json:"is_wa_no"`
	WaNo            *string `db:"wa_no" json:"wa_no"`
	ContactType     *string `db:"contact_type" json:"contact_type"`
	CreditLimitType *string `db:"credit_limit_type" json:"credit_limit_type"`
	Phone           *string `db:"phone" json:"phone" `
	FaxNumber       *string `db:"fax_number" json:"fax_number" `
	TaxIdentifierNo *string `db:"tax_identifier_no" json:"tax_identifier_no" `
	Nitku           *string `db:"nitku" json:"nitku" `
	TaxAddress      *string `db:"tax_address" json:"tax_address" `

	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
