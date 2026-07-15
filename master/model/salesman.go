package model

import (
	"time"
)

type Salesman struct {
	CustId           string     `db:"cust_id" json:"cust_id"`
	EmpId            int64      `db:"emp_id" json:"emp_id"`
	SalesName        string     `db:"sales_name" json:"sales_name"`
	SalesTeamId      int64      `db:"sales_team_id" json:"sales_team_id"`
	OprType          string     `db:"opr_type" json:"opr_type"`
	IsBonusRep       bool       `db:"is_bonus_rep" json:"is_bonus_rep"`
	TransDate        *time.Time `db:"trans_date" json:"trans_date"`
	WhId             int64      `db:"wh_id" json:"wh_id"`
	WhIdTackingOrder int64      `db:"wh_id" json:"wh_id_tacking_order"`
	IncGrpId         int64      `db:"inc_grp_id" json:"inc_grp_id"`
	OfficialId       int64      `db:"official_id" json:"official_id"`
	SaleSystem       string     `db:"sale_system" json:"sale_system"`
	SmIsTransfer     bool       `db:"sm_is_transfer" json:"sm_is_transfer"`
	SmValidRoute     bool       `db:"sm_valid_route" json:"sm_valid_route"`
	SmGeolocValid    bool       `db:"sm_geoloc_valid" json:"sm_geoloc_valid"`
	SmRadius         int64      `db:"sm_radius" json:"sm_radius"`
	SmPassword       string     `db:"sm_password,omitempty" json:"sm_password"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	IsTakingOrder    bool       `db:"is_taking_order" json:"is_taking_order"`
	ImageUrl         string     `db:"image_url" json:"image_url"`
	SmIsBarcode      bool       `db:"sm_is_barcode" json:"sm_is_barcode"`
	SmIsPhotoProfile bool       `db:"sm_is_photo_profile" json:"sm_is_photo_profile"`
	JobType          string     `db:"job_type" json:"job_type"`
	AllowInputPrice  bool       `json:"allow_input_price" db:"allow_input_price" `
	TaxOption        string     `db:"tax_option" json:"tax_option"`
	StartDate        time.Time  `db:"start_date" json:"start_date"`
	EndDate          *time.Time `db:"end_date" json:"end_date"`
	IsDel            bool       `db:"is_del" json:"is_del"`
	CreatedBy        *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt        *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy        *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName    *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt        *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy        *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt        *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type SalesmanList struct {
	CustId           string     `db:"cust_id" json:"cust_id"`
	EmpId            int64      `db:"emp_id" json:"emp_id"`
	EmpCode          string     `db:"emp_code" json:"emp_code"`
	EmpName          string     `db:"emp_name" json:"emp_name"`
	Email            *string    `db:"email" json:"email"`
	PhoneNo          *string    `db:"phone_no" json:"phone_no"`
	LastEducation    *string    `db:"last_education" json:"last_education"`
	Address          *string    `db:"address" json:"address"`
	SalesName        string     `db:"sales_name" json:"sales_name"`
	SalesTeamId      int64      `db:"sales_team_id" json:"sales_team_id"`
	SalesTeamCode    *string    `db:"sales_team_code" json:"sales_team_code"`
	SalesTeamName    *string    `db:"sales_team_name" json:"sales_team_name"`
	OprType          string     `db:"opr_type" json:"opr_type"`
	OprTypeCanvas    *string    `db:"opr_type_canvas" json:"opr_type_canvas"`
	IsBonusRep       bool       `db:"is_bonus_rep" json:"is_bonus_rep"`
	TransDate        *time.Time `db:"trans_date" json:"trans_date"`
	WhId             int64      `db:"wh_id" json:"wh_id"`
	WhCode           *string    `db:"wh_code" json:"wh_code"`
	WhName           *string    `db:"wh_name" json:"wh_name"`
	WhNameCanvas     *string    `db:"wh_name_canvas" json:"wh_name_canvas"`
	WhCanvasID       *int64     `db:"wh_canvas_id" json:"wh_canvas_id"`
	VehicleId        *int64     `db:"vehicle_id" json:"vehicle_id"`
	VehicleName      *string    `db:"vehicle_name" json:"vehicle_name"`
	DriverName       *string    `db:"driver_name" json:"driver_name"`
	IncGrpId         int64      `db:"inc_grp_id" json:"inc_grp_id"`
	IncGrpName       *string    `db:"inc_grp_name" json:"inc_grp_name"`
	OfficialId       int64      `db:"official_id" json:"official_id"`
	OfficialName     *string    `db:"official_name" json:"official_name"`
	OfficialType     *int64     `db:"official_type" json:"official_type"`
	HierarchyCode    *string    `db:"hierarchy_code" json:"hierarchy_code"`
	SaleSystem       string     `db:"sale_system" json:"sale_system"`
	SmIsTransfer     bool       `db:"sm_is_transfer" json:"sm_is_transfer"`
	SmValidRoute     bool       `db:"sm_valid_route" json:"sm_valid_route"`
	SmGeolocValid    bool       `db:"sm_geoloc_valid" json:"sm_geoloc_valid"`
	SmRadius         int64      `db:"sm_radius" json:"sm_radius"`
	SmPassword       string     `db:"sm_password,omitempty" json:"sm_password"`
	ImageUrl         *string    `db:"image_url,omitempty" json:"image_url"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	UpdatedByName    *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt        *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	SmIsBarcode      bool       `db:"sm_is_barcode" json:"sm_is_barcode"`
	SmIsPhotoProfile bool       `db:"sm_is_photo_profile" json:"sm_is_photo_profile"`
	IsActiveCanvas   *bool      `db:"is_active_canvas" json:"is_active_canvas"`
	IsTakingOrder    *bool      `db:"is_taking_order" json:"is_taking_order"`
	AllowInputPrice  *bool      `db:"allow_input_price" json:"allow_input_price" `
	JobType          *string    `db:"job_type" json:"job_type"`
	TaxOption        *string    `db:"tax_option" json:"tax_option"`
	StartDate        *string    `db:"start_date" json:"start_date"`
	EndDate          *string    `db:"end_date" json:"end_date"`
}

