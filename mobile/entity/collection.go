package entity

import "time"

type CollectionQueryFilter struct {
	SalesmanId   []int `query:"salesman_id"`
	OutletID     []int `query:"outlet_id"`
	Status       []int `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	RoFrom       *int64 `query:"ro_date_from" validate:"required_with=RoTo,omitempty,gte=1000000000"`
	RoTo         *int64 `query:"ro_date_to" validate:"required_with=RoFrom,omitempty,lte=9999999999,gtefield=RoFrom"`
	InvoiceFrom  *int64 `query:"inv_date_from" validate:"required_with=InvoiceTo,omitempty,gte=1000000000"`
	InvoiceTo    *int64 `query:"inv_date_to" validate:"required_with=InvoiceFrom,omitempty,lte=9999999999,gtefield=InvoiceFrom"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type NoCollectionQueryFilter struct {
	SalesmanId   []int `query:"salesman_id"`
	OutletID     []int `query:"outlet_id"`
	Status       []int `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	RoFrom       *int64 `query:"ro_date_from" validate:"required_with=RoTo,omitempty,gte=1000000000"`
	RoTo         *int64 `query:"ro_date_to" validate:"required_with=RoFrom,omitempty,lte=9999999999,gtefield=RoFrom"`
	InvoiceFrom  *int64 `query:"inv_date_from" validate:"required_with=InvoiceTo,omitempty,gte=1000000000"`
	InvoiceTo    *int64 `query:"inv_date_to" validate:"required_with=InvoiceFrom,omitempty,lte=9999999999,gtefield=InvoiceFrom"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type CreateCollectionBody struct {
	CustId        string   `json:"cust_id"`
	RoDate        *string  `json:"ro_date"`
	ValDate       *string  `json:"val_date"`
	DueDate       *string  `json:"due_date"`
	SalesmanId    int64    `json:"salesman_id" validate:"required"`
	WhId          *int64   `json:"wh_id"`
	OutletID      int64    `json:"outlet_id" validate:"required"`
	DeliveryDate  *string  `json:"delivery_date"`
	CollectionNo  *string  `json:"order_no"`
	PoNo          *string  `json:"po_no"`
	VehicleNo     *string  `json:"vehicle_no"`
	PayType       *int64   `json:"pay_type"`
	ReffNo        *string  `json:"reff_no"`
	MobileID      *int64   `json:"mobile_id"`
	SubTotal      *float64 `json:"sub_total"`
	Disc          *float64 `json:"disc"`
	DiscValue     *float64 `json:"disc_value"`
	PromoValue    *float64 `json:"promo_value"`
	CashDiscValue *float64 `json:"cash_disc_value"`
	TotDisc1      *float64 `json:"tot_disc1"`
	TotDisc2      *float64 `json:"tot_disc2"`
	Vat           *float64 `json:"vat"`
	VatValue      *float64 `json:"vat_value"`
	Total         *float64 `json:"total"`
	DataStatus    *int64   `json:"data_status"`
	CreatedBy     *int64   `json:"created_by"`
	DataSource    *int64   `json:"data_source"`
	// Details       CollectionDetWithGroup `json:"details"`
	TrCode      *string `json:"tr_code"`
	IsClosed    bool    `json:"is_closed"`
	Notes       *string `json:"notes"`
	InvoiceNo   *string `json:"invoice_no"`
	InvoiceDate *string `json:"invoice_date"`
}

type CreateNoCollectionBody struct {
	CustId             string  `json:"cust_id"`
	SalesmanId         int64   `json:"salesman_id" validate:"required"`
	NoCollectionDate   *string `json:"no_order_date"`
	OutletId           int64   `json:"outlet_id" validate:"required"`
	TakingCollectionId int64   `json:"taking_order_id" validate:"required"`
	Reason             *string `json:"reason"`
	CreatedBy          int64   `json:"created_by"`
	CreatedAt          *string `json:"created_at"`
}

type DetailCollectionParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}
type DeleteCollectionParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type UpdateCollectionParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type CollectionResponse struct {
	CustID          string                  `json:"cust_id"`
	CollectionNo    *string                 `json:"collection_no"`
	CollectionDate  *string                 `json:"collection_date"`
	EmpID           *int64                  `json:"emp_id"`
	EmpCode         *string                 `json:"emp_code"`
	EmpName         *string                 `json:"emp_name"`
	OtGrpID         *int64                  `json:"ot_grp_id"`
	OtGrpCode       *string                 `json:"ot_grp_code"`
	OtGrpName       *string                 `json:"ot_grp_name"`
	Notes           *string                 `json:"notes"`
	TotalAmount     *float64                `json:"total_amount"`
	RemainingAmount *float64                `json:"remaining_amount"`
	InvoiceDateFrom *string                 `json:"invoice_date_from"`
	InvoiceDateTo   *string                 `json:"invoice_date_to"`
	DueDateFrom     *string                 `json:"due_date_from"`
	DueDateTo       *string                 `json:"due_date_to"`
	CreatedBy       *int64                  `json:"created_by"`
	CreatedByName   *string                 `json:"created_by_name"`
	CreatedAt       *string                 `json:"created_at"`
	UpdatedBy       *int64                  `json:"updated_by"`
	UpdatedByName   *string                 `json:"updated_by_name"`
	UpdatedAt       *string                 `json:"updated_at"`
	DeletedBy       *int64                  `json:"deleted_by"`
	DeletedByName   *string                 `json:"deleted_by_name"`
	DeletedAt       *string                 `json:"deleted_at"`
	PrintedBy       *int64                  `json:"printed_by"`
	PrintedByName   *string                 `json:"printed_by_name"`
	PrintedAt       *string                 `json:"printed_at"`
	Details         []CollectionDetResponse `json:"details"`
}

type CollectionListResponse struct {
	CollectionNo   *string `json:"collection_no"`
	CollectionDate *string `json:"collection_date"`
	InvoiceNo      string  `json:"invoice_no"`
	EmpID          *int64  `json:"emp_id"`
	EmpName        *string `json:"sales_name"`
	// EmpCode         *string  `json:"emp_code"`
	// EmpGrpID        *int64   `json:"emp_grp_id"`
	// EmpGrpCode      *string  `json:"emp_grp_code"`
	// EmpGrpName      *string  `json:"emp_grp_name"`
	RemainingAmount *float64 `json:"remaining_amount"`
}

type CollectionListV2Response struct {
	CollectionNo    string  `json:"collection_no"`
	EmpID           *int64  `json:"emp_id"`
	InvoiceNo       *string `json:"invoice_no"`
	InvoiceAmount   float64 `json:"invoice_amount"`
	PaidAmount      float64 `json:"paid_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	InvoiceDate     *string `json:"invoice_date"`
	DueDate         *string `json:"due_date"`
	RoNo            string  `json:"ro_no"`
	OrderNo         string  `json:"order_no"`
}

type CollectionDetResponse struct {
	CollectionDetID     int64    `json:"collection_det_id"`
	CollectionNo        string   `json:"collection_no"`
	InvoiceNo           string   `json:"invoice_no"`
	SalesOrder          string   `json:"sales_order"`
	InvoiceDate         string   `json:"invoice_date"`
	DueDate             string   `json:"due_date"`
	SalesmanId          int64    `json:"salesman_id"`
	SalesmanName        string   `json:"salesman_name"`
	SalesmanCode        string   `json:"salesman_code"`
	OutletId            int64    `json:"outlet_id"`
	OutletCode          string   `json:"outlet_code"`
	OutletName          string   `json:"outlet_name"`
	InvoiceAmount       *float64 `json:"invoice_amount"`
	RemainingAmount     *float64 `json:"remaining_amount"`
	PaidAmount          *float64 `json:"paid_amount"`
	TotalInvoicePayment *float64 `json:"total_invoice_amount"`
	CreatedBy           *int64   `json:"created_by"`
	CreatedByName       *string  `json:"created_by_name"`
	CreatedAt           *string  `json:"created_at"`
}

type CreateDepositBodyByCollection struct {
	CustID              string          `json:"cust_id"`
	DepositNo           string          `json:"deposit_no"`
	DepositDate         string          `json:"deposit_date"`
	CollectionNo        *string         `json:"collection_no"`
	EmpGrpID            *int            `json:"emp_grp_id"`
	EmpID               *int            `json:"emp_id"`
	DepositStatus       int             `json:"deposit_status"`
	RemainingAmount     float64         `json:"remaining_amount"`
	TotalDiscount       float64         `json:"total_discount"`
	TotalMaterai        float64         `json:"total_materai"`
	TotalPaymentBalance float64         `json:"total_payment_balance"`
	TotalPayment        float64         `json:"total_payment"`
	Details             []DepositDetail `json:"detail" validate:"required,dive,required"`
	CreatedBy           *int64          `json:"created_by"`
	CreatedAt           *time.Time      `json:"created_at"`
	UpdatedBy           *int64          `json:"updated_by"`
	UpdatedAt           *time.Time      `json:"updated_at"`
}

type CreateCollectionNoPaymentRequest struct {
	CustID                 string  `json:"cust_id"`
	SalesmanID             *int64  `json:"salesman_id,omitempty"`
	OutletID               *int64  `json:"outlet_id,omitempty"`
	CollectionNo           string  `json:"collection_no"`
	InvoiceNo              string  `json:"invoice_no"`
	MissedPaymentReasonsID *int64  `json:"missed_payment_reasons_id,omitempty"`
	Reason                 *string `json:"reason,omitempty"`
	PaymentDate            string  `json:"payment_date,omitempty"`

	CreatedBy *int64     `json:"created_by"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DepositDetail struct {
	DepositDetailID  int64            `json:"deposit_detail_id"`
	DepositNo        string           `json:"deposit_no"`
	InvoiceNo        string           `json:"invoice_no"`
	Discount         float64          `json:"discount"`
	PaymentBalance   float64          `json:"payment_balance"`
	Materai          float64          `json:"materai"`
	InvoiceAmount    float64          `json:"invoice_amount"`
	TotalPayment     float64          `json:"total_payment"`
	RemainingPayment float64          `json:"remaining_payment"`
	IsCollection     bool             `json:"is_collection"`
	PayType          int16            `json:"pay_type"`
	InvoiceDate      string           `json:"invoice_date"`
	DocumentNo       string           `json:"document_no"`
	Payment          []DepositPayment `json:"payment" validate:"required,dive,required"`
}

type DepositPayment struct {
	DepositPaymentID int64   `json:"deposit_payment_id"`
	DepositNo        string  `json:"deposit_no"`
	InvoiceNo        string  `json:"invoice_no"`
	PayType          int16   `json:"pay_type"`
	DocumentNo       string  `json:"document_no"`
	Balance          float64 `json:"balance"`
	PaymentAmount    float64 `json:"payment_amount"`

	Images []DepositPaymentImage `json:"images"`
}

type DepositPaymentImage struct {
	DepositImageID int    `json:"deposit_image_id"`
	DepositNo      string `json:"deposit_no"`
	InvoiceNo      string `json:"invoice_no"`
	ImageUrl       string `json:"image_url"`
}

type DetailDepositParams struct {
	DepositNo string `params:"deposit_no" validate:"required"`
}
type DetailInvoiceParams struct {
	InvoiceNo string `params:"invoice_no" validate:"required"`
}

type DepositDetailResponse struct {
	DepositResponse
	Details []DepositDetail         `json:"details"`
	Cash    []DepositPaymentInvoice `json:"cash"`
	Cek     []DepositPaymentInvoice `json:"cek"`
	Trasfer []DepositPaymentInvoice `json:"transfer"`
	Return  []DepositPaymentInvoice `json:"return"`
	CNDN    []DepositPaymentInvoice `json:"cndn"`
}

type DepositInvoiceDetailResponse struct {
	// DepositResponse
	// Details []DepositDetail         `json:"details"`
	Cash    []DepositPaymentInvoice `json:"cash"`
	Cek     []DepositPaymentInvoice `json:"cek"`
	Trasfer []DepositPaymentInvoice `json:"transfer"`
	Return  []DepositPaymentInvoice `json:"return"`
	CNDN    []DepositPaymentInvoice `json:"cndn"`
}

type DepositResponse struct {
	CustID              string     `json:"cust_id"`
	DepositNo           string     `json:"deposit_no"`
	DepositDate         *string    `json:"deposit_date"`
	CollectionNo        *string    `json:"collection_no"`
	CollectionDate      *string    `json:"collection_date"`
	EmpGrpID            *int       `json:"emp_grp_id"`
	EmpGrpName          *string    `json:"emp_grp_name"`
	EmpID               *int       `json:"emp_id"`
	EmpName             *string    `json:"emp_name"`
	EmpCode             *string    `json:"emp_code"`
	OutletGroupID       *int       `json:"outlet_group_id"`
	OutletGroupCode     *string    `json:"outlet_group_code"`
	OutletGroupName     *string    `json:"outlet_group_name"`
	SalesmanID          *int       `json:"salesman_id"`
	SalesmanName        *string    `json:"salesman_name"`
	InvoiceDateFrom     *string    `json:"invoice_date_from"`
	InvoiceDateTo       *string    `json:"invoice_date_to"`
	DueDateFrom         *string    `json:"due_date_from"`
	DueDateTo           *string    `json:"due_date_to"`
	DepositStatus       int        `json:"deposit_status"`
	DepositStatusName   string     `json:"deposit_status_name"`
	RemainingAmount     float64    `json:"remaining_amount"`
	TotalDiscount       float64    `json:"total_discount"`
	TotalMaterai        float64    `json:"total_materai"`
	TotalPaymentBalance float64    `json:"total_payment_balance"`
	TotalPayment        float64    `json:"total_payment"`
	IsApproved          bool       `json:"is_approved"`
	ApprovedBy          *int       `json:"approved_by"`
	ApprovedAt          *time.Time `json:"approved_at"`
	CreatedBy           *int       `json:"created_by"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedBy           *int       `json:"updated_by"`
	UpdatedByName       *int       `json:"updated_by_name"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedBy           *int       `json:"deleted_by"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

type DepositPaymentInvoice struct {
	DepositPayment
	SalesmanId     *int64  `json:"salesman_id"`
	SalesmanCode   *string `json:"salesman_code"`
	SalesmanName   *string `json:"salesman_name"`
	OutletID       *int64  `json:"outlet_id"`
	OutletCode     *string `json:"outlet_code"`
	OutletName     *string `json:"outlet_name"`
	InvoiceDate    string  `json:"invoice_date"`
	Materai        float64 `json:"materai"`
	Discount       float64 `json:"discount"`
	PaymentBalance float64 `json:"payment_balance"`
	TotalPayment   float64 `json:"total_payment"`
}

func ConvStatus(data map[int]string, param int) string {
	statusString, ok := data[int(param)]
	if !ok {
		statusString = "Unknown"
	}
	return statusString
}

var StatusDeposit = map[int]string{
	1: "In Review",
	2: "Approved",
	3: "Rejected",
}

func ConvDepositStatus(data map[int]string, param int) string {
	statusString, ok := data[int(param)]
	if !ok {
		statusString = "Unknown"
	}
	return statusString
}

type MissedPaymentReasonResp struct {
	MissedPaymentId   int64  `json:"missed_payment_reasons_id"`
	MissedPaymentName string `json:"missed_payment_reasons_name"`
	ImageUrl          string `json:"image_url"`
}

type CreateCollectionListBody struct {
	EmpID        int64  `json:"emp_id" validate:"required"`
	CustID       string `json:"cust_id" validate:"required"`
	ParentCustID string `json:"parent_cust_id" validate:"required"`
	UserID       int64  `json:"user_id" validate:"required"`
	IsCollection bool   `json:"is_collection"`
}
