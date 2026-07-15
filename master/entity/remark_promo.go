package entity

import "time"

type RemarkPromoResponse struct {
	RemPromoId    int        `json:"rem_promo_id"`
	RemPromoCode  string     `json:"rem_promo_code"`
	RemPromoName  string     `json:"rem_promo_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type CreateRemarkPromoBody struct {
	CustId       string `json:"cust_id"`
	RemPromoCode string `json:"rem_promo_code" validate:"max=5,alphanumericSpace"`
	RemPromoName string `json:"rem_promo_name" validate:"max=100"`
	IsActive     bool   `json:"is_active"`
	CreatedBy    int64  `json:"created_by" validate:"required"`
	UpdatedBy    int64  `json:"updated_by"`
}

type DetailRemarkPromoParams struct {
	RemPromoId int `params:"rem_promo_id" validate:"required"`
}

type UpdateRemarkPromoParams struct {
	RemPromoId int `params:"rem_promo_id" validate:"required"`
}

type DeleteRemarkPromoParams struct {
	RemPromoId int `params:"rem_promo_id" validate:"required"`
}

type UpdateRemarkPromoRequest struct {
	CustId       string `json:"cust_id"`
	RemPromoCode string `json:"rem_promo_code" validate:"max=5,alphanumericSpace"`
	RemPromoName string `json:"rem_promo_name" validate:"max=100"`
	IsActive     bool   `json:"is_active"`
	UpdatedBy    int64  `json:"updated_by"`
}
