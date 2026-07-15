package entity

type CreateArPayDetBody struct {
	SoNo         *string  `json:"so_no"`
	CashAmt      *float64 `json:"cash_amt"`
	ChqTrNo      *string  `json:"chq_tr_no"`
	ChqBlc       *float64 `json:"chq_blc"`
	ChqAmt       *float64 `json:"chq_amt"`
	TrfTrNo      *string  `json:"trf_tr_no"`
	TrfBlc       *float64 `json:"trf_blc"`
	TrfAmt       *float64 `json:"trf_amt"`
	ReturnNo     *string  `json:"return_no"`
	ReturnBlc    *float64 `json:"return_blc"`
	ReturnAmt    *float64 `json:"return_amt"`
	CndnNo       *string  `json:"cndn_no"`
	CndnBlc      *float64 `json:"cndn_blc"`
	CndnAmt      *float64 `json:"cndn_amt"`
	DiscAmt      *float64 `json:"disc_amt"`
	DutyStampAmt *float64 `json:"duty_stamp_amt"`
	PayAmt       *float64 `json:"pay_amt"`
	PayDiff      *float64 `json:"pay_diff"`
	PayAmtRound  *float64 `json:"pay_amt_round"`
}

type ArPayDetResponse struct {
	ArPayDetID   int64    `json:"ar_pay_det_id"`
	SoNo         *string  `json:"so_no"`
	CashAmt      *float64 `json:"cash_amt"`
	ChqTrNo      *string  `json:"chq_tr_no"`
	ChqBlc       *float64 `json:"chq_blc"`
	ChqAmt       *float64 `json:"chq_amt"`
	TrfTrNo      *string  `json:"trf_tr_no"`
	TrfBlc       *float64 `json:"trf_blc"`
	TrfAmt       *float64 `json:"trf_amt"`
	ReturnNo     *string  `json:"return_no"`
	ReturnBlc    *float64 `json:"return_blc"`
	ReturnAmt    *float64 `json:"return_amt"`
	CndnNo       *string  `json:"cndn_no"`
	CndnBlc      *float64 `json:"cndn_blc"`
	CndnAmt      *float64 `json:"cndn_amt"`
	DiscAmt      *float64 `json:"disc_amt"`
	DutyStampAmt *float64 `json:"duty_stamp_amt"`
	PayAmt       *float64 `json:"pay_amt"`
	PayDiff      *float64 `json:"pay_diff"`
	PayAmtRound  *float64 `json:"pay_amt_round"`
}

type UpdateArPayDetBody struct {
	ArPayDetID   *int64   `json:"ar_pay_det_id"`
	SoNo         *string  `json:"so_no"`
	CashAmt      *float64 `json:"cash_amt"`
	ChqTrNo      *string  `json:"chq_tr_no"`
	ChqBlc       *float64 `json:"chq_blc"`
	ChqAmt       *float64 `json:"chq_amt"`
	TrfTrNo      *string  `json:"trf_tr_no"`
	TrfBlc       *float64 `json:"trf_blc"`
	TrfAmt       *float64 `json:"trf_amt"`
	ReturnNo     *string  `json:"return_no"`
	ReturnBlc    *float64 `json:"return_blc"`
	ReturnAmt    *float64 `json:"return_amt"`
	CndnNo       *string  `json:"cndn_no"`
	CndnBlc      *float64 `json:"cndn_blc"`
	CndnAmt      *float64 `json:"cndn_amt"`
	DiscAmt      *float64 `json:"disc_amt"`
	DutyStampAmt *float64 `json:"duty_stamp_amt"`
	PayAmt       *float64 `json:"pay_amt"`
	PayDiff      *float64 `json:"pay_diff"`
	PayAmtRound  *float64 `json:"pay_amt_round"`
}
