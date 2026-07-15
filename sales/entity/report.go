package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"
)

const (
	TYPE_REPORT_SALESMAN_ACTIVITY_REPORT = "SalesmanActivityReport"
	REPORT_NAME_SECONDARY_SALES          = "SecondarySales"
	REPORT_NAME_DOWNLOAD_SALES_ORDER     = "DownloadSalesOrder"

	SECONDARY_SALES_GROUP_OUTLET           = "outlet"
	SECONDARY_SALES_GROUP_SALESMAN         = "salesman"
	SECONDARY_SALES_GROUP_PRODUCT_CATEGORY = "product_category"
	SECONDARY_SALES_GROUP_PRODUCT          = "product"

	ACTIVITY_SALESMAN_GROUP_SALES  = "sales"
	ACTIVITY_SALESMAN_GROUP_RETURN = "return"
)

type ReportQueryFilter struct {
	ReportType   []string `query:"report_type"`
	FileStatus   []int    `query:"file_status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
}

type SecondarySalesReportQueryFilter struct {
	// CustID and ParentCustID are set from JWT locals only; json:"-" prevents body spoofing
	// but they ARE included in RMQ serialization via explicit field names below.
	CustID           string             `json:"_cust_id"`
	ParentCustID     string             `json:"_parent_cust_id"`
	RequestedCustID  string             `json:"-" validate:"omitempty,alphanum,max=20"`
	RequestedCustIDs StringListOrScalar `json:"-"`
	CustIDs          []string           `json:"cust_ids,omitempty"`
	From             *int64             `json:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To               *int64             `json:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Sort             string             `json:"sort"`
	Page             int                `json:"page"`
	Limit            int                `json:"limit"`
	ExportBy         string
	DistributorIDs   []int64 `json:"distributor_ids"`
	SalesmanIDs      []int64 `json:"salesman_ids"`
	OutletIDs        []int64 `json:"outlet_ids"`
	ProIDs           []int64 `json:"pro_ids"`
	ExportDate       string  `query:"export_date"`
	ReportID         string  `json:"report_id"`
}

type SecondarySalesReportResponse struct {
	DistributorCode string  `json:"distributor_code"`
	DistributorName string  `json:"distributor_name"`
	TrxType         string  `json:"trx_type"` // <-- tambahan
	InvoiceNo       string  `json:"invoice_no"`
	InvoiceDate     string  `json:"invoice_date"`
	DocumentNo      string  `json:"document_no"`
	DocumentDate    string  `json:"document_date"`
	OutletID        int64   `json:"outlet_id"`
	OutletCode      string  `json:"outlet_code"`
	OutletName      string  `json:"outlet_name"`
	SalesmanID      int64   `json:"salesman_id"`
	EmpCode         string  `json:"emp_code"`
	EmpName         string  `json:"emp_name"`
	ProID           int64   `json:"pro_id"`
	SupCode         string  `json:"sup_code"`
	SupName         string  `json:"sup_name"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	UnitID1         string  `json:"unit_id1"`
	UnitID2         string  `json:"unit_id2"`
	UnitID3         string  `json:"unit_id3"`
	ConvUnit2       float64 `json:"conv_unit2"`
	ConvUnit3       float64 `json:"conv_unit3"`
	Qty1            float64 `json:"qty1"`
	Qty2            float64 `json:"qty2"`
	Qty3            float64 `json:"qty3"`
	GrossSales      float64 `json:"gross_sales"`
	SpecialDiscount float64 `json:"special_discount"`
	Discount        float64 `json:"discount"`
	NetSalesExcPPN  float64 `json:"net_sales_exc_ppn"`
	PPN             float64 `json:"ppn"`
	NetSalesIncPPN  float64 `json:"net_sales_inc_ppn"`
	Qty1Return      float64 `json:"qty1_return"`
	Qty2Return      float64 `json:"qty2_return"`
	Qty3Return      float64 `json:"qty3_return"`
	SellPrice1      float64 `json:"sell_price1"`
	SellPrice2      float64 `json:"sell_price2"`
	SellPrice3      float64 `json:"sell_price3"`
}
type DetailReportSecondarySalesParams struct {
	InvoiceNo string `params:"invoice_no" validate:"required"`
}

