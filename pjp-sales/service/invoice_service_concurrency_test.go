package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"sales/entity"
	"sales/model"
)

type mockInvoiceRepositoryConcurrency struct {
	mu                  sync.Mutex
	sequence            int
	generated           map[string]struct{}
	orders              map[string]model.InvoiceList
	details             map[string][]model.InvoiceDetRead
	failOnceByInvoiceNo map[string]bool
}

func newMockInvoiceRepositoryConcurrency() *mockInvoiceRepositoryConcurrency {
	return &mockInvoiceRepositoryConcurrency{
		generated:           make(map[string]struct{}),
		orders:              make(map[string]model.InvoiceList),
		details:             make(map[string][]model.InvoiceDetRead),
		failOnceByInvoiceNo: make(map[string]bool),
	}
}

func (m *mockInvoiceRepositoryConcurrency) FindByNo(roNo string, custId string) (model.InvoiceList, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	order, ok := m.orders[roNo]
	if !ok {
		return model.InvoiceList{}, errors.New("not found")
	}
	return order, nil
}

func (m *mockInvoiceRepositoryConcurrency) FindDetail(roNo string, custId string) ([]model.InvoiceDetRead, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.details[roNo], nil
}

func (m *mockInvoiceRepositoryConcurrency) FindAllByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockInvoiceRepositoryConcurrency) FindAllByInvoiceNombersAndCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, error) {
	return nil, nil
}

func (m *mockInvoiceRepositoryConcurrency) GenerateInvoiceNo(c context.Context, custId string, invoiceDate time.Time) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sequence++
	invoiceNo := fmt.Sprintf("INV%s%04d", invoiceDate.Format("060102"), m.sequence)
	return invoiceNo, nil
}

func (m *mockInvoiceRepositoryConcurrency) Update(c context.Context, roNo, custID string, data model.Invoice) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if data.InvoiceNo == nil {
		return errors.New("invoice_no is required")
	}

	invoiceNo := *data.InvoiceNo
	if m.failOnceByInvoiceNo[invoiceNo] {
		delete(m.failOnceByInvoiceNo, invoiceNo)
		return errors.New("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)")
	}

	if _, exists := m.generated[invoiceNo]; exists {
		return errors.New("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)")
	}
	m.generated[invoiceNo] = struct{}{}
	return nil
}

func (m *mockInvoiceRepositoryConcurrency) Print(c context.Context, custId string, invoiceNo string, printedBy int64) error {
	return nil
}

func (m *mockInvoiceRepositoryConcurrency) UpdateOutletStatusFromPreDormantIfSet(c context.Context, custID string, outletID int64, updatedBy int64) error {
	return nil
}

type mockStockRepositoryConcurrency struct{}

func (m *mockStockRepositoryConcurrency) StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error {
	return nil
}

func (m *mockStockRepositoryConcurrency) SalesStockUpdates(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
	return nil
}

func (m *mockStockRepositoryConcurrency) InvoiceSalesStockUpdates(c context.Context, stockUpdates []*entity.InvoiceSalesStockUpdate) error {
	return nil
}

func (m *mockStockRepositoryConcurrency) CancelSalesStockUpdates(c context.Context, orderNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
	return nil
}

func (m *mockStockRepositoryConcurrency) GetCancelStockBasis(c context.Context, custID string, orderNo string) ([]entity.CancelStockBasis, error) {
	return nil, nil
}

func (m *mockStockRepositoryConcurrency) UpdateOnCustomerOrder(c context.Context, custId string, whId int64, proId int64, delta float64) error {
	return nil
}

func (m *mockStockRepositoryConcurrency) GetCurrentStock(c context.Context, custId string, whId int64, proId int64) (float64, error) {
	return 0, nil
}

type mockTransactionPassThrough struct{}

func (m *mockTransactionPassThrough) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return tFunc(ctx)
}

