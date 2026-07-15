package entity

const (
	STATUS_TAXES_RESERVED  = 1
	STATUS_TAXES_ACTIVE    = 2
	STATUS_TAXES_COMPLETED = 3
	STATUS_TAXES_INACTIVE  = 9
)

type MTaxesCreateReq struct {
	CustID                string `json:"cust_id"`
	Year                  int    `json:"year" validate:"required"`
	TransactionStatusCode string `json:"transaction_status_code" validate:"required"`
	SerialCode            string `json:"serial_code" validate:"required"`
	From                  int    `json:"from"`
	To                    int    `json:"to" validate:"required,gtfield=From"`
	TaxNumberAlert        int    `json:"tax_number_alert"`
	CreatedBy             int64  `json:"created_by"`
}

type MTaxesResp struct {
	CustID                string `json:"cust_id"`
	MTaxID                int64  `json:"m_tax_id" `
	Year                  int    `json:"year" `
	TransactionStatusCode string `json:"transaction_status_code"`
	SerialCode            string `json:"serial_code"`
	From                  int    `json:"from"`
	To                    int    `json:"to"`
	Sequence              int    `json:"sequence"`
	Status                int    `json:"status"`
	LastGeneratedTax      string `json:"last_generated_tax"`
	RemainingQty          int    `json:"remaining_qty"`
	TotalTaxNo            int    `json:"total_tax_no"`
	TaxNumberAlert        int    `json:"tax_number_alert"`
	UsedTotal             int    `json:"used_total"`
	DeletedTotal          int    `json:"deleted_total"`
	CreatedBy             int64  `json:"created_by"`
}

type DetailMTaxParams struct {
	MTaxID int64 `params:"m_taxes_id" validate:"required"`
}
type UpdateMTaxParams struct {
	MTaxID int64 `params:"m_taxes_id" validate:"required"`
}
type MTaxQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64  `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64  `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit" validate:""`
	Query        string  `query:"q"`
	Mode         string  `query:"mode"`
	Sort         string  `query:"sort"`
	IsActive     *int    `query:"is_active"`
	TrCode       string  `query:"tr_code"`
	Status       int     `query:"status"`
	CollectionNo *string `query:"collection_no"`
	Year         int     `query:"year"`
}

type MTaxesUpdateReq struct {
	CustID                string `json:"cust_id"`
	Year                  int    `json:"year" validate:"required"`
	TransactionStatusCode string `json:"transaction_status_code" validate:"required"`
	SerialCode            string `json:"serial_code" validate:"required"`
	From                  int    `json:"from"`
	To                    int    `json:"to" validate:"required,gtfield=From"`
	TaxNumberAlert        int    `json:"tax_number_alert"`
	UpdatedBy             int64  `json:"updated_by"`
	Status                int    `json:"status"`
}
type DeleteTaxesParams struct {
	TaxesID int64 `params:"taxes_id" validate:"required"`
}

type BulkDeleteTaxesParams struct {
	TaxesID []int64 `json:"taxes_id"`
}
