package entity

type StockReturnQueryFilter struct {
	CustId       string
	ParentCustId string
	SalesmanId   []int  `query:"salesman_id"`
	OutletID     []int  `query:"outlet_id"`
	Status       []int  `query:"status"`
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

const (
	IN_REVIEW   = 1
	NEED_REVIEW = 2
	PROCESSED   = 3
	IN_PICKUP   = 4
	PICKED_UP   = 5
	COMPLETED   = 6
	CANCELED    = 9
)

var dataReturnStatusName = map[int64]string{
	IN_REVIEW:   "In Review",
	NEED_REVIEW: "Need Review",
	PROCESSED:   "Processed",
	IN_PICKUP:   "In Pickup",
	PICKED_UP:   "Picked Up",
	COMPLETED:   "Completed",
	CANCELED:    "Canceled",
}

type StockReturnListResponse struct {
	RefferenceNo   *string `json:"refference_no"`
	ReturnNo       string  `json:"return_no"`
	InvoiceNo      *string `json:"invoice_no"`
	InvoiceDate    *string `json:"invoice_date"`
	SalesmanID     *int64  `json:"salesman_id"`
	SalesmanCode   *string `json:"salesman_code"`
	SalesmanName   *string `json:"salesman_name"`
	OutletID       *int64  `json:"outlet_id"`
	OutletCode     *string `json:"outlet_code"`
	OutletName     *string `json:"outlet_name"`
	DataStatus     *int64  `json:"data_status"`
	DataStatusName *string `json:"data_status_name"`
	CreatedBy      *int64  `json:"created_by"`
	CreatedByName  *string `json:"created_by_name"`
	CreatedAt      *string `json:"created_at"`
	ReviewedBy     *int64  `json:"reviewed_by"`
	ReviewedByName *string `json:"reviewed_by_name"`
	ReviewedAt     *string `json:"reviewed_at"`
}

func (rtn StockReturnListResponse) GenerateReturnStatusName() string {
	if rtn.DataStatus != nil {
		return dataReturnStatusName[*rtn.DataStatus]
	}
	return ""
}

type StockReturnResponse struct {
	RefferenceNo   *string                     `json:"refference_no"`
	ReturnNo       string                      `json:"return_no"`
	ReturnDate     *string                     `json:"return_date"`
	InvoiceNo      *string                     `json:"invoice_no"`
	InvoiceDate    *string                     `json:"invoice_date"`
	SalesmanID     *int64                      `json:"salesman_id"`
	SalesmanCode   *string                     `json:"salesman_code"`
	SalesmanName   *string                     `json:"salesman_name"`
	OutletID       *int64                      `json:"outlet_id"`
	OutletCode     *string                     `json:"outlet_code"`
	OutletName     *string                     `json:"outlet_name"`
	TprCashValue   *float64                    `json:"tpr_cash_value"`
	TprItemValue   *float64                    `json:"tpr_item_value"`
	Discount       *float64                    `json:"discount"`
	DiscountValue  *float64                    `json:"discount_value"`
	Vat            *float64                    `json:"vat"`
	VatValue       *float64                    `json:"vat_value"`
	SubTotal       *float64                    `json:"sub_total"`
	Total          *float64                    `json:"total"`
	DataStatus     *int64                      `json:"data_status"`
	DataStatusName *string                     `json:"data_status_name"`
	Details        []StockReturnDetailResponse `json:"details"`
}

func (rtn StockReturnResponse) GenerateReturnStatusName() string {
	if rtn.DataStatus != nil {
		return dataReturnStatusName[*rtn.DataStatus]
	}
	return ""
}

type StockDetailReturnParams struct {
	ReturnNo string `params:"return_no" validate:"required"`
}

type StockReturnUpdateBody struct {
	CustID     string                        `json:"cust_id"`
	DataStatus int                           `json:"data_status" validate:"required"`
	UpdatedBy  int64                         `json:"updated_by"`
	Details    []StockReturnDetailUpdateBody `json:"details"`
}

type UpdateReturnParams struct {
	ReturnNo string `params:"return_no" validate:"required"`
}

type StockReturnUpdateBatchBody struct {
	ReturnsNo  []string `json:"returns_no"`
	CustID     string   `json:"cust_id"`
	UpdatedBy  int64    `json:"updated_by"`
	DataStatus int      `json:"data_status" validate:"required"`
}
