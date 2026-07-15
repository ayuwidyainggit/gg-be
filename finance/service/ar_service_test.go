package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"testing"
)

type arRepositoryMock struct {
	repository.ArRepository
	countInvoicePaidAmountFn         func(invoiceNo string, custId string) (model.InvoicePaidAmount, error)
	storeCollectionFn                func(c context.Context, data *model.Collection) error
	storeCollectionDetailFn          func(c context.Context, data *model.CollectionDet) error
	updateCollectionFn               func(c context.Context, collectionNo string, data model.Collection) error
	updateCollectionRemainingFn      func(c context.Context, collectionNo string, custId string, remainingAmount float64) error
	deleteCollectionDetailNotInIDsFn func(c context.Context, collectionNo string, IDs []int64, custId string) error
	updateCollectionDetailFn         func(c context.Context, details *model.CollectionDet) error
	deleteAllCollectionDetailsFn     func(c context.Context, collectionNo string, custId string) error
}

func (m *arRepositoryMock) CountInvoicePaidAmount(invoiceNo string, custId string) (model.InvoicePaidAmount, error) {
	if m.countInvoicePaidAmountFn != nil {
		return m.countInvoicePaidAmountFn(invoiceNo, custId)
	}
	return model.InvoicePaidAmount{}, nil
}

func (m *arRepositoryMock) StoreCollection(c context.Context, data *model.Collection) error {
	if m.storeCollectionFn != nil {
		return m.storeCollectionFn(c, data)
	}
	return nil
}

func (m *arRepositoryMock) StoreCollectionDetail(c context.Context, data *model.CollectionDet) error {
	if m.storeCollectionDetailFn != nil {
		return m.storeCollectionDetailFn(c, data)
	}
	return nil
}

func (m *arRepositoryMock) UpdateCollection(c context.Context, collectionNo string, data model.Collection) error {
	if m.updateCollectionFn != nil {
		return m.updateCollectionFn(c, collectionNo, data)
	}
	return nil
}

func (m *arRepositoryMock) UpdateCollectionRemainingAmount(c context.Context, collectionNo string, custId string, remainingAmount float64) error {
	if m.updateCollectionRemainingFn != nil {
		return m.updateCollectionRemainingFn(c, collectionNo, custId, remainingAmount)
	}
	return nil
}

func (m *arRepositoryMock) DeleteCollectionDetailNotInIDs(c context.Context, collectionNo string, IDs []int64, custId string) error {
	if m.deleteCollectionDetailNotInIDsFn != nil {
		return m.deleteCollectionDetailNotInIDsFn(c, collectionNo, IDs, custId)
	}
	return nil
}

func (m *arRepositoryMock) UpdateCollectionDetail(c context.Context, details *model.CollectionDet) error {
	if m.updateCollectionDetailFn != nil {
		return m.updateCollectionDetailFn(c, details)
	}
	return nil
}

func (m *arRepositoryMock) DeleteAllCollectionDetails(c context.Context, collectionNo string, custId string) error {
	if m.deleteAllCollectionDetailsFn != nil {
		return m.deleteAllCollectionDetailsFn(c, collectionNo, custId)
	}
	return nil
}

func TestMapCollectionDetailAmounts(t *testing.T) {
	testCases := []struct {
		name                           string
		detail                         model.CollectionDetList
		expectedInvoiceAmount          float64
		expectedPaidAmount             float64
		expectedRemainingAmount        float64
		expectedInvoicePayment         float64
		expectedPaidAmountByCollection float64
	}{
		{
			name: "all values mapped directly",
			detail: model.CollectionDetList{
				InvoiceAmount:       float64Ptr(2327500),
				RemainingAmount:     float64Ptr(700000),
				PaidAmount:          float64Ptr(500000),
				TotalInvoicePayment: float64Ptr(2127500),
			},
			expectedInvoiceAmount:          2327500,
			expectedPaidAmount:             500000,
			expectedRemainingAmount:        700000,
			expectedInvoicePayment:         2127500,
			expectedPaidAmountByCollection: 500000,
		},
		{
			name: "nil detail values treated as zero",
			detail: model.CollectionDetList{
				InvoiceAmount:       nil,
				RemainingAmount:     nil,
				PaidAmount:          nil,
				TotalInvoicePayment: nil,
			},
			expectedInvoiceAmount:          0,
			expectedPaidAmount:             0,
			expectedRemainingAmount:        0,
			expectedInvoicePayment:         0,
			expectedPaidAmountByCollection: 0,
		},
		{
			name: "remaining amount is not derived from header total",
			detail: model.CollectionDetList{
				InvoiceAmount:       float64Ptr(2548000),
				RemainingAmount:     float64Ptr(2548000),
				PaidAmount:          float64Ptr(0),
				TotalInvoicePayment: float64Ptr(0),
			},
			expectedInvoiceAmount:          2548000,
			expectedPaidAmount:             0,
			expectedRemainingAmount:        2548000,
			expectedInvoicePayment:         0,
			expectedPaidAmountByCollection: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			invoiceAmount, paidAmount, remainingAmount, invoicePayment, paidAmountByCollection := mapCollectionDetailAmounts(testCase.detail)

			if invoiceAmount != testCase.expectedInvoiceAmount {
				t.Fatalf("invoiceAmount mismatch: expected %v, got %v", testCase.expectedInvoiceAmount, invoiceAmount)
			}

			if paidAmount != testCase.expectedPaidAmount {
				t.Fatalf("paidAmount mismatch: expected %v, got %v", testCase.expectedPaidAmount, paidAmount)
			}

			if remainingAmount != testCase.expectedRemainingAmount {
				t.Fatalf("remainingAmount mismatch: expected %v, got %v", testCase.expectedRemainingAmount, remainingAmount)
			}

			if invoicePayment != testCase.expectedInvoicePayment {
				t.Fatalf("invoicePayment mismatch: expected %v, got %v", testCase.expectedInvoicePayment, invoicePayment)
			}

			if paidAmountByCollection != testCase.expectedPaidAmountByCollection {
				t.Fatalf("paidAmountByCollection mismatch: expected %v, got %v", testCase.expectedPaidAmountByCollection, paidAmountByCollection)
			}
		})
	}
}

