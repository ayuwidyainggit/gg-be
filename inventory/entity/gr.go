package entity

type GrQueryFilter struct {
	CustId              string
	ParentCustId        string
	From                *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To                  *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page                int    `query:"page"`
	Limit               int    `query:"limit" validate:"required"`
	Query               string `query:"q"`
	Mode                string `query:"mode"`
	Sort                string `query:"sort"`
	IsActive            *int   `query:"is_active"`
	InvoiceNo           string `query:"invoice_no"`
	GrType              *int   `query:"gr_type"`
	WhID                *int   `query:"wh_id"`
	SupId               []int  `query:"supplier_id"`
	IsAp                *int   `query:"is_ap"`
	DataStatus          *int   `query:"data_status"`
	GrNo                string `query:"gr_no"`
	ExcludeEmptyInvoice bool   `query:"exclude_empty_invoice"`
}

type GrDetailQuery struct {
	IsAp bool `query:"is_ap"`
}

type GrSupplierQueryFilter struct {
	StartDate *int64 `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	EndDate   *int64 `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
	Query     string `query:"q"`
	Sort      string `query:"sort"`
}

type Gr struct {
	CustID        string  `json:"cust_id"`
	GrDate        string  `json:"gr_date"`
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
	GrType        int64   `json:"gr_type"`
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
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
}
type GrResponse struct {
	GrNo          string  `json:"gr_no"`
	GrDate        string  `json:"gr_date"`
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
	WhName        string  `json:"wh_name"`
	Notes         string  `json:"notes"`
	TotEmbInc     float64 `json:"tot_emb_inc"`
	TotEmbExc     float64 `json:"tot_emb_exc"`
	GrType        int64   `json:"gr_type"`
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
	UpdatedByName string  `json:"updated_by_name"`
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
	Details       []GrDet `json:"details"`
}

type GrListResponse struct {
	GrNo          string `json:"gr_no"`
	GrDate        string `json:"gr_date"`
	DeliveryDate  string `json:"delivery_date"`
	DeliveryNo    string `json:"delivery_no"`
	InvoiceNo     string `json:"invoice_no"`
	InvoiceDate   string `json:"invoice_date"`
	VehicleNo     string `json:"vehicle_no"`
	PoNo          string `json:"po_no"`
	PoDnNo        string `json:"po_dn_no"`
	SupID         int64  `json:"sup_id"`
	SupCode       string `json:"sup_code"`
	SupName       string `json:"sup_name"`
	WhID          int64  `json:"wh_id"`
	WhCode        string `json:"wh_code"`
	WhName        string `json:"wh_name"`
	Notes         string `json:"notes"`
	WithReference bool   `json:"with_reference"`
	DataStatus    int64  `json:"data_status"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedByName string `json:"updated_by_name"`
	IsClosed      bool   `json:"is_closed"`
	ClosedBy      int64  `json:"closed_by"`
	ClosedByName  string `json:"closed_by_name"`
	ClosedAt      string `json:"closed_at"`
}

type GrSupplierListResponse struct {
	SupID   int64  `json:"sup_id"`
	SupCode string `json:"sup_code"`
	SupName string `json:"sup_name"`
}

type GrDistributorListResponse struct {
	CustID          string `json:"cust_id"`
	DistributorId   int64  `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
}

type DetailGrParams struct {
	GrNo string `params:"gr_no" validate:"required"`
}

type DetailGrInvoiceParams struct {
	InvoiceNo string `params:"invoice_no" validate:"required"`
}
type UpdateGrParams struct {
	GrNo string `params:"gr_no" validate:"required"`
}
type DeleteGrParams struct {
	GrNo string `params:"gr_no" validate:"required"`
}

type CreateGrBody struct {
	CustID          string `json:"cust_id" validate:"required"`
	ParentCustID    string
	DistributorID   int64
	DeliveryDate    string         `json:"delivery_date" validate:"required"`
	DeliveryNo      string         `json:"delivery_no" validate:"omitempty,max=25,alphanum"`
	InvoiceNo       string         `json:"invoice_no" validate:"max=20"`
	InvoiceDate     *string        `json:"invoice_date"`
	VehicleNo       string         `json:"vehicle_no" validate:"omitempty,max=25,alphanumspace"`
	PoNo            string         `json:"po_no" validate:"omitempty,max=25,alphanum"`
	PoDnNo          string         `json:"po_dn_no" validate:"max=25"`
	SupID           int64          `json:"sup_id"`
	WhID            int64          `json:"wh_id"`
	Notes           string         `json:"notes" validate:"omitempty,max=40"`
	GoodReceiptType *string        `json:"good_receipt_type" validate:"omitempty"`
	WithReference   *bool          `json:"with_reference"`
	DeliveryFee     *float64       `json:"delivery_fee" validate:"omitempty,gte=0"`
	SoNo            *string        `json:"so_no" validate:"omitempty,max=25,alphanum"`
	CreatedBy       int64          `json:"created_by"`
	UpdatedBy       int64          `json:"updated_by"`
	Details         GrDetWithGroup `json:"details" validate:"required,dive,required"`
}

