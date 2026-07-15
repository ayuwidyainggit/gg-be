package entity

import (
	"time"
)

type SalesmanQueryFilter struct {
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	Mode          string `query:"mode"`
	Sort          string `query:"sort"`
	IsActive      *int   `query:"is_active"`
	SalesTeamId   string   `query:"sales_team_id"`
	DistributorID []int    `query:"-"`
	CustIds       []string `query:"-"`
}

type SalesmanResponse struct {
	EmpId         int64   `json:"emp_id"`
	EmpCode       string  `json:"emp_code"`
	EmpName       string  `json:"emp_name"`
	Email         string  `json:"email"`
	PhoneNo       string  `json:"phone_no"`
	LastEducation string  `json:"last_education"`
	Address       string  `json:"address"`
	SalesTeamId   int64   `json:"sales_team_id"`
	SalesTeamCode string  `json:"sales_team_code"`
	SalesTeamName string  `json:"sales_team_name"`
	SalesName     string  `json:"sales_name"`
	OprType       string  `json:"opr_type"`
	OprTypeCanvas string  `json:"opr_type_canvas"`
	IsBonusRep    bool    `json:"is_bonus_rep"`
	TransDate     *string `json:"trans_date"`
	// WhId             int64            `json:"wh_id"`
	WhId             int              `json:"wh_id"`
	WhCode           string           `json:"wh_code"`
	WhName           string           `json:"wh_name"`
	WhNameCanvas     string           `json:"wh_name_canvas"`
	VehicleId        int64            `json:"vehicle_id"`
	VehicleName      string           `json:"vehicle_name"`
	DriverName       string           `json:"driver_name"`
	IncGrpId         int64            `json:"inc_grp_id"`
	IncGrpName       *string          `json:"inc_grp_name"`
	OfficialId       int64            `json:"official_id"`
	OfficialName     *string          `json:"official_name"`
	OfficialType     int64            `json:"official_type"`
	HierarchyCode    string           `json:"hierarchy_code"`
	SaleSystem       string           `json:"sale_system"`
	SmIsTransfer     bool             `json:"sm_is_transfer"`
	SmValidRoute     bool             `json:"sm_valid_route"`
	SmGeolocValid    bool             `json:"sm_geoloc_valid"`
	SmRadius         int64            `json:"sm_radius"`
	SmPassword       string           `json:"sm_password"`
	ImageUrl         string           `json:"image_url"`
	IsActive         *bool            `json:"is_active"`
	IsActiveCanvas   *bool            `json:"is_active_canvas"`
	UpdatedAt        *time.Time       `json:"updated_at"`
	UpdatedByName    string           `json:"updated_by_name"`
	Details          SalesmanDetGroup `json:"details"`
	SmIsBarcode      bool             `json:"sm_is_barcode"`
	SmIsPhotoProfile bool             `json:"sm_is_photo_profile"`
	IsTakingOrder    bool             `json:"is_taking_order"`
	JobType          *string          `json:"job_type"`
	JobTypeName      *string          `json:"job_type_name"`
	TaxOption        *string          `json:"tax_option"`
	TaxOptionName    *string          `json:"tax_option_name"`
	AllowInputPrice  *bool            `json:"allow_input_price"`
	StartDate        *string          `json:"start_date"`
	EndDate          *string          `json:"end_date"`
}

type SalesmanListResponse struct {
	EmpId             int64            `json:"emp_id"`
	EmpCode           string           `json:"emp_code"`
	EmpName           string           `json:"emp_name"`
	SalesTeamId       int64            `json:"sales_team_id"`
	SalesTeamCode     string           `json:"sales_team_code"`
	SalesTeamName     string           `json:"sales_team_name"`
	SalesName         string           `json:"sales_name"`
	OprType           string           `json:"opr_type"`
	OprTypeCanvas     string           `json:"opr_type_canvas"`
	OprTypeText       string           `json:"opr_type_name"`
	OprTypeTextCanvas *string          `json:"opr_type_name_canvas"`
	IsBonusRep        bool             `json:"is_bonus_rep"`
	TransDate         *string          `json:"trans_date"`
	WhId              int64            `json:"wh_id"`
	WhCode            *string          `json:"wh_code"`
	WhName            string           `json:"wh_name"`
	WhCanvasID        *int64           `json:"wh_canvas_id"`
	WhNameCanvas      string           `json:"wh_name_canvas"`
	WhNameView        string           `json:"wh_name_view"`
	IncGrpId          int64            `json:"inc_grp_id"`
	OfficialId        int64            `json:"official_id"`
	OfficialType      int64            `json:"official_type"`
	HierarchyCode     string           `json:"hierarchy_code"`
	SaleSystem        string           `json:"sale_system"`
	SmIsTransfer      bool             `json:"sm_is_transfer"`
	SmValidRoute      bool             `json:"sm_valid_route"`
	SmGeolocValid     bool             `json:"sm_geoloc_valid"`
	SmRadius          int64            `json:"sm_radius"`
	SmPassword        string           `json:"sm_password"`
	IsActive          bool             `json:"is_active"`
	UpdatedAt         *time.Time       `json:"updated_at"`
	UpdatedByName     string           `json:"updated_by_name"`
	ImageUrl          string           `json:"image_url"`
	Details           SalesmanDetGroup `json:"details"`
	SmIsBarcode       bool             `json:"sm_is_barcode"`
	SmIsPhotoProfile  bool             `json:"sm_is_photo_profile"`
	IsActiveCanvas    *bool            `json:"is_active_canvas"`
	IsTakingOrder     *bool            `json:"is_taking_order"`
}

