package model

import (
	"time"
)

// ExpenseType represents acf.expense_type table.
type ExpenseType struct {
	ExpenseTypeID   int        `gorm:"column:expense_type_id;primaryKey;autoIncrement" json:"expense_type_id"`
	CustID          string     `gorm:"column:cust_id;type:varchar(50)" json:"cust_id"`
	ExpenseTypeCode string     `gorm:"column:expense_type_code;type:varchar(20)" json:"expense_type_code"`
	ExpenseTypeName string     `gorm:"column:expense_type_name;type:varchar(50)" json:"expense_type_name"`
	IsActive        bool       `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedBy       int        `gorm:"column:created_by;type:int4" json:"created_by"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy       *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	DeletedBy       *int       `gorm:"column:deleted_by;type:int4" json:"deleted_by"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
	IsDel           bool       `gorm:"column:is_del;default:false" json:"is_del"`
}

func (ExpenseType) TableName() string {
	return "acf.expense_type"
}

// ExpenseTypeList represents read model for expense_type with join to sys.m_user
type ExpenseTypeList struct {
	ExpenseTypeID   int        `gorm:"column:expense_type_id;primaryKey" json:"expense_type_id"`
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	ExpenseTypeCode string     `gorm:"column:expense_type_code" json:"expense_type_code"`
	ExpenseTypeName string     `gorm:"column:expense_type_name" json:"expense_type_name"`
	IsActive        bool       `gorm:"column:is_active" json:"is_active"`
	CreatedBy       int        `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy       *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName   *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel           bool       `gorm:"column:is_del" json:"is_del"`
	DeletedBy       *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ExpenseTypeList) TableName() string {
	return "acf.expense_type"
}
