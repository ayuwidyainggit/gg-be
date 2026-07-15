package entity

import "finance/pkg/conversion"

const (
	AP_TYPE_INVOICE    = "I"
	AP_TYPE_RETURN     = "R"
	AP_DET_TYPE_NORMAL = 1
	AP_DET_TYPE_PROMO  = 2
)

type ApSupplierInvoiceReturnCreate struct {
	CustId             string `json:"cust_id"`
	CustIdParam        string `json:"cust_id_param"`
	ParentCustID       string
	AccountPayableDate *string `json:"account_payable_date"`
	ApType             string  `json:"ap_type"`
	SupId              int64   `json:"sup_id"`
	InvoiceNo          string  `json:"invoice_no"`
	InvoiceDate        *string `json:"invoice_date"`
	DocumentNo         string  `json:"document_no" validate:"required"`
	TaxInvoiceDate     *string `json:"tax_invoice_date"`
	TaxInvoiceNo       *string `json:"tax_invoice_no"`
	DueDate            *string `json:"due_date"`
	ReturnDate         *string `json:"return_date"`
	DiscountPercent    float64 `json:"discount_percent"`
	Materai            float64 `json:"materai"`
	CreatedBy          int64   `json:"created_by"`
	UpdatedBy          int64   `json:"updated_by"`
	TotalSkuprice      float64
	DiscountValue      float64
	TotalVatValue      float64
	TotalVatLgValue    float64
	Subtotal           float64
	Total              float64
	ProductLists       []*ProductListCreate `json:"product_list"`
}

func (a *ApSupplierInvoiceReturnCreate) Calculate() {
	var skuPrice, subTotal, totalVatValue, totalValLgValue float64
	for _, productList := range a.ProductLists {
		skuPrice += productList.NetAmount
		subTotal += productList.NetAmountAfterInvoiceDiscount
		totalVatValue += productList.VatValue
		totalValLgValue += productList.VatLgValue

	}
	a.DiscountValue = skuPrice * (a.DiscountPercent / 100)
	a.TotalSkuprice = skuPrice
	a.Subtotal = subTotal
	a.TotalVatValue = totalVatValue
	a.TotalVatLgValue = totalValLgValue
	a.Total = a.Subtotal + a.TotalVatValue + a.TotalVatLgValue + a.Materai
}

type ApSupplierInvoiceReturnRespone struct {
	CustId             string                `json:"cust_id"`
	PoNo               string                `json:"po_no"`
	AccountPayableDate string                `json:"account_payable_date"`
	ApType             string                `json:"ap_type"`
	SupId              int64                 `json:"sup_id"`
	SupName            string                `json:"sup_name"`
	SupCode            string                `json:"sup_code"`
	DistributorId      int64                 `json:"distributor_id"`
	DistributorName    string                `json:"distributor"`
	DistributorCode    string                `json:"distributor_code"`
	InvoiceNo          string                `json:"invoice_no"`
	InvoiceDate        string                `json:"invoice_date"`
	DocumentNo         string                `json:"document_no"`
	TaxInvoiceDate     string                `json:"tax_invoice_date"`
	TaxInvoiceNo       string                `json:"tax_invoice_no"`
	DueDate            string                `json:"due_date"`
	Amount             float64               `json:"amount"`
	DiscountRp         float64               `json:"discount_rp"`
	DiscountPercent    float64               `json:"discount_percent"`
	SubTotal           float64               `json:"sub_total"`
	Vat                float64               `json:"vat"`
	VatValue           float64               `json:"vat_value"`
	VatLg              float64               `json:"vat_lg"`
	VatLgValue         float64               `json:"vat_lg_value"`
	Materai            float64               `json:"materai"`
	Total              float64               `json:"total"`
	CreatedBy          int64                 `json:"created_by"`
	CreatedByName      string                `json:"created_by_name"`
	UpdatedBy          int64                 `json:"updated_by"`
	UpdatedByName      string                `json:"updated_by_name"`
	ProductList        []ProductListRespone  `json:"product_list"`
	ProductPromo       []ProductPromoRespone `json:"product_promo"`
}

