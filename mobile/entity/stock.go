package entity

type StockQueryFilter struct {
	Page        int    `query:"page"`
	Limit       int    `query:"limit" validate:"required"`
	Query       string `query:"q"`
	Mode        string `query:"mode"`
	Sort        string `query:"sort"`
	ProductName string `query:"product_name"`
	ProductCode string `query:"product_code"`
}

type StockGudangUtamaListResponse struct {
	CustId         string                     `json:"cust_id"`
	EmpId          *int64                     `json:"emp_id"`
	SalesName      string                     `json:"sales_name"`
	WhId           int64                      `json:"wh_id"`
	WhCode         string                     `json:"wh_code"`
	WhName         string                     `json:"wh_name"`
	DetailsProduct []DetilsGudangUtamaProduct `json:"details_product"`
}

type DetilsGudangUtamaProduct struct {
	WhId               int64   `json:"wh_id"`
	ProId              int64   `json:"pro_id"`
	ProCode            string  `json:"pro_code"`
	ProName            string  `json:"pro_name"`
	UnitId1            string  `json:"unit_id1"`
	UnitId2            string  `json:"unit_id2"`
	UnitId3            string  `json:"unit_id3"`
	PurchPrice1        float64 `json:"purch_price1"`
	PurchPrice2        float64 `json:"purch_price2"`
	PurchPrice3        float64 `json:"purch_price3"`
	SellPrice1         float64 `json:"sell_price1"`
	SellPrice2         float64 `json:"sell_price2"`
	SellPrice3         float64 `json:"sell_price3"`
	TotalQty           float64 `json:"total_qty"`
	Qty1               float64 `json:"qty1"`
	Qty2               float64 `json:"qty2"`
	Qty3               float64 `json:"qty3"`
	TotalQtyOrder      float64 `json:"total_qty_order"`
	QtyOrder1          float64 `json:"qty_order1"`
	QtyOrder2          float64 `json:"qty_order2"`
	QtyOrder3          float64 `json:"qty_order3"`
	TotalQtyIncOnOrder float64 `json:"total_qty_inc_on_order"`
	QtyIncOnOrder1     float64 `json:"qty_inc_on_order1"`
	QtyIncOnOrder2     float64 `json:"qty_inc_on_order2"`
	QtyIncOnOrder3     float64 `json:"qty_inc_on_order3"`
	ConvUnit2          int     `json:"conv_unit2"`
	ConvUnit3          int     `json:"conv_unit3"`
	IsActive           bool    `json:"is_active"`
	Vat                float64 `json:"vat"`
	VatLgPurch         float64 `json:"vat_lg_purch"`
	VatLgSell          float64 `json:"vat_lg_sell"`
}

type StockGudangCanvasistResponse struct {
	CustId         string                      `json:"cust_id"`
	EmpId          *int64                      `json:"emp_id"`
	SalesName      string                      `json:"sales_name"`
	WhId           int64                       `json:"wh_id"`
	WhCode         string                      `json:"wh_code"`
	WhName         string                      `json:"wh_name"`
	DetailsProduct []DetilsGudangCanvasProduct `json:"details_product"`
}

type DetilsGudangCanvasProduct struct {
	WhId              int64   `json:"wh_id"`
	ProId             int64   `json:"pro_id"`
	ProCode           string  `json:"pro_code"`
	ProName           string  `json:"pro_name"`
	UnitId1           string  `json:"unit_id1"`
	UnitId2           string  `json:"unit_id2"`
	UnitId3           string  `json:"unit_id3"`
	PurchPrice1       float64 `json:"purch_price1"`
	PurchPrice2       float64 `json:"purch_price2"`
	PurchPrice3       float64 `json:"purch_price3"`
	SellPrice1        float64 `json:"sell_price1"`
	SellPrice2        float64 `json:"sell_price2"`
	SellPrice3        float64 `json:"sell_price3"`
	TotalQtyAvailable float64 `json:"total_qty_available"`
	Qty1Available     float64 `json:"qty1_available"`
	Qty2Available     float64 `json:"qty2_available"`
	Qty3Available     float64 `json:"qty3_available"`
	TotalQtyStock     float64 `json:"total_qty_stock"`
	Qty1Stock         float64 `json:"qty1_stock"`
	Qty2Stock         float64 `json:"qty2_stock"`
	Qty3Stock         float64 `json:"qty3_stock"`
	// TotalQtyOrder      float64 `json:"total_qty_order"`
	// QtyOrder1 float64 `json:"qty_order1"`
	// QtyOrder2 float64 `json:"qty_order2"`
	// QtyOrder3 float64 `json:"qty_order3"`
	// TotalQtyIncOnOrder float64 `json:"total_qty_inc_on_order"`
	// QtyIncOnOrder1     float64 `json:"qty_inc_on_order1"`
	// QtyIncOnOrder2     float64 `json:"qty_inc_on_order2"`
	// QtyIncOnOrder3     float64 `json:"qty_inc_on_order3"`
	ConvUnit2  int     `json:"conv_unit2"`
	ConvUnit3  int     `json:"conv_unit3"`
	IsActive   bool    `json:"is_active"`
	Vat        float64 `json:"vat"`
	VatLgPurch float64 `json:"vat_lg_purch"`
	VatLgSell  float64 `json:"vat_lg_sell"`
}
