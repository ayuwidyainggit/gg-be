package service

import (
	"testing"

	"master/model"
	"master/repository"
)

// stubProductRepository is a minimal stub that implements repository.ProductRepository.
// Only FindDistributor and StoreDist are functional; all other methods panic to catch
// unexpected calls during StoreDistProduct / BulkStoreDistProduct tests.
type stubProductRepository struct {
	repository.ProductRepository // embed to satisfy the interface without listing every method

	findDistributorFn func(custId string) ([]model.MCustomer, error)
	storeDistFn       func(productsDist []model.ProductDistCreate) error
	storeDistCalls    [][]model.ProductDistCreate
}

func (s *stubProductRepository) FindDistributor(custId string) ([]model.MCustomer, error) {
	return s.findDistributorFn(custId)
}

func (s *stubProductRepository) StoreDist(productsDist []model.ProductDistCreate) error {
	s.storeDistCalls = append(s.storeDistCalls, productsDist)
	return s.storeDistFn(productsDist)
}

// The methods below are called by productServiceImpl.Store / BulkStore paths that we do NOT
// exercise in these tests. Embedding repository.ProductRepository (nil) would panic on any
// call, which is the desired behaviour — it surfaces unexpected code paths immediately.

func newStubProductRepo(distributors []model.MCustomer) *stubProductRepository {
	return &stubProductRepository{
		findDistributorFn: func(custId string) ([]model.MCustomer, error) {
			return distributors, nil
		},
		storeDistFn: func(productsDist []model.ProductDistCreate) error {
			return nil
		},
	}
}

func newProductServiceWithStub(stub *stubProductRepository) *productServiceImpl {
	return &productServiceImpl{
		ProductRepository: stub,
	}
}

// ---------------------------------------------------------------------------
// TestStoreDistProduct_PrincipalContext_SetsParentProId
// ---------------------------------------------------------------------------

func TestStoreDistProduct_PrincipalContext_SetsParentProId(t *testing.T) {
	distributors := []model.MCustomer{
		{CustId: "C260020001", CustName: "Distributor One", ParentCustId: "C26002"},
	}
	stub := newStubProductRepo(distributors)
	svc := newProductServiceWithStub(stub)

	const (
		distID    = "C26002"
		userID    = int64(1)
		productID = int64(100)
	)

	err := svc.StoreDistProduct(distID, userID, productID)
	if err != nil {
		t.Fatalf("StoreDistProduct returned unexpected error: %v", err)
	}

	// StoreDist must have been called exactly once
	if len(stub.storeDistCalls) != 1 {
		t.Fatalf("expected StoreDist called 1 time, got %d", len(stub.storeDistCalls))
	}

	items := stub.storeDistCalls[0]
	if len(items) != 1 {
		t.Fatalf("expected 1 item passed to StoreDist, got %d", len(items))
	}

	item := items[0]
	if item.ParentProId == nil {
		t.Fatal("expected ParentProId to be non-nil, got nil")
	}
	if *item.ParentProId != productID {
		t.Fatalf("expected ParentProId == %d, got %d", productID, *item.ParentProId)
	}
}

// ---------------------------------------------------------------------------
// TestBulkStoreDistProduct_PrincipalContext_SetsParentProId
// ---------------------------------------------------------------------------

func TestBulkStoreDistProduct_PrincipalContext_SetsParentProId(t *testing.T) {
	distributors := []model.MCustomer{
		{CustId: "C260020001", CustName: "Distributor One", ParentCustId: "C26002"},
	}
	stub := newStubProductRepo(distributors)
	svc := newProductServiceWithStub(stub)

	const (
		distID = "C26002"
		userID = int64(1)
	)
	productIDs := []int64{200, 201}

	err := svc.BulkStoreDistProduct(distID, userID, productIDs)
	if err != nil {
		t.Fatalf("BulkStoreDistProduct returned unexpected error: %v", err)
	}

	// StoreDist must have been called once per productID
	if len(stub.storeDistCalls) != len(productIDs) {
		t.Fatalf("expected StoreDist called %d times, got %d", len(productIDs), len(stub.storeDistCalls))
	}

	for i, pid := range productIDs {
		items := stub.storeDistCalls[i]
		if len(items) != 1 {
			t.Fatalf("call %d: expected 1 item, got %d", i, len(items))
		}
		item := items[0]
		if item.ParentProId == nil {
			t.Fatalf("call %d: expected ParentProId to be non-nil, got nil", i)
		}
		if *item.ParentProId != pid {
			t.Fatalf("call %d: expected ParentProId == %d, got %d", i, pid, *item.ParentProId)
		}
	}
}

// ---------------------------------------------------------------------------
// Compile-time guard: stubProductRepository must satisfy repository.ProductRepository.
// The embed already covers the full interface; these explicit overrides confirm the
// two methods we care about are correctly typed.
// ---------------------------------------------------------------------------

var _ repository.ProductRepository = (*stubProductRepository)(nil)