type ReportList struct {
	ReportID       string    `json:"report_id"`
	ReportName     string    `json:"report_name"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	FileStatus     int       `json:"file_status"`
	FileStatusName string    `json:"file_status_name"`
	FileURL        string    `json:"file_url"`
	FileBase64     string    `json:"file_base64,omitempty"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

var FileStatusName = map[int]string{
	1: "Ready",
	2: "Processing",
	3: "Failed",
	4: "Expired",
}

const (
	FILE_STATUS_READY      = 1
	FILE_STATUS_PROCESSING = 2
	FILE_STATUS_FAILED     = 3
	FILE_STATUS_EXPIRED    = 4
)

func (report ReportList) GetFileStatusName() string {
	return FileStatusName[report.FileStatus]
}

type PublishByRmqSecondarySalesReq struct {
	PriceID      string `json:"price_id"`
	CustID       string `json:"cust_id"`
	ParentCustID string `json:"parent_cust_id"`
	Status       int    `json:"status"`
	UpdatedBy    string `json:"updated_by"`
}

type ActivityReportQueryFilter struct {
	CustID           string   `json:"cust_id,omitempty"`
	RequestedCustID  string   `json:"requested_cust_id,omitempty"`
	RequestedCustIDs []string `json:"requested_cust_ids,omitempty"`
	CustIDs          []string `json:"cust_ids,omitempty"`
	ParentCustID     string   `json:"parent_cust_id"`
	AuthCustID       string   `json:"auth_cust_id,omitempty"`
	IsAdmin         bool   `json:"-"`
	DistPriceGrpID  int    `json:"-"`
	SalesmanIDs      []int    `json:"salesman_ids"`
	DistributorCodes []string `json:"distributor_code,omitempty"`
	FromDate         string   `json:"from" validate:"required"`
	ToDate            string   `json:"to" validate:"required"`
	Sort              string   `json:"sort"`
	Page              int      `json:"page"`
	Limit             int      `json:"limit"`
	ExportBy          string
	ExportDate      string `query:"export_date"`
	ReportID        string `json:"report_id"`
}

type ActivityReportSalesmanListQueryFilter struct {
	CustID       string
	FromDate     string `query:"from" validate:"required"`
	ToDate       string `query:"to" validate:"required"`
	SalesmanName string `query:"salesman_name"`
}

type ActivityReportQueryFilterList struct {
	CustID          string
	AuthCustID      string
	CustIDs         []string
	RequestedCustID string `query:"cust_id"`
	ParentCustID    string
	IsAdmin         bool
	DistPriceGrpID  int
	SalesmanIDs      []int  `query:"salesman_ids"`
	DistributorCode  string `query:"distributor_code"`
	DistributorCodes []string
	FromDate         string `query:"from" validate:"required"`
	ToDate           string `query:"to" validate:"required"`
	Sort             string `query:"sort"`
	Page             int    `query:"page"`
	Limit            int    `query:"limit"`
}

type ActivityReportListResp struct {
	BusinessUnitCode    string  `json:"business_unit_code"`
	BusinessUnitName    string  `json:"business_unit_name"`
	DistributorCode     string  `json:"distributor_code"`
	DistributorName     string  `json:"distributor_name"`
	PJPCode             string  `json:"pjp_code"`
	EmployeeCode        string  `json:"employee_code"`
	SalesmanName        string  `json:"salesman_name"`
	OutletCode          string  `json:"outlet_code"`
	OutletPrincipalCode string  `json:"outlet_principal_code"`
	OutletName          string  `json:"outlet_name"`
	Date                string  `json:"date"`
	ClockInTime         string  `json:"clock_in_time"`
	ClockOutTime        string  `json:"clock_out_time"`
	CheckinTime         string  `json:"check_in_time"`
	CheckoutTime        string  `json:"check_out_time"`
	Duration            string  `json:"duration"`
	PjpStatus           string  `json:"pjp_status"`
	VisitStatus         string  `json:"visit_status"`
	Compliance          string  `json:"compliance"`
	SalesValue          float64 `json:"sales_value"`
	ReturnValue         float64 `json:"return_value"`
	PaymentCollected    float64 `json:"payment_collected"`
	LocationMaster      string  `json:"location_master"`
	LocationActual      string  `json:"location_actual"`
	GeotagStatus        *int    `json:"geotag_status"`
	GeotagStatusDesc    string  `json:"geotag_status_desc"`
	Remarks             string  `json:"remarks"`
}

type SecondarySalesReportDashboardExtractQueryFilter struct {
	Date time.Time
}

type SecondarySalesReportExtractPayload struct {
	Day   int `json:"day" validate:"required"`
	Month int `json:"month" validate:"required"`
	Year  int `json:"year" validate:"required"`
}

type SecondarySalesReportDashboardSumPayload struct {
	Month       int      `query:"month" validate:"required,gte=1,lte=12"`
	Year        *int     `query:"year" validate:"omitempty,gte=2000,lte=9999"`
	CustID      string   `query:"cust_id" validate:"omitempty"`
	CustIDs     []string `json:"-"`
	From        *int64   `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To          *int64   `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	OutletIDs   []int64  `query:"outlet_ids"`
	SalesmanIDs []int64  `query:"salesman_ids"`
	ProIDs      []int64  `query:"pro_ids"`
}

type SecondarySalesReportTrensSalesSumPayload struct {
	Year    int      `query:"year" validate:"required"`
	CustID  string   `json:"cust_id,omitempty" validate:"omitempty,alphanum,max=20"`
	CustIDs []string `json:"-"`
}

type SecondarySalesReportDashboardGroupPayload struct {
	Month   int      `query:"month" validate:"required,gte=1,lte=12"`
	Year    *int     `query:"year" validate:"omitempty,gte=2000,lte=9999"`
	CustID  string   `query:"cust_id" validate:"omitempty"`
	CustIDs []string `json:"-"`
	GroupBy string   `query:"group_by" validate:"omitempty"`
}
type SumReportByMonthModelResp struct {
	TotalGrossSales    float64    `json:"total_gross_sale"`
	TotalDiscountPromo float64    `json:"total_discount_promo"`
	TotalPPN           float64    `json:"total_ppn"`
	NetSalesExcPPN     float64    `json:"net_sales_exc_ppn"`
	NetSales           float64    `json:"net_sales"`
	TotalSalesman      int32      `json:"total_salesman"`
	TotalOutlet        int32      `json:"total_outlet"`
	TotalProduct       int32      `json:"total_product"`
	Qty                int64      `json:"qty"`
	QtyReturn          int64      `json:"qty_return"`
	ReturnRate         float64    `json:"return_rate"`
	NetSalesReturn     float64    `json:"net_sales_return"`
	LastUpdate         *time.Time `json:"last_update"`
}

type SecondarySalesReportGroupResp struct {
	ID       int     `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	NetSales float64 `json:"net_sales"`
}

type SumReportTrendSalesResp struct {
	Month              int     `json:"month"`
	TotalGrossSales    float64 `json:"total_gross_sale"`
	TotalDiscountPromo float64 `json:"total_discount_promo"`
	NetSales           float64 `json:"net_sales"`
}

type StringListOrScalar []string

func (s *StringListOrScalar) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*s = nil
		return nil
	}

	if len(trimmed) > 0 && trimmed[0] == '[' {
		var values []string
		if err := json.Unmarshal(trimmed, &values); err != nil {
			return err
		}
		normalized, err := NormalizeStringList(values)
		if err != nil {
			return err
		}
		*s = normalized
		return nil
	}

	var value string
	if err := json.Unmarshal(trimmed, &value); err != nil {
		return err
	}
	normalized, err := NormalizeStringList([]string{value})
	if err != nil {
		return err
	}
	*s = normalized
	return nil
}

func NormalizeStringList(raw []string) ([]string, error) {
	result := make([]string, 0)
	seen := make(map[string]struct{})

	for _, value := range raw {
		parts := strings.Split(value, ",")
		for _, part := range parts {
			cleaned := strings.TrimSpace(part)
			if cleaned == "" {
				continue
			}
			for _, r := range cleaned {
				if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
					return nil, fmt.Errorf("invalid cust_id value %q", cleaned)
				}
			}
			if _, exists := seen[cleaned]; exists {
				continue
			}
			seen[cleaned] = struct{}{}
			result = append(result, cleaned)
		}
	}

	return result, nil
}

func NormalizeDistributorCodeList(raw []string) ([]string, error) {
	result := make([]string, 0)
	seen := make(map[string]struct{})

	for _, value := range raw {
		parts := strings.Split(value, ",")
		for _, part := range parts {
			cleaned := strings.TrimSpace(part)
			if cleaned == "" {
				continue
			}
			if _, exists := seen[cleaned]; exists {
				continue
			}
			seen[cleaned] = struct{}{}
			result = append(result, cleaned)
		}
	}

	return result, nil
}

type ActivityReportTrendSalesPayload struct {
	Year    int      `query:"year" validate:"required,gte=2000,lte=9999"`
	CustID  string   `query:"cust_id" validate:"omitempty"`
	CustIDs []string `json:"-"`
}

type ActivityReportTrendSalesResp struct {
	Month        int     `json:"month"`
	TotalInvoice float64 `json:"total_invoice"`
	TotalReturn  float64 `json:"total_return"`
	NetSales     float64 `json:"net_sales"`
}

type SalesmanActivityReportDashboardSumPayload struct {
	Month   int      `query:"month"`
	Year    *int     `query:"year" validate:"omitempty,gte=2000,lte=9999"`
	CustID  string   `query:"cust_id" validate:"omitempty"`
	CustIDs []string `json:"-"`
}
type SalesmanActivityReportByMonthModelResp struct {
	TotalSales    float64    `json:"total_sales"`
	TotalReturn   float64    `json:"total_return"`
	SalesmanTotal int32      `json:"salesman_total"`
	LastUpdate    *time.Time `json:"last_update"`
}

type SalesmanActivityReportDashboardGroupPayload struct {
	Month        int      `query:"month"`
	Year         *int     `query:"year" validate:"omitempty,gte=2000,lte=9999"`
	ActivityType string   `query:"activity_type"`
	CustID       string   `query:"cust_id" validate:"omitempty"`
	CustIDs      []string `json:"-"`
}

type SalesmanActivityReportSalesmanListResp struct {
	SalesmanID   int64  `json:"salesman_id"`
	SalesmanCode string `json:"salesman_code"`
	SalesmanName string `json:"salesman_name"`
}

type ActivityReportGeotagPayload struct {
	Year    int      `query:"year" validate:"required,gte=2000,lte=9999"`
	CustID  string   `query:"cust_id" validate:"omitempty"`
	CustIDs []string `json:"-"`
	EmpID   *int     `query:"emp_id" validate:"omitempty"`
}

type ActivityReportGeotagResp struct {
	TotalGeotagMatchPercentage   float64                        `json:"total_geotag_match_percentage"`
	TotalGeotagUnmatchPercentage float64                        `json:"total_geotag_unmatch_percentage"`
	Details                      []ActivityReportGeotagDetailResp `json:"details"`
}

type ActivityReportGeotagDetailResp struct {
	SalesmanCode            string  `json:"salesman_code"`
	SalesmanName            string  `json:"salesman_name"`
	TotalVisit              int64   `json:"total_visit"`
	GeotagMatchCount        int64   `json:"geotag_match_count"`
	GeotagUnmatchCount      int64   `json:"geotag_unmatch_count"`
	GeotagMatchPercentage   float64 `json:"geotag_match_percentage"`
	GeotagUnmatchPercentage float64 `json:"geotag_unmatch_percentage"`
}
