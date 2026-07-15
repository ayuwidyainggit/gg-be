package service

import (
	"context"
	"testing"
	"time"

	"sales/entity"
	"sales/model"
)

type mockInvoiceRepositoryFinal struct {
	invoiceByNo        map[string]model.InvoiceList
	detailsByRoNo      map[string][]model.InvoiceDetRead
	listRows           []model.InvoiceList
	invoiceNumberRows  []model.InvoiceList
	generatedInvoiceNo string
	updatedByRoNo      map[string]model.Invoice
}

func newMockInvoiceRepositoryFinal() *mockInvoiceRepositoryFinal {
	return &mockInvoiceRepositoryFinal{
		invoiceByNo:        make(map[string]model.InvoiceList),
		detailsByRoNo:      make(map[string][]model.InvoiceDetRead),
		generatedInvoiceNo: "INV2606150001",
		updatedByRoNo:      make(map[string]model.Invoice),
	}
}

func (m *mockInvoiceRepositoryFinal) FindByNo(roNo string, custId string) (model.InvoiceList, error) {
	return m.invoiceByNo[roNo], nil
}

func (m *mockInvoiceRepositoryFinal) FindDetail(roNo string, custId string) ([]model.InvoiceDetRead, error) {
	return m.detailsByRoNo[roNo], nil
}

func (m *mockInvoiceRepositoryFinal) FindAllByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, int64, int, error) {
	return m.listRows, int64(len(m.listRows)), 1, nil
}

func (m *mockInvoiceRepositoryFinal) FindAllByInvoiceNombersAndCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, error) {
	return m.invoiceNumberRows, nil
}

func (m *mockInvoiceRepositoryFinal) GenerateInvoiceNo(c context.Context, custId string, invoiceDate time.Time) (string, error) {
	return m.generatedInvoiceNo, nil
}

func (m *mockInvoiceRepositoryFinal) Update(c context.Context, roNo, custID string, data model.Invoice) error {
	m.updatedByRoNo[roNo] = data
	return nil
}

func (m *mockInvoiceRepositoryFinal) UpdateOutletStatusFromPreDormantIfSet(c context.Context, custID string, outletID int64, updatedBy int64) error {
	return nil
}

func (m *mockInvoiceRepositoryFinal) Print(c context.Context, custId string, invoiceNo string, printedBy int64) error {
	return nil
}

type mockStockRepositoryFinal struct{}

func (m *mockStockRepositoryFinal) StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error {
	return nil
}

func (m *mockStockRepositoryFinal) SalesStockUpdates(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
	return nil
}

func (m *mockStockRepositoryFinal) InvoiceSalesStockUpdates(c context.Context, stockUpdates []*entity.InvoiceSalesStockUpdate) error {
	return nil
}

func (m *mockStockRepositoryFinal) CancelSalesStockUpdates(c context.Context, orderNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
	return nil
}

func (m *mockStockRepositoryFinal) GetCancelStockBasis(c context.Context, custID string, orderNo string) ([]entity.CancelStockBasis, error) {
	return nil, nil
}

func (m *mockStockRepositoryFinal) UpdateOnCustomerOrder(c context.Context, custId string, whId int64, proId int64, delta float64) error {
	return nil
}

func (m *mockStockRepositoryFinal) GetCurrentStock(c context.Context, custId string, whId int64, proId int64) (float64, error) {
	return 0, nil
}

func TestInvoiceFinalLineAmountUsesFinalOrderFieldsAndNullPromo(t *testing.T) {
	detail := invoiceFinalRegressionDetails()[0]

	line := calculateInvoiceFinalLineAmount(detail)

	if line.Gross != 16500000 {
		t.Fatalf("expected final gross 16500000, got %v", line.Gross)
	}
	if line.PromoPrimary != 0 || line.PromoSecondary != 0 {
		t.Fatalf("expected nil final promo fields to sum to 0, got primary=%v secondary=%v", line.PromoPrimary, line.PromoSecondary)
	}
	if line.Discount != 0 {
		t.Fatalf("expected final discount 0, got %v", line.Discount)
	}
	if line.VAT != 1650000 {
		t.Fatalf("expected final vat 1650000, got %v", line.VAT)
	}
	if line.Net != 18150000 {
		t.Fatalf("expected final net 18150000, got %v", line.Net)
	}
}

