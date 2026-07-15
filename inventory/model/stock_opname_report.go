package model

import "time"

type StockOpnameReport struct {
	CustID        string    `gorm:"column:cust_id" json:"cust_id"`
	StockReportID string    `gorm:"column:stock_report_id" json:"stock_report_id"`
	DocNo         string    `gorm:"column:doc_no" json:"doc_no"`
	Status        int       `gorm:"column:status" json:"status"`
	CreatedBy     *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy     *int64    `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (StockOpnameReport) TableName() string {
	return "inv.stock_opname_reports"
}
