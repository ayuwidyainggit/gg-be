package entity

import "time"

type ApPaymentQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:""`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	SuppId       int    `query:"sup_id"`
	DocumentNo   string `query:"document_no"`
	Type         string `query:"Type"`
}

type CreateAccountPayablePaymentBody struct {
	CustId                    string                                  `json:"cust_id"`
	AccountPayablePaymentNo   string                                  `json:"account_payable_payment_no"`
	AccountPayablePaymentDate *string                                 `json:"account_payable_payment_date"`
	SupId                     *int64                                  `json:"sup_id"`
	TotalDiscount             *float64                                `json:"total_discount"`
	TotalPaymentBalance       *float64                                `json:"total_payment_balance"`
	TotalMateri               *float64                                `json:"total_materai"`
	TotalPayment              *float64                                `json:"total_payment"`
	Details                   []CreateAccountPayablePaymentDetailBody `json:"detail" validate:"required,dive,required"`
	CreatedBy                 int64                                   `json:"created_by"`
	UpdatedBy                 int64                                   `json:"updated_by"`
}

type CreateAccountPayablePaymentDetailBody struct {
	AccountPayablePaymentNo *string                                  `json:"account_payable_payment_no"`
	InvoiceNo               *string                                  `json:"invoice_no"`
	InvoiceDate             *string                                  `json:"invoice_date"`
	InvoiceAmount           *float64                                 `json:"invoice_amount"`
	PaidAmount              *float64                                 `json:"paid_amount"`
	RemainingAmount         *float64                                 `json:"remaining_amount"`
	Discount                *float64                                 `json:"discount"`
	PaymentBalance          *float64                                 `json:"payment_balance"`
	Materai                 *float64                                 `json:"materai"`
	TotalPayment            *float64                                 `json:"total_payment"`
	Payment                 []CreateAccountPayablePaymentOptionsBody `json:"payment" validate:"required,dive,required"`
}

type CreateAccountPayablePaymentOptionsBody struct {
	// AccountPayablePaymentNo *string  `json:"account_payable_payment_no"`
	InvoiceNo     *string  `json:"invoice_no"`
	PayType       *int64   `json:"pay_type"`
	DocumentNo    *string  `json:"document_no"`
	Balance       *float64 `json:"balance"`
	PaymentAmount *float64 `json:"payment_amount"`
}

type UpdateAccountPayablePaymentBody struct {
	CustId string `json:"cust_id"`
	// AccountPayablePaymentNo   *string                                 `json:"account_payable_payment_no"`
	AccountPayablePaymentDate *string                                 `json:"account_payable_payment_date"`
	SupId                     *int64                                  `json:"sup_id"`
	TotalDiscount             *float64                                `json:"total_discount"`
	TotalPaymentBalance       *float64                                `json:"total_payment_balance"`
	TotalMateri               *float64                                `json:"total_materai"`
	TotalPayment              *float64                                `json:"total_payment"`
	Details                   []UpdateAccountPayablePaymentDetailBody `json:"detail" validate:"required,dive,required"`
	UpdatedBy                 int64                                   `json:"updated_by"`
}

type UpdateAccountPayablePaymentDetailBody struct {
	// AccountPayablePaymentNo *string                                  `json:"account_payable_payment_no"`
	InvoiceNo       *string                                  `json:"invoice_no"`
	InvoiceDate     *string                                  `json:"invoice_date"`
	InvoiceAmount   *float64                                 `json:"invoice_amount"`
	PaidAmount      *float64                                 `json:"paid_amount"`
	RemainingAmount *float64                                 `json:"remaining_amount"`
	Discount        *float64                                 `json:"discount"`
	PaymentBalance  *float64                                 `json:"payment_balance"`
	Materai         *float64                                 `json:"materai"`
	TotalPayment    *float64                                 `json:"total_payment"`
	Payment         []UpdateAccountPayablePaymentOptionsBody `json:"payment" validate:"required,dive,required"`
}

type UpdateAccountPayablePaymentOptionsBody struct {
	AccountPayablePaymentNo *string  `json:"account_payable_payment_no"`
	InvoiceNo               *string  `json:"invoice_no"`
	PayType                 *int64   `json:"pay_type"`
	DocumentNo              *string  `json:"document_no"`
	Balance                 *float64 `json:"balance"`
	PaymentAmount           *float64 `json:"payment_amount"`
}

type AccountPayablePaymentList struct {
	CustId                    string     `json:"cust_id"`
	DocumentNo                *string    `json:"document_no"`
	AccountPayablePaymentNo   *string    `json:"account_payable_payment_no"`
	AccountPayablePaymentDate *string    `json:"account_payable_payment_date"`
	SupId                     *int64     `json:"sup_id"`
	SupName                   *string    `json:"sup_name"`
	SupCode                   *string    `json:"sup_code"`
	DistributorId             *int64     `json:"distributor_id"`
	DistributorName           *string    `json:"distributor"`
	DistributorCode           *string    `json:"distributor_code"`
	TotalPayment              *float64   `json:"total_payment"`
	CreatedBy                 *int       `json:"created_by"`
	CreatedByName             *string    `json:"created_by_name"`
	CreatedAt                 *time.Time `json:"created_at"`
	UpdatedBy                 *int       `json:"updated_by"`
	UpdatedByName             *string    `json:"updated_by_name"`
	UpdatedAt                 *time.Time `json:"updated_at"`
}

type AccountPayablePaymentRespone struct {
	CustId                    string   `json:"cust_id"`
	AccountPayablePaymentNo   *string  `json:"account_payable_payment_no"`
	AccountPayablePaymentDate *string  `json:"account_payable_payment_date"`
	SupId                     *int64   `json:"sup_id"`
	SupName                   *string  `json:"sup_name"`
	SupCode                   *string  `json:"sup_code"`
	DistributorId             *int64   `json:"distributor_id"`
	DistributorName           *string  `json:"distributor"`
	DistributorCode           *string  `json:"distributor_code"`
	TotalPayment              *float64 `json:"total_payment"`
	CreatedBy                 int64    `json:"created_by"`
	UpdatedBy                 int64    `json:"updated_by"`
}

type AccountPayablePaymentDetailRespone struct {
	AccountPayablePaymentNo *string                               `json:"account_payable_payment_no"`
	InvoiceNo               *string                               `json:"invoice_no"`
	InvoiceDate             *string                               `json:"invoice_date"`
	InvoiceAmount           *float64                              `json:"invoice_amount"`
	PaidAmount              *float64                              `json:"paid_amount"`
	RemainingAmount         *float64                              `json:"remaining_amount"`
	PaymentBalance          *float64                              `json:"payment_balance"`
	TotalPayment            *float64                              `json:"total_payment"`
	Payment                 []AccountPayablePaymentOptionsRespone `json:"payment" validate:"required,dive,required"`
}

type AccountPayablePaymentOptionsRespone struct {
	AccountPayablePaymentNo *string  `json:"account_payable_payment_no"`
	InvoiceNo               *string  `json:"invoice_no"`
	InvoiceDate             *string  `json:"invoice_date"`
	PayType                 *int64   `json:"pay_type"`
	PayTypeName             string   `json:"pay_type_name"`
	DocumentNo              *string  `json:"document_no"`
	Balance                 *float64 `json:"balance"`
	PaymentAmount           float64  `json:"payment_amount"`
	PaymentBalance          float64  `json:"payment_balance"`
	RemainingAmount         *float64 `json:"remaining_amount"`
}

type AccountPayablePaymentDetailResponse struct {
	AccountPayablePaymentRespone
	Details       []AccountPayablePaymentDetailRespone  `json:"details"`
	Cash          []AccountPayablePaymentOptionsRespone `json:"cash"`
	TotalCash     float64                               `json:"total_cash"`
	Cek           []AccountPayablePaymentOptionsRespone `json:"cek"`
	TotalCek      float64                               `json:"total_cek"`
	Trasfer       []AccountPayablePaymentOptionsRespone `json:"transfer"`
	TotalTransfer float64                               `json:"total_transfer"`
	Return        []AccountPayablePaymentOptionsRespone `json:"return"`
	TotalRetrun   float64                               `json:"total_return"`
	Cndn          []AccountPayablePaymentOptionsRespone `json:"cndn"`
	TotalCndn     float64                               `json:"total_cndn"`
}

type DetailAccountPayablePaymentParams struct {
	AccountPayablePaymentNo string `params:"account_payable_payment_no" validate:"required" json:"account_payable_payment_no"`
}

type UpdateAccountPayablePaymentParams struct {
	AccountPayablePaymentNo string `params:"account_payable_payment_no" validate:"required" json:"account_payable_payment_no"`
}

type DeleteAccountPayablePaymentParams struct {
	AccountPayablePaymentNo string `params:"account_payable_payment_no" validate:"required" json:"account_payable_payment_no"`
}

type ApLookupSupplierInoviceReturnQueryFilter struct {
	CustId              string
	ParentCustId        string
	From                *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To                  *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page                int    `query:"page"`
	Limit               int    `query:"limit" validate:""`
	Query               string `query:"q"`
	Mode                string `query:"mode"`
	Sort                string `query:"sort"`
	SuppId              int    `query:"sup_id"`
	DocumentNo          string `query:"document_no"`
	Type                string `query:"Type"`
	InvoiceNo           string `query:"invoice_no"`
	ExcludeEmptyInvoice bool   `query:"exclude_empty_invoice"`
}

type ApLookupSupplierInvoiceReturnResponeList struct {
	ApType          string  `json:"ap_type"`
	SupId           int64   `json:"sup_id"`
	SupName         string  `json:"sup_name"`
	InvoiceNo       string  `json:"invoice_no"`
	DocumentNo      string  `json:"document_no"`
	Amount          float64 `json:"amount"`
	SubTotal        float64 `json:"sub_total"`
	PaidAmount      float64 `json:"paid_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Total           float64 `json:"total"`
	CreatedBy       int64   `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	UpdatedBy       int64   `json:"updated_by"`
	UpdatedByName   string  `json:"updated_by_name"`
}
