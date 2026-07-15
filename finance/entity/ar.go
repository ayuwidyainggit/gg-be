package entity

type CreateArBody struct {
	CustID     string            `json:"cust_id"`
	ArDate     *string           `json:"ar_date"`
	TrCode     *string           `json:"tr_code"`
	SalesmanID *int64            `json:"salesman_id"`
	RefNo      *string           `json:"ref_no"`
	DataStatus *int64            `json:"data_status"`
	CreatedBy  *int64            `json:"created_by"`
	IsPosted   *bool             `json:"is_posted"`
	Details    []CreateArDetBody `json:"details"`
}

type ArResponse struct {
	InvoiceNo         *string             `json:"invoice_no"`
	InvoiceDate       *string             `json:"invoice_date"`
	DueDate           *string             `json:"due_date"`
	InvoiceAmount     float64             `json:"invoice_amount"`
	PaidAmount        float64             `json:"paid_amount"`
	RemainingAmount   float64             `json:"remaining_amount"`
	SalesmanId        int64               `json:"salesman_id"`
	SalesmanCode      *string             `json:"salesman_code"`
	SalesmanName      *string             `json:"salesman_name"`
	OutletID          int64               `json:"outlet_id"`
	OutletCode        *string             `json:"outlet_code"`
	OutletName        *string             `json:"outlet_name"`
	InvoiceStatus     int64               `json:"invoice_status"`
	InvoiceStatusName string              `json:"invoice_status_name"`
	DueDateStatus     *int64              `json:"due_date_status"`
	DueDateStatusName *string             `json:"due_date_status_name"`
	Aging             *int64              `json:"aging"`
	Details           []ArPaymentResponse `json:"details"`
}

type ArListResponse struct {
	InvoiceNo         *string `json:"invoice_no"`
	InvoiceDate       *string `json:"invoice_date"`
	DueDate           *string `json:"due_date"`
	InvoiceAmount     float64 `json:"invoice_amount"`
	PaidAmount        float64 `json:"paid_amount"`
	RemainingAmount   float64 `json:"remaining_amount"`
	SalesmanId        int64   `json:"salesman_id"`
	SalesmanCode      *string `json:"salesman_code"`
	SalesmanName      *string `json:"salesman_name"`
	OutletID          int64   `json:"outlet_id"`
	OutletCode        *string `json:"outlet_code"`
	OutletName        *string `json:"outlet_name"`
	InvoiceStatus     int64   `json:"invoice_status"`
	InvoiceStatusName string  `json:"invoice_status_name"`
	DueDateStatus     *int64  `json:"due_date_status"`
	DueDateStatusName *string `json:"due_date_status_name"`
	Aging             *int64  `json:"aging"`
}

type UpdateArbody struct {
	CustID     string            `json:"cust_id"`
	ArDate     *string           `json:"ar_date"`
	TrCode     *string           `json:"tr_code"`
	SalesmanID *int64            `json:"salesman_id"`
	RefNo      *string           `json:"ref_no"`
	DataStatus *int64            `json:"data_status"`
	IsPosted   *bool             `json:"is_posted"`
	UpdatedBy  int64             `json:"updated_by"`
	Details    []UpdateArDetBody `json:"details"`
}

type DetailArParams struct {
	InvoiceNo string `params:"invoice_no" validate:"required" json:"invoice_no"`
}
type DeleteArParams struct {
	ArNo string `params:"ar_no" validate:"required" json:"ar_no"`
}
type UpdateArParams struct {
	ArNo string `params:"ar_no" validate:"required" json:"ar_no"`
}

type CreateCollectionBody struct {
	CustID          string                    `json:"cust_id"`
	CollectionNo    *string                   `json:"collection_no"`
	CollectionDate  *string                   `json:"collection_date"`
	EmpID           *int64                    `json:"emp_id"`
	OtGrpID         *int64                    `json:"ot_grp_id"`
	Notes           *string                   `json:"notes"`
	TotalAmount     *float64                  `json:"total_amount"`
	RemainingAmount *float64                  `json:"remaining_amount"`
	InvoiceDateFrom *string                   `json:"invoice_date_from"`
	InvoiceDateTo   *string                   `json:"invoice_date_to"`
	DueDateFrom     *string                   `json:"due_date_from"`
	DueDateTo       *string                   `json:"due_date_to"`
	CreatedBy       *int64                    `json:"created_by"`
	CreatedAt       *string                   `json:"created_at"`
	Details         []CreateCollectionDetBody `json:"details"`
}

type DetailCollectionParams struct {
	CollectionNo string `params:"collection_no" validate:"required" json:"collection_no"`
}

type DeleteCollectionParams struct {
	CollectionNo string `params:"collection_no" validate:"required" json:"collection_no"`
}
type UpdateCollectionParams struct {
	CollectionNo string `params:"collection_no" validate:"required" json:"collection_no"`
}
type PrintCollectionParams struct {
	CollectionNo string `params:"collection_no" validate:"required" json:"collection_no"`
}

