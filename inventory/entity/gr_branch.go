package entity

const (
	GR_BRANCH_SUBMITTED = 1
	GR_BRANCH_PROCESSED = 2
	GR_BRANCH_COMPLETED = 3
	GR_BRANCH_REJECTED  = 4
)

var grBranchDataStatusNameList = map[int64]string{
	GR_BRANCH_SUBMITTED: "Submitted",
	GR_BRANCH_PROCESSED: "Processed",
	GR_BRANCH_COMPLETED: "Completed",
	GR_BRANCH_REJECTED:  "Rejected",
}

const (
	ORDER_BOOKING_TYPE_APPROVAL_INTERNAL  = 1
	ORDER_BOOKING_TYPE_APPROVAL_EKSTERNAL = 2
)

var orderBookingTypeApprovalNameList = map[int64]string{
	ORDER_BOOKING_TYPE_APPROVAL_INTERNAL:  "Internal",
	ORDER_BOOKING_TYPE_APPROVAL_EKSTERNAL: "Eksternal",
}

type GrBranchQueryFilter struct {
	CustId        string
	ParentCustId  string
	From          *int64   `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To            *int64   `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page          int      `query:"page"`
	Limit         int      `query:"limit" validate:"required"`
	Query         string   `query:"q"`
	Mode          string   `query:"mode"`
	Sort          string   `query:"sort"`
	IsActive      *int     `query:"is_active"`
	WhID          []int    `query:"wh_id"`
	SupId         []int    `query:"supplier_id"`
	DistributorId []string `query:"cust_id"`
	DataStatus    []int    `query:"data_status"`
}

type GrBranchDetailQuery struct {
	GrBranchCustId string `query:"cust_id" validate:"required"`
}

type GrBranchDetailInvoiceQuery struct {
	IsAp bool `query:"is_ap"`
}

type GrBranchSupplierQueryFilter struct {
	StartDate *int64 `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	EndDate   *int64 `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
	Query     string `query:"q"`
	Sort      string `query:"sort"`
}

type GrBranchDistributorQueryFilter struct {
	StartDate *int64 `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	EndDate   *int64 `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
	Query     string `query:"q"`
	Sort      string `query:"sort"`
}

type GrBranch struct {
	CustID        string  `json:"cust_id"`
	GrBranchDate  string  `json:"gr_branch_date"`
	TrCode        string  `json:"tr_code"`
	DeliveryDate  string  `json:"delivery_date"`
	DeliveryNo    string  `json:"delivery_no"`
	InvoiceNo     string  `json:"invoice_no"`
	VehicleNo     string  `json:"vehicle_no"`
	PoNo          string  `json:"po_no"`
	PoDnNo        string  `json:"po_dn_no"`
	SupID         int64   `json:"sup_id"`
	SupName       string  `json:"sup_name"`
	WhID          int64   `json:"wh_id"`
	Notes         string  `json:"notes"`
	TotEmbInc     float64 `json:"tot_emb_inc"`
	TotEmbExc     float64 `json:"tot_emb_exc"`
	SubTotal      float64 `json:"sub_total"`
	Vat           float64 `json:"vat"`
	VatValue      float64 `json:"vat_value"`
	VatLg         float64 `json:"vat_lg"`
	VatLgValue    float64 `json:"vat_lg_value"`
	Total         float64 `json:"total"`
	VatBg         float64 `json:"vat_bg"`
	VatBgValue    float64 `json:"vat_bg_value"`
	WeekIDOb      int     `json:"week_id_ob"`
	DataStatus    int64   `json:"data_status"`
	CreatedBy     int64   `json:"created_by"`
	UpdatedAt     string  `json:"updated_by"`
	UpdatedByName string  `json:"updated_by_name"`
	// IsClosed      bool    `json:"is_closed"`
	// ClosedBy      int64   `json:"closed_by"`
	// ClosedByName  string  `json:"closed_by_name"`
	// ClosedAt      string  `json:"closed_at"`
}
type GrBranchResponse struct {
	CustId        string        `json:"cust_id"`
	GrBranchNo    string        `json:"gr_branch_no"`
	GrBranchDate  string        `json:"gr_branch_date"`
	TrCode        string        `json:"tr_code"`
	DeliveryDate  string        `json:"delivery_date"`
	DeliveryNo    string        `json:"delivery_no"`
	DeliveryFee   float64       `json:"delivery_fee"`
	InvoiceNo     string        `json:"invoice_no"`
	VehicleNo     string        `json:"vehicle_no"`
	PoNo          string        `json:"po_no"`
	PoDnNo        string        `json:"po_dn_no"`
	SupID         int64         `json:"sup_id"`
	SupName       string        `json:"sup_name"`
	WhID          int64         `json:"wh_id"`
	WhName        string        `json:"wh_name"`
	Notes         string        `json:"notes"`
	TotEmbInc     float64       `json:"tot_emb_inc"`
	TotEmbExc     float64       `json:"tot_emb_exc"`
	SubTotal      float64       `json:"sub_total"`
	Vat           float64       `json:"vat"`
	VatValue      float64       `json:"vat_value"`
	VatLg         float64       `json:"vat_lg"`
	VatLgValue    float64       `json:"vat_lg_value"`
	Total         float64       `json:"total"`
	VatBg         float64       `json:"vat_bg"`
	VatBgValue    float64       `json:"vat_bg_value"`
	WeekIDOb      int           `json:"week_id_ob"`
	DataStatus    int64         `json:"data_status"`
	CreatedBy     int64         `json:"created_by"`
	UpdatedByName string        `json:"updated_by_name"`
	IsPrinted     bool          `json:"is_printed"`
	PrintedBy     int64         `json:"printed_by"`
	PrintedByName string        `json:"printed_by_name"`
	PrintedAt     string        `json:"printed_at"`
	Details       []GrBranchDet `json:"details"`
}

