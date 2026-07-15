package entity

// Query filter for GET /v1/outlet-list
type OutletListQueryFilter struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Sort         string `query:"sort"`
	OutletStatus []int  `query:"outlet_status"`
	IsActive     *int   `query:"is_active"`
}

// Response item for GET /v1/outlet-list
type OutletListItem struct {
	OutletId     int64  `json:"outlet_id"`
	OutletCode   string `json:"outlet_code"`
	OutletName   string `json:"outlet_name"`
	Address1     string `json:"address1"`
	Longitude    string `json:"longitude"`
	Latitude     string `json:"latitude"`
	OutletStatus int    `json:"outlet_status"`
}

// Path params for DELETE and PATCH
type OutletListParams struct {
	OutletId int64 `params:"outlet_id" validate:"required"`
}

// Request body for PATCH /v1/m-outlets/:outlet_id
type UpdateOutletBody struct {
	OutletId     int64  `json:"-"`
	CustId       string `json:"-"`
	ParentCustId string `json:"-"`
	UpdatedBy    int64  `json:"-"`

	OutletName  *string             `json:"outlet_name,omitempty" validate:"omitempty,max=150"`
	Address     *string             `json:"address,omitempty" validate:"omitempty,max=150"`
	PhoneNo     *string             `json:"phone_no,omitempty" validate:"omitempty,max=20"`
	BuildingOwn *int                `json:"building_own,omitempty"`
	Latitude    *string             `json:"latitude,omitempty"`
	Longitude   *string             `json:"longitude,omitempty"`
	FileUrl     *string             `json:"file_url,omitempty" validate:"omitempty,max=500"`
	Details     *UpdateOutletDetail `json:"details,omitempty"`
}

type UpdateOutletDetail struct {
	Contact []UpdateOutletContact `json:"contact,omitempty" validate:"dive"`
}

type UpdateOutletContact struct {
	OutletContactId *int64  `json:"outlet_contact_id,omitempty"`
	ContactName     *string `json:"contact_name,omitempty" validate:"omitempty,max=150"`
	JobTitle        *string `json:"job_title,omitempty" validate:"omitempty,max=100"`
	PhoneNo         *string `json:"phone_no,omitempty" validate:"omitempty,max=20"`
	WaNo            *string `json:"wa_no,omitempty" validate:"omitempty,max=20"`
	Email           *string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	IdentityNo      *string `json:"identity_no,omitempty" validate:"omitempty,max=100"`
}
