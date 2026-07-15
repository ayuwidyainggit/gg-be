package entity

import "time"

type GrBranchDet struct {
	ID         int      `json:"gr_branch_det_id"`
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	ItemType   int      `json:"item_type"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	UnitPrice1 *float64 `json:"unit_price1"`
	UnitPrice2 *float64 `json:"unit_price2"`
	UnitPrice3 *float64 `json:"unit_price3"`
	UnitPrice4 *float64 `json:"unit_price4"`
	UnitPrice5 *float64 `json:"unit_price5"`
	UnitId1    *string  `json:"unit_id1"`
	UnitId2    *string  `json:"unit_id2"`
	UnitId3    *string  `json:"unit_id3"`
	UnitId4    *string  `json:"unit_id4"`
	UnitId5    *string  `json:"unit_id5"`
	EmbInc     *float64 `json:"emb_inc"`
	EmbExc     *float64 `json:"emb_exc"`
	InvoiceNo  *string  `json:"invoice_no"`
	BatchNo    *string  `json:"batch_no"`
	ExpDate    string   `json:"exp_date"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	ConvUnit4  float64  `json:"conv_unit4"`
	ConvUnit5  float64  `json:"conv_unit5"`
}

type GrBranchDetailRequest struct {
	GrBranchDetId *int     `json:"gr_branch_det_id"`
	CustID        string   `json:"cust_id"`
	SeqNo         int      `json:"seq_no"`
	ProID         int64    `json:"pro_id" validate:"required"`
	Unit1         string   `json:"unit_id1"`
	Unit2         string   `json:"unit_id2"`
	Unit3         string   `json:"unit_id3"`
	UnitPrice1    float64  `json:"unit_price1"`
	UnitPrice2    float64  `json:"unit_price2"`
	UnitPrice3    float64  `json:"unit_price3"`
	Vat           *float64 `json:"vat"`
	VatLgPurch    *float64 `json:"vat_lg_purch"`
	VatBg         *float64 `json:"vat_bg"`
	QtyShip1      int      `json:"qty_ship1"`
	QtyShip2      int      `json:"qty_ship2"`
	QtyShip3      int      `json:"qty_ship3"`
	QtyShip       int      `json:"qty_ship"`
	QtyReceived1  int      `json:"qty_received1"`
	QtyReceived2  int      `json:"qty_received2"`
	QtyReceived3  int      `json:"qty_received3"`
	QtyReceived   int      `json:"qty_received"`
	ConvUnit2     int      `json:"conv_unit2"`
	ConvUnit3     int      `json:"conv_unit3"`
}

type GrBranchDetUpdateRequest struct {
	GrBranchDetId *int       `json:"gr_branch_det_id"`
	SeqNo         int        `json:"seq_no"`
	CustID        string     `json:"cust_id"`
	ProID         int        `json:"pro_id"`
	ItemType      int        `json:"item_type"`
	Qty           *float64   `json:"qty"`
	QtyStr        *string    `json:"qty_str"`
	UnitPrice1    *float64   `json:"unit_price1"`
	UnitPrice2    *float64   `json:"unit_price2"`
	UnitPrice3    *float64   `json:"unit_price3"`
	UnitPrice4    *float64   `json:"unit_price4"`
	UnitPrice5    *float64   `json:"unit_price5"`
	UnitId1       *string    `json:"unit_id1"`
	UnitId2       *string    `json:"unit_id2"`
	UnitId3       *string    `json:"unit_id3"`
	UnitId4       *string    `json:"unit_id4"`
	UnitId5       *string    `json:"unit_id5"`
	EmbInc        *float64   `json:"emb_inc"`
	EmbExc        *float64   `json:"emb_exc"`
	InvoiceNo     *string    `json:"invoice_no"`
	BatchNo       *string    `json:"batch_no"`
	ExpDate       *time.Time `json:"exp_date"`
	Vat           *float64   `json:"vat"`
	VatBg         *float64   `json:"vat_bg"`
	VatVgPurch    *float64   `json:"vat_lg_purch"`
	ExciseRate    *float64   `json:"excise_rate"`
	ExciseTax     *float64   `json:"excise_tax"`
	ConvUnit2     float64    `json:"conv_unit2"`
	ConvUnit3     float64    `json:"conv_unit3"`
	ConvUnit4     float64    `json:"conv_unit4"`
	ConvUnit5     float64    `json:"conv_unit5"`
	Qty1          *float64   `json:"qty1"`
	Qty2          *float64   `json:"qty2"`
	Qty3          *float64   `json:"qty3"`
	QtyShip1      *float64   `json:"qty_ship1"`
	QtyShip2      *float64   `json:"qty_ship2"`
	QtyShip3      *float64   `json:"qty_ship3"`
}
type GrBranchDetListGroup struct {
	Normal []GrBranchDetailList `json:"normal"`
	Promo  []GrBranchDetailList `json:"promo"`
}

