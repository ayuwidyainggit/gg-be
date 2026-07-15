package entity

import "time"

type StockQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	SupId        *int   `query:"sup_id"`
	WhId         *int   `query:"wh_id"`
}

type CheckStockQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64  `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64  `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit" validate:"required"`
	Query        string  `query:"q"`
	Mode         string  `query:"mode"`
	Sort         string  `query:"sort"`
	IsActive     *int    `query:"is_active"`
	SupId        *int    `query:"sup_id"`
	WhId         *int    `query:"wh_id"`
	ProId        *int    `query:"pro_id"`
	ProCode      *string `query:"pro_code"`
	ProName      *string `query:"pro_name"`
}

type Stock struct {
	CustID        string  `json:"cust_id"`
	GrDate        string  `json:"gr_date"`
	TrCode        string  `json:"tr_code"`
	DeliveryDate  string  `json:"delivery_date"`
	DeliveryNo    string  `json:"delivery_no"`
	InvoiceNo     string  `json:"invoice_no"`
	VehicleNo     string  `json:"vehicle_no"`
	PoNo          string  `json:"po_no"`
	PoDnNo        string  `json:"po_dn_no"`
	SupID         int64   `json:"sup_id"`
	SupName       string  `json:"sup_name"`
	WhID          int64   `json:"wh_id"`
	Notes         string  `json:"notes"`
	TotEmbInc     float64 `json:"tot_emb_inc"`
	TotEmbExc     float64 `json:"tot_emb_exc"`
	GrType        int64   `json:"gr_type"`
	SubTotal      float64 `json:"sub_total"`
	Vat           float64 `json:"vat"`
	VatValue      float64 `json:"vat_value"`
	VatLg         float64 `json:"vat_lg"`
	VatLgValue    float64 `json:"vat_lg_value"`
	Total         float64 `json:"total"`
	VatBg         float64 `json:"vat_bg"`
	VatBgValue    float64 `json:"vat_bg_value"`
	WeekIDOb      int     `json:"week_id_ob"`
	DataStatus    int64   `json:"data_status"`
	CreatedBy     int64   `json:"created_by"`
	UpdatedAt     string  `json:"updated_by"`
	UpdatedByName string  `json:"updated_by_name"`
}
type StockResponse struct {
	GrNo          string  `json:"gr_no"`
	GrDate        string  `json:"gr_date"`
	TrCode        string  `json:"tr_code"`
	DeliveryDate  string  `json:"delivery_date"`
	DeliveryNo    string  `json:"delivery_no"`
	InvoiceNo     string  `json:"invoice_no"`
	VehicleNo     string  `json:"vehicle_no"`
	PoNo          string  `json:"po_no"`
	PoDnNo        string  `json:"po_dn_no"`
	SupID         int64   `json:"sup_id"`
	SupName       string  `json:"sup_name"`
	WhID          int64   `json:"wh_id"`
	WhName        string  `json:"wh_name"`
	Notes         string  `json:"notes"`
	TotEmbInc     float64 `json:"tot_emb_inc"`
	TotEmbExc     float64 `json:"tot_emb_exc"`
	GrType        int64   `json:"gr_type"`
	SubTotal      float64 `json:"sub_total"`
	Vat           float64 `json:"vat"`
	VatValue      float64 `json:"vat_value"`
	VatLg         float64 `json:"vat_lg"`
	VatLgValue    float64 `json:"vat_lg_value"`
	Total         float64 `json:"total"`
	VatBg         float64 `json:"vat_bg"`
	VatBgValue    float64 `json:"vat_bg_value"`
	WeekIDOb      int     `json:"week_id_ob"`
	DataStatus    int64   `json:"data_status"`
	CreatedBy     int64   `json:"created_by"`
	UpdatedByName string  `json:"updated_by_name"`
	Details       []GrDet `json:"details"`
}

type StockList struct {
	StockID   int64   `json:"stock_id"`
	StockDate string  `json:"stock_date"`
	TrCode    string  `json:"tr_code"`
	TrNo      string  `json:"tr_no"`
	WhID      int64   `json:"wh_id"`
	ProID     int64   `json:"pro_id"`
	ItemCdn   int64   `json:"item_cdn"`
	QtyIn     float64 `json:"qty_in"`
	QtyOut    float64 `json:"qty_out"`
	UnitPrice float64 `json:"unit_price"`
	Cogs      float64 `json:"cogs"`
	RefDetId  int64   `json:"ref_det_id"`
}

type CreateStock struct {
	CustID       string  `json:"cust_id" validate:"required"`
	ParentCustID string  `json:"parent_cust_id" validate:"required"`
	StockDate    string  `json:"stock_date" validate:"required"`
	TrCode       string  `json:"tr_code" validate:"required,max=10,alphanum"`
	TrNo         string  `json:"tr_no" validate:"required,max=30,alphanum"`
	WhID         int64   `json:"wh_id" validate:"min=0"`
	ProID        int64   `json:"pro_id" validate:"required"`
	ItemCdn      int64   `json:"item_cdn" validate:"required,oneof=1 2"`
	QtyIn        float64 `json:"qty_in" validate:"min=0"`
	QtyOut       float64 `json:"qty_out" validate:"min=0"`
	UnitPrice    float64 `json:"unit_price" validate:"min=0"`
	Cogs         float64 `json:"cogs" validate:"min=0"`
	RefDetId     int64   `json:"ref_det_id" validate:"required"`
}