func TestInvoiceFinalLineAmountSubtractsAllFinalPromoFields(t *testing.T) {
	detail := invoiceFinalRegressionDetails()[0]
	detail.PromoFinal1 = float64Ptr(100)
	detail.PromoFinal2 = float64Ptr(200)
	detail.PromoFinal3 = float64Ptr(300)
	detail.PromoFinal4 = nil
	detail.PromoFinal5 = float64Ptr(400)

	line := calculateInvoiceFinalLineAmount(detail)

	if line.PromoPrimary != 100 {
		t.Fatalf("expected final primary promo 100, got %v", line.PromoPrimary)
	}
	if line.PromoSecondary != 900 {
		t.Fatalf("expected final secondary promo 900, got %v", line.PromoSecondary)
	}
	if line.Net != 18149000 {
		t.Fatalf("expected final net after all promo fields 18149000, got %v", line.Net)
	}
}

func TestInvoiceDetailUsesFinalOrderFieldsAndHeaderTotals(t *testing.T) {
	repo := newMockInvoiceRepositoryFinal()
	stockRepo := &mockStockRepositoryFinal{}
	trx := &mockTransactionPassThrough{}
	service := NewInvoiceService(repo, stockRepo, trx)

	custID := "C220010001"
	roNo := "RO-FINAL-DETAIL"
	repo.invoiceByNo[roNo] = invoiceListFixture(roNo)
	repo.detailsByRoNo[roNo] = invoiceFinalRegressionDetails()

	response, err := service.Detail(roNo, custID)
	if err != nil {
		t.Fatalf("detail should succeed: %v", err)
	}

	assertInvoiceFinalHeader(t, response.SubTotal, response.PromoValue, response.DiscValue, response.VatValue, response.Total)
	if len(response.Details.Normal) != 2 {
		t.Fatalf("expected 2 normal details, got %d", len(response.Details.Normal))
	}
	line := response.Details.Normal[0]
	if line.Qty1 != 11 {
		t.Fatalf("expected final qty1 11, got %v", line.Qty1)
	}
	if line.SellPrice1 != 1500000 {
		t.Fatalf("expected final sell_price1 1500000, got %v", line.SellPrice1)
	}
	if line.Amount != 16500000 {
		t.Fatalf("expected final gross amount 16500000, got %v", line.Amount)
	}
	if line.VatValue != 1650000 {
		t.Fatalf("expected final vat value 1650000, got %v", line.VatValue)
	}
	if line.NetValue != 18150000 {
		t.Fatalf("expected final net value 18150000, got %v", line.NetValue)
	}
	second := response.Details.Normal[1]
	if second.Amount != 2100000 {
		t.Fatalf("expected second detail final gross amount 2100000, got %v", second.Amount)
	}
}

