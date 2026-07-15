package entity

import "inventory/pkg/conversion"

type CreateSupplierReturnDetBody struct {
	SeqNo          int     `json:"seq_no"`
	ProID          int64   `json:"pro_id" validate:"required"`
	Qty1           float64 `json:"qty1"`
	Qty2           float64 `json:"qty2"`
	Qty3           float64 `json:"qty3"`
	ItemCdn        *int64  `json:"item_cdn" validate:"required"`
	ReturnReasonID *int64  `json:"return_reason_id" validate:"required"`
	Qty            float64
	UnitPrice1     float64
	UnitPrice2     float64
	UnitPrice3     float64
	ConvUnit2      float64
	ConvUnit3      float64
	Vat            float64
	VatValue       float64
	VatLg          float64
	VatLgValue     float64
	VatBg          float64
	VatBgValue     float64
	Disc           float64
	DiscValue      float64
	Subtotal       float64
	Total          float64
	QtyRemaining   float64
}

func (r *CreateSupplierReturnDetBody) Calculate() error {
	QtyUnitGrDet := &conversion.QtyUnit{
		Qty1:      int(r.Qty1),
		Qty2:      int(r.Qty2),
		Qty3:      int(r.Qty3),
		ConvUnit2: int(r.ConvUnit2),
		ConvUnit3: int(r.ConvUnit3),
	}

	totalQtyGrDet, err := QtyUnitGrDet.ToTotalQuantity()
	if err != nil {
		return err
	}

	r.Qty = float64(totalQtyGrDet)
	r.Subtotal = (r.UnitPrice1 * r.Qty1) + (r.UnitPrice2 * r.Qty2) + (r.UnitPrice3 * r.Qty3)
	r.DiscValue = r.Subtotal * (r.Disc / 100)
	subtotalAfterDiscount := r.Subtotal - r.DiscValue
	r.VatValue = subtotalAfterDiscount * (r.Vat / 100)
	r.VatLgValue = subtotalAfterDiscount * (r.VatLg / 100)
	r.VatBgValue = subtotalAfterDiscount * (r.VatBg / 100)
	r.Total = subtotalAfterDiscount + r.VatValue + r.VatLgValue + r.VatBgValue

	return nil
}

type SupplierReturnGetDetResp struct {
	SeqNo            int     `json:"seq_no"`
	ProID            int64   `json:"pro_id" validate:"required"`
	ProCode          string  `json:"pro_code"`
	ProName          string  `json:"pro_name"`
	Qty1             int     `json:"qty1" validate:"required"`
	Qty2             int     `json:"qty2" validate:"required"`
	Qty3             int     `json:"qty3" validate:"required"`
	UnitPrice1       float64 `json:"unit_price1"`
	UnitPrice2       float64 `json:"unit_price2"`
	UnitPrice3       float64 `json:"unit_price3"`
	UnitId1          string  `json:"unit_id1"`
	UnitId2          string  `json:"unit_id2"`
	UnitId3          string  `json:"unit_id3"`
	ConvUnit1        float64 `json:"conv_unit1"`
	ConvUnit2        float64 `json:"conv_unit2"`
	ConvUnit3        float64 `json:"conv_unit3"`
	InvoiceQty1      int     `json:"invoice_qty1"`
	InvoiceQty2      int     `json:"invoice_qty2"`
	InvoiceQty3      int     `json:"invoice_qty3"`
	RemainingQty1    int     `json:"remaining_qty1"`
	RemainingQty2    int     `json:"remaining_qty2"`
	RemainingQty3    int     `json:"remaining_qty3"`
	WhQty1           int     `json:"wh_qty1"`
	WhQty2           int     `json:"wh_qty2"`
	WhQty3           int     `json:"wh_qty3"`
	ItemCdn          *int64  `json:"item_cdn" validate:"required"`
	ReturnReasonID   int64   `json:"return_reason_id" validate:"required"`
	ReturnReasonName *string `json:"return_reason_name"`
	Vat              float64 `json:"vat"`
	VatValue         float64 `json:"vat_value"`
	VatLg            float64 `json:"vat_lg"`
	VatLgValue       float64 `json:"vat_lg_value"`
	VatBg            float64 `json:"vat_bg"`
	VatBgValue       float64 `json:"vat_bg_value"`
	Discount         float64 `json:"discount"`
	DiscountValue    float64 `json:"discount_value"`
	SubTotal         float64 `json:"sub_total"`
	Nett             float64 `json:"nett"`
	Total            float64 `json:"total"`
}