type GrBranchListResponse struct {
	CustId            string  `json:"cust_id"`
	CustName          string  `json:"cust_name"`
	GrBranchNo        string  `json:"gr_branch_no"`
	GrBranchDate      string  `json:"gr_branch_date"`
	DeliveryNo        *string `json:"delivery_no"`
	DeliveryDate      *string `json:"delivery_date"`
	InvoiceNo         *string `json:"invoice_no"`
	InvoiceDate       *string `json:"invoice_date"`
	InvoiceNoBranch   *string `json:"invoice_no_branch"`
	InvoiceDateBranch *string `json:"invoice_date_branch"`
	PoNo              *string `json:"po_no"`
	SoNo              *string `json:"so_no"`
	VehicleNo         *string `json:"vehicle_no"`
	SupID             int64   `json:"sup_id"`
	SupCode           string  `json:"sup_code"`
	SupName           string  `json:"sup_name"`
	WhID              int64   `json:"wh_id"`
	WhCode            string  `json:"wh_code"`
	WhName            string  `json:"wh_name"`
	UpdatedAt         string  `json:"updated_at"`
	UpdatedBy         int64   `json:"updated_by"`
	UpdatedByName     string  `json:"updated_by_name"`
	IsPrint           bool    `json:"is_print"`
	PrintedAt         string  `json:"printed_at"`
	PrintedBy         int64   `json:"printed_by"`
	PrintedByName     string  `json:"printed_by_name"`
	DataStatus        int64   `json:"data_status"`
	DataStatusName    string  `json:"data_status_name"`
}

type GrBranchSupplierListResponse struct {
	SupID   int64  `json:"sup_id"`
	SupCode string `json:"sup_code"`
	SupName string `json:"sup_name"`
}

type GrBranchDistributorListResponse struct {
	CustID   string `json:"cust_id"`
	CustName string `json:"cust_name"`
}

type DetailGrBranchParams struct {
	GrBranchNo string `params:"gr_branch_no" validate:"required"`
}

type DetailGrBranchInvoiceParams struct {
	InvoiceNo string `params:"invoice_no" validate:"required"`
}
type UpdateGrBranchParams struct {
	GrBranchNo string `params:"gr_branch_no" validate:"required"`
}
type DeleteGrBranchParams struct {
	GrBranchNo string `params:"gr_branch_no" validate:"required"`
}

type CreateGrBranchBody struct {
	CustID            string `json:"cust_id" validate:"required"`
	ParentCustID      string
	DeliveryNo        string               `json:"delivery_no" validate:"max=20"`
	DeliveryDate      string               `json:"delivery_date" validate:"required"`
	InvoiceNo         *string              `json:"invoice_no" validate:"omitempty,max=20"`
	InvoiceDate       *string              `json:"invoice_date"`
	VehicleNo         string               `json:"vehicle_no" validate:"max=15"`
	PoNo              string               `json:"po_no" validate:"max=20"`
	SoNo              string               `json:"so_no" validate:"max=20"`
	SupID             int64                `json:"sup_id"`
	WhID              int64                `json:"wh_id"`
	Notes             string               `json:"notes" validate:"max=100"`
	CreatedBy         int64                `json:"created_by"`
	UpdatedBy         int64                `json:"updated_by"`
	Details           GrBranchDetWithGroup `json:"details" validate:"required,dive,required"`
	InvoiceNoBranch   *string              `json:"invoice_no_branch"`
	InvoiceDateBranch *string              `json:"invoice_date_branch"`
}

