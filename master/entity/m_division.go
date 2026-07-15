package entity

import "time"

type MDivisionDetailsResponse struct {
	DivisionID    int64      `json:"division_id"`
	DivisionCode  string     `json:"division_code"`
	DivisionName  string     `json:"division_name"`
	IsActive      bool       `json:"is_active"`
	CreatedBy     *int64     `json:"created_by"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName string     `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
	IsDel         bool       `json:"is_del"`
	DeletedBy     *int64     `json:"deleted_by"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type MDivisionLookupResponse struct {
	DivisionID   int64      `json:"division_id"`
	DivisionCode string     `json:"division_code"`
	DivisionName string     `json:"division_name"`
	IsActive     bool       `json:"is_active"`
	UpdatedBy    *int64     `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type DetailMDivisionParams struct {
	MDivisionId int64 `params:"division_id" validate:"required"`
}
type UpdateMDivisionParams struct {
	MDivisionId int64 `params:"division_id" validate:"required"`
}

type DeleteMDivisionParams struct {
	MDivisionId int64 `params:"division_id" validate:"required"`
}
type CreateDivisionBody struct {
	CustId       string `json:"cust_id"`
	DivisionCode string `json:"division_code" validate:"required,number,max=5"`
	DivisionName string `json:"division_name" validate:"required,alphanumericSpace,max=50"`
	IsActive     bool   `json:"is_active"`
	CreatedBy    int64  `json:"created_by"`
}

type UpdateDivisionBody struct {
	CustId       string `json:"cust_id"`
	DivisionCode string `json:"division_code" validate:"required,number,max=5"`
	DivisionName string `json:"division_name" validate:"required,alphanumericSpace,max=50"`
	IsActive     *bool  `json:"is_active"`
	UpdatedBy    int64  `json:"updated_by"`
}
