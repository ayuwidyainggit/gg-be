package service

import "sales/model"

type invoiceFinalLineAmount struct {
	Gross          float64
	PromoPrimary   float64
	PromoSecondary float64
	Discount       float64
	VAT            float64
	Net            float64
}

type invoiceFinalTotals struct {
	Gross          float64
	PromoPrimary   float64
	PromoSecondary float64
	Discount       float64
	VAT            float64
	Net            float64
}

func (totals invoiceFinalTotals) PromoTotal() float64 {
	return totals.PromoPrimary + totals.PromoSecondary
}

func calculateInvoiceFinalLineAmount(detail model.InvoiceDetRead) invoiceFinalLineAmount {
	gross := (detail.Qty1Final * detail.SellPriceFinal1) +
		(detail.Qty2Final * detail.SellPriceFinal2) +
		(detail.Qty3Final * detail.SellPriceFinal3)
	promoPrimary := getValueOrDefault(detail.PromoFinal1, 0)
	promoSecondary := getValueOrDefault(detail.PromoFinal2, 0) +
		getValueOrDefault(detail.PromoFinal3, 0) +
		getValueOrDefault(detail.PromoFinal4, 0) +
		getValueOrDefault(detail.PromoFinal5, 0)
	discount := getValueOrDefault(detail.DiscValueFinal, 0)
	vat := getValueOrDefault(detail.VatValueFinal, 0)
	promoTotal := promoPrimary + promoSecondary

	return invoiceFinalLineAmount{
		Gross:          gross,
		PromoPrimary:   promoPrimary,
		PromoSecondary: promoSecondary,
		Discount:       discount,
		VAT:            vat,
		Net:            gross - promoTotal - discount + vat,
	}
}

func calculateInvoiceFinalTotals(details []model.InvoiceDetRead) invoiceFinalTotals {
	var totals invoiceFinalTotals
	for _, detail := range details {
		line := calculateInvoiceFinalLineAmount(detail)
		totals.Gross += line.Gross
		totals.PromoPrimary += line.PromoPrimary
		totals.PromoSecondary += line.PromoSecondary
		totals.Discount += line.Discount
		totals.VAT += line.VAT
		totals.Net += line.Net
	}
	return totals
}
