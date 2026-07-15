package model

import "time"

type StockOpnameDetail struct {
	CustID    string    `gorm:"column:cust_id" json:"cust_id"`
	DocNo     string    `gorm:"column:doc_no" json:"doc_no"`
	ProID     int64     `gorm:"column:pro_id" json:"pro_id"`
	QtyStock  float32   `gorm:"column:qty_stock" json:"qty_stock"`
	QtyOpname float32   `gorm:"column:qty_opname" json:"qty_opname"`
	CreatedBy *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
}

func (StockOpnameDetail) TableName() string {
	return "inv.stock_opname_details"
}
