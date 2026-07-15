package entity

import "time"

type MSpPriceQueryFilter struct {
	CustID       string
	ParentCustID string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Sort         string `query:"sort"`
	Status       int    `query:"status"`
}

type PreviewMSpPriceResp struct {
	StartDate        string           `json:"start_date"`
	EndDate          string           `json:"end_date"`
	PriceGrpID       int64            `json:"price_grp_id"`
	PriceGrpCode     string           `json:"price_grp_code"`
	PriceGrpName     string           `json:"price_grp_name"`
	ProID            int64            `json:"pro_id"`
	ProCode          string           `json:"pro_code"`
	ProName          string           `json:"pro_name"`
	UnitId1          string           `json:"unit_id1"`
	UnitId2          string           `json:"unit_id2"`
	UnitId3          string           `json:"unit_id3"`
	PurchPrice1      float64          `json:"purch_price1"`
	PurchPrice2      float64          `json:"purch_price2"`
	PurchPrice3      float64          `json:"purch_price3"`
	MasterSellPrice1 float64          `json:"master_sell_price1"`
	MasterSellPrice2 float64          `json:"master_sell_price2"`
	MasterSellPrice3 float64          `json:"master_sell_price3"`
	NewSellPrice1    float64          `json:"new_sell_price1"`
	NewSellPrice2    float64          `json:"new_sell_price2"`
	NewSellPrice3    float64          `json:"new_sell_price3"`
	ConvUnit2        float32          `json:"conv_unit2"`
	ConvUnit3        float32          `json:"conv_unit3"`
	Status           int              `json:"status"`
	StatusDesc       string           `json:"status_desc"`
	Details          MSpPriceDetGroup `json:"details"`
	OutputSprice     []OutputSprice   `json:"output"`
}

type OutputSprice struct {
	OutletID         int64   `json:"outlet_id"`
	OutletCode       string  `json:"outlet_code"`
	OutletName       string  `json:"outlet_name"`
	MasterSellPrice1 float64 `json:"master_sell_price1"`
	MasterSellPrice2 float64 `json:"master_sell_price2"`
	MasterSellPrice3 float64 `json:"master_sell_price3"`
	NewSellPrice1    float64 `json:"new_sell_price1"`
	NewSellPrice2    float64 `json:"new_sell_price2"`
	NewSellPrice3    float64 `json:"new_sell_price3"`
}

type CreateMSpPriceBody struct {
	ParentCustID    string           `json:"parent_cust_id"`
	CustID          string           `json:"cust_id"`
	SpPriceID       string           `json:"sp_price_id"`
	StartDate       string           `json:"start_date" validate:"required,max=10"`
	EndDate         string           `json:"end_date" validate:"required,max=10"`
	PriceGrpID      int64            `json:"price_grp_id" validate:"required"`
	PriceGrpCode    string           `json:"price_grp_code"`
	PriceGrpName    string           `json:"price_grp_name"`
	ProID           int64            `json:"pro_id" validate:"required"`
	ProCode         string           `json:"pro_code"`
	ProName         string           `json:"pro_name"`
	UnitId1         string           `json:"unit_id1"`
	UnitId2         string           `json:"unit_id2"`
	UnitId3         string           `json:"unit_id3"`
	SellPrice1      float64          `json:"sell_price1"`
	SellPrice2      float64          `json:"sell_price2"`
	SellPrice3      float64          `json:"sell_price3"`
	NewSellPrice1   float64          `json:"new_sell_price1" validate:"required"`
	NewSellPrice2   float64          `json:"new_sell_price2" validate:"required"`
	NewSellPrice3   float64          `json:"new_sell_price3" validate:"required"`
	ConvUnit2       float32          `json:"conv_unit2"`
	ConvUnit3       float32          `json:"conv_unit3"`
	Status          int              `json:"status"`
	CreatedBy       string           `json:"created_by"`
	Details         MSpPriceDetGroup `json:"details"`
	ExpirationMs    int              `json:"expiration_ms"`
	EndExpirationMs int              `json:"end_expiration_ms"`
}

