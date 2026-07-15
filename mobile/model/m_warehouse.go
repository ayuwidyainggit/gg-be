package model

type MWarehouse struct {
	WHID   int    `gorm:"column:wh_id" json:"wh_id"`
	WHName string `gorm:"column:wh_name" json:"wh_name"`
	WHCode string `gorm:"column:wh_code" json:"wh_code"`
	Status string `gorm:"column:status" json:"status"`
}
