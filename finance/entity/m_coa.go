package entity

type CreateMCoaBody struct {
	CustID    string  `json:"cust_id"`
	CoaID     *int64  `json:"coa_id"`
	CoaCode   *string `json:"coa_code"`
	CoaName   *string `json:"coa_name"`
	Level     *int64  `json:"level"`
	CoaTypeID *int64  `json:"coa_type_id"`
	ParentID  *int64  `json:"parent_id"`
	CashType  *int64  `json:"cash_type"`
	DefBlc    *string `json:"def_blc"`
	IsDetail  bool    `json:"is_detail"`
	CreatedBy int64   `json:"created_by"`
	UpdatedBy int64   `json:"updated_by"`
}

type UpdateMCoaBody struct {
	CustID    string  `json:"cust_id"`
	CoaCode   *string `json:"coa_code"`
	CoaName   *string `json:"coa_name"`
	Level     *int64  `json:"level"`
	CoaTypeID *int64  `json:"coa_type_id"`
	ParentID  *int64  `json:"parent_id"`
	CashType  *int64  `json:"cash_type"`
	DefBlc    *string `json:"def_blc"`
	IsDetail  *bool   `json:"is_detail"`
	UpdatedBy int64   `json:"updated_by"`
}
type MCoaResponse struct {
	CoaID         int     `json:"coa_id"`
	CoaCode       *string `json:"coa_code"`
	CoaName       *string `json:"coa_name"`
	Level         *int64  `json:"level"`
	CoaTypeID     *int64  `json:"coa_type_id"`
	CoaTypeName   string  `json:"coa_type_name"`
	ParentID      *int64  `json:"parent_id"`
	CashType      *int64  `json:"cash_type"`
	DefBlc        *string `json:"def_blc"`
	IsDetail      bool    `json:"is_detail"`
	UpdatedAt     *string `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
}
type MCoaListResponse struct {
	CoaID         int64   `json:"coa_id"`
	CoaCode       *string `json:"coa_code"`
	CoaName       *string `json:"coa_name"`
	Level         *int64  `json:"level"`
	CoaTypeID     *int64  `json:"coa_type_id"`
	CoaTypeName   string  `json:"coa_type_name"`
	ParentID      *int64  `json:"parent_id"`
	CashType      *int64  `json:"cash_type"`
	DefBlc        *string `json:"def_blc"`
	IsDetail      bool    `json:"is_detail"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     *string `json:"updated_at"`
}

type MCoaLookupListResponse struct {
	CoaID       int64   `json:"coa_id"`
	CoaCode     *string `json:"coa_code"`
	CoaName     *string `json:"coa_name"`
	Level       *int64  `json:"level"`
	CoaTypeID   *int64  `json:"coa_type_id"`
	CoaTypeName string  `json:"coa_type_name"`
	ParentID    *int64  `json:"parent_id"`
	CashType    *int64  `json:"cash_type"`
	DefBlc      *string `json:"def_blc"`
	IsDetail    bool    `json:"is_detail"`
}

type DetailMCoaParams struct {
	CoaID int64 `params:"coa_id" validate:"required"`
}

type UpdateMCoaParams struct {
	CoaID int64 `params:"coa_id" validate:"required"`
}
type DeleteMCoaParams struct {
	CoaID int64 `params:"coa_id" validate:"required"`
}
