package entity

type CreateApPayDetBody struct {
	PayAmount *float64 `json:"pay_amount"`
}

type ApPayDetResponse struct {
	ApPayNo    string   `json:"ap_pay_no"`
	ApPayDetId int      `json:"ap_pay_det_id"`
	PayAmount  *float64 `json:"pay_amount"`
}

type UpdateApPayDetBody struct {
	ApPayDetId *int64   `json:"ap_pay_det_id"`
	PayAmount  *float64 `json:"pay_amount"`
}
