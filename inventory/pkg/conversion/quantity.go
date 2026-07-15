package conversion

type Qty struct {
	Qty       int
	ConvUnit2 int
	ConvUnit3 int
}

type QtyConversionResult struct {
	Qty1 int
	Qty2 int
	Qty3 int
}

func (q *Qty) ConvToQtyConversion() (qtyConv QtyConversionResult) {
	qtyConv.Qty3 = q.Qty / (q.ConvUnit2 * q.ConvUnit3)
	q.Qty -= qtyConv.Qty3 * (q.ConvUnit2 * q.ConvUnit3)

	qtyConv.Qty2 = q.Qty / q.ConvUnit2
	q.Qty -= qtyConv.Qty2 * q.ConvUnit2

	qtyConv.Qty1 = q.Qty
	return
}
