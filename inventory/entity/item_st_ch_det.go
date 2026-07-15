package entity

type CreateItemStChDetBody struct {
	IscNo   string   `json:"isc_no"`
	CustID  string   `json:"cust_id"`
	SeqNo   int      `json:"seq_no"`
	ProID   int      `json:"pro_id"`
	Qty     *float64 `json:"qty"`
	QtyStr  *string  `json:"qty_str"`
	BatchNo *string  `json:"batch_no"`
	ExpDate *string  `json:"exp_date"`
}

type ItemStChDetResponse struct {
	ItemStChDetId *int     `json:"isc_det_id"`
	SeqNo         int      `json:"seq_no"`
	ProID         int      `json:"pro_id"`
	ProCode       string   `json:"pro_code"`
	ProName       string   `json:"pro_name"`
	Qty           *float64 `json:"qty"`
	QtyStr        *string  `json:"qty_str"`
	BatchNo       *string  `json:"batch_no"`
	ExpDate       *string  `json:"exp_date"`
}

type UpdateItemStChDetBody struct {
	IscDetID *int     `json:"isc_det_id"`
	IscNo    string   `json:"isc_no"`
	CustID   string   `json:"cust_id"`
	SeqNo    int      `json:"seq_no"`
	ProID    int      `json:"pro_id"`
	Qty      *float64 `json:"qty"`
	QtyStr   *string  `json:"qty_str"`
	BatchNo  *string  `json:"batch_no"`
	ExpDate  *string  `json:"exp_date"`
}
