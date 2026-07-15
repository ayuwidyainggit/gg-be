package entity

type CreateReturnDetBodyGroup struct {
	Normal []CreateReturnDetBody `json:"normal"`
	Promo  []CreateReturnDetBody `json:"promo"`
}
type ReturnDetReponseBodyGroup struct {
	Normal []ReturnDetReponse `json:"normal"`
	Promo  []ReturnDetReponse `json:"promo"`
}
type UpdateReturnDetBodyGroup struct {
	Normal []UpdateReturnDetBody `json:"normal"`
	Promo  []UpdateReturnDetBody `json:"promo"`
}
type CreateReturnDetBody struct {
	SeqNo          int      `json:"seq_no"`
	ProID          int      `json:"pro_id"`
	ItemType       int      `json:"item_type"`
	ItemCnd        *int64   `json:"item_cnd"`
	SellPrice1     *float64 `json:"sell_price_1"`
	SellPrice2     *float64 `json:"sell_price_2"`
	SellPrice3     *float64 `json:"sell_price_3"`
	Qty            *float64 `json:"qty"`
	QtyStr         *string  `json:"qty_str"`
	TotAmount      *float64 `json:"tot_amount"`
	Disc1P         *float64 `json:"disc_1_p"`
	Disc1Value     *float64 `json:"disc_1_value"`
	Disc2P         *float64 `json:"disc_2_p"`
	Disc2Value     *float64 `json:"disc_2_value"`
	DiscTotal      *float64 `json:"disc_total"`
	TotAmountNet   *float64 `json:"tot_amount_net"`
	ReturnReasonID *int64   `json:"return_reason_id"`
	BatchNo        *string  `json:"batch_no"`
	ExpDate        *string  `json:"exp_date"`
}

type ReturnDetReponse struct {
	ReturnDetID    int64    `json:"return_det_id"`
	SeqNo          int      `json:"seq_no"`
	ProID          int      `json:"pro_id"`
	ProCode        string   `json:"pro_code"`
	ProName        string   `json:"pro_name"`
	ItemType       int      `json:"item_type"`
	ItemCnd        *int64   `json:"item_cnd"`
	SellPrice1     *float64 `json:"sell_price_1"`
	SellPrice2     *float64 `json:"sell_price_2"`
	SellPrice3     *float64 `json:"sell_price_3"`
	Qty            *float64 `json:"qty"`
	QtyStr         *string  `json:"qty_str"`
	TotAmount      *float64 `json:"tot_amount"`
	Disc1P         *float64 `json:"disc_1_p"`
	Disc1Value     *float64 `json:"disc_1_value"`
	DiscValue      *float64 `json:"disc_value"`
	PromoValue     *float64 `json:"promo_value"`
	VatValue       *float64 `json:"vat_value"`
	Disc2P         *float64 `json:"disc_2_p"`
	Disc2Value     *float64 `json:"disc_2_value"`
	DiscTotal      *float64 `json:"disc_total"`
	TotAmountNet   *float64 `json:"tot_amount_net"`
	ReturnReasonID *int64   `json:"return_reason_id"`
	BatchNo        *string  `json:"batch_no"`
	ExpDate        *string  `json:"exp_date"`
}
type UpdateReturnDetBody struct {
	ReturnDetID    *int64   `json:"return_det_id"`
	SeqNo          int      `json:"seq_no"`
	ProID          int      `json:"pro_id"`
	ItemType       int      `json:"item_type"`
	ItemCnd        *int64   `json:"item_cnd"`
	SellPrice1     *float64 `json:"sell_price_1"`
	SellPrice2     *float64 `json:"sell_price_2"`
	SellPrice3     *float64 `json:"sell_price_3"`
	Qty            *float64 `json:"qty"`
	QtyStr         *string  `json:"qty_str"`
	TotAmount      *float64 `json:"tot_amount"`
	Disc1P         *float64 `json:"disc_1_p"`
	Disc1Value     *float64 `json:"disc_1_value"`
	Disc2P         *float64 `json:"disc_2_p"`
	Disc2Value     *float64 `json:"disc_2_value"`
	DiscTotal      *float64 `json:"disc_total"`
	TotAmountNet   *float64 `json:"tot_amount_net"`
	ReturnReasonID *int64   `json:"return_reason_id"`
	BatchNo        *string  `json:"batch_no"`
	ExpDate        *string  `json:"exp_date"`
}

