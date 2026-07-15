package entity

import "time"

type CreateSupplierReturnBody struct {
	CustID          string `json:"cust_id"`
	SupID           int64  `json:"sup_id" validate:"required"`
	WhID            int64  `json:"wh_id" validate:"required"`
	InvoiceNo       string `json:"invoice_no" validate:"required"`
	Notes           string `json:"notes"`
	CreatedBy       int64
	UpdatedBy       int64
	DiscountValue   float64
	TotalVatValue   float64
	TotalVatLgValue float64
	TotalVatBgValue float64
	Subtotal        float64
	Total           float64
	Details         []*CreateSupplierReturnDetBody `json:"details" validate:"required,dive"`
}

func (s *CreateSupplierReturnBody) Calculate() {
	for _, Detail := range s.Details {
		s.Subtotal += Detail.Subtotal
		s.DiscountValue += Detail.DiscValue
		s.TotalVatValue += Detail.VatValue
		s.TotalVatLgValue += Detail.VatLgValue
		s.TotalVatBgValue += Detail.VatBgValue
		s.Total += Detail.Total
	}
}

type SupplierReturnListResponse struct {
	SupplierReturnNo   string     `json:"supplier_return_no"`
	SupplierReturnDate string     `json:"supplier_return_date"`
	InvoiceNo          string     `json:"invoice_no"`
	InvoiceDate        string     `json:"invoice_date"`
	TaxInvoiceDate     string     `json:"tax_invoice_date"`
	TaxInvoiceNo       string     `json:"tax_invoice_no"`
	DueDate            string     `json:"due_date"`
	SupID              int64      `json:"sup_id"`
	SupCode            string     `json:"sup_code"`
	SupName            string     `json:"sup_name"`
	WhID               int64      `json:"wh_id"`
	WhCode             string     `json:"wh_code"`
	WhName             string     `json:"wh_name"`
	Notes              string     `json:"notes"`
	Total              float64    `json:"total"`
	DataStatus         int64      `json:"status"`
	UpdatedAt          *time.Time `json:"updated_at"`
	UpdatedByName      string     `json:"updated_by_name"`
	IsClosed           bool       `json:"is_closed"`
	ClosedBy           int64      `json:"closed_by"`
	ClosedByName       string     `json:"closed_by_name"`
	ClosedAt           *time.Time `json:"closed_at"`
}

type SupplierReturnGetResp struct {
	SupplierReturnNo   string                     `json:"supplier_return_no"`
	SupplierReturnDate string                     `json:"supplier_return_date"`
	InvoiceNo          string                     `json:"invoice_no"`
	InvoiceDate        string                     `json:"invoice_date"`
	TaxInvoiceDate     string                     `json:"tax_invoice_date"`
	TaxInvoiceNo       string                     `json:"tax_invoice_no"`
	DueDate            string                     `json:"due_date"`
	SupID              int64                      `json:"sup_id"`
	SupCode            string                     `json:"sup_code"`
	SupName            string                     `json:"sup_name"`
	WhID               int64                      `json:"wh_id"`
	WhCode             string                     `json:"wh_code"`
	WhName             string                     `json:"wh_name"`
	Notes              string                     `json:"notes"`
	DataStatus         int64                      `json:"status"`
	UpdatedAt          *time.Time                 `json:"updated_at"`
	UpdatedByName      string                     `json:"updated_by_name"`
	IsClosed           bool                       `json:"is_closed"`
	ClosedBy           int64                      `json:"closed_by"`
	ClosedByName       string                     `json:"closed_by_name"`
	ClosedAt           *time.Time                 `json:"closed_at"`
	VatValue           float64                    `json:"vat_value"`
	VatLgValue         float64                    `json:"vat_lg_value"`
	VatBgValue         float64                    `json:"vat_bg_value"`
	SubTotal           float64                    `json:"sub_total"`
	TotalSkuPrice      float64                    `json:"total_sku_price"`
	Total              float64                    `json:"total"`
	DiscountValue      float64                    `json:"discount_value"`
	Details            []SupplierReturnGetDetResp `json:"details"`
}
type ReturnSupplierQueryFilter struct {
	StartDate *int64 `query:"startDate" validate:"required_with=EndDate,omitempty,gte=1000000000"`
	EndDate   *int64 `query:"endDate" validate:"required_with=StartDate,omitempty,lte=9999999999,gtefield=StartDate"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
}
type SupplierReturnResponse struct {
	SupplierReturnNo string `params:"supplier_return_no"`
}
type DetailSupplierReturnParams struct {
	SupplierReturnNo string `json:"return_supplier_no" validate:"required"`
}
type SupplierReturnQueryFilter struct {
	CustId           string
	ParentCustId     string
	From             *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To               *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page             int    `query:"page"`
	Limit            int    `query:"limit" validate:"required"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	IsActive         *int   `query:"is_active"`
	TrCode           string `query:"tr_code"`
	SupplierReturnNo string `query:"supplier_return_no"`
	GrType           *int   `query:"gr_type"`
	SupId            []int  `query:"supplier_id"`
	IsAp             *int   `query:"is_ap"`
	Status           []int  `query:"status"`
}
type SupplierReturnSupplierListResponse struct {
	SupID   int64  `json:"sup_id"`
	SupCode string `json:"sup_code"`
	SupName string `json:"sup_name"`
}

type UpdateSupplierReturnParams struct {
	SupplierReturnNo string `params:"supplier_return_no"`
}
type UpdateSupplierReturnStatusBody struct {
	CustID       string `json:"cust_id"`
	ParentCustID string
	DataStatus   int   `json:"status" validate:"oneof=2 9"`
	UpdatedBy    int64 `json:"updated_by"`
}
