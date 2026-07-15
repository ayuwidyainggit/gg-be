package model

import "time"

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
	IsDel            bool       `db:"is_del" json:"is_del"`
	CreatedBy        *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt        *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy        *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName    *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt        *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy        *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt        *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

func (Salesman) TableName() string {
	return "mst.m_salesman"
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

func (SalesmanCanvas) TableName() string {
	return "mst.m_salesman_canvas"
}
