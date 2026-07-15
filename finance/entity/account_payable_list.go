package entity

type AccountPayableListResponse struct {
	InvDate         *string `json:"inv_date"`
	InvDueDate      *string `json:"inv_due_date"`
	InvNo           string  `json:"inv_no"`
	InvAmount       float64 `json:"inv_amount"`
	AmountPaid      float64 `json:"amount_paid"`
	RemainingAmount float64 `json:"remaining_amount"`
	SupplierId      int64   `json:"supplier_id"`
	SupplierCode    string  `json:"supplier_code"`
	Supplier        string  `json:"supplier"`
	DistributorId   int64   `json:"distributor_id"`
	DistributorCode string  `json:"distributor_code"`
	Distributor     string  `json:"distributor"`
	InvStatus       string  `json:"inv_status"`
	DueDateStatus   string  `json:"due_date_status"`
	Aging           int64   `json:"aging"`
	CreatedAt       *string `json:"created_at"`
	CreatedByName   *string `json:"created_by_name"`
	UpdatedAt       *string `json:"updated_at"`
	UpdatedByName   string  `json:"updated_by_name"`
}

type AccountPayableListDetailResponse struct {
	PoNo            string                                     `json:"po_no"`
	InvDate         *string                                    `json:"inv_date"`
	InvDueDate      *string                                    `json:"inv_due_date"`
	InvNo           *string                                    `json:"inv_no"`
	InvAmount       float64                                    `json:"inv_amount"`
	AmountPaid      float64                                    `json:"amount_paid"`
	RemainingAmount *float64                                   `json:"remaining_amount"`
	SupplierId      int64                                      `json:"supplier_id"`
	SupplierCode    string                                     `json:"supplier_code"`
	Supplier        string                                     `json:"supplier"`
	DistributorId   int64                                      `json:"distributor_id"`
	DistributorCode string                                     `json:"distributor_code"`
	Distributor     string                                     `json:"distributor"`
	InvStatus       string                                     `json:"inv_status"`
	DueDateStatus   *string                                    `json:"due_date_status"`
	Aging           int64                                      `json:"aging"`
	CreatedAt       *string                                    `json:"created_at"`
	CreatedByName   *string                                    `json:"created_by_name"`
	UpdatedAt       string                                     `json:"updated_at"`
	UpdatedByName   string                                     `json:"updated_by_name"`
	PaymentHistory  []AccountPayableListPaymentHistoryResponse `json:"payment_history"`
}

type AccountPayableListPaymentHistoryResponse struct {
	PaymentMethod     *int     `json:"payment_method"`
	PaymentMethodName string   `json:"payment_method_name"`
	PaymentDate       *string  `json:"payment_date"`
	PaymentBalance    *float64 `json:"payment_balance"`
	DocumentNo        string   `json:"document_no"`
	Amount            float64  `json:"amount"`
	UpdatedBy         int64    `json:"updated_by"`
	UpdatedByName     string   `json:"updated_by_name"`
	UpdatedAt         string   `json:"updated_at"`
}

type DetailAccountPayableListParams struct {
	InvNo string `params:"inv_no" validate:"required" json:"inv_no"`
}

type AccountPayableListQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:""`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	InvoiceNo    string `query:"inv_no"`
	Supplier     int    `query:"supplier_id"`
}