type CollectionResponse struct {
	CustID              string                  `json:"cust_id"`
	CollectionNo        *string                 `json:"collection_no"`
	CollectionDate      *string                 `json:"collection_date"`
	EmpID               *int64                  `json:"emp_id"`
	EmpCode             *string                 `json:"emp_code"`
	EmpName             *string                 `json:"emp_name"`
	OtGrpID             *int64                  `json:"ot_grp_id"`
	OtGrpCode           *string                 `json:"ot_grp_code"`
	OtGrpName           *string                 `json:"ot_grp_name"`
	Notes               *string                 `json:"notes"`
	TotalAmount         *float64                `json:"total_amount"`
	TotalInvoicePayment *float64                `json:"total_invoice_payment"`
	RemainingAmount     *float64                `json:"remaining_amount"`
	InvoiceDateFrom     *string                 `json:"invoice_date_from"`
	InvoiceDateTo       *string                 `json:"invoice_date_to"`
	DueDateFrom         *string                 `json:"due_date_from"`
	DueDateTo           *string                 `json:"due_date_to"`
	CreatedBy           *int64                  `json:"created_by"`
	CreatedByName       *string                 `json:"created_by_name"`
	CreatedAt           *string                 `json:"created_at"`
	UpdatedBy           *int64                  `json:"updated_by"`
	UpdatedByName       *string                 `json:"updated_by_name"`
	UpdatedAt           *string                 `json:"updated_at"`
	DeletedBy           *int64                  `json:"deleted_by"`
	DeletedByName       *string                 `json:"deleted_by_name"`
	DeletedAt           *string                 `json:"deleted_at"`
	PrintedBy           *int64                  `json:"printed_by"`
	PrintedByName       *string                 `json:"printed_by_name"`
	PrintedAt           *string                 `json:"printed_at"`
	Details             []CollectionDetResponse `json:"details"`
}
type CollectionListResponse struct {
	CollectionNo    *string  `json:"collection_no"`
	CollectionDate  *string  `json:"collection_date"`
	EmpID           *int64   `json:"emp_id"`
	EmpCode         *string  `json:"emp_code"`
	EmpName         *string  `json:"emp_name"`
	EmpGrpID        *int64   `json:"emp_grp_id"`
	EmpGrpCode      *string  `json:"emp_grp_code"`
	EmpGrpName      *string  `json:"emp_grp_name"`
	RemainingAmount *float64 `json:"remaining_amount"`
}
type UpdateCollectionBody struct {
	CustID          string                    `json:"cust_id"`
	CollectionDate  *string                   `json:"collection_date"`
	UpdatedBy       int64                     `json:"updated_by"`
	EmpID           *int64                    `json:"emp_id"`
	OtGrpID         *int64                    `json:"ot_grp_id"`
	TotalAmount     *float64                  `json:"total_amount"`
	RemainingAmount *float64                  `json:"remaining_amount"`
	InvoiceDateFrom *string                   `json:"invoice_date_from"`
	InvoiceDateTo   *string                   `json:"invoice_date_to"`
	DueDateFrom     *string                   `json:"due_date_from"`
	DueDateTo       *string                   `json:"due_date_to"`
	Notes           *string                   `json:"notes"`
	Details         []UpdateCollectionDetBody `json:"details"`
}

type CollectionQueryFilter struct {
	CustId             string
	ParentCustId       string
	EmpId              []int  `query:"emp_id"`
	CollectionDateFrom *int64 `query:"collection_date_from" validate:"required_with=CollectionDateTo,omitempty,gte=1000000000"`
	CollectionDateTo   *int64 `query:"collection_date_to" validate:"required_with=CollectionDateFrom,omitempty,lte=9999999999,gtefield=CollectionDateFrom"`
	Page               int    `query:"page"`
	Limit              int    `query:"limit" validate:"required"`
	Query              string `query:"q"`
	Mode               string `query:"mode"`
	Sort               string `query:"sort"`
	IsActive           *int   `query:"is_active"`
}

type EmployeeGroupLookupResponse struct {
	EmpGroupId   int    `json:"emp_grp_id"`
	EmpGroupCode string `json:"emp_grp_code"`
	EmpGroupName string `json:"emp_grp_name"`
}

type OutletGroupLookupResponse struct {
	OutletGroupId   int    `json:"ot_grp_id"`
	OutletGroupCode string `json:"ot_grp_code"`
	OutletGroupName string `json:"ot_grp_name"`
}

type EmployeeListQueryFilter struct {
	EmpGrpID     []int `query:"emp_grp_id"`
	InvoiceNo    []int `query:"inv_no"`
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
	TrCode       string `query:"tr_code"`
}

