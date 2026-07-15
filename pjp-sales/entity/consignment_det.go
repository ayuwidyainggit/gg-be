package entity

type CreateConsignDetBody struct {
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	SellPrice1 *float64 `json:"sell_price_1"`
	SellPrice2 *float64 `json:"sell_price_2"`
	SellPrice3 *float64 `json:"sell_price_3"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	TotAmount  *float64 `json:"tot_amount"`
	BatchNo    *string  `json:"batch_no"`
	ExpDate    *string  `json:"exp_date"`
}

type ConsignDetResponse struct {
	ConsDetID  int      `json:"cons_det_id"`
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	ProCode    string   `json:"pro_code"`
	ProName    string   `json:"pro_name"`
	SellPrice1 *float64 `json:"sell_price_1"`
	SellPrice2 *float64 `json:"sell_price_2"`
	SellPrice3 *float64 `json:"sell_price_3"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	TotAmount  *float64 `json:"tot_amount"`
	BatchNo    *string  `json:"batch_no"`
	ExpDate    *string  `json:"exp_date"`
}
type UpdateConsignDetBody struct {
	ConsDetID  *int64   `json:"cons_det_id"`
	SeqNo      int      `json:"seq_no"`
	ProID      int      `json:"pro_id"`
	SellPrice1 *float64 `json:"sell_price_1"`
	SellPrice2 *float64 `json:"sell_price_2"`
	SellPrice3 *float64 `json:"sell_price_3"`
	Qty        *float64 `json:"qty"`
	QtyStr     *string  `json:"qty_str"`
	TotAmount  *float64 `json:"tot_amount"`
	BatchNo    *string  `json:"batch_no"`
	ExpDate    *string  `json:"exp_date"`
}
