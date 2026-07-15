package model

import "time"

type SalesTeam struct {
	SalesTeamID   int    `json:"sales_team_id"`
	SalesTeamCode string `json:"sales_team_code"`
	SalesTeamName string `json:"sales_team_name"`
}

type Salesman struct {
	EmployeeID    int    `json:"emp_id"`
	EmployeeCode  string `json:"emp_code"`
	EmployeeName  string `json:"emp_name"`
	SalesTeamID   int    `json:"sales_team_id"`
	SalesTeamCode string `json:"sales_team_code"`
	SalesTeamName string `json:"sales_team_name"`
	SalesName     string `json:"sales_name"`
}

type NewSalesman struct {
	EmployeeID          int    `json:"emp_id"`
	EmployeeCode        string `json:"emp_code"`
	EmployeeName        string `json:"emp_name"`
	SalesTeamID         int    `json:"sales_team_id"`
	SalesTeamCode       string `json:"sales_team_code"`
	SalesTeamName       string `json:"sales_team_name"`
	SalesName           string `json:"sales_name"`
	OperationType       string `json:"opr_type"`
	OperationTypeCanvas string `json:"opr_type_canvas"`
	OperationTypeName   string `json:"opr_type_name"`
	IsActiveCanvas      bool   `json:"is_active_canvas"`
	IsTakingOrder       bool   `json:"is_taking_order"`
	WarehouseID         int    `json:"wh_id"`
	WarehouseCode       string `json:"wh_code"`
	WarehouseName       string `json:"wh_name"`
}

type Warehouse struct {
	WarehouseID   int    `json:"wh_id"`
	WarehouseCode string `json:"wh_code"`
	WarehouseName string `json:"wh_name"`
}

type Outlet struct {
	DestinationID     int     `json:"outlet_id"`
	DestinationCode   string  `json:"outlet_code"`
	DestinationName   string  `json:"outlet_name"`
	DestinationStatus int     `json:"outlet_status"`
	Address1          string  `json:"address1"`
	Latitude          string  `json:"latitude"`
	Longitude         string  `json:"longitude"`
	AvgSalesWeek      float64 `json:"avg_sales_week"`
}

type Distributor struct {
	DistributorID     int     `json:"distributor_id"`
	DistributorCode   string  `json:"distributor_code"`
	DistributorName   string  `json:"distributor_name"`
	DistributorStatus int     `json:"distributor_status"`
	Address           string  `json:"address"`
	Latitude          string  `json:"latitude"`
	Longitude         string  `json:"longitude"`
	AvgSalesWeek      float64 `json:"avg_sales_week"`
}

type OutletList struct {
	DestinationID     int     `json:"outlet_id"`
	DestinationCode   string  `json:"outlet_code"`
	DestinationName   string  `json:"outlet_name"`
	DestinationStatus int     `json:"outlet_status"`
	Address1          string  `json:"address1"`
	Latitude          string  `json:"latitude"`
	Longitude         string  `json:"longitude"`
	OtGrpID           int     `json:"ot_grp_id"`
	OtGrpCode         string  `json:"ot_grp_code"`
	OtGrpName         string  `json:"ot_grp_name"`
	OtTypeID          int     `json:"ot_type_id"`
	OtTypeName        string  `json:"ot_type_name"`
	AvgSalesWeek      float64 `json:"avg_sales_week"`
}

type DmsQueryFilter struct {
	Page            string `query:"page"`
	Limit           string `query:"limit"`
	Query           string `query:"q"`
	DestinationCode string `query:"outlet_code"`
	DestinationID   int    `query:"outlet_id"`
	Mode            string `query:"mode"`
	Sort            string `query:"sort"`
	IsActive        string `query:"is_active"`
	SalesTeamID     string `query:"sales_team_id"`
}

type OutletQueryFilter struct {
	Page            string `query:"page"`
	Limit           string `query:"limit"`
	Query           string `query:"q"`
	OutletCode      string `query:"outlet_code"`
	OutletID        int    `query:"outlet_id"`
	OutletTypeID    int    `query:"outlet_type_id"`
	OutletTypeName  string `query:"outlet_type_name"`
	OutletGroupName string `query:"outlet_group_name"`
	Sort            string `query:"sort"`
	IsActive        string `query:"is_active"`
	SalesTeamID     string `query:"sales_team_id"`
}

