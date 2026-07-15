package entity

type CreateMChequeRejectBody struct {
	CustId        string  `json:"cust_id"`
	ChqRejectName *string `json:"chq_reject_name"`
	CreatedBy     *int64  `json:"created_by"`
}

type MChequeRejectResponse struct {
	ChqRejectId   *int    `json:"chq_reject_id"`
	ChqRejectName *string `json:"chq_reject_name"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     string  `json:"updated_at"`
}

type MChequeRejectListResponse struct {
	ChqRejectId   *int    `json:"chq_reject_id"`
	ChqRejectName *string `json:"chq_reject_name"`
	UpdatedByName string  `json:"updated_by_name"`
	UpdatedAt     string  `json:"updated_at"`
}

type UpdateMChequeRejectBody struct {
	CustId        string  `json:"cust_id"`
	ChqRejectId   *int    `json:"chq_reject_id"`
	ChqRejectName *string `json:"chq_reject_name"`
	CreatedBy     *int64  `json:"created_by"`
	CreatedAt     *string `json:"created_at"`
	UpdatedBy     int64   `json:"updated_by"`
}

type DetailMChequeRejectParams struct {
	ChqRejectId int `params:"chq_reject_id" validate:"required"`
}
type DeleteMChequeRejectParams struct {
	ChqRejectId int `params:"chq_reject_id" validate:"required"`
}
type UpdateMChequeRejectParams struct {
	ChqRejectId int `params:"chq_reject_id" validate:"required"`
}
