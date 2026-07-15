package entity

import (
	"time"
)

type DistPriceQueryFilter struct {
	SupId            *int   `query:"sup_id"`
	Page             int    `query:"page"`
	Limit            int    `query:"limit" validate:"required"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	IsActive         *int   `query:"is_active"`
	DistPriceGroupId *int   `query:"dist_price_group_id"`
	WhDate           string `query:"wh_date"`
	ProId            *int   `query:"pro_id"`
	WhId             *int   `query:"wh_id"`
}

type DistPriceResponse struct {
	DistPricePriceId      int64      `json:"dist_price_id"`
	DistPricePriceGroupId int64      `json:"dist_price_group_id"`
	DistPricePriceGrpCode string     `json:"dist_price_grp_code"`
	DistPricePriceGrpName string     `json:"dist_price_grp_name"`
	StartDate             string     `json:"start_date"`
	EndDate               string     `json:"end_date"`
	ProId                 int64      `json:"pro_id"`
	UnitId1               string     `json:"unit_id1"`
	UnitId2               string     `json:"unit_id2"`
	UnitId3               string     `json:"unit_id3"`
	UnitId4               *string    `json:"unit_id4"`
	UnitId5               *string    `json:"unit_id5"`
	ConvUnit2             float32    `json:"conv_unit2"`
	ConvUnit3             float32    `json:"conv_unit3"`
	ConvUnit4             float32    `json:"conv_unit4"`
	ConvUnit5             float32    `json:"conv_unit5"`
	PurchPrice1           float64    `json:"purch_price1"`
	PurchPrice2           float64    `json:"purch_price2"`
	PurchPrice3           float64    `json:"purch_price3"`
	PurchPrice4           *float64   `json:"purch_price4"`
	PurchPrice5           *float64   `json:"purch_price5"`
	SellPrice1            float64    `json:"sell_price1"`
	SellPrice2            float64    `json:"sell_price2"`
	SellPrice3            float64    `json:"sell_price3"`
	SellPrice4            *float64   `json:"sell_price4"`
	SellPrice5            *float64   `json:"sell_price5"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedByName         *string    `json:"updated_by_name"`
}