type DistributorQueryFilter struct {
	Page                 string `query:"page"`
	Limit                string `query:"limit"`
	Query                string `query:"q"`
	DistributorCode      string `query:"distributor_code"`
	DistributorID        int    `query:"distributor_id"`
	DistributorTypeID    int    `query:"distributor_type_id"`
	DistributorTypeName  string `query:"distributor_type_name"`
	DistributorGroupName string `query:"distributor_group_name"`
	Sort                 string `query:"sort"`
	IsActive             string `query:"is_active"`
	SalesTeamID          string `query:"sales_team_id"`
}

type Meta struct {
	Limit     int `json:"limit"`
	Page      int `json:"page"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type OutletBySalesman struct {
	SalesmanCode []string `json:"salesman_code"`
}

type OutletBySalesmanId struct {
	SalesmanId string `json:"salesman_id"`
	Search     string `json:"search"`
}

type OutletNew struct {
	DestinationID       int         `json:"outlet_id"`
	DestinationCode     string      `json:"outlet_code"`
	DestinationName     string      `json:"outlet_name"`
	Barcode             string      `json:"barcode"`
	DestinationStatus   int         `json:"outlet_status"`
	Address1            string      `json:"address1"`
	Address2            string      `json:"address2"`
	City                string      `json:"city"`
	ZipCode             string      `json:"zip_code"`
	PhoneNo             string      `json:"phone_no"`
	WaNo                string      `json:"wa_no"`
	FaxNo               string      `json:"fax_no"`
	Email               string      `json:"email"`
	DiscGrpID           int         `json:"disc_grp_id"`
	DiscGrpCode         string      `json:"disc_grp_code"`
	DiscGrpName         string      `json:"disc_grp_name"`
	OtLocID             int         `json:"ot_loc_id"`
	OtLocCode           string      `json:"ot_loc_code"`
	OtLocName           string      `json:"ot_loc_name"`
	OtGrpID             int         `json:"ot_grp_id"`
	OtGrpCode           string      `json:"ot_grp_code"`
	OtGrpName           string      `json:"ot_grp_name"`
	PriceGrpID          int         `json:"price_grp_id"`
	PriceGrpCode        string      `json:"price_grp_code"`
	PriceGrpName        string      `json:"price_grp_name"`
	DistrictID          int         `json:"district_id"`
	DistrictCode        string      `json:"district_code"`
	DistrictName        string      `json:"district_name"`
	BeatID              int         `json:"beat_id"`
	BeatCode            string      `json:"beat_code"`
	BeatName            string      `json:"beat_name"`
	SbeatID             int         `json:"sbeat_id"`
	SbeatCode           string      `json:"sbeat_code"`
	SbeatName           string      `json:"sbeat_name"`
	OtClassID           int         `json:"ot_class_id"`
	OtClassCode         string      `json:"ot_class_code"`
	OtClassName         string      `json:"ot_class_name"`
	IndustryID          int         `json:"industry_id"`
	IndustryCode        string      `json:"industry_code"`
	IndustryName        string      `json:"industry_name"`
	MarketID            int         `json:"market_id"`
	MarketCode          string      `json:"market_code"`
	MarketName          string      `json:"market_name"`
	Top                 int         `json:"top"`
	DueDate             string      `json:"due_date"`
	PaymentType         int         `json:"payment_type"`
	IsContraBon         bool        `json:"is_contra_bon"`
	PluGrpID            int         `json:"plu_grp_id"`
	PluGrpCode          string      `json:"plu_grp_code"`
	PluGrpName          string      `json:"plu_grp_name"`
	ConvGrpID           int         `json:"conv_grp_id"`
	ConvGrpCode         string      `json:"conv_grp_code"`
	ConvGrpName         string      `json:"conv_grp_name"`
	DiscInvID           int         `json:"disc_inv_id"`
	DiscInvCode         string      `json:"disc_inv_code"`
	DiscInvName         string      `json:"disc_inv_name"`
	AgentFrom           string      `json:"agent_from"`
	CreditLimitType     int         `json:"credit_limit_type"`
	CreditLimit         int         `json:"credit_limit"`
	SalesInvLimitType   int         `json:"sales_inv_limit_type"`
	SalesInvLimit       int         `json:"sales_inv_limit"`
	AvgSalesWeek        int         `json:"avg_sales_week"`
	AvgSalesMonth       int         `json:"avg_sales_month"`
	FirstTransDate      string      `json:"first_trans_date"`
	LastTransDate       string      `json:"last_trans_date"`
	FirstWeekNo         int         `json:"first_week_no"`
	OtStartDate         string      `json:"ot_start_date"`
	OtRegDate           string      `json:"ot_reg_date"`
	BuildingOwn         int         `json:"building_own"`
	Dob                 string      `json:"dob"`
	ArStatus            int         `json:"ar_status"`
	ArTotal             int         `json:"ar_total"`
	ClosedDate          string      `json:"closed_date"`
	IsEmbBail           bool        `json:"is_emb_bail"`
	TaxName             string      `json:"tax_name"`
	TaxAddr1            string      `json:"tax_addr1"`
	TaxAddr2            string      `json:"tax_addr2"`
	TaxCity             string      `json:"tax_city"`
	TaxNo               string      `json:"tax_no"`
	TaxInvoiceForm      int         `json:"tax_invoice_form"`
	TaxInvoiceFormName  string      `json:"tax_invoice_form_name"`
	OwnerName           string      `json:"owner_name"`
	OwnerAddr1          string      `json:"owner_addr1"`
	OwnerAddr2          string      `json:"owner_addr2"`
	OwnerCity           string      `json:"owner_city"`
	OwnerPhoneNo        string      `json:"owner_phone_no"`
	OwnerIDNo           string      `json:"owner_id_no"`
	DelvAddr1           string      `json:"delv_addr1"`
	DelvAddr2           string      `json:"delv_addr2"`
	DelvCity            string      `json:"delv_city"`
	InvAddr1            string      `json:"inv_addr1"`
	InvAddr2            string      `json:"inv_addr2"`
	InvCity             string      `json:"inv_city"`
	IsActive            bool        `json:"is_active"`
	UpdatedBy           int         `json:"updated_by"`
	UpdatedAt           time.Time   `json:"updated_at"`
	UpdatedByName       string      `json:"updated_by_name"`
	Latitude            string      `json:"latitude"`
	Longitude           string      `json:"longitude"`
	OtTypeID            int         `json:"ot_type_id"`
	OtTypeCode          interface{} `json:"ot_type_code"`
	OtTypeName          interface{} `json:"ot_type_name"`
	IsObs               bool        `json:"is_obs"`
	Obs                 int         `json:"obs"`
	OutletWardID        string      `json:"outlet_ward_id"`
	OutletWard          string      `json:"outlet_ward"`
	OutletSubDistrictID string      `json:"outlet_sub_district_id"`
	OutletSubDistrict   string      `json:"outlet_sub_district"`
	OutletRegencyID     string      `json:"outlet_regency_id"`
	OutletRegency       string      `json:"outlet_regency"`
	OutletProvinceID    string      `json:"outlet_province_id"`
	OutletProvince      string      `json:"outlet_province"`
	IsWaNo              interface{} `json:"is_wa_no"`
	DelvWardID          interface{} `json:"delv_ward_id"`
	DelvWard            interface{} `json:"delv_ward"`
	DelvSubDistrictID   interface{} `json:"delv_sub_district_id"`
	DelvSubDistrict     interface{} `json:"delv_sub_district"`
	DelvRegencyID       interface{} `json:"delv_regency_id"`
	DelvRegency         interface{} `json:"delv_regency"`
	DelvProvinceID      interface{} `json:"delv_province_id"`
	DelvProvince        interface{} `json:"delv_province"`
	DelvZipCode         interface{} `json:"delv_zip_code"`
	DelvIsSameAddr      interface{} `json:"delv_is_same_addr"`
	InvWardID           interface{} `json:"inv_ward_id"`
	InvWard             interface{} `json:"inv_ward"`
	InvSubDistrictID    interface{} `json:"inv_sub_district_id"`
	InvSubDistrict      interface{} `json:"inv_sub_district"`
	InvRegencyID        interface{} `json:"inv_regency_id"`
	InvRegency          interface{} `json:"inv_regency"`
	InvProvinceID       interface{} `json:"inv_province_id"`
	InvProvince         interface{} `json:"inv_province"`
	InvZipCode          interface{} `json:"inv_zip_code"`
	InvIsSameAddr       interface{} `json:"inv_is_same_addr"`
	VerificationStatus  interface{} `json:"verification_status"`
	ImageURL            string      `json:"image_url,omitempty"`
}
