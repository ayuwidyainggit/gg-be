package model

type ApPayDet struct {
	ApPayDetId int64    `gorm:"column:ap_pay_det_id;primaryKey" json:"ap_pay_det_id"`
	CustID     string   `gorm:"column:cust_id" json:"cust_id"`
	ApPayNo    string   `gorm:"column:ap_pay_no" json:"ap_pay_no"`
	PayAmount  *float64 `gorm:"column:pay_amount" json:"pay_amount"`
}

func (ApPayDet) TableName() string {
	return "acf.ap_pay_det"
}