type DistPriceDetailResp struct {
	DistPricePriceId      int64      `json:"dist_price_id"`
	DistPricePriceGroupId int64      `json:"dist_price_group_id"`
	DistPricePriceGrpCode string     `json:"dist_price_grp_code"`
	DistPricePriceGrpName string     `json:"dist_price_grp_name"`
	StartDate             string     `json:"start_date"`
	EndDate               string     `json:"end_date"`
	ProId                 int64      `json:"pro_id"`
	ProCode               string     `json:"pro_code"`
	ProName               string     `json:"pro_name"`
	UnitId1               string     `json:"unit_id1"`
	UnitId2               string     `json:"unit_id2"`
	UnitId3               string     `json:"unit_id3"`
	UnitId4               *string    `json:"unit_id4"`
	UnitId5               *string    `json:"unit_id5"`
	UnitName1             string     `json:"unit_name1"`
	UnitName2             string     `json:"unit_name2"`
	UnitName3             string     `json:"unit_name3"`
	UnitName4             string     `json:"unit_name4"`
	UnitName5             string     `json:"unit_name5"`
	ConvUnit2             float32    `json:"conv_unit2"`
	ConvUnit3             float32    `json:"conv_unit3"`
	ConvUnit4             float32    `json:"conv_unit4"`
	ConvUnit5             float32    `json:"conv_unit5"`
	PurchPrice1           float64    `json:"purch_price1"`
	PurchPrice2           float64    `json:"purch_price2"`
	PurchPrice3           float64    `json:"purch_price3"`
	PurchPrice4           *float64   `json:"purch_price4"`
	PurchPrice5           *float64   `json:"purch_price5"`
	SellPrice1            float64    `json:"sell_price1"`
	SellPrice2            float64    `json:"sell_price2"`
	SellPrice3            float64    `json:"sell_price3"`
	SellPrice4            *float64   `json:"sell_price4"`
	SellPrice5            *float64   `json:"sell_price5"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedByName         *string    `json:"updated_by_name"`
	DistPricePriceIdOld   *int64     `json:"dist_price_id_old"`
	DistPriceGroupIdOld   int64      `json:"dist_price_group_id_old"`
	DistPriceGrpCodeOld   string     `json:"dist_price_grp_code_old"`
	DistPriceGrpNameOld   string     `json:"dist_price_grp_name_old"`
	StartDateOld          *time.Time `json:"start_date_old"`
	EndDateOld            *time.Time `json:"end_date_old"`
	ProIdOld              int64      `json:"pro_id_old"`
	UnitId1Old            string     `json:"unit_id1_old"`
	UnitId2Old            string     `json:"unit_id2_old"`
	UnitId3Old            string     `json:"unit_id3_old"`
	UnitId4Old            *string    `json:"unit_id4_old"`
	UnitId5Old            *string    `json:"unit_id5_old"`
	UnitName1Old          *string    `json:"unit_name1_old"`
	UnitName2Old          *string    `json:"unit_name2_old"`
	UnitName3Old          *string    `json:"unit_name3_old"`
	UnitName4Old          *string    `json:"unit_name4_old"`
	UnitName5Old          *string    `json:"unit_name5_old"`
	ConvUnit2Old          float32    `json:"conv_unit2_old"`
	ConvUnit3Old          float32    `json:"conv_unit3_old"`
	ConvUnit4Old          float32    `json:"conv_unit4_old"`
	ConvUnit5Old          float32    `json:"conv_unit5_old"`
	PurchPrice1Old        float64    `json:"purch_price1_old"`
	PurchPrice2Old        float64    `json:"purch_price2_old"`
	PurchPrice3Old        float64    `json:"purch_price3_old"`
	PurchPrice4Old        *float64   `json:"purch_price4_old"`
	PurchPrice5Old        *float64   `json:"purch_price5_old"`
	SellPrice1Old         float64    `json:"sell_price1_old"`
	SellPrice2Old         float64    `json:"sell_price2_old"`
	SellPrice3Old         float64    `json:"sell_price3_old"`
	SellPrice4Old         *float64   `json:"sell_price4_old"`
	SellPrice5Old         *float64   `json:"sell_price5_old"`
	Status                int        `json:"status"`
	StatusDesc            string     `json:"status_desc"`
}

type DistPriceListResponse struct {
	DistPricePriceId      int64      `json:"dist_price_id"`
	DistPricePriceGroupId int64      `json:"dist_price_group_id"`
	DistPricePriceGrpCode string     `json:"dist_price_grp_code"`
	DistPricePriceGrpName string     `json:"dist_price_grp_name"`
	StartDate             string     `json:"start_date"`
	EndDate               string     `json:"end_date"`
	ProId                 int64      `json:"pro_id"`
	ProCode               string     `json:"pro_code"`
	ProName               string     `json:"pro_name"`
	UnitId1               string     `json:"unit_id1"`
	UnitId2               string     `json:"unit_id2"`
	UnitId3               string     `json:"unit_id3"`
	UnitId4               *string    `json:"unit_id4"`
	UnitId5               *string    `json:"unit_id5"`
	ConvUnit2             float32    `json:"conv_unit2"`
	ConvUnit3             float32    `json:"conv_unit3"`
	ConvUnit4             float32    `json:"conv_unit4"`
	ConvUnit5             float32    `json:"conv_unit5"`
	PurchPrice1           float64    `json:"purch_price1"`
	PurchPrice2           float64    `json:"purch_price2"`
	PurchPrice3           float64    `json:"purch_price3"`
	PurchPrice4           *float64   `json:"purch_price4"`
	PurchPrice5           *float64   `json:"purch_price5"`
	SellPrice1            float64    `json:"sell_price1"`
	SellPrice2            float64    `json:"sell_price2"`
	SellPrice3            float64    `json:"sell_price3"`
	SellPrice4            *float64   `json:"sell_price4"`
	SellPrice5            *float64   `json:"sell_price5"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedByName         *string    `json:"updated_by_name"`
	DistPricePriceIdOld   *int64     `json:"dist_price_id_old"`
	Status                int        `json:"status"`
	StatusDesc            string     `json:"status_desc"`
}

