package entity

type CreateMemoJrBody struct {
	CustID     string                `json:"cust_id"`
	MjDate     *string               `json:"mj_date"`
	TrCode     *string               `json:"tr_code"`
	Notes      *string               `json:"notes"`
	DataStatus *int64                `json:"data_status"`
	CreatedBy  *int64                `json:"created_by"`
	Details    []CreateMemoJrDetBody `json:"details"`
}
type MemoJrResponse struct {
	MjNo          string              `json:"mj_no"`
	MjDate        *string             `json:"mj_date"`
	TrCode        *string             `json:"tr_code"`
	Notes         *string             `json:"notes"`
	DataStatus    *int64              `json:"data_status"`
	UpdatedByName string              `json:"updated_by_name"`
	UpdatedAt     string              `json:"updated_at"`
	Details       []MemoJrDetResponse `json:"details"`
}
type MemoJrListResponse struct {
	MjNo          string  `json:"mj_no"`
	MjDate        *string `json:"mj_date"`
	TrCode        *string `json:"tr_code"`
	Notes         *string `json:"notes"`
	DataStatus    *int64  `json:"data_status"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     string  `json:"updated_at"`
}
type UpdateMemoJrBody struct {
	MjNo       string                `json:"mj_no"`
	CustID     string                `json:"cust_id"`
	MjDate     *string               `json:"mj_date"`
	TrCode     *string               `json:"tr_code"`
	Notes      *string               `json:"notes"`
	DataStatus *int64                `json:"data_status"`
	CreatedBy  *int64                `json:"created_by"`
	UpdatedBy  int64                 `json:"updated_by"`
	Details    []UpdateMemoJrDetBody `json:"details"`
}
type DetailMemoJrParams struct {
	MjNo string `params:"mj_no" validate:"required"`
}

type UpdateMemoJrParams struct {
	MjNo string `params:"mj_no" validate:"required"`
}
type DeleteMemoJrParams struct {
	MjNo string `params:"mj_no" validate:"required"`
}
