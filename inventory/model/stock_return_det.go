package model

type StockReturnDetailRead struct {
	ReturnDetailID int64   `gorm:"column:return_detail_id;primaryKey" json:"return_detail_id"`
	OrderDetailID  *int64  `gorm:"column:order_detail_id" json:"order_detail_id"`
	ReturnNo       string  `gorm:"column:return_no" json:"return_no"`
	ProductID      int64   `gorm:"column:product_id" json:"product_id"`
	ProductCode    *string `gorm:"column:product_code" json:"product_code"`
	ProductName    *string `gorm:"column:product_name" json:"product_name"`
	ItemCnd        int64   `gorm:"column:item_cnd" json:"item_cnd"`
	Qty1           float64 `gorm:"column:qty1" json:"qty1"`
	Qty2           float64 `gorm:"column:qty2" json:"qty2"`
	Qty3           float64 `gorm:"column:qty3" json:"qty3"`
	InvoiceQty1    float64 `gorm:"column:invoice_qty1" json:"invoice_qty1"`
	InvoiceQty2    float64 `gorm:"column:invoice_qty2" json:"invoice_qty2"`
	InvoiceQty3    float64 `gorm:"column:invoice_qty3" json:"invoice_qty3"`
	// RemainingQty1    *float64 `gorm:"column:remaining_qty1" json:"remaining_qty1"`
	// RemainingQty2    *float64 `gorm:"column:remaining_qty2" json:"remaining_qty2"`
	// RemainingQty3    *float64 `gorm:"column:remaining_qty3" json:"remaining_qty3"`
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
	SubTotal         float64  `gorm:"column:sub_total" json:"sub_total"`
	Total            float64  `gorm:"column:total" json:"total"`
	ReturnReasonID   int64    `gorm:"column:return_reason_id" json:"return_reason_id"`
	ReturnReasonCode *string  `gorm:"column:return_reason_code" json:"return_reason_code"`
	ReturnReasonName *string  `gorm:"column:return_reason_name" json:"return_reason_name"`
	WhId             *int64   `gorm:"column:wh_id" json:"wh_id"`
	WhCode           string   `gorm:"column:wh_code" json:"wh_code"`
	WhName           string   `gorm:"column:wh_name" json:"wh_name"`
}

func (StockReturnDetailRead) TableName() string {
	return "sls.return_det"
}

type StockReturnedDetailRead struct {
	RemainingQty1 float64 `gorm:"column:remaining_qty1" json:"remaining_qty1"`
	RemainingQty2 float64 `gorm:"column:remaining_qty2" json:"remaining_qty2"`
	RemainingQty3 float64 `gorm:"column:remaining_qty3" json:"remaining_qty3"`
}

func (StockReturnedDetailRead) TableName() string {
	return "sls.return_det"
}

type StockReturnDetail struct {
	CustID         string   `gorm:"column:cust_id" json:"cust_id"`
	ReturnDetailID *int     `gorm:"column:return_detail_id;primaryKey" json:"return_detail_id"`
	OrderDetailID  *int64   `gorm:"column:order_detail_id" json:"order_detail_id"`
	ProductID      int      `gorm:"column:product_id" json:"product_id"`
	ItemCnd        int64    `gorm:"column:item_cnd" json:"item_cnd"`
	Qty1           *float64 `gorm:"column:qty1" json:"qty1"`
	Qty2           *float64 `gorm:"column:qty2" json:"qty2"`
	Qty3           *float64 `gorm:"column:qty3" json:"qty3"`
	WhID           *int64   `gorm:"column:wh_id" json:"wh_id"`
	ReturnReasonID int64    `gorm:"column:return_reason_id" json:"return_reason_id"`
}

func (StockReturnDetail) TableName() string {
	return "sls.return_det"
}
