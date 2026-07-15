package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/config/env"
	"inventory/pkg/validation"
	"inventory/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockStockDisposalRepository struct {
	findProductByIDFn   func(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error)
	findWarehouseByIDFn func(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error)
	findSupplierByIDFn  func(ctx context.Context, supID int64, custID string) (*model.GrSupplier, error)
	getAvailableStockFn func(ctx context.Context, custID string, whID int64, proID int64) (float64, error)
	storeFn             func(ctx context.Context, stockDisposal *model.StockDisposal) error
	createDetailFn      func(ctx context.Context, detail *model.StockDisposalDetail) (*model.StockDisposalDetail, error)
	findByNumberFn      func(ctx context.Context, sdNumber string, custID, parentCustID, warehouseCustID string) (*model.StockDisposalList, error)
}

func (m *mockStockDisposalRepository) Store(ctx context.Context, stockDisposal *model.StockDisposal) error {
	if m.storeFn != nil {
		return m.storeFn(ctx, stockDisposal)
	}
	return nil
}

func (m *mockStockDisposalRepository) CreateDetail(ctx context.Context, detail *model.StockDisposalDetail) (*model.StockDisposalDetail, error) {
	if m.createDetailFn != nil {
		return m.createDetailFn(ctx, detail)
	}
	return detail, nil
}

func (m *mockStockDisposalRepository) FindByID(ctx context.Context, sdID int64, custID string) (*model.StockDisposal, error) {
	return &model.StockDisposal{SdID: sdID, SdNumber: "SD-DETAIL", CustID: custID}, nil
}

func (m *mockStockDisposalRepository) FindByNumber(ctx context.Context, sdNumber string, custID, parentCustID, warehouseCustID string) (*model.StockDisposalList, error) {
	if m.findByNumberFn != nil {
		return m.findByNumberFn(ctx, sdNumber, custID, parentCustID, warehouseCustID)
	}
	return &model.StockDisposalList{CustID: custID, SdNumber: sdNumber, WhID: 1, StockType: "GOOD", SubTotal: 100, VatValue: 11}, nil
}

func (m *mockStockDisposalRepository) FindDetail(ctx context.Context, sdID int64, custID string) ([]model.StockDisposalDetailList, error) {
	return nil, nil
}

func (m *mockStockDisposalRepository) FindAllByCustId(ctx context.Context, dataFilter entity.StockDisposalQueryFilter, custId, parentCustId string) ([]model.StockDisposalList, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockStockDisposalRepository) FindProductByID(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error) {
	if m.findProductByIDFn != nil {
		return m.findProductByIDFn(ctx, proID, custID, parentCustID)
	}
	return nil, nil
}

func (m *mockStockDisposalRepository) FindWarehouseByID(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error) {
	if m.findWarehouseByIDFn != nil {
		return m.findWarehouseByIDFn(ctx, whID, custID)
	}
	return &model.WarehouseStockWhList{CustID: custID, WhID: whID, StockType: "GOOD"}, nil
}

func (m *mockStockDisposalRepository) FindSupplierByID(ctx context.Context, supID int64, custID string) (*model.GrSupplier, error) {
	if m.findSupplierByIDFn != nil {
		return m.findSupplierByIDFn(ctx, supID, custID)
	}
	return &model.GrSupplier{}, nil
}

func (m *mockStockDisposalRepository) GetAvailableStock(ctx context.Context, custID string, whID int64, proID int64) (float64, error) {
	if m.getAvailableStockFn != nil {
		return m.getAvailableStockFn(ctx, custID, whID, proID)
	}
	return 0, nil
}

func (m *mockStockDisposalRepository) FindProductsForLookup(ctx context.Context, dataFilter entity.StockDisposalProductLookupQueryFilter, custId, parentCustId string) ([]model.StockDisposalProductLookup, int64, int, error) {
	return nil, 0, 0, nil
}

type mockStockRepository struct {
	stockUpdatesFn func(c context.Context, stockUpdates []*entity.StockUpdate) error
}