type UpdateGrRequest struct {
	CustID       string  `json:"cust_id"`
	ParentCustID string  `json:"parent_cust_id"`
	GrNo         *string `json:"gr_no"`
	SeqNo        *string `json:"seq_no"`
	// GrDate       *string        `json:"gr_date"`
	TrCode       *string        `json:"tr_code"`
	DeliveryDate *string        `json:"delivery_date"`
	DeliveryNo   *string        `json:"delivery_no"`
	InvoiceNo    *string        `json:"invoice_no"`
	VehicleNo    *string        `json:"vehicle_no"`
	PoNo         *string        `json:"po_no"`
	PoDnNo       *string        `json:"po_dn_no"`
	SupID        *int64         `json:"sup_id"`
	WhID         *int64         `json:"wh_id"`
	Notes        *string        `json:"notes"`
	TotEmbInc    *float64       `json:"tot_emb_inc"`
	TotEmbExc    *float64       `json:"tot_emb_exc"`
	GrType       *int64         `json:"gr_type"`
	SubTotal     *float64       `json:"sub_total"`
	Vat          *float64       `json:"vat"`
	VatValue     *float64       `json:"vat_value"`
	VatLg        *float64       `json:"vat_lg"`
	VatLgValue   *float64       `json:"vat_lg_value"`
	Total        *float64       `json:"total"`
	VatBg        *float64       `json:"vat_bg"`
	VatBgValue   *float64       `json:"vat_bg_value"`
	WeekIDOb     *int           `json:"week_id_ob"`
	DataStatus   *int64         `json:"data_status"`
	UpdatedBy    int64          `json:"updated_by"`
	Details      GrDetWithGroup `json:"details"`
}

type GrWithDetailResponse struct {
	GrNo            string         `json:"gr_no"`
	GrDate          string         `json:"gr_date"`
	DeliveryDate    string         `json:"delivery_date"`
	DeliveryNo      string         `json:"delivery_no"`
	InvoiceNo       string         `json:"invoice_no"`
	InvoiceDate     string         `json:"invoice_date"`
	VehicleNo       string         `json:"vehicle_no"`
	PoNo            string         `json:"po_no"`
	PoDnNo          string         `json:"po_dn_no"`
	SupID           int64          `json:"sup_id"`
	SupCode         string         `json:"sup_code"`
	SupName         string         `json:"sup_name"`
	WhID            int64          `json:"wh_id"`
	WhCode          string         `json:"wh_code"`
	WhName          string         `json:"wh_name"`
	Notes           string         `json:"notes"`
	GoodReceiptType string         `json:"good_receipt_type"`
	SoNo            string         `json:"so_no"`
	SubTotal        float64        `json:"sub_total"`
	DeliveryFee     float64        `json:"delivery_fee"`
	TotalSkuPrice   float64        `json:"total_sku_price"`
	DiscountValue   *float64       `json:"discount_value,omitempty"`
	Total           float64        `json:"total"`
	TotalVat        float64        `json:"total_vat"`
	TotalVatLgPurch float64        `json:"total_vat_lg_purch"`
	TotalVatBg      float64        `json:"total_vat_bg"`
	DataStatus      int64          `json:"data_status"`
	UpdatedAt       string         `json:"updated_at"`
	UpdatedByName   string         `json:"updated_by_name"`
	IsClosed        bool           `json:"is_closed"`
	ClosedBy        int64          `json:"closed_by"`
	ClosedByName    string         `json:"closed_by_name"`
	ClosedAt        string         `json:"closed_at"`
	Details         GrDetListGroup `json:"details"`
}

type GrWarehouseQueryFilter struct {
	StartDate *int64  `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	EndDate   *int64  `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	Page      int     `query:"page"`
	Limit     int     `query:"limit" validate:"required"`
	SupID     []int64 `query:"sup_id"`
	Query     string  `query:"q"`
	Sort      string  `query:"sort"`
}

type GrWarehouseListResponse struct {
	WhID   int64  `json:"wh_id"`
	WhCode string `json:"wh_code"`
	WhName string `json:"wh_name"`
}

type GrLookupQueryFilter struct {
	CustId       string
	ParentCustId string
	Sort         string `query:"sort"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	GrNo         string `query:"gr_no"`
	SupId        []int  `query:"sup_id"`
	CustIdParam  string `query:"cust_id_param"`
}

type GrLookupResponse struct {
	GrNo string `json:"gr_no"`
}

type GrDownloadQueryFilter struct {
	GrNo string `query:"gr_no" validate:"required"`
}

type GrDownloadResponse struct {
	ReportName  string `json:"report_name"`
	FileBase64  string `json:"file_base64"`
}
