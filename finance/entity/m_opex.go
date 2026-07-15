package entity

type CreateMOpexBody struct {
	CustId    string  `json:"cust_id"`
	OpexCode  *string `json:"opex_code"`
	OpexName  *string `json:"opex_name"`
	CoaId     *int64  `json:"coa_id"`
	IsActive  *bool   `json:"is_active"`
	CreatedBy *int64  `json:"created_by"`
	UpdatedBy *int64  `json:"updated_by"`
}

type MOpexResponse struct {
	OpexId        *int   `json:"opex_id"`
	OpexCode      string `json:"opex_code"`
	OpexName      string `json:"opex_name"`
	CoaId         *int64 `json:"coa_id"`
	CoaCode       string `json:"coa_code"`
	CoaName       string `json:"coa_name"`
	IsActive      *bool  `json:"is_active"`
	CreatedBy     *int64 `json:"created_by"`
	UpdatedByName string `json:"updated_by_name"`
	UpdatedAt     string `json:"updated_at"`
}

type DetailMOpexParams struct {
	OpexId int `params:"opex_id" validate:"required"`
}
type DeleteMOpexParams struct {
	OpexId int `params:"opex_id" validate:"required"`
}

type UpdateMOpexParams struct {
	OpexId int `params:"opex_id" validate:"required"`
}

type MOpexListResponse struct {
	OpexId        int     `json:"opex_id"`
	OpexCode      *string `json:"opex_code"`
	OpexName      *string `json:"opex_name"`
	CoaId         *int64  `json:"coa_id"`
	CoaCode       string  `json:"coa_code"`
	CoaName       string  `json:"coa_name"`
	IsActive      *bool   `json:"is_active"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     string  `json:"updated_at"`
}

type UpdateMOpexBody struct {
	CustId    string  `json:"cust_id"`
	OpexCode  *string `json:"opex_code"`
	OpexName  *string `json:"opex_name"`
	CoaId     *int64  `json:"coa_id"`
	IsActive  *bool   `json:"is_active"`
	CreatedBy *int64  `json:"created_by"`
	CreatedAt *string `json:"created_at"`
	UpdatedBy int64   `json:"updated_by"`
}
