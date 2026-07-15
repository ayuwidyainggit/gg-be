package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Consignment struct {
	CustID     string          `gorm:"column:cust_id" json:"cust_id"`
	ConsNo     string          `gorm:"column:cons_no" json:"cons_no"`
	ConsDate   *time.Time      `gorm:"column:cons_date" json:"cons_date"`
	ConsType   *int64          `gorm:"column:cons_type" json:"cons_type"`
	OutletID   *int64          `gorm:"column:outlet_id" json:"outlet_id"`
	SalesmanID *int64          `gorm:"column:salesman_id" json:"salesman_id"`
	DeliveryNo *string         `gorm:"column:delivery_no" json:"delivery_no"`
	WhID       *int64          `gorm:"column:wh_id" json:"wh_id"`
	Notes      *string         `gorm:"column:notes" json:"notes"`
	SubTotal   *float64        `gorm:"column:sub_total" json:"sub_total"`
	Vat        *float64        `gorm:"column:vat" json:"vat"`
	VatValue   *float64        `gorm:"column:vat_value" json:"vat_value"`
	VatLg      *float64        `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue *float64        `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total      *float64        `gorm:"column:total" json:"total"`
	VatBg      *float64        `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue *float64        `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus *int64          `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time       `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool            `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Consignment) TableName() string {
	return "sls.consign"
}
func (m *Consignment) BeforeCreate(trx *gorm.DB) (err error) {
	intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	m.ConsNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type ConsignmentList struct {
	CustID        string          `gorm:"column:cust_id" json:"cust_id"`
	ConsNo        string          `gorm:"column:cons_no" json:"cons_no"`
	ConsDate      *time.Time      `gorm:"column:cons_date" json:"cons_date"`
	ConsType      *int64          `gorm:"column:cons_type" json:"cons_type"`
	OutletID      *int64          `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    *string         `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string         `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanID    *int64          `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  *string         `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  *string         `gorm:"column:salesman_name" json:"salesman_name"`
	DeliveryNo    *string         `gorm:"column:delivery_no" json:"delivery_no"`
	WhID          *int64          `gorm:"column:wh_id" json:"wh_id"`
	Notes         *string         `gorm:"column:notes" json:"notes"`
	SubTotal      *float64        `gorm:"column:sub_total" json:"sub_total"`
	Vat           *float64        `gorm:"column:vat" json:"vat"`
	VatValue      *float64        `gorm:"column:vat_value" json:"vat_value"`
	VatLg         *float64        `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue    *float64        `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total         *float64        `gorm:"column:total" json:"total"`
	VatBg         *float64        `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue    *float64        `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus    *int64          `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ConsignmentList) TableName() string {
	return "sls.consign"
}