type MSpPriceResponse struct {
	SpPriceID     string     `json:"sp_price_id"`
	StartDate     *string    `json:"start_date"`
	EndDate       *string    `json:"end_date"`
	PriceGrpID    *int64     `json:"price_grp_id"`
	PriceGrpCode  *string    `json:"price_grp_code"`
	PriceGrpName  *string    `json:"price_grp_name"`
	ProID         *int64     `json:"pro_id"`
	ProCode       *string    `json:"pro_code"`
	ProName       *string    `json:"pro_name"`
	UnitId1       *string    `json:"unit_id1"`
	UnitId2       *string    `json:"unit_id2"`
	UnitId3       *string    `json:"unit_id3"`
	UnitName1     *string    `json:"unit_name1"`
	UnitName2     *string    `json:"unit_name2"`
	UnitName3     *string    `json:"unit_name3"`
	SellPrice1    float64    `json:"sell_price1"`
	SellPrice2    float64    `json:"sell_price2"`
	SellPrice3    float64    `json:"sell_price3"`
	NewSellPrice1 float64    `json:"new_sell_price1"`
	NewSellPrice2 float64    `json:"new_sell_price2"`
	NewSellPrice3 float64    `json:"new_sell_price3"`
	ConvUnit2     float32    `json:"conv_unit2"`
	ConvUnit3     float32    `json:"conv_unit3"`
	Status        int        `json:"status"`
	StatusDesc    string     `json:"status_desc"`
	UpdatedBy     string     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type MSpPriceParams struct {
	CustID       string
	ParentCustID string
	SpPriceID    string `params:"sp_price_id" validate:"required"`
}

type MSpPriceCancelParams struct {
	CustID       string
	ParentCustID string
	SpPriceID    string `params:"sp_price_id" validate:"required"`
	UpdatedBy    string
}

type MSpPriceDeleteParams struct {
	SpPriceID string `params:"sp_price_id" validate:"required"`
}

type MSpPriceListResponse struct {
	SpPriceID     string     `json:"sp_price_id"`
	StartDate     *string    `json:"start_date"`
	EndDate       *string    `json:"end_date"`
	PriceGrpID    *int64     `json:"price_grp_id"`
	PriceGrpCode  *string    `json:"price_grp_code"`
	PriceGrpName  *string    `json:"price_grp_name"`
	ProID         *int64     `json:"pro_id"`
	ProCode       *string    `json:"pro_code"`
	ProName       *string    `json:"pro_name"`
	UnitId1       *string    `json:"unit_id1"`
	UnitId2       *string    `json:"unit_id2"`
	UnitId3       *string    `json:"unit_id3"`
	SellPrice1    float64    `json:"sell_price1"`
	SellPrice2    float64    `json:"sell_price2"`
	SellPrice3    float64    `json:"sell_price3"`
	NewSellPrice1 float64    `json:"new_sell_price1"`
	NewSellPrice2 float64    `json:"new_sell_price2"`
	NewSellPrice3 float64    `json:"new_sell_price3"`
	ConvUnit2     float32    `json:"conv_unit2"`
	ConvUnit3     float32    `json:"conv_unit3"`
	Status        string     `json:"status"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedBy     string     `json:"updated_by"`
}

type UpdateMSpPriceBody struct {
	ParentCustID  string           `json:"parent_cust_id"`
	CustID        string           `json:"cust_id"`
	SpPriceID     string           `json:"sp_price_id"`
	StartDate     string           `json:"start_date,omitempty" validate:"required,max=10"`
	EndDate       string           `json:"end_date,omitempty" validate:"required,max=10"`
	PriceGrpID    int64            `json:"price_grp_id,omitempty" validate:"required"`
	ProID         int64            `json:"pro_id,omitempty" validate:"required"`
	UnitId1       string           `json:"unit_id1,omitempty"`
	UnitId2       string           `json:"unit_id2,omitempty"`
	UnitId3       string           `json:"unit_id3,omitempty"`
	SellPrice1    float64          `json:"sell_price1,omitempty"`
	SellPrice2    float64          `json:"sell_price2,omitempty"`
	SellPrice3    float64          `json:"sell_price3,omitempty"`
	NewSellPrice1 float64          `json:"new_sell_price1,omitempty"`
	NewSellPrice2 float64          `json:"new_sell_price2,omitempty"`
	NewSellPrice3 float64          `json:"new_sell_price3,omitempty"`
	ConvUnit2     float32          `json:"conv_unit2,omitempty"`
	ConvUnit3     float32          `json:"conv_unit3,omitempty"`
	Status        int              `json:"status"`
	UpdatedBy     string           `json:"updated_by" validate:"required"`
	UpdatedAt     *time.Time       `json:"updated_at"`
	Details       MSpPriceDetGroup `json:"details"`
}

type MSpPriceStatusDescSlice []MSpPriceStatus

func (s MSpPriceStatusDescSlice) Len() int {
	return len(s)
}

func (s MSpPriceStatusDescSlice) Less(i, j int) bool {
	return s[i].StatusID < s[j].StatusID
}

func (p MSpPriceStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

var MSpPriceStatusDesc = map[int]string{
	1: "Scheduled", 5: "Cancelled", 7: "Inactive", 10: "Published",
}

type MSpPriceStatus struct {
	StatusID   int    `json:"status_id"`
	StatusDesc string `json:"status_desc"`
}

func (price MSpPriceResponse) GetSpPriceStatusDesc() string {
	return MSpPriceStatusDesc[price.Status]
}

type MSpPriceWithDetailResp struct {
	SpPriceID     string               `json:"sp_price_id"`
	StartDate     *string              `json:"start_date"`
	EndDate       *string              `json:"end_date"`
	PriceGrpID    *int64               `json:"price_grp_id"`
	PriceGrpCode  *string              `json:"price_grp_code"`
	PriceGrpName  *string              `json:"price_grp_name"`
	ProID         *int64               `json:"pro_id"`
	ProCode       *string              `json:"pro_code"`
	ProName       *string              `json:"pro_name"`
	UnitId1       *string              `json:"unit_id1"`
	UnitId2       *string              `json:"unit_id2"`
	UnitId3       *string              `json:"unit_id3"`
	UnitName1     *string              `json:"unit_name1"`
	UnitName2     *string              `json:"unit_name2"`
	UnitName3     *string              `json:"unit_name3"`
	PurchPrice1   float64              `json:"purch_price1"`
	PurchPrice2   float64              `json:"purch_price2"`
	PurchPrice3   float64              `json:"purch_price3"`
	SellPrice1    float64              `json:"sell_price1"`
	SellPrice2    float64              `json:"sell_price2"`
	SellPrice3    float64              `json:"sell_price3"`
	NewSellPrice1 float64              `json:"new_sell_price1"`
	NewSellPrice2 float64              `json:"new_sell_price2"`
	NewSellPrice3 float64              `json:"new_sell_price3"`
	ConvUnit2     float32              `json:"conv_unit2"`
	ConvUnit3     float32              `json:"conv_unit3"`
	Status        int                  `json:"status"`
	StatusDesc    string               `json:"status_desc"`
	UpdatedBy     string               `json:"updated_by"`
	UpdatedAt     *time.Time           `json:"updated_at"`
	Details       MSpPriceDetRespGroup `json:"details"`
	OutputSprice  []OutputSprice       `json:"output"`
}

type PublishUnpublishSPriceReq struct {
	SpPriceID    string `json:"sp_price_id"`
	CustID       string `json:"cust_id"`
	ParentCustID string `json:"parent_cust_id"`
	Status       int    `json:"status"`
	UpdatedBy    string `json:"updated_by"`
}