type SalesmanUpdate struct {
	SalesName        *string    `json:"sales_name,omitempty" sql:"sales_name"`
	SalesTeamId      *int64     `json:"sales_team_id,omitempty" sql:"sales_team_id"`
	OprType          *string    `json:"opr_type" sql:"opr_type"`
	IsBonusRep       *bool      `json:"is_bonus_rep" sql:"is_bonus_rep"`
	TransDate        *string    `json:"trans_date" sql:"trans_date"`
	WhId             *int64     `json:"wh_id" sql:"wh_id"`
	IncGrpId         *int64     `json:"inc_grp_id" sql:"inc_grp_id"`
	OfficialId       *int64     `json:"official_id" sql:"official_id"`
	SalesSystem      *string    `json:"sales_system" sql:"sales_system"`
	SmIsTransfer     *bool      `json:"sm_is_transfer" sql:"sm_is_transfer"`
	SmValidRoute     *bool      `json:"sm_valid_route" sql:"sm_valid_route"`
	SmGeolocValid    *bool      `json:"sm_geoloc_valid" sql:"sm_geoloc_valid"`
	SmRadius         *int64     `json:"sm_radius" sql:"sm_radius"`
	SmPassword       *string    `json:"sm_password" sql:"sm_password"`
	ImageUrl         *string    `json:"image_url" sql:"image_url"`
	SmIsBarcode      *bool      `json:"sm_is_barcode" sql:"sm_is_barcode"`
	SmIsPhotoProfile *bool      `json:"sm_is_photo_profile" sql:"sm_is_photo_profile"`
	IsActive         *bool      `json:"is_active" sql:"is_active"`
	IsTakingOrder    *bool      `json:"is_taking_order" sql:"is_taking_order" `
	AllowInputPrice  *bool      `json:"allow_input_price" sql:"allow_input_price" `
	JobType          *string    `sql:"job_type" json:"job_type"`
	TaxOption        *string    `sql:"tax_option" json:"tax_option"`
	StartDate        *string    `sql:"start_date" json:"start_date"`
	EndDate          *string    `sql:"end_date" json:"end_date"`
	UpdatedAt        *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy        *int64     `json:"updated_by" sql:"updated_by"`
}

type SalesmanUpdateTakingOrder struct {
	OprType   *string    `json:"opr_type" sql:"opr_type"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}

type SalesmanDetailRead struct {
	MSalesmanProductTypeID int64  `db:"m_salesman_product_type_id" json:"m_salesman_product_type_id"`
	GroupType              int    `db:"group_type" json:"group_type"`
	RefID                  int64  `db:"ref_id" json:"ref_id"`
	RefCode                string `db:"ref_code" json:"ref_code"`
	RefName                string `db:"ref_name" json:"ref_name"`
	PlId                   int    `db:"pl_id" json:"pl_id"`
}

type SalesmanDetail struct {
	MSalesmanProductTypeID *int64  `db:"m_salesman_product_type_id" json:"m_salesman_product_type_id"`
	CustId                 *string `db:"cust_id" json:"cust_id"`
	EmpId                  *int64  `db:"emp_id" json:"emp_id"`
	GroupType              *int    `db:"group_type" json:"group_type"`
	RefID                  *int64  `db:"ref_id" json:"ref_id"`
}

type SalesmanCanvas struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	EmpId         int64      `db:"emp_id" json:"emp_id"`
	WhId          int        `db:"wh_id" json:"wh_id"`
	VehicleId     int64      `db:"vehicle_id" json:"vehicle_id"`
	IsActive      bool       `db:"is_active" json:"is_active_canvas"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	OprTypeCanvas string     `json:"opr_type_canvas" sql:"opr_type"`
}

type SalesmanCanvasUpdate struct {
	IsActive      *bool      `db:"is_active" json:"is_active_canvas"`
	VehicleId     *int64     `db:"vehicle_id" json:"vehicle_id"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}