type UpdateGrBranchRequest struct {
	CustID       string  `json:"cust_id"`
	ParentCustID string  `json:"parent_cust_id"`
	GrBranchNo   *string `json:"gr_branch_no"`
	SeqNo        *string `json:"seq_no"`
	// GrBranchDate       *string        `json:"gr_branch_date"`
	TrCode       *string              `json:"tr_code"`
	DeliveryDate *string              `json:"delivery_date"`
	DeliveryNo   *string              `json:"delivery_no"`
	InvoiceNo    *string              `json:"invoice_no" validate:"omitempty,max=20"`
	VehicleNo    *string              `json:"vehicle_no"`
	PoNo         *string              `json:"po_no"`
	PoDnNo       *string              `json:"po_dn_no"`
	SupID        *int64               `json:"sup_id"`
	WhID         *int64               `json:"wh_id"`
	Notes        *string              `json:"notes"`
	TotEmbInc    *float64             `json:"tot_emb_inc"`
	TotEmbExc    *float64             `json:"tot_emb_exc"`
	SubTotal     *float64             `json:"sub_total"`
	Vat          *float64             `json:"vat"`
	VatValue     *float64             `json:"vat_value"`
	VatLg        *float64             `json:"vat_lg"`
	VatLgValue   *float64             `json:"vat_lg_value"`
	Total        *float64             `json:"total"`
	VatBg        *float64             `json:"vat_bg"`
	VatBgValue   *float64             `json:"vat_bg_value"`
	WeekIDOb     *int                 `json:"week_id_ob"`
	DataStatus   *int64               `json:"data_status"`
	UpdatedBy    int64                `json:"updated_by"`
	Details      GrBranchDetWithGroup `json:"details"`
}

type GrBranchWithDetailResponse struct {
	CustId       string   `json:"cust_id"`
	CustName     string   `json:"cust_name"`
	GrBranchNo   string   `json:"gr_branch_no"`
	GrBranchDate string   `json:"gr_branch_date"`
	DeliveryDate string   `json:"delivery_date"`
	DeliveryNo   string   `json:"delivery_no"`
	DeliveryFee  *float64 `json:"delivery_fee"`
	InvoiceNo    *string  `json:"invoice_no"`
	InvoiceDate  *string  `json:"invoice_date"`
	VehicleNo    string   `json:"vehicle_no"`
	PoNo         string   `json:"po_no"`
	SoNo         string   `json:"so_no"`
	SupID        int64    `json:"sup_id"`
	SupCode      string   `json:"sup_code"`
	SupName      string   `json:"sup_name"`
	WhID         int64    `json:"wh_id"`
	WhCode       string   `json:"wh_code"`
	WhName       string   `json:"wh_name"`
	SubTotal     *float64 `json:"sub_total"`
	VatValue     *float64 `json:"vat_value"`
	Total        *float64 `json:"total"`
	Notes        string   `json:"notes"`
	// SubTotal        float64  `json:"sub_total"`
	// TotalSkuPrice   float64  `json:"total_sku_price"`
	// DiscountValue   *float64 `json:"discount_value,omitempty"`
	// Total           float64  `json:"total"`
	// TotalVat        float64  `json:"total_vat"`
	// TotalVatLgPurch float64  `json:"total_vat_lg_purch"`
	// TotalVatBg      float64  `json:"total_vat_bg"`
	DataStatus           int64                `json:"data_status"`
	DataStatusName       string               `json:"data_status_name"`
	UpdatedAt            string               `json:"updated_at"`
	UpdatedByName        string               `json:"updated_by_name"`
	IsPrint              *bool                `json:"is_print"`
	PrintedBy            *int64               `json:"printed_by"`
	PrintedByName        *string              `json:"printed_by_name"`
	PrintedAt            *string              `json:"printed_at"`
	TypeApproval         *int                 `json:"type_approval"`
	TypeApprovalName     *string              `json:"type_approval_name"`
	InvoiceNoBranch      *string              `json:"invoice_no_branch"`
	InvoiceDateBranch    *string              `json:"invoice_date_branch"`
	InvoiceDueDateBranch *string              `json:"invoice_due_date_branch"`
	Details              GrBranchDetListGroup `json:"details"`
}

