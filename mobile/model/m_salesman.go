package model

import (
	"time"

	"gorm.io/gorm"
)

type MSalesman struct {
	CustId        string     `gorm:"column:cust_id" json:"cust_id"`
	EmpId         int64      `gorm:"column:emp_id" json:"emp_id"`
	SalesName     string     `gorm:"column:sales_name" json:"sales_name"`
	SalesTeamId   int64      `gorm:"column:sales_team_id" json:"sales_team_id"`
	OprType       string     `gorm:"column:opr_type" json:"opr_type"`
	IsBonusRep    bool       `gorm:"column:is_bonus_rep" json:"is_bonus_rep"`
	TransDate     *time.Time `gorm:"column:trans_date" json:"trans_date"`
	WhId          int64      `gorm:"column:wh_id" json:"wh_id"`
	IncGrpId      int64      `gorm:"column:inc_grp_id" json:"inc_grp_id"`
	OfficialId    int64      `gorm:"column:official_id" json:"official_id"`
	SaleSystem    string     `gorm:"column:sale_system" json:"sale_system"`
	IsPharmacy    bool       `gorm:"column:is_pharmacy" json:"is_pharmacy"`
	SmIsTransfer  bool       `gorm:"column:sm_is_transfer" json:"sm_is_transfer"`
	SmValidRoute  bool       `gorm:"column:sm_valid_route" json:"sm_valid_route"`
	SmGeolocValid bool       `gorm:"column:sm_geoloc_valid" json:"sm_geoloc_valid"`
	SmRadius      int64      `gorm:"column:sm_radius" json:"sm_radius"`
	SmPassword    string     `gorm:"column:sm_password,omitempty" json:"sm_password"`
	IsActive      bool       `gorm:"column:is_active" json:"is_active"`
	IsDel         bool       `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64     `gorm:"column:created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `gorm:"column:updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `gorm:"column:updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `gorm:"column:deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `gorm:"column:deleted_at,omitempty" json:"deleted_at"`
}

func (MSalesman) TableName() string {
	return "mst.m_salesman"
}

func (m *MSalesman) BeforeUpdate(trx *gorm.DB) (err error) {
	timeNow := time.Now()
	m.UpdatedAt = &timeNow
	return nil
}

type MSalesmanRead struct {
	CustId               string     `gorm:"column:cust_id" json:"cust_id"`
	EmpId                int64      `gorm:"column:emp_id" json:"emp_id"`
	SalesName            string     `gorm:"column:sales_name" json:"sales_name"`
	SalesTeamId          int64      `gorm:"column:sales_team_id" json:"sales_team_id"`
	SalesTeamName        string     `gorm:"column:sales_team_name" json:"sales_team_name"`
	OprType              string     `gorm:"column:opr_type" json:"opr_type"`
	IsBonusRep           bool       `gorm:"column:is_bonus_rep" json:"is_bonus_rep"`
	TransDate            *time.Time `gorm:"column:trans_date" json:"trans_date"`
	WhId                 int64      `gorm:"column:wh_id" json:"wh_id"`
	IncGrpId             int64      `gorm:"column:inc_grp_id" json:"inc_grp_id"`
	OfficialId           int64      `gorm:"column:official_id" json:"official_id"`
	SaleSystem           string     `gorm:"column:sale_system" json:"sale_system"`
	IsPharmacy           bool       `gorm:"column:is_pharmacy" json:"is_pharmacy"`
	SmIsTransfer         bool       `gorm:"column:sm_is_transfer" json:"sm_is_transfer"`
	SmValidRoute         bool       `gorm:"column:sm_valid_route" json:"sm_valid_route"`
	SmGeolocValid        bool       `gorm:"column:sm_geoloc_valid" json:"sm_geoloc_valid"`
	SmRadius             int64      `gorm:"column:sm_radius" json:"sm_radius"`
	SmPassword           string     `gorm:"column:sm_password,omitempty" json:"sm_password"`
	IsActive             bool       `gorm:"column:is_active" json:"is_active"`
	IsDel                bool       `gorm:"column:is_del" json:"is_del"`
	CreatedBy            *int64     `gorm:"column:created_by,omitempty" json:"created_by"`
	CreatedAt            *time.Time `gorm:"column:created_at,omitempty" json:"created_at"`
	UpdatedBy            *int64     `gorm:"column:updated_by,omitempty" json:"updated_by"`
	UpdatedByName        *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt            *time.Time `gorm:"column:updated_at,omitempty" json:"updated_at"`
	DeletedBy            *int64     `gorm:"column:deleted_by,omitempty" json:"deleted_by"`
	DeletedAt            *time.Time `gorm:"column:deleted_at,omitempty" json:"deleted_at"`
	OprTypeOrderTaking   string     `gorm:"column:opr_type_order_taking" json:"opr_type_order_taking"`
	OprTypeCanvas        string     `gorm:"column:opr_type_canvas" json:"opr_type_canvas"`
	AllowInputPrice      bool       `gorm:"column:allow_input_price" json:"allow_input_price"`
	TaxOption            string     `gorm:"column:tax_option" json:"tax_option"`
	IsActiveGudangCanvas bool       `gorm:"column:is_active_gudang_canvas" json:"is_active_gudang_canvas"`
	IsActiveGudangUtama  bool       `gorm:"column:is_active_gudang_utama" json:"is_active_gudang_utama"`
}

func (MSalesmanRead) TableName() string {
	return "mst.m_salesman"
}

type PjpSalesmanDetail struct {
	CustId                 string `gorm:"column:cust_id" json:"cust_id"`
	CustName               string `gorm:"column:cust_name" json:"cust_name"`
	DistributorId          *int   `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorCode        string `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName        string `gorm:"column:distributor_name" json:"distributor_name"`
	EmpId                  int64  `gorm:"column:emp_id" json:"emp_id"`
	SalesName              string `gorm:"column:sales_name" json:"sales_name"`
	SalesTeamId            int64  `gorm:"column:sales_team_id" json:"sales_team_id"`
	SalesTeamCode          string `gorm:"column:sales_team_code" json:"sales_team_code"`
	SalesTeamName          string `gorm:"column:sales_team_name" json:"sales_team_name"`
	IsTakingOrder          bool   `gorm:"column:is_taking_order" json:"is_taking_order"`
	OprType                string `gorm:"column:opr_type" json:"opr_type"`
	WHID                   int64  `gorm:"column:wh_id" json:"wh_id"`
	IsActiveSalesmanCanvas bool   `gorm:"column:is_active_salesman_canvas" json:"is_active_salesman_canvas"`
	WHIDCanvas             int64  `gorm:"column:wh_id_canvas" json:"wh_id_canvas"`
	OprTypeCanvas          string `gorm:"column:opr_type_canvas" json:"opr_type_canvas"`
}

type PjpSalesmanWarehouse struct {
	OprType string `gorm:"column:opr_type" json:"opr_type"`
	WhId    int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode  string `gorm:"column:wh_code" json:"wh_code"`
	WhName  string `gorm:"column:wh_name" json:"wh_name"`
}
