package entity

import "time"

type CreateApPayBody struct {
	CustID             string                  `json:"cust_id"`
	ApPayDate          *string                 `json:"ap_pay_date"`
	TrCode             *string                 `json:"tr_code"`
	SupId              *int64                  `json:"sup_id"`
	CashAmt            *float64                `json:"cash_amt"`
	CndnAmt            *float64                `json:"cndn_amt"`
	ReturnAmt          *float64                `json:"return_amt"`
	ChequeAmt          *float64                `json:"cheque_amt"`
	TransferAmt        *float64                `json:"transfer_amt"`
	TotalAmt           *float64                `json:"total_amt"`
	DataStatus         *int64                  `json:"data_status"`
	CreatedBy          *int64                  `json:"created_by"`
	IsPosted           *bool                   `json:"is_posted"`
	Details            []CreateApPayDetBody    `json:"details"`
	ApPayMethodDetails []CreateApPayMethodBody `json:"ap_pay_method_details"`
}

type ApPayRespone struct {
	ApPayNo            string               `json:"ap_pay_no"`
	ApPayDate          *string              `json:"ap_pay_date"`
	TrCode             *string              `json:"tr_code"`
	SupId              *int64               `json:"sup_id"`
	SupCode            string               `json:"sup_code"`
	SupName            string               `json:"sup_name"`
	CashAmt            *float64             `json:"cash_amt"`
	CndnAmt            *float64             `json:"cndn_amt"`
	ReturnAmt          *float64             `json:"return_amt"`
	ChequeAmt          *float64             `json:"cheque_amt"`
	TransferAmt        *float64             `json:"transfer_amt"`
	TotalAmt           *float64             `json:"total_amt"`
	DataStatus         *int64               `json:"data_status"`
	CreatedBy          *int64               `json:"created_by"`
	UpdatedBy          *int64               `json:"updated_by"`
	IsPosted           *bool                `json:"is_posted"`
	UpdatedAt          time.Time            `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName      *string              `json:"updated_by_name"`
	Details            []ApPayDetResponse   `json:"details"`
	ApPayMethodDetails []ApPayMethodRespone `json:"ap_pay_method_details"`
}

type ApPayListResponse struct {
	ApPayNo       string    `json:"ap_pay_no"`
	ApPayDate     *string   `json:"ap_pay_date"`
	TrCode        *string   `json:"tr_code"`
	SupId         *int64    `json:"sup_id"`
	SupCode       string    `json:"sup_code"`
	SupName       string    `json:"sup_name"`
	CashAmt       *float64  `json:"cash_amt"`
	CndnAmt       *float64  `json:"cndn_amt"`
	ReturnAmt     *float64  `json:"return_amt"`
	ChequeAmt     *float64  `json:"cheque_amt"`
	TransferAmt   *float64  `json:"transfer_amt"`
	TotalAmt      *float64  `json:"total_amt"`
	DataStatus    *int64    `json:"data_status"`
	CreatedBy     *int64    `json:"created_by"`
	UpdatedBy     *int64    `json:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedByName *string   `json:"updated_by_name"`
	IsPosted      *bool     `json:"is_posted"`
}

type UpdateApPayBody struct {
	CustID             string                  `json:"cust_id"`
	ApPayDate          *string                 `json:"ap_pay_date"`
	TrCode             *string                 `json:"tr_code"`
	SupId              *int64                  `json:"sup_id"`
	CashAmt            *float64                `json:"cash_amt"`
	CndnAmt            *float64                `json:"cndn_amt"`
	ReturnAmt          *float64                `json:"return_amt"`
	ChequeAmt          *float64                `json:"cheque_amt"`
	TransferAmt        *float64                `json:"transfer_amt"`
	TotalAmt           *float64                `json:"total_amt"`
	DataStatus         *int64                  `json:"data_status"`
	CreatedBy          *int64                  `json:"created_by"`
	UpdatedBy          int64                   `json:"updated_by"`
	IsPosted           *bool                   `json:"is_posted"`
	Details            []UpdateApPayDetBody    `json:"details"`
	ApPayMethodDetails []UpdateApPayMethodBody `json:"ap_pay_method_details"`
}

type DetailApPayParams struct {
	ApPayNo string `params:"ap_pay_no" validate:"required" json:"ap_pay_no"`
}

type DeleteApPayParams struct {
	ApPayNo string `params:"ap_pay_no" validate:"required" json:"ap_pay_no"`
}
type UpdateApPayParams struct {
	ApPayNo string `params:"ap_pay_no" validate:"required" json:"ap_pay_no"`
}