type EmployeeLookupResponse struct {
	EmpId   int    `json:"emp_id"`
	EmpCode string `json:"emp_code"`
	EmpName string `json:"emp_name"`
}

type SalesmanListQueryFilter struct {
	InvoiceNo       []int `query:"inv_no"`
	CustId          string
	ParentCustId    string
	From            *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To              *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	InvoiceDateFrom *int64 `query:"invoice_date_from" validate:"required_with=InvoiceDateTo,omitempty,gte=1000000000"`
	InvoiceDateTo   *int64 `query:"invoice_date_to" validate:"required_with=InvoiceDateFrom,omitempty,lte=9999999999,gtefield=InvoiceDateFrom"`
	DueDateFrom     *int64 `query:"due_date_from" validate:"required_with=DueDateTo,omitempty,gte=1000000000"`
	DueDateTo       *int64 `query:"due_date_to" validate:"required_with=DueDateFrom,omitempty,lte=9999999999,gtefield=DueDateFrom"`
	Page            int    `query:"page"`
	Limit           int    `query:"limit" validate:""`
	Query           string `query:"q"`
	Mode            string `query:"mode"`
	Sort            string `query:"sort"`
	IsActive        *int   `query:"is_active"`
	TrCode          string `query:"tr_code"`
}

type SalesmanLookupResponse struct {
	SalesmanId   int    `json:"salesman_id"`
	SalesmanCode string `json:"salesman_code"`
	SalesmanName string `json:"salesman_name"`
}

type InvoiceQueryFilter struct {
	SalesmanId    []int  `query:"salesman_id"`
	OutletID      []int  `query:"outlet_id"`
	OutletGroupID []int  `query:"ot_grp_id"`
	InvoiceStatus *int64 `query:"invoice_status" validate:"omitempty,oneof=1 2"`
	// Status        []int `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	RoFrom       *int64 `query:"ro_date_from" validate:"required_with=RoTo,omitempty,gte=1000000000"`
	RoTo         *int64 `query:"ro_date_to" validate:"required_with=RoFrom,omitempty,lte=9999999999,gtefield=RoFrom"`
	InvoiceFrom  *int64 `query:"inv_date_from" validate:"required_with=InvoiceTo,omitempty,gte=1000000000"`
	InvoiceTo    *int64 `query:"inv_date_to" validate:"required_with=InvoiceFrom,omitempty,lte=9999999999,gtefield=InvoiceFrom"`
	DueFrom      *int64 `query:"due_date_from" validate:"required_with=DueTo,omitempty,gte=1000000000"`
	DueTo        *int64 `query:"due_date_to" validate:"required_with=DueFrom,omitempty,lte=9999999999,gtefield=DueFrom"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type InvoiceListResponse struct {
	InvoiceNo       *string  `json:"invoice_no"`
	InvoiceDate     *string  `json:"invoice_date"`
	RoNo            *string  `json:"ro_no"`
	DueDate         *string  `json:"due_date"`
	OutletID        *int64   `json:"outlet_id"`
	OutletCode      *string  `json:"outlet_code"`
	OutletName      *string  `json:"outlet_name"`
	SalesmanId      *int64   `json:"salesman_id"`
	SalesmanCode    *string  `json:"salesman_code"`
	SalesmanName    *string  `json:"salesman_name"`
	InvoiceAmount   *float64 `json:"invoice_amount"`
	RemainingAmount *float64 `json:"remaining_amount"`
	PaidAmount      *float64 `json:"paid_amount"`
}

type OutletLookupResponse struct {
	OutletId   int    `json:"outlet_id"`
	OutletCode string `json:"outlet_code"`
	OutletName string `json:"outlet_name"`
}

type InvoiceStatus int64

const (
	InvoiceStatusPaid        InvoiceStatus = 1
	InvoiceStatusOutstanding InvoiceStatus = 2
)

var dataInvoiceStatusName = map[InvoiceStatus]string{
	InvoiceStatusPaid:        "Paid",
	InvoiceStatusOutstanding: "Outstanding",
}

var dataDueDateStatusName = map[int64]string{
	1: "On Schedule",
	2: "Overdue",
}

func (ar ArResponse) GenerateDataInvoiceStatusName() string {
	return dataInvoiceStatusName[InvoiceStatus(ar.InvoiceStatus)]
}

func (ar ArResponse) GenerateDataDueDateStatusName() string {
	if ar.DueDateStatus != nil {
		return dataDueDateStatusName[*ar.DueDateStatus]
	}

	return ""
}

func (ar ArListResponse) GenerateDataInvoiceStatusName() string {
	return dataInvoiceStatusName[InvoiceStatus(ar.InvoiceStatus)]
}

func (ar ArListResponse) GenerateDataDueDateStatusName() string {
	if ar.DueDateStatus != nil {
		return dataDueDateStatusName[*ar.DueDateStatus]
	}

	return ""
}