type CreateBulkStock struct {
	Stock []CreateStock `json:"stock" validate:"required,dive"`
}

type StockUpdate struct {
	CustID    string    `json:"cust_id"`
	WhID      int64     `json:"wh_id"`
	ProID     int64     `json:"pro_id"`
	StockDate time.Time `json:"stock_date"`
	TrCode    string    `json:"tr_code"`
	TrNo      string    `json:"tr_no"`
	ItemCdn   int64     `json:"item_cdn"`
	QtyIn     float64   `json:"qty_in"`
	QtyOut    float64   `json:"qty_out"`
	UnitPrice float64   `json:"unit_price"`
	RefDetId  int64     `json:"ref_det_id"`
}

type SalesOrderStockUpdate struct {
	CustID         string    `json:"cust_id"`
	WhID           int64     `json:"wh_id"`
	ProID          int64     `json:"pro_id"`
	StockDate      time.Time `json:"stock_date"`
	TrCode         string    `json:"tr_code"`
	TrNo           string    `json:"tr_no"`
	ItemCdn        int64     `json:"item_cdn"`
	QtyOrder       float64   `json:"qty_order"`
	QtyOrderBefore *float64  `json:"qty_order_before"`
	UnitPrice      float64   `json:"unit_price"`
	RefDetId       int64     `json:"ref_det_id"`
}

type InvoiceSalesStockUpdate struct {
	CustID         string    `json:"cust_id"`
	WhID           int64     `json:"wh_id"`
	ProID          int64     `json:"pro_id"`
	StockDate      time.Time `json:"stock_date"`
	TrCode         string    `json:"tr_code"`
	TrNo           string    `json:"tr_no"`
	ItemCdn        int64     `json:"item_cdn"`
	QtyOrderBefore float64   `json:"qty_order_before"`
	UnitPrice      float64   `json:"unit_price"`
	RefDetId       int64     `json:"ref_det_id"`
}

type StockReportQueryFilter struct {
	CustID            string  `query:"cust_id" json:"cust_id,omitempty"`
	ParentCustID      string  `query:"parent_cust_id" json:"parent_cust_id,omitempty"`
	Date              string  `query:"date" json:"date" validate:"omitempty,len=10,yyyyMmDdDate"`
	OrderDate         string  `query:"order_date" json:"order_date" validate:"omitempty,len=10,yyyyMmDdDate"`
	Page              int     `query:"page" json:"page"`
	Limit             int     `query:"limit" json:"limit" validate:""`
	Query             string  `query:"q" json:"q"`
	Sort              string  `query:"sort" json:"sort" validate:"oneof='pro_id:asc' 'pro_id:desc' 'pro_code:asc' 'pro_code:desc' 'pro_name:asc' 'pro_name:desc'"`
	ProID             []int64 `query:"pro_id" json:"pro_id"`
	WhID              []int   `query:"wh_id" json:"wh_id" validate:"required,min=1"`
	SupID             []int   `query:"sup_id" json:"sup_id" validate:""`
	ShowPrice         string  `query:"show_price" json:"show_price" validate:"omitempty,oneof='true' 'false'"`
	IncludeZeroStock  string  `query:"include_zero_stock" json:"include_zero_stock" validate:"omitempty,oneof='true' 'false'"`
	ActiveProductOnly string  `query:"active_product_only" json:"active_product_only" validate:"omitempty,oneof='true' 'false'"`
	BrandID           []int   `query:"brand_id" json:"brand_id" validate:""`
	PCatID            []int   `query:"pcat_id" json:"pcat_id" validate:""`
	PLID              []int   `query:"pl_id" json:"pl_id" validate:""`
	OutletID          int64   `query:"outlet_id" json:"outlet_id" validate:"omitempty"`
}

type StockReport struct {
	ProID              int64   `json:"pro_id"`
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

type StockOpnameLookupQueryFilter struct {
	CustID             string `query:"cust_id" json:"cust_id,omitempty"`
	ParentCustID       string `query:"parent_cust_id" json:"parent_cust_id,omitempty"`
	Date               string `query:"date" json:"date" validate:"omitempty,len=10,yyyyMmDdDate"`
	WhID               []int  `query:"wh_id" json:"wh_id" validate:"required,min=1"`
	ProductHierarchy   int    `query:"product_hierarchy" json:"product_hierarchy" validate:"required,min=1,max=4"`
	IncludeZeroStock   string `query:"include_zero_stock" json:"include_zero_stock" validate:"omitempty,oneof='true' 'false'"`
	IsShowCurrentStock string `query:"is_show_current_stock" json:"is_show_current_stock" validate:"omitempty,oneof='true' 'false'"`
}

type StockOpnameLookup struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
