package model

import "time"

// StockMovementWarehouseTotalStock model for warehouse total stock query result
type StockMovementWarehouseTotalStock struct {
	WhID          int64  `gorm:"column:wh_id"`
	WhCode        string `gorm:"column:wh_code"`
	WhName        string `gorm:"column:wh_name"`
	OpeningStock  int64  `gorm:"column:opening_stock"`
	ChangingStock int64  `gorm:"column:changing_stock"`
	ClosingStock  int64  `gorm:"column:closing_stock"`
}

// StockMovementTransactionType model for transaction type query result
type StockMovementTransactionType struct {
	TrCode  string `gorm:"column:tr_code"`
	TrName  string `gorm:"column:tr_name"`
	NoOfDoc int64  `gorm:"column:no_of_doc"`
}

// StockMovementTopProduct model for top product query result (with qty in smallest unit)
type StockMovementTopProduct struct {
	ProID     int64   `gorm:"column:pro_id"`
	ProCode   string  `gorm:"column:pro_code"`
	ProName   string  `gorm:"column:pro_name"`
	TotalQty  float64 `gorm:"column:total_qty"` // in smallest unit
	ConvUnit2 float64 `gorm:"column:conv_unit2"`
	ConvUnit3 float64 `gorm:"column:conv_unit3"`
}

type StockLedgerRequest struct {
	StartDate        string   `json:"start_date"`
	EndDate          string   `json:"end_date"`
	WarehouseIDs     []int64  `json:"wh_id"`
	TransactionTypes []string `json:"transaction_type"`
	SupIDs           []int64  `json:"sup_id"`
	PrincipalIDs     []int64  `json:"principal_id"`
	ProductIDs       []int64  `json:"pro_id"`
	Page             int      `json:"page"`
	Limit            int      `json:"limit"`
}

type StockLedgerRow struct {
	TotalRecord int64 `gorm:"column:total_record" json:"-"`

	DistributorID   int64  `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorCode string `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName string `gorm:"column:distributor_name" json:"distributor_name"`

	WhID      int64     `gorm:"column:wh_id" json:"wh_id"`
	WhCode    string    `gorm:"column:wh_code" json:"wh_code"`
	Warehouse string    `gorm:"column:warehouse" json:"warehouse"`
	Date      time.Time `gorm:"column:date" json:"date"`

	ProID       int64  `gorm:"column:pro_id" json:"pro_id"`
	ProductCode string `gorm:"column:product_code" json:"product_code"`
	ProductName string `gorm:"column:product_name" json:"product_name"`

	OpeningStock       int64 `gorm:"column:opening_stock" json:"opening_stock"`
	OpeningStockLarge  int64 `gorm:"column:opening_stock_large" json:"opening_stock_large"`
	OpeningStockMedium int64 `gorm:"column:opening_stock_medium" json:"opening_stock_medium"`
	OpeningStockSmall  int64 `gorm:"column:opening_stock_small" json:"opening_stock_small"`

	Updates       int64 `gorm:"column:updates" json:"updates"`
	UpdatesLarge  int64 `gorm:"column:updates_large" json:"updates_large"`
	UpdatesMedium int64 `gorm:"column:updates_medium" json:"updates_medium"`
	UpdatesSmall  int64 `gorm:"column:updates_small" json:"updates_small"`

	ClosingStock       int64 `gorm:"column:closing_stock" json:"closing_stock"`
	ClosingStockLarge  int64 `gorm:"column:closing_stock_large" json:"closing_stock_large"`
	ClosingStockMedium int64 `gorm:"column:closing_stock_medium" json:"closing_stock_medium"`
	ClosingStockSmall  int64 `gorm:"column:closing_stock_small" json:"closing_stock_small"`

	TransactionType string `gorm:"column:transaction_type" json:"transaction_type"`
	ReferenceNo     string `gorm:"column:reference_no" json:"reference_no"`
	Remarks         string `gorm:"column:remarks" json:"remarks"`
}

// ReportList represents a single record in report.list used for background/download reports.
type ReportList struct {
	CustID     string    `gorm:"column:cust_id" json:"cust_id"`
	ReportID   string    `gorm:"column:report_id" json:"report_id"`
	ReportName string    `gorm:"column:report_name" json:"report_name"`
	StartDate  time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate    time.Time `gorm:"column:end_date" json:"end_date"`
	FileStatus int       `gorm:"column:file_status" json:"file_status"`
	FileURL    string    `gorm:"column:file_url" json:"file_url"`
	FileBase64 string    `gorm:"column:file_base64" json:"file_base64"`
	CreatedBy  string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ReportList) TableName() string {
	return "report.list"
}