func TestInvoiceListUsesFinalOrderTotals(t *testing.T) {
	repo := newMockInvoiceRepositoryFinal()
	stockRepo := &mockStockRepositoryFinal{}
	trx := &mockTransactionPassThrough{}
	service := NewInvoiceService(repo, stockRepo, trx)

	row := invoiceListFixture("RO-FINAL-LIST")
	repo.listRows = []model.InvoiceList{row}
	repo.detailsByRoNo[row.OrderNo] = invoiceFinalRegressionDetails()

	response, total, lastPage, err := service.List(entity.InvoiceQueryFilter{CustId: "C220010001", Limit: 10, Page: 1})
	if err != nil {
		t.Fatalf("list should succeed: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("unexpected pagination total=%d lastPage=%d", total, lastPage)
	}
	if len(response) != 1 {
		t.Fatalf("expected 1 invoice row, got %d", len(response))
	}

	assertInvoiceFinalHeader(t, response[0].SubTotal, response[0].PromoValue, response[0].DiscValue, response[0].VatValue, response[0].Total)
	if len(response[0].Details) != 2 {
		t.Fatalf("expected 2 list detail rows, got %d", len(response[0].Details))
	}
	if response[0].Details[0].Amount != 16500000 {
		t.Fatalf("expected first list detail final amount 16500000, got %v", response[0].Details[0].Amount)
	}
	if response[0].Details[1].Amount != 2100000 {
		t.Fatalf("expected second list detail final amount 2100000, got %v", response[0].Details[1].Amount)
	}
}

func TestInvoiceDetailsUsesFinalOrderTotals(t *testing.T) {
	repo := newMockInvoiceRepositoryFinal()
	stockRepo := &mockStockRepositoryFinal{}
	trx := &mockTransactionPassThrough{}
	service := NewInvoiceService(repo, stockRepo, trx)

	row := invoiceListFixture("RO-FINAL-DETAILS")
	repo.invoiceNumberRows = []model.InvoiceList{row}
	repo.detailsByRoNo[row.OrderNo] = invoiceFinalRegressionDetails()

	response, err := service.Details(entity.InvoiceQueryFilter{CustId: "C220010001", InvoiceNo: []string{"INV2606150001"}})
	if err != nil {
		t.Fatalf("details should succeed: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("expected 1 invoice row, got %d", len(response))
	}

	assertInvoiceFinalHeader(t, response[0].SubTotal, response[0].PromoValue, response[0].DiscValue, response[0].VatValue, response[0].Total)
	if response[0].Details[0].NetValue != 18150000 {
		t.Fatalf("expected first details final net 18150000, got %v", response[0].Details[0].NetValue)
	}
	if response[0].Details[1].NetValue != 2310000 {
		t.Fatalf("expected second details final net 2310000, got %v", response[0].Details[1].NetValue)
	}
}

func TestInvoiceBulkUpdatePersistsFinalHeaderTotals(t *testing.T) {
	repo := newMockInvoiceRepositoryFinal()
	stockRepo := &mockStockRepositoryFinal{}
	trx := &mockTransactionPassThrough{}
	service := NewInvoiceService(repo, stockRepo, trx)

	custID := "C220010001"
	roNo := "RO-FINAL-BULK"
	whID := int64(63)
	orderDate := time.Now()
	repo.invoiceByNo[roNo] = model.InvoiceList{
		WhId:        &whID,
		RoDate:      &orderDate,
		PaymentType: entity.PAY_TYPE_CASH_ON_DELIVERY,
	}
	bulkDetails := invoiceFinalRegressionDetails()
	for i := range bulkDetails {
		bulkDetails[i].ConvUnit2 = 1
		bulkDetails[i].ConvUnit3 = 1
	}
	repo.detailsByRoNo[roNo] = bulkDetails

	roDate := "2026-06-15"
	valDate := "2026-06-15"
	deliveryDate := "2026-06-15"
	staleSubTotal := 3960000.0
	staleDiscValue := 10.0
	stalePromoValue := 20.0
	staleVatValue := 30.0
	staleTotal := 3960000.0

	err := service.BulkUpdate(custID, entity.BulkUpdateInvoiceBody{Orders: []entity.UpdateInvoiceBody{{
		RoNo:         roNo,
		RoDate:       &roDate,
		ValDate:      &valDate,
		DeliveryDate: &deliveryDate,
		SubTotal:     &staleSubTotal,
		DiscValue:    &staleDiscValue,
		PromoValue:   &stalePromoValue,
		VatValue:     &staleVatValue,
		Total:        &staleTotal,
	}}})
	if err != nil {
		t.Fatalf("bulk update should succeed: %v", err)
	}

	updated := repo.updatedByRoNo[roNo]
	if updated.SubTotalFinal == nil || *updated.SubTotalFinal != 18600000 {
		t.Fatalf("expected final subtotal 18600000, got %+v", updated.SubTotalFinal)
	}
	if updated.PromoValueFinal == nil || *updated.PromoValueFinal != 0 {
		t.Fatalf("expected final promo 0, got %+v", updated.PromoValueFinal)
	}
	if updated.DiscValueFinal == nil || *updated.DiscValueFinal != 0 {
		t.Fatalf("expected final discount 0, got %+v", updated.DiscValueFinal)
	}
	if updated.VatValueFinal == nil || *updated.VatValueFinal != 1860000 {
		t.Fatalf("expected final vat 1860000, got %+v", updated.VatValueFinal)
	}
	if updated.TotalFinal == nil || *updated.TotalFinal != 20460000 {
		t.Fatalf("expected final total 20460000, got %+v", updated.TotalFinal)
	}
	if updated.SubTotal != nil || updated.PromoValue != nil || updated.DiscValue != nil || updated.VatValue != nil || updated.Total != nil {
		t.Fatalf("expected stale non-final monetary fields cleared, got %+v", updated)
	}
}

func invoiceFinalRegressionDetails() []model.InvoiceDetRead {
	firstOrderDetID := 101
	firstVatValueFinal := 1650000.0
	secondOrderDetID := 102
	secondVatValueFinal := 210000.0

	return []model.InvoiceDetRead{
		{
			CustId:          "C220010001",
			RoNo:            "RO-FINAL",
			SeqNo:           1,
			OrderDetID:      &firstOrderDetID,
			ProId:           748,
			ProCode:         "TP-012",
			ProName:         "Toothpaste",
			ItemType:        1,
			Qty1:            1,
			Qty2:            1,
			Qty3:            0,
			Qty1Final:       11,
			Qty2Final:       0,
			Qty3Final:       0,
			SellPrice1:      1500000,
			SellPrice2:      2460000,
			SellPrice3:      0,
			SellPriceFinal1: 1500000,
			SellPriceFinal2: 0,
			SellPriceFinal3: 0,
			Amount:          3960000,
			AmountFinal:     123,
			DiscValue:       50,
			DiscValueFinal:  float64Ptr(0),
			PromoFinal1:     nil,
			PromoFinal2:     nil,
			PromoFinal3:     nil,
			PromoFinal4:     nil,
			PromoFinal5:     nil,
			VatValue:        396000,
			VatValueFinal:   &firstVatValueFinal,
		},
		{
			CustId:          "C220010001",
			RoNo:            "RO-FINAL",
			SeqNo:           2,
			OrderDetID:      &secondOrderDetID,
			ProId:           749,
			ProCode:         "TP-013",
			ProName:         "Toothbrush",
			ItemType:        1,
			Qty1:            0,
			Qty2:            0,
			Qty3:            0,
			Qty1Final:       1,
			Qty2Final:       0,
			Qty3Final:       0,
			SellPrice1:      0,
			SellPrice2:      0,
			SellPrice3:      0,
			SellPriceFinal1: 2100000,
			SellPriceFinal2: 0,
			SellPriceFinal3: 0,
			Amount:          0,
			AmountFinal:     0,
			DiscValue:       0,
			DiscValueFinal:  float64Ptr(0),
			PromoFinal1:     nil,
			PromoFinal2:     nil,
			PromoFinal3:     nil,
			PromoFinal4:     nil,
			PromoFinal5:     nil,
			VatValue:        0,
			VatValueFinal:   &secondVatValueFinal,
		},
	}
}

func invoiceListFixture(roNo string) model.InvoiceList {
	staleSubTotal := 3960000.0
	stalePromoValue := 100.0
	staleDiscValue := 200.0
	staleVatValue := 396000.0
	staleTotal := 3960000.0
	invoiceDate := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	return model.InvoiceList{
		CustID:      "C220010001",
		OrderNo:     roNo,
		SubTotal:    &staleSubTotal,
		PromoValue:  &stalePromoValue,
		DiscValue:   &staleDiscValue,
		VatValue:    &staleVatValue,
		Total:       &staleTotal,
		InvoiceDate: &invoiceDate,
	}
}

func assertInvoiceFinalHeader(t *testing.T, subTotal, promoValue, discValue, vatValue, total *float64) {
	t.Helper()
	if subTotal == nil || *subTotal != 18600000 {
		t.Fatalf("expected final subtotal 18600000, got %+v", subTotal)
	}
	if promoValue == nil || *promoValue != 0 {
		t.Fatalf("expected final promo 0, got %+v", promoValue)
	}
	if discValue == nil || *discValue != 0 {
		t.Fatalf("expected final discount 0, got %+v", discValue)
	}
	if vatValue == nil || *vatValue != 1860000 {
		t.Fatalf("expected final vat 1860000, got %+v", vatValue)
	}
	if total == nil || *total != 20460000 {
		t.Fatalf("expected final total 20460000, got %+v", total)
	}
}
