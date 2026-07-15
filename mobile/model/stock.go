package model

type StockWareHouseList struct {
	CustID    string  `gorm:"cust_id" json:"cust_id"`
	EmpId     *int64  `gorm:"emp_id" json:"emp_id"`
	SalesName *string `gorm:"sales_name" json:"sales_name"`
	WhId      *int64  `gorm:"wh_id" json:"wh_id"`
	WhCode    *string `gorm:"wh_code" json:"wh_code"`
	WhName    *string `gorm:"wh_name" json:"wh_name"`
}

func (StockWareHouseList) TableName() string {
	return "mst.m_salesman"
}

type DetilProductStock struct {
	// CustId      string  `json:"cust_id" gorm:"column:cust_id"`
	// WhId        *int64  `json:"wh_id" gorm:"column:wh_id"`
	// ProductId   *int64  `json:"pro_id" gorm:"column:pro_id"`
	// ProductCode *string `json:"pro_code" gorm:"column:pro_code"`
	// ProductName *string `json:"pro_name" gorm:"column:pro_name"`
	// UnitId1     string  `json:"unit_id1" gorm:"column:unit_id1"`
	// UnitId2     string  `json:"unit_id2" gorm:"column:unit_id2"`
	// UnitId3     string  `json:"unit_id3" gorm:"column:unit_id3"`
	// UnitId4     *string `json:"unit_id4" gorm:"column:unit_id4"`
	// UnitId5     *string `json:"unit_id5" gorm:"column:unit_id5"`
	// ConvUnit2   float32 `json:"conv_unit2" gorm:"column:conv_unit2"`
	// ConvUnit3   float32 `json:"conv_unit3" gorm:"column:conv_unit3"`
	// ConvUnit4   float32 `json:"conv_unit4" gorm:"column:conv_unit4"`
	// ConvUnit5   float32 `json:"conv_unit5" gorm:"column:conv_unit5"`
	// PurchPrice1 float64 `json:"purch_price1" gorm:"column:purch_price1"`
	// PurchPrice2 float64 `json:"purch_price2" gorm:"column:purch_price2"`
	// PurchPrice3 float64 `json:"purch_price3" gorm:"column:purch_price3"`
	// PurchPrice4 float64 `json:"purch_price4" gorm:"column:purch_price4"`
	// PurchPrice5 float64 `json:"purch_price5" gorm:"column:purch_price5"`
	// SellPrice1  float64 `json:"sell_price1" gorm:"column:sell_price1"`
	// SellPrice2  float64 `json:"sell_price2" gorm:"column:sell_price2"`
	// SellPrice3  float64 `json:"sell_price3" gorm:"column:sell_price3"`
	// SellPrice4  float64 `json:"sell_price4" gorm:"column:sell_price4"`
	// SellPrice5  float64 `json:"sell_price5" gorm:"column:sell_price5"`

	WhId        *int64  `json:"wh_id" gorm:"column:wh_id"`
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	UnitId1     string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     string  `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2   int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   int     `gorm:"column:conv_unit3" json:"conv_unit3"`
	PurchPrice1 float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2 float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3 float64 `gorm:"column:purch_price3" json:"purch_price3"`
	SellPrice1  float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2  float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3  float64 `gorm:"column:sell_price3" json:"sell_price3"`
	Qty         float64 `gorm:"column:qty" json:"qty"`
	QtyStock    float64 `gorm:"column:qty_stock" json:"qty_stock"`
	QtyOrder    float64 `gorm:"column:qty_order" json:"qty_order"`
	QtyOrder1   float64 `gorm:"column:qty_order1" json:"qty_order1"`
	QtyOrder2   float64 `gorm:"column:qty_order2" json:"qty_order2"`
	QtyOrder3   float64 `gorm:"column:qty_order3" json:"qty_order3"`
	IsActive    bool    `gorm:"column:is_active" json:"is_active"`
	Vat         float64 `gorm:"column:vat" json:"vat"`
	VatLgPurch  float64 `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	VatLgSell   float64 `gorm:"column:vat_lg_sell" json:"vat_lg_sell"`
}

func (DetilProductStock) TableName() string {
	return "mst.m_product"
}

type Stock struct {
	CustID      string  `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	StockID     int64   `gorm:"column:stock_id;primaryKey;autoIncrement" json:"stock_id"`
	StockDate   string  `gorm:"column:stock_date" json:"stock_date"`
	TrCode      string  `gorm:"column:tr_code" json:"tr_code"`
	TrNo        string  `gorm:"column:tr_no" json:"tr_no"`
	WhID        int64   `gorm:"column:wh_id" json:"wh_id"`
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ItemCdn     int16   `gorm:"column:item_cdn" json:"item_cdn"`
	QtyIn       float64 `gorm:"column:qty_in" json:"qty_in"`
	QtyOut      float64 `gorm:"column:qty_out" json:"qty_out"`
	UnitPrice   float64 `gorm:"column:unit_price" json:"unit_price"`
	Cogs        float64 `gorm:"column:cogs" json:"cogs"`
	RefDetID    int64   `gorm:"column:ref_det_id" json:"ref_det_id"`
	CreatedAt   int64   `gorm:"column:created_at" json:"created_at"`
	QtyInOrder  float64 `gorm:"column:qty_in_order" json:"qty_in_order"`
	QtyOutOrder float64 `gorm:"column:qty_out_order" json:"qty_out_order"`
}

func (Stock) TableName() string {
	return "inv.stock"
}
