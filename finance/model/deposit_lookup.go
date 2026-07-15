package model

import (
	"time"
)

type CollectionNoLookup struct {
	CollectionNo string `db:"collection_no" json:"collection_no"`
}

func (CollectionNoLookup) TableName() string {
	return "acf.deposit"
}

type DepositNoLookup struct {
	DepositNo string `db:"deposit_no" json:"deposit_no"`
}

func (DepositNoLookup) TableName() string {
	return "acf.deposit"
}

type DepositStatusLookup struct {
	DepositStatus     int `db:"deposit_status" json:"deposit_status"`
	DepositStatusName int `db:"deposit_status_name" json:"deposit_status_name"`
}

func (DepositStatusLookup) TableName() string {
	return "acf.deposit"
}

type InvoiceCollectionList struct {
	CustID          string     `gorm:"cust_id" json:"cust_id"`
	RoNo            string     `gorm:"ro_no" json:"ro_no"`
	SalesmanId      *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode    *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesmanName    *string    `gorm:"salesman_name" json:"salesman_name"`
	OutletID        *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode      *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      *string    `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceAmount   *float64   `gorm:"invoice_amount" json:"invoice_amount"`
	RemainingAmount *float64   `gorm:"remaining_amount" json:"remaining_amount"`
	InvoiceNo       *string    `gorm:"invoice_no" json:"invoice_no"`
	CollectionNo    string     `gorm:"collection_no" json:"collection_no"`
	InvoiceDate     *time.Time `gorm:"invoice_date" json:"invoice_date"`
}

func (InvoiceCollectionList) TableName() string {
	return "sls.order"
}
