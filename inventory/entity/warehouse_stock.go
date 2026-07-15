package entity

type DistributorStockQueryFilter struct {
	CustID        string  `query:"cust_id" json:"cust_id,omitempty"`
	ParentCustID  string  `query:"parent_cust_id" json:"parent_cust_id,omitempty"`
	Page          int     `query:"page" json:"page"`
	Limit         int     `query:"limit" json:"limit" validate:"required"`
	Query         string  `query:"q" json:"q"`
	Sort          string  `query:"sort" json:"sort" validate:"oneof='pro_id:asc' 'pro_id:desc' 'pro_code:asc' 'pro_code:desc' 'pro_name:asc' 'pro_name:desc'"`
	ProID         []int64 `query:"pro_id" json:"pro_id"`
	WhID          int     `query:"wh_id" json:"wh_id" validate:"required"`
	SupID         []int   `query:"sup_id" json:"sup_id" validate:""`
	ShowPrice     string  `query:"show_price" json:"show_price" validate:"omitempty,oneof='true' 'false'"`
	ZeroStock     string  `query:"zero_stock" json:"zero_stock" validate:"omitempty,oneof='true' 'false'"`
	ActiveProduct string  `query:"active_product" json:"active_product" validate:"omitempty,oneof='true' 'false'"`
}

type WarehouseStockWhListQueryFilter struct {
	CustID       string
	ParentCustID string
	From         *int64  `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64  `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit" validate:"required"`
	Query        string  `query:"q"`
	Sort         string  `query:"sort"`
	ProID        []int64 `query:"pro_id"`
	WhID         []int   `query:"wh_id"`
	SupID        []int   `query:"sup_id"`
}

type ProductWarehouseListQueryFilter struct {
	CustID        string
	ParentCustID  string
	DistributorID int64
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	WhID          int    `query:"wh_id"`
	ProID         int64  `query:"pro_id"`
}

type WarehouseStock struct {
	CustID        string  `json:"cust_id"`
	ParentCustID  string  `json:"parent_cust_id"`
	WhId          int64   `json:"wh_id"`
	ProId         int64   `json:"pro_id"`
	Qty           float64 `json:"qty"`
	QtyOnOrder    float64 `json:"qty_on_order"`
	QtyOnShipping float64 `json:"qty_on_shipping"`
	QtyBs         float64 `json:"qty_bs"`
	QtyExp        float64 `json:"qty_exp"`
}

type UpsertWarehouseStock struct {
	CustID        string  `json:"cust_id" validate:"required"`
	ParentCustID  string  `json:"parent_cust_id" validate:"required"`
	WhId          int64   `json:"wh_id" validate:""`
	ProId         int64   `json:"pro_id" validate:"required"`
	Qty           float64 `json:"qty" validate:"min=0"`
	QtyOnOrder    float64 `json:"qty_on_order" validate:"min=0"`
	QtyOnShipping float64 `json:"qty_on_shipping" validate:"min=0"`
	QtyBs         float64 `json:"qty_bs" validate:"min=0"`
	QtyExp        float64 `json:"qty_exp" validate:"min=0"`
}

type UpsertBulkWarehouseStock struct {
	WarehouseStock []UpsertWarehouseStock `json:"warehouse_stock" validate:"required"`
}

type DistributorStockList struct {
	ProID         int64   `json:"pro_id"`
	ProCode       string  `json:"pro_code"`
	ProName       string  `json:"pro_name"`
	UnitId1       string  `json:"unit_id1"`
	UnitId2       string  `json:"unit_id2"`
	UnitId3       string  `json:"unit_id3"`
	PurchPrice1   float64 `json:"purch_price1"`
	PurchPrice2   float64 `json:"purch_price2"`
	PurchPrice3   float64 `json:"purch_price3"`
	SellPrice1    float64 `json:"sell_price1"`
	SellPrice2    float64 `json:"sell_price2"`
	SellPrice3    float64 `json:"sell_price3"`
	SupID         int64   `json:"sup_id"`
	SupCode       string  `json:"sup_code"`
	SupName       string  `json:"sup_name"`
	TotalQty      float64 `json:"total_qty"`
	Qty1          float64 `json:"qty1"`
	Qty2          float64 `json:"qty2"`
	Qty3          float64 `json:"qty3"`
	ConvUnit2     float64 `json:"conv_unit2"`
	ConvUnit3     float64 `json:"conv_unit3"`
	QtyOnOrder    float64 `json:"qty_on_order"`
	QtyOnShipping float64 `json:"qty_on_shipping"`
	QtyBs         float64 `json:"qty_bs"`
	QtyExp        float64 `json:"qty_exp"`
	UpdatedAt     int64   `json:"updated_at"`
}

type ProductWarehouseList struct {
	ProID       int64   `json:"pro_id"`
	ProCode     string  `json:"pro_code"`
	ProName     string  `json:"pro_name"`
	Qty1        float64 `json:"qty1"`
	Qty2        float64 `json:"qty2"`
	Qty3        float64 `json:"qty3"`
	SellPrice1  float64 `json:"sell_price1"`
	SellPrice2  float64 `json:"sell_price2"`
	SellPrice3  float64 `json:"sell_price3"`
	UnitId1     string  `json:"unit_id1"`
	UnitId2     string  `json:"unit_id2"`
	UnitId3     string  `json:"unit_id3"`
	Vat         float64 `json:"vat"`
	VatLgPurch  float64 `json:"vat_lg_purch"`
	VatLgSell   float64 `json:"vat_lg_sell"`
	VatBg       float64 `json:"vat_bg"`
	PurchPrice1 float64 `json:"purch_price1"`
	PurchPrice2 float64 `json:"purch_price2"`
	PurchPrice3 float64 `json:"purch_price3"`
}

type WarehouseStockWhList struct {
	WhID      int64  `json:"wh_id"`
	WhCode    string `json:"wh_code"`
	WhName    string `json:"wh_name"`
	StockType string `json:"stock_type"`
}
