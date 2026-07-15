package entity

type CreateApPayMethodBody struct {
	PayMethodType *int64   `json:"pay_method_type"`
	RefNo         *string  `json:"ref_no"`
	Amount        *float64 `json:"amount"`
}

type ApPayMethodRespone struct {
	ApPayNo       string   `json:"ap_pay_no"`
	ApPayMethodId int64    `json:"ap_pay_method_id"`
	PayMethodType *int64   `json:"pay_method_type"`
	RefNo         *string  `json:"ref_no"`
	Amount        *float64 `json:"amount"`
}

type UpdateApPayMethodBody struct {
	ApPayMethodId *int64   `json:"ap_pay_method_id"`
	PayMethodType *int64   `json:"pay_method_type"`
	RefNo         *string  `json:"ref_no"`
	Amount        *float64 `json:"amount"`
}