func (m *mockStockRepository) FindAllByCustId(dataFilter entity.StockQueryFilter, custId, parentCustId string) ([]model.Stock, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockStockRepository) Report(dataFilter entity.StockReportQueryFilter) ([]model.StockReport, int64, int, error) {
	return nil, 0, 0, nil
}

func (m *mockStockRepository) Store(c context.Context, data *model.Stock) error {
	return nil
}

func (m *mockStockRepository) StoreBulk(c context.Context, data []*model.Stock) error {
	return nil
}

func (m *mockStockRepository) StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error {
	if m.stockUpdatesFn != nil {
		return m.stockUpdatesFn(c, stockUpdates)
	}
	return nil
}

func (m *mockStockRepository) OpnameLookup(dataFilter entity.StockOpnameLookupQueryFilter) ([]model.StockOpnameLookup, int64, int, error) {
	return nil, 0, 0, nil
}

type mockTransaction struct{}

func (m *mockTransaction) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return tFunc(ctx)
}

type mockConfig struct{}

func (m *mockConfig) Get(key string) string {
	return ""
}

func newTestStockDisposalService(stockDisposalRepo repository.StockDisposalRepository, stockRepo repository.StockRepository) *stockDisposalServiceImpl {
	return NewStockDisposalService(
		stockDisposalRepo,
		stockRepo,
		&mockTransaction{},
		validation.NewValidator(),
		&mockConfig{},
	)
}

func TestStockDisposalValidateProducts_SupportsQty3Conversion(t *testing.T) {
	service := newTestStockDisposalService(&mockStockDisposalRepository{
		findProductByIDFn: func(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error) {
			return &model.Product{
				ProductId: proID,
				ConvUnit2: 10,
				ConvUnit3: 5,
			}, nil
		},
		getAvailableStockFn: func(ctx context.Context, custID string, whID int64, proID int64) (float64, error) {
			return 50, nil
		},
	}, &mockStockRepository{})

	request := entity.CreateStockDisposalBody{
		CustID:       "DIST01",
		ParentCustID: "PARENT01",
		WhID:         1,
		Products: []entity.CreateStockDisposalProductBody{{
			ProID:   10743,
			Qty3:    1,
			UnitID1: "PCS",
			UnitID2: "BOX",
			UnitID3: "CTN",
		}},
	}

	productMap, validationErrors := service.validateProducts(request, request.CustID, context.Background())

	require.Empty(t, validationErrors)
	require.Contains(t, productMap, int64(10743))
}

