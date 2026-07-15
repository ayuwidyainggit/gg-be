package model

import (
	"time"

	"gorm.io/gorm"
)

type AccountPayableList struct {
	PoNo            string         `gorm:"column:po_no_doc" json:"po_no"`
	InvDate         *time.Time     `gorm:"column:inv_date" json:"inv_date"`
	InvDueDate      *time.Time     `gorm:"column:inv_due_date" json:"inv_due_date"`
	InvNo           string         `gorm:"column:inv_no" json:"inv_no"`
	InvAmount       float64        `gorm:"column:inv_amount" json:"inv_amount"`
	AmountPaid      float64        `gorm:"column:amount_paid" json:"amount_paid"`
	RemainingAmount float64        `gorm:"column:remaining_amount" json:"remaining_amount"`
	SupplierId      int64          `gorm:"column:supplier_id" json:"supplier_id"`
	SupplierCode    string         `gorm:"column:supplier_code" json:"supplier_code"`
	Supplier        string         `gorm:"column:supplier" json:"supplier"`
	DistributorId   int64          `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorCode string         `gorm:"column:distributor_code" json:"distributor_code"`
	Distributor     string         `gorm:"column:distributor" json:"distributor"`
	InvStatus       string         `gorm:"column:inv_status" json:"inv_status"`
	DueDateStatus   string         `gorm:"column:due_date_status" json:"due_date_status"`
	Aging           int64          `gorm:"column:aging" json:"aging"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedByName   *string        `gorm:"column:created_by_name" json:"created_by_name"`
	CreatedAt       *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName   *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy       *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AccountPayableList) TableName() string {
	return "acf.account_payable"
}

type AccountPayableListDet struct {
	PaymentMethod  *int       `gorm:"column:payment_method" json:"payment_method"`
	PaymentDate    *time.Time `gorm:"column:payment_date" json:"payment_date"`
	PaymentBalance *float64   `gorm:"column:payment_balance" json:"payment_balance"`
	DocumentNo     *string    `gorm:"column:document_no" json:"document_no"`
	Amount         *float64   `gorm:"column:amount" json:"amount"`
	UpdatedBy      *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName  *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AccountPayableListDet) TableName() string {
	return "acf.account_payable"
}
