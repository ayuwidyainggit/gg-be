package entity

type ConversionResponse struct {
	Data struct {
		Qty1     int `json:"qty1"`
		Qty2     int `json:"qty2"`
		Qty3     int `json:"qty3"`
		TotalQty int `json:"total_qty"`
	} `json:"data"`
}
