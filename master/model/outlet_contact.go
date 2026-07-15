package model

type MOutletContact struct {
	CustID                  string `db:"cust_id" json:"cust_id"`
	OutletID                int64  `db:"outlet_id" json:"outlet_id"`
	ContactName             string `db:"contact_name" json:"contact_name"`
	JobTitle                string `db:"job_title" json:"job_title"`
	PhoneNo                 string `db:"phone_no" json:"phone_no"`
	WaNo                    string `db:"wa_no" json:"wa_no"`
	IdentityNo              string `db:"identity_no" json:"identity_no"`
	IsWaNo                  bool   `json:"is_wa_no" db:"is_wa_no"`
	Email                   string `db:"email" json:"email"`
	OutletContactId         *int64 `db:"outlet_contact_id" json:"outlet_contact_id"`
	IdentityType            string `db:"identity_no" json:"identity_type"`
	FaxNumber               string `db:"fax_number" json:"fax_number"`
	OutletEstablishmentDate string `db:"outlet_establishment_date" json:"outlet_establishment_date"`
}

type OutletContactList struct {
	ContactName             string `db:"contact_name" json:"contact_name"`
	JobTitle                string `db:"job_title" json:"job_title"`
	PhoneNo                 string `db:"phone_no" json:"phone_no"`
	WaNo                    string `db:"wa_no" json:"wa_no"`
	Email                   string `db:"email" json:"email"`
	IdentityNo              string `db:"identity_no" json:"identity_no"`
	IsWaNo                  bool   `json:"is_wa_no" db:"is_wa_no"`
	OutletContactId         *int64 `db:"outlet_contact_id" json:"outlet_contact_id"`
	IdentityType            string `db:"identity_no" json:"identity_type"`
	FaxNumber               string `db:"fax_number" json:"fax_number"`
	OutletEstablishmentDate string `db:"outlet_establishment_date" json:"outlet_establishment_date"`
}

type MOutletContactUpdate struct {
	ContactName             *string `db:"contact_name" json:"contact_name" sql:"contact_name"`
	JobTitle                *string `db:"job_title" json:"job_title" sql:"job_title"`
	PhoneNo                 *string `db:"phone_no" json:"phone_no" sql:"phone_no"`
	WaNo                    *string `db:"wa_no" json:"wa_no" sql:"wa_no"`
	IdentityNo              *string `db:"identity_no" json:"identity_no" sql:"identity_no"`
	IsWaNo                  *bool   `json:"is_wa_no" db:"is_wa_no" sql:"is_wa_no"`
	Email                   *string `db:"email" json:"email" sql:"email"`
	IdentityType            *string `db:"identity_type" json:"identity_type" sql:"identity_type"`
	FaxNumber               *string `db:"fax_number" json:"fax_number" sql:"fax_number"`
	OutletEstablishmentDate *string `db:"outlet_establishment_date" json:"outlet_establishment_date" sql:"outlet_establishment_date"`
}
type MOutletContactRead struct {
	CustID                  string  `db:"cust_id" json:"cust_id"`
	OutletID                int64   `db:"outlet_id" json:"outlet_id"`
	ContactName             *string `db:"contact_name" json:"contact_name"`
	JobTitle                *string `db:"job_title" json:"job_title"`
	PhoneNo                 *string `db:"phone_no" json:"phone_no"`
	WaNo                    *string `db:"wa_no" json:"wa_no"`
	Email                   *string `db:"email" json:"email"`
	IdentityNo              *string `db:"identity_no" json:"identity_no"`
	IsWaNo                  *bool   `json:"is_wa_no" db:"is_wa_no"`
	OutletContactId         *int64  `db:"outlet_contact_id" json:"outlet_contact_id"`
	IdentityType            *string `db:"identity_type" json:"identity_type"`
	FaxNumber               *string `db:"fax_number" json:"fax_number"`
	OutletEstablishmentDate *string `db:"outlet_establishment_date" json:"outlet_establishment_date"`
}
