package response

import "time"

type SalesmanDetail struct {
	// Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id"`
}

type ListSalesmanResponse struct {
	EmpID             int             `json:"emp_id"`
	EmpCode           string          `json:"emp_code"`
	EmpName           string          `json:"emp_name"`
	SalesTeamID       int             `json:"sales_team_id"`
	SalesTeamCode     string          `json:"sales_team_code"`
	SalesTeamName     string          `json:"sales_team_name"`
	SalesName         string          `json:"sales_name"`
	OprType           string          `json:"opr_type"`
	OprTypeCanvas     string          `json:"opr_type_canvas"`
	OprTypeName       string          `json:"opr_type_name"`
	OprTypeNameCanvas string          `json:"opr_type_name_canvas"`
	IsBonusRep        bool            `json:"is_bonus_rep"`
	TransDate         *time.Time      `json:"trans_date"` // nullable
	WhID              int             `json:"wh_id"`
	WhCode            string          `json:"wh_code"`
	WhName            string          `json:"wh_name"`
	IncGrpID          int             `json:"inc_grp_id"`
	OfficialID        int             `json:"official_id"`
	OfficialType      int             `json:"official_type"`
	HierarchyCode     string          `json:"hierarchy_code"`
	SaleSystem        string          `json:"sale_system"`
	SmIsTransfer      bool            `json:"sm_is_transfer"`
	SmValidRoute      bool            `json:"sm_valid_route"`
	SmGeolocValid     bool            `json:"sm_geoloc_valid"`
	SmRadius          int             `json:"sm_radius"`
	SmPassword        string          `json:"sm_password"`
	IsActive          bool            `json:"is_active"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedByName     string          `json:"updated_by_name"`
	ImageURL          string          `json:"image_url"`
	Details           SalesmanDetails `json:"details"`
	SmIsBarcode       bool            `json:"sm_is_barcode"`
	SmIsPhotoProfile  bool            `json:"sm_is_photo_profile"`
	IsActiveCanvas    bool            `json:"is_active_canvas"`
	IsTakingOrder     bool            `json:"is_taking_order"`
}

type SalesmanDetails struct {
	ProductLine []ProductLine `json:"product_line"`
	Brand       []Brand       `json:"brand"`
	SubBrand    []SubBrand    `json:"sub_brand"`
}

type ProductLine struct {
	PlID                   int    `json:"pl_id"`
	MSalesmanProductTypeID int    `json:"m_salesman_product_type_id"`
	RefID                  int    `json:"ref_id"`
	PlCode                 string `json:"pl_code"`
	PlName                 string `json:"pl_name"`
}

type Brand struct {
	PlID                   int    `json:"pl_id"`
	MSalesmanProductTypeID int    `json:"m_salesman_product_type_id"`
	RefID                  int    `json:"ref_id"`
	BrandCode              string `json:"brand_code"`
	BrandName              string `json:"brand_name"`
}

type SubBrand struct {
	PlID                   int    `json:"pl_id"`
	MSalesmanProductTypeID int    `json:"m_salesman_product_type_id"`
	RefID                  int    `json:"ref_id"`
	SBrand1Code            string `json:"sbrand1_code"`
	SBrand1Name            string `json:"sbrand1_name"`
}

type ListSalesmanAPIResponse struct {
	Message   string                 `json:"message"`
	Data      []ListSalesmanResponse `json:"data"`
	Paging    PagingInfo             `json:"paging"`
	RequestID string                 `json:"request_id"`
}

type PagingInfo struct {
	TotalRecord int `json:"total_record"`
	PageCurrent int `json:"page_current"`
	PageLimit   int `json:"page_limit"`
	PageTotal   int `json:"page_total"`
}
