package model

type DiscountPrincipal struct {
	CustID      string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	DiscountID  string `gorm:"column:discount_id;primaryKey;autoIncrement:false" json:"discount_id"`
	PrincipalID int64  `gorm:"column:principal_id;primaryKey;autoIncrement:false" json:"principal_id"`
}

func (DiscountPrincipal) TableName() string {
	return "sls.discount_principals"
}

type DiscountPrincipalDetail struct {
	CustID        string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	DiscountID    string `gorm:"column:discount_id;primaryKey;autoIncrement:false" json:"discount_id"`
	PrincipalID   int64  `gorm:"column:principal_id;primaryKey;autoIncrement:false" json:"principal_id"`
	PrincipalCode string `gorm:"column:principal_code" json:"principal_code"`
	PrincipalName string `gorm:"column:principal_name" json:"principal_name"`
}

func (DiscountPrincipalDetail) TableName() string {
	return "sls.discount_principals"
}
