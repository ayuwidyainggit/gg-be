package model

import (
	"strconv"
	"strings"
)

type OutletAPIResponse struct {
	Data   []OutletAPI     `json:"data"`
	Paging OutletAPIPaging `json:"paging"`
}

type OutletAPIPaging struct {
	TotalRecord int `json:"total_record"`
	PageCurrent int `json:"page_current"`
	PageLimit   int `json:"page_limit"`
	PageTotal   int `json:"page_total"`
}

type OutletAPI struct {
	OutletID            int         `json:"outlet_id"`
	OutletCode          string      `json:"outlet_code"`
	OutletName          string      `json:"outlet_name"`
	Barcode             string      `json:"barcode"`
	OutletStatus        string      `json:"outlet_status"`
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
	CreditLimit         float64     `json:"credit_limit"`
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
	UpdatedAt           string      `json:"updated_at"`
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

func parseOutletStatus(raw string) int {
	outletStatus, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0
	}

	return outletStatus
}

func (o OutletAPI) ToOutlet() Outlet {
	return Outlet{
		OutletID:     o.OutletID,
		OutletCode:   o.OutletCode,
		OutletName:   o.OutletName,
		OutletStatus: parseOutletStatus(o.OutletStatus),
		Address1:     o.Address1,
		Latitude:     o.Latitude,
		Longitude:    o.Longitude,
		AvgSalesWeek: float64(o.AvgSalesWeek),
	}
}

