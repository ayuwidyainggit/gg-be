package entity

type CreateCollectionPayRequest struct {
	CustID        string
	EmpID         int64
	UserID        int64
	CollectionNo  string                `json:"collection_no"`
	OprType       string                `json:"opr_type"`
	OutletID      int64                 `json:"outlet_id"`
	PaymentOption string                `json:"payment_option"`
	Detail        []CollectionPayDetail `json:"detail"`
}

type CollectionPayDetail struct {
	InvoiceNo          string                    `json:"invoice_no"`
	OrderNo            string                    `json:"order_no"`
	Discount           float64                   `json:"discount"`
	PaymentBalance     float64                   `json:"payment_balance"`
	InvoiceAmount      float64                   `json:"invoice_amount"`
	PaidAmount         float64                   `json:"paid_amount"`
	NewRemainingAmount float64                   `json:"new_remaining_amount"`
	IsCollection       bool                      `json:"is_collection"`
	Payment            []PaymentCollectionDetail `json:"payment"`
	Notes              string                    `json:"notes"`
	Image              []string                  `json:"image"`
}

type PaymentCollectionDetail struct {
	InvoiceNo     string  `json:"invoice_no"`
	PayType       int     `json:"pay_type"`
	PaymentAmount float64 `json:"payment_amount"`
	CndnJenis     int     `json:"cndn_jenis"`
}

type CollectionPayQueryFilter struct {
	InvoiceNo string `query:"invoice_no" validate:"required"`
	OutletID  int64  `query:"outlet_id" validate:"required"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
	Sort      string `query:"sort"`
}

type CollectionPayResponse struct {
	PaymentTrxID     int     `json:"payment_trx_id"`
	OutletID         int64   `json:"outlet_id"`
	EmpID            int64   `json:"emp_id"`
	PONumber         string  `json:"po_number"`
	DocumentNo       string  `json:"document_no"`
	TotalTransaction float64 `json:"total_transaction"`
	PaymentAmount    float64 `json:"payment_amount"`
	RemainingAmount  float64 `json:"remaining_amount"`
	PayType          int     `json:"pay_type"`
	Amount           float64 `json:"amount"`
	PaymentDate      string  `json:"payment_date"`
}

type CreateNoPaymentRequest struct {
	EmpID       int64
	CustID      string
	UserID      int64
	Reason      string `json:"reason" validate:"required"`
	ReasonID    int64  `json:"reason_id" validate:"required"`
	OutletID    int64  `json:"outlet_id" validate:"required"`
	PaymentDate string `json:"payment_date" validate:"required,datetime=2006-01-02"`
}

type StoreCollectionPayResponse struct {
	InvoiceNo string `json:"invoice_no"`
	DocNoBank string `json:"doc_no_bank"`
	DocNoCash string `json:"doc_no_cash,omitempty"`
}
