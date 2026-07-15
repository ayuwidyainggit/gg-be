package entity

type CreateWhSoBody struct {
	CustID     string              `json:"cust_id"`
	WhSoNo     string              `json:"wh_so_no"`
	WhSoDate   *string             `json:"wh_so_date"`
	TrCode     *string             `json:"tr_code" validate:"required,len=3"`
	WhID       *int64              `json:"wh_id"`
	SoType     *int64              `json:"so_type"`
	Notes      *string             `json:"notes"`
	DataStatus *int64              `json:"data_status"`
	CreatedBy  *int64              `json:"created_by"`
	Details    []CreateWhSoDetBody `json:"details"`
}
type DetailWhSoParams struct {
	WhSoNo string `params:"wh_so_no" validate:"required"`
}
type WhSoResponse struct {
	WhSoNo        string            `json:"wh_so_no"`
	WhSoDate      *string           `json:"wh_so_date"`
	TrCode        *string           `json:"tr_code"`
	WhID          *int64            `json:"wh_id"`
	WhCode        string            `json:"wh_code"`
	WhName        string            `json:"wh_name"`
	SoType        *int64            `json:"so_type"`
	Notes         *string           `json:"notes"`
	DataStatus    *int64            `json:"data_status"`
	UpdatedAt     string            `json:"updated_at"`
	UpdatedByName string            `json:"updated_by_name"`
	IsClosed      bool              `json:"is_closed"`
	ClosedBy      int64             `json:"closed_by"`
	ClosedByName  string            `json:"closed_by_name"`
	ClosedAt      string            `json:"closed_at"`
	Details       []WhSoDetResponse `json:"details"`
}
type WhSoListResponse struct {
	WhSoNo        string  `json:"wh_so_no"`
	WhSoDate      *string `json:"wh_so_date"`
	TrCode        *string `json:"tr_code"`
	WhID          *int64  `json:"wh_id"`
	WhCode        string  `json:"wh_code"`
	WhName        string  `json:"wh_name"`
	SoType        *int64  `json:"so_type"`
	Notes         *string `json:"notes"`
	DataStatus    *int64  `json:"data_status"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
}
type DeleteWhSoParams struct {
	WhSoNo string `params:"wh_so_no" validate:"required"`
}
type UpdateWhSoParams struct {
	WhSoNo string `params:"wh_so_no" validate:"required"`
}
type UpdateWhSoBody struct {
	CustID     string              `json:"cust_id"`
	WhSoNo     string              `json:"wh_so_no"`
	WhSoDate   *string             `json:"wh_so_date"`
	TrCode     *string             `json:"tr_code"`
	WhID       *int64              `json:"wh_id"`
	SoType     *int64              `json:"so_type"`
	Notes      *string             `json:"notes"`
	DataStatus *int64              `json:"data_status"`
	CreatedBy  *int64              `json:"created_by"`
	UpdatedBy  *int64              `json:"updated_by"`
	Details    []UpdateWhSoDetBody `json:"details"`
}
