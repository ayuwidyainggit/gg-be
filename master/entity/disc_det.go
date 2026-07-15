package entity

type DiscDet struct {
	// CustID    string   `json:"cust_id"`
	DiscID    int      `json:"disc_id"`
	RowNo     int      `json:"row_no"`
	MinValue  *float64 `json:"min_value"`
	MaxValue  *float64 `json:"max_value"`
	DiscType  *int64   `json:"disc_type"`
	DiscPerc  *float64 `json:"disc_perc"`
	DiscValue *float64 `json:"disc_value"`
	DiscDetID *int64   `json:"disc_det_id"`
}

type CreateDiscDetBody struct {
	RowNo     int      `json:"row_no"`
	MinValue  *float64 `json:"min_value"`
	MaxValue  *float64 `json:"max_value"`
	DiscType  *int64   `json:"disc_type"`
	DiscPerc  *float64 `json:"disc_perc"`
	DiscValue *float64 `json:"disc_value"`
}
type UpdateDiscDetBody struct {
	DiscDetID *int64   `json:"disc_det_id"`
	RowNo     int      `json:"row_no"`
	MinValue  *float64 `json:"min_value"`
	MaxValue  *float64 `json:"max_value"`
	DiscType  *int64   `json:"disc_type"`
	DiscPerc  *float64 `json:"disc_perc"`
	DiscValue *float64 `json:"disc_value"`
}
