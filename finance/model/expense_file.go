package model

import "time"

// ExpenseFile maps to acf.expense_file.
// file_key uses DB type acf.media_category_type (e.g. 'expense').
type ExpenseFile struct {
	CustID        string  `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	ExpenseFileID int64   `gorm:"column:expense_file_id;primaryKey;autoIncrement" json:"expense_file_id"`
	ExpenseID     int64   `gorm:"column:expense_id" json:"expense_id"`
	FileName       string     `gorm:"column:file_name" json:"file_name"`
	FileURL        string     `gorm:"column:file_url" json:"file_url"`
	FileKey        string     `gorm:"column:file_key" json:"file_key"` // acf.media_category_type
	MediaCategory  *string    `gorm:"column:media_category" json:"media_category"`
	FileSize      int64    `gorm:"column:file_size" json:"file_size"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ExpenseFile) TableName() string {
	return "acf.expense_file"
}
