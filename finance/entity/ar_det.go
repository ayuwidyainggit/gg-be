package entity

type CreateArDetBody struct {
	ArNo     string   `json:"ar_no"`
	SoNo     string   `json:"so_no"`
	ArAmount *float64 `json:"ar_amount"`
	ArPaid   *float64 `json:"ar_paid"`
}

type ArPaymentResponse struct {
	DepositPaymentID       int64   `json:"deposit_payment_id"`
	VisitDate              string  `json:"visit_date"`
	DepositDate            string  `json:"deposit_date"`
	DepositNo              string  `json:"deposit_no"`
	EmpId                  int64   `json:"emp_id"`
	EmpCode                *string `json:"emp_code"`
	EmpName                *string `json:"emp_name"`
	EmpGrpId               int64   `json:"emp_grp_id"`
	EmpGrpCode             *string `json:"emp_grp_code"`
	EmpGrpName             *string `json:"emp_grp_name"`
	CollectionStatus       int64   `json:"collection_status"`
	CollectionStatusName   string  `json:"collection_status_name"`
	PaymentOption          int64   `json:"payment_option"`
	PaymentOptionName      string  `json:"payment_option_name"`
	PaymentMethod          int64   `json:"payment_method"`
	PaymentMethodName      string  `json:"payment_method_name"`
	Amount                 float64 `json:"amount"`
	VerificationStatus     int64   `json:"verification_status"`
	VerificationStatusName string  `json:"verification_status_name"`
	VerifiedBy             *int64  `json:"verified_by"`
	VerifiedByName         *string `json:"verified_by_name"`
	VerifiedDate           *string `json:"verified_date"`
	Reason                 *string `json:"reason"`
	AdditionalInfo         *string `json:"additional_info"`
}

type UpdateArDetBody struct {
	ArDetID  *int64   `json:"ar_det_id"`
	ArNo     *string  `json:"ar_no"`
	SoNo     *string  `json:"so_no"`
	ArAmount *float64 `json:"ar_amount"`
	ArPaid   *float64 `json:"ar_paid"`
}

type CreateCollectionDetBody struct {
	CustID          string   `json:"cust_id"`
	CollectionNo    string   `json:"collection_no"`
	InvoiceNo       string   `json:"invoice_no"`
	SalesmanID      *int64   `json:"salesman_id"`
	InvoiceAmount   *float64 `json:"invoice_amount"`
	RemainingAmount *float64 `json:"remaining_amount"`
	PaidAmount      *float64 `json:"paid_amount"`
	CreatedBy       *int64   `json:"created_by"`
	CreatedAt       *string  `json:"created_at"`
}

type CollectionDetResponse struct {
	CollectionDetID        int64    `json:"collection_det_id"`
	CollectionNo           string   `json:"collection_no"`
	InvoiceNo              string   `json:"invoice_no"`
	SalesOrder             string   `json:"sales_order"`
	InvoiceDate            string   `json:"invoice_date"`
	DueDate                string   `json:"due_date"`
	SalesmanId             int64    `json:"salesman_id"`
	SalesmanName           string   `json:"salesman_name"`
	SalesmanCode           string   `json:"salesman_code"`
	OutletId               int64    `json:"outlet_id"`
	OutletCode             string   `json:"outlet_code"`
	OutletName             string   `json:"outlet_name"`
	InvoiceAmount          *float64 `json:"invoice_amount"`
	RemainingAmount        *float64 `json:"remaining_amount"`
	PaidAmount             *float64 `json:"paid_amount"`
	TotalInvoicePayment    *float64 `json:"total_invoice_amount"`
	InvoicePayment         *float64 `json:"invoice_payment"`
	PaidAmountByCollection *float64 `json:"paid_amount_by_collection"`
	CreatedBy              *int64   `json:"created_by"`
	CreatedByName          *string  `json:"created_by_name"`
	CreatedAt              *string  `json:"created_at"`
}

type UpdateCollectionDetBody struct {
	CollectionDetID *int64   `json:"collection_det_id"`
	CollectionNo    *string  `json:"collection_no"`
	InvoiceNo       *string  `json:"invoice_no"`
	SalesmanID      *int64   `json:"salesman_id"`
	InvoiceAmount   *float64 `json:"invoice_amount"`
	RemainingAmount *float64 `json:"remaining_amount"`
	PaidAmount      *float64 `json:"paid_amount"`
}

var dataCollectionStatusName = map[int64]string{
	1: "Paid",
	2: "Unpaid",
}

var dataPaymentOptionName = map[int64]string{
	1: "Full",
	2: "Partial",
}

var dataPaymentMethodName = map[int64]string{
	1: "Cash",
	2: "Cheque",
	3: "Transfer",
	4: "Return",
	5: "Credit",
}

var dataVerificationStatusName = map[int64]string{
	1: "In Review",
	2: "Approved",
	3: "Rejected",
}

func (arPayment ArPaymentResponse) GenerateDataCollectionStatusName() string {
	return dataCollectionStatusName[arPayment.CollectionStatus]
}

func (arPayment ArPaymentResponse) GenerateDataPaymentOptionName() string {
	return dataPaymentOptionName[arPayment.PaymentOption]
}

func (arPayment ArPaymentResponse) GenerateDataPaymentMethodName() string {
	return dataPaymentMethodName[arPayment.PaymentMethod]
}

func (arPayment ArPaymentResponse) GenerateDataVerificationStatusName() string {
	return dataVerificationStatusName[arPayment.VerificationStatus]
}
