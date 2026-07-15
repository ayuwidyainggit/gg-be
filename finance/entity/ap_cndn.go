package entity

type CreateApCndnBody struct {
	CustID     string   `json:"cust_id"`
	ApCndnDate *string  `json:"ap_cndn_date"`
	TrCode     *string  `json:"tr_code"`
	DocNo      *string  `json:"doc_no"`
	CndnType   *string  `json:"cndn_type"`
	CndnID     *int64   `json:"cndn_id"`
	SupID      *int64   `json:"sup_id"`
	ApNo       *string  `json:"ap_no"`
	CndnValue  *float64 `json:"cndn_value"`
	CndnUsed   *float64 `json:"cndn_used"`
	Notes      *string  `json:"notes"`
	DataStatus *int64   `json:"data_status"`
	CreatedBy  *int64   `json:"created_by"`
}
type ApCndnResponse struct {
	ApCndnNo      string   `json:"ap_cndn_no"`
	ApCndnDate    *string  `json:"ap_cndn_date"`
	TrCode        *string  `json:"tr_code"`
	DocNo         *string  `json:"doc_no"`
	CndnType      *string  `json:"cndn_type"`
	CndnID        *int64   `json:"cndn_id"`
	CndnCode      string   `json:"cndn_code"`
	CndnName      string   `json:"cndn_name"`
	SupID         *int64   `json:"sup_id"`
	SupCode       *string  `json:"sup_code"`
	SupName       *string  `json:"sup_name"`
	ApNo          *string  `json:"ap_no"`
	CndnValue     *float64 `json:"cndn_value"`
	CndnUsed      *float64 `json:"cndn_used"`
	Notes         *string  `json:"notes"`
	DataStatus    *int64   `json:"data_status"`
	IsPosted      bool     `json:"is_posted"`
	PostedAt      *string  `json:"posted_at"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
}
type ApCndnListResponse struct {
	ApCndnNo      string   `json:"ap_cndn_no"`
	ApCndnDate    *string  `json:"ap_cndn_date"`
	TrCode        *string  `json:"tr_code"`
	DocNo         *string  `json:"doc_no"`
	CndnType      *string  `json:"cndn_type"`
	CndnID        *int64   `json:"cndn_id"`
	CndnCode      string   `json:"cndn_code"`
	CndnName      string   `json:"cndn_name"`
	SupID         *int64   `json:"sup_id"`
	SupCode       *string  `json:"sup_code"`
	SupName       *string  `json:"sup_name"`
	ApNo          *string  `json:"ap_no"`
	CndnValue     *float64 `json:"cndn_value"`
	CndnUsed      *float64 `json:"cndn_used"`
	Notes         *string  `json:"notes"`
	DataStatus    *int64   `json:"data_status"`
	IsPosted      bool     `json:"is_posted"`
	PostedAt      *string  `json:"posted_at"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
}
type UpdateApCndnBody struct {
	CustID     string   `json:"cust_id"`
	ApCndnDate *string  `json:"ap_cndn_date"`
	TrCode     *string  `json:"tr_code"`
	DocNo      *string  `json:"doc_no"`
	CndnType   *string  `json:"cndn_type"`
	CndnID     *int64   `json:"cndn_id"`
	SupID      *int64   `json:"sup_id"`
	ApNo       *string  `json:"ap_no"`
	CndnValue  *float64 `json:"cndn_value"`
	CndnUsed   *float64 `json:"cndn_used"`
	Notes      *string  `json:"notes"`
	DataStatus *int64   `json:"data_status"`
	UpdatedBy  int64    `json:"updated_by"`
}
type DetailApCndnParams struct {
	ApCndnNo string `params:"ap_cndn_no" validate:"required" json:"ap_cndn_no"`
}
type DeleteApCndnParams struct {
	ApCndnNo string `params:"ap_cndn_no" validate:"required" json:"ap_cndn_no"`
}
type UpdateApCndnParams struct {
	ApCndnNo string `params:"ap_cndn_no" validate:"required" json:"ap_cndn_no"`
}
