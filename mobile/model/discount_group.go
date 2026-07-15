package model

type DiscountGroup struct {
	CustID      string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	DiscountID  string `gorm:"column:discount_id;primaryKey;autoIncrement:false" json:"discount_id"`
	DiscGrpID   int    `gorm:"column:disc_grp_id;primaryKey;autoIncrement:false" json:"disc_grp_id"`
}

func (DiscountGroup) TableName() string {
	return "sls.discount_groups"
}

type DiscountGroupDetail struct {
	CustID      string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	DiscountID  string `gorm:"column:discount_id;primaryKey;autoIncrement:false" json:"discount_id"`
	DiscGrpID   int    `gorm:"column:disc_grp_id;primaryKey;autoIncrement:false" json:"disc_grp_id"`
	DiscGrpCode string `gorm:"column:disc_grp_code" json:"disc_grp_code"`
	DiscGrpName string `gorm:"column:disc_grp_name" json:"disc_grp_name"`
}

func (DiscountGroupDetail) TableName() string {
	return "sls.discount_groups"
}