func TestCalculateTotalInvoicePayment(t *testing.T) {
	testCases := []struct {
		name     string
		details  []model.CollectionDetList
		expected float64
	}{
		{
			name: "sums invoice payment from all details",
			details: []model.CollectionDetList{
				{TotalInvoicePayment: float64Ptr(2128500)},
				{TotalInvoicePayment: float64Ptr(0)},
			},
			expected: 2128500,
		},
		{
			name: "ignores nil invoice payment values",
			details: []model.CollectionDetList{
				{TotalInvoicePayment: float64Ptr(1000)},
				{TotalInvoicePayment: nil},
				{TotalInvoicePayment: float64Ptr(2000)},
			},
			expected: 3000,
		},
		{
			name:     "empty details returns zero",
			details:  []model.CollectionDetList{},
			expected: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := calculateTotalInvoicePayment(testCase.details)
			if result != testCase.expected {
				t.Fatalf("totalInvoicePayment mismatch: expected %v, got %v", testCase.expected, result)
			}
		})
	}
}

func TestCalculateTotalInvoicePaymentFromResponse(t *testing.T) {
	details := []entity.CollectionDetResponse{
		{InvoicePayment: float64Ptr(2128500)},
		{InvoicePayment: float64Ptr(0)},
	}

	result := calculateTotalInvoicePaymentFromResponse(details)
	if result != 2128500 {
		t.Fatalf("totalInvoicePaymentFromResponse mismatch: expected %v, got %v", 2128500, result)
	}
}

func TestCalculateTotalRemainingAmountFromResponse(t *testing.T) {
	details := []entity.CollectionDetResponse{
		{RemainingAmount: float64Ptr(199000)},
		{RemainingAmount: float64Ptr(2327500)},
	}

	result := calculateTotalRemainingAmountFromResponse(details)
	if result != 2526500 {
		t.Fatalf("totalRemainingAmountFromResponse mismatch: expected %v, got %v", 2526500, result)
	}
}

