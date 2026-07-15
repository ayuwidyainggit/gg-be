package entity

type GeneralQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	TrCode       string `query:"tr_code"`
	GrType       *int   `query:"gr_type"`
}

type Pagination struct {
	TotalRecord int64 `json:"total_record"`
	PageCurrent int   `json:"page_current"`
	PageLimit   int   `json:"page_limit"`
	PageTotal   int   `json:"page_total"`
}

type ApiResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Paging    interface{} `json:"paging,omitempty"`
	RequestId string      `json:"request_id"`
}