type ReturnDetailResponse struct {
	ReturnDetailID   int64    `json:"return_detail_id"`
	OrderDetailID    *int64   `json:"order_detail_id"`
	ReturnNo         string   `json:"return_no"`
	ProductID        int64    `json:"product_id"`
	ProductCode      *string  `json:"product_code"`
	ProductName      *string  `json:"product_name"`
	WhID             int64    `json:"wh_id"`
	WhCode           *string  `json:"wh_code"`
	WhName           *string  `json:"wh_name"`
	ItemType         int64    `json:"item_type"`
	ItemCnd          int64    `json:"item_cnd"`
	ItemCndName      *string  `json:"item_cnd_name"`
	Qty              float64  `json:"qty"`
	Qty1             float64  `json:"qty1"`
	Qty2             float64  `json:"qty2"`
	Qty3             float64  `json:"qty3"`
	Volume1          float64  `json:"volume1"`
	Volume2          float64  `json:"volume2"`
	Volume3          float64  `json:"volume3"`
	Weight1          float64  `json:"weight1"`
	Weight2          float64  `json:"weight2"`
	Weight3          float64  `json:"weight3"`
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
	DiscValue        *float64 `json:"disc_value"`
	PromoValue       *float64 `json:"promo_value"`
	VatValue         *float64 `json:"vat_value"`
	SubTotal         float64  `json:"sub_total"`
	Total            float64  `json:"total"`
	ReturnReasonID   int64    `json:"return_reason_id"`
	ReturnReasonCode *string  `json:"return_reason_code"`
	ReturnReasonName *string  `json:"return_reason_name"`
	Volume           float64  `json:"volume"`
	Weight           float64  `json:"weight"`
}