func TestInvoiceBulkUpdateConcurrentNoDuplicateInvoiceNo(t *testing.T) {
	repo := newMockInvoiceRepositoryConcurrency()
	stockRepo := &mockStockRepositoryConcurrency{}
	trx := &mockTransactionPassThrough{}
	service := NewInvoiceService(repo, stockRepo, trx)

	custID := "C220010001"
	whID := int64(63)
	roDate := "2026-03-03"
	valDate := "2026-03-03"
	deliveryDate := "2026-03-03"

	totalOrder := 25
	for i := 0; i < totalOrder; i++ {
		roNo := fmt.Sprintf("RO%05d", i+1)
		orderDate := time.Now()
		repo.orders[roNo] = model.InvoiceList{
			WhId:        &whID,
			RoDate:      &orderDate,
			PaymentType: entity.PAY_TYPE_CASH_ON_DELIVERY,
		}
		orderDetID := i + 1
		repo.details[roNo] = []model.InvoiceDetRead{{
			CustId:     custID,
			ProId:      1000 + i,
			Qty1Final:  1,
			Qty2Final:  0,
			Qty3Final:  0,
			ConvUnit2:  1,
			ConvUnit3:  1,
			SellPrice1: 1000,
			OrderDetID: &orderDetID,
		}}
	}

	var wg sync.WaitGroup
	errCh := make(chan error, totalOrder)

	for i := 0; i < totalOrder; i++ {
		wg.Add(1)
		roNo := fmt.Sprintf("RO%05d", i+1)
		go func(ro string) {
			defer wg.Done()
			err := service.BulkUpdate(custID, entity.BulkUpdateInvoiceBody{
				Orders: []entity.UpdateInvoiceBody{{
					RoNo:         ro,
					RoDate:       &roDate,
					ValDate:      &valDate,
					DeliveryDate: &deliveryDate,
				}},
			})
			if err != nil {
				errCh <- err
			}
		}(roNo)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("bulk update should succeed without duplicate invoice_no: %v", err)
		}
	}

	if len(repo.generated) != totalOrder {
		t.Fatalf("expected %d unique invoice numbers, got %d", totalOrder, len(repo.generated))
	}
}

func TestInvoiceBulkUpdateRetryOnUniqueViolation(t *testing.T) {
	repo := newMockInvoiceRepositoryConcurrency()
	stockRepo := &mockStockRepositoryConcurrency{}
	trx := &mockTransactionPassThrough{}
	service := NewInvoiceService(repo, stockRepo, trx)

	custID := "C220010001"
	roNo := "RO-RETRY-1"
	whID := int64(63)
	roDate := "2026-03-03"
	valDate := "2026-03-03"
	deliveryDate := "2026-03-03"

	orderDate := time.Now()
	repo.orders[roNo] = model.InvoiceList{
		WhId:        &whID,
		RoDate:      &orderDate,
		PaymentType: entity.PAY_TYPE_CASH_ON_DELIVERY,
	}
	orderDetID := 123
	repo.details[roNo] = []model.InvoiceDetRead{{
		CustId:     custID,
		ProId:      1001,
		Qty1Final:  1,
		Qty2Final:  0,
		Qty3Final:  0,
		ConvUnit2:  1,
		ConvUnit3:  1,
		SellPrice1: 1000,
		OrderDetID: &orderDetID,
	}}

	dateNow := time.Now().Format("060102")
	repo.failOnceByInvoiceNo[fmt.Sprintf("INV%s%04d", dateNow, 1)] = true

	err := service.BulkUpdate(custID, entity.BulkUpdateInvoiceBody{
		Orders: []entity.UpdateInvoiceBody{{
			RoNo:         roNo,
			RoDate:       &roDate,
			ValDate:      &valDate,
			DeliveryDate: &deliveryDate,
		}},
	})
	if err != nil {
		t.Fatalf("bulk update should retry once and succeed: %v", err)
	}

	if len(repo.generated) != 1 {
		t.Fatalf("expected one persisted invoice number after retry success, got %d", len(repo.generated))
	}
}
