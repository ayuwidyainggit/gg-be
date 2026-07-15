package model

/*
type ReturnDet struct {
	CustID         string     `gorm:"column:cust_id" json:"cust_id"`
	ReturnNo       string     `gorm:"column:return_no" json:"return_no"`
	SeqNo          int        `gorm:"column:seq_no" json:"seq_no"`
	ReturnDetID    *int       `gorm:"column:return_det_id;primaryKey" json:"return_det_id"`
	ProID          int        `gorm:"column:pro_id" json:"pro_id"`
	ItemType       int        `gorm:"column:item_type" json:"item_type"`
	ItemCnd        *int64     `gorm:"column:item_cnd" json:"item_cnd"`
	SellPrice1     *float64   `gorm:"column:sell_price1" json:"sell_price_1"`
	SellPrice2     *float64   `gorm:"column:sell_price2" json:"sell_price_2"`
	SellPrice3     *float64   `gorm:"column:sell_price3" json:"sell_price_3"`
	Qty            *float64   `gorm:"column:qty" json:"qty"`
	QtyStr         *string    `gorm:"column:qty_str" json:"qty_str"`
	TotAmount      *float64   `gorm:"column:tot_amount" json:"tot_amount"`
	Disc1P         *float64   `gorm:"column:disc1_p" json:"disc_1_p"`
	Disc1Value     *float64   `gorm:"column:disc1_value" json:"disc_1_value"`
	Disc2P         *float64   `gorm:"column:disc2_p" json:"disc_2_p"`
	Disc2Value     *float64   `gorm:"column:disc2_value" json:"disc_2_value"`
	DiscTotal      *float64   `gorm:"column:disc_total" json:"disc_total"`
	TotAmountNet   *float64   `gorm:"column:tot_amount_net" json:"tot_amount_net"`
	ReturnReasonID *int64     `gorm:"column:return_reason_id" json:"return_reason_id"`
	BatchNo        *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate        *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (ReturnDet) TableName() string {
	return "sls.return_det"
}

type ReturnDetRead struct {
	CustID         string     `gorm:"column:cust_id" json:"cust_id"`
	ReturnNo       string     `gorm:"column:return_no" json:"return_no"`
	SeqNo          int        `gorm:"column:seq_no" json:"seq_no"`
	ReturnDetID    *int       `gorm:"column:return_det_id;primaryKey" json:"return_det_id"`
	ProID          int        `gorm:"column:pro_id" json:"pro_id"`
	ProCode        string     `gorm:"column:pro_code" json:"pro_code"`
	ProName        string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType       int        `gorm:"column:item_type" json:"item_type"`
	ItemCnd        *int64     `gorm:"column:item_cnd" json:"item_cnd"`
	SellPrice1     *float64   `gorm:"column:sell_price1" json:"sell_price_1"`
	SellPrice2     *float64   `gorm:"column:sell_price2" json:"sell_price_2"`
	SellPrice3     *float64   `gorm:"column:sell_price3" json:"sell_price_3"`
	Qty            *float64   `gorm:"column:qty" json:"qty"`
	QtyStr         *string    `gorm:"column:qty_str" json:"qty_str"`
	TotAmount      *float64   `gorm:"column:tot_amount" json:"tot_amount"`
	Disc1P         *float64   `gorm:"column:disc1_p" json:"disc_1_p"`
	Disc1Value     *float64   `gorm:"column:disc1_value" json:"disc_1_value"`
	Disc2P         *float64   `gorm:"column:disc2_p" json:"disc_2_p"`
	Disc2Value     *float64   `gorm:"column:disc2_value" json:"disc_2_value"`
	DiscTotal      *float64   `gorm:"column:disc_total" json:"disc_total"`
	TotAmountNet   *float64   `gorm:"column:tot_amount_net" json:"tot_amount_net"`
	ReturnReasonID *int64     `gorm:"column:return_reason_id" json:"return_reason_id"`
	BatchNo        *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate        *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (ReturnDetRead) TableName() string {
	return "sls.return_det"
}
*/
type ReturnDetailRead struct {
	ReturnDetailID   int64    `gorm:"column:return_detail_id;primaryKey" json:"return_detail_id"`
	OrderDetailID    *int64   `gorm:"column:order_detail_id" json:"order_detail_id"`
	ReturnNo         string   `gorm:"column:return_no" json:"return_no"`
	ProductID        int64    `gorm:"column:product_id" json:"product_id"`
	ProductCode      *string  `gorm:"column:product_code" json:"product_code"`
	ProductName      *string  `gorm:"column:product_name" json:"product_name"`
	ItemType         int64    `gorm:"column:item_type" json:"item_type"`
	ItemCnd          int64    `gorm:"column:item_cnd" json:"item_cnd"`
	WhID             int64    `gorm:"column:wh_id" json:"wh_id"`
	WhCode           *string  `gorm:"column:wh_code" json:"wh_code"`
	WhName           *string  `gorm:"column:wh_name" json:"wh_name"`
	Qty              float64  `gorm:"column:qty" json:"qty"`
	Qty1             float64  `gorm:"column:qty1" json:"qty1"`
	Qty2             float64  `gorm:"column:qty2" json:"qty2"`
	Qty3             float64  `gorm:"column:qty3" json:"qty3"`
	InvoiceQty1      float64  `gorm:"column:invoice_qty1" json:"invoice_qty1"`
	InvoiceQty2      float64  `gorm:"column:invoice_qty2" json:"invoice_qty2"`
	InvoiceQty3      float64  `gorm:"column:invoice_qty3" json:"invoice_qty3"`
	RemainingQty1    *float64 `gorm:"column:remaining_qty1" json:"remaining_qty1"`
	RemainingQty2    *float64 `gorm:"column:remaining_qty2" json:"remaining_qty2"`
	RemainingQty3    *float64 `gorm:"column:remaining_qty3" json:"remaining_qty3"`
	SellPrice1       float64  `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2       float64  `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3       float64  `gorm:"column:sell_price3" json:"sell_price3"`
	UnitId1          string   `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2          string   `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3          string   `gorm:"column:unit_id3" json:"unit_id3"`
	UnitName1        *string  `gorm:"column:unit_name1" json:"unit_name1"`
	UnitName2        *string  `gorm:"column:unit_name2" json:"unit_name2"`
	UnitName3        *string  `gorm:"column:unit_name3" json:"unit_name3"`
	ConvUnit2        *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3        *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	Vat              *float64 `gorm:"column:vat" json:"vat"`
	VatValue         *float64 `gorm:"column:vat_value" json:"vat_value"`
	DiscValue        *float64 `gorm:"column:disc_value" json:"disc_value"`
	PromoValue       *float64 `gorm:"column:promo_value" json:"promo_value"`
	SubTotal         float64  `gorm:"column:sub_total" json:"sub_total"`
	Total            float64  `gorm:"column:total" json:"total"`
	ReturnReasonID   int64    `gorm:"column:return_reason_id" json:"return_reason_id"`
	ReturnReasonCode *string  `gorm:"column:return_reason_code" json:"return_reason_code"`
	ReturnReasonName *string  `gorm:"column:return_reason_name" json:"return_reason_name"`
	Volume           float64  `gorm:"column:volume" json:"volume"`
	Weight           float64  `gorm:"column:weight" json:"weight"`
	Volume1          float64  `gorm:"column:volume1" json:"volume1"`
	Volume2          float64  `gorm:"column:volume2" json:"volume2"`
	Volume3          float64  `gorm:"column:volume3" json:"volume3"`
	Weight1          float64  `gorm:"column:weight1" json:"weight1"`
	Weight2          float64  `gorm:"column:weight2" json:"weight2"`
	Weight3          float64  `gorm:"column:weight3" json:"weight3"`
}

