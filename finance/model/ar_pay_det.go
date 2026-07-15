package model

type ArPayDet struct {
	ArPayDetID   *int64   `gorm:"column:ar_pay_det_id;primaryKey" json:"ar_pay_det_id"`
	CustID       string   `gorm:"column:cust_id" json:"cust_id"`
	ArPayNo      string   `gorm:"column:ar_pay_no" json:"ar_pay_no"`
	SoNo         *string  `gorm:"column:so_no" json:"so_no"`
	CashAmt      *float64 `gorm:"column:cash_amt" json:"cash_amt"`
	ChqTrNo      *string  `gorm:"column:chq_tr_no" json:"chq_tr_no"`
	ChqBlc       *float64 `gorm:"column:chq_blc" json:"chq_blc"`
	ChqAmt       *float64 `gorm:"column:chq_amt" json:"chq_amt"`
	TrfTrNo      *string  `gorm:"column:trf_tr_no" json:"trf_tr_no"`
	TrfBlc       *float64 `gorm:"column:trf_blc" json:"trf_blc"`
	TrfAmt       *float64 `gorm:"column:trf_amt" json:"trf_amt"`
	ReturnNo     *string  `gorm:"column:return_no" json:"return_no"`
	ReturnBlc    *float64 `gorm:"column:return_blc" json:"return_blc"`
	ReturnAmt    *float64 `gorm:"column:return_amt" json:"return_amt"`
	CndnNo       *string  `gorm:"column:cndn_no" json:"cndn_no"`
	CndnBlc      *float64 `gorm:"column:cndn_blc" json:"cndn_blc"`
	CndnAmt      *float64 `gorm:"column:cndn_amt" json:"cndn_amt"`
	DiscAmt      *float64 `gorm:"column:disc_amt" json:"disc_amt"`
	DutyStampAmt *float64 `gorm:"column:duty_stamp_amt" json:"duty_stamp_amt"`
	PayAmt       *float64 `gorm:"column:pay_amt" json:"pay_amt"`
	PayDiff      *float64 `gorm:"column:pay_diff" json:"pay_diff"`
	PayAmtRound  *float64 `gorm:"column:pay_amt_round" json:"pay_amt_round"`
}

func (ArPayDet) TableName() string {
	return "acf.ar_pay_det"
}
