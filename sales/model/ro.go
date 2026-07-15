package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Ro struct {
	CustID        string          `gorm:"cust_id" json:"cust_id"`
	RoNo          string          `gorm:"ro_no" json:"ro_no"`
	RoDate        *time.Time      `gorm:"ro_date" json:"ro_date"`
	ValDate       *time.Time      `gorm:"val_date" json:"val_date"`
	DueDate       *time.Time      `gorm:"due_date" json:"due_date"`
	SalesmanId    *int64          `gorm:"salesman_id" json:"salesman_id"`
	WhId          *int64          `gorm:"wh_id" json:"wh_id"`
	OutletID      *int64          `gorm:"outlet_id" json:"outlet_id"`
	DeliveryDate  *time.Time      `gorm:"delivery_date" json:"delivery_date"`
	OrderNo       *string         `gorm:"order_no" json:"order_no"`
	PoNo          *string         `gorm:"po_no" json:"po_no"`
	VehicleNo     *string         `gorm:"vehicle_no" json:"vehicle_no"`
	PayType       *int64          `gorm:"pay_type" json:"pay_type"`
	ReffNo        *string         `gorm:"reff_no" json:"reff_no"`
	MobileID      *int64          `gorm:"mobile_id" json:"mobile_id"`
	SubTotal      *float64        `gorm:"sub_total" json:"sub_total"`
	Disc          *float64        `gorm:"disc" json:"disc"`
	DiscValue     *float64        `gorm:"disc_value" json:"disc_value"`
	PromoValue    *float64        `gorm:"promo_value" json:"promo_value"`
	CashDiscValue *float64        `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1      *float64        `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2      *float64        `gorm:"tot_disc2" json:"tot_disc2"`
	Vat           *float64        `gorm:"vat" json:"vat"`
	VatValue      *float64        `gorm:"vat_value" json:"vat_value"`
	Total         *float64        `gorm:"total" json:"total"`
	DataStatus    *int64          `gorm:"data_status" json:"data_status"`
	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
	DataSource    *int64          `gorm:"data_source" json:"data_source"`
}

func (Ro) TableName() string {
	return "sls.ro"
}

func (m *Ro) BeforeCreate(trx *gorm.DB) (err error) {
	intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	m.RoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type RoList struct {
	CustID        string          `gorm:"cust_id" json:"cust_id"`
	RoNo          string          `gorm:"ro_no" json:"ro_no"`
	RoDate        *time.Time      `gorm:"ro_date" json:"ro_date"`
	ValDate       *time.Time      `gorm:"val_date" json:"val_date"`
	DueDate       *time.Time      `gorm:"due_date" json:"due_date"`
	SalesmanId    *int64          `gorm:"salesman_id" json:"salesman_id"`
	SalesName     *string         `gorm:"sales_name" json:"sales_name"`
	WhId          *int64          `gorm:"wh_id" json:"wh_id"`
	WhCode        *string         `gorm:"wh_code" json:"wh_code"`
	WhName        *string         `gorm:"wh_name" json:"wh_name"`
	OutletID      *int64          `gorm:"outlet_id" json:"outlet_id"`
	OutletCode    *string         `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string         `gorm:"column:outlet_name" json:"outlet_name"`
	DeliveryDate  *time.Time      `gorm:"delivery_date" json:"delivery_date"`
	OrderNo       *string         `gorm:"order_no" json:"order_no"`
	PoNo          *string         `gorm:"po_no" json:"po_no"`
	VehicleNo     *string         `gorm:"vehicle_no" json:"vehicle_no"`
	PayType       *int64          `gorm:"pay_type" json:"pay_type"`
	ReffNo        *string         `gorm:"reff_no" json:"reff_no"`
	MobileID      *int64          `gorm:"mobile_id" json:"mobile_id"`
	SubTotal      *float64        `gorm:"sub_total" json:"sub_total"`
	Disc          *float64        `gorm:"disc" json:"disc"`
	DiscValue     *float64        `gorm:"disc_value" json:"disc_value"`
	PromoValue    *float64        `gorm:"promo_value" json:"promo_value"`
	CashDiscValue *float64        `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1      *float64        `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2      *float64        `gorm:"tot_disc2" json:"tot_disc2"`
	Vat           *float64        `gorm:"vat" json:"vat"`
	VatValue      *float64        `gorm:"vat_value" json:"vat_value"`
	Total         *float64        `gorm:"total" json:"total"`
	DataStatus    *int64          `gorm:"data_status" json:"data_status"`
	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
	DataSource    *int64          `gorm:"data_source" json:"data_source"`
}

func (RoList) TableName() string {
	return "sls.ro"
}