func TestStockDisposalStore_ProductNotFoundMessageIsClean(t *testing.T) {
	service := newTestStockDisposalService(&mockStockDisposalRepository{
		findWarehouseByIDFn: func(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error) {
			return &model.WarehouseStockWhList{CustID: custID, WhID: whID, StockType: "GOOD"}, nil
		},
		findSupplierByIDFn: func(ctx context.Context, supID int64, custID string) (*model.GrSupplier, error) {
			return &model.GrSupplier{}, nil
		},
		findProductByIDFn: func(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}, &mockStockRepository{})

	request := entity.CreateStockDisposalBody{
		CustID:       "DIST01",
		ParentCustID: "PARENT01",
		CreatedBy:    99,
		Date:         "2026-04-21",
		SupID:        7,
		WhID:         1,
		Products: []entity.CreateStockDisposalProductBody{{
			ProID:       10743,
			UnitID1:     "PCS",
			UnitID2:     "BOX",
			UnitID3:     "CTN",
			Qty1:        1,
			PurchPrice1: 100,
			PurchPrice2: 1000,
			PurchPrice3: 5000,
			GrossPrice:  100,
			Vat:         11,
			VatValue:    11,
			SubTotal:    111,
		}},
	}

	_, err := service.Store(request)

	require.Error(t, err)
	require.Contains(t, err.Error(), "product index 0 (pro_id: 10743): product not found")
	require.NotContains(t, err.Error(), "record not found")
}

func TestStockDisposalStore_SuccessWithValidProduct(t *testing.T) {
	stockUpdated := false
	service := newTestStockDisposalService(&mockStockDisposalRepository{
		findWarehouseByIDFn: func(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error) {
			return &model.WarehouseStockWhList{CustID: "C260020001", WhID: whID, StockType: "GOOD"}, nil
		},
		findSupplierByIDFn: func(ctx context.Context, supID int64, custID string) (*model.GrSupplier, error) {
			return &model.GrSupplier{}, nil
		},
		findProductByIDFn: func(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error) {
			return &model.Product{
				ProductId: proID,
				ConvUnit2: 10,
				ConvUnit3: 5,
			}, nil
		},
		getAvailableStockFn: func(ctx context.Context, custID string, whID int64, proID int64) (float64, error) {
			return 50, nil
		},
		storeFn: func(ctx context.Context, stockDisposal *model.StockDisposal) error {
			stockDisposal.SdID = 10
			stockDisposal.SdNumber = "SD260421001"
			require.Equal(t, "C260020001", stockDisposal.CustID)
			return nil
		},
		createDetailFn: func(ctx context.Context, detail *model.StockDisposalDetail) (*model.StockDisposalDetail, error) {
			detail.SdDetailID = 20
			return detail, nil
		},
		findByNumberFn: func(ctx context.Context, sdNumber string, custID, parentCustID, warehouseCustID string) (*model.StockDisposalList, error) {
			now := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)
			return &model.StockDisposalList{
				CustID:       custID,
				SdID:         10,
				SdNumber:     sdNumber,
				WhID:         1,
				StockType:    "GOOD",
				SubTotal:     5000,
				VatValue:     550,
				DisposalDate: &now,
			}, nil
		},
	}, &mockStockRepository{
		stockUpdatesFn: func(c context.Context, stockUpdates []*entity.StockUpdate) error {
			stockUpdated = true
			require.Len(t, stockUpdates, 1)
			require.Equal(t, "C260020001", stockUpdates[0].CustID)
			require.Equal(t, 50.0, stockUpdates[0].QtyOut)
			return nil
		},
	})

	request := entity.CreateStockDisposalBody{
		CustID:       "DIST01",
		ParentCustID: "PARENT01",
		CreatedBy:    99,
		Date:         "2026-04-21",
		SupID:        7,
		WhID:         1,
		Products: []entity.CreateStockDisposalProductBody{{
			ProID:       10743,
			UnitID1:     "PCS",
			UnitID2:     "BOX",
			UnitID3:     "CTN",
			Qty3:        1,
			PurchPrice1: 100,
			PurchPrice2: 1000,
			PurchPrice3: 5000,
			GrossPrice:  5000,
			Vat:         11,
			VatValue:    550,
			SubTotal:    5550,
		}},
	}

	response, err := service.Store(request)

	require.NoError(t, err)
	require.True(t, stockUpdated)
	require.Equal(t, "SD260421001", response.SdNumber)
	require.Equal(t, "2026-04-21", response.DisposalDate)
	require.Equal(t, 5550.0, response.Total)
}

func TestStockDisposalStore_InvalidWarehouseRelation(t *testing.T) {
	service := newTestStockDisposalService(&mockStockDisposalRepository{
		findWarehouseByIDFn: func(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}, &mockStockRepository{})

	request := entity.CreateStockDisposalBody{
		CustID:       "DIST01",
		ParentCustID: "PARENT01",
		CreatedBy:    99,
		Date:         "2026-04-21",
		SupID:        7,
		WhID:         999,
		Products: []entity.CreateStockDisposalProductBody{{
			ProID:       10743,
			UnitID1:     "PCS",
			UnitID2:     "BOX",
			UnitID3:     "CTN",
			Qty1:        1,
			PurchPrice1: 100,
			PurchPrice2: 1000,
			PurchPrice3: 5000,
			GrossPrice:  100,
			Vat:         11,
			VatValue:    11,
			SubTotal:    111,
		}},
	}

	_, err := service.Store(request)

	require.Error(t, err)
	require.Equal(t, "warehouse not found or inactive", err.Error())
}

var _ repository.StockDisposalRepository = (*mockStockDisposalRepository)(nil)
var _ repository.StockRepository = (*mockStockRepository)(nil)
var _ repository.Dbtransaction = (*mockTransaction)(nil)
var _ env.ConfigEnv = (*mockConfig)(nil)
