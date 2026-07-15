package entity

type CreateWhAdjDetBody struct {
	WhAdjDetId   *int    `json:"wh_adj_det_id"`
	StockAdjNo   string  `json:"stock_adjustment_no"`
	SeqNo        int     `json:"seq_no"`
	ProID        int64   `json:"pro_id"`
	Qty1         float64 `json:"qty1"`
	Qty2         float64 `json:"qty2"`
	Qty3         float64 `json:"qty3"`
	WhAdjDetType int     `json:"wh_adj_det_type"`
}

type WhAdjDetresponse struct {
	WhAdjDetId   *int    `json:"wh_adj_det_id"`
	SeqNo        int     `json:"seq_no"`
	ProID        int     `json:"pro_id"`
	ProCode      string  `json:"pro_code"`
	ProName      string  `json:"pro_name"`
	Qty1         int     `json:"qty1"`
	Qty2         int     `json:"qty2"`
	Qty3         int     `json:"qty3"`
	UnitId1      *string `json:"unit_id1"`
	UnitId2      *string `json:"unit_id2"`
	UnitId3      *string `json:"unit_id3"`
	WhAdjDetType int     `json:"wh_adj_det_type"`
}
type UpdateWhAdjDetBody struct {
	WhAdjDetId *int     `json:"wh_adj_det_id"`
	CustID     string   `json:"cust_id"`
	AdjNo      string   `json:"adj_no"`
	SeqNo      *int     `json:"seq_no"`
	ProID      *int     `json:"pro_id"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	BatchNo    *string  `json:"batch_no"`
	ExpDate    *string  `json:"exp_date"`
}
