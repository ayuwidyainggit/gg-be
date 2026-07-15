package conversion

import "errors"

type QtyUnit struct {
	Qty1      int
	Qty2      int
	Qty3      int
	ConvUnit2 int
	ConvUnit3 int
}

func (q *QtyUnit) Get() *QtyUnit {
	return q
}

func (q *QtyUnit) DoConversion() {
	rQty2 := q.Qty1 / q.ConvUnit2
	if rQty2 > 0 {
		q.Qty1 = q.Qty1 % q.ConvUnit2
		q.Qty2 += rQty2
	}

	rQty3 := q.Qty2 / q.ConvUnit3
	if rQty3 > 0 {
		q.Qty2 = q.Qty2 % q.ConvUnit3
		q.Qty3 += rQty3
	}
}

func (q *QtyUnit) ToTotalQuantity() (int, error) {
	if q.ConvUnit2 == 0 || q.ConvUnit3 == 0 {
		return 0, errors.New("invalid conv unit")
	}

	q.DoConversion()

	totalQty := (q.ConvUnit2*q.ConvUnit3)*q.Qty3 + q.ConvUnit2*q.Qty2 + q.Qty1
	return totalQty, nil
}