type ApSupplierInvoiceReturnResponeList struct {
	AccountPayableID   uint    `json:"account_payable_id"`
	AccountPayableDate string  `json:"account_payable_date"`
	ApType             string  `json:"ap_type"`
	SupId              int64   `json:"sup_id"`
	SupName            string  `json:"sup_name"`
	SupCode            string  `json:"sup_code"`
	DistributorId      int64   `json:"distributor_id"`
	DistributorName    string  `json:"distributor"`
	DistributorCode    string  `json:"distributor_code"`
	InvoiceNo          string  `json:"invoice_no"`
	InvoiceDate        string  `json:"invoice_date"`
	DocumentNo         string  `json:"document_no"`
	TaxInvoiceDate     string  `json:"tax_invoice_date"`
	TaxInvoiceNo       string  `json:"tax_invoice_no"`
	DueDate            string  `json:"due_date"`
	Amount             float64 `json:"amount"`
	DiscountRp         float64 `json:"discount_rp"`
	DiscountPercent    float64 `json:"discount_percent"`
	SubTotal           float64 `json:"sub_total"`
	Vat                float64 `json:"vat"`
	VatValue           float64 `json:"vat_value"`
	VatLg              float64 `json:"vat_lg"`
	VatLgValue         float64 `json:"vat_lg_value"`
	Materai            float64 `json:"materai"`
	Total              float64 `json:"total"`
	CreatedBy          int64   `json:"created_by"`
	CreatedByName      string  `json:"created_by_name"`
	UpdatedBy          int64   `json:"updated_by"`
	UpdatedByName      string  `json:"updated_by_name"`
}

type ProductListCreate struct {
	ID                            *int64
	ProId                         int64
	Qty                           int
	Qty1                          float64
	Qty2                          float64
	Qty3                          float64
	UnitPrice1                    float64
	UnitPrice2                    float64
	UnitPrice3                    float64
	ConvUnit2                     float64
	ConvUnit3                     float64
	Gross                         float64
	Disc                          float64
	DiscValue                     float64
	NetAmount                     float64
	Vat                           float64
	VatValue                      float64
	VatLg                         float64
	VatLgValue                    float64
	Total                         float64
	InvoiceDisc                   float64
	InvoiceDiscValue              float64
	NetAmountAfterInvoiceDiscount float64
	Type                          int
}

func (p *ProductListCreate) Calculate() {
	qty := &conversion.Qty{
		Qty:       int(p.Qty),
		ConvUnit2: int(p.ConvUnit2),
		ConvUnit3: int(p.ConvUnit3),
	}
	qtyConversion := qty.ConvToQtyConversion()

	p.Qty1 = float64(qtyConversion.Qty1)
	p.Qty2 = float64(qtyConversion.Qty2)
	p.Qty3 = float64(qtyConversion.Qty3)

	p.Gross = (p.Qty1 * p.UnitPrice1) + (p.Qty2 * p.UnitPrice2) + (p.Qty3 * p.UnitPrice3)

	if p.Disc > 0 {
		p.DiscValue = p.Gross * (p.Disc / 100)
	}

	p.NetAmount = p.Gross - p.DiscValue
	p.Total = p.NetAmount + p.VatValue + p.VatLgValue
	p.InvoiceDiscValue = p.NetAmount * (p.InvoiceDisc / 100)
	p.NetAmountAfterInvoiceDiscount = p.NetAmount - p.InvoiceDiscValue
	p.VatValue = p.NetAmountAfterInvoiceDiscount * (p.Vat / 100)
	p.VatLgValue = p.NetAmountAfterInvoiceDiscount * (p.VatLg / 100)
}

