package entity

type SoQueryFilter struct {
	SalesmanId   []int `query:"salesman_id"`
	OutletID     []int `query:"outlet_id"`
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
