package model

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ExpenseID     int64          `gorm:"column:expense_id;primaryKey;autoIncrement" json:"expense_id"`
	ExpenseTypeID int            `gorm:"column:expense_type_id" json:"expense_type_id"`
	Date          *time.Time     `gorm:"column:date" json:"date"`
	Amount        *float64       `gorm:"column:amount" json:"amount"`
	Note          *string        `gorm:"column:note" json:"note"`
	CreatedBy     int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DocNo         *string        `gorm:"column:doc_no" json:"doc_no"`
	Source        *int           `gorm:"column:source" json:"source"`
	Balance       *float64       `gorm:"column:balance" json:"balance"`
	CollectorID   *int64         `gorm:"column:collector_id" json:"collector_id"`
}

type ExpenseList struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	ExpenseID       int64      `gorm:"column:expense_id" json:"expense_id"`
	ExpenseTypeID   int        `gorm:"column:expense_type_id" json:"expense_type_id"`
	Date            *time.Time `gorm:"column:date" json:"date"`
	Amount          *float64   `gorm:"column:amount" json:"amount"`
	Note            *string    `gorm:"column:note" json:"note"`
	DocNo           *string    `gorm:"column:doc_no" json:"doc_no"`
	Balance         *float64   `gorm:"column:balance" json:"balance"`
	CollectorID     *int64     `gorm:"column:collector_id" json:"collector_id"`
	CreatedBy       int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	ExpenseTypeCode *string    `gorm:"column:expense_type_code" json:"expense_type_code"`
	ExpenseTypeName *string    `gorm:"column:expense_type_name" json:"expense_type_name"`
	CollectorName   *string    `gorm:"column:collector_name" json:"collector_name"`
	RemainingAmount *float64   `gorm:"column:remaining_amount" json:"remaining_amount"`
	Source          *int       `gorm:"column:source" json:"source"`
}

func (Expense) TableName() string {
	return "acf.expense"
}
