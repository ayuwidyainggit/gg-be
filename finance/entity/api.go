package entity

import (
	"fmt"
	"time"
)

type GeneralQueryFilter struct {
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
	CollectionNo *string `query:"collection_no"`
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

func ConvStatus(data map[int]string, param int) string {
	statusString, ok := data[int(param)]
	if !ok {
		statusString = "Unknown"
	}
	return statusString
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
