package model

import "gorm.io/gorm"

const (
	STATUS_MANAGE_PRICE_NON_ACTIVE = 0
	STATUS_MANAGE_PRICE_SUBMIT     = 1
	STATUS_MANAGE_PRICE_ACTIVE     = 2
)

type ManageMinimumPrice struct {
	ManageMinimumPriceId     int             `gorm:"column:manage_minimum_price_id" json:"manage_minimum_price_id"`
	BasePrice                int             `gorm:"column:base_price" json:"base_price"`
	LimitAction              int             `gorm:"column:limit_action" json:"limit_action"`
	Threshold                float64         `gorm:"column:threshold" json:"threshold"`
	StatusManageMinimumPrice int             `gorm:"column:status_manage_minimum_price" json:"status_manage_minimum_price"`
	ProId                    int             `gorm:"column:pro_id" json:"pro_id"`
	Price1                   float64         `gorm:"column:price1" json:"price1"`
	Price2                   float64         `gorm:"column:price2" json:"price2"`
	Price3                   float64         `gorm:"column:price3" json:"price3"`
	Price4                   float64         `gorm:"column:price4" json:"price4"`
	Price5                   float64         `gorm:"column:price5" json:"price5"`
	PriceMinimum1            float64         `gorm:"column:price1_minimum" json:"price1_minimum"`
	PriceMinimum2            float64         `gorm:"column:price2_minimum" json:"price2_minimum"`
	PriceMinimum3            float64         `gorm:"column:price3_minimum" json:"price3_minimum"`
	PriceMinimum4            float64         `gorm:"column:price4_minimum" json:"price4_minimum"`
	PriceMinimum5            float64         `gorm:"column:price5_minimum" json:"price5_minimum"`
	UnitId1                  string          `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2                  string          `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3                  string          `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4                  string          `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5                  string          `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2                int             `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3                int             `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4                int             `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5                int             `gorm:"column:conv_unit5" json:"conv_unit5"`
	DeletedAt                *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (ManageMinimumPrice) TableName() string {
	return "mst.manage_minimum_price"
}