type SalesmanLookupResponse struct {
	EmpId            int64      `json:"emp_id"`
	SalesTeamId      int64      `json:"sales_team_id"`
	SalesName        string     `json:"sales_name"`
	OprType          string     `json:"opr_type"`
	IsBonusRep       bool       `json:"is_bonus_rep"`
	TransDate        *string    `json:"trans_date"`
	WhId             int64      `json:"wh_id"`
	IncGrpId         int64      `json:"inc_grp_id"`
	OfficialId       int64      `json:"official_id"`
	SaleSystem       string     `json:"sale_system"`
	SmIsTransfer     bool       `json:"sm_is_transfer"`
	SmValidRoute     bool       `json:"sm_valid_route"`
	SmGeolocValid    bool       `json:"sm_geoloc_valid"`
	SmRadius         int64      `json:"sm_radius"`
	SmPassword       string     `json:"sm_password"`
	IsActive         bool       `json:"is_active"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedByName    string     `json:"updated_by_name"`
	ImageUrl         string     `json:"image_url"`
	SmIsBarcode      bool       `json:"sm_is_barcode"`
	SmIsPhotoProfile bool       `json:"sm_is_photo_profile"`
}

type CreateSalesmanBody struct {
	CustId       string
	ParentCustId string
	CreatedBy    int64   `json:"created_by" validate:"required"`
	EmpId        int64   `json:"emp_id" validate:"required"`
	SalesName    string  `json:"sales_name" validate:"required,max=150"`
	SalesTeamId  int64   `json:"sales_team_id" validate:"required"`
	OprType      string  `json:"opr_type" validate:"oneof='C' 'O' 'S'"`
	IsBonusRep   bool    `json:"is_bonus_rep"`
	TransDate    *string `json:"trans_date"`
	// WhId             int64                     `json:"wh_id" validate:"required"`
	WhId             int64                     `json:"wh_id"`
	WhIdTackingOrder int64                     `json:"wh_id_tacking_order"`
	WarehouseCode    string                    `json:"wh_code" validate:"max=3"`
	IncGrpId         int64                     `json:"inc_grp_id"`
	OfficialId       int64                     `json:"official_id"`
	SaleSystem       string                    `json:"sale_system" validate:"oneof='N' 'K'"`
	SmIsTransfer     bool                      `json:"sm_is_transfer"`
	SmValidRoute     bool                      `json:"sm_valid_route"`
	SmGeolocValid    bool                      `json:"sm_geoloc_valid"`
	SmRadius         int64                     `json:"sm_radius"`
	SmPassword       string                    `json:"sm_password"`
	IsActive         bool                      `json:"is_active"`
	IsActiveCanvas   bool                      `json:"is_active_canvas"`
	ImageUrl         string                    `json:"image_url"`
	SmIsBarcode      bool                      `json:"sm_is_barcode"`
	SmIsPhotoProfile bool                      `json:"sm_is_photo_profile"`
	WarehouseName    string                    `json:"wh_name"`
	IsCanvas         bool                      `json:"is_canvas"`
	IsTakingOrder    bool                      `json:"is_taking_order"`
	VehicleId        int64                     `json:"vehicle_id"`
	OprTypeCanvas    string                    `json:"opr_type_canvas"`
	JobType          string                    `json:"job_type"`
	TaxOption        string                    `json:"tax_option"`
	AllowInputPrice  bool                      `json:"allow_input_price"`
	StartDate        string                    `json:"start_date"`
	EndDate          *string                   `json:"end_date"`
	Details          SalesmanDetCreateDetGroup `json:"details"`
}

type DetailSalesmanParams struct {
	CustId       string
	ParentCustId string
	EmpId        int64 `params:"emp_id" validate:"required"`
}

type SalesTeamIdParams struct {
	SalesTeamId int64 `params:"sales_team_id" validate:"required"`
}

type UpdateSalesmanParams struct {
	EmpId int64 `params:"emp_id" validate:"required"`
}

type DeleteSalesmanParams struct {
	EmpId int64 `params:"emp_id" validate:"required"`
}

type UpdateSalesmanRequest struct {
	CustId       string
	ParentCustId string
	UpdatedBy    int64  `json:"updated_by" validate:"required"`
	SalesName    string `json:"sales_name" validate:"required,max=150"`
	SalesTeamId  int64  `json:"sales_team_id"`
	// OprType          string                 `json:"opr_type" validate:"oneof='C' 'O' 'S'"`
	OprType          string                 `json:"opr_type"`
	IsBonusRep       bool                   `json:"is_bonus_rep"`
	TransDate        string                 `json:"trans_date"`
	WhId             int64                  `json:"wh_id"`
	IncGrpId         int64                  `json:"inc_grp_id"`
	OfficialId       int64                  `json:"official_id"`
	SaleSystem       string                 `json:"sale_system" validate:"oneof='N' 'K'"`
	SmIsTransfer     bool                   `json:"sm_is_transfer"`
	SmValidRoute     bool                   `json:"sm_valid_route"`
	SmGeolocValid    bool                   `json:"sm_geoloc_valid"`
	SmRadius         int64                  `json:"sm_radius"`
	SmPassword       string                 `json:"sm_password"`
	IsActive         *bool                  `json:"is_active,omitempty"`
	IsActiveCanvas   *bool                  `json:"is_active_canvas,omitempty"`
	IsTakingOrder    bool                   `json:"is_taking_order"`
	ImageUrl         string                 `json:"image_url"`
	SmIsBarcode      bool                   `json:"sm_is_barcode"`
	SmIsPhotoProfile bool                   `json:"sm_is_photo_profile"`
	WarehouseName    string                 `json:"wh_name"`
	VehicleId        int64                  `json:"vehicle_id"`
	JobType          string                 `json:"job_type"`
	TaxOption        string                 `json:"tax_option"`
	AllowInputPrice  bool                   `json:"allow_input_price"`
	StartDate        *string                `json:"start_date"`
	EndDate          *string                `json:"end_date"`
	Details          SalesmanDetGroupUpdate `json:"details"`
}

type UpdateSalesmanCanvasRequest struct {
	CustId         string
	ParentCustId   string
	UpdatedBy      int64 `json:"updated_by" validate:"required"`
	IsActiveCanvas *bool `json:"is_active_canvas,omitempty"`
}

var OprType = map[string]string{
	"C": "Canvas",
	"O": "Taking Order",
	"S": "Shop Sales",
}

func ConvStringString(data map[string]string, param string) string {
	result, ok := data[string(param)]
	if !ok {
		result = "Unknown"
	}
	return result
}

type UpdateTakingOrder struct {
	CustId       string
	OprType      string `json:"opr_type"`
	ParentCustId string
	UpdatedBy    int64 `json:"updated_by" validate:"required"`
}

type JobTypeLookupResponse struct {
	JobTypeId   string `json:"job_type"`
	JobTypeName string `json:"job_type_name"`
}

var JobType = []JobTypeLookupResponse{
	{JobTypeId: "P", JobTypeName: "Permanent"},
	{JobTypeId: "T", JobTypeName: "Temporary"},
}

type TaxOptionLookupResponse struct {
	TaxOptionId   string `json:"tax_option"`
	TaxOptionName string `json:"tax_option_name"`
}

var TaxOption = []TaxOptionLookupResponse{
	{TaxOptionId: "I", TaxOptionName: "Include Tax"},
	{TaxOptionId: "E", TaxOptionName: "Exclude Tax"},
}

func GetTaxOptionName(id string) string {
	for _, option := range TaxOption {
		if option.TaxOptionId == id {
			return option.TaxOptionName
		}
	}
	return "-"
}

func GetJobTypeName(id string) string {
	for _, option := range JobType {
		if option.JobTypeId == id {
			return option.JobTypeName
		}
	}
	return "-"
}

type UpdateIsActiveRequest struct {
	CustId string `json:"cust_id"`
	UserId int64  `json:"user_id"`
	EmpId  int64  `json:"emp_id"`
}

type CustomSchedulerRequest struct {
	StartDate string `json:"start_date"`
	Url       string `json:"url"`
}