func (ReturnDetailRead) TableName() string {
	return "sls.return_det"
}

type ReturnedDetailRead struct {
	ReturnedQty1 float64 `gorm:"column:returned_qty1" json:"returned_qty1"`
	ReturnedQty2 float64 `gorm:"column:returned_qty2" json:"returned_qty2"`
	ReturnedQty3 float64 `gorm:"column:returned_qty3" json:"returned_qty3"`
}

func (ReturnedDetailRead) TableName() string {
	return "sls.return_det"
}

type ReturnDetail struct {
	CustID         string   `gorm:"column:cust_id" json:"cust_id"`
	ReturnDetailID *int     `gorm:"column:return_detail_id;primaryKey" json:"return_detail_id"`
	ReturnNo       string   `gorm:"column:return_no" json:"return_no"`
	OrderDetailID  *int64   `gorm:"column:order_detail_id" json:"order_detail_id"`
	ProductID      int      `gorm:"column:product_id" json:"product_id"`
	ItemType       int64    `gorm:"column:item_type" json:"item_type"`
	ItemCnd        int64    `gorm:"column:item_cnd" json:"item_cnd"`
	WhId           int64    `gorm:"column:wh_id" json:"wh_id"`
	Qty            float64  `gorm:"column:qty" json:"qty"`
	Qty1           float64  `gorm:"column:qty1" json:"qty1"`
	Qty2           float64  `gorm:"column:qty2" json:"qty2"`
	Qty3           float64  `gorm:"column:qty3" json:"qty3"`
	UnitId1        string   `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2        string   `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3        string   `gorm:"column:unit_id3" json:"unit_id3"`
	SellPrice1     float64  `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2     float64  `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3     float64  `gorm:"column:sell_price3" json:"sell_price3"`
	PromoValue     *float64 `gorm:"column:promo_value" json:"promo_value"`
	DiscValue      *float64 `gorm:"column:disc_value" json:"disc_value"`
	ConvUnit2      *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3      *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	Vat            float64  `gorm:"column:vat" json:"vat"`
	VatValue       float64  `gorm:"column:vat_value" json:"vat_value"`
	SubTotal       float64  `gorm:"column:sub_total" json:"sub_total"`
	Total          float64  `gorm:"column:total" json:"total"`
	ReturnReasonID int64    `gorm:"column:return_reason_id" json:"return_reason_id"`
}

func (ReturnDetail) TableName() string {
	return "sls.return_det"
}

type ReturnQuantity struct {
	CustID         string  `gorm:"column:cust_id" json:"cust_id"`
	ReturnDetailID *int    `gorm:"column:return_detail_id;primaryKey" json:"return_detail_id"`
	Qty1           float64 `gorm:"column:qty1" json:"qty1"`
	Qty2           float64 `gorm:"column:qty2" json:"qty2"`
	Qty3           float64 `gorm:"column:qty3" json:"qty3"`
	SellPrice1     float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2     float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3     float64 `gorm:"column:sell_price3" json:"sell_price3"`
	Vat            float64 `gorm:"column:vat" json:"vat"`
	VatValue       float64 `gorm:"column:vat_value" json:"vat_value"`
	SubTotal       float64 `gorm:"column:sub_total" json:"sub_total"`
	Total          float64 `gorm:"column:total" json:"total"`
}

func (ReturnQuantity) TableName() string {
	return "sls.return_det"
}
