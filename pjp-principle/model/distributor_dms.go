package model

import "time"

type DistributorDms struct {
	DistributorID      int       `gorm:"column:distributor_id;type:int" json:"distributor_id"`
	DistributorCode    string    `gorm:"column:distributor_code;type:varchar(125)" json:"distributor_code"`
	DistributorName    string    `gorm:"column:distributor_name;type:varchar(125)" json:"distributor_name"`
	Longitude          string    `gorm:"column:longitude;type:varchar(125)" json:"longitude"`
	Latitude           string    `gorm:"column:latitude;type:varchar(125)" json:"latitude"`
	DistributorStatus  string    `gorm:"column:distributor_status;type:varchar(125)" json:"distributor_status"`
	DistributorAddress string    `gorm:"column:address;type:varchar(125);null" json:"address"`
	AvgSalesWeek       float64   `gorm:"column:avg_sales_week;type:numeric(10,2)" json:"avg_sales_week"`
	CustID             string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (DistributorDms) TableName() string {
	return "mst.m_distributor"
}