type DistPriceLookupResponse struct {
	DistPricePriceId      int64    `json:"dist_price_id"`
	DistPricePriceGroupId int64    `json:"dist_price_group_id"`
	DistPricePriceGrpCode string   `json:"dist_price_grp_code"`
	DistPricePriceGrpName string   `json:"dist_price_grp_name"`
	StartDate             string   `json:"start_date"`
	EndDate               string   `json:"end_date"`
	ProId                 int64    `json:"pro_id"`
	ProCode               string   `json:"pro_code"`
	ProName               string   `json:"pro_name"`
	UnitId1               string   `json:"unit_id1"`
	UnitId2               string   `json:"unit_id2"`
	UnitId3               string   `json:"unit_id3"`
	UnitId4               *string  `json:"unit_id4"`
	UnitId5               *string  `json:"unit_id5"`
	UnitName1             string   `json:"unit_name1"`
	UnitName2             string   `json:"unit_name2"`
	UnitName3             string   `json:"unit_name3"`
	UnitName4             string   `json:"unit_name4"`
	UnitName5             string   `json:"unit_name5"`
	ConvUnit2             float32  `json:"conv_unit2"`
	ConvUnit3             float32  `json:"conv_unit3"`
	ConvUnit4             float32  `json:"conv_unit4"`
	ConvUnit5             float32  `json:"conv_unit5"`
	PurchPrice1           float64  `json:"purch_price1"`
	PurchPrice2           float64  `json:"purch_price2"`
	PurchPrice3           float64  `json:"purch_price3"`
	PurchPrice4           *float64 `json:"purch_price4"`
	PurchPrice5           *float64 `json:"purch_price5"`
	SellPrice1            float64  `json:"sell_price1"`
	SellPrice2            float64  `json:"sell_price2"`
	SellPrice3            float64  `json:"sell_price3"`
	SellPrice4            *float64 `json:"sell_price4"`
	SellPrice5            *float64 `json:"sell_price5"`
	SupId                 *int     `json:"sup_id"`
	SupCode               *string  `json:"sup_code"`
	SupName               *string  `json:"sup_name"`
	DistPricePriceIdOld   *int64   `json:"dist_price_id_old"`
	Vat                   *float64 `json:"vat"`
	VatBg                 *float64 `json:"vat_bg"`
	VatLgPurch            *float64 `json:"vat_lg_purch"`
	VatLgSell             *float64 `json:"vat_lg_sell"`
	ExciseRate            float64  `json:"excise_rate"`
	ExciseTax             float64  `json:"excise_tax"`
	Qty                   float64  `json:"qty"`
	QtyOnShipping         float64  `json:"qty_on_shipping"`
	QtyOnOrder            float64  `json:"qty_on_order"`
	QtyBs                 float64  `json:"qty_bs"`
	QtyExp                float64  `json:"qty_exp"`
}

type DistPriceLookupProResp struct {
	DistPricePriceId      int64    `json:"dist_price_id"`
	DistPricePriceGroupId int64    `json:"dist_price_group_id"`
	DistPricePriceGrpCode string   `json:"dist_price_grp_code"`
	DistPricePriceGrpName string   `json:"dist_price_grp_name"`
	ProId                 int64    `json:"pro_id"`
	ProCode               string   `json:"pro_code"`
	ProName               string   `json:"pro_name"`
	UnitId1               string   `json:"unit_id1"`
	UnitId2               string   `json:"unit_id2"`
	UnitId3               string   `json:"unit_id3"`
	UnitId4               *string  `json:"unit_id4"`
	UnitId5               *string  `json:"unit_id5"`
	UnitName1             string   `json:"unit_name1"`
	UnitName2             string   `json:"unit_name2"`
	UnitName3             string   `json:"unit_name3"`
	UnitName4             string   `json:"unit_name4"`
	UnitName5             string   `json:"unit_name5"`
	ConvUnit2             float32  `json:"conv_unit2"`
	ConvUnit3             float32  `json:"conv_unit3"`
	ConvUnit4             float32  `json:"conv_unit4"`
	ConvUnit5             float32  `json:"conv_unit5"`
	PurchPrice1           float64  `json:"purch_price1"`
	PurchPrice2           float64  `json:"purch_price2"`
	PurchPrice3           float64  `json:"purch_price3"`
	PurchPrice4           *float64 `json:"purch_price4"`
	PurchPrice5           *float64 `json:"purch_price5"`
	SellPrice1            float64  `json:"sell_price1"`
	SellPrice2            float64  `json:"sell_price2"`
	SellPrice3            float64  `json:"sell_price3"`
	SellPrice4            *float64 `json:"sell_price4"`
	SellPrice5            *float64 `json:"sell_price5"`
	DistPricePriceIdOld   *int64   `json:"dist_price_id_old"`
}

