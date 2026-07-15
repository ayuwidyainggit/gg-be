package entity

const (
	CORETAX_TYPE_INVOICE = "I"
	CORETAX_TYPE_RETURN  = "R"

	TAX_IDENTIFIER_TYPE_TIN = "TIN"
)

type CoretaxVatExtractParams struct {
	CoretaxVatExtractID int64 `params:"coretax_vat_extract_id" validate:"required"`
}

type CoreTaxVatExtractQueryFilter struct {
	InvoiceType      string `query:"invoice_type" validate:"required,oneof=I R"`
	ExtractionStatus string `query:"extraction_status"` // E, NE
	SalesmanId       []int  `query:"salesman_id"`
	OutletID         []int  `query:"outlet_id"`
	CustId           string
	ParentCustId     string
	From             *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To               *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	InvoiceFrom      *int64 `query:"invoice_from" validate:"required_with=InvoiceTo,omitempty,gte=1000000000"`
	InvoiceTo        *int64 `query:"invoice_to" validate:"required_with=InvoiceFrom,omitempty,lte=9999999999,gtefield=InvoiceFrom"`
	Page             int    `query:"page"`
	Limit            int    `query:"limit"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	IsActive         *int   `query:"is_active"`
}
type CoreTaxVatExtractListResponse struct {
	TransactionID     string  `json:"transaction_id"`
	InvoiceNo         string  `json:"invoice_no"`
	InvoiceDate       string  `json:"invoice_date"`
	TaxType           string  `json:"tax_type"`
	TaxNo             string  `json:"tax_no"`
	NITKU             string  `json:"nitku"`
	SalesId           *int64  `json:"salesman_id"`
	SalesCode         string  `json:"salesman_code"`
	SalesName         string  `json:"sales_name"`
	OutletID          *int64  `json:"outlet_id"`
	OutletCode        string  `json:"outlet_code"`
	OutletName        string  `json:"outlet_name"`
	OutletAddress1    string  `json:"outlet_address1"`
	OutletAddress2    string  `json:"outlet_address2"`
	OutletTaxAddress1 string  `json:"outlet_tax_address1"`
	OutletTaxAddress2 string  `json:"outlet_tax_address2"`
	DPP               float64 `json:"dpp"`
	DPPFinal          float64 `json:"dpp_final"`
	PPN               float64 `json:"ppn"`
	PPNValue          float64 `json:"ppn_value"`
	PPNFinalValue     float64 `json:"ppn_final_value"`
	PPNBM             float64 `json:"ppnbm"`
	TaxExtractDate    string  `json:"tax_extract_date"`
}

type CoreTaxExtractReq struct {
	TransactionID []string `json:"transaction_id"`
	InvoiceType   string   `json:"invoice_type" validate:"required,oneof=I R"`
	CustID        string
	CreatedBy     int64 `json:"created_by"`
}

type CoretaxVatExtractResp struct {
	ID int64 `json:"id"`
}
type CoretaxVatExtractDownload struct {
	NPWPSeller     string                    `json:"npwp_seller"`
	ExtractResults []CoretaxVatExtractResult `json:"extract_results"`
}
type CoretaxVatExtractResult struct {
	Row                   int                             `json:"row"`
	FakturDate            string                          `json:"faktur_date"`
	FakturType            string                          `json:"faktur_type"`
	TransactionCode       string                          `json:"transaction_code"`
	AdditionalDescription string                          `json:"additional_description"`
	DocumentSupport       string                          `json:"document_support"`
	Reference             string                          `json:"reference"`
	FacilityMark          string                          `json:"facility_mark"`
	SellerTKUId           string                          `json:"seller_tku_id"`
	BuyerNPWPorNIK        string                          `json:"buyer_npwp_or_nik"`
	BuyerTypeID           string                          `json:"buyer_type_id"`
	BuyerCountry          string                          `json:"buyer_country"`
	BuyerDocumentNo       string                          `json:"buyer_document_support_no"`
	BuyerName             string                          `json:"buyer_name"`
	BuyerAddress          string                          `json:"buyer_address"`
	BuyerEmail            string                          `json:"buyer_email"`
	BuyerIDTKU            string                          `json:"buyer_id_tku"`
	Lists                 []CoretaxVatExtractResultDetail `json:"lists"`
}

type CoretaxVatExtractResultDetail struct {
	Item          string  `json:"item"`
	ItemCode      string  `json:"item_code"`
	ItemName      string  `json:"item_name"`
	Qty           float64 `json:"qty"`
	Price         float64 `json:"price"`
	UnitId        string  `json:"unit_id"`
	UnitIdCoretax string  `json:"unit_id_coretax"`
	DPP           float64 `json:"dpp"`
	DPPOther      float64 `json:"dpp_other"`
	PPN           float64 `json:"ppn"`
	PPNValue      float64 `json:"ppn_value"`
	PPNBM         float64 `json:"ppnbm"`
	Discount      float64 `json:"discount"`
}
