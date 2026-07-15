package model

import "time"

type OutletDms struct {
	OutletID      int       `gorm:"column:outlet_id;type:int" json:"outlet_id"`
	OutletCode    string    `gorm:"column:outlet_code;type:varchar(125)" json:"outlet_code"`
	OutletName    string    `gorm:"column:outlet_name;type:varchar(125)" json:"outlet_name"`
	Longitude     string    `gorm:"column:longitude;type:varchar(125)" json:"longitude"`
	Latitude      string    `gorm:"column:latitude;type:varchar(125)" json:"latitude"`
	OutletStatus  string    `gorm:"column:outlet_status;type:varchar(125)" json:"outlet_status"`
	OutletAddress string    `gorm:"column:address1;type:varchar(125);null" json:"address1"`
	AvgSalesWeek  float64   `gorm:"column:avg_sales_week;type:numeric(10,2)" json:"avg_sales_week"`
	CustID        string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (OutletDms) TableName() string {
	return "mst.m_outlet"
}
