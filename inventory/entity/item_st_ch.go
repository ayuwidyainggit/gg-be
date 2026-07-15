package entity

type CreateItemStChBody struct {
	CustID      string                  `json:"cust_id"`
	IscNo       string                  `json:"isc_no"`
	IscDate     *string                 `json:"isc_date"`
	TrCode      *string                 `json:"tr_code" validate:"required,len=3"`
	WhID        *int64                  `json:"wh_id"`
	ItemCdnFrom *int64                  `json:"item_cdn_from"`
	ItemCdnTo   *int64                  `json:"item_cdn_to"`
	Notes       *string                 `json:"notes"`
	DataStatus  *int64                  `json:"data_status"`
	CreatedBy   int64                   `json:"created_by"`
	Details     []CreateItemStChDetBody `json:"details" validate:"required,dive,required"`
}

type ItemStChResponse struct {
	IscNo         string                `json:"isc_no"`
	IscDate       *string               `json:"isc_date"`
	TrCode        *string               `json:"tr_code"`
	WhID          *int64                `json:"wh_id"`
	WhCode        string                `json:"wh_code"`
	WhName        string                `json:"wh_name"`
	ItemCdnFrom   *int64                `json:"item_cdn_from"`
	ItemCdnTo     *int64                `json:"item_cdn_to"`
	Notes         *string               `json:"notes"`
	DataStatus    *int64                `json:"data_status"`
	UpdatedAt     string                `json:"updated_at"`
	UpdatedByName string                `json:"updated_by_name"`
	IsClosed      bool                  `json:"is_closed"`
	ClosedBy      int64                 `json:"closed_by"`
	ClosedByName  string                `json:"closed_by_name"`
	ClosedAt      string                `json:"closed_at"`
	Details       []ItemStChDetResponse `json:"details"`
}

type ItemStChListResponse struct {
	IscNo         string  `json:"isc_no"`
	IscDate       *string `json:"isc_date"`
	TrCode        *string `json:"tr_code"`
	WhID          *int64  `json:"wh_id"`
	WhCode        string  `json:"wh_code"`
	WhName        string  `json:"wh_name"`
	ItemCdnFrom   *int64  `json:"item_cdn_from"`
	ItemCdnTo     *int64  `json:"item_cdn_to"`
	Notes         *string `json:"notes"`
	DataStatus    *int64  `json:"data_status"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
}

type DetailIscParams struct {
	IscNo string `params:"isc_no" validate:"required"`
}

type UpdateItemStChBody struct {
	CustID      string                  `json:"cust_id"`
	IscNo       string                  `json:"isc_no"`
	IscDate     *string                 `json:"isc_date"`
	WhID        *int64                  `json:"wh_id"`
	ItemCdnFrom *int64                  `json:"item_cdn_from"`
	ItemCdnTo   *int64                  `json:"item_cdn_to"`
	Notes       *string                 `json:"notes"`
	DataStatus  *int64                  `json:"data_status"`
	CreatedBy   *int64                  `json:"created_by"`
	UpdatedBy   int64                   `json:"updated_by"`
	Details     []UpdateItemStChDetBody `json:"details"`
}

type UpdateIscParams struct {
	IscNo string `params:"isc_no" validate:"required"`
}