func TestStoreCollection_RecalculatesRemainingAmountsFromServerSideFormula(t *testing.T) {
	var storedHeader model.Collection
	var storedDetails []model.CollectionDet
	var updatedRemainingAmount float64

	repoMock := &arRepositoryMock{
		countInvoicePaidAmountFn: func(invoiceNo string, custId string) (model.InvoicePaidAmount, error) {
			switch invoiceNo {
			case "INV2603310011":
				return model.InvoicePaidAmount{PaidAmount: 300000}, nil
			case "INV2603300005":
				return model.InvoicePaidAmount{PaidAmount: 480750}, nil
			default:
				return model.InvoicePaidAmount{}, nil
			}
		},
		storeCollectionFn: func(c context.Context, data *model.Collection) error {
			data.CollectionNo = "CL-TEST-001"
			storedHeader = *data
			return nil
		},
		storeCollectionDetailFn: func(c context.Context, data *model.CollectionDet) error {
			storedDetails = append(storedDetails, *data)
			return nil
		},
		updateCollectionRemainingFn: func(c context.Context, collectionNo string, custId string, remainingAmount float64) error {
			updatedRemainingAmount = remainingAmount
			return nil
		},
	}

	service := &arServiceImpl{Repository: repoMock, Transaction: &transactionPassThroughMock{}}

	err := service.StoreCollection(entity.CreateCollectionBody{
		CustID:          "C220010001",
		CollectionDate:  stringPtr("2026-03-31"),
		CreatedBy:       int64Ptr(99),
		RemainingAmount: float64Ptr(3465000),
		Details: []entity.CreateCollectionDetBody{
			{
				InvoiceNo:       "INV2603310011",
				InvoiceAmount:   float64Ptr(2109000),
				PaidAmount:      float64Ptr(809000),
				RemainingAmount: float64Ptr(1809000),
			},
			{
				InvoiceNo:       "INV2603300005",
				InvoiceAmount:   float64Ptr(2136750),
				PaidAmount:      float64Ptr(0),
				RemainingAmount: float64Ptr(1656000),
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if storedHeader.CollectionNo != "CL-TEST-001" {
		t.Fatalf("expected stored collection number to be assigned, got %s", storedHeader.CollectionNo)
	}

	if len(storedDetails) != 2 {
		t.Fatalf("expected 2 stored details, got %d", len(storedDetails))
	}

	if storedDetails[0].RemainingAmount != 1000000 {
		t.Fatalf("expected first detail remaining amount 1000000, got %v", storedDetails[0].RemainingAmount)
	}

	if storedDetails[1].RemainingAmount != 1656000 {
		t.Fatalf("expected second detail remaining amount 1656000, got %v", storedDetails[1].RemainingAmount)
	}

	if updatedRemainingAmount != 2656000 {
		t.Fatalf("expected header remaining amount 2656000, got %v", updatedRemainingAmount)
	}
}

func TestUpdateCollection_RecalculatesRemainingAmountsFromServerSideFormula(t *testing.T) {
	var updatedDetails []model.CollectionDet
	var updatedRemainingAmount float64

	repoMock := &arRepositoryMock{
		countInvoicePaidAmountFn: func(invoiceNo string, custId string) (model.InvoicePaidAmount, error) {
			switch invoiceNo {
			case "INV2603310011":
				return model.InvoicePaidAmount{PaidAmount: 300000}, nil
			case "INV2603300005":
				return model.InvoicePaidAmount{PaidAmount: 480750}, nil
			default:
				return model.InvoicePaidAmount{}, nil
			}
		},
		updateCollectionFn: func(c context.Context, collectionNo string, data model.Collection) error {
			return nil
		},
		deleteCollectionDetailNotInIDsFn: func(c context.Context, collectionNo string, IDs []int64, custId string) error {
			return nil
		},
		updateCollectionDetailFn: func(c context.Context, details *model.CollectionDet) error {
			updatedDetails = append(updatedDetails, *details)
			return nil
		},
		updateCollectionRemainingFn: func(c context.Context, collectionNo string, custId string, remainingAmount float64) error {
			updatedRemainingAmount = remainingAmount
			return nil
		},
	}

	service := &arServiceImpl{Repository: repoMock, Transaction: &transactionPassThroughMock{}}

	err := service.UpdateCollection("CL2603310001", entity.UpdateCollectionBody{
		CustID:          "C220010001",
		CollectionDate:  stringPtr("2026-03-31"),
		UpdatedBy:       99,
		RemainingAmount: float64Ptr(3465000),
		Details: []entity.UpdateCollectionDetBody{
			{
				CollectionDetID: int64Ptr(1),
				InvoiceNo:       stringPtr("INV2603310011"),
				InvoiceAmount:   float64Ptr(2109000),
				PaidAmount:      float64Ptr(809000),
				RemainingAmount: float64Ptr(1809000),
			},
			{
				CollectionDetID: int64Ptr(2),
				InvoiceNo:       stringPtr("INV2603300005"),
				InvoiceAmount:   float64Ptr(2136750),
				PaidAmount:      float64Ptr(0),
				RemainingAmount: float64Ptr(1656000),
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(updatedDetails) != 2 {
		t.Fatalf("expected 2 updated details, got %d", len(updatedDetails))
	}

	if updatedDetails[0].RemainingAmount != 1000000 {
		t.Fatalf("expected first updated detail remaining amount 1000000, got %v", updatedDetails[0].RemainingAmount)
	}

	if updatedDetails[1].RemainingAmount != 1656000 {
		t.Fatalf("expected second updated detail remaining amount 1656000, got %v", updatedDetails[1].RemainingAmount)
	}

	if updatedRemainingAmount != 2656000 {
		t.Fatalf("expected updated header remaining amount 2656000, got %v", updatedRemainingAmount)
	}
}

func float64Ptr(value float64) *float64 {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}

func stringPtr(value string) *string {
	return &value
}