type CreateReturnDetailBody struct {
	// CustID string `json:"cust_id"`
	// ReturnDetailID int64  `json:"return_detail_id"`
	OrderDetailID *int64   `json:"order_detail_id"`
	InvoiceNo     *string  `json:"invoice_no"`
	InvoiceDate   *string  `json:"invoice_date"`
	SalesmanID    *int64   `json:"salesman_id"`
	WhID          *int64   `json:"wh_id"`
	OutletID      *int64   `json:"outlet_id"`
	ReturnNo      *string  `json:"return_no"`
	SeqNo         int      `json:"seq_no"`
	ProductID     int64    `json:"product_id"`
	Qty1          float64  `json:"qty1"`
	Qty2          float64  `json:"qty2"`
	Qty3          float64  `json:"qty3"`
	ItemCnd       int64    `json:"item_cnd"`
	SellPrice1    float64  `json:"sell_price1"`
	SellPrice2    float64  `json:"sell_price2"`
	SellPrice3    float64  `json:"sell_price3"`
	UnitId1       string   `json:"unit_id1"`
	UnitId2       string   `json:"unit_id2"`
	UnitId3       string   `json:"unit_id3"`
	ConvUnit2     *float64 `json:"conv_unit2"`
	ConvUnit3     *float64 `json:"conv_unit3"`
	Vat           *float64 `json:"vat"`
	VatValue      float64  `json:"vat_value"`
	SubTotal      float64  `json:"sub_total"`
	Total         float64  `json:"total"`
	// CreatedBy      *int64   `json:"created_by"`
	ReturnReasonID *int64 `json:"return_reason_id"`
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

type UpdateReturnDetailBody struct {
	ReturnDetailID int64    `json:"return_detail_id"`
	OrderDetailID  *int64   `json:"order_detail_id"`
	ReturnNo       *string  `json:"return_no"`
	ProductID      int64    `json:"product_id"`
	WhID           int64    `json:"wh_id"`
	Qty1           float64  `json:"qty1"`
	Qty2           float64  `json:"qty2"`
	Qty3           float64  `json:"qty3"`
	ItemCnd        int64    `json:"item_cnd"`
	SellPrice1     float64  `json:"sell_price1"`
	SellPrice2     float64  `json:"sell_price2"`
	SellPrice3     float64  `json:"sell_price3"`
	UnitId1        string   `json:"unit_id1"`
	UnitId2        string   `json:"unit_id2"`
	UnitId3        string   `json:"unit_id3"`
	ConvUnit2      *float64 `json:"conv_unit2"`
	ConvUnit3      *float64 `json:"conv_unit3"`
	Vat            *float64 `json:"vat"`
	VatValue       float64  `json:"vat_value"`
	SubTotal       float64  `json:"sub_total"`
	Total          float64  `json:"total"`
	ReturnReasonID *int64   `json:"return_reason_id"`
}

type UpdateReturnQuantityBody struct {
	ReturnDetailID int64    `json:"return_detail_id"`
	Qty1           float64  `json:"qty1"`
	Qty2           float64  `json:"qty2"`
	Qty3           float64  `json:"qty3"`
	SellPrice1     float64  `json:"sell_price1"`
	SellPrice2     float64  `json:"sell_price2"`
	SellPrice3     float64  `json:"sell_price3"`
	Vat            *float64 `json:"vat"`
	VatValue       float64  `json:"vat_value"`
	SubTotal       float64  `json:"sub_total"`
	Total          float64  `json:"total"`
}

type ApproveReturnDetailBody struct {
	ReturnDetailID int64    `json:"return_detail_id"`
	OrderDetailID  *int64   `json:"order_detail_id"`
	ReturnNo       *string  `json:"return_no"`
	ProductID      int64    `json:"product_id"`
	WhID           int64    `json:"wh_id"`
	Qty1           float64  `json:"qty1"`
	Qty2           float64  `json:"qty2"`
	Qty3           float64  `json:"qty3"`
	ItemCnd        int64    `json:"item_cnd"`
	SellPrice1     float64  `json:"sell_price1"`
	SellPrice2     float64  `json:"sell_price2"`
	SellPrice3     float64  `json:"sell_price3"`
	UnitId1        string   `json:"unit_id1"`
	UnitId2        string   `json:"unit_id2"`
	UnitId3        string   `json:"unit_id3"`
	ConvUnit2      *float64 `json:"conv_unit2"`
	ConvUnit3      *float64 `json:"conv_unit3"`
	Vat            *float64 `json:"vat"`
	VatValue       float64  `json:"vat_value"`
	SubTotal       float64  `json:"sub_total"`
	Total          float64  `json:"total"`
	ReturnReasonID *int64   `json:"return_reason_id"`
}

type UpdateStatusReturnDetailBody struct {
	ReturnNo   string `json:"return_no"`
	DataStatus int64  `json:"status"`
}

type UpdateAssignReturnDetailBody struct {
	ReturnNo   string `json:"return_no"`
	DataStatus int64  `json:"status"`
	EmpId      int64  `json:"emp_id"`
}

type CreateReturnDetailRequestBody struct {
	CustID string `json:"cust_id"`
	// ReturnDetailID int      `json:"return_detail_id"`
	// ReturnNo       string   `json:"return_no"`
	OrderDetailID *int64 `json:"order_detail_id"`
	ProductID     int    `json:"product_id"`
	ItemCnd       int64  `json:"item_cnd"`
	// ItemType       *int64   `json:"item_type"`
	WhId           int64    `json:"wh_id"`
	Qty1           float64  `json:"qty1"`
	Qty2           float64  `json:"qty2"`
	Qty3           float64  `json:"qty3"`
	UnitId1        string   `json:"unit_id1"`
	UnitId2        string   `json:"unit_id2"`
	UnitId3        string   `json:"unit_id3"`
	SellPrice1     float64  `json:"sell_price1"`
	SellPrice2     float64  `json:"sell_price2"`
	SellPrice3     float64  `json:"sell_price3"`
	ConvUnit2      *float64 `json:"conv_unit2"`
	ConvUnit3      *float64 `json:"conv_unit3"`
	Vat            float64  `json:"vat"`
	VatValue       float64  `json:"vat_value"`
	SubTotal       float64  `json:"sub_total"`
	Total          float64  `json:"total"`
	ReturnReasonID int64    `json:"return_reason_id"`
}
