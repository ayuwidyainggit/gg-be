package model

type DiscDet struct {
	CustID    string   `db:"cust_id" json:"cust_id"`
	DiscID    int64    `db:"disc_id" json:"disc_id"`
	DiscCode  int64    `db:"disc_code" json:"disc_code"`
	DiscName  int64    `db:"disc_name" json:"disc_name"`
	RowNo     int      `db:"row_no" json:"row_no"`
	MinValue  *float64 `db:"min_value" json:"min_value"`
	MaxValue  *float64 `db:"max_value" json:"max_value"`
	DiscType  *int64   `db:"disc_type" json:"disc_type"`
	DiscPerc  *float64 `db:"disc_perc" json:"disc_perc"`
	DiscValue *float64 `db:"disc_value" json:"disc_value"`
	DiscDetID *int64   `db:"disc_det_id" json:"disc_det_id"`
}
type DiscDetUpdate struct {
	RowNo     *int     `db:"row_no" json:"row_no"`
	MinValue  *float64 `db:"min_value" json:"min_value"`
	MaxValue  *float64 `db:"max_value" json:"max_value"`
	DiscType  *int64   `db:"disc_type" json:"disc_type"`
	DiscPerc  *float64 `db:"disc_perc" json:"disc_perc"`
	DiscValue *float64 `db:"disc_value" json:"disc_value"`
}