func (o OutletAPI) ToOutletNew() OutletNew {
	return OutletNew{
		OutletID:            o.OutletID,
		OutletCode:          o.OutletCode,
		OutletName:          o.OutletName,
		Barcode:             o.Barcode,
		OutletStatus:        parseOutletStatus(o.OutletStatus),
		Address1:            o.Address1,
		Address2:            o.Address2,
		City:                o.City,
		ZipCode:             o.ZipCode,
		PhoneNo:             o.PhoneNo,
		WaNo:                o.WaNo,
		FaxNo:               o.FaxNo,
		Email:               o.Email,
		DiscGrpID:           o.DiscGrpID,
		DiscGrpCode:         o.DiscGrpCode,
		DiscGrpName:         o.DiscGrpName,
		OtLocID:             o.OtLocID,
		OtLocCode:           o.OtLocCode,
		OtLocName:           o.OtLocName,
		OtGrpID:             o.OtGrpID,
		OtGrpCode:           o.OtGrpCode,
		OtGrpName:           o.OtGrpName,
		PriceGrpID:          o.PriceGrpID,
		PriceGrpCode:        o.PriceGrpCode,
		PriceGrpName:        o.PriceGrpName,
		DistrictID:          o.DistrictID,
		DistrictCode:        o.DistrictCode,
		DistrictName:        o.DistrictName,
		BeatID:              o.BeatID,
		BeatCode:            o.BeatCode,
		BeatName:            o.BeatName,
		SbeatID:             o.SbeatID,
		SbeatCode:           o.SbeatCode,
		SbeatName:           o.SbeatName,
		OtClassID:           o.OtClassID,
		OtClassCode:         o.OtClassCode,
		OtClassName:         o.OtClassName,
		IndustryID:          o.IndustryID,
		IndustryCode:        o.IndustryCode,
		IndustryName:        o.IndustryName,
		MarketID:            o.MarketID,
		MarketCode:          o.MarketCode,
		MarketName:          o.MarketName,
		Top:                 o.Top,
		DueDate:             o.DueDate,
		PaymentType:         o.PaymentType,
		IsContraBon:         o.IsContraBon,
		PluGrpID:            o.PluGrpID,
		PluGrpCode:          o.PluGrpCode,
		PluGrpName:          o.PluGrpName,
		ConvGrpID:           o.ConvGrpID,
		ConvGrpCode:         o.ConvGrpCode,
		ConvGrpName:         o.ConvGrpName,
		DiscInvID:           o.DiscInvID,
		DiscInvCode:         o.DiscInvCode,
		DiscInvName:         o.DiscInvName,
		AgentFrom:           o.AgentFrom,
		CreditLimitType:     o.CreditLimitType,
		CreditLimit:         o.CreditLimit,
		SalesInvLimitType:   o.SalesInvLimitType,
		SalesInvLimit:       o.SalesInvLimit,
		AvgSalesWeek:        o.AvgSalesWeek,
		AvgSalesMonth:       o.AvgSalesMonth,
		FirstTransDate:      o.FirstTransDate,
		LastTransDate:       o.LastTransDate,
		FirstWeekNo:         o.FirstWeekNo,
		OtStartDate:         o.OtStartDate,
		OtRegDate:           o.OtRegDate,
		BuildingOwn:         o.BuildingOwn,
		Dob:                 o.Dob,
		ArStatus:            o.ArStatus,
		ArTotal:             o.ArTotal,
		ClosedDate:          o.ClosedDate,
		IsEmbBail:           o.IsEmbBail,
		TaxName:             o.TaxName,
		TaxAddr1:            o.TaxAddr1,
		TaxAddr2:            o.TaxAddr2,
		TaxCity:             o.TaxCity,
		TaxNo:               o.TaxNo,
		TaxInvoiceForm:      o.TaxInvoiceForm,
		TaxInvoiceFormName:  o.TaxInvoiceFormName,
		OwnerName:           o.OwnerName,
		OwnerAddr1:          o.OwnerAddr1,
		OwnerAddr2:          o.OwnerAddr2,
		OwnerCity:           o.OwnerCity,
		OwnerPhoneNo:        o.OwnerPhoneNo,
		OwnerIDNo:           o.OwnerIDNo,
		DelvAddr1:           o.DelvAddr1,
		DelvAddr2:           o.DelvAddr2,
		DelvCity:            o.DelvCity,
		InvAddr1:            o.InvAddr1,
		InvAddr2:            o.InvAddr2,
		InvCity:             o.InvCity,
		IsActive:            o.IsActive,
		UpdatedBy:           o.UpdatedBy,
		Latitude:            o.Latitude,
		Longitude:           o.Longitude,
		OtTypeID:            o.OtTypeID,
		OtTypeCode:          o.OtTypeCode,
		OtTypeName:          o.OtTypeName,
		IsObs:               o.IsObs,
		Obs:                 o.Obs,
		OutletWardID:        o.OutletWardID,
		OutletWard:          o.OutletWard,
		OutletSubDistrictID: o.OutletSubDistrictID,
		OutletSubDistrict:   o.OutletSubDistrict,
		OutletRegencyID:     o.OutletRegencyID,
		OutletRegency:       o.OutletRegency,
		OutletProvinceID:    o.OutletProvinceID,
		OutletProvince:      o.OutletProvince,
		IsWaNo:              o.IsWaNo,
		DelvWardID:          o.DelvWardID,
		DelvWard:            o.DelvWard,
		DelvSubDistrictID:   o.DelvSubDistrictID,
		DelvSubDistrict:     o.DelvSubDistrict,
		DelvRegencyID:       o.DelvRegencyID,
		DelvRegency:         o.DelvRegency,
		DelvProvinceID:      o.DelvProvinceID,
		DelvProvince:        o.DelvProvince,
		DelvZipCode:         o.DelvZipCode,
		DelvIsSameAddr:      o.DelvIsSameAddr,
		InvWardID:           o.InvWardID,
		InvWard:             o.InvWard,
		InvSubDistrictID:    o.InvSubDistrictID,
		InvSubDistrict:      o.InvSubDistrict,
		InvRegencyID:        o.InvRegencyID,
		InvRegency:          o.InvRegency,
		InvProvinceID:       o.InvProvinceID,
		InvProvince:         o.InvProvince,
		InvZipCode:          o.InvZipCode,
		InvIsSameAddr:       o.InvIsSameAddr,
		VerificationStatus:  o.VerificationStatus,
		ImageURL:            o.ImageURL,
	}
}

func ConvertOutletAPIsToOutlets(outlets []OutletAPI) []Outlet {
	responses := make([]Outlet, 0, len(outlets))
	for _, outlet := range outlets {
		responses = append(responses, outlet.ToOutlet())
	}

	return responses
}

func ConvertOutletAPIsToOutletNews(outlets []OutletAPI) []OutletNew {
	responses := make([]OutletNew, 0, len(outlets))
	for _, outlet := range outlets {
		responses = append(responses, outlet.ToOutletNew())
	}

	return responses
}
