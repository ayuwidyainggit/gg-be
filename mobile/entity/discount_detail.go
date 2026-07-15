package entity

type DiscountPrincipal struct {
	CustID        string `json:"cust_id,omitempty"`
	DiscountID    string `json:"discount_id,omitempty" validate:""`
	PrincipalID   int64  `json:"principal_id" validate:"required"`
	PrincipalCode string `json:"principal_code,omitempty"`
	PrincipalName string `json:"principal_name,omitempty"`
}

type DiscountGroup struct {
	CustID      string `json:"cust_id,omitempty"`
	DiscountID  string `json:"discount_id,omitempty" validate:""`
	DiscGrpID   int    `json:"disc_grp_id" validate:"required"`
	DiscGrpCode string `json:"disc_grp_code,omitempty"`
	DiscGrpName string `json:"disc_grp_name,omitempty"`
}

type DetailDiscountGrp struct {
	OutletId   int    `json:"outlet_id"`
	OutletCode string `json:"outlet_code"`
	OutletName string `json:"outlet_name"`
	Address    string `json:"address"`
	PhoneNo    string `json:"phone_no"`
	WaNo       string `json:"wa_no"`
	FaxNo      string `json:"fax_no"`
	Email      string `json:"email"`
}
