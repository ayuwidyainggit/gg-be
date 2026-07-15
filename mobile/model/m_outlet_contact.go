package model

type MOutletContact struct {
	CustID          string `db:"cust_id" json:"cust_id"`
	OutletID        int64  `db:"outlet_id" json:"outlet_id"`
	ContactName     string `db:"contact_name" json:"contact_name"`
	JobTitle        string `db:"job_title" json:"job_title"`
	PhoneNo         string `db:"phone_no" json:"phone_no"`
	WaNo            string `db:"wa_no" json:"wa_no"`
	IdentityNo      string `db:"identity_no" json:"identity_no"`
	Email           string `db:"email" json:"email"`
	OutletContactId *int64 `db:"outlet_contact_id" json:"outlet_contact_id"`
}
