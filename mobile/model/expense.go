package model

import (
	"time"

	"gorm.io/gorm"
)

// ExpenseType represents acf.expense_type table
// Global table, tidak ada cust_id
type ExpenseType struct {
	ExpenseTypeID   int        `gorm:"column:expense_type_id;primaryKey;autoIncrement" json:"expense_type_id"`
	ExpenseTypeCode *string    `gorm:"column:expense_type_code;type:varchar(10)" json:"expense_type_code"`
	ExpenseTypeName *string    `gorm:"column:expense_type_name;type:varchar(100)" json:"expense_type_name"`
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

// Expense represents acf.expense table (Header)
// Composite primary key: (cust_id, expense_id)
type Expense struct {
	CustID        string     `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	ExpenseID     int64      `gorm:"column:expense_id;type:bigserial;primaryKey;autoIncrement" json:"expense_id"`
	DocNo         *string    `gorm:"column:doc_no;type:varchar(50)" json:"doc_no"`
	ExpenseTypeID int        `gorm:"column:expense_type_id;type:int4;not null" json:"expense_type_id"`
	Date          time.Time  `gorm:"column:date;type:date;not null" json:"date"`
	Amount        float64    `gorm:"column:amount;type:numeric(20,4);default:0" json:"amount"`
	Balance       float64    `gorm:"column:balance;type:numeric(20,4);default:0" json:"balance"`
	Note          *string    `gorm:"column:note;type:varchar(100)" json:"note"`
	CreatedBy     int        `gorm:"column:created_by;type:int4;not null" json:"created_by"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy     *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt     *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	DeletedBy     *int       `gorm:"column:deleted_by;type:int4" json:"deleted_by"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
	IsDel         bool       `gorm:"column:is_del;default:false" json:"is_del"`

	// Relations
	ExpenseType ExpenseType   `gorm:"foreignKey:ExpenseTypeID;references:ExpenseTypeID" json:"expense_type,omitempty"`
	Details     []ExpenseDet  `gorm:"foreignKey:CustID,ExpenseID;references:CustID,ExpenseID" json:"details,omitempty"`
	Files       []ExpenseFile `gorm:"foreignKey:CustID,ExpenseID;references:CustID,ExpenseID" json:"files,omitempty"`
	Source      *int          `gorm:"column:source;type:int4" json:"source"`
	CollectorID *int          `gorm:"column:collector_id;type:int4" json:"collector_id"`
}

func (Expense) TableName() string {
	return "acf.expense"
}

// ExpenseDet represents acf.expense_det table (Detail Outlet)
// Composite primary key: (cust_id, expense_det_id)
type ExpenseDet struct {
	CustID       string `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	ExpenseDetID int64  `gorm:"column:expense_det_id;type:bigserial;primaryKey;autoIncrement" json:"expense_det_id"`
	ExpenseID    int64  `gorm:"column:expense_id;type:int8;not null" json:"expense_id"`
	OutletID     int    `gorm:"column:outlet_id;type:int4;not null" json:"outlet_id"`
}

func (ExpenseDet) TableName() string {
	return "acf.expense_det"
}

// ExpenseFile represents acf.expense_file table (Lampiran File)
// Composite primary key: (cust_id, expense_file_id)
type ExpenseFile struct {
	CustID        string    `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	ExpenseFileID int64     `gorm:"column:expense_file_id;type:bigserial;primaryKey;autoIncrement" json:"expense_file_id"`
	ExpenseID     int64     `gorm:"column:expense_id;type:int8;not null" json:"expense_id"`
	FileName      string    `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FileURL       string    `gorm:"column:file_url;type:varchar(500);not null" json:"file_url"`
	FileKey       string    `gorm:"column:file_key;type:acf.media_category_type;not null" json:"file_key"`
	MediaCategory *string   `gorm:"column:media_category;type:text" json:"media_category"`
	FileSize      int64     `gorm:"column:file_size;type:bigint;not null" json:"file_size"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (ExpenseFile) TableName() string {
	return "acf.expense_file"
}

// ExpenseListRead represents read model for expense list with joins
type ExpenseListRead struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	ExpenseID       int64      `gorm:"column:expense_id" json:"expense_id"`
	DocNo           *string    `gorm:"column:doc_no" json:"doc_no"`
	ExpenseTypeID   int        `gorm:"column:expense_type_id" json:"expense_type_id"`
	ExpenseTypeName *string    `gorm:"column:expense_type_name" json:"expense_type_name"`
	Date            time.Time  `gorm:"column:date" json:"date"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	Note            *string    `gorm:"column:note" json:"note"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ExpenseDetailRead represents read model for expense detail with all relations
type ExpenseDetailRead struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	ExpenseID       int64      `gorm:"column:expense_id" json:"expense_id"`
	DocNo           *string    `gorm:"column:doc_no" json:"doc_no"`
	ExpenseTypeID   int        `gorm:"column:expense_type_id" json:"expense_type_id"`
	ExpenseTypeCode *string    `gorm:"column:expense_type_code" json:"expense_type_code"`
	ExpenseTypeName *string    `gorm:"column:expense_type_name" json:"expense_type_name"`
	Date            time.Time  `gorm:"column:date" json:"date"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	Note            *string    `gorm:"column:note" json:"note"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// ExpenseDetRead represents read model for expense_det with outlet info
type ExpenseDetRead struct {
	ExpenseDetID   int64  `gorm:"column:expense_det_id" json:"expense_det_id"`
	ExpenseID      int64  `gorm:"column:expense_id" json:"expense_id"`
	OutletID       int    `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode     string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName     string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletAddress1 string `gorm:"column:outlet_address1" json:"outlet_address1"`
}

// ExpenseFileRead represents read model for expense_file
type ExpenseFileRead struct {
	ExpenseFileID int64     `gorm:"column:expense_file_id" json:"expense_file_id"`
	ExpenseID     int64     `gorm:"column:expense_id" json:"expense_id"`
	FileName      string    `gorm:"column:file_name" json:"file_name"`
	FileURL       string    `gorm:"column:file_url" json:"file_url"`
	FileKey       string    `gorm:"column:file_key" json:"file_key"`
	MediaCategory *string   `gorm:"column:media_category" json:"media_category"`
	FileSize      int64     `gorm:"column:file_size" json:"file_size"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

// OutletLookupPJP represents outlet lookup from PJP
type OutletLookupPJP struct {
	OutletID   int    `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName string `gorm:"column:outlet_name" json:"outlet_name"`
}

// BeforeCreate hook untuk Expense
func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	return nil
}

// BeforeUpdate hook untuk Expense
func (e *Expense) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	if e.UpdatedAt == nil {
		e.UpdatedAt = &now
	}
	return nil
}

// ExpenseListItem represents a simplified view of an expense (doc_no and amount)
type ExpenseListItem struct {
	ExpenseID   int64   `json:"expense_id"`
	ExpenseName string  `json:"expense_name"`
	Amount      float64 `json:"amount"`
	DocNo       *string `json:"doc_no"`
}
