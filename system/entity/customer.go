package entity

type CustomerResponse struct {
	CustId         string  `json:"cust_id"`
	CustName       *string `json:"cust_name"`
	Street1        *string `json:"street1"`
	Street2        *string `json:"street2"`
	City           *string `json:"city"`
	StateId        *int    `json:"state_id"`
	CountryId      *int    `json:"country_id"`
	ZipCode        *string `json:"zip_code"`
	ContactName    *string `json:"contact_name"`
	ContactEmail   *string `json:"contact_email"`
	ContactPhoneNo *string `json:"contact_phone_no"`
	Notes          *string `json:"notes"`
	ParentCustId   *string `json:"parent_cust_id"`
	Domain         *string `json:"domain"`
	DistPriceGrpId *int64  `json:"dist_price_grp_id"`
	Npwp           *string `json:"npwp"`
}