type ProductListRespone struct {
	ProId             *int64   `json:"pro_id"`
	ProCode           *string  `json:"pro_code"`
	ProName           *string  `json:"pro_name"`
	Qty1              int      `json:"qty1"`
	Qty2              int      `json:"qty2"`
	Qty3              int      `json:"qty3"`
	UnitID1           string   `json:"unit_id1"`
	UnitID2           string   `json:"unit_id2"`
	UnitID3           string   `json:"unit_id3"`
	UnitPrice1        *float64 `json:"unit_price1"`
	UnitPrice2        *float64 `json:"unit_price2"`
	UnitPrice3        *float64 `json:"unit_price3"`
	ConvUnit2         float64  `json:"conv_unit2"`
	ConvUnit3         float64  `json:"conv_unit3"`
	SubTotal          *float64 `json:"sub_total"`
	Disc              *float64 `json:"disc"`
	DiscValue         *float64 `json:"disc_value"`
	SubTotalBeforePpn *float64 `json:"sub_total_before_ppn"`
	Vat               *float64 `json:"vat"`
	VatValue          *float64 `json:"vat_value"`
	Total             *float64 `json:"total"`
	VatLg             *float64 `json:"vat_lg"`
	VatLgValue        *float64 `json:"vat_lg_value"`
	Qty               float64  `json:"qty"`
	QtyRemaining1     int      `json:"qty_remaining1"`
	QtyRemaining2     int      `json:"qty_remaining2"`
	QtyRemaining3     int      `json:"qty_remaining3"`
	WhQty1            int      `json:"wh_qty1"`
	WhQty2            int      `json:"wh_qty2"`
	WhQty3            int      `json:"wh_qty3"`
}

type ProductPromoCreate struct {
	ProId      int64
	Qty        int
	ConvUnit2  float64
	ConvUnit3  float64
	UnitPrice1 float64
	UnitPrice2 float64
	UnitPrice3 float64
}

type ProductPromoRespone struct {
	InvoiceNo  *string  `json:"invoice_no"`
	ProId      *int64   `json:"pro_id"`
	ProCode    *string  `json:"pro_code"`
	ProName    *string  `json:"pro_name"`
	Qty1       int      `json:"qty1"`
	Qty2       int      `json:"qty2"`
	Qty3       int      `json:"qty3"`
	UnitID1    string   `json:"unit_id1"`
	UnitID2    string   `json:"unit_id2"`
	UnitID3    string   `json:"unit_id3"`
	UnitPrice1 *float64 `json:"unit_price1"`
	UnitPrice2 *float64 `json:"unit_price2"`
	UnitPrice3 *float64 `json:"unit_price3"`
	ConvUnit2  float64  `json:"conv_unit2"`
	ConvUnit3  float64  `json:"conv_unit3"`
	Qty        float64  `json:"qty"`
}

type UpdateGr struct {
	InvoiceNo   *string `json:"invoice_no"`
	InvoiceDate *string `json:"invoice_date"`
	IsAp        bool    `json:"is_ap"`
}

type ApSupplierInoviceReturnQueryFilter struct {
	CustId              string
	ParentCustId        string
	From                *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To                  *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page                int    `query:"page"`
	Limit               int    `query:"limit" validate:""`
	Query               string `query:"q"`
	Mode                string `query:"mode"`
	Sort                string `query:"sort"`
	SuppId              int    `query:"sup_id"`
	DocumentNo          string `query:"document_no"`
	Type                string `query:"Type"`
	InvoiceNo           string `query:"invoice_no"`
	ExcludeEmptyInvoice bool   `query:"exclude_empty_invoice"`
}

type DetailApSupplierInvoiceReturnParams struct {
	AccountPayableID uint `params:"account_payable_id" validate:"required" json:"account_payable_id"`
}

type UpdateApSupplierInvoiceReturnParams struct {
	AccountPayableID uint `params:"account_payable_id" validate:"required" json:"account_payable_id"`
}

type DeleteApSupplierInvoiceReturnParams struct {
	AccountPayableID uint `params:"account_payable_id" validate:"required" json:"account_payable_id"`
}
