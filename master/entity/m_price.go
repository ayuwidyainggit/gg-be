package entity

import (
	"time"
)

type MPriceQueryFilter struct {
	Page               int     `query:"page"`
	Limit              int     `query:"limit" validate:"required"`
	Query              string  `query:"q"`
	Sort               string  `query:"sort"`
	Status             []int   `query:"-"`
	FileType           string  `query:"file_type"`
	EffectiveDateStart string  `query:"effective_date_start"`
	EffectiveDateEnd   string  `query:"effective_date_end"`
	DistributorIDs     []int64 `query:"-"`
}

type MPriceResponse struct {
	CustID         string                      `json:"cust_id,omitempty"`
	ParentCustID   string                      `json:"parent_cust_id,omitempty"`
	PriceID        string                      `json:"price_id"`
	Status         int                         `json:"status"`
	StatusDesc     string                      `json:"status_desc"`
	Coverage       string                      `json:"coverage"`
	EffectiveDate  string                      `json:"effective_date"`
	ProID          int64                       `json:"pro_id"`
	ProCode        string                      `json:"pro_code"`
	ProName        string                      `json:"pro_name"`
	UnitId1        string                      `json:"unit_id1"`
	UnitId2        string                      `json:"unit_id2"`
	UnitId3        string                      `json:"unit_id3"`
	ConvUnit2      int                         `json:"conv_unit2"`
	ConvUnit3      int                         `json:"conv_unit3"`
	PurchPrice1    float64                     `json:"purch_price1"`
	PurchPrice2    float64                     `json:"purch_price2"`
	PurchPrice3    float64                     `json:"purch_price3"`
	SellPrice1     float64                     `json:"sell_price1"`
	SellPrice2     float64                     `json:"sell_price2"`
	SellPrice3     float64                     `json:"sell_price3"`
	NewPurchPrice1 float64                     `json:"new_purch_price1"`
	NewPurchPrice2 float64                     `json:"new_purch_price2"`
	NewPurchPrice3 float64                     `json:"new_purch_price3"`
	NewSellPrice1  float64                     `json:"new_sell_price1"`
	NewSellPrice2  float64                     `json:"new_sell_price2"`
	NewSellPrice3  float64                     `json:"new_sell_price3"`
	CreatedByID    *int64                      `json:"created_by_id,omitempty"`
	CreatedBy      string                      `json:"created_by"`
	CreatedAt      *time.Time                  `json:"created_at,omitempty"`
	UpdatedByID    *int64                      `json:"updated_by_id,omitempty"`
	UpdatedBy      string                      `json:"updated_by"`
	UpdatedAt      *time.Time                  `json:"updated_at,omitempty"`
	DistributorIDs []int64                     `json:"distributor_ids,omitempty"`
	Details        []DistributorAreaRegionData `json:"details,omitempty"`
}

type MPriceLookupResponse struct {
	PriceId          int     `json:"price_id"`
	Coverage         string  `json:"coverage"`
	DataCoverage     int     `json:"data_coverage"`
	DataCoverageCode *string `json:"data_coverage_code"`
	DataCoverageName *string `json:"data_coverage_name"`
	EffectiveDate    string  `json:"effective_date"`
	EffectiveDateEnd string  `json:"effective_date_end"`
	ProId            int     `json:"pro_id"`
	ProductName      string  `json:"pro_name"`
	UnitId1          string  `json:"unit_id1"`
	UnitId2          string  `json:"unit_id2"`
	UnitId3          string  `json:"unit_id3"`
	UnitId4          *string `json:"unit_id4"`
	UnitId5          *string `json:"unit_id5"`
	UnitName1        string  `json:"unit_name1"`
	UnitName2        string  `json:"unit_name2"`
	UnitName3        string  `json:"unit_name3"`
	UnitName4        *string `json:"unit_name4"`
	UnitName5        *string `json:"unit_name5"`
	ConvUnit2        int     `json:"conv_unit2"`
	ConvUnit3        int     `json:"conv_unit3"`
	ConvUnit4        int     `json:"conv_unit4"`
	ConvUnit5        int     `json:"conv_unit5"`
	PurchPrice1      float64 `json:"purch_price1"`
	PurchPrice2      float64 `json:"purch_price2"`
	PurchPrice3      float64 `json:"purch_price3"`
	PurchPrice4      float64 `json:"purch_price4"`
	PurchPrice5      float64 `json:"purch_price5"`
	SellPrice1       float64 `json:"sell_price1"`
	SellPrice2       float64 `json:"sell_price2"`
	SellPrice3       float64 `json:"sell_price3"`
	SellPrice4       float64 `json:"sell_price4"`
	SellPrice5       float64 `json:"sell_price5"`
	NewPurchPrice1   float64 `json:"new_purch_price1"`
	NewPurchPrice2   float64 `json:"new_purch_price2"`
	NewPurchPrice3   float64 `json:"new_purch_price3"`
	NewPurchPrice4   float64 `json:"new_purch_price4"`
	NewPurchPrice5   float64 `json:"new_purch_price5"`
	NewSellPrice1    float64 `json:"new_sell_price1"`
	NewSellPrice2    float64 `json:"new_sell_price2"`
	NewSellPrice3    float64 `json:"new_sell_price3"`
	NewSellPrice4    float64 `json:"new_sell_price4"`
	NewSellPrice5    float64 `json:"new_sell_price5"`
}

