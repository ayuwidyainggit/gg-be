package entity

import "time"

type InvoicesListReq struct {
	OutletID int64 `query:"outlet_id" validate:"required"`
}

// type InvoicesCreate struct {
// 	InvoiceNumber string          `json:"invoice_number"`
// 	OrderNumber   string          `json:"order_number"`
// 	PaymentOption string          `json:"payment_option"`
// 	PaymentDetail []PaymentDetail `json:"payment_detail"`
// 	Note          string          `json:"note"`
// 	Images        []string        `json:"images"`
// }

type InvoicesPaymentCreate struct {
	OprType             OperationType         `json:"opr_type" validate:"oneof=O C"`
	CollectionNo        string                `json:"collection_no"`
	TotalDiscount       float64               `json:"total_discount"`
	TotalMaterai        float64               `json:"total_materai"`
	TotalPaymentBalance float64               `json:"total_payment_balance" validate:"required"`
	TotalPayment        float64               `json:"total_payment" validate:"required"`
	OutletID            int64                 `json:"outlet_id" validate:"required"`
	PaymentOption       string                `json:"payment_option" validate:"oneof=installment full"`
	Notes               *string               `json:"notes"`
	Files               []string              `json:"files"`
	Detail              InvoicesPaymentDetail `json:"detail" validate:"required"`

	// Internal fields
	CustID     string     `json:"-"`
	EmpGrpID   int64      `json:"-"`
	EmpID      int64      `json:"-"`
	SalesmanID int64      `json:"-"`
	CreatedBy  int64      `json:"-"`
	CreatedAt  *time.Time `json:"-"`
	UpdatedBy  *int64     `json:"-"`
	UpdatedAt  *time.Time `json:"-"`
}

type InvoicesGetResp struct {
	InvoiceNumber  string          `json:"invoice_number"`
	OrderNumber    string          `json:"order_number"`
	InvoiceDate    string          `json:"invoice_date"`
	DueDate        string          `json:"due_date"`
	PaymentOption  string          `json:"payment_option"`
	SettlementDate string          `json:"settlement_date"`
	InvoiceAmount  int             `json:"invoice_amount"`
	PaymentDetail  []PaymentDetail `json:"payment_detail"`
	Images         []string        `json:"images"`
	CollectHistory []any           `json:"collect_history"`
}
type PaymentDetail struct {
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

type InvoicesPaymentDetail struct {
	InvoiceNo      string                   `json:"invoice_no" validate:"required"`
	Discount       float64                  `json:"discount"`
	Materai        *float64                 `json:"materai"`
	PaymentBalance float64                  `json:"payment_balance" validate:"required"`
	InvoiceAmount  float64                  `json:"invoice_amount" validate:"required"`
	TotalPayment   float64                  `json:"total_payment" validate:"required"`
	IsCollection   bool                     `json:"is_collection"`
	Payment        []InvoicesPaymentDeposit `json:"payment" validate:"required,dive,required"`

	// Internal fields
	DepositDetailID  int64   `json:"-"`
	DepositNo        string  `json:"-"`
	OutletID         int64   `json:"-"`
	RemainingPayment float64 `json:"-"`
	InvoiceDate      string  `json:"-"`
	DocumentNo       string  `json:"-"`
}

type InvoicesPaymentDeposit struct {
	InvoiceNo     string                         `json:"invoice_no" validate:"required"`
	PayType       PaymentType                    `json:"pay_type" validate:"required"`
	CndnJenis     int                            `json:"cndn_jenis"`
	DocumentNo    string                         `json:"document_no"`
	Balance       float64                        `json:"balance"`
	PaymentAmount float64                        `json:"payment_amount" validate:"required"`
	Images        []InvoicesPaymentDepositImages `json:"images"`

	// Internal fields
	DepositPaymentID int64  `json:"-"`
	DepositNo        string `json:"-"`
}

type InvoicesPaymentDepositImages struct {
	ImageURL string `json:"image_url"`

	// Internal fields
	DepositImageID int    `json:"-"`
	DepositNo      string `json:"-"`
	InvoiceNo      string `json:"-"`
}

type DetilPaymentDepositInvoiceDetailResponse struct {
	// DepositResponse
	// Details []DepositDetail         `json:"details"`
	Cash    []DetilDepositPaymentInvoice `json:"cash"`
	Cek     []DetilDepositPaymentInvoice `json:"cek"`
	Trasfer []DetilDepositPaymentInvoice `json:"transfer"`
	Return  []DetilDepositPaymentInvoice `json:"return"`
	CNDN    []DetilDepositPaymentInvoice `json:"cndn"`
}

type DepositInvoicePayment struct {
	DepositPaymentID int64   `json:"deposit_payment_id"`
	DepositNo        string  `json:"deposit_no"`
	InvoiceNo        string  `json:"invoice_no"`
	PayType          int16   `json:"pay_type"`
	DocumentNo       string  `json:"document_no"`
	Balance          float64 `json:"balance"`
	PaymentAmount    float64 `json:"payment_amount"`

	Images []DepositInvoicePaymentImage `json:"images"`
}

type DepositInvoicePaymentImage struct {
	DepositImageID int    `json:"deposit_image_id"`
	DepositNo      string `json:"deposit_no"`
	InvoiceNo      string `json:"invoice_no"`
	ImageUrl       string `json:"image_url"`
}

type DetilDepositPaymentInvoice struct {
	DepositInvoicePayment
	SalesmanId     *int64  `json:"salesman_id"`
	SalesmanCode   *string `json:"salesman_code"`
	SalesmanName   *string `json:"salesman_name"`
	OutletID       *int64  `json:"outlet_id"`
	OutletCode     *string `json:"outlet_code"`
	OutletName     *string `json:"outlet_name"`
	InvoiceDate    string  `json:"invoice_date"`
	Materai        float64 `json:"materai"`
	Discount       float64 `json:"discount"`
	PaymentBalance float64 `json:"payment_balance"`
	TotalPayment   float64 `json:"total_payment"`
}
