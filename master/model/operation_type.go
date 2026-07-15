package model

import (
	"time"
)

type OperationType struct {
	OperationTypeCode string     `db:"operation_type_code" json:"operation_type_code"`
	OperationTypeName string     `db:"operation_type_name" json:"operation_type_name"`
	CreatedBy         *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt         *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy         *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt         *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy         *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt         *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName     *string    `db:"updated_by_name" json:"updated_by_name"`
}

type OperationTypeUpdate struct {
	OperationTypeCode *string    `json:"operation_type_code,omitempty" sql:"operation_type_code"`
	OperationTypeName *string    `json:"operation_type_name,omitempty" sql:"operation_type_name"`
	UpdatedAt         *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy         *int64     `json:"updated_by" sql:"updated_by"`
}
