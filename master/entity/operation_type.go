package entity

import (
	"time"
)

type OperationTypeResponse struct {
	OperationTypeCode string     `json:"operation_type_code"`
	OperationTypeName string     `json:"operation_type_name"`
	IsActive          bool       `json:"is_active"`
	UpdatedBy         *int64     `json:"updated_by"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type OperationTypeListResponse struct {
	OperationTypeCode string     `json:"operation_type_code"`
	OperationTypeName string     `json:"operation_type_name"`
	UpdatedBy         *int64     `json:"updated_by"`
	UpdatedAt         *time.Time `json:"updated_at"`
	UpdatedByName     string     `json:"updated_by_name"`
}

type OperationTypeLookupResponse struct {
	OperationTypeCode string `json:"operation_type_code"`
	OperationTypeName string `json:"operation_type_name"`
}

type CreateOperationTypeBody struct {
	CreatedBy         int64  `json:"created_by" validate:"required"`
	OperationTypeCode string `json:"operation_type_code" validate:"required,max=5,alphanumericSpace"`
	OperationTypeName string `json:"operation_type_name" validate:"required,max=50"`
	IsActive          bool   `json:"is_active"`
}

type DetailOperationTypeParams struct {
	OperationTypeCode string `params:"operation_type_code" validate:"required"`
}

type UpdateOperationTypeParams struct {
	OperationTypeCode string `params:"operation_type_code" validate:"required"`
}

type DeleteOperationTypeParams struct {
	OperationTypeCode string `params:"operation_type_code" validate:"required"`
}

type UpdateOperationTypeRequest struct {
	UpdatedBy         int64  `json:"updated_by" validate:"required"`
	OperationTypeCode string `json:"operation_type_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	OperationTypeName string `json:"operation_type_name,omitempty" validate:"max=50,omitempty"`
	IsActive          *bool  `json:"is_active,omitempty"`
}