type GrBranchDetWithGroup struct {
	Normal []GrBranchDetailRequest `json:"normal" validate:"required,dive"`
	Promo  []GrBranchDetailRequest `json:"promo"`
}

type GrBranchDetailList struct {
	GrBranchDetId int    `json:"gr_branch_det_id"`
	SeqNo         int    `json:"seq_no"`
	ProID         int    `json:"pro_id"`
	ProCode       string `json:"pro_code"`
	ProName       string `json:"pro_name"`
	// Qty1            int      `json:"qty1"`
	// Qty2            int      `json:"qty2"`
	// Qty3            int      `json:"qty3"`
	// Qty             float64  `json:"qty"`
	UnitId1     *string  `json:"unit_id1"`
	UnitId2     *string  `json:"unit_id2"`
	UnitId3     *string  `json:"unit_id3"`
	PurchPrice1 *float64 `json:"purch_price1" validate:"required"`
	Purchrice2  *float64 `json:"purch_price2" validate:"required"`
	PurchPrice3 *float64 `json:"purch_price3" validate:"required"`
	// Vat             float64  `json:"vat"`
	// VatValue        float64  `json:"vat_value"`
	// VatLgPurch      float64  `json:"vat_lg_purch"`
	// VatLgPurchValue float64  `json:"vat_lg_purch_value"`
	// VatBg           *float64 `json:"vat_bg"`
	QtyShip      float64 `json:"qty_ship"`
	QtyShip1     int     `json:"qty_ship1"`
	QtyShip2     int     `json:"qty_ship2"`
	QtyShip3     int     `json:"qty_ship3"`
	QtyReceived  float64 `json:"qty_received"`
	QtyReceived1 int     `json:"qty_received1"`
	QtyReceived2 int     `json:"qty_received2"`
	QtyReceived3 int     `json:"qty_received3"`
	// SubTotal        float64  `json:"sub_total"`
	// Nett            float64  `json:"nett"`
	// Total           float64  `json:"total"`
	// ConvUnit1 float64 `json:"conv_unit1"`
	ConvUnit2 float64 `json:"conv_unit2"`
	ConvUnit3 float64 `json:"conv_unit3"`
	// QtyRemaining1   int      `json:"qty_remaining1"`
	// QtyRemaining2   int      `json:"qty_remaining2"`
	// QtyRemaining3   int      `json:"qty_remaining3"`
	Qty1Alloc *float64 `json:"qty1_alloc"`
	Qty2Alloc *float64 `json:"qty2_alloc"`
	Qty3Alloc *float64 `json:"qty3_alloc"`
	// WhQty1          int      `json:"wh_qty1"`
	// WhQty2          int      `json:"wh_qty2"`
	// WhQty3          int      `json:"wh_qty3"`
	Amount   *float64 `json:"amount"`
	VatValue *float64 `json:"vat_value"`
}
