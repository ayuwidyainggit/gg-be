package model

import (
	"time"

	"gorm.io/gorm"
)

type SmcMCustomer struct {
	CustId         string         `gorm:"column:cust_id" json:"cust_id"`
	CustName       string         `gorm:"column:cust_name" json:"cust_name"`
	CompanyName    string         `gorm:"column:company_name" json:"company_name"`
	CompanyCode    string         `gorm:"column:company_code" json:"company_code"`
	Street1        string         `gorm:"column:street1" json:"street1"`
	Street2        string         `gorm:"column:street2" json:"street2"`
	City           string         `gorm:"column:city" json:"city"`
	StateID        int64          `gorm:"column:state_id" json:"state_id"`
	CountryID      int64          `gorm:"column:country_id" json:"country_id"`
	ZipCode        string         `gorm:"column:zip_code" json:"zip_code"`
	ContactName    string         `gorm:"column:contact_name" json:"contact_name"`
	ContactEmail   string         `gorm:"column:contact_email" json:"contact_email"`
	ContactPhoneNo string         `gorm:"column:contact_phone_no" json:"contact_phone_no"`
	Notes          string         `gorm:"column:notes" json:"notes"`
	ParentCustID   string         `gorm:"column:parent_cust_id" json:"parent_cust_id"`
	DistPriceGrpId int            `gorm:"column:dist_price_grp_id" json:"dist_price_grp_id"`
	Domain         string         `gorm:"column:domain" json:"domain"`
	DistributorID  *int8          `gorm:"column:distributor_id" json:"distributor_id"`
	IsActive       bool           `gorm:"column:is_active" json:"is_active"`
	CreatedBy      *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy      *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel          bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy      *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SmcMCustomer) TableName() string {
	return "smc.m_customer"
}