type GrBranchOrderBookingDetailResponse struct {
	OrderBookingId       int      `json:"order_booking_id"`
	OrderBookingDetailId int      `json:"order_booking_detail_id"`
	ProId                int      `json:"pro_id"`
	ProCode              string   `json:"pro_code"`
	ProName              string   `json:"pro_name"`
	ItemType             int      `json:"item_type"`
	Qty                  *float64 `json:"qty_bo"`
	QtyAlloc             *float64 `json:"qty_alloc"`
	Qty1                 *float64 `json:"qty1"`
	Qty2                 *float64 `json:"qty2"`
	Qty3                 *float64 `json:"qty3"`
	Qty4                 *float64 `json:"qty4"`
	Qty5                 *float64 `json:"qty5"`
	Qty1Alloc            *float64 `json:"qty1_alloc"`
	Qty2Alloc            *float64 `json:"qty2_alloc"`
	Qty3Alloc            *float64 `json:"qty3_alloc"`
	Qty4Alloc            *float64 `json:"qty4_alloc"`
	Qty5Alloc            *float64 `json:"qty5_alloc"`
	PurchPrice1          *float64 `json:"purch_price1"`
	PurchPrice2          *float64 `json:"purch_price2"`
	PurchPrice3          *float64 `json:"purch_price3"`
	PurchPrice4          *float64 `json:"purch_price4"`
	PurchPrice5          *float64 `json:"purch_price5"`
	SellPrice1           *float64 `json:"sell_price1"`
	SellPrice2           *float64 `json:"sell_price2"`
	SellPrice3           *float64 `json:"sell_price3"`
	SellPrice4           *float64 `json:"sell_price4"`
	SellPrice5           *float64 `json:"sell_price5"`
	Amount               *float64 `json:"amount"`
	AmountAlloc          *float64 `json:"amount_alloc"`
	Vat                  *float64 `json:"vat"`
	VatValue             *float64 `json:"vat_value"`
	VatValueAlloc        *float64 `json:"vat_value_alloc"`
	UnitId1              *string  `json:"unit_id1"`
	UnitId2              *string  `json:"unit_id2"`
	UnitId3              *string  `json:"unit_id3"`
	UnitId4              *string  `json:"unit_id4"`
	UnitId5              *string  `json:"unit_id5"`
	ConvUnit2            *int     `json:"conv_unit2"`
	ConvUnit3            *int     `json:"conv_unit3"`
	ConvUnit4            *int     `json:"conv_unit4"`
	ConvUnit5            *int     `json:"conv_unit5"`
}

type GrBranchWarehouseQueryFilter struct {
	// StartDate *int64  `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	// EndDate   *int64  `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	PoNo         string  `query:"po_no" validate:"required"`
	TypeApproval int     `query:"type_approval" validate:"required"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit"`
	SupID        []int64 `query:"sup_id"`
	Query        string  `query:"q"`
	Sort         string  `query:"sort"`
}

type GrBranchPrintWarehouseQueryFilter struct {
	// StartDate *int64  `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	// EndDate   *int64  `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	CustId       string `query:"cust_id" validate:"required"`
	TypeApproval int    `query:"type_approval" validate:"required"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	// SupID        []int64 `query:"sup_id"`
	Query string `query:"q"`
	Sort  string `query:"sort"`
}

type GrBranchWarehouseListResponse struct {
	WhID   int64  `json:"wh_id"`
	WhCode string `json:"wh_code"`
	WhName string `json:"wh_name"`
}

type GrBranchOrderBookingDetailParams struct {
	OrderBookingId int `params:"order_booking_id" validate:"required"`
}

func (grBranch GrBranchWithDetailResponse) GenerateDataStatusName() string {
	return grBranchDataStatusNameList[grBranch.DataStatus]
}

func (orderBooking GrBranchWithDetailResponse) GenerateOrderBookingTypeApprovalName() string {
	if orderBooking.TypeApproval != nil {
		return orderBookingTypeApprovalNameList[int64(*orderBooking.TypeApproval)]
	}
	return ""
}

func (grBranch GrBranchListResponse) GenerateDataStatusName() string {
	return grBranchDataStatusNameList[grBranch.DataStatus]
}

type GrBranchOrderBookingListQueryFilter struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Query string `query:"q"`
	Sort  string `query:"sort"`
}

type GrBranchOrderBookingListResponse struct {
	OrderBookingId int     `json:"order_booking_id"`
	PoNo           *string `json:"po_no"`
	TypeApproval   *int    `json:"type_approval"`
	SoNo           *string `json:"so_no"`
	SupID          int64   `json:"sup_id"`
	SupCode        string  `json:"sup_code"`
	SupName        string  `json:"sup_name"`
}

type GrBranchUpdateDataStatusBody struct {
	CustId            string  `json:"cust_id"`
	GrBranchNo        string  `json:"gr_branch_no"`
	DataStatus        *int64  `json:"data_status"`
	InvoiceNoBranch   *string `json:"invoice_no_branch"`
	InvoiceDateBranch *string `json:"invoice_date_branch"`
	UpdatedBy         int64   `json:"updated_by"`
}

type GrBranchBulkUpdateDataStatus struct {
	GrBranches []GrBranchUpdateDataStatusBody `json:"gr_branches" validate:"min=1"`
}

type GrBranchBulkPrintBody struct {
	CustId     string `json:"cust_id"`
	GrBranchNo string `json:"gr_branch_no"`
	WhId       int64  `json:"wh_id"`
}

type GrBranchBulkPrint struct {
	GrBranches []GrBranchBulkPrintBody `json:"gr_branches" validate:"min=1"`
}
