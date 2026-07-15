package entity

type CreateMCoaTypeBody struct {
	CoaTypeName *string `json:"coa_type_name"`
	CoaGroup    *string `json:"coa_group"`
	DefBlc      *string `json:"def_blc"`
	SortIndex   *int64  `json:"sort_index"`
	CoaKind     *int64  `json:"coa_kind"`
	IsActive    bool    `json:"is_active"`
	CreatedBy   int64   `json:"created_by"`
}

type UpdateMCoaTypeBody struct {
	CoaTypeName *string `json:"coa_type_name"`
	CoaGroup    *string `json:"coa_group"`
	DefBlc      *string `json:"def_blc"`
	SortIndex   *int64  `json:"sort_index"`
	CoaKind     *int64  `json:"coa_kind"`
	IsActive    *bool   `json:"is_active"`
	UpdatedBy   int64   `json:"updated_by"`
}

type MCoaTypeListResponse struct {
	CoaTypeID     int64   `json:"coa_type_id"`
	CoaTypeName   *string `json:"coa_type_name"`
	CoaGroup      *string `json:"coa_group"`
	DefBlc        *string `json:"def_blc"`
	SortIndex     *int64  `json:"sort_index"`
	CoaKind       *int64  `json:"coa_kind"`
	IsActive      bool    `json:"is_active"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     string  `json:"updated_at"`
}

type MCoaTypeLookupListResponse struct {
	CoaTypeID   int64   `json:"coa_type_id"`
	CoaTypeName *string `json:"coa_type_name"`
	CoaGroup    *string `json:"coa_group"`
	DefBlc      *string `json:"def_blc"`
	SortIndex   *int64  `json:"sort_index"`
	CoaKind     *int64  `json:"coa_kind"`
}

type MCoaTypeResponse struct {
	CoaTypeID     int64   `json:"coa_type_id"`
	CoaTypeName   *string `json:"coa_type_name"`
	CoaGroup      *string `json:"coa_group"`
	DefBlc        *string `json:"def_blc"`
	SortIndex     *int64  `json:"sort_index"`
	CoaKind       *int64  `json:"coa_kind"`
	IsActive      bool    `json:"is_active"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     string  `json:"updated_at"`
}
type DetailMCoaTypeParams struct {
	CoaTypeID int64 `params:"coa_type_id" validate:"required"`
}
type DeleteMCoaTypeParams struct {
	CoaTypeID int64 `params:"coa_type_id" validate:"required"`
}
type UpdateMCoaTypeParams struct {
	CoaTypeID int64 `params:"coa_type_id" validate:"required"`
}
