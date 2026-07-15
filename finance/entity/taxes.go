package entity

import (
	"errors"
	"fmt"
)

var (
	TYPE_INVOICE_STANDART = 1
	TYPE_INVOICE_GABUNGAN = 2
)

type TaxesGenerateReq struct {
	Invoices  []string `json:"invoices"`
	CustID    string   `json:"cust_id"`
	CreatedBy int64    `json:"created_by"`
}

type OutletInvoice struct {
	OutletID    int64
	InvoiceType int
	Invoices    []string
}

type MapOutletInvoice map[int64]*OutletInvoice

func (m MapOutletInvoice) MapOutletInvoice(OutletID int64, InvoiceType int, invoice string) {
	if obj, exists := m[OutletID]; exists {
		obj.Invoices = append(obj.Invoices, invoice)
	} else {
		invoices := []string{invoice}
		m[OutletID] = &OutletInvoice{
			OutletID:    OutletID,
			InvoiceType: InvoiceType,
			Invoices:    invoices,
		}
	}
}

type TaxesObj struct {
	TransactionStatusCode string
	SerialCode            string
	OutletInvoices        []OutletInvoice
	RemainingQty          int
	SerialFrom            int
	SerialTo              int
	Start                 int
	InvoiceGenerated      []InvoiceMap
	Status                int
}

type InvoiceMap struct {
	Invoice string
	Tax     string
}

func (t *TaxesObj) GenerateInvoice() error {
	var totalInvoiceWillGenerate int

	for _, OutletInvoice := range t.OutletInvoices {
		if OutletInvoice.InvoiceType == TYPE_INVOICE_STANDART {
			totalInvoiceWillGenerate += len(OutletInvoice.Invoices)
		} else {
			totalInvoiceWillGenerate++
		}
	}

	if t.RemainingQty < totalInvoiceWillGenerate {
		return errors.New(fmt.Sprintf("only %v taxes can be generate", t.RemainingQty))
	}

	t.Start = (t.SerialTo - t.RemainingQty) + 1
	for _, OutletInvoice := range t.OutletInvoices {
		for index, invoice := range OutletInvoice.Invoices {
			if OutletInvoice.InvoiceType == TYPE_INVOICE_STANDART {
				t.InvoiceGenerated = append(t.InvoiceGenerated, InvoiceMap{
					Invoice: invoice,
					Tax:     fmt.Sprintf("%v.%v.%v", t.TransactionStatusCode, t.SerialCode, fmt.Sprintf("%09d", t.Start)),
				})

				t.Start++
				t.RemainingQty--
			} else {
				t.InvoiceGenerated = append(t.InvoiceGenerated, InvoiceMap{
					Invoice: invoice,
					Tax:     fmt.Sprintf("%v.%v.%v", t.TransactionStatusCode, t.SerialCode, fmt.Sprintf("%09d", t.Start)),
				})

				if index == len(OutletInvoice.Invoices)-1 {
					t.Start++
					t.RemainingQty--
				}
			}
		}

		if t.RemainingQty > 0 {
			t.Status = STATUS_TAXES_ACTIVE
		} else {
			t.Status = STATUS_TAXES_COMPLETED
		}

	}
	return nil
}

// func (t *TaxesObj) GenerateInvoice() error {

// 	if t.RemainingQty < len(t.Invoices) {
// 		return errors.New(fmt.Sprintf("only %v taxes can be generate", t.RemainingQty))
// 	}

// 	t.Start = (t.To - t.RemainingQty) + 1
// 	for _, invoice := range t.Invoices {
// 		t.InvoiceGenerated = append(t.InvoiceGenerated, InvoiceMap{
// 			Invoice: invoice,
// 			Tax:     fmt.Sprintf("%v.%v.%v", t.TransactionStatusCode, t.SerialCode, fmt.Sprintf("%09d", t.Start)),
// 		})
// 		t.Start++
// 		t.RemainingQty--
// 	}

// 	if t.RemainingQty > 0 {
// 		t.Status = STATUS_TAXES_ACTIVE
// 	} else {
// 		t.Status = STATUS_TAXES_COMPLETED
// 	}

// 	return nil
// }

type TaxesQueryFilter struct {
	InvoiceNo    []string `query:"invoice_no"`
	SalesmanId   []int    `query:"salesman_id"`
	OutletID     []int    `query:"outlet_id"`
	Status       []int    `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	Taxes        bool   `query:"taxes"`
	Type         string `query:"type"`
}

type TaxesGenerateQueryFilter struct {
	CustId       string
	ParentCustId string
	MTaxID       int64  `query:"m_tax_id"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type TaxesResponse struct {
	SalesmanId      *int64   `json:"salesman_id"`
	SalesmanCode    string   `json:"salesman_code"`
	SalesName       string   `json:"sales_name"`
	OutletID        *int64   `json:"outlet_id"`
	OutletCode      string   `json:"outlet_code"`
	OutletName      string   `json:"outlet_name"`
	OutletAddress   string   `json:"outlet_address"`
	OutletLatitude  string   `json:"outlet_latitude"`
	OutletLongitude string   `json:"outlet_longitude"`
	PayType         *int64   `json:"pay_type"`
	PayTypeName     string   `json:"pay_type_name"`
	MobileID        *int64   `json:"mobile_id"`
	SubTotal        *float64 `json:"sub_total"`
	Vat             *float64 `json:"vat"`
	VatValue        *float64 `json:"vat_value"`
	Total           *float64 `json:"total"`
	DataStatus      *int64   `json:"data_status"`
	DataStatusName  string   `json:"data_status_name"`
	DueDate         *string  `json:"due_date"`
	InvoiceNo       string   `json:"invoice_no"`
	InvoiceDate     string   `json:"invoice_date"`
	Npwpc           string   `json:"npwp"`
	TaxNo           string   `json:"tax_no"`
	TaxesId         int      `json:"taxes_id"`
	TaxDate         string   `json:"tax_date"`
	Type            string   `json:"type"`
}

type TaxesGenerateResponse struct {
	TaxId       int64  `json:"taxes_id"`
	TaxNo       string `json:"tax_no"`
	InvoiceNo   string `json:"invoice_no"`
	InvoiceDate string `json:"invoice_date"`
	Status      int    `json:"status"`
	CustID      string `json:"cust_id"`
}
