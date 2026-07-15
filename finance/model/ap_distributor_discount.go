package model

import (
	"time"

	"gorm.io/gorm"
)

type ApDistributorDiscount struct {
	CustID                string `gorm:"column:cust_id" json:"cust_id"`
	DistributorDiscountId *int64 `gorm:"column:distributor_discount_id;primaryKey" json:"distributor_discount_id"`
	ProId                 *int64 `gorm:"column:pro_id" json:"pro_id"`
	// ProCode               *string        `gorm:"column:pro_code" json:"pro_code"`
	// ProName               *string        `gorm:"column:pro_name" json:"pro_name"`
	PurchPrice *float64       `gorm:"column:purch_price" json:"purch_price"`
	Discount   *float64       `gorm:"column:discount" json:"discount"`
	NetPrice   *float64       `gorm:"column:net_price" json:"net_price"`
	IsActive   *bool          `gorm:"column:is_active" json:"is_active"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *ApDistributorDiscount) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()

	// intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	// m.ChqTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (ApDistributorDiscount) TableName() string {
	return "acf.account_payable_discounts"
}

type ApDistributorDiscountList struct {
	CustID                string         `gorm:"column:cust_id" json:"cust_id"`
	DistributorDiscountId *int64         `gorm:"column:distributor_discount_id;primaryKey" json:"distributor_discount_id"`
	ProId                 *int64         `gorm:"column:pro_id" json:"pro_id"`
	ProCode               *string        `gorm:"column:pro_code" json:"pro_code"`
	ProName               *string        `gorm:"column:pro_name" json:"pro_name"`
	PurchPrice            *float64       `gorm:"column:purch_price" json:"purch_price"`
	Discount              *float64       `gorm:"column:discount" json:"discount"`
	NetPrice              *float64       `gorm:"column:net_price" json:"net_price"`
	IsActive              *bool          `gorm:"column:is_active" json:"is_active"`
	CreatedBy             *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt             time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy             *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt             time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName         *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel                 bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy             *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ApDistributorDiscountList) TableName() string {
	return "acf.account_payable_discounts"
}
