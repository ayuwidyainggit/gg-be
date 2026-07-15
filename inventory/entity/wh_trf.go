package entity

import (
	"errors"
	"fmt"
)

type CreateWhTrfBody struct {
	CustID    string               `json:"cust_id"`
	WhIDFrom  int64                `json:"wh_id_from"`
	WhIDTo    int64                `json:"wh_id_to"`
	Notes     string               `json:"notes"`
	CreatedBy int64                `json:"created_by"`
	UpdatedBy *int64               `json:"updated_by"`
	Details   []CreateWhTrfDetBody `json:"details"`
}
type WhTrfResponse struct {
	WhTrfNo         string            `json:"stock_trf_no"`
	WhTrfDate       *string           `json:"stock_trf_date"`
	WhIDFrom        int64             `json:"wh_id_from"`
	WhCodeFrom      string            `json:"wh_code_from"`
	WhNameFrom      string            `json:"wh_name_from"`
	WhIDTo          int64             `json:"wh_id_to"`
	WhCodeTo        string            `json:"wh_code_to"`
	WhNameTo        string            `json:"wh_name_to"`
	SubTotal        float64           `json:"sub_total"`
	Total           float64           `json:"total"`
	TotalVat        float64           `json:"total_vat"`
	TotalVatLgPurch float64           `json:"total_vat_lg_purch"`
	TotalVatBg      float64           `json:"total_vat_bg"`
	Notes           string            `json:"notes"`
	DataStatus      int64             `json:"data_status"`
	CreatedBy       int64             `json:"created_by"`
	UpdatedAt       string            `json:"updated_at"`
	UpdatedByName   string            `json:"updated_by_name"`
	IsClosed        bool              `json:"is_closed"`
	ClosedBy        int64             `json:"closed_by"`
	ClosedByName    string            `json:"closed_by_name"`
	ClosedAt        string            `json:"closed_at"`
	Details         []WhTrfDetRespose `json:"details"`
}
type DetailWhTrfParams struct {
	WhTrfNo string `params:"wh_trf_no" validate:"required"`
}
type UpdateWhTrfParams struct {
	WhTrfNo string `params:"wh_trf_no" validate:"required"`
}
type WhTrfListResponse struct {
	WhTrfNo       string  `json:"stock_trf_no"`
	WhTrfDate     *string `json:"stock_trf_date"`
	WhIDFrom      int64   `json:"wh_id_from"`
	WhCodeFrom    string  `json:"wh_code_from"`
	WhNameFrom    string  `json:"wh_name_from"`
	WhIDTo        int64   `json:"wh_id_to"`
	WhCodeTo      string  `json:"wh_code_to"`
	WhNameTo      string  `json:"wh_name_to"`
	Notes         string  `json:"notes"`
	DataStatus    int64   `json:"data_status"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
}
type UpdateWhTrfBody struct {
	CustID     string               `json:"cust_id"`
	WhTrfNo    string               `json:"wh_trf_no"`
	WhTrfDate  *string              `json:"wh_trf_date"`
	TrCode     *string              `json:"tr_code"`
	DeliveryNo *string              `json:"delivery_no"`
	WhIDFrom   *int64               `json:"wh_id_from"`
	WhIDTo     *int64               `json:"wh_id_to"`
	CustIDTo   *string              `json:"cust_id_to"`
	SubTotal   *float64             `json:"sub_total"`
	Vat        *float64             `json:"vat"`
	VatValue   *float64             `json:"vat_value"`
	VatLg      *float64             `json:"vat_lg"`
	VatLgValue *float64             `json:"vat_lg_value"`
	Total      *float64             `json:"total"`
	VatBg      *float64             `json:"vat_bg"`
	VatBgValue *float64             `json:"vat_bg_value"`
	TotEmbInc  *float64             `json:"tot_emb_inc"`
	TotEmbExc  *float64             `json:"tot_emb_exc"`
	Notes      *string              `json:"notes"`
	DataStatus *int64               `json:"data_status"`
	UpdatedBy  int64                `json:"updated_by"`
	Details    []UpdateWhTrfDetBody `json:"details"`
}

type WhQueryFilter struct {
	CustId          string
	ParentCustId    string
	From            *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To              *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page            int    `query:"page"`
	Limit           int    `query:"limit" validate:"required"`
	Query           string `query:"q"`
	Mode            string `query:"mode"`
	Sort            string `query:"sort"`
	IsActive        *int   `query:"is_active"`
	TrCode          string `query:"tr_code"`
	StockTransferNo string `query:"stock_transfer_no"`
	GrType          *int   `query:"gr_type"`
	WhIdFrom        []int  `query:"warehouse_id_from"`
	WhIdTo          []int  `query:"warehouse_id_to"`
}

type StockTranferWarehouseQueryFilter struct {
	StartDate *int64 `query:"startDate"`
	EndDate   *int64 `query:"endDate"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Method    string
}
type StockTranferWarehouseQueryFilterParams struct {
	Method string `params:"method" validate:"oneof='from to'"`
}

type StockTranferWarehouse struct {
	WhId   *int64  `json:"wh_id"`
	WhCode *string `json:"wh_code"`
	WhName *string `json:"wh_name"`
}

type ProductListDistPriceResponse struct {
	Message string `json:"message"`
	Data    []ProductListDistPriceDataResp
}

type ProductListDistPriceDataResp struct {
	ProID       int64   `json:"pro_id"`
	PurchPrice1 float64 `json:"purch_price1"`
	PurchPrice2 float64 `json:"purch_price2"`
	PurchPrice3 float64 `json:"purch_price3"`
	SellPrice1  float64 `json:"sell_price1"`
	SellPrice2  float64 `json:"sell_price2"`
	SellPrice3  float64 `json:"sell_price3"`
}

type MapProductDistPrice map[int64]ProductListDistPriceDataResp

func (m MapProductDistPrice) SetProduct(id int64, product ProductListDistPriceDataResp) {
	m[id] = product
}

func (m MapProductDistPrice) GetByID(id int64) (product ProductListDistPriceDataResp, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return product, errors.New(fmt.Sprintf("Product ID %v Not Found", id))
	}

	return val, nil
}
