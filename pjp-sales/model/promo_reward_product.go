package model

type PromoRewardProduct struct {
	CustID        string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoID       string `gorm:"column:promo_id" json:"promo_id"`
	PromoRewardID *int64 `gorm:"column:promo_reward_id;primaryKey;autoIncrement:true" json:"promo_reward_id"`
	ProID         int64  `gorm:"column:pro_id" json:"pro_id"`
}

func (PromoRewardProduct) TableName() string {
	return "sls.promo_reward_products"
}

type PromoRewardProductRead struct {
	CustID        string  `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoID       string  `gorm:"column:promo_id" json:"promo_id"`
	PromoRewardID *int64  `gorm:"column:promo_reward_id;primaryKey;autoIncrement:true" json:"promo_reward_id"`
	ProID         int64   `gorm:"column:pro_id" json:"pro_id"`
	UnitId        string  `gorm:"column:unit_id" json:"unit_id"`
	QtyStock      float64 `gorm:"column:qty_stock" json:"qty_stock"`
}

func (PromoRewardProductRead) TableName() string {
	return "promo.promotion_reward_products"
}

type RewardProductDetail struct {
	CustID    string  `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	ProID     int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode   string  `gorm:"column:pro_code" json:"pro_code"`
	ProName   string  `gorm:"column:pro_name" json:"pro_name"`
	ConvUnit2 float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3 float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
}

func (RewardProductDetail) TableName() string {
	return "mst.m_product"
}
