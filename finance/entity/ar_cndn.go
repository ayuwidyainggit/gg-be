package entity

type CreateArCndnBody struct {
	CustID     string   `json:"cust_id"`
	ArCndnNo   string   `json:"ar_cndn_no"`
	ArCndnDate *string  `json:"ar_cndn_date"`
	TrCode     *string  `json:"tr_code"`
	DocNo      *string  `json:"doc_no"`
	CndnType   *string  `json:"cndn_type" validate:"required,oneof='CN' 'DN'"`
	CndnId     *int64   `json:"cndn_id"`
	OutletId   *int64   `json:"outlet_id"`
	CndnValue  *float64 `json:"cndn_value"`
	CndnUsed   *float64 `json:"cndn_used"`
	Notes      *string  `json:"notes"`
	DataStatus *int64   `json:"data_status"`
	IsPosted   bool     `json:"is_posted"`
	CreatedBy  *int64   `json:"created_by"`
}

type ArCndnResponse struct {
	ArCndnId      int      `json:"ar_cndn_id"`
	ArCndnNo      string   `json:"ar_cndn_no"`
	ArCndnDate    *string  `json:"ar_cndn_date"`
	TrCode        *string  `json:"tr_code"`
	DocNo         *string  `json:"doc_no"`
	CndnType      *string  `json:"cndn_type" validate:"required,oneof='CN' 'DN'"`
	CndnId        *int64   `json:"cndn_id"`
	CndnCode      string   `json:"cndn_code"`
	CndnName      string   `json:"cndn_name"`
	OutletId      *int64   `json:"outlet_id"`
	OutletCode    string   `json:"outlet_code"`
	OutletName    string   `json:"outlet_name"`
	CndnValue     *float64 `json:"cndn_value"`
	CndnUsed      *float64 `json:"cndn_used"`
	Notes         *string  `json:"notes"`
	DataStatus    *int64   `json:"data_status"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
	IsPosted      bool     `json:"is_posted"`
	PostedAt      *string  `json:"posted_at"`
}

type ArCndnListResponse struct {
	ArCndnNo      string   `json:"ar_cndn_no"`
	ArCndnId      *int     `json:"ar_cndn_id"`
	ArCndnDate    *string  `json:"ar_cndn_date"`
	TrCode        *string  `json:"tr_code"`
	DocNo         *string  `json:"doc_no"`
	CndnType      *string  `json:"cndn_type" validate:"required,oneof='CN' 'DN'"`
	CndnId        *int64   `json:"cndn_id"`
	CndnCode      string   `json:"cndn_code"`
	CndnName      string   `json:"cndn_name"`
	OutletId      *int64   `json:"outlet_id"`
	OutletCode    string   `json:"outlet_code"`
	OutletName    string   `json:"outlet_name"`
	CndnValue     *float64 `json:"cndn_value"`
	CndnUsed      *float64 `json:"cndn_used"`
	Notes         *string  `json:"notes"`
	DataStatus    *int64   `json:"data_status"`
	CreatedBy     *int64   `json:"created_by"`
	UpdatedAt     string   `json:"updated_at"`
	UpdatedByName string   `json:"updated_by_name"`
	IsPosted      *bool    `json:"is_posted"`
}

type UpdateArCndnBody struct {
	CustID     string   `json:"cust_id"`
	ArCndnNo   string   `json:"ar_cndn_no"`
	ArCndnId   *int     `json:"ar_cndn_id"`
	ArCndnDate *string  `json:"ar_cndn_date"`
	TrCode     *string  `json:"tr_code"`
	DocNo      *string  `json:"doc_no"`
	CndnType   *string  `json:"cndn_type" validate:"required,oneof='CN' 'DN'"`
	CndnId     *int64   `json:"cndn_id"`
	OutletId   *int64   `json:"outlet_id"`
	CndnValue  *float64 `json:"cndn_value"`
	CndnUsed   *float64 `json:"cndn_used"`
	Notes      *string  `json:"notes"`
	DataStatus *int64   `json:"data_status"`
	UpdatedBy  int64    `json:"updated_by"`
	IsPosted   *bool    `json:"is_posted"`
}

type DetailArCndnParams struct {
	ArCndnId int `params:"ar_cndn_id" validate:"required"`
}
type DeleteArCndnParams struct {
	ArCndnId int `params:"ar_cndn_id" validate:"required"`
}
type UpdateArCndnParams struct {
	ArCndnId int `params:"ar_cndn_id" validate:"required"`
}