type CreateMPriceBody struct {
	ParentCustID   string  `json:"-"`
	CustID         string  `json:"-"`
	CreatedBy      string  `json:"-"`
	CreatedByID    *int64  `json:"-"`
	DistributorID  int64   `json:"-"`
	Coverage       string  `json:"coverage" validate:"required,oneof='N' 'D'"`
	DistributorIDs []int64 `json:"distributor_ids" validate:"omitempty"`
	EffectiveDate  string  `json:"effective_date" validate:"required"`
	ProID          int64   `json:"pro_id" validate:"required"`
	UnitID1        string  `json:"unit_id1"`
	UnitID2        string  `json:"unit_id2"`
	UnitID3        string  `json:"unit_id3"`
	ConvUnit2      int     `json:"conv_unit2"`
	ConvUnit3      int     `json:"conv_unit3"`
	PurchPrice1    float64 `json:"purch_price1"`
	PurchPrice2    float64 `json:"purch_price2"`
	PurchPrice3    float64 `json:"purch_price3"`
	SellPrice1     float64 `json:"sell_price1"`
	SellPrice2     float64 `json:"sell_price2"`
	SellPrice3     float64 `json:"sell_price3"`
	NewPurchPrice1 float64 `json:"new_purch_price1" validate:"required"`
	NewPurchPrice2 float64 `json:"new_purch_price2" validate:"required"`
	NewPurchPrice3 float64 `json:"new_purch_price3" validate:"required"`
	NewSellPrice1  float64 `json:"new_sell_price1" validate:"required"`
	NewSellPrice2  float64 `json:"new_sell_price2" validate:"required"`
	NewSellPrice3  float64 `json:"new_sell_price3" validate:"required"`
	Status         int     `json:"status"`
	ExpirationMs   int     `json:"expiration_ms"`
}

type UpdateMPriceRequest struct {
	CustID         string   `json:"-"`
	ParentCustID   string   `json:"-"`
	UpdatedBy      string   `json:"-"`
	UpdatedByID    *int64   `json:"-"`
	DistributorID  int64    `json:"-"`
	EffectiveDate  *string  `json:"effective_date" validate:"required"`
	Coverage       *string  `json:"coverage" validate:"required,oneof='N' 'D'"`
	DistributorIDs *[]int64 `json:"distributor_ids"`
	ProID          *int64   `json:"pro_id" validate:"required"`
	UnitID1        *string  `json:"unit_id1"`
	UnitID2        *string  `json:"unit_id2"`
	UnitID3        *string  `json:"unit_id3"`
	ConvUnit2      *int     `json:"conv_unit2"`
	ConvUnit3      *int     `json:"conv_unit3"`
	PurchPrice1    *float64 `json:"purch_price1"`
	PurchPrice2    *float64 `json:"purch_price2"`
	PurchPrice3    *float64 `json:"purch_price3"`
	SellPrice1     *float64 `json:"sell_price1"`
	SellPrice2     *float64 `json:"sell_price2"`
	SellPrice3     *float64 `json:"sell_price3"`
	NewPurchPrice1 *float64 `json:"new_purch_price1,omitempty" validate:"required"`
	NewPurchPrice2 *float64 `json:"new_purch_price2,omitempty" validate:"required"`
	NewPurchPrice3 *float64 `json:"new_purch_price3,omitempty" validate:"required"`
	NewSellPrice1  *float64 `json:"new_sell_price1,omitempty" validate:"required"`
	NewSellPrice2  *float64 `json:"new_sell_price2,omitempty" validate:"required"`
	NewSellPrice3  *float64 `json:"new_sell_price3,omitempty" validate:"required"`
}

type DetailMPriceParams struct {
	ParentCustID string
	CustID       string
	PriceID      string `params:"price_id" validate:"required"`
}

type UpdateMPriceParams struct {
	ParentCustID string
	CustID       string
	PriceID      string `params:"price_id" validate:"required"`
}

type CancelMPriceParams struct {
	UpdatedBy    string
	UpdatedByID  *int64
	ParentCustID string
	CustID       string
	PriceID      string `params:"price_id" validate:"required"`
}

type PublishMPriceParams struct {
	UpdatedBy     string
	UpdatedByID   *int64
	ParentCustID  string
	CustID        string
	DistributorID int64
	PriceID       string `params:"price_id" validate:"required"`
}

type DeleteMPriceParams struct {
	PriceId string `params:"price_id" validate:"required"`
}

type MPriceStatusDescSlice []MPriceStatus

func (s MPriceStatusDescSlice) Len() int {
	return len(s)
}

func (s MPriceStatusDescSlice) Less(i, j int) bool {
	return s[i].StatusID < s[j].StatusID
}

func (p MPriceStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

var MPriceStatusDesc = map[int]string{
	1: "Scheduled", 5: "Cancelled", 10: "Published",
}

type MPriceStatus struct {
	StatusID   int    `json:"status_id"`
	StatusDesc string `json:"status_desc"`
}

func (price MPriceResponse) GetPriceStatusDesc() string {
	return MPriceStatusDesc[price.Status]
}

type PublishByRmqMPriceReq struct {
	PriceID       string `json:"price_id"`
	CustID        string `json:"cust_id"`
	ParentCustID  string `json:"parent_cust_id"`
	DistributorID int64  `json:"distributor_id"`
	Status        int    `json:"status"`
	UpdatedBy     string `json:"updated_by"`
	UpdatedByID   *int64 `json:"updated_by_id,omitempty"`
}

type MPriceImportRequest struct {
	FileURL string `json:"file_url" validate:"required"`
}

type MPriceImportResponse struct {
	FileURL     string   `json:"file_url"`
	TotalRow    int      `json:"total_row"`
	SuccessRow  int      `json:"success_row"`
	FailedRow   int      `json:"failed_row"`
	ProcessedAt string   `json:"processed_at"`
	FailedRows  []string `json:"failed_rows"`
}
