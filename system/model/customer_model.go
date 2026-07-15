package model

type Customer struct {
	CustId         string  `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	CustName       *string `gorm:"column:cust_name" json:"cust_name"`
	Street1        *string `gorm:"column:street1" json:"street1"`
	Street2        *string `gorm:"column:street2" json:"street2"`
	City           *string `gorm:"column:city" json:"city"`
	StateId        *int    `gorm:"column:state_id" json:"state_id"`
	CountryId      *int    `gorm:"column:country_id" json:"country_id"`
	ZipCode        *string `gorm:"column:zip_code" json:"zip_code"`
	ContactName    *string `gorm:"column:contact_name" json:"contact_name"`
	ContactEmail   *string `gorm:"column:contact_email" json:"contact_email"`
	ContactPhoneNo *string `gorm:"column:contact_phone_no" json:"contact_phone_no"`
	Notes          *string `gorm:"column:notes" json:"notes"`
	ParentCustId   *string `gorm:"column:parent_cust_id" json:"parent_cust_id"`
	Domain         *string `gorm:"column:domain" json:"domain"`
	DistPriceGrpId *int64  `gorm:"column:dist_price_grp_id" json:"dist_price_grp_id"`
	Npwp           *string `gorm:"column:npwp" json:"npwp"`
}

func (Customer) TableName() string {
	return "smc.m_customer"
}
