package model

import (
	"time"

	"gorm.io/gorm"
)

type Discount struct {
	CustID           string    `gorm:"column:cust_id" json:"cust_id"`
	DiscountID       string    `gorm:"column:discount_id" json:"discount_id"`
	DiscountDesc     string    `gorm:"column:discount_desc" json:"discount_desc"`
	DiscountStatusID int       `gorm:"column:discount_status_id" json:"discount_status_id"`
	PublishStatusID  int       `gorm:"column:publish_status_id" json:"publish_status_id"`
	EffectiveFrom    time.Time `gorm:"column:effective_from" json:"effective_from"`
	EffectiveTo      time.Time `gorm:"column:effective_to" json:"effective_to"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy        string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy        string    `gorm:"column:updated_by" json:"updated_by"`
	Remarks          string    `gorm:"column:remarks" json:"remarks"`
}

func (Discount) TableName() string {
	return "sls.discounts"
}

func (m *Discount) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now().UTC()
	m.UpdatedAt = now

	return nil
}

type OutletRead struct {
	CustId                  string          `gorm:"column:cust_id" json:"cust_id"`
	OutletId                int             `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode              string          `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName              string          `gorm:"column:outlet_name" json:"outlet_name"`
	Address1                *string         `gorm:"column:address1" json:"address1"`
	Address2                *string         `gorm:"column:address2" json:"address2"`
	DiscGrpId               int             `gorm:"column:disc_grp_id" json:"disc_grp_id"`
	CreditLimitAction       *int            `gorm:"column:credit_limit_action" json:"credit_limit_action"`
	CreditLimitActionName   string          `gorm:"column:credit_limit_action_name" json:"credit_limit_action_name"`
	SalesInvLimitAction     *int            `gorm:"column:sales_inv_limit_action" json:"sales_inv_limit_action"`
	SalesInvLimitActionName string          `gorm:"column:sales_inv_limit_action_name" json:"sales_inv_limit_action_name"`
	ObsLimitAction          *int            `gorm:"column:obs_limit_action" json:"obs_limit_action"`
	ObsLimitActionName      string          `gorm:"column:obs_limit_action_name" json:"obs_limit_action_name"`
	OtGrpId                 int             `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	OtClassId               int             `gorm:"column:ot_class_id" json:"ot_class_id"`
	OtTypeId                int             `gorm:"column:ot_type_id" json:"ot_type_id"`
	IsActive                bool            `gorm:"column:is_active" json:"is_active"`
	IsDel                   bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy               *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt               *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy               *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt               *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName           *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy               *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt               *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletRead) TableName() string {
	return "mst.m_outlet"
}

type ProductRead struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	ProId         int             `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string          `gorm:"column:pro_code" json:"pro_code"`
	ProName       string          `gorm:"column:pro_name" json:"pro_name"`
	PrincipalId   int             `gorm:"column:principal_id" json:"principal_id"`
	ParentProId   *int            `gorm:"column:parent_pro_id" json:"parent_pro_id"`
	UnitId1       string          `gorm:"column:unit_id1" json:"unit_id1" `
	UnitId2       string          `gorm:"column:unit_id2" json:"unit_id2" `
	UnitId3       string          `gorm:"column:unit_id3" json:"unit_id3" `
	UnitId4       string          `gorm:"column:unit_id4" json:"unit_id4" `
	UnitId5       string          `gorm:"column:unit_id5" json:"unit_id5" `
	ConvUnit2     float32         `gorm:"column:conv_unit2" json:"conv_unit2" `
	ConvUnit3     float32         `gorm:"column:conv_unit3" json:"conv_unit3" `
	ConvUnit4     float32         `gorm:"column:conv_unit4" json:"conv_unit4" `
	ConvUnit5     float32         `gorm:"column:conv_unit5" json:"conv_unit5" `
	PurchPrice1   float64         `gorm:"column:purch_price1" json:"purch_price1" `
	PurchPrice2   float64         `gorm:"column:purch_price2" json:"purch_price2" `
	PurchPrice3   float64         `gorm:"column:purch_price3" json:"purch_price3" `
	PurchPrice4   float64         `gorm:"column:purch_price4" json:"purch_price4" `
	PurchPrice5   float64         `gorm:"column:purch_price5" json:"purch_price5" `
	SellPrice1    float64         `gorm:"column:sell_price1" json:"sell_price1" `
	SellPrice2    float64         `gorm:"column:sell_price2" json:"sell_price2" `
	SellPrice3    float64         `gorm:"column:sell_price3" json:"sell_price3" `
	SellPrice4    float64         `gorm:"column:sell_price4" json:"sell_price4" `
	SellPrice5    float64         `gorm:"column:sell_price5" json:"sell_price5" `
	Vat           float64         `gorm:"column:vat" json:"vat"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ProductRead) TableName() string {
	return "mst.m_product"
}

type DiscountRead struct {
	CustId        string     `gorm:"column:cust_id" json:"cust_id"`
	DiscountId    string     `gorm:"column:discount_id" json:"discount_id"`
	DiscountDesc  string     `gorm:"column:discount_desc" json:"discount_desc"`
	CreatedBy     *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
}

func (DiscountRead) TableName() string {
	return "sls.discounts"
}
