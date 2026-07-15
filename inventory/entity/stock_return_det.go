package entity

type ReturnDetailResponse struct {
	ReturnDetailID   int64    `json:"return_detail_id"`
	OrderDetailID    *int64   `json:"order_detail_id"`
	ReturnNo         string   `json:"return_no"`
	ProductID        int64    `json:"product_id"`
	ProductCode      *string  `json:"product_code"`
	ProductName      *string  `json:"product_name"`
	ItemCnd          int64    `json:"item_cnd"`
	ItemCndName      *string  `json:"item_cnd_name"`
	Qty1             float64  `json:"qty1"`
	Qty2             float64  `json:"qty2"`
	Qty3             float64  `json:"qty3"`
	InvoiceQty1      float64  `json:"invoice_qty1"`
	InvoiceQty2      float64  `json:"invoice_qty2"`
	InvoiceQty3      float64  `json:"invoice_qty3"`
	RemainingQty1    float64  `json:"remaining_qty1"`
	RemainingQty2    float64  `json:"remaining_qty2"`
	RemainingQty3    float64  `json:"remaining_qty3"`
	SellPrice1       float64  `json:"sell_price1"`
	SellPrice2       float64  `json:"sell_price2"`
	SellPrice3       float64  `json:"sell_price3"`
	UnitId1          string   `json:"unit_id1"`
	UnitId2          string   `json:"unit_id2"`
	UnitId3          string   `json:"unit_id3"`
	UnitName1        *string  `json:"unit_name1"`
	UnitName2        *string  `json:"unit_name2"`
	UnitName3        *string  `json:"unit_name3"`
	ConvUnit2        *float64 `json:"conv_unit2"`
	ConvUnit3        *float64 `json:"conv_unit3"`
	Vat              *float64 `json:"vat"`
	VatValue         *float64 `json:"vat_value"`
	SubTotal         float64  `json:"sub_total"`
	Total            float64  `json:"total"`
	ReturnReasonID   int64    `json:"return_reason_id"`
	ReturnReasonCode *string  `json:"return_reason_code"`
	ReturnReasonName *string  `json:"return_reason_name"`
}

var dataItemConditionName = map[int64]string{
	1: "Good",
	2: "Bad",
	3: "Expired",
}

func (returnDetail ReturnDetailResponse) GenerateItemConditionName() string {
	if returnDetail.ItemCnd != 0 {
		return dataItemConditionName[returnDetail.ItemCnd]
	}
	return ""
}

type StockReturnDetailResponse struct {
	ReturnDetailID   int64    `json:"return_detail_id"`
	OrderDetailID    *int64   `json:"order_detail_id"`
	ReturnNo         string   `json:"return_no"`
	ProductID        int64    `json:"product_id"`
	ProductCode      *string  `json:"product_code"`
	ProductName      *string  `json:"product_name"`
	ItemCnd          int64    `json:"item_cnd"`
	ItemCndName      *string  `json:"item_cnd_name"`
	Qty1             float64  `json:"qty1"`
	Qty2             float64  `json:"qty2"`
	Qty3             float64  `json:"qty3"`
	InvoiceQty1      float64  `json:"invoice_qty1"`
	InvoiceQty2      float64  `json:"invoice_qty2"`
	InvoiceQty3      float64  `json:"invoice_qty3"`
	RemainingQty1    float64  `json:"remaining_qty1"`
	RemainingQty2    float64  `json:"remaining_qty2"`
	RemainingQty3    float64  `json:"remaining_qty3"`
	SellPrice1       float64  `json:"sell_price1"`
	SellPrice2       float64  `json:"sell_price2"`
	SellPrice3       float64  `json:"sell_price3"`
	UnitId1          string   `json:"unit_id1"`
	UnitId2          string   `json:"unit_id2"`
	UnitId3          string   `json:"unit_id3"`
	UnitName1        *string  `json:"unit_name1"`
	UnitName2        *string  `json:"unit_name2"`
	UnitName3        *string  `json:"unit_name3"`
	ConvUnit2        *float64 `json:"conv_unit2"`
	ConvUnit3        *float64 `json:"conv_unit3"`
	Vat              *float64 `json:"vat"`
	VatValue         *float64 `json:"vat_value"`
	SubTotal         float64  `json:"sub_total"`
	Total            float64  `json:"total"`
	ReturnReasonID   int64    `json:"return_reason_id"`
	ReturnReasonCode *string  `json:"return_reason_code"`
	ReturnReasonName *string  `json:"return_reason_name"`
	WhID             *int64   `json:"wh_id"`
	WhCode           string   `json:"wh_code"`
	Whname           string   `json:"wh_name"`
}

func (returnDetail StockReturnDetailResponse) GenerateItemConditionName() string {
	if returnDetail.ItemCnd != 0 {
		return dataItemConditionName[returnDetail.ItemCnd]
	}
	return ""
}

type StockReturnDetailUpdateBody struct {
	ReturnDetailID int64   `json:"return_detail_id"`
	ProductID      int64   `json:"product_id"`
	Qty1           float64 `json:"qty1"`
	Qty2           float64 `json:"qty2"`
	Qty3           float64 `json:"qty3"`
	ItemCnd        int64   `json:"item_cnd"`
	ReturnReasonID *int64  `json:"return_reason_id"`
	WhID           *int64  `json:"wh_id"`
}
