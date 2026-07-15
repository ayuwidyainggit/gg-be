package model

type DistributorContact struct {
	CustId               string  `db:"cust_id" json:"cust_id"`
	DistributorContactId int     `db:"distributor_contact_id" json:"distributor_contact_id"`
	DistributorId        int     `db:"distributor_id" json:"distributor_id"`
	ContactName          string  `db:"contact_name" json:"contact_name"`
	JobTitle             string  `db:"job_title" json:"job_title"`
	PhoneNo              string  `db:"phone_no" json:"phone_no"`
	IsWaNo               bool    `db:"is_wa_no" json:"is_wa_no"`
	WaNo                 *string `db:"wa_no" json:"wa_no"`
	Email                string  `db:"email" json:"email"`
	IdentityNo           string  `db:"identity_no" json:"identity_no"`
	IdentityType         *string `db:"identity_type" json:"identity_type"`
}

type DistributorContactUpdate struct {
	ContactName  *string `db:"contact_name" json:"contact_name"`
	JobTitle     *string `db:"job_title" json:"job_title"`
	PhoneNo      *string `db:"phone_no" json:"phone_no"`
	IsWaNo       *bool   `db:"is_wa_no" json:"is_wa_no"`
	WaNo         *string `db:"wa_no" json:"wa_no"`
	Email        *string `db:"email" json:"email"`
	IdentityNo   *string `db:"identity_no" json:"identity_no"`
	IdentityType *string `db:"identity_type" json:"identity_type"`
}
