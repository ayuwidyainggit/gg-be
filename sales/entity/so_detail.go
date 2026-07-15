package entity

type CreateSoDetBody struct {
	SoNo        string   `json:"so_no"`
	SeqNo       int      `json:"seq_no"`
	ProID       int      `json:"pro_id"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price_1"`
	PurchPrice2 *float64 `json:"purch_price_2"`
	PurchPrice3 *float64 `json:"purch_price_3"`
	SellPrice1  *float64 `json:"sell_price_1"`
	SellPrice2  *float64 `json:"sell_price_2"`
	SellPrice3  *float64 `json:"sell_price_3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}

type CreateSoDetBodyWithGroup struct {
	Normal []CreateSoDetBody `json:"normal"`
	Promo  []CreateSoDetBody `json:"promo"`
}

type UpdateSoDetBodyWithGroup struct {
	Normal []UpdateSoDetBody `json:"normal"`
	Promo  []UpdateSoDetBody `json:"promo"`
}

type SoDetResponseWithGroup struct {
	Normal []SoDetResponse `json:"normal"`
	Promo  []SoDetResponse `json:"promo"`
}
type UpdateSoDetBody struct {
	SeqNo       int      `json:"seq_no"`
	SoDetID     *int64   `json:"so_det_id"`
	ProID       int      `json:"pro_id"`
	ProCode     string   `json:"pro_code"`
	ProName     string   `json:"pro_name"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price_1"`
	PurchPrice2 *float64 `json:"purch_price_2"`
	PurchPrice3 *float64 `json:"purch_price_3"`
	SellPrice1  *float64 `json:"sell_price_1"`
	SellPrice2  *float64 `json:"sell_price_2"`
	SellPrice3  *float64 `json:"sell_price_3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}
type SoDetResponse struct {
	SoNo        string   `json:"so_no"`
	SeqNo       int      `json:"seq_no"`
	SoDetID     int      `json:"so_det_id"`
	ProID       int      `json:"pro_id"`
	ProCode     string   `json:"pro_code"`
	ProName     string   `json:"pro_name"`
	ItemType    int      `json:"item_type"`
	Qty         *float64 `json:"qty"`
	QtyStr      *string  `json:"qty_str"`
	PurchPrice1 *float64 `json:"purch_price_1"`
	PurchPrice2 *float64 `json:"purch_price_2"`
	PurchPrice3 *float64 `json:"purch_price_3"`
	SellPrice1  *float64 `json:"sell_price_1"`
	SellPrice2  *float64 `json:"sell_price_2"`
	SellPrice3  *float64 `json:"sell_price_3"`
	Amount      *float64 `json:"amount"`
	DiscValue   *float64 `json:"disc_value"`
	BatchNo     *string  `json:"batch_no"`
	ExpDate     *string  `json:"exp_date"`
}
