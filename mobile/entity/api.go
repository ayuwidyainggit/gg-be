package entity

import (
	"fmt"
	"time"
)

type GeneralQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:""`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}
type Pagination struct {
	TotalRecord int64 `json:"total_record,omitempty"`
	PageCurrent int   `json:"page_current,omitempty"`
	PageLimit   int   `json:"page_limit,omitempty"`
	PageTotal   int   `json:"page_total,omitempty"`
}
type ApiResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	Paging    interface{} `json:"paging"`
	RequestId string      `json:"request_id"`
}

func GenerateNumber(initial string, seq int, roDate *time.Time) string {
	// Format tanggal menjadi yy, mm, dan dd
	yy := roDate.Format("06") // 2 digit tahun
	mm := roDate.Format("01") // 2 digit bulan
	dd := roDate.Format("02") // 2 digit hari
	// Format urutan menjadi 4 digit
	seqFormatted := fmt.Sprintf("%04d", seq+1)

	// Gabungkan semuanya untuk membuat nomor faktur
	roNumber := fmt.Sprintf("%s%s%s%s%s", initial, yy, mm, dd, seqFormatted)
	return roNumber
}

type PaginationWithTotalAmount struct {
	TotalRecord  int64   `json:"total_record,omitempty"`
	TotalInvoice float64 `json:"total_invoice,omitempty"`
	PageCurrent  int     `json:"page_current,omitempty"`
	PageLimit    int     `json:"page_limit,omitempty"`
	PageTotal    int     `json:"page_total,omitempty"`
}
