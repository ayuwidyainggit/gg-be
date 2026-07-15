package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"testing"
	"time"
)

type depositRepositoryMock struct {
	repository.DepositRepository
	capturedFilter entity.DepositNumberListQueryFilter
	findFn         func(dataFilter entity.DepositNumberListQueryFilter) ([]model.DepositNumberList, int64, int, error)
}

func (m *depositRepositoryMock) FindDepositNumberListByCustId(dataFilter entity.DepositNumberListQueryFilter) ([]model.DepositNumberList, int64, int, error) {
	m.capturedFilter = dataFilter
	if m.findFn != nil {
		return m.findFn(dataFilter)
	}
	return []model.DepositNumberList{}, 0, 0, nil
}

func TestDepositService_ListDepositNumber_DefaultFilterAndMapping(t *testing.T) {
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	repoMock := &depositRepositoryMock{
		findFn: func(dataFilter entity.DepositNumberListQueryFilter) ([]model.DepositNumberList, int64, int, error) {
			return []model.DepositNumberList{
				{
					DepositNo:   "DEP-2026-001",
					CollectorID: 123,
					DepositDate: &now,
				},
			}, 1, 1, nil
		},
	}

	service := &DepositServiceImpl{DepositRepository: repoMock}

	result, total, lastPage, err := service.ListDepositNumber(entity.DepositNumberListQueryFilter{
		CustId:       "C001",
		CollectorIDs: []int{123},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.capturedFilter.Page != 1 {
		t.Fatalf("expected default page 1, got %d", repoMock.capturedFilter.Page)
	}
	if repoMock.capturedFilter.Limit != 20 {
		t.Fatalf("expected default limit 20, got %d", repoMock.capturedFilter.Limit)
	}
	if repoMock.capturedFilter.Sort != "created_date:desc" {
		t.Fatalf("expected default sort created_date:desc, got %s", repoMock.capturedFilter.Sort)
	}

	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if lastPage != 1 {
		t.Fatalf("expected lastPage 1, got %d", lastPage)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].DepositNo != "DEP-2026-001" {
		t.Fatalf("unexpected deposit_no: %s", result[0].DepositNo)
	}
	if result[0].CollectorID != 123 {
		t.Fatalf("unexpected collector_id: %d", result[0].CollectorID)
	}
	if result[0].DepositDate != "2026-02-01T00:00:00Z" {
		t.Fatalf("unexpected deposit_date format: %s", result[0].DepositDate)
	}
}

func TestDepositService_ListDepositNumber_MaxLimit(t *testing.T) {
	repoMock := &depositRepositoryMock{}
	service := &DepositServiceImpl{DepositRepository: repoMock}

	_, _, _, err := service.ListDepositNumber(entity.DepositNumberListQueryFilter{
		CustId:       "C001",
		CollectorIDs: []int{123},
		Limit:        10000,
		Page:         1,
		Sort:         "deposit_date:desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.capturedFilter.Limit != 9999 {
		t.Fatalf("expected max limit 9999, got %d", repoMock.capturedFilter.Limit)
	}
}

type depositRepositoryModeAwareMock struct {
	repository.DepositRepository
	calcCollectionPaidByInvoiceCalls int
}

func (m *depositRepositoryModeAwareMock) CountAllByCustId(custId string, depositDate string) (int, error) {
	return 0, nil
}

func (m *depositRepositoryModeAwareMock) StoreDetail(c context.Context, data *model.DepositDetail) (int, error) {
	return 1, nil
}

func (m *depositRepositoryModeAwareMock) Store(c context.Context, data *model.Deposit) error {
	return nil
}

func (m *depositRepositoryModeAwareMock) Update(c context.Context, depositNo string, custId string, data model.Deposit) error {
	return nil
}

func (m *depositRepositoryModeAwareMock) CalcCollectionPaidByInvoice(c context.Context, data *model.DepositDetail) error {
	m.calcCollectionPaidByInvoiceCalls++
	return nil
}

func (m *depositRepositoryModeAwareMock) CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error) {
	return 0, nil
}

type transactionPassThroughMock struct{}

func (m *transactionPassThroughMock) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return tFunc(ctx)
}

func TestDepositService_ModeAwareRecalc_StoreCollectionCallsCalcCollectionPaidByInvoice(t *testing.T) {
	repoMock := &depositRepositoryModeAwareMock{}
	service := &DepositServiceImpl{
		DepositRepository: repoMock,
		Transaction:       &transactionPassThroughMock{},
	}

	err := service.StoreCollection(entity.CreateDepositBodyByCollection{
		CustID:      "C001",
		DepositDate: "2026-02-01",
		Details: []entity.DepositDetail{
			{
				InvoiceNo:        "INV-001",
				InvoiceAmount:    100,
				TotalPayment:     10,
				RemainingPayment: 90,
				Payment:          []entity.DepositPayment{},
			},
		},
		Expense: []entity.DepositExpense{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.calcCollectionPaidByInvoiceCalls != 1 {
		t.Fatalf("expected CalcCollectionPaidByInvoice called once, got %d", repoMock.calcCollectionPaidByInvoiceCalls)
	}
}

func TestDepositService_ModeAwareRecalc_StoreInvoiceDoesNotCallCalcCollectionPaidByInvoice(t *testing.T) {
	repoMock := &depositRepositoryModeAwareMock{}
	service := &DepositServiceImpl{
		DepositRepository: repoMock,
		Transaction:       &transactionPassThroughMock{},
	}

	err := service.StoreInvoice(entity.CreateDepositBodyByInvoice{
		CustID:          "C001",
		DepositDate:     "2026-02-01",
		InvoiceDateFrom: "2026-02-01",
		InvoiceDateTo:   "2026-02-01",
		DueDateFrom:     "2026-02-01",
		DueDateTo:       "2026-02-01",
		Details: []entity.DepositDetail{
			{
				InvoiceNo:     "INV-001",
				InvoiceAmount: 100,
				TotalPayment:  10,
				Payment:       []entity.DepositPayment{},
			},
		},
		Expense: []entity.DepositExpense{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.calcCollectionPaidByInvoiceCalls != 1 {
		t.Fatalf("expected CalcCollectionPaidByInvoice called once for invoice flow, got %d", repoMock.calcCollectionPaidByInvoiceCalls)
	}
}
