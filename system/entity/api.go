package entity

type GeneralQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	IsActive *int   `query:"is_active"`
}
type Pagination struct {
	TotalRecord int64 `json:"total_record,omitempty"`
	PageCurrent int   `json:"page_current,omitempty"`
	PageLimit   int   `json:"page_limit,omitempty"`
	PageTotal   int   `json:"page_total,omitempty"`
}
type ApiResponse struct {
	Message   string
	Data      interface{}
	Errors    interface{}
	Paging    interface{}
	RequestId string
}