type CreateDistPriceBody struct {
	CustId                string   `json:"cust_id" validate:"required,max=10"`
	CreatedBy             int64    `json:"created_by" validate:"required"`
	UpdatedBy             int64    `json:"updated_by"`
	DistPricePriceId      int64    `json:"dist_price_id"`
	DistPricePriceGroupId int64    `json:"dist_price_group_id"`
	StartDate             string   `json:"start_date"`
	EndDate               string   `json:"end_date"`
	ProId                 int64    `json:"pro_id"`
	UnitId1               string   `json:"unit_id1"`
	UnitId2               string   `json:"unit_id2"`
	UnitId3               string   `json:"unit_id3"`
	UnitId4               *string  `json:"unit_id4"`
	UnitId5               *string  `json:"unit_id5"`
	PurchPrice1           *float64 `json:"purch_price1"`
	PurchPrice2           *float64 `json:"purch_price2"`
	PurchPrice3           *float64 `json:"purch_price3"`
	PurchPrice4           *float64 `json:"purch_price4"`
	PurchPrice5           *float64 `json:"purch_price5"`
	ConvUnit2             float32  `json:"conv_unit2"`
	ConvUnit3             float32  `json:"conv_unit3"`
	ConvUnit4             float32  `json:"conv_unit4"`
	ConvUnit5             float32  `json:"conv_unit5"`
	SellPrice1            float64  `json:"sell_price1"`
	SellPrice2            float64  `json:"sell_price2"`
	SellPrice3            float64  `json:"sell_price3"`
	SellPrice4            *float64 `json:"sell_price4"`
	SellPrice5            *float64 `json:"sell_price5"`
	DistPricePriceIdOld   *int64   `json:"dist_price_id_old"`
}

type DetailDistPriceParams struct {
	DistPriceId int64 `params:"dist_price_id" validate:"required"`
}

type UpdateDistPriceParams struct {
	DistPriceId int64 `params:"dist_price_id" validate:"required"`
}

type DeleteDistPriceParams struct {
	DistPriceId int64 `params:"dist_price_id" validate:"required"`
}

type UpdateDistPriceRequest struct {
	CustId                string   `json:"cust_id"`
	UpdatedBy             int64    `json:"updated_by"`
	DistPricePriceGroupId *int64   `json:"dist_price_group_id"`
	StartDate             *string  `json:"start_date"`
	EndDate               *string  `json:"end_date"`
	ProId                 *int64   `json:"pro_id"`
	UnitId1               *string  `json:"unit_id1"`
	UnitId2               *string  `json:"unit_id2"`
	UnitId3               *string  `json:"unit_id3"`
	UnitId4               *string  `json:"unit_id4"`
	UnitId5               *string  `json:"unit_id5"`
	ConvUnit2             *float32 `json:"conv_unit2"`
	ConvUnit3             *float32 `json:"conv_unit3"`
	ConvUnit4             *float32 `json:"conv_unit4"`
	ConvUnit5             *float32 `json:"conv_unit5"`
	PurchPrice1           *float64 `json:"purch_price1"`
	PurchPrice2           *float64 `json:"purch_price2"`
	PurchPrice3           *float64 `json:"purch_price3"`
	PurchPrice4           *float64 `json:"purch_price4"`
	PurchPrice5           *float64 `json:"purch_price5"`
	SellPrice1            *float64 `json:"sell_price1"`
	SellPrice2            *float64 `json:"sell_price2"`
	SellPrice3            *float64 `json:"sell_price3"`
	SellPrice4            *float64 `json:"sell_price4"`
	SellPrice5            *float64 `json:"sell_price5"`
	DistPricePriceIdOld   *int64   `json:"dist_price_id_old"`
}

var DistPriceStatusDesc = map[int]string{
	1: "Scheduled", 5: "Cancelled", 7: "Inactive", 10: "Published",
}

type DistPriceStatus struct {
	StatusID   int    `json:"status_id"`
	StatusDesc string `json:"status_desc"`
}

func (price DistPriceListResponse) GetDistPriceStatusDesc() string {
	return DistPriceStatusDesc[price.Status]
}

type PublishUnpublishDistPriceReq struct {
	DistPriceID  int64  `json:"dist_price_id"`
	CustID       string `json:"cust_id"`
	Status       int    `json:"status"`
	UpdatedBy    string `json:"updated_by"`
}
