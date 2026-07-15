package entity

import "time"

type PluProductResponse struct {
	PluProId      int        `json:"plu_pro_id"`
	PluGrpId      int        `json:"plu_grp_id"`
	ProId         int        `json:"pro_id"`
	ProCode       string     `json:"pro_code"`
	ProName       string     `json:"pro_name"`
	PluNo         string     `json:"plu_no"`
	UpdatedBy     *int       `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type CreatePluProductBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	PluGrpId  int    `json:"plu_grp_id" validate:"required"`
	ProId     int    `json:"pro_id" validate:"required"`
	PluNo     string `json:"plu_no" validate:"required,max=50,alphanumericSpace"`
	CreatedBy int64  `json:"created_by"`
	UpdatedBy int64  `json:"updated_by"`
}

type DetailPluProductParams struct {
	PluProId int `params:"plu_pro_id" validate:"required"`
}

type UpdatePluProductParams struct {
	PluProId int `params:"plu_pro_id" validate:"required"`
}

type DeletePluProductParams struct {
	PluProId int `params:"plu_pro_id" validate:"required"`
}

type UpdatePluProductRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	PluGrpId  int    `json:"plu_grp_id" validate:"required"`
	ProId     int    `json:"pro_id" validate:"required"`
	PluNo     string `json:"plu_no" validate:"required,max=50,alphanumericSpace"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
}
