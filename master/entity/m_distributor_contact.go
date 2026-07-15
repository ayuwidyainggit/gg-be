package entity

type DistributorContactCreate struct {
	DistributorContact []DistributorContact `json:"distributor_contact"`
}

type DistributorContact struct {
	DistributorContactId int64  `json:"distributor_contact_id"`
	DistributorId        int64  `json:"distributor_id"`
	ContactName          string `json:"contact_name" validate:"required,max=50"`
	JobTitle             string `json:"job_title" validate:"required,max=20"`
	PhoneNo              string `json:"phone_no" validate:"required,numeric,max=20"`
	IsWaNo               bool   `json:"is_wa_no"`
	WaNo                 string `json:"wa_no" validate:"required,numeric,max=20"`
	Email                string `json:"email" validate:"omitempty,email,max=100"`
	IdentityNo           string `json:"identity_no" validate:"required,max=20"`
	IdentityType         string `json:"identity_type" validate:"required,oneof='National ID' Passport 'Others ID'"`
}

type DistributorContactUpdates struct {
	DistributorContact []DistributorContactUpdate `json:"distributor_contact"`
}

type DistributorContactUpdate struct {
	CustId               string  `json:"cust_id"`
	DistributorContactId *int64  `json:"distributor_contact_id"`
	DistributorId        int     `json:"distributor_id"`
	ContactName          *string `json:"contact_name" validate:"required,max=50"`
	JobTitle             *string `json:"job_title" validate:"required,max=20"`
	PhoneNo              *string `json:"phone_no" validate:"required,numeric,max=20"`
	IsWaNo               *bool   `json:"is_wa_no"`
	WaNo                 *string `json:"wa_no" validate:"required,numeric,max=20"`
	Email                *string `json:"email" validate:"omitempty,email,max=100"`
	IdentityNo           *string `json:"identity_no" validate:"required,max=20"`
	IdentityType         *string `json:"identity_type" validate:"required,oneof='National ID' Passport 'Others ID'"`
}
