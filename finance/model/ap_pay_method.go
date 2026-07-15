package model

type ApPayMethod struct {
	ApPayMethodId *int64   `gorm:"column:ap_pay_method_id;primaryKey" json:"ap_pay_method_id"`
	CustID        string   `gorm:"column:cust_id" json:"cust_id"`
	ApPayNo       string   `gorm:"column:ap_pay_no" json:"ap_pay_no"`
	PayMethodType *int64   `gorm:"column:pay_method_type" json:"pay_method_type"`
	RefNo         *string  `gorm:"column:ref_no" json:"ref_no"`
	Amount        *float64 `gorm:"column:amount" json:"amount"`
}

func (ApPayMethod) TableName() string {
	return "acf.ap_pay_method"
}
