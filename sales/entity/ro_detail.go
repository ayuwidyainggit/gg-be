package entity

type CreateRoDetBody struct {
	SeqNo       int      `json:"seq_no"`
	ProId       int      `json:"pro_id"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price1"`
	PurchPrice2 *float64 `json:"purch_price2"`
	PurchPrice3 *float64 `json:"purch_price3"`
	SellPrice1  *float64 `json:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}

type RoDetResponse struct {
	RoDetId     int64    `json:"ro_det_id"`
	SeqNo       int      `json:"seq_no"`
	ProId       int      `json:"pro_id"`
	ProCode     string   `json:"pro_code"`
	ProName     string   `json:"pro_name"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price1"`
	PurchPrice2 *float64 `json:"purch_price2"`
	PurchPrice3 *float64 `json:"purch_price3"`
	SellPrice1  *float64 `json:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}

type UpdateRoDetBody struct {
	SeqNo       int      `json:"seq_no"`
	RoDetId     *int64   `json:"ro_det_id"`
	ProId       int      `json:"pro_id"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price1"`
	PurchPrice2 *float64 `json:"purch_price2"`
	PurchPrice3 *float64 `json:"purch_price3"`
	SellPrice1  *float64 `json:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}

type RoDetWithGroup struct {
	Normal []CreateRoDetBody `json:"normal"`
	Promo  []CreateRoDetBody `json:"promo"`
}
type RoDetReadWithGroup struct {
	Normal []RoDetResponse `json:"normal"`
	Promo  []RoDetResponse `json:"promo"`
}
type UpdateRoDetWithGroup struct {
	Normal []UpdateRoDetBody `json:"normal"`
	Promo  []UpdateRoDetBody `json:"promo"`
}
