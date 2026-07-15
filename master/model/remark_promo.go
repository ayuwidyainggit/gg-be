package model

import "time"

type RemarkPromo struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	RemPromoId    int        `json:"rem_promo_id" db:"rem_promo_id"`
	RemPromoCode  string     `json:"rem_promo_code" db:"rem_promo_code"`
	RemPromoName  *string    `json:"rem_promo_name" db:"rem_promo_name"`
	IsActive      *bool      `json:"is_active" db:"is_active"`
	CreatedBy     *int64     `json:"created_by" db:"created_by"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
	IsDel         bool       `json:"is_del" db:"is_del"`
	DeletedBy     *int64     `json:"deleted_by" db:"deleted_by"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
}

type RemarkPromoUpdate struct {
	RemPromoCode *string    `json:"rem_promo_code" sql:"rem_promo_code"`
	RemPromoName *string    `json:"rem_promo_name" sql:"rem_promo_name"`
	IsActive     *bool      `json:"is_active" sql:"is_active"`
	UpdatedBy    *int64     `json:"updated_by" sql:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at" sql:"updated_at"`
}
