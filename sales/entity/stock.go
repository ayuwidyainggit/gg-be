package entity

import "time"

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

type CancelStockBasis struct {
	CustID            string    `json:"cust_id" gorm:"column:cust_id"`
	WhID              int64     `json:"wh_id" gorm:"column:wh_id"`
	ProID             int64     `json:"pro_id" gorm:"column:pro_id"`
	RefDetID          int64     `json:"ref_det_id" gorm:"column:ref_det_id"`
	StockDate         time.Time `json:"stock_date" gorm:"column:stock_date"`
	UnitPrice         float64   `json:"unit_price" gorm:"column:unit_price"`
	QtyFinal          float64   `json:"qty_final" gorm:"column:qty_final"`
	QtyOutSO          float64   `json:"qty_out_so" gorm:"column:qty_out_so"`
	QtyOutOrderCancel float64   `json:"qty_out_order_cancel" gorm:"column:qty_out_order_cancel"`
	QtyOutstanding    float64   `json:"qty_outstanding" gorm:"column:qty_outstanding"`
	QtyOutSmallest    float64   `json:"qty_out_smallest" gorm:"column:qty_out_smallest"`
	IsMissingSource   bool      `json:"is_missing_source" gorm:"column:is_missing_source"`
	IsAmbiguous       bool      `json:"is_ambiguous" gorm:"column:is_ambiguous"`
}

type CancelStockWrite struct {
	CustID      string  `json:"cust_id"`
	RoNo        string  `json:"ro_no"`
	WhID        int64   `json:"wh_id"`
	ProID       int64   `json:"pro_id"`
	RefDetID    int64   `json:"ref_det_id"`
	QtySmallest float64 `json:"qty_smallest"`
	UnitPrice   float64 `json:"unit_price"`
}
