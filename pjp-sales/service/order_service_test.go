package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"sales/entity"
	"sales/model"
	"sales/repository"
)

type mockOrderRepository struct {
	repository.OrderRepository
	findByNoFn                       func(roNo string, custID string) (model.OrderList, error)
	findOrderDetailByDetailIDFn      func(detailID int64, custID string) (model.OrderDetailRead, error)
	findOrderDetailsByIDsFn          func(detailIDs []int64, custID string) ([]model.OrderDetailRead, error)
	findProductByIDFn                func(productID int) (model.ProductRead, error)
	findOutletByIDFn                 func(outletID int, custID string, parentCustID string) (model.OutletRead, error)
	findDiscountByProductAndOutletFn func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error)
	findDiscountCriteriaBySubTotalFn func(discountID string, subTotal int) (model.DiscountCriteria, error)
	updateDetailPartialFn            func(c context.Context, orderDetailId int64, custID string, updates map[string]interface{}) error
	findOrderDetailsForProformaFn    func(ctx context.Context, roNos []string, custID string) ([]model.OrderDetailRead, error)
	updateFn                         func(c context.Context, roNo, custID string, data model.Order) error
	storeDetailFn                    func(c context.Context, data *model.OrderDetail) error
	storeRewardFn                    func(c context.Context, data *model.OrderReward) error
	lockOrderByScopeFn               func(ctx context.Context, custId string, roDates []time.Time) error
	deleteOrderDetailByScopeFn       func(ctx context.Context, custId string, roDates []time.Time) (int64, error)
	deleteOrderByScopeFn             func(ctx context.Context, custId string, roDates []time.Time) (int64, error)
	storeFn                          func(c context.Context, data *model.Order) error
	deletePromoDetailsFn             func(c context.Context, roNo string, custID string) error
	deleteRewardsFn                  func(c context.Context, roNo string, custID string) error
	findDetailFn                     func(roNo string, custID string) ([]model.OrderDetailRead, error)
	findDetailByNotInDetailIDsFn     func(detailIDs []int64, roNo string, custID string) ([]model.OrderDetailRead, error)
	deleteDetailNotInIDsFn           func(c context.Context, roNo string, custID string, ids []int64) error
	findDetailByDetailIDFn           func(detailID int64, roNo string, custID string) (model.OrderDetailRead, error)
	updateDetailFn                   func(c context.Context, details *model.OrderDetail) error
	syncFinalOrderFieldsFn           func(c context.Context, orderDetailId int64) error
}

func (m *mockOrderRepository) FindByNo(roNo string, custID string) (model.OrderList, error) {
	return m.findByNoFn(roNo, custID)
}

func (m *mockOrderRepository) FindOrderDetailByDetailID(detailID int64, custID string) (model.OrderDetailRead, error) {
	return m.findOrderDetailByDetailIDFn(detailID, custID)
}

func (m *mockOrderRepository) FindOrderDetailsByIDs(detailIDs []int64, custID string) ([]model.OrderDetailRead, error) {
	if m.findOrderDetailsByIDsFn != nil {
		return m.findOrderDetailsByIDsFn(detailIDs, custID)
	}

	if m.findOrderDetailByDetailIDFn == nil {
		return []model.OrderDetailRead{}, nil
	}

	details := make([]model.OrderDetailRead, 0, len(detailIDs))
	for _, detailID := range detailIDs {
		detail, err := m.findOrderDetailByDetailIDFn(detailID, custID)
		if err != nil {
			return nil, err
		}

		details = append(details, detail)
	}

	return details, nil
}

func (m *mockOrderRepository) FindProductByID(productID int) (model.ProductRead, error) {
	return m.findProductByIDFn(productID)
}

func (m *mockOrderRepository) FindOutletByID(outletID int, custID string, parentCustID string) (model.OutletRead, error) {
	return m.findOutletByIDFn(outletID, custID, parentCustID)
}

func (m *mockOrderRepository) FindDiscountByProductAndOutlet(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
	return m.findDiscountByProductAndOutletFn(product, outlet)
}

func (m *mockOrderRepository) FindDiscountCriteriaBySubTotal(discountID string, subTotal int) (model.DiscountCriteria, error) {
	return m.findDiscountCriteriaBySubTotalFn(discountID, subTotal)
}

func (m *mockOrderRepository) UpdateDetailPartial(c context.Context, orderDetailId int64, custID string, updates map[string]interface{}) error {
	return m.updateDetailPartialFn(c, orderDetailId, custID, updates)
}

func (m *mockOrderRepository) FindOrderDetailsForProforma(ctx context.Context, roNos []string, custID string) ([]model.OrderDetailRead, error) {
	return m.findOrderDetailsForProformaFn(ctx, roNos, custID)
}

func (m *mockOrderRepository) Update(c context.Context, roNo, custID string, data model.Order) error {
	return m.updateFn(c, roNo, custID, data)
}

func (m *mockOrderRepository) Store(c context.Context, data *model.Order) error {
	if m.storeFn != nil {
		return m.storeFn(c, data)
	}
	return nil
}

func (m *mockOrderRepository) StoreDetail(c context.Context, data *model.OrderDetail) error {
	if m.storeDetailFn != nil {
		return m.storeDetailFn(c, data)
	}
	return nil
}

func (m *mockOrderRepository) LockOrderByScope(ctx context.Context, custId string, roDates []time.Time) error {
	if m.lockOrderByScopeFn != nil {
		return m.lockOrderByScopeFn(ctx, custId, roDates)
	}
	return nil
}

func (m *mockOrderRepository) DeleteOrderDetailByScope(ctx context.Context, custId string, roDates []time.Time) (int64, error) {
	if m.deleteOrderDetailByScopeFn != nil {
		return m.deleteOrderDetailByScopeFn(ctx, custId, roDates)
	}
	return 0, nil
}

func (m *mockOrderRepository) DeleteOrderByScope(ctx context.Context, custId string, roDates []time.Time) (int64, error) {
	if m.deleteOrderByScopeFn != nil {
		return m.deleteOrderByScopeFn(ctx, custId, roDates)
	}
	return 0, nil
}

func (m *mockOrderRepository) StoreReward(c context.Context, data *model.OrderReward) error {
	if m.storeRewardFn != nil {
		return m.storeRewardFn(c, data)
	}
	return nil
}

func (m *mockOrderRepository) DeletePromoDetails(c context.Context, roNo string, custID string) error {
	if m.deletePromoDetailsFn != nil {
		return m.deletePromoDetailsFn(c, roNo, custID)
	}
	return nil
}

func (m *mockOrderRepository) DeleteRewards(c context.Context, roNo string, custID string) error {
	if m.deleteRewardsFn != nil {
		return m.deleteRewardsFn(c, roNo, custID)
	}
	return nil
}

func (m *mockOrderRepository) FindDetail(roNo string, custID string) ([]model.OrderDetailRead, error) {
	if m.findDetailFn != nil {
		return m.findDetailFn(roNo, custID)
	}
	return []model.OrderDetailRead{}, nil
}

func (m *mockOrderRepository) FindDetailByNotInDetailIDs(detailIDs []int64, roNo string, custID string) ([]model.OrderDetailRead, error) {
	if m.findDetailByNotInDetailIDsFn != nil {
		return m.findDetailByNotInDetailIDsFn(detailIDs, roNo, custID)
	}
	return []model.OrderDetailRead{}, nil
}

func (m *mockOrderRepository) DeleteDetailNotInIDs(c context.Context, roNo string, custID string, ids []int64) error {
	if m.deleteDetailNotInIDsFn != nil {
		return m.deleteDetailNotInIDsFn(c, roNo, custID, ids)
	}
	return nil
}

func (m *mockOrderRepository) FindDetailByDetailID(detailID int64, roNo string, custID string) (model.OrderDetailRead, error) {
	if m.findDetailByDetailIDFn != nil {
		return m.findDetailByDetailIDFn(detailID, roNo, custID)
	}
	return model.OrderDetailRead{}, nil
}

func (m *mockOrderRepository) UpdateDetail(c context.Context, details *model.OrderDetail) error {
	if m.updateDetailFn != nil {
		return m.updateDetailFn(c, details)
	}
	return nil
}

func (m *mockOrderRepository) SyncFinalOrderFields(c context.Context, orderDetailId int64) error {
	if m.syncFinalOrderFieldsFn != nil {
		return m.syncFinalOrderFieldsFn(c, orderDetailId)
	}
	return nil
}

type mockOrderRepositoryDetailV2 struct {
	repository.OrderRepository
	findByNoFn                                 func(roNo string, custID string) (model.OrderList, error)
	findByNoNoCustIDFn                         func(roNo string, custIDOrigin string) (model.OrderList, error)
	findDetailFn                               func(roNo string, custID string) ([]model.OrderDetailRead, error)
	findDetailNoCustIDFn                       func(roNo string, custIDOrigin string) ([]model.OrderDetailRead, error)
	findRewardFn                               func(roNo string, custID string) ([]model.OrderRewardRead, error)
	findRewardNoCustIDFn                       func(roNo string, custIDOrigin string) ([]model.OrderRewardRead, error)
	findOrderApprovalRequestDetailByRoAndEmpFn func(roNo string, empID int64) (model.OrderApprovalRequestDetailRead, error)
	findWarehouseStockByWhIdAndProIds          func(custID string, whID int64, proIDs []int64) (map[int64]float64, error)
	findProductByIDFn                          func(productID int) (model.ProductRead, error)
	findProductByListIDFn                      func(productIDs []int64) ([]model.Product, error)
}

func (m *mockOrderRepositoryDetailV2) FindByNo(roNo string, custID string) (model.OrderList, error) {
	return m.findByNoFn(roNo, custID)
}

func (m *mockOrderRepositoryDetailV2) FindByNoNoCustID(roNo string, custIDOrigin string) (model.OrderList, error) {
	if m.findByNoNoCustIDFn != nil {
		return m.findByNoNoCustIDFn(roNo, custIDOrigin)
	}
	return model.OrderList{}, nil
}

func (m *mockOrderRepositoryDetailV2) FindDetail(roNo string, custID string) ([]model.OrderDetailRead, error) {
	return m.findDetailFn(roNo, custID)
}

func (m *mockOrderRepositoryDetailV2) FindDetailNoCustID(roNo string, custIDOrigin string) ([]model.OrderDetailRead, error) {
	if m.findDetailNoCustIDFn != nil {
		return m.findDetailNoCustIDFn(roNo, custIDOrigin)
	}
	return []model.OrderDetailRead{}, nil
}

func (m *mockOrderRepositoryDetailV2) FindReward(roNo string, custID string) ([]model.OrderRewardRead, error) {
	return m.findRewardFn(roNo, custID)
}

func (m *mockOrderRepositoryDetailV2) FindRewardNoCustID(roNo string, custIDOrigin string) ([]model.OrderRewardRead, error) {
	if m.findRewardNoCustIDFn != nil {
		return m.findRewardNoCustIDFn(roNo, custIDOrigin)
	}
	return []model.OrderRewardRead{}, nil
}

func (m *mockOrderRepositoryDetailV2) FindOrderApprovalRequestDetailByRoAndEmp(roNo string, empID int64) (model.OrderApprovalRequestDetailRead, error) {
	if m.findOrderApprovalRequestDetailByRoAndEmpFn != nil {
		return m.findOrderApprovalRequestDetailByRoAndEmpFn(roNo, empID)
	}
	return model.OrderApprovalRequestDetailRead{}, errors.New("not found")
}

func (m *mockOrderRepositoryDetailV2) FindWarehouseStockByWhIdAndProIds(custID string, whID int64, proIDs []int64) (map[int64]float64, error) {
	return m.findWarehouseStockByWhIdAndProIds(custID, whID, proIDs)
}

func (m *mockOrderRepositoryDetailV2) FindProductByID(productID int) (model.ProductRead, error) {
	return m.findProductByIDFn(productID)
}

func (m *mockOrderRepositoryDetailV2) FindProductByListID(productIDs []int64) ([]model.Product, error) {
	if m.findProductByListIDFn != nil {
		return m.findProductByListIDFn(productIDs)
	}

	if m.findProductByIDFn == nil {
		return []model.Product{}, nil
	}

	products := make([]model.Product, 0, len(productIDs))
	for _, productID := range productIDs {
		productRead, err := m.findProductByIDFn(int(productID))
		if err != nil {
			return nil, err
		}

		products = append(products, model.Product{
			ProductId: int64(productRead.ProId),
			UnitId1:   productRead.UnitId1,
			UnitId2:   productRead.UnitId2,
			UnitId3:   productRead.UnitId3,
			UnitId4:   stringPtr(productRead.UnitId4),
			UnitId5:   stringPtr(productRead.UnitId5),
		})
	}

	return products, nil
}

type mockOrderRepositoryStore struct {
	repository.OrderRepository
	countAllRoByCustIdFn                func(custId string, roDate string) (int, error)
	findOutletByIDFn                    func(outletID int, custId string, parentCustId string) (model.OutletRead, error)
	findProductByIDFn                   func(productID int) (model.ProductRead, error)
	findDiscountByProductAndOutletFn    func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error)
	findDiscountCriteriaBySubTotalFn    func(discountID string, subTotal int) (model.DiscountCriteria, error)
	storeFn                             func(c context.Context, data *model.Order) error
	storeDetailFn                       func(c context.Context, data *model.OrderDetail) error
	storeRewardFn                       func(c context.Context, data *model.OrderReward) error
	findWarehouseStockByWhIdAndProIdsFn func(custID string, whID int64, proIDs []int64) (map[int64]float64, error)
}

func (m *mockOrderRepositoryStore) CountAllRoByCustId(custId string, roDate string) (int, error) {
	return m.countAllRoByCustIdFn(custId, roDate)
}

func (m *mockOrderRepositoryStore) FindOutletByID(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
	return m.findOutletByIDFn(outletID, custId, parentCustId)
}

func (m *mockOrderRepositoryStore) FindProductByID(productID int) (model.ProductRead, error) {
	return m.findProductByIDFn(productID)
}

func (m *mockOrderRepositoryStore) FindDiscountByProductAndOutlet(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
	return m.findDiscountByProductAndOutletFn(product, outlet)
}

func (m *mockOrderRepositoryStore) FindDiscountCriteriaBySubTotal(discountID string, subTotal int) (model.DiscountCriteria, error) {
	return m.findDiscountCriteriaBySubTotalFn(discountID, subTotal)
}

func (m *mockOrderRepositoryStore) Store(c context.Context, data *model.Order) error {
	return m.storeFn(c, data)
}

func (m *mockOrderRepositoryStore) StoreDetail(c context.Context, data *model.OrderDetail) error {
	return m.storeDetailFn(c, data)
}

func (m *mockOrderRepositoryStore) StoreReward(c context.Context, data *model.OrderReward) error {
	return m.storeRewardFn(c, data)
}

func (m *mockOrderRepositoryStore) FindWarehouseStockByWhIdAndProIds(custID string, whID int64, proIDs []int64) (map[int64]float64, error) {
	if m.findWarehouseStockByWhIdAndProIdsFn != nil {
		return m.findWarehouseStockByWhIdAndProIdsFn(custID, whID, proIDs)
	}
	stock := make(map[int64]float64, len(proIDs))
	for _, proID := range proIDs {
		stock[proID] = 1_000_000_000
	}
	return stock, nil
}

type mockPromotionRepositoryEnhance struct {
	repository.PromotionRepository
	findProductByIDAndCustIDFn func(productID int64, custID string) (model.ProductRead, error)
	findProductAndPriceByIDFn  func(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error)
}

func (m *mockPromotionRepositoryEnhance) FindProductByIDAndCustID(productID int64, custID string) (model.ProductRead, error) {
	return m.findProductByIDAndCustIDFn(productID, custID)
}

func (m *mockPromotionRepositoryEnhance) FindProductAndPriceByID(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error) {
	if m.findProductAndPriceByIDFn != nil {
		return m.findProductAndPriceByIDFn(productID, distributorID, transDate, custID, parentCustID)
	}
	return model.Product{ProductId: productID, SellPrice1: 50, SellPrice2: 100, SellPrice3: 200, ConvUnit2: 10, ConvUnit3: 5}, nil
}

type mockPromotionRepositoryDetailV2 struct {
	repository.PromotionRepository
	findProductByIDAndCustIDFn func(productID int64, custID string) (model.ProductRead, error)
	findProductAndPriceByIDFn  func(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error)
}

func (m *mockPromotionRepositoryDetailV2) FindProductByIDAndCustID(productID int64, custID string) (model.ProductRead, error) {
	return m.findProductByIDAndCustIDFn(productID, custID)
}

func (m *mockPromotionRepositoryDetailV2) FindProductAndPriceByID(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error) {
	if m.findProductAndPriceByIDFn != nil {
		return m.findProductAndPriceByIDFn(productID, distributorID, transDate, custID, parentCustID)
	}
	return model.Product{ProductId: productID, SellPrice1: 50, SellPrice2: 100, SellPrice3: 200, ConvUnit2: 10, ConvUnit3: 5}, nil
}

type mockPromotionRepositoryStore struct {
	repository.PromotionRepository
	findProductByIDAndCustIDFn func(productID int64, custID string) (model.ProductRead, error)
	findProductAndPriceByIDFn  func(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error)
}

func (m *mockPromotionRepositoryStore) FindProductByIDAndCustID(productID int64, custID string) (model.ProductRead, error) {
	return m.findProductByIDAndCustIDFn(productID, custID)
}

func (m *mockPromotionRepositoryStore) FindProductAndPriceByID(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error) {
	if m.findProductAndPriceByIDFn != nil {
		return m.findProductAndPriceByIDFn(productID, distributorID, transDate, custID, parentCustID)
	}
	return model.Product{ProductId: productID, SellPrice1: 50, SellPrice2: 100, SellPrice3: 200, ConvUnit2: 10, ConvUnit3: 5}, nil
}

type mockPromotionV2RepositoryEnhance struct {
	repository.PromotionV2Repository
	findOutletByIDFn                 func(outletID int64, custID string) (model.OutletPromo, error)
	findSalesmanByIDFn               func(salesmanID int64, custID string) (model.SalesmanPromo, error)
	findWarehouseByIDFn              func(warehouseID int64, custID string) (model.WarehousePromo, error)
	findActivePromotionsByOutletFn   func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error)
	findProductCriteriasByPromoIDsFn func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error)
	findSlabsByPromoIDsFn            func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error)
	findStratasByPromoIDsFn          func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error)
	getAllRewardProductFromStockV2Fn func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error)
}

func (m *mockPromotionV2RepositoryEnhance) FindOutletByID(outletID int64, custID string) (model.OutletPromo, error) {
	return m.findOutletByIDFn(outletID, custID)
}

func (m *mockPromotionV2RepositoryEnhance) FindSalesmanByID(salesmanID int64, custID string) (model.SalesmanPromo, error) {
	return m.findSalesmanByIDFn(salesmanID, custID)
}

func (m *mockPromotionV2RepositoryEnhance) FindWarehouseByID(warehouseID int64, custID string) (model.WarehousePromo, error) {
	return m.findWarehouseByIDFn(warehouseID, custID)
}

func (m *mockPromotionV2RepositoryEnhance) FindActivePromotionsByOutletCriteria(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
	return m.findActivePromotionsByOutletFn(req, outlet, salesman)
}

func (m *mockPromotionV2RepositoryEnhance) FindProductCriteriasByPromoIDs(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) {
	return m.findProductCriteriasByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryEnhance) FindSlabsByPromoIDs(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
	return m.findSlabsByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryEnhance) FindStratasByPromoIDs(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) {
	return m.findStratasByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryEnhance) GetAllRewardProductFromStockV2(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
	return m.getAllRewardProductFromStockV2Fn(req, rewardCtx)
}

type mockPromotionV2RepositoryDetailV2 struct {
	repository.PromotionV2Repository
	findOutletByIDFn                 func(outletID int64, custID string) (model.OutletPromo, error)
	findSalesmanByIDFn               func(salesmanID int64, custID string) (model.SalesmanPromo, error)
	findWarehouseByIDFn              func(warehouseID int64, custID string) (model.WarehousePromo, error)
	findActivePromotionsByOutletFn   func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error)
	findProductCriteriasByPromoIDsFn func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error)
	findSlabsByPromoIDsFn            func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error)
	findStratasByPromoIDsFn          func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error)
	getAllRewardProductFromStockV2Fn func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error)
}

func (m *mockPromotionV2RepositoryDetailV2) FindOutletByID(outletID int64, custID string) (model.OutletPromo, error) {
	return m.findOutletByIDFn(outletID, custID)
}

func (m *mockPromotionV2RepositoryDetailV2) FindSalesmanByID(salesmanID int64, custID string) (model.SalesmanPromo, error) {
	return m.findSalesmanByIDFn(salesmanID, custID)
}

func (m *mockPromotionV2RepositoryDetailV2) FindWarehouseByID(warehouseID int64, custID string) (model.WarehousePromo, error) {
	return m.findWarehouseByIDFn(warehouseID, custID)
}

func (m *mockPromotionV2RepositoryDetailV2) FindActivePromotionsByOutletCriteria(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
	return m.findActivePromotionsByOutletFn(req, outlet, salesman)
}

func (m *mockPromotionV2RepositoryDetailV2) FindProductCriteriasByPromoIDs(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) {
	return m.findProductCriteriasByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryDetailV2) FindSlabsByPromoIDs(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
	return m.findSlabsByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryDetailV2) FindStratasByPromoIDs(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) {
	return m.findStratasByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryDetailV2) GetAllRewardProductFromStockV2(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
	return m.getAllRewardProductFromStockV2Fn(req, rewardCtx)
}

type mockPromotionV2RepositoryStore struct {
	repository.PromotionV2Repository
	findOutletByIDFn                 func(outletID int64, custID string) (model.OutletPromo, error)
	findSalesmanByIDFn               func(salesmanID int64, custID string) (model.SalesmanPromo, error)
	findWarehouseByIDFn              func(warehouseID int64, custID string) (model.WarehousePromo, error)
	findActivePromotionsByOutletFn   func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error)
	findProductCriteriasByPromoIDsFn func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error)
	findSlabsByPromoIDsFn            func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error)
	findStratasByPromoIDsFn          func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error)
	getAllRewardProductFromStockV2Fn func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error)
}

func (m *mockPromotionV2RepositoryStore) FindOutletByID(outletID int64, custID string) (model.OutletPromo, error) {
	return m.findOutletByIDFn(outletID, custID)
}

func (m *mockPromotionV2RepositoryStore) FindSalesmanByID(salesmanID int64, custID string) (model.SalesmanPromo, error) {
	return m.findSalesmanByIDFn(salesmanID, custID)
}

func (m *mockPromotionV2RepositoryStore) FindWarehouseByID(warehouseID int64, custID string) (model.WarehousePromo, error) {
	return m.findWarehouseByIDFn(warehouseID, custID)
}

func (m *mockPromotionV2RepositoryStore) FindActivePromotionsByOutletCriteria(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
	return m.findActivePromotionsByOutletFn(req, outlet, salesman)
}

func (m *mockPromotionV2RepositoryStore) FindProductCriteriasByPromoIDs(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) {
	return m.findProductCriteriasByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryStore) FindSlabsByPromoIDs(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
	return m.findSlabsByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryStore) FindStratasByPromoIDs(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) {
	return m.findStratasByPromoIDsFn(promoIDs, custID)
}

func (m *mockPromotionV2RepositoryStore) GetAllRewardProductFromStockV2(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
	return m.getAllRewardProductFromStockV2Fn(req, rewardCtx)
}

type mockStockRepository struct {
	repository.StockRepository
	salesStockUpdatesFn       func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error
	getCurrentStockFn         func(c context.Context, custID string, whID int64, proID int64) (float64, error)
	getCancelStockBasisFn     func(c context.Context, custID string, orderNo string) ([]entity.CancelStockBasis, error)
	cancelSalesStockUpdatesFn func(c context.Context, orderNo string, stockDate time.Time, rows []entity.CancelStockWrite) error
}

func (m *mockStockRepository) SalesStockUpdates(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
	if m.salesStockUpdatesFn == nil {
		return nil
	}
	return m.salesStockUpdatesFn(c, stockUpdates)
}

func (m *mockStockRepository) GetCurrentStock(c context.Context, custID string, whID int64, proID int64) (float64, error) {
	if m.getCurrentStockFn == nil {
		return 0, nil
	}
	return m.getCurrentStockFn(c, custID, whID, proID)
}

func (m *mockStockRepository) GetCancelStockBasis(c context.Context, custID string, orderNo string) ([]entity.CancelStockBasis, error) {
	if m.getCancelStockBasisFn == nil {
		return []entity.CancelStockBasis{}, nil
	}
	return m.getCancelStockBasisFn(c, custID, orderNo)
}

func (m *mockStockRepository) CancelSalesStockUpdates(c context.Context, orderNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
	if m.cancelSalesStockUpdatesFn == nil {
		return nil
	}
	return m.cancelSalesStockUpdatesFn(c, orderNo, stockDate, rows)
}

type mockDbtransaction struct {
	repository.Dbtransaction
}

func (m *mockDbtransaction) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return tFunc(ctx)
}

func TestValidateCancelTransition(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus int64
		wantAllowed   bool
	}{
		{name: "need review to cancelled", currentStatus: 1, wantAllowed: true},
		{name: "processed to cancelled", currentStatus: 2, wantAllowed: true},
		{name: "completed to cancelled", currentStatus: 7, wantAllowed: false},
		{name: "already cancelled", currentStatus: 9, wantAllowed: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := validateCancelTransition(tt.currentStatus)
			if allowed != tt.wantAllowed {
				t.Fatalf("unexpected transition validation: got=%v want=%v", allowed, tt.wantAllowed)
			}
		})
	}
}

func TestBuildCancelStockWriteCommands_MultiSKU(t *testing.T) {
	stockDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)
	basis := []cancelOrderStockBasis{
		{CustID: "C220010001", RoNo: "SO2602090002", WhID: 63, ProID: 474, RefDetID: 1001, QtyOutSmallest: 2400, UnitPrice: 6432, StockDate: stockDate},
		{CustID: "C220010001", RoNo: "SO2602090002", WhID: 63, ProID: 475, RefDetID: 1002, QtyOutSmallest: 1200, UnitPrice: 5000, StockDate: stockDate},
	}

	commands := buildCancelStockWriteCommands(basis)
	if len(commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(commands))
	}
	if commands[0].RoNo != "SO2602090002" || commands[0].RefDetID != 1001 || commands[0].QtySmallest != 2400 {
		t.Fatalf("unexpected command[0]: %+v", commands[0])
	}
	if commands[1].RoNo != "SO2602090002" || commands[1].RefDetID != 1002 || commands[1].QtySmallest != 1200 {
		t.Fatalf("unexpected command[1]: %+v", commands[1])
	}
}

func TestBulkUpdateStatus_Cancel_ConsistentBasisShouldApplyReversal(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603220001"
	processed := int64(entity.PROCESSED)
	cancelled := int64(entity.CANCELLED)
	roDate := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)

	updateCalled := 0
	cancelCalled := 0

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, DataStatus: &processed, RoDate: &roDate}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			updateCalled++
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		getCancelStockBasisFn: func(c context.Context, inputCustID string, inputRoNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{{
				CustID:         custID,
				WhID:           63,
				ProID:          474,
				RefDetID:       1001,
				StockDate:      roDate,
				UnitPrice:      6432,
				QtyOutSmallest: 12,
			}}, nil
		},
		cancelSalesStockUpdatesFn: func(c context.Context, inputRoNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
			cancelCalled++
			if len(rows) != 1 {
				t.Fatalf("expected 1 cancel row, got %d", len(rows))
			}
			if rows[0].RefDetID != 1001 || rows[0].QtySmallest != 12 {
				t.Fatalf("unexpected cancel row %+v", rows[0])
			}
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.BulkUpdateStatus(custID, entity.BulkUpdateStatusOrder{Orders: []entity.UpdateDataStatusBody{{RoNo: roNo, DataStatus: &cancelled}}})
	if err != nil {
		t.Fatalf("BulkUpdateStatus returned error: %v", err)
	}
	if cancelCalled != 1 {
		t.Fatalf("expected cancel stock updates to be called once, got %d", cancelCalled)
	}
	if updateCalled != 1 {
		t.Fatalf("expected order status update to be called once, got %d", updateCalled)
	}
}

func TestBulkUpdateStatus_Cancel_MissingBasisShouldFailWithoutReversal(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603220002"
	processed := int64(entity.PROCESSED)
	cancelled := int64(entity.CANCELLED)
	roDate := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)

	cancelCalled := 0

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, DataStatus: &processed, RoDate: &roDate}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		getCancelStockBasisFn: func(c context.Context, inputCustID string, inputRoNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{}, nil
		},
		cancelSalesStockUpdatesFn: func(c context.Context, inputRoNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
			cancelCalled++
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.BulkUpdateStatus(custID, entity.BulkUpdateStatusOrder{Orders: []entity.UpdateDataStatusBody{{RoNo: roNo, DataStatus: &cancelled}}})
	if err == nil {
		t.Fatalf("expected cancel to fail when basis is missing")
	}
	if cancelCalled != 0 {
		t.Fatalf("expected cancel stock reversal not to be called, got %d", cancelCalled)
	}
}

func TestBulkUpdateStatus_Cancel_AmbiguousBasisShouldFailWithoutReversal(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603220003"
	processed := int64(entity.PROCESSED)
	cancelled := int64(entity.CANCELLED)
	roDate := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)

	cancelCalled := 0
	updateCalled := 0

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, DataStatus: &processed, RoDate: &roDate}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			updateCalled++
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		getCancelStockBasisFn: func(c context.Context, inputCustID string, inputRoNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{{
				CustID:         custID,
				WhID:           63,
				ProID:          474,
				RefDetID:       1001,
				StockDate:      roDate,
				UnitPrice:      6432,
				QtyOutstanding: 12,
				QtyOutSmallest: 12,
				IsAmbiguous:    true,
			}}, nil
		},
		cancelSalesStockUpdatesFn: func(c context.Context, inputRoNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
			cancelCalled++
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.BulkUpdateStatus(custID, entity.BulkUpdateStatusOrder{Orders: []entity.UpdateDataStatusBody{{RoNo: roNo, DataStatus: &cancelled}}})
	if err == nil {
		t.Fatalf("expected cancel to fail when basis is ambiguous")
	}
	if cancelCalled != 0 {
		t.Fatalf("expected cancel stock reversal not to be called, got %d", cancelCalled)
	}
	if updateCalled != 0 {
		t.Fatalf("expected order update not to be called, got %d", updateCalled)
	}
}

func TestBulkUpdateStatus_Cancel_InvalidOutstandingShouldFailWithoutReversal(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603220004"
	processed := int64(entity.PROCESSED)
	cancelled := int64(entity.CANCELLED)
	roDate := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)

	cancelCalled := 0

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, DataStatus: &processed, RoDate: &roDate}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		getCancelStockBasisFn: func(c context.Context, inputCustID string, inputRoNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{{
				CustID:         custID,
				WhID:           63,
				ProID:          474,
				RefDetID:       1001,
				StockDate:      roDate,
				UnitPrice:      6432,
				QtyOutstanding: -1,
				QtyOutSmallest: 0,
			}}, nil
		},
		cancelSalesStockUpdatesFn: func(c context.Context, inputRoNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
			cancelCalled++
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.BulkUpdateStatus(custID, entity.BulkUpdateStatusOrder{Orders: []entity.UpdateDataStatusBody{{RoNo: roNo, DataStatus: &cancelled}}})
	if err == nil {
		t.Fatalf("expected cancel to fail when stock basis outstanding is invalid")
	}
	if cancelCalled != 0 {
		t.Fatalf("expected cancel stock reversal not to be called, got %d", cancelCalled)
	}
}

func TestOrchestratePromoConsultByTabs_IdenticalSignatureShouldSingleConsult(t *testing.T) {
	payloads := map[string]entity.ConsultPromoV2Req{
		"normal":   {OrderDate: "2026-02-24", Details: []entity.ConPromoV2Det{{ProID: 10, Qty1: 2, GrossValue: 2000}}},
		"final":    {OrderDate: "2026-02-24", Details: []entity.ConPromoV2Det{{ProID: 10, Qty1: 2, GrossValue: 2000}}},
		"purchase": {OrderDate: "2026-02-24", Details: []entity.ConPromoV2Det{{ProID: 10, Qty1: 2, GrossValue: 2000}}},
	}
	signatures := map[string]string{
		"normal":   "abc",
		"final":    "abc",
		"purchase": "abc",
	}

	callCount := 0
	consultFn := func(req entity.ConsultPromoV2Req) ([]entity.ConsultPromoResp, error) {
		callCount++
		return []entity.ConsultPromoResp{{PromoID: "PROMO-1", RewardPercentage: []entity.PromoRewardPercentage{{ProID: 10, Promo1: 100}}}}, nil
	}

	result := orchestratePromoConsultByTabs(payloads, signatures, consultFn)

	if callCount != 1 {
		t.Fatalf("single consult expected, got %d", callCount)
	}
	if len(result["normal"]) == 0 || len(result["final"]) == 0 || len(result["purchase"]) == 0 {
		t.Fatalf("all tab results must be populated in single consult path")
	}
}

func TestOrchestratePromoConsultByTabs_DifferentSignatureShouldPerTabConsult(t *testing.T) {
	payloads := map[string]entity.ConsultPromoV2Req{
		"normal":   {OrderDate: "2026-02-24", Details: []entity.ConPromoV2Det{{ProID: 11, Qty1: 1, GrossValue: 1000}}},
		"final":    {OrderDate: "2026-02-24", Details: []entity.ConPromoV2Det{{ProID: 12, Qty1: 1, GrossValue: 2000}}},
		"purchase": {OrderDate: "2026-02-24", Details: []entity.ConPromoV2Det{{ProID: 13, Qty1: 1, GrossValue: 3000}}},
	}
	signatures := map[string]string{
		"normal":   "sig-n",
		"final":    "sig-f",
		"purchase": "sig-p",
	}

	callCount := 0
	consultFn := func(req entity.ConsultPromoV2Req) ([]entity.ConsultPromoResp, error) {
		callCount++
		return []entity.ConsultPromoResp{{PromoID: "PROMO-2"}}, nil
	}

	result := orchestratePromoConsultByTabs(payloads, signatures, consultFn)

	if callCount != 3 {
		t.Fatalf("per tab consult expected 3, got %d", callCount)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 tab results, got %d", len(result))
	}
}

func TestAggregatePromoByProduct_RewardPercentage(t *testing.T) {
	consultResp := []entity.ConsultPromoResp{{
		PromoID:          "PROMO-PERCENT",
		ProductsEligible: []int{1001},
		RewardPercentage: []entity.PromoRewardPercentage{{ProID: 1001, Promo1: 10, Promo2: 20, Promo3: 30, Promo4: 40, Promo5: 50}},
	}}

	aggr := aggregatePromoByProduct(consultResp)
	row := aggr[1001]

	if row.Promo1 != 10 || row.Promo2 != 20 || row.Promo3 != 30 || row.Promo4 != 40 || row.Promo5 != 50 {
		t.Fatalf("unexpected percentage mapping %+v", row)
	}
	if row.PromoTotal != 150 {
		t.Fatalf("unexpected promo_total, got %.2f", row.PromoTotal)
	}
	if len(row.Remarks) != 1 || row.Remarks[0] != "PROMO-PERCENT" {
		t.Fatalf("unexpected remarks %+v", row.Remarks)
	}
}

func TestAggregatePromoByProduct_RewardValue(t *testing.T) {
	consultResp := []entity.ConsultPromoResp{{
		PromoID:          "PROMO-VALUE",
		RewardValue:      []entity.PromoRewardValue{{ProID: 1002, Promo1: 11, Promo2: 12, Promo3: 13, Promo4: 14, Promo5: 15}},
		ProductsEligible: []int{1002},
	}}

	aggr := aggregatePromoByProduct(consultResp)
	row := aggr[1002]

	if row.Promo1 != 11 || row.Promo2 != 12 || row.Promo3 != 13 || row.Promo4 != 14 || row.Promo5 != 15 {
		t.Fatalf("unexpected value mapping %+v", row)
	}
}

func TestAggregatePromoByProduct_RewardProduct(t *testing.T) {
	consultResp := []entity.ConsultPromoResp{{
		PromoID:          "PROMO-PRODUCT",
		ProductsEligible: []int{1003},
		RewardProduct:    []entity.PromoRewardProductDet{{ProID: 1004, Qty1: 1, Promo1: 2, Promo2: 3, Promo3: 4, Promo4: 5, Promo5: 6}},
	}}

	aggr := aggregatePromoByProduct(consultResp)
	rewardRow := aggr[1004]
	eligibleRow := aggr[1003]

	if rewardRow.Promo1 != 2 || rewardRow.Promo2 != 3 || rewardRow.Promo3 != 4 || rewardRow.Promo4 != 5 || rewardRow.Promo5 != 6 {
		t.Fatalf("unexpected reward_product mapping %+v", rewardRow)
	}
	if len(rewardRow.Remarks) != 1 || rewardRow.Remarks[0] != "PROMO-PRODUCT" {
		t.Fatalf("expected reward row remarks to contain promo id, got %+v", rewardRow.Remarks)
	}
	if len(eligibleRow.Remarks) != 0 {
		t.Fatalf("eligible row must not receive reward-product remarks, got %+v", eligibleRow.Remarks)
	}
}

func TestOrchestratePromoConsultByTabs_ConsultErrorShouldFallback(t *testing.T) {
	payloads := map[string]entity.ConsultPromoV2Req{
		"normal":   {OrderDate: "2026-02-24"},
		"final":    {OrderDate: "2026-02-24"},
		"purchase": {OrderDate: "2026-02-24"},
	}
	signatures := map[string]string{
		"normal":   "sig-n",
		"final":    "sig-f",
		"purchase": "sig-p",
	}

	consultFn := func(req entity.ConsultPromoV2Req) ([]entity.ConsultPromoResp, error) {
		if len(req.Details) == 0 {
			return nil, errors.New("consult failed")
		}
		return []entity.ConsultPromoResp{{PromoID: "PROMO-OK"}}, nil
	}

	result := orchestratePromoConsultByTabs(payloads, signatures, consultFn)

	if result == nil {
		t.Fatalf("result should not be nil")
	}
	if len(result["normal"]) != 0 || len(result["final"]) != 0 || len(result["purchase"]) != 0 {
		t.Fatalf("fallback should return empty result for failed consult")
	}
}

func TestAggregatePromoByProduct_RewardProductMarksProductPromotion(t *testing.T) {
	consultResp := []entity.ConsultPromoResp{{
		PromoID:          "PROMO-PRODUCT",
		ProductsEligible: []int{1003},
		RewardProduct:    []entity.PromoRewardProductDet{{ProID: 1003, Qty1: 1, Promo1: 2, Promo2: 3, Promo3: 4, Promo4: 5, Promo5: 6}},
	}}

	aggr := aggregatePromoByProduct(consultResp)
	row := aggr[1003]

	if !row.IsProductPromotion {
		t.Fatalf("reward product promo must mark row as product promotion")
	}
}

func TestApplyPersistedPromoSnapshotToItems_UsesSalesOrderSnapshot(t *testing.T) {
	items := []entity.OrderDetResponse{{
		ProId:                1001,
		PromoSo1:             11,
		PromoSo2:             12,
		PromoSo3:             13,
		PromoSo4:             14,
		PromoSo5:             15,
		PromoRemarksSo:       []string{"PROMO-SO-1", "PROMO-SO-2"},
		IsProductPromotionSo: true,
	}}

	mapped := applyPersistedPromoSnapshotToItems(items, promoSnapshotTabSalesOrder)
	item := mapped[0]

	if item.Promo1 != 11 || item.Promo2 != 12 || item.Promo3 != 13 || item.Promo4 != 14 || item.Promo5 != 15 {
		t.Fatalf("unexpected sales snapshot mapping %+v", item)
	}
	if item.PromoTotal != 65 {
		t.Fatalf("unexpected promo total %.2f", item.PromoTotal)
	}
	if len(item.Remarks) != 2 || item.Remarks[0] != "PROMO-SO-1" || item.Remarks[1] != "PROMO-SO-2" {
		t.Fatalf("unexpected remarks %+v", item.Remarks)
	}
	if !item.IsProductPromotion {
		t.Fatalf("sales snapshot should mark product promotion flag")
	}
}

func TestInjectPromoToOrderItems_PropagatesIsProductPromotion(t *testing.T) {
	firstDetailID := int64(7231)
	secondDetailID := int64(7321)
	items := []entity.OrderDetResponse{{OrderDetId: firstDetailID, ProId: 723}, {OrderDetId: secondDetailID, ProId: 732, IsProductPromotion: true, Remarks: []string{"STALE"}}}
	promoMap := map[int]promoAggregateRow{
		7231: {
			Promo1:             4000000,
			Promo2:             8000000,
			Promo3:             12000000,
			PromoTotal:         24000000,
			Remarks:            []string{"VALUE-GET-QTY"},
			IsProductPromotion: true,
		},
	}

	mapped := injectPromoToOrderItems(items, promoMap)

	if !mapped[0].IsProductPromotion {
		t.Fatalf("expected reward item to propagate is_product_promotion")
	}
	if len(mapped[0].Remarks) != 1 || mapped[0].Remarks[0] != "VALUE-GET-QTY" {
		t.Fatalf("expected reward item remarks from promo map, got %+v", mapped[0].Remarks)
	}
	if mapped[1].IsProductPromotion {
		t.Fatalf("expected item without promo map to reset is_product_promotion")
	}
	if len(mapped[1].Remarks) != 0 {
		t.Fatalf("expected item without promo map to clear stale remarks, got %+v", mapped[1].Remarks)
	}
}

func TestHasPersistedPromoSnapshot_TrueWhenHeaderOrDetailHasSnapshot(t *testing.T) {
	remarks := []string{"PROMO-SO-1"}
	items := []entity.OrderDetResponse{{ProId: 1001, PromoSo3: 7}}

	if !hasPersistedPromoSnapshot(remarks, items, promoSnapshotTabSalesOrder) {
		t.Fatalf("persisted snapshot should be detected when header/detail snapshot exists")
	}
}

func TestHasPersistedPromoSnapshot_FalseForLegacyRows(t *testing.T) {
	items := []entity.OrderDetResponse{{ProId: 1001}}

	if hasPersistedPromoSnapshot(nil, items, promoSnapshotTabSalesOrder) {
		t.Fatalf("legacy row without snapshot should not be detected as persisted snapshot")
	}
}

func TestDetailV2_UsesPersistedSalesOrderSnapshotWithoutConsultFallback(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603110001"
	oprType := "O"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	qty := 52.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	promoValueFinal := 65.0
	discValueFinal := 5.0
	vat := 11.0
	amountFinal := 1326.0
	promoRemarks := model.JSONStringArray{"PROMO-SO-1", "PROMO-SO-2"}
	isProductPromotion := true

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, OprType: &oprType, CustID: custID, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID, WhId: &whID, PromoRemarksSo: promoRemarks, PromoRemarksFinal: promoRemarks}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{
				OrderDetailID:           intPtr(1),
				RoNo:                    roNo,
				ProId:                   748,
				ProCode:                 "PRO-748",
				ProName:                 "Product 748",
				ItemType:                1,
				Qty:                     &qty,
				Qty1:                    &qty1,
				Qty2:                    &qty2,
				Qty3:                    &qty3,
				QtyFinal:                &qty,
				Qty1Final:               &qty1,
				Qty2Final:               &qty2,
				Qty3Final:               &qty3,
				SellPrice1:              &sellPrice1,
				SellPrice2:              &sellPrice2,
				SellPrice3:              &sellPrice3,
				SellPriceFinal1:         &sellPrice1,
				SellPriceFinal2:         &sellPrice2,
				SellPriceFinal3:         &sellPrice3,
				PromoValueFinal:         &promoValueFinal,
				DiscValueFinal:          &discValueFinal,
				Vat:                     &vat,
				AmountFinal:             &amountFinal,
				MpConvUnit2:             &convUnit2,
				MpConvUnit3:             &convUnit3,
				PromoSo1:                float64Ptr(11),
				PromoSo2:                float64Ptr(12),
				PromoSo3:                float64Ptr(13),
				PromoSo4:                float64Ptr(14),
				PromoSo5:                float64Ptr(15),
				PromoRemarksSo:          promoRemarks,
				IsProductPromotionSo:    &isProductPromotion,
				PromoFinal1:             float64Ptr(11),
				PromoFinal2:             float64Ptr(12),
				PromoFinal3:             float64Ptr(13),
				PromoFinal4:             float64Ptr(14),
				PromoFinal5:             float64Ptr(15),
				PromoRemarksFinal:       promoRemarks,
				IsProductPromotionFinal: &isProductPromotion,
			}}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ProCode: "PRO-748", ProName: "Product 748", SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if len(response.Details.Normal) != 1 {
		t.Fatalf("expected 1 normal detail, got %d", len(response.Details.Normal))
	}
	item := response.Details.Normal[0]
	if item.Promo1 != 11 || item.Promo2 != 12 || item.Promo3 != 13 || item.Promo4 != 14 || item.Promo5 != 15 {
		t.Fatalf("expected persisted sales promo snapshot, got %+v", item)
	}
	if len(response.Details.FinalRemarks) != 2 || response.Details.FinalRemarks[0] != "PROMO-SO-1" {
		t.Fatalf("expected persisted sales remarks, got %+v", response.Details.FinalRemarks)
	}
	if !item.IsProductPromotion {
		t.Fatalf("expected persisted product promotion flag")
	}
	if response.OprType == nil || *response.OprType != oprType {
		t.Fatalf("expected opr_type=%s, got %+v", oprType, response.OprType)
	}
}

func TestDetailV2_UsesPersistedRewardProductsWhenSnapshotExists(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603110004"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	qty := 52.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	rewardQty1 := 1.0
	rewardSellPrice1 := 50.0
	promoRemarks := model.JSONStringArray{"PROMO-RWD"}
	isProductPromotion := true
	unit1 := "PCS"
	unit2 := "BOX"
	unit3 := "CRT"

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID, WhId: &whID, PromoRemarksSo: promoRemarks, PromoRemarksFinal: promoRemarks}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{
				{OrderDetailID: intPtr(1), RoNo: roNo, ProId: 748, ProCode: "PRO-748", ProName: "Product 748", ItemType: 1, Qty: &qty, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, QtyFinal: &qty, Qty1Final: &qty1, Qty2Final: &qty2, Qty3Final: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3},
				{OrderDetailID: intPtr(2), RoNo: roNo, ProId: 990, ProCode: "PRM-990", ProName: "Promo Product", ItemType: 2, Qty: &rewardQty1, Qty1: &rewardQty1, QtyFinal: &rewardQty1, Qty1Final: &rewardQty1, SellPrice1: &rewardSellPrice1, SellPriceFinal1: &rewardSellPrice1, PromoSo1: float64Ptr(50), PromoFinal1: float64Ptr(50), PromoRemarksSo: promoRemarks, PromoRemarksFinal: promoRemarks, IsProductPromotionSo: &isProductPromotion, IsProductPromotionFinal: &isProductPromotion, UnitId1: &unit1, UnitId2: &unit2, UnitId3: &unit3, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3},
			}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0, 990: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ProCode: "PRM-990", ProName: "Promo Product", SellPrice1: rewardSellPrice1}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if len(response.Details.Promo) != 0 {
		t.Fatalf("expected reward item to be moved out from sales promo list, got %d", len(response.Details.Promo))
	}
	if len(response.DetailsFinal.Promo) != 0 {
		t.Fatalf("expected reward item to be moved out from final promo list, got %d", len(response.DetailsFinal.Promo))
	}
	if len(response.PurchaseDetails.Promo) != 0 {
		t.Fatalf("expected reward item to be moved out from purchase promo list, got %d", len(response.PurchaseDetails.Promo))
	}
	if len(response.Details.Normal) != 2 {
		t.Fatalf("expected reward detail to be moved into sales normal list, got %d", len(response.Details.Normal))
	}
	if len(response.DetailsFinal.Normal) != 2 {
		t.Fatalf("expected reward detail to be moved into final normal list, got %d", len(response.DetailsFinal.Normal))
	}
	if len(response.PurchaseDetails.Normal) != 2 {
		t.Fatalf("expected reward detail to be moved into purchase normal list, got %d", len(response.PurchaseDetails.Normal))
	}
	if len(response.Details.RewardProducts) != 1 {
		t.Fatalf("expected persisted reward_products for sales tab, got %d", len(response.Details.RewardProducts))
	}
	if len(response.DetailsFinal.RewardProducts) != 1 {
		t.Fatalf("expected persisted reward_products for final tab, got %d", len(response.DetailsFinal.RewardProducts))
	}
	if response.Details.RewardProducts[0].ProID != 990 || response.Details.RewardProducts[0].Promo1 != 50 {
		t.Fatalf("unexpected persisted sales reward product %+v", response.Details.RewardProducts[0])
	}
	if response.Details.RewardProducts[0].UnitId1 == nil || *response.Details.RewardProducts[0].UnitId1 != "PCS" {
		t.Fatalf("expected persisted sales reward product unit_id1=PCS, got %+v", response.Details.RewardProducts[0].UnitId1)
	}
	if response.Details.RewardProducts[0].UnitId2 == nil || *response.Details.RewardProducts[0].UnitId2 != "BOX" {
		t.Fatalf("expected persisted sales reward product unit_id2=BOX, got %+v", response.Details.RewardProducts[0].UnitId2)
	}
	if response.Details.RewardProducts[0].UnitId3 == nil || *response.Details.RewardProducts[0].UnitId3 != "CRT" {
		t.Fatalf("expected persisted sales reward product unit_id3=CRT, got %+v", response.Details.RewardProducts[0].UnitId3)
	}
	if response.DetailsFinal.RewardProducts[0].ProID != 990 || response.DetailsFinal.RewardProducts[0].Promo1 != 50 {
		t.Fatalf("unexpected persisted final reward product %+v", response.DetailsFinal.RewardProducts[0])
	}
	if response.DetailsFinal.RewardProducts[0].UnitId1 == nil || *response.DetailsFinal.RewardProducts[0].UnitId1 != "PCS" {
		t.Fatalf("expected persisted final reward product unit_id1=PCS, got %+v", response.DetailsFinal.RewardProducts[0].UnitId1)
	}

	var salesRewardRow *entity.OrderDetResponse
	for i := range response.Details.Normal {
		if response.Details.Normal[i].ProId == 990 {
			salesRewardRow = &response.Details.Normal[i]
			break
		}
	}
	if salesRewardRow == nil {
		t.Fatalf("expected reward detail in sales normal list")
	}
	if !salesRewardRow.IsProductPromotion {
		t.Fatalf("expected flattened sales reward detail to carry is_product_promotion")
	}
	if salesRewardRow.Promo1 != 50 {
		t.Fatalf("expected flattened sales reward promo1=50, got %+v", salesRewardRow)
	}
	if len(salesRewardRow.Remarks) != 1 || salesRewardRow.Remarks[0] != "PROMO-RWD" {
		t.Fatalf("expected flattened sales reward remarks, got %+v", salesRewardRow.Remarks)
	}

	var finalRewardRow *entity.OrderDetResponse
	for i := range response.DetailsFinal.Normal {
		if response.DetailsFinal.Normal[i].ProId == 990 {
			finalRewardRow = &response.DetailsFinal.Normal[i]
			break
		}
	}
	if finalRewardRow == nil {
		t.Fatalf("expected reward detail in final normal list")
	}
	if !finalRewardRow.IsProductPromotion {
		t.Fatalf("expected flattened final reward detail to carry is_product_promotion")
	}
	if finalRewardRow.Promo1 != 50 {
		t.Fatalf("expected flattened final reward promo1=50, got %+v", finalRewardRow)
	}
}

func TestDetailV2_PostRolloutWithoutSnapshot_UsesConsultV2ByTab(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603120001"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	qty := 52.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0
	consultCalled := 0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, CreatedAt: createdAt, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID, WhId: &whID}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{OrderDetailID: intPtr(2), RoNo: roNo, ProId: 748, ProCode: "PRO-748", ProName: "Product 748", ItemType: 1, Qty: &qty, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, QtyFinal: &qty, Qty1Final: &qty1, Qty2Final: &qty2, Qty3Final: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, Vat: &vat, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3}}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ProCode: "PRO-748", ProName: "Product 748", SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
	}

	promotionRepo := &mockPromotionRepositoryDetailV2{findProductByIDAndCustIDFn: func(productID int64, custID string) (model.ProductRead, error) {
		return model.ProductRead{ProId: int(productID), ConvUnit2: 10, ConvUnit3: 5}, nil
	}}

	promotionV2Repo := &mockPromotionV2RepositoryDetailV2{
		findOutletByIDFn: func(outletID int64, custID string) (model.OutletPromo, error) {
			consultCalled++
			return model.OutletPromo{OutletID: int(outletID)}, nil
		},
		findSalesmanByIDFn: func(salesmanID int64, custID string) (model.SalesmanPromo, error) {
			return model.SalesmanPromo{WhId: int(whID)}, nil
		},
		findWarehouseByIDFn: func(warehouseID int64, custID string) (model.WarehousePromo, error) {
			return model.WarehousePromo{WhID: int(warehouseID)}, nil
		},
		findActivePromotionsByOutletFn: func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
			return []model.PromotionV2{{PromoID: "PROMO-V2-POST", PromoDesc: "Promo Post Rollout", PromoType: model.PromotionTypeSlab}}, nil
		},
		findProductCriteriasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) {
			return nil, nil
		},
		findSlabsByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
			rewardValue := 10.0
			perScope := string(model.PerScopeProduct)
			return []model.PromotionV2Slabs{{PromoID: "PROMO-V2-POST", RuleType: model.RuleTypeQuantity, RewardType: model.RewardTypeFixedValue, RewardValue: &rewardValue, PerScope: &perScope, RangeTo: 200}}, nil
		},
		findStratasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) { return nil, nil },
		getAllRewardProductFromStockV2Fn: func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
			return nil, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, PromotionRepository: promotionRepo, PromotionV2Repository: promotionV2Repo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if consultCalled == 0 {
		t.Fatalf("post-rollout rows without snapshot must consult v2")
	}
	if len(response.Details.PromoRemarksSo) == 0 {
		t.Fatalf("post-rollout rows without snapshot must expose promo remarks from v2 consult")
	}
	if response.Details.Normal[0].Promo1 == 0 {
		t.Fatalf("post-rollout rows without snapshot must inject promo from consult, got %+v", response.Details.Normal[0])
	}
}

func TestBuildRewardProducts_MapsUnitFieldsFromProductMeta(t *testing.T) {
	unit1 := "PCS"
	unit2 := "BOX"
	unit3 := "CRT"
	sellPrice1 := 50.0
	sellPrice2 := 100.0
	sellPrice3 := 200.0

	rewardProducts := buildRewardProducts(
		[]entity.ConsultPromoResp{{
			PromoID: "PROMO-RUNTIME",
			RewardProduct: []entity.PromoRewardProductDet{{
				ProID:      990,
				Qty1:       1,
				Qty2:       2,
				Qty3:       3,
				GrossValue: 850,
				Promo1:     50,
			}},
		}},
		map[int]entity.OrderDetResponse{
			990: {
				ProId:      990,
				ProCode:    "PRM-990",
				ProName:    "Promo Product",
				UnitId1:    &unit1,
				UnitId2:    &unit2,
				UnitId3:    &unit3,
				SellPrice1: &sellPrice1,
				SellPrice2: &sellPrice2,
				SellPrice3: &sellPrice3,
			},
		},
		nil,
	)

	if len(rewardProducts) != 1 {
		t.Fatalf("expected 1 reward product, got %d", len(rewardProducts))
	}
	if rewardProducts[0].UnitId1 == nil || *rewardProducts[0].UnitId1 != "PCS" {
		t.Fatalf("expected runtime reward unit_id1=PCS, got %+v", rewardProducts[0].UnitId1)
	}
	if rewardProducts[0].UnitId2 == nil || *rewardProducts[0].UnitId2 != "BOX" {
		t.Fatalf("expected runtime reward unit_id2=BOX, got %+v", rewardProducts[0].UnitId2)
	}
	if rewardProducts[0].UnitId3 == nil || *rewardProducts[0].UnitId3 != "CRT" {
		t.Fatalf("expected runtime reward unit_id3=CRT, got %+v", rewardProducts[0].UnitId3)
	}
}

func TestBuildRewardProducts_FallbackProductLookupMapsUnitFields(t *testing.T) {
	unit1 := "PCS"
	unit2 := "BOX"

	rewardProducts := buildRewardProducts(
		[]entity.ConsultPromoResp{{
			PromoID: "PROMO-RUNTIME",
			RewardProduct: []entity.PromoRewardProductDet{{
				ProID:      991,
				Qty1:       1,
				GrossValue: 50,
				Promo1:     50,
			}},
		}},
		map[int]entity.OrderDetResponse{},
		func(productID int) (model.ProductRead, error) {
			return model.ProductRead{
				ProId:      productID,
				ProCode:    "PRM-991",
				ProName:    "Promo Product Fallback",
				UnitId1:    unit1,
				UnitId2:    unit2,
				SellPrice1: 50,
			}, nil
		},
	)

	if len(rewardProducts) != 1 {
		t.Fatalf("expected 1 fallback reward product, got %d", len(rewardProducts))
	}
	if rewardProducts[0].UnitId1 == nil || *rewardProducts[0].UnitId1 != "PCS" {
		t.Fatalf("expected fallback reward unit_id1=PCS, got %+v", rewardProducts[0].UnitId1)
	}
	if rewardProducts[0].UnitId2 == nil || *rewardProducts[0].UnitId2 != "BOX" {
		t.Fatalf("expected fallback reward unit_id2=BOX, got %+v", rewardProducts[0].UnitId2)
	}
}

func TestBuildRewardProductsFromPersistedDetails_MapsNilUnitFields(t *testing.T) {
	rewardQty1 := 1.0
	rewardSellPrice1 := 50.0

	rewardProducts := buildRewardProductsFromPersistedDetails([]model.OrderDetailRead{{
		ProId:      992,
		ProCode:    "PRM-992",
		ProName:    "Promo Product Nil Unit",
		ItemType:   2,
		Qty1:       &rewardQty1,
		SellPrice1: &rewardSellPrice1,
		PromoSo1:   float64Ptr(50),
		UnitId1:    nil,
		UnitId2:    nil,
		UnitId3:    nil,
		UnitId4:    nil,
		UnitId5:    nil,
	}}, promoSnapshotTabSalesOrder, nil)

	if len(rewardProducts) != 1 {
		t.Fatalf("expected 1 persisted reward product, got %d", len(rewardProducts))
	}
	if rewardProducts[0].UnitId1 != nil || rewardProducts[0].UnitId2 != nil || rewardProducts[0].UnitId3 != nil || rewardProducts[0].UnitId4 != nil || rewardProducts[0].UnitId5 != nil {
		t.Fatalf("expected nil unit fields to stay nil, got %+v", rewardProducts[0])
	}
}

func TestBuildRewardProductsFromPersistedDetails_FallbackToProductMaster(t *testing.T) {
	rewardQty1 := 1.0
	rewardSellPrice1 := 50.0
	emptyUnit := ""

	rewardProducts := buildRewardProductsFromPersistedDetails([]model.OrderDetailRead{{
		ProId:      723,
		ProCode:    "PRM-723",
		ProName:    "Promo Product Empty Unit",
		ItemType:   2,
		Qty1:       &rewardQty1,
		SellPrice1: &rewardSellPrice1,
		PromoSo1:   float64Ptr(50),
		UnitId1:    &emptyUnit,
		UnitId2:    &emptyUnit,
		UnitId3:    &emptyUnit,
	}}, promoSnapshotTabSalesOrder, map[int]model.Product{
		723: {
			ProductId: 723,
			UnitId1:   "PCS",
			UnitId2:   "CRT",
		},
	})

	if len(rewardProducts) != 1 {
		t.Fatalf("expected 1 persisted reward product with fallback, got %d", len(rewardProducts))
	}
	if rewardProducts[0].UnitId1 == nil || *rewardProducts[0].UnitId1 != "PCS" {
		t.Fatalf("expected fallback persisted reward unit_id1=PCS, got %+v", rewardProducts[0].UnitId1)
	}
	if rewardProducts[0].UnitId2 == nil || *rewardProducts[0].UnitId2 != "CRT" {
		t.Fatalf("expected fallback persisted reward unit_id2=CRT, got %+v", rewardProducts[0].UnitId2)
	}
}

func TestDetailV2_MapsPromoRemarksPerTabFieldsFromPersistedSnapshot(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603120002"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	qty := 52.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0
	salesRemarks := model.JSONStringArray{"PROMO-SO-1"}
	finalRemarks := model.JSONStringArray{"PROMO-FINAL-1"}
	purchaseRemarks := model.JSONStringArray{"PROMO-PO-1"}
	isProductPromotion := true

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, CreatedAt: createdAt, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID, WhId: &whID, PromoRemarksSo: salesRemarks, PromoRemarksFinal: finalRemarks, PromoRemarksPo: purchaseRemarks}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{OrderDetailID: intPtr(1), RoNo: roNo, ProId: 748, ProCode: "PRO-748", ProName: "Product 748", ItemType: 1, Qty: &qty, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, QtyFinal: &qty, Qty1Final: &qty1, Qty2Final: &qty2, Qty3Final: &qty3, QtyPo1: &qty1, QtyPo2: &qty2, QtyPo3: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, SellPricePo1: &sellPrice1, SellPricePo2: &sellPrice2, SellPricePo3: &sellPrice3, Vat: &vat, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3, PromoSo1: float64Ptr(11), PromoFinal1: float64Ptr(22), PromoPo1: float64Ptr(33), PromoRemarksSo: salesRemarks, PromoRemarksFinal: finalRemarks, PromoRemarksPo: purchaseRemarks, IsProductPromotionSo: &isProductPromotion, IsProductPromotionFinal: &isProductPromotion, IsProductPromotionPo: &isProductPromotion}}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ProCode: "PRO-748", ProName: "Product 748", SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if len(response.Details.PromoRemarksSo) != 1 || response.Details.PromoRemarksSo[0] != "PROMO-SO-1" {
		t.Fatalf("expected sales promo_remarks_so, got %+v", response.Details.PromoRemarksSo)
	}
	if len(response.DetailsFinal.PromoRemarksFinal) != 1 || response.DetailsFinal.PromoRemarksFinal[0] != "PROMO-FINAL-1" {
		t.Fatalf("expected final promo_remarks_final, got %+v", response.DetailsFinal.PromoRemarksFinal)
	}
	if len(response.PurchaseDetails.PromoRemarksPo) != 1 || response.PurchaseDetails.PromoRemarksPo[0] != "PROMO-PO-1" {
		t.Fatalf("expected purchase promo_remarks_po, got %+v", response.PurchaseDetails.PromoRemarksPo)
	}
	if len(response.Details.FinalRemarks) != 1 || response.Details.FinalRemarks[0] != "PROMO-SO-1" {
		t.Fatalf("expected compatibility final_remarks for sales tab, got %+v", response.Details.FinalRemarks)
	}
	if len(response.DetailsFinal.FinalRemarks) != 1 || response.DetailsFinal.FinalRemarks[0] != "PROMO-FINAL-1" {
		t.Fatalf("expected compatibility final_remarks for final tab, got %+v", response.DetailsFinal.FinalRemarks)
	}
	if len(response.PurchaseDetails.FinalRemarks) != 1 || response.PurchaseDetails.FinalRemarks[0] != "PROMO-PO-1" {
		t.Fatalf("expected compatibility final_remarks for purchase tab, got %+v", response.PurchaseDetails.FinalRemarks)
	}
}

func TestStore_PersistsPromoSnapshotForSalesAndFinalTabs(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-03-11"
	zero := 0.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0

	var storedOrder model.Order
	var storedDetails []model.OrderDetail
	var stockUpdates []*entity.SalesOrderStockUpdate

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findWarehouseStockByWhIdAndProIdsFn: func(custID string, whID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 100, 990: 100}, nil
		},
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error {
			storedOrder = *data
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := len(storedDetails) + 1
			data.OrderDetailID = &id
			storedDetails = append(storedDetails, *data)
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	promotionV2Repo := &mockPromotionV2RepositoryStore{
		findOutletByIDFn: func(outletID int64, custID string) (model.OutletPromo, error) {
			return model.OutletPromo{OutletID: int(outletID)}, nil
		},
		findSalesmanByIDFn: func(salesmanID int64, custID string) (model.SalesmanPromo, error) {
			return model.SalesmanPromo{WhId: int(whID)}, nil
		},
		findWarehouseByIDFn: func(warehouseID int64, custID string) (model.WarehousePromo, error) {
			return model.WarehousePromo{WhID: int(warehouseID)}, nil
		},
		findActivePromotionsByOutletFn: func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
			return []model.PromotionV2{{PromoID: "PROMO-SO", PromoDesc: "Promo SO", PromoType: model.PromotionTypeSlab}}, nil
		},
		findProductCriteriasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) {
			return nil, nil
		},
		findSlabsByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
			rewardValue := 10.0
			perScope := string(model.PerScopeProduct)
			return []model.PromotionV2Slabs{{PromoID: "PROMO-SO", RuleType: model.RuleTypeQuantity, RewardType: model.RewardTypeFixedValue, RewardValue: &rewardValue, PerScope: &perScope, RangeTo: 1}}, nil
		},
		findStratasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) { return nil, nil },
		getAllRewardProductFromStockV2Fn: func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
			return []model.PromotionRewardProduct{{ProID: 990, QtyStock: 100, ConvUnit2: 10, ConvUnit3: 5}}, nil
		},
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, updates []*entity.SalesOrderStockUpdate) error {
			stockUpdates = updates
			return nil
		},
	}

	promotionRepo := &mockPromotionRepositoryStore{findProductByIDAndCustIDFn: func(productID int64, custID string) (model.ProductRead, error) {
		return model.ProductRead{ProId: int(productID), ConvUnit2: 10, ConvUnit3: 5}, nil
	}}

	service := &orderServiceImpl{
		OrderRepository:       orderRepo,
		PromotionRepository:   promotionRepo,
		PromotionV2Repository: promotionV2Repo,
		StockRepository:       stockRepo,
		Transaction:           &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	response, err := service.Store(request, validateResponse)
	if err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if response.RoNo == "" {
		t.Fatalf("ro number must be generated")
	}
	if len(storedDetails) != 1 {
		t.Fatalf("expected 1 stored detail, got %d", len(storedDetails))
	}
	if len(stockUpdates) != 1 {
		t.Fatalf("expected 1 stock update, got %d", len(stockUpdates))
	}
	if len(storedOrder.PromoRemarksSo) == 0 || storedOrder.PromoRemarksSo[0] != "PROMO-SO" {
		t.Fatalf("expected persisted sales promo remarks, got %+v", storedOrder.PromoRemarksSo)
	}
	if len(storedOrder.PromoRemarksFinal) == 0 || storedOrder.PromoRemarksFinal[0] != "PROMO-SO" {
		t.Fatalf("expected persisted final promo remarks, got %+v", storedOrder.PromoRemarksFinal)
	}
	if got := getValueOrDefault(storedDetails[0].PromoSo1, 0); got == 0 {
		t.Fatalf("expected sales promo snapshot to be persisted")
	}
	if got := getValueOrDefault(storedDetails[0].PromoFinal1, 0); got == 0 {
		t.Fatalf("expected final promo snapshot to be persisted")
	}
}

func TestStore_DeterminesProcessedStatusWhenValidationSuccess(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-03-11"
	zero := 0.0
	qty1 := 1.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	initialStatus := entity.NEED_REVIEW

	var storedOrder model.Order

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error {
			storedOrder = *data
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1
			data.OrderDetailID = &id
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{},
		Transaction:     &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		DataStatus:   intPtrForTest(int64(initialStatus)),
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}
	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}

	if storedOrder.DataStatus == nil || *storedOrder.DataStatus != int64(entity.PROCESSED) {
		t.Fatalf("expected processed data_status, got %+v", storedOrder.DataStatus)
	}
}

func TestStore_DeterminesNeedReviewForRestrictedCreditLimit(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-03-11"
	zero := 0.0
	qty1 := 1.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0
	restricted := LIMIT_ACTION_RESTRICTED

	var storedOrder model.Order
	stockWriteCalled := false

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findWarehouseStockByWhIdAndProIdsFn: func(custID string, whID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 100, 990: 100}, nil
		},
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address"), CreditLimitAction: &restricted}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error {
			storedOrder = *data
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1
			data.OrderDetailID = &id
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled = true
			return nil
		}},
		Transaction: &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: false, Validate3Success: true, Validate4Success: true, IsSuccessValidate: false}
	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}

	if storedOrder.DataStatus == nil || *storedOrder.DataStatus != int64(entity.NEED_REVIEW) {
		t.Fatalf("expected need review data_status, got %+v", storedOrder.DataStatus)
	}
	if stockWriteCalled {
		t.Fatal("expected no stock mutation for restricted need review order")
	}
}

func TestStore_PersistsRewardProductSnapshotForSalesAndFinalTabs(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-03-11"
	zero := 0.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0

	var storedDetails []model.OrderDetail

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			if productID == 990 {
				return model.ProductRead{ProId: productID, ProCode: "PRM-990", ProName: "Promo Product", SellPrice1: 50, SellPrice2: 100, SellPrice3: 200, ConvUnit2: 10, ConvUnit3: 5}, nil
			}
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error { return nil },
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := len(storedDetails) + 1
			data.OrderDetailID = &id
			storedDetails = append(storedDetails, *data)
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	promotionV2Repo := &mockPromotionV2RepositoryStore{
		findOutletByIDFn: func(outletID int64, custID string) (model.OutletPromo, error) {
			return model.OutletPromo{OutletID: int(outletID)}, nil
		},
		findSalesmanByIDFn: func(salesmanID int64, custID string) (model.SalesmanPromo, error) {
			return model.SalesmanPromo{WhId: int(whID)}, nil
		},
		findWarehouseByIDFn: func(warehouseID int64, custID string) (model.WarehousePromo, error) {
			return model.WarehousePromo{WhID: int(warehouseID)}, nil
		},
		findActivePromotionsByOutletFn: func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
			return []model.PromotionV2{{PromoID: "PROMO-RWD", PromoDesc: "Promo Reward Product", PromoType: model.PromotionTypeSlab}}, nil
		},
		findProductCriteriasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) { return nil, nil },
		findSlabsByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
			rewardUom := model.UomTypeSmallest
			return []model.PromotionV2Slabs{{PromoID: "PROMO-RWD", RuleType: model.RuleTypeQuantity, RewardType: model.RewardTypeProduct, RewardUom: &rewardUom, RangeTo: 1}}, nil
		},
		findStratasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) { return nil, nil },
		getAllRewardProductFromStockV2Fn: func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
			return []model.PromotionRewardProduct{{ProID: 990, QtyStock: 100, ConvUnit2: 10, ConvUnit3: 5}}, nil
		},
	}

	stockRepo := &mockStockRepository{salesStockUpdatesFn: func(c context.Context, updates []*entity.SalesOrderStockUpdate) error { return nil }}
	promotionRepo := &mockPromotionRepositoryStore{findProductByIDAndCustIDFn: func(productID int64, custID string) (model.ProductRead, error) {
		return model.ProductRead{ProId: int(productID), ConvUnit2: 10, ConvUnit3: 5}, nil
	}}

	service := &orderServiceImpl{
		OrderRepository:       orderRepo,
		PromotionRepository:   promotionRepo,
		PromotionV2Repository: promotionV2Repo,
		StockRepository:       stockRepo,
		Transaction:           &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:       748,
			Qty1:        &qty1,
			Qty2:        &qty2,
			Qty3:        &qty3,
			ConvUnit2:   intPtr(10),
			ConvUnit3:   intPtr(5),
			SellPrice1:  &sellPrice1,
			SellPrice2:  &sellPrice2,
			SellPrice3:  &sellPrice3,
			PromoValue:  &zero,
			DiscValue:   &zero,
			Vat:         &vat,
			VatValue:    &zero,
			Amount:      &zero,
			AmountFinal: &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	_, err := service.Store(request, validateResponse)
	if err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if len(storedDetails) < 2 {
		t.Fatalf("expected reward detail to be persisted, got %d detail rows", len(storedDetails))
	}

	var rewardDetail *model.OrderDetail
	for i := range storedDetails {
		if storedDetails[i].ItemType == 2 {
			rewardDetail = &storedDetails[i]
			break
		}
	}
	if rewardDetail == nil {
		t.Fatalf("expected persisted reward detail row")
	}
	if got := getValueOrDefault(rewardDetail.PromoSo1, 0); got != 50 {
		t.Fatalf("expected reward promo_so1 to be persisted, got %v", got)
	}
	if got := getValueOrDefault(rewardDetail.PromoFinal1, 0); got != 50 {
		t.Fatalf("expected reward promo_final1 to be persisted, got %v", got)
	}
	if rewardDetail.IsProductPromotionSo == nil || !*rewardDetail.IsProductPromotionSo {
		t.Fatalf("expected reward detail is_product_promotion_so to be true")
	}
	if rewardDetail.IsProductPromotionFinal == nil || !*rewardDetail.IsProductPromotionFinal {
		t.Fatalf("expected reward detail is_product_promotion_final to be true")
	}
	if len(rewardDetail.PromoRemarksSo) != 1 || rewardDetail.PromoRemarksSo[0] != "PROMO-RWD" {
		t.Fatalf("expected reward promo remarks so to be persisted, got %+v", rewardDetail.PromoRemarksSo)
	}
	if len(rewardDetail.PromoRemarksFinal) != 1 || rewardDetail.PromoRemarksFinal[0] != "PROMO-RWD" {
		t.Fatalf("expected reward promo remarks final to be persisted, got %+v", rewardDetail.PromoRemarksFinal)
	}
}

func TestUpdateEnhance_FinalOrderBuildsStockDelta(t *testing.T) {
	roNo := "SO2603060001"
	custID := "C220010001"
	whID := int64(301)
	roDate := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	orderDetailID := int64(1001)
	orderDetailIDInt := 1001
	convUnit2 := 10
	convUnit3 := 1
	oldFinalQtySmallest := 10.0
	sellPrice := 100.0

	var capturedPartialUpdates []map[string]interface{}
	var capturedHeaderUpdate model.Order
	var capturedStockUpdates []*entity.SalesOrderStockUpdate

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{CustID: custID, WhId: &whID, RoDate: &roDate}, nil
		},
		findOrderDetailByDetailIDFn: func(detailID int64, inputCustID string) (model.OrderDetailRead, error) {
			return model.OrderDetailRead{OrderDetailID: &orderDetailIDInt, RoNo: roNo, ProId: 748, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3, QtyFinal: &oldFinalQtySmallest, SellPrice1: &sellPrice, SellPriceFinal1: &sellPrice}, nil
		},
		updateDetailPartialFn: func(c context.Context, detailID int64, inputCustID string, updates map[string]interface{}) error {
			copied := map[string]interface{}{}
			for k, v := range updates {
				copied[k] = v
			}
			capturedPartialUpdates = append(capturedPartialUpdates, copied)
			return nil
		},
		findOrderDetailsForProformaFn: func(ctx context.Context, roNos []string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{OrderDetailID: intPtr(1001), ProId: 748, Qty1Final: float64Ptr(1), Qty2Final: float64Ptr(0), Qty3Final: float64Ptr(2), SellPriceFinal1: &sellPrice, SellPriceFinal2: float64Ptr(0), SellPriceFinal3: &sellPrice, DiscValueFinal: float64Ptr(0), VatValueFinal: float64Ptr(0), Vat: float64Ptr(0)}}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			capturedHeaderUpdate = data
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			capturedStockUpdates = stockUpdates
			return nil
		},
		getCurrentStockFn: func(c context.Context, custID string, whID int64, proID int64) (float64, error) { return 0, nil },
		getCancelStockBasisFn: func(c context.Context, inputCustID string, orderNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{{
				CustID:         custID,
				WhID:           whID,
				ProID:          748,
				RefDetID:       orderDetailID,
				StockDate:      roDate,
				UnitPrice:      sellPrice,
				QtyOutstanding: 1,
				QtyOutSmallest: 1,
			}}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}

	req := entity.EditOrderEnhanceBody{CustId: custID, FinalOrder: []entity.EditFinalOrderDetail{{OrderDetailId: orderDetailID, Qty1Final: float64Ptr(1), Qty2Final: float64Ptr(0), Qty3Final: float64Ptr(2)}}}

	err := service.UpdateEnhance(context.Background(), roNo, req)
	if err != nil {
		t.Fatalf("UpdateEnhance returned error: %v", err)
	}

	if len(capturedPartialUpdates) == 0 {
		t.Fatal("expected final order detail update to be executed")
	}
	firstUpdate := capturedPartialUpdates[0]
	if got := firstUpdate["qty_final"]; got != float64(21) {
		t.Fatalf("unexpected qty_final update: got=%v want=%v", got, float64(21))
	}
	if got := firstUpdate["qty1_final"]; got != float64(1) {
		t.Fatalf("unexpected qty1_final update: got=%v want=%v", got, float64(1))
	}
	if got := firstUpdate["qty2_final"]; got != float64(0) {
		t.Fatalf("unexpected qty2_final update: got=%v want=%v", got, float64(0))
	}
	if got := firstUpdate["qty3_final"]; got != float64(2) {
		t.Fatalf("unexpected qty3_final update: got=%v want=%v", got, float64(2))
	}

	if len(capturedStockUpdates) != 1 {
		t.Fatalf("expected 1 stock update, got %d", len(capturedStockUpdates))
	}

	stockUpdate := capturedStockUpdates[0]
	if stockUpdate.QtyOrderBefore == nil {
		t.Fatal("expected qty_order_before to be populated")
	}
	if *stockUpdate.QtyOrderBefore != 10 {
		t.Fatalf("unexpected qty_order_before: got=%v want=%v", *stockUpdate.QtyOrderBefore, float64(10))
	}
	if stockUpdate.QtyOrder != 21 {
		t.Fatalf("unexpected qty_order: got=%v want=%v", stockUpdate.QtyOrder, float64(21))
	}
	if stockUpdate.ProID != 748 {
		t.Fatalf("unexpected pro_id: got=%v want=%v", stockUpdate.ProID, int64(748))
	}
	if stockUpdate.WhID != whID {
		t.Fatalf("unexpected wh_id: got=%v want=%v", stockUpdate.WhID, whID)
	}
	if stockUpdate.RefDetId != orderDetailID {
		t.Fatalf("unexpected ref_det_id: got=%v want=%v", stockUpdate.RefDetId, orderDetailID)
	}

	if capturedHeaderUpdate.SubTotalFinal == nil || *capturedHeaderUpdate.SubTotalFinal != 300 {
		t.Fatalf("unexpected header subtotal final: %+v", capturedHeaderUpdate.SubTotalFinal)
	}
	if capturedHeaderUpdate.TotalFinal == nil || *capturedHeaderUpdate.TotalFinal != 300 {
		t.Fatalf("unexpected header total final: %+v", capturedHeaderUpdate.TotalFinal)
	}
}

func TestNormalizeEnhancePromoFlags_MapsSalesOrderGenericToSpecific(t *testing.T) {
	generic := true
	request := entity.EditOrderEnhanceBody{
		SalesOrder: []entity.EditSalesOrderDetail{{
			OrderDetailId:      1001,
			IsProductPromotion: &generic,
		}},
	}

	if err := normalizeEnhancePromoFlags(&request); err != nil {
		t.Fatalf("normalizeEnhancePromoFlags returned error: %v", err)
	}
	if request.SalesOrder[0].IsProductPromotionSo == nil || !*request.SalesOrder[0].IsProductPromotionSo {
		t.Fatalf("expected generic sales promotion flag to map into is_product_promotion_so")
	}
}

func TestNormalizeEnhancePromoFlags_MapsFinalOrderGenericToSpecific(t *testing.T) {
	generic := true
	request := entity.EditOrderEnhanceBody{
		FinalOrder: []entity.EditFinalOrderDetail{{
			OrderDetailId:      1002,
			IsProductPromotion: &generic,
		}},
	}

	if err := normalizeEnhancePromoFlags(&request); err != nil {
		t.Fatalf("normalizeEnhancePromoFlags returned error: %v", err)
	}
	if request.FinalOrder[0].IsProductPromotionFinal == nil || !*request.FinalOrder[0].IsProductPromotionFinal {
		t.Fatalf("expected generic final promotion flag to map into is_product_promotion_final")
	}
}

func TestNormalizeEnhancePromoFlags_MergesAliasDetailArrays(t *testing.T) {
	flag := true
	request := entity.EditOrderEnhanceBody{
		SalesOrderDetails: []entity.EditSalesOrderDetail{{
			OrderDetailId:      1001,
			IsProductPromotion: &flag,
		}},
		FinalOrderDetails: []entity.EditFinalOrderDetail{{
			OrderDetailId:      1002,
			IsProductPromotion: &flag,
		}},
	}

	if err := normalizeEnhancePromoFlags(&request); err != nil {
		t.Fatalf("normalizeEnhancePromoFlags returned error: %v", err)
	}
	if len(request.SalesOrder) != 1 {
		t.Fatalf("expected sales_order alias to be merged into sales_order, got %d", len(request.SalesOrder))
	}
	if len(request.FinalOrder) != 1 {
		t.Fatalf("expected final_order alias to be merged into final_order, got %d", len(request.FinalOrder))
	}
	if request.SalesOrder[0].IsProductPromotionSo == nil || !*request.SalesOrder[0].IsProductPromotionSo {
		t.Fatalf("expected merged sales order alias to map generic promo flag")
	}
	if request.FinalOrder[0].IsProductPromotionFinal == nil || !*request.FinalOrder[0].IsProductPromotionFinal {
		t.Fatalf("expected merged final order alias to map generic promo flag")
	}
}

func TestNormalizeEnhancePromoFlags_ReturnsErrorOnGenericSpecificConflict(t *testing.T) {
	generic := true
	specific := false
	request := entity.EditOrderEnhanceBody{
		SalesOrder: []entity.EditSalesOrderDetail{{
			OrderDetailId:        1001,
			IsProductPromotion:   &generic,
			IsProductPromotionSo: &specific,
		}},
	}

	err := normalizeEnhancePromoFlags(&request)
	if err == nil {
		t.Fatalf("expected conflict error when generic and sales specific promo flags differ")
	}
	if !strings.Contains(err.Error(), "is_product_promotion") || !strings.Contains(err.Error(), "is_product_promotion_so") {
		t.Fatalf("unexpected conflict error message: %v", err)
	}
}

func TestRecomputePromoStateForTab_PreservesRowSpecificPromoFlagsForSameProduct(t *testing.T) {
	custID := "C220010001"
	parentCustID := "C220010001"
	roNo := "SO2603250001"
	firstDetailID := 1001
	secondDetailID := 1002
	proID := 748
	qty := 1.0
	zero := 0.0
	price := 100.0
	vat := 0.0
	roDate := time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC)

	type partialUpdateCapture struct {
		detailID int64
		updates  map[string]interface{}
	}
	captured := []partialUpdateCapture{}

	orderRepo := &mockOrderRepository{
		updateDetailPartialFn: func(c context.Context, detailID int64, inputCustID string, updates map[string]interface{}) error {
			copied := map[string]interface{}{}
			for k, v := range updates {
				copied[k] = v
			}
			captured = append(captured, partialUpdateCapture{detailID: detailID, updates: copied})
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}
	_, err := service.recomputePromoStateForTab(
		context.Background(),
		model.OrderList{RoNo: roNo, RoDate: &roDate},
		custID,
		parentCustID,
		[]model.OrderDetailRead{
			{OrderDetailID: &firstDetailID, ProId: proID, ItemType: 1, Qty1: &qty, Qty2: &zero, Qty3: &zero, SellPrice1: &price, SellPrice2: &zero, SellPrice3: &zero, Vat: &vat},
			{OrderDetailID: &secondDetailID, ProId: proID, ItemType: 1, Qty1: &qty, Qty2: &zero, Qty3: &zero, SellPrice1: &price, SellPrice2: &zero, SellPrice3: &zero, Vat: &vat},
		},
		promoSnapshotTabSalesOrder,
		map[int64]promoFlagOverride{
			int64(firstDetailID):  {SalesOrder: boolPtr(true)},
			int64(secondDetailID): {SalesOrder: boolPtr(false)},
		},
	)
	if err != nil {
		t.Fatalf("recomputePromoStateForTab returned error: %v", err)
	}

	if len(captured) != 2 {
		t.Fatalf("expected 2 partial updates for duplicated product rows, got %d", len(captured))
	}

	flags := map[int64]bool{}
	for _, row := range captured {
		value, ok := row.updates["is_product_promotion_so"].(bool)
		if !ok {
			t.Fatalf("expected is_product_promotion_so bool in updates, got %+v", row.updates)
		}
		flags[row.detailID] = value
	}

	if !flags[int64(firstDetailID)] {
		t.Fatalf("expected first row to preserve explicit true promo flag")
	}
	if flags[int64(secondDetailID)] {
		t.Fatalf("expected second row to preserve explicit false promo flag")
	}
}

func TestBuildDetailPromoSnapshotUpdates_SalesOrderMapsSnapshotFields(t *testing.T) {
	remarks := []string{"PROMO-SO-1"}
	orderDetailID := 7481
	aggregate := map[int]promoAggregateRow{7481: {Promo1: 10, Promo2: 20, Promo3: 30, Promo4: 40, Promo5: 50, PromoTotal: 150, Remarks: []string{"PROMO-SO-1"}, IsProductPromotion: true}}

	updates := buildDetailPromoSnapshotUpdates(model.OrderDetailRead{OrderDetailID: &orderDetailID, ProId: 748}, aggregate, promoSnapshotTabSalesOrder)

	if updates["promo_so1"] != 10.0 || updates["promo_so5"] != 50.0 {
		t.Fatalf("unexpected sales snapshot promo mapping %+v", updates)
	}
	if updates["is_product_promotion_so"] != true {
		t.Fatalf("sales snapshot must persist product promotion flag")
	}
	if got, ok := updates["promo_remarks_so"].(model.JSONStringArray); !ok || len(got) != 1 || got[0] != remarks[0] {
		t.Fatalf("unexpected sales snapshot remarks %+v", updates["promo_remarks_so"])
	}
}

func TestBuildHeaderPromoSnapshotUpdate_FinalOrderMapsRemarks(t *testing.T) {
	remarks := []string{"PROMO-FINAL-1", "PROMO-FINAL-2"}
	update := buildHeaderPromoSnapshotUpdate(remarks, promoSnapshotTabFinalOrder)

	if got, ok := update["promo_remarks_final"].(model.JSONStringArray); !ok || len(got) != 2 || got[0] != "PROMO-FINAL-1" || got[1] != "PROMO-FINAL-2" {
		t.Fatalf("unexpected final header snapshot %+v", update)
	}
}

func TestBuildEnhancePromoPayload_UsesProjectedWholeOrderDetails(t *testing.T) {
	roDate := time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC)
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	firstDetailID := 1001
	secondDetailID := 1002
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	otherQty1 := 5.0
	otherQty2 := 0.0
	otherQty3 := 0.0
	price1 := 100.0
	price2 := 0.0
	price3 := 1000.0
	otherPrice1 := 200.0

	payload := buildEnhancePromoPayload(
		model.OrderList{RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID, WhId: &whID},
		"C220010001",
		"C220010001",
		[]model.OrderDetailRead{{OrderDetailID: &firstDetailID, ProId: 748, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, SellPrice1: &price1, SellPrice2: &price2, SellPrice3: &price3}, {OrderDetailID: &secondDetailID, ProId: 749, Qty1: &otherQty1, Qty2: &otherQty2, Qty3: &otherQty3, SellPrice1: &otherPrice1}},
		int64(firstDetailID),
		promoSnapshotTabSalesOrder,
		3,
		0,
		1,
		150,
		0,
		900,
	)

	if len(payload.Details) != 2 {
		t.Fatalf("expected projected payload to include 2 details, got %d", len(payload.Details))
	}
	if payload.Details[0].ProID != 748 || payload.Details[1].ProID != 749 {
		t.Fatalf("unexpected payload product order %+v", payload.Details)
	}
	if payload.Details[0].Qty1 != 3 || payload.Details[0].Qty3 != 1 {
		t.Fatalf("expected edited detail quantities to be projected, got %+v", payload.Details[0])
	}
	if payload.Details[0].GrossValue != 1350 {
		t.Fatalf("expected edited detail gross value 1350, got %+v", payload.Details[0].GrossValue)
	}
	if payload.Details[1].GrossValue != 1000 {
		t.Fatalf("expected untouched detail to remain in payload, got %+v", payload.Details[1].GrossValue)
	}
}

func TestMergeUniqueRemarks_DeduplicatesAndSorts(t *testing.T) {
	merged := mergeUniqueRemarks([]string{"PROMO-B", "PROMO-A"}, []string{"PROMO-C", "PROMO-A", ""})

	if len(merged) != 3 {
		t.Fatalf("expected 3 merged remarks, got %+v", merged)
	}
	if merged[0] != "PROMO-A" || merged[1] != "PROMO-B" || merged[2] != "PROMO-C" {
		t.Fatalf("unexpected merged remarks %+v", merged)
	}
}

func TestBuildRewardProductStockDeltas_DetectsRowBasedDeleteAndInsertForRewardProduct(t *testing.T) {
	roDate := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	oldQty := 2.0
	rewardDetID := int64(9001)

	existingRewards := []model.OrderDetailRead{{OrderDetailID: intPtrForTest(rewardDetID), ProId: 888, ItemType: 2, QtyFinal: &oldQty, SellPriceFinal1: float64Ptr(0)}}
	newRewards := []entity.OrderRewardProductResponse{{ProID: 888, Qty1: 5}}

	updates := buildRewardProductStockDeltas("C220010001", "SO2603060001", 301, roDate, existingRewards, newRewards)

	if len(updates) != 2 {
		t.Fatalf("expected 2 row-based reward stock deltas, got %d", len(updates))
	}

	var insertDelta *entity.SalesOrderStockUpdate
	var deleteDelta *entity.SalesOrderStockUpdate
	for _, update := range updates {
		if update.QtyOrderBefore == nil {
			insertDelta = update
			continue
		}
		deleteDelta = update
	}

	if deleteDelta == nil || deleteDelta.RefDetId != rewardDetID || *deleteDelta.QtyOrderBefore != 2 || deleteDelta.QtyOrder != 0 {
		t.Fatalf("unexpected reward delete delta %+v", deleteDelta)
	}
	if insertDelta == nil || insertDelta.RefDetId != 0 || insertDelta.QtyOrder != 5 {
		t.Fatalf("unexpected reward insert delta %+v", insertDelta)
	}
	if insertDelta.ProID != 888 || deleteDelta.ProID != 888 {
		t.Fatalf("unexpected reward pro ids %+v", updates)
	}
}

func TestBuildRewardProductStockDeltasFromModels_DetectsInsertUpdateAndDelete(t *testing.T) {
	roDate := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	oldQty := 2.0
	oldPrice := 10.0
	newPrice := 15.0
	unchangedQty := 3.0
	unchangedPrice := 20.0
	deletedQty := 4.0
	deletedPrice := 30.0
	newQty := 5.0
	insertedPrice := 40.0
	updatedID := 9001
	unchangedID := 9002
	deletedID := 9003
	insertedID := 9004

	existingRewards := []model.OrderDetailRead{
		{OrderDetailID: &updatedID, ProId: 888, ItemType: 2, QtyFinal: &oldQty, SellPriceFinal1: &oldPrice},
		{OrderDetailID: &unchangedID, ProId: 889, ItemType: 2, QtyFinal: &unchangedQty, SellPriceFinal1: &unchangedPrice},
		{OrderDetailID: &deletedID, ProId: 890, ItemType: 2, QtyFinal: &deletedQty, SellPriceFinal1: &deletedPrice},
	}
	newRewardModels := []model.OrderDetail{
		{OrderDetailID: &updatedID, ProId: 888, ItemType: 2, QtyFinal: newQty, SellPriceFinal1: &newPrice},
		{OrderDetailID: &unchangedID, ProId: 889, ItemType: 2, QtyFinal: unchangedQty, SellPriceFinal1: &unchangedPrice},
		{OrderDetailID: &insertedID, ProId: 891, ItemType: 2, QtyFinal: newQty, SellPriceFinal1: &insertedPrice},
	}

	updates := buildRewardProductStockDeltasFromModels("C220010001", "SO2603060001", 301, roDate, existingRewards, newRewardModels)

	if len(updates) != 6 {
		t.Fatalf("expected 6 row-based reward stock deltas, got %d", len(updates))
	}

	if updates[0].ProID != 888 || updates[0].RefDetId != int64(updatedID) || updates[0].QtyOrderBefore == nil || *updates[0].QtyOrderBefore != 2 || updates[0].QtyOrder != 0 {
		t.Fatalf("unexpected updated reward delete delta %+v", updates[0])
	}
	if updates[1].ProID != 888 || updates[1].RefDetId != int64(updatedID) || updates[1].QtyOrderBefore != nil || updates[1].QtyOrder != 5 {
		t.Fatalf("unexpected updated reward insert delta %+v", updates[1])
	}
	if updates[2].ProID != 889 || updates[2].RefDetId != int64(unchangedID) || updates[2].QtyOrderBefore == nil || *updates[2].QtyOrderBefore != 3 || updates[2].QtyOrder != 0 {
		t.Fatalf("unexpected unchanged reward delete delta %+v", updates[2])
	}
	if updates[3].ProID != 889 || updates[3].RefDetId != int64(unchangedID) || updates[3].QtyOrderBefore != nil || updates[3].QtyOrder != 3 {
		t.Fatalf("unexpected unchanged reward insert delta %+v", updates[3])
	}
	if updates[4].ProID != 890 || updates[4].RefDetId != int64(deletedID) || updates[4].QtyOrderBefore == nil || *updates[4].QtyOrderBefore != 4 || updates[4].QtyOrder != 0 {
		t.Fatalf("unexpected deleted reward delta %+v", updates[4])
	}
	if updates[5].ProID != 891 || updates[5].RefDetId != int64(insertedID) || updates[5].QtyOrderBefore != nil || updates[5].QtyOrder != 5 {
		t.Fatalf("unexpected inserted reward delta %+v", updates[5])
	}
}

func TestBuildRewardProductStockDeltasFromModels_PreservesDuplicateRowsByOrderDetailID(t *testing.T) {
	roDate := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	oldQtyA := 2.0
	oldQtyB := 3.0
	newQtyA := 4.0
	newQtyB := 1.0
	oldPriceA := 10.0
	oldPriceB := 12.0
	newPriceA := 11.0
	newPriceB := 13.0
	firstID := 9001
	secondID := 9002

	existingRewards := []model.OrderDetailRead{
		{OrderDetailID: &firstID, ProId: 888, ItemType: 2, QtyFinal: &oldQtyA, SellPriceFinal1: &oldPriceA},
		{OrderDetailID: &secondID, ProId: 888, ItemType: 2, QtyFinal: &oldQtyB, SellPriceFinal1: &oldPriceB},
	}
	newRewardModels := []model.OrderDetail{
		{OrderDetailID: &firstID, ProId: 888, ItemType: 2, QtyFinal: newQtyA, SellPriceFinal1: &newPriceA},
		{OrderDetailID: &secondID, ProId: 888, ItemType: 2, QtyFinal: newQtyB, SellPriceFinal1: &newPriceB},
	}

	updates := buildRewardProductStockDeltasFromModels("C220010001", "SO2603060001", 301, roDate, existingRewards, newRewardModels)

	if len(updates) != 4 {
		t.Fatalf("expected 4 row-based stock deltas for duplicate reward rows, got %+v", updates)
	}
	if updates[0].RefDetId != int64(firstID) || updates[0].QtyOrderBefore == nil || *updates[0].QtyOrderBefore != oldQtyA || updates[0].QtyOrder != 0 {
		t.Fatalf("unexpected delete delta for first row %+v", updates[0])
	}
	if updates[1].RefDetId != int64(firstID) || updates[1].QtyOrderBefore != nil || updates[1].QtyOrder != newQtyA {
		t.Fatalf("unexpected insert delta for first row %+v", updates[1])
	}
	if updates[2].RefDetId != int64(secondID) || updates[2].QtyOrderBefore == nil || *updates[2].QtyOrderBefore != oldQtyB || updates[2].QtyOrder != 0 {
		t.Fatalf("unexpected delete delta for second row %+v", updates[2])
	}
	if updates[3].RefDetId != int64(secondID) || updates[3].QtyOrderBefore != nil || updates[3].QtyOrder != newQtyB {
		t.Fatalf("unexpected insert delta for second row %+v", updates[3])
	}
}

func TestCreateOrderDetailFromSalesOrder_DefaultsFinalPromoFlagIndependently(t *testing.T) {
	salesFlag := true
	var storedDetail *model.OrderDetail

	orderRepo := &mockOrderRepository{
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			stored := *data
			storedDetail = &stored
			id := 9001
			data.OrderDetailID = &id
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}
	_, err := service.createOrderDetailFromSalesOrder(context.Background(), "SO2603170002", "C220010001", 301, time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC), entity.AddSalesOrderDetail{
		ProId:                748,
		Qty1:                 1,
		Qty2:                 0,
		Qty3:                 0,
		SellPriceSystem1:     100,
		SellPriceSystem2:     0,
		SellPriceSystem3:     0,
		SellPrice1:           100,
		SellPrice2:           0,
		SellPrice3:           0,
		UnitId1:              "PCS",
		UnitId2:              "BOX",
		UnitId3:              "CTN",
		Qty1Stock:            1,
		Qty2Stock:            0,
		Qty3Stock:            0,
		IsProductPromotionSo: &salesFlag,
	})
	if err != nil {
		t.Fatalf("createOrderDetailFromSalesOrder returned error: %v", err)
	}
	if storedDetail == nil {
		t.Fatalf("expected sales order detail to be stored")
	}
	if storedDetail.IsProductPromotionSo == nil || !*storedDetail.IsProductPromotionSo {
		t.Fatalf("expected sales promo flag to be persisted")
	}
	if storedDetail.IsProductPromotionFinal == nil || *storedDetail.IsProductPromotionFinal {
		t.Fatalf("expected final promo flag to default independently to false, got %+v", storedDetail.IsProductPromotionFinal)
	}
}

func TestCreateOrderDetailFromPurchaseOrder_InheritsUOMFromProductMaster(t *testing.T) {
	var storedDetail *model.OrderDetail

	orderRepo := &mockOrderRepository{
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{
				ProId:      productID,
				ConvUnit2:  10,
				ConvUnit3:  5,
				UnitId1:    "PCS",
				UnitId2:    "CRT",
				UnitId3:    "BOX",
				SellPrice1: 100,
				SellPrice2: 50,
				SellPrice3: 10,
			}, nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			stored := *data
			storedDetail = &stored
			id := 9101
			data.OrderDetailID = &id
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}
	_, err := service.createOrderDetailFromPurchaseOrder(
		context.Background(),
		"SO2603180007",
		"C220010001",
		301,
		time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
		entity.AddPurchaseOrderDetail{
			ProId:            748,
			QtyPo1:           1,
			QtyPo2:           0,
			QtyPo3:           1,
			SellPriceSystem1: 100,
			SellPriceSystem2: 50,
			SellPriceSystem3: 10,
			SellPricePo1:     100,
			SellPricePo2:     50,
			SellPricePo3:     10,
		},
	)
	if err != nil {
		t.Fatalf("createOrderDetailFromPurchaseOrder returned error: %v", err)
	}
	if storedDetail == nil {
		t.Fatalf("expected purchase order detail to be stored")
	}
	if storedDetail.UnitId1 == nil || *storedDetail.UnitId1 != "PCS" {
		t.Fatalf("expected unit_id1 to inherit from product master, got %+v", storedDetail.UnitId1)
	}
	if storedDetail.UnitId2 == nil || *storedDetail.UnitId2 != "CRT" {
		t.Fatalf("expected unit_id2 to inherit from product master, got %+v", storedDetail.UnitId2)
	}
	if storedDetail.UnitId3 == nil || *storedDetail.UnitId3 != "BOX" {
		t.Fatalf("expected unit_id3 to inherit from product master, got %+v", storedDetail.UnitId3)
	}
}

func TestCreateOrderDetailFromFinalOrder_InheritsUOMFromProductMaster(t *testing.T) {
	var storedDetail *model.OrderDetail

	orderRepo := &mockOrderRepository{
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{
				ProId:      productID,
				ConvUnit2:  10,
				ConvUnit3:  5,
				UnitId1:    "PCS",
				UnitId2:    "CRT",
				UnitId3:    "BOX",
				SellPrice1: 100,
				SellPrice2: 50,
				SellPrice3: 10,
			}, nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			stored := *data
			storedDetail = &stored
			id := 9201
			data.OrderDetailID = &id
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}
	_, err := service.createOrderDetailFromFinalOrder(
		context.Background(),
		"SO2603180008",
		"C220010001",
		301,
		time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
		entity.AddFinalOrderDetail{
			ProId:            749,
			Qty1Final:        2,
			Qty2Final:        0,
			Qty3Final:        1,
			SellPriceSystem1: 100,
			SellPriceSystem2: 50,
			SellPriceSystem3: 10,
			SellPriceFinal1:  100,
			SellPriceFinal2:  50,
			SellPriceFinal3:  10,
		},
	)
	if err != nil {
		t.Fatalf("createOrderDetailFromFinalOrder returned error: %v", err)
	}
	if storedDetail == nil {
		t.Fatalf("expected final order detail to be stored")
	}
	if storedDetail.UnitId1 == nil || *storedDetail.UnitId1 != "PCS" {
		t.Fatalf("expected unit_id1 to inherit from product master, got %+v", storedDetail.UnitId1)
	}
	if storedDetail.UnitId2 == nil || *storedDetail.UnitId2 != "CRT" {
		t.Fatalf("expected unit_id2 to inherit from product master, got %+v", storedDetail.UnitId2)
	}
	if storedDetail.UnitId3 == nil || *storedDetail.UnitId3 != "BOX" {
		t.Fatalf("expected unit_id3 to inherit from product master, got %+v", storedDetail.UnitId3)
	}
}

func TestUpdateEnhance_SalesOrderPersistsPromoSnapshotAndHeaderRemarks(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603110101"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC)
	orderDetailID := int64(1001)
	orderDetailIDInt := 1001
	convUnit2 := 10
	convUnit3 := 5
	qty := 52.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	oldPromoValueFinal := 0.0
	vat := 11.0
	remarks := model.JSONStringArray{"PROMO-SO-1"}

	var detailUpdates []map[string]interface{}
	var headerUpdates []model.Order
	var stockUpdates []*entity.SalesOrderStockUpdate

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{CustID: custID, WhId: &whID, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID}, nil
		},
		findOutletByIDFn: func(outletID int, custID string, parentCustID string) (model.OutletRead, error) {
			return model.OutletRead{OutletId: outletID}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		findOrderDetailByDetailIDFn: func(detailID int64, inputCustID string) (model.OrderDetailRead, error) {
			return model.OrderDetailRead{OrderDetailID: &orderDetailIDInt, RoNo: roNo, ProId: 748, Qty: &qty, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, QtyFinal: &qty, Qty1Final: &qty1, Qty2Final: &qty2, Qty3Final: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, PromoValueFinal: &oldPromoValueFinal, Vat: &vat, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3, PromoRemarksFinal: remarks, PromoRemarksSo: remarks, IsProductPromotionFinal: boolPtr(true), IsProductPromotionSo: boolPtr(true)}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5, SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3}, nil
		},
		updateDetailPartialFn: func(c context.Context, detailID int64, inputCustID string, updates map[string]interface{}) error {
			copied := map[string]interface{}{}
			for k, v := range updates {
				copied[k] = v
			}
			detailUpdates = append(detailUpdates, copied)
			return nil
		},
		findOrderDetailsForProformaFn: func(ctx context.Context, roNos []string, inputCustID string) ([]model.OrderDetailRead, error) {
			addedID := intPtrForTest(1002)
			addedQty := 1.0
			zero := 0.0
			addedPrice := 50.0
			return []model.OrderDetailRead{
				{OrderDetailID: &orderDetailIDInt, ProId: 748, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, Qty1Final: &qty1, Qty2Final: &qty2, Qty3Final: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, DiscValueFinal: float64Ptr(5), VatValueFinal: float64Ptr(131), Vat: &vat},
				{OrderDetailID: addedID, ProId: 749, Qty1: &addedQty, Qty2: &zero, Qty3: &zero, Qty1Final: &addedQty, Qty2Final: &zero, Qty3Final: &zero, SellPrice1: &addedPrice, SellPrice2: &zero, SellPrice3: &zero, SellPriceFinal1: &addedPrice, SellPriceFinal2: &zero, SellPriceFinal3: &zero, DiscValueFinal: float64Ptr(0), VatValueFinal: float64Ptr(0), Vat: &vat},
			}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			headerUpdates = append(headerUpdates, data)
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1002
			data.OrderDetailID = &id
			return nil
		},
	}

	promotionRepo := &mockPromotionRepositoryEnhance{findProductByIDAndCustIDFn: func(productID int64, custID string) (model.ProductRead, error) {
		return model.ProductRead{ProId: int(productID), ConvUnit2: 10, ConvUnit3: 5}, nil
	}}

	promotionV2Repo := &mockPromotionV2RepositoryEnhance{
		findOutletByIDFn: func(outletID int64, custID string) (model.OutletPromo, error) {
			return model.OutletPromo{OutletID: int(outletID)}, nil
		},
		findSalesmanByIDFn: func(salesmanID int64, custID string) (model.SalesmanPromo, error) {
			return model.SalesmanPromo{WhId: int(whID)}, nil
		},
		findWarehouseByIDFn: func(warehouseID int64, custID string) (model.WarehousePromo, error) {
			return model.WarehousePromo{WhID: int(warehouseID)}, nil
		},
		findActivePromotionsByOutletFn: func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
			return []model.PromotionV2{{PromoID: "PROMO-SO-1", PromoDesc: "Promo SO 1", PromoType: model.PromotionTypeSlab}}, nil
		},
		findProductCriteriasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) { return nil, nil },
		findSlabsByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
			rewardValue := 10.0
			perScope := string(model.PerScopeProduct)
			return []model.PromotionV2Slabs{{PromoID: "PROMO-SO-1", RuleType: model.RuleTypeQuantity, RewardType: model.RewardTypeFixedValue, RewardValue: &rewardValue, PerScope: &perScope, RangeTo: 1}}, nil
		},
		findStratasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) { return nil, nil },
		getAllRewardProductFromStockV2Fn: func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
			return []model.PromotionRewardProduct{{ProID: 990, QtyStock: 100, ConvUnit2: 10, ConvUnit3: 5}}, nil
		},
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdatesInput []*entity.SalesOrderStockUpdate) error {
			stockUpdates = stockUpdatesInput
			return nil
		},
		getCurrentStockFn: func(c context.Context, custID string, whID int64, proID int64) (float64, error) { return 100, nil },
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, PromotionRepository: promotionRepo, PromotionV2Repository: promotionV2Repo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}

	req := entity.EditOrderEnhanceBody{CustId: custID, ParentCustId: custID, SalesOrder: []entity.EditSalesOrderDetail{{OrderDetailId: orderDetailID, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3}}, AddSalesOrder: []entity.AddSalesOrderDetail{{ProId: 749, Qty1: 1, Qty2: 0, Qty3: 0, SellPriceSystem1: 50, SellPriceSystem2: 0, SellPriceSystem3: 0, SellPrice1: 50, SellPrice2: 0, SellPrice3: 0, UnitId1: "PCS", UnitId2: "BOX", UnitId3: "CRT"}}}

	if err := service.UpdateEnhance(context.Background(), roNo, req); err != nil {
		t.Fatalf("UpdateEnhance returned error: %v", err)
	}
	if len(detailUpdates) == 0 {
		t.Fatalf("expected detail updates to be executed")
	}
	foundSnapshot := false
	for _, updates := range detailUpdates {
		if _, ok := updates["promo_remarks_so"]; ok {
			foundSnapshot = true
			break
		}
	}
	if !foundSnapshot {
		t.Fatalf("expected promo snapshot updates in sales order flow")
	}
	foundHeaderSnapshot := false
	for _, update := range headerUpdates {
		if len(update.PromoRemarksSo) > 0 {
			foundHeaderSnapshot = true
			if update.PromoRemarksSo[0] != "PROMO-SO-1" {
				t.Fatalf("expected header sales promo remarks, got %+v", update.PromoRemarksSo)
			}
			if len(update.PromoRemarksFinal) == 0 || update.PromoRemarksFinal[0] != "PROMO-SO-1" {
				t.Fatalf("expected header final promo remarks refreshed from sales tab, got %+v", update.PromoRemarksFinal)
			}
		}
	}
	if !foundHeaderSnapshot {
		t.Fatalf("expected header promo snapshot update in sales order flow")
	}
	if len(stockUpdates) == 0 {
		t.Fatalf("expected stock update for sales order tab")
	}
}

func TestUpdateEnhance_SalesOrderDeletedEffectiveQtyZero_ExcludesRowFromFinalPromoRecompute(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603180001"
	whID := int64(301)
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	orderDetailAID := int64(1001)
	orderDetailAIDInt := 1001
	convUnit2 := 10
	convUnit3 := 5
	qtyAOld := 52.0
	qtyAZero := 0.0
	sellPriceA1 := 100.0
	sellPriceA2 := 0.0
	sellPriceA3 := 1000.0
	vat := 11.0

	type detailPartialCapture struct {
		id      int64
		updates map[string]interface{}
	}
	var partialCaptures []detailPartialCapture
	var headerUpdates []model.Order
	var stockUpdates []*entity.SalesOrderStockUpdate

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{CustID: custID, WhId: &whID, RoDate: &roDate}, nil
		},
		findOrderDetailByDetailIDFn: func(detailID int64, inputCustID string) (model.OrderDetailRead, error) {
			if detailID != orderDetailAID {
				return model.OrderDetailRead{}, errors.New("unexpected detail id")
			}
			return model.OrderDetailRead{
				OrderDetailID:   &orderDetailAIDInt,
				RoNo:            roNo,
				ProId:           748,
				Qty:             &qtyAOld,
				Qty1:            float64Ptr(2),
				Qty2:            float64Ptr(0),
				Qty3:            float64Ptr(1),
				QtyFinal:        &qtyAOld,
				Qty1Final:       float64Ptr(2),
				Qty2Final:       float64Ptr(0),
				Qty3Final:       float64Ptr(1),
				SellPrice1:      &sellPriceA1,
				SellPrice2:      &sellPriceA2,
				SellPrice3:      &sellPriceA3,
				SellPriceFinal1: &sellPriceA1,
				SellPriceFinal2: &sellPriceA2,
				SellPriceFinal3: &sellPriceA3,
				Vat:             &vat,
				MpConvUnit2:     &convUnit2,
				MpConvUnit3:     &convUnit3,
			}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5, SellPrice1: 50, SellPrice2: 0, SellPrice3: 0}, nil
		},
		updateDetailPartialFn: func(c context.Context, detailID int64, inputCustID string, updates map[string]interface{}) error {
			copied := map[string]interface{}{}
			for k, v := range updates {
				copied[k] = v
			}
			partialCaptures = append(partialCaptures, detailPartialCapture{id: detailID, updates: copied})
			return nil
		},
		findOrderDetailsForProformaFn: func(ctx context.Context, roNos []string, inputCustID string) ([]model.OrderDetailRead, error) {
			rowBID := intPtrForTest(1002)
			qtyB := 10.0
			zero := 0.0
			priceB := 50.0
			return []model.OrderDetailRead{
				{
					OrderDetailID:   &orderDetailAIDInt,
					ProId:           748,
					Qty1Final:       &qtyAZero,
					Qty2Final:       &qtyAZero,
					Qty3Final:       &qtyAZero,
					SellPriceFinal1: &sellPriceA1,
					SellPriceFinal2: &sellPriceA2,
					SellPriceFinal3: &sellPriceA3,
					Vat:             &vat,
				},
				{
					OrderDetailID:   rowBID,
					ProId:           749,
					Qty1Final:       &qtyB,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceB,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vat,
				},
			}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			headerUpdates = append(headerUpdates, data)
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1002
			data.OrderDetailID = &id
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdatesInput []*entity.SalesOrderStockUpdate) error {
			stockUpdates = stockUpdatesInput
			return nil
		},
		getCurrentStockFn: func(c context.Context, custID string, whID int64, proID int64) (float64, error) {
			return 0, nil
		},
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: stockRepo,
		Transaction:     &mockDbtransaction{},
	}

	req := entity.EditOrderEnhanceBody{
		CustId: custID,
		SalesOrder: []entity.EditSalesOrderDetail{{
			OrderDetailId: orderDetailAID,
			Qty1:          &qtyAZero,
			Qty2:          &qtyAZero,
			Qty3:          &qtyAZero,
		}},
		AddSalesOrder: []entity.AddSalesOrderDetail{{
			ProId:            749,
			Qty1:             10,
			Qty2:             0,
			Qty3:             0,
			SellPriceSystem1: 50,
			SellPriceSystem2: 0,
			SellPriceSystem3: 0,
			SellPrice1:       50,
			SellPrice2:       0,
			SellPrice3:       0,
			UnitId1:          "PCS",
			UnitId2:          "BOX",
			UnitId3:          "CRT",
		}},
	}

	if err := service.UpdateEnhance(context.Background(), roNo, req); err != nil {
		t.Fatalf("UpdateEnhance returned error: %v", err)
	}

	if len(stockUpdates) < 2 {
		t.Fatalf("expected stock updates for row A release and row B insert, got %d", len(stockUpdates))
	}
	foundRowARelease := false
	for _, update := range stockUpdates {
		if update.RefDetId == orderDetailAID {
			if update.QtyOrderBefore != nil && *update.QtyOrderBefore == qtyAOld && update.QtyOrder == 0 {
				foundRowARelease = true
			}
		}
	}
	if !foundRowARelease {
		t.Fatalf("expected stock release for row A qty 0 transition, got %+v", stockUpdates)
	}

	if len(headerUpdates) == 0 {
		t.Fatalf("expected header update after recompute")
	}
	lastHeader := headerUpdates[len(headerUpdates)-1]
	if lastHeader.SubTotalFinal == nil || *lastHeader.SubTotalFinal != 500 {
		t.Fatalf("expected final subtotal to only include active row B, got %+v", lastHeader.SubTotalFinal)
	}
	if lastHeader.VatValueFinal == nil || *lastHeader.VatValueFinal != 55 {
		t.Fatalf("expected final vat to only include active row B, got %+v", lastHeader.VatValueFinal)
	}
	if lastHeader.TotalFinal == nil || *lastHeader.TotalFinal != 555 {
		t.Fatalf("expected final total to only include active row B, got %+v", lastHeader.TotalFinal)
	}

	rowAFinalSnapshotRecomputeCalls := 0
	for _, capture := range partialCaptures {
		if capture.id == orderDetailAID {
			if _, exists := capture.updates["promo_value_final"]; exists {
				rowAFinalSnapshotRecomputeCalls++
			}
		}
	}
	if rowAFinalSnapshotRecomputeCalls > 1 {
		t.Fatalf("deleted-effective row A must not be recomputed into final promo snapshot, got %d calls", rowAFinalSnapshotRecomputeCalls)
	}
}

func TestUpdateEnhance_SalesOrderDeletedEffectiveQtyZero_VatUsesOnlyNewActiveRows(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603180004"
	whID := int64(301)
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	orderDetailAID := int64(1001)
	orderDetailAIDInt := 1001
	convUnit2 := 10
	convUnit3 := 5
	qtyAOld := 10.0
	qtyZero := 0.0
	priceA1 := 100.0
	priceNew1 := 200.0
	vatA := 0.0
	vatB := 11.0
	zero := 0.0

	var headerUpdates []model.Order
	var stockUpdates []*entity.SalesOrderStockUpdate

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{CustID: custID, WhId: &whID, RoDate: &roDate}, nil
		},
		findOrderDetailByDetailIDFn: func(detailID int64, inputCustID string) (model.OrderDetailRead, error) {
			if detailID != orderDetailAID {
				return model.OrderDetailRead{}, errors.New("unexpected detail id")
			}
			return model.OrderDetailRead{
				OrderDetailID:   &orderDetailAIDInt,
				RoNo:            roNo,
				ProId:           748,
				Qty:             &qtyAOld,
				Qty1:            &qtyAOld,
				Qty2:            &qtyZero,
				Qty3:            &qtyZero,
				QtyFinal:        &qtyAOld,
				Qty1Final:       &qtyAOld,
				Qty2Final:       &qtyZero,
				Qty3Final:       &qtyZero,
				SellPrice1:      &priceA1,
				SellPriceFinal1: &priceA1,
				Vat:             &vatA,
				MpConvUnit2:     &convUnit2,
				MpConvUnit3:     &convUnit3,
			}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5, SellPrice1: priceNew1, SellPrice2: 0, SellPrice3: 0}, nil
		},
		updateDetailPartialFn: func(c context.Context, detailID int64, inputCustID string, updates map[string]interface{}) error {
			return nil
		},
		findOrderDetailsForProformaFn: func(ctx context.Context, roNos []string, inputCustID string) ([]model.OrderDetailRead, error) {
			rowBID := intPtrForTest(1002)
			qtyB := 5.0
			return []model.OrderDetailRead{
				{
					OrderDetailID:   &orderDetailAIDInt,
					ProId:           748,
					Qty1Final:       &qtyZero,
					Qty2Final:       &qtyZero,
					Qty3Final:       &qtyZero,
					SellPriceFinal1: &priceA1,
					Vat:             &vatA,
				},
				{
					OrderDetailID:   rowBID,
					ProId:           749,
					Qty1Final:       &qtyB,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceNew1,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vatB,
				},
			}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			headerUpdates = append(headerUpdates, data)
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1002
			data.OrderDetailID = &id
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdatesInput []*entity.SalesOrderStockUpdate) error {
			stockUpdates = stockUpdatesInput
			return nil
		},
		getCurrentStockFn: func(c context.Context, custID string, whID int64, proID int64) (float64, error) {
			return 0, nil
		},
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: stockRepo,
		Transaction:     &mockDbtransaction{},
	}

	req := entity.EditOrderEnhanceBody{
		CustId: custID,
		SalesOrder: []entity.EditSalesOrderDetail{{
			OrderDetailId: orderDetailAID,
			Qty1:          &qtyZero,
			Qty2:          &qtyZero,
			Qty3:          &qtyZero,
		}},
		AddSalesOrder: []entity.AddSalesOrderDetail{{
			ProId:            749,
			Qty1:             5,
			Qty2:             0,
			Qty3:             0,
			SellPriceSystem1: 200,
			SellPriceSystem2: 0,
			SellPriceSystem3: 0,
			SellPrice1:       200,
			SellPrice2:       0,
			SellPrice3:       0,
			UnitId1:          "PCS",
			UnitId2:          "BOX",
			UnitId3:          "CRT",
		}},
	}

	if err := service.UpdateEnhance(context.Background(), roNo, req); err != nil {
		t.Fatalf("UpdateEnhance returned error: %v", err)
	}

	if len(headerUpdates) == 0 {
		t.Fatalf("expected header update after recompute")
	}
	lastHeader := headerUpdates[len(headerUpdates)-1]
	if lastHeader.SubTotalFinal == nil || *lastHeader.SubTotalFinal != 1000 {
		t.Fatalf("expected final subtotal to only include new vat row, got %+v", lastHeader.SubTotalFinal)
	}
	if lastHeader.VatValueFinal == nil || *lastHeader.VatValueFinal != 110 {
		t.Fatalf("expected final vat to only include active vat row, got %+v", lastHeader.VatValueFinal)
	}
	if lastHeader.TotalFinal == nil || *lastHeader.TotalFinal != 1110 {
		t.Fatalf("expected final total to only include active vat row, got %+v", lastHeader.TotalFinal)
	}

	foundRelease := false
	foundInsert := false
	for _, update := range stockUpdates {
		if update.RefDetId == orderDetailAID && update.QtyOrderBefore != nil && *update.QtyOrderBefore == qtyAOld && update.QtyOrder == 0 {
			foundRelease = true
		}
		if update.ProID == 749 && update.QtyOrderBefore == nil && update.QtyOrder == 5 {
			foundInsert = true
		}
	}
	if !foundRelease {
		t.Fatalf("expected old non-vat row stock release, got %+v", stockUpdates)
	}
	if !foundInsert {
		t.Fatalf("expected new vat row stock insert, got %+v", stockUpdates)
	}
}

func TestDetailV2_FinalOrderFiltersDeletedEffectiveQtyZeroRows_MultiItemScenario(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603180005"
	whID := int64(301)
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	zero := 0.0
	vat := 11.0
	priceA1 := 100.0
	priceB1 := 120.0
	priceC1 := 50.0
	qtyA := 2.0
	qtyC := 4.0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, WhId: &whID}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			rowAID := intPtrForTest(1001)
			rowBID := intPtrForTest(1002)
			rowCID := intPtrForTest(1003)
			return []model.OrderDetailRead{
				{
					OrderDetailID:   rowAID,
					RoNo:            roNo,
					ProId:           748,
					ProCode:         "PRO-748",
					ProName:         "Product A",
					ItemType:        1,
					Qty:             &qtyA,
					Qty1:            &qtyA,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceA1,
					SellPrice2:      &zero,
					SellPrice3:      &zero,
					QtyFinal:        &qtyA,
					Qty1Final:       &qtyA,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceA1,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
				{
					OrderDetailID:   rowBID,
					RoNo:            roNo,
					ProId:           749,
					ProCode:         "PRO-749",
					ProName:         "Product B",
					ItemType:        1,
					Qty:             &zero,
					Qty1:            &zero,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceB1,
					SellPrice2:      &zero,
					SellPrice3:      &zero,
					QtyFinal:        &zero,
					Qty1Final:       &zero,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceB1,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
				{
					OrderDetailID:   rowCID,
					RoNo:            roNo,
					ProId:           750,
					ProCode:         "PRO-750",
					ProName:         "Product C",
					ItemType:        1,
					Qty:             &qtyC,
					Qty1:            &qtyC,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceC1,
					SellPrice2:      &zero,
					SellPrice3:      &zero,
					QtyFinal:        &qtyC,
					Qty1Final:       &qtyC,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceC1,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
			}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) {
			return nil, nil
		},
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0, 749: 0, 750: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}

	if len(response.DetailsFinal.Normal) != 2 {
		t.Fatalf("expected final normal details to include only active rows A and C, got %d", len(response.DetailsFinal.Normal))
	}
	if response.DetailsFinal.Normal[0].ProId != 748 || response.DetailsFinal.Normal[1].ProId != 750 {
		t.Fatalf("expected final normal details to keep active rows A and C only, got %+v", response.DetailsFinal.Normal)
	}
	if response.DetailsFinal.Normal[0].VatValueFinal == nil || *response.DetailsFinal.Normal[0].VatValueFinal != 22 {
		t.Fatalf("expected row A final VAT 22, got %+v", response.DetailsFinal.Normal[0].VatValueFinal)
	}
	if response.DetailsFinal.Normal[1].VatValueFinal == nil || *response.DetailsFinal.Normal[1].VatValueFinal != 22 {
		t.Fatalf("expected row C final VAT 22, got %+v", response.DetailsFinal.Normal[1].VatValueFinal)
	}
}

func TestDetailV2_FinalOrderFiltersDeletedEffectiveQtyZeroRows(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603180002"
	whID := int64(301)
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	zero := 0.0
	vat := 11.0
	priceA1 := 100.0
	priceA2 := 0.0
	priceA3 := 1000.0
	priceB1 := 50.0
	qtyB := 1.0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, WhId: &whID}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			rowAID := intPtrForTest(1001)
			rowBID := intPtrForTest(1002)
			return []model.OrderDetailRead{
				{
					OrderDetailID:   rowAID,
					RoNo:            roNo,
					ProId:           748,
					ProCode:         "PRO-748",
					ProName:         "Product A",
					ItemType:        1,
					Qty:             &zero,
					Qty1:            &zero,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceA1,
					SellPrice2:      &priceA2,
					SellPrice3:      &priceA3,
					QtyFinal:        &zero,
					Qty1Final:       &zero,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceA1,
					SellPriceFinal2: &priceA2,
					SellPriceFinal3: &priceA3,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
				{
					OrderDetailID:   rowBID,
					RoNo:            roNo,
					ProId:           749,
					ProCode:         "PRO-749",
					ProName:         "Product B",
					ItemType:        1,
					Qty:             &qtyB,
					Qty1:            &qtyB,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceB1,
					SellPrice2:      &zero,
					SellPrice3:      &zero,
					QtyFinal:        &qtyB,
					Qty1Final:       &qtyB,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceB1,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
			}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) {
			return nil, nil
		},
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0, 749: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}

	if len(response.DetailsFinal.Normal) != 1 {
		t.Fatalf("expected final normal details to include only active rows, got %d", len(response.DetailsFinal.Normal))
	}
	if response.DetailsFinal.Normal[0].ProId != 749 {
		t.Fatalf("expected only row B in final normal details, got %+v", response.DetailsFinal.Normal)
	}
	if response.DetailsFinal.Normal[0].VatValueFinal == nil || *response.DetailsFinal.Normal[0].VatValueFinal != 6 {
		t.Fatalf("expected final VAT to be derived from row B only, got %+v", response.DetailsFinal.Normal[0].VatValueFinal)
	}

	if response.DetailsFinal.Normal[0].QtyFinal == nil || *response.DetailsFinal.Normal[0].QtyFinal != 1 {
		t.Fatalf("expected final qty_order to remain on active row B only, got %+v", response.DetailsFinal.Normal[0].QtyFinal)
	}

	if response.DetailsFinal.Normal[0].Qty1Stok == nil || *response.DetailsFinal.Normal[0].Qty1Stok != 0 {
		t.Fatalf("expected final initial_stock (qty1_stok) on active row B to be 1, got %+v", response.DetailsFinal.Normal[0].Qty1Stok)
	}
	if response.DetailsFinal.Normal[0].Qty2Stok == nil || *response.DetailsFinal.Normal[0].Qty2Stok != 0 {
		t.Fatalf("expected final initial_stock (qty2_stok) on active row B to be 0, got %+v", response.DetailsFinal.Normal[0].Qty2Stok)
	}
	if response.DetailsFinal.Normal[0].Qty3Stok == nil || *response.DetailsFinal.Normal[0].Qty3Stok != 1 {
		t.Fatalf("expected final initial_stock (qty3_stok) on active row B to be 0, got %+v", response.DetailsFinal.Normal[0].Qty3Stok)
	}
}

func TestDetailNoCustID_FiltersDeletedEffectiveQtyZeroRows(t *testing.T) {
	custIDOrigin := "C220010001"
	roNo := "SO2603180003"
	oprType := "C"
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	zero := 0.0
	qtyB := 10.0
	priceA1 := 100.0
	priceA2 := 0.0
	priceA3 := 1000.0
	priceB1 := 50.0
	vat := 11.0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoNoCustIDFn: func(inputRoNo string, inputCustIDOrigin string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, OprType: &oprType, CustID: custIDOrigin, RoDate: &roDate}, nil
		},
		findDetailNoCustIDFn: func(inputRoNo string, inputCustIDOrigin string) ([]model.OrderDetailRead, error) {
			rowAID := intPtrForTest(1001)
			rowBID := intPtrForTest(1002)
			return []model.OrderDetailRead{
				{
					OrderDetailID:   rowAID,
					RoNo:            roNo,
					ProId:           748,
					ProCode:         "PRO-748",
					ProName:         "Product A",
					ItemType:        1,
					Qty:             &zero,
					Qty1:            &zero,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceA1,
					SellPrice2:      &priceA2,
					SellPrice3:      &priceA3,
					QtyFinal:        &zero,
					Qty1Final:       &zero,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceA1,
					SellPriceFinal2: &priceA2,
					SellPriceFinal3: &priceA3,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
				{
					OrderDetailID:   rowBID,
					RoNo:            roNo,
					ProId:           749,
					ProCode:         "PRO-749",
					ProName:         "Product B",
					ItemType:        1,
					Qty:             &qtyB,
					Qty1:            &qtyB,
					Qty2:            &zero,
					Qty3:            &zero,
					SellPrice1:      &priceB1,
					SellPrice2:      &zero,
					SellPrice3:      &zero,
					QtyFinal:        &qtyB,
					Qty1Final:       &qtyB,
					Qty2Final:       &zero,
					Qty3Final:       &zero,
					SellPriceFinal1: &priceB1,
					SellPriceFinal2: &zero,
					SellPriceFinal3: &zero,
					Vat:             &vat,
					MpConvUnit2:     &convUnit2,
					MpConvUnit3:     &convUnit3,
				},
			}, nil
		},
		findRewardNoCustIDFn: func(inputRoNo string, inputCustIDOrigin string) ([]model.OrderRewardRead, error) {
			return nil, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailNoCustID(roNo, custIDOrigin, nil)
	if err != nil {
		t.Fatalf("DetailNoCustID returned error: %v", err)
	}

	if len(response.Details.Normal) != 1 {
		t.Fatalf("expected sales normal details to include only active rows, got %d", len(response.Details.Normal))
	}
	if response.Details.Normal[0].ProId != 749 {
		t.Fatalf("expected only row B in sales normal details, got %+v", response.Details.Normal)
	}
	if len(response.DetailsFinal.Normal) != 1 {
		t.Fatalf("expected final normal details to include only active rows, got %d", len(response.DetailsFinal.Normal))
	}
	if response.DetailsFinal.Normal[0].ProId != 749 {
		t.Fatalf("expected only row B in final normal details, got %+v", response.DetailsFinal.Normal)
	}
	if response.OprType == nil || *response.OprType != oprType {
		t.Fatalf("expected opr_type=%s, got %+v", oprType, response.OprType)
	}
}

func TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603180999"
	whID := int64(301)
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	dataStatus := int64(entity.CANCELLED)
	convUnit2 := 10
	convUnit3 := 5
	vat := 11.0

	qty52 := 52.0
	qty1Two := 2.0
	qty2Zero := 0.0
	qty3One := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, WhId: &whID, DataStatus: &dataStatus}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			rowID := intPtrForTest(1001)
			return []model.OrderDetailRead{{
				OrderDetailID:   rowID,
				RoNo:            roNo,
				ProId:           748,
				ProCode:         "PRO-748",
				ProName:         "Product A",
				ItemType:        1,
				Qty:             &qty52,
				Qty1:            &qty1Two,
				Qty2:            &qty2Zero,
				Qty3:            &qty3One,
				QtyFinal:        &qty52,
				Qty1Final:       &qty1Two,
				Qty2Final:       &qty2Zero,
				Qty3Final:       &qty3One,
				SellPrice1:      &sellPrice1,
				SellPrice2:      &sellPrice2,
				SellPrice3:      &sellPrice3,
				SellPriceFinal1: &sellPrice1,
				SellPriceFinal2: &sellPrice2,
				SellPriceFinal3: &sellPrice3,
				Vat:             &vat,
				MpConvUnit2:     &convUnit2,
				MpConvUnit3:     &convUnit3,
			}}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 52}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if len(response.Details.Normal) != 1 {
		t.Fatalf("expected 1 sales normal detail, got %d", len(response.Details.Normal))
	}
	if len(response.DetailsFinal.Normal) != 1 {
		t.Fatalf("expected 1 final normal detail, got %d", len(response.DetailsFinal.Normal))
	}

	sales := response.Details.Normal[0]
	if sales.Qty1Stok == nil || *sales.Qty1Stok != 1 {
		t.Fatalf("cancelled sales qty1_stok must use warehouse current only, got %+v", sales.Qty1Stok)
	}
	if sales.Qty2Stok == nil || *sales.Qty2Stok != 0 {
		t.Fatalf("cancelled sales qty2_stok must use warehouse current only, got %+v", sales.Qty2Stok)
	}
	if sales.Qty3Stok == nil || *sales.Qty3Stok != 2 {
		t.Fatalf("cancelled sales qty3_stok must use warehouse current only, got %+v", sales.Qty3Stok)
	}

	final := response.DetailsFinal.Normal[0]
	if final.Qty1Stok == nil || *final.Qty1Stok != 1 {
		t.Fatalf("cancelled final qty1_stok must use warehouse current only, got %+v", final.Qty1Stok)
	}
	if final.Qty2Stok == nil || *final.Qty2Stok != 0 {
		t.Fatalf("cancelled final qty2_stok must use warehouse current only, got %+v", final.Qty2Stok)
	}
	if final.Qty3Stok == nil || *final.Qty3Stok != 2 {
		t.Fatalf("cancelled final qty3_stok must use warehouse current only, got %+v", final.Qty3Stok)
	}
}

func TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603181000"
	whID := int64(301)
	roDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	dataStatus := int64(entity.PROCESSED)
	convUnit2 := 10
	convUnit3 := 5
	vat := 11.0

	qty52 := 52.0
	qty1Two := 2.0
	qty2Zero := 0.0
	qty3One := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, WhId: &whID, DataStatus: &dataStatus}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			rowID := intPtrForTest(1001)
			return []model.OrderDetailRead{{
				OrderDetailID:   rowID,
				RoNo:            roNo,
				ProId:           748,
				ProCode:         "PRO-748",
				ProName:         "Product A",
				ItemType:        1,
				Qty:             &qty52,
				Qty1:            &qty1Two,
				Qty2:            &qty2Zero,
				Qty3:            &qty3One,
				QtyFinal:        &qty52,
				Qty1Final:       &qty1Two,
				Qty2Final:       &qty2Zero,
				Qty3Final:       &qty3One,
				SellPrice1:      &sellPrice1,
				SellPrice2:      &sellPrice2,
				SellPrice3:      &sellPrice3,
				SellPriceFinal1: &sellPrice1,
				SellPriceFinal2: &sellPrice2,
				SellPriceFinal3: &sellPrice3,
				Vat:             &vat,
				MpConvUnit2:     &convUnit2,
				MpConvUnit3:     &convUnit3,
			}}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 52}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if len(response.Details.Normal) != 1 {
		t.Fatalf("expected 1 sales normal detail, got %d", len(response.Details.Normal))
	}
	if len(response.DetailsFinal.Normal) != 1 {
		t.Fatalf("expected 1 final normal detail, got %d", len(response.DetailsFinal.Normal))
	}

	sales := response.Details.Normal[0]
	if sales.Qty1Stok == nil || *sales.Qty1Stok != 2 {
		t.Fatalf("non-cancelled sales qty1_stok must use canonical large mapping, got %+v", sales.Qty1Stok)
	}
	if sales.Qty2Stok == nil || *sales.Qty2Stok != 0 {
		t.Fatalf("non-cancelled sales qty2_stok must use canonical medium mapping, got %+v", sales.Qty2Stok)
	}
	if sales.Qty3Stok == nil || *sales.Qty3Stok != 4 {
		t.Fatalf("non-cancelled sales qty3_stok must use canonical small mapping, got %+v", sales.Qty3Stok)
	}

	final := response.DetailsFinal.Normal[0]
	if final.Qty1Stok == nil || *final.Qty1Stok != 2 {
		t.Fatalf("non-cancelled final qty1_stok must use canonical large mapping, got %+v", final.Qty1Stok)
	}
	if final.Qty2Stok == nil || *final.Qty2Stok != 0 {
		t.Fatalf("non-cancelled final qty2_stok must use canonical medium mapping, got %+v", final.Qty2Stok)
	}
	if final.Qty3Stok == nil || *final.Qty3Stok != 4 {
		t.Fatalf("non-cancelled final qty3_stok must use canonical small mapping, got %+v", final.Qty3Stok)
	}
}

func TestStore_DoesNotPersistStockSnapshotDuringInitialCreate(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-03-11"
	zero := 0.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0

	var storedDetails []model.OrderDetail

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findWarehouseStockByWhIdAndProIdsFn: func(custID string, whID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 100, 990: 100}, nil
		},
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error { return nil },
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := len(storedDetails) + 1
			data.OrderDetailID = &id
			storedDetails = append(storedDetails, *data)
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	stockRepo := &mockStockRepository{salesStockUpdatesFn: func(c context.Context, updates []*entity.SalesOrderStockUpdate) error { return nil }}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(5),
			ConvUnit3:       intPtr(1),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if len(storedDetails) != 1 {
		t.Fatalf("expected 1 stored detail, got %d", len(storedDetails))
	}
	if storedDetails[0].Qty1Stok != nil || storedDetails[0].Qty2Stok != nil || storedDetails[0].Qty3Stok != nil {
		t.Fatalf("expected create flow to avoid persisting stock snapshot on initial store, got %+v %+v %+v", storedDetails[0].Qty1Stok, storedDetails[0].Qty2Stok, storedDetails[0].Qty3Stok)
	}
}

func TestUpdateEnhance_FinalOrderPersistsFinalPromoSnapshot(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603110102"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC)
	orderDetailID := int64(1002)
	orderDetailIDInt := 1002
	convUnit2 := 10
	convUnit3 := 5
	qty1Final := 2.0
	qty2Final := 0.0
	qty3Final := 1.0
	qtyFinal := 52.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0

	var detailUpdates []map[string]interface{}
	var headerUpdates []model.Order

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{CustID: custID, WhId: &whID, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID}, nil
		},
		findOutletByIDFn: func(outletID int, custID string, parentCustID string) (model.OutletRead, error) {
			return model.OutletRead{OutletId: outletID}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		findOrderDetailByDetailIDFn: func(detailID int64, inputCustID string) (model.OrderDetailRead, error) {
			return model.OrderDetailRead{OrderDetailID: &orderDetailIDInt, RoNo: roNo, ProId: 748, QtyFinal: &qtyFinal, Qty1Final: &qty1Final, Qty2Final: &qty2Final, Qty3Final: &qty3Final, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, Vat: &vat, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5, SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3}, nil
		},
		updateDetailPartialFn: func(c context.Context, detailID int64, inputCustID string, updates map[string]interface{}) error {
			copied := map[string]interface{}{}
			for k, v := range updates {
				copied[k] = v
			}
			detailUpdates = append(detailUpdates, copied)
			return nil
		},
		findOrderDetailsForProformaFn: func(ctx context.Context, roNos []string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{OrderDetailID: &orderDetailIDInt, ProId: 748, Qty1Final: &qty1Final, Qty2Final: &qty2Final, Qty3Final: &qty3Final, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, DiscValueFinal: float64Ptr(0), VatValueFinal: float64Ptr(131), Vat: &vat}}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			headerUpdates = append(headerUpdates, data)
			return nil
		},
	}

	promotionRepo := &mockPromotionRepositoryEnhance{findProductByIDAndCustIDFn: func(productID int64, custID string) (model.ProductRead, error) {
		return model.ProductRead{ProId: int(productID), ConvUnit2: 10, ConvUnit3: 5}, nil
	}}

	promotionV2Repo := &mockPromotionV2RepositoryEnhance{
		findOutletByIDFn: func(outletID int64, custID string) (model.OutletPromo, error) {
			return model.OutletPromo{OutletID: int(outletID)}, nil
		},
		findSalesmanByIDFn: func(salesmanID int64, custID string) (model.SalesmanPromo, error) {
			return model.SalesmanPromo{WhId: int(whID)}, nil
		},
		findWarehouseByIDFn: func(warehouseID int64, custID string) (model.WarehousePromo, error) {
			return model.WarehousePromo{WhID: int(warehouseID)}, nil
		},
		findActivePromotionsByOutletFn: func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
			return []model.PromotionV2{{PromoID: "PROMO-FINAL-1", PromoDesc: "Promo Final 1", PromoType: model.PromotionTypeSlab}}, nil
		},
		findProductCriteriasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) { return nil, nil },
		findSlabsByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
			rewardValue := 10.0
			perScope := string(model.PerScopeProduct)
			return []model.PromotionV2Slabs{{PromoID: "PROMO-FINAL-1", RuleType: model.RuleTypeQuantity, RewardType: model.RewardTypeFixedValue, RewardValue: &rewardValue, PerScope: &perScope, RangeTo: 1}}, nil
		},
		findStratasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) { return nil, nil },
		getAllRewardProductFromStockV2Fn: func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
			return []model.PromotionRewardProduct{{ProID: 990, QtyStock: 100, ConvUnit2: 10, ConvUnit3: 5}}, nil
		},
	}

	stockRepo := &mockStockRepository{
		getCancelStockBasisFn: func(c context.Context, inputCustID string, orderNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{{
				CustID:         custID,
				WhID:           whID,
				ProID:          748,
				RefDetID:       orderDetailID,
				StockDate:      roDate,
				UnitPrice:      sellPrice1,
				QtyOutstanding: 1,
				QtyOutSmallest: 1,
			}}, nil
		},
		salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error { return nil },
		getCurrentStockFn:   func(c context.Context, custID string, whID int64, proID int64) (float64, error) { return 100, nil },
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, PromotionRepository: promotionRepo, PromotionV2Repository: promotionV2Repo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}

	req := entity.EditOrderEnhanceBody{CustId: custID, ParentCustId: custID, FinalOrder: []entity.EditFinalOrderDetail{{OrderDetailId: orderDetailID, Qty1Final: &qty1Final, Qty2Final: &qty2Final, Qty3Final: &qty3Final, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3}}}

	if err := service.UpdateEnhance(context.Background(), roNo, req); err != nil {
		t.Fatalf("UpdateEnhance returned error: %v", err)
	}
	if len(detailUpdates) == 0 {
		t.Fatalf("expected final order detail updates to be executed")
	}
	foundFinalSnapshot := false
	for _, updates := range detailUpdates {
		if _, ok := updates["promo_final1"]; ok {
			foundFinalSnapshot = true
			if updates["promo_final1"] == float64(0) {
				t.Fatalf("expected final promo snapshot to be persisted")
			}
		}
	}
	if !foundFinalSnapshot {
		t.Fatalf("expected final promo snapshot updates")
	}
	foundHeaderSnapshot := false
	for _, update := range headerUpdates {
		if len(update.PromoRemarksFinal) > 0 {
			foundHeaderSnapshot = true
			if update.PromoRemarksFinal[0] != "PROMO-FINAL-1" {
				t.Fatalf("expected header final promo remarks, got %+v", update.PromoRemarksFinal)
			}
		}
	}
	if !foundHeaderSnapshot {
		t.Fatalf("expected header final promo snapshot update")
	}
}

func TestUpdateEnhance_AddFinalOrder_MissingBasisShouldFailAndSkipWrites(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603220005"
	whID := int64(301)
	roDate := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)

	storeDetailCalled := 0
	stockWriteCalled := 0
	updateHeaderCalled := 0

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, WhId: &whID, RoDate: &roDate}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5, UnitId1: "PCS", UnitId2: "BOX", UnitId3: "CRT"}, nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			storeDetailCalled++
			id := 9991
			data.OrderDetailID = &id
			return nil
		},
		findOrderDetailsForProformaFn: func(ctx context.Context, roNos []string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			updateHeaderCalled++
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		getCancelStockBasisFn: func(c context.Context, inputCustID string, orderNo string) ([]entity.CancelStockBasis, error) {
			return []entity.CancelStockBasis{}, nil
		},
		salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled++
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.UpdateEnhance(context.Background(), roNo, entity.EditOrderEnhanceBody{
		CustId: custID,
		AddFinalOrder: []entity.AddFinalOrderDetail{{
			ProId:            748,
			Qty1Final:        1,
			Qty2Final:        0,
			Qty3Final:        0,
			SellPriceSystem1: 100,
			SellPriceSystem2: 0,
			SellPriceSystem3: 0,
			SellPriceFinal1:  100,
			SellPriceFinal2:  0,
			SellPriceFinal3:  0,
		}},
	})
	if err == nil {
		t.Fatalf("expected add_final_order to fail when stock basis is not synchronized")
	}
	if !strings.Contains(err.Error(), "stock basis is not synchronized") {
		t.Fatalf("expected synchronized stock basis error, got %v", err)
	}
	if storeDetailCalled != 0 {
		t.Fatalf("expected no detail write when stock basis is inconsistent, got %d", storeDetailCalled)
	}
	if stockWriteCalled != 0 {
		t.Fatalf("expected no stock write when stock basis is inconsistent, got %d", stockWriteCalled)
	}
	if updateHeaderCalled != 0 {
		t.Fatalf("expected no header update when stock basis is inconsistent, got %d", updateHeaderCalled)
	}
}

func TestBuildCreateOrderRewardDetails_BuildsPromoLinesFromConsultV2(t *testing.T) {
	consultResp := []entity.ConsultPromoResp{{
		PromoID: "PROMO-RWD",
		RewardProduct: []entity.PromoRewardProductDet{{
			ProID:      990,
			Qty1:       1,
			Qty2:       2,
			Qty3:       3,
			GrossValue: 450,
		}},
	}}

	details, promoBgTotal, err := buildCreateOrderRewardDetails(consultResp, func(productID int) (model.ProductRead, error) {
		return model.ProductRead{ProId: productID, SellPrice1: 50, SellPrice2: 100, SellPrice3: 200, ConvUnit2: 10, ConvUnit3: 5}, nil
	})
	if err != nil {
		t.Fatalf("buildCreateOrderRewardDetails returned error: %v", err)
	}
	if len(details) != 1 {
		t.Fatalf("expected 1 reward detail, got %d", len(details))
	}
	detail := details[0]
	if detail.ItemType != 2 {
		t.Fatalf("expected reward detail item_type=2, got %d", detail.ItemType)
	}
	if detail.PromoID == nil || *detail.PromoID != "PROMO-RWD" {
		t.Fatalf("expected promo id to be mapped, got %+v", detail.PromoID)
	}
	if getValueOrDefault(detail.Qty1, 0) != 1 || getValueOrDefault(detail.Qty2, 0) != 2 || getValueOrDefault(detail.Qty3, 0) != 3 {
		t.Fatalf("unexpected reward qty mapping: %+v", detail)
	}
	if getValueOrDefault(detail.PromoValue, 0) != 450 || getValueOrDefault(detail.PromoValueFinal, 0) != 450 {
		t.Fatalf("expected reward promo value to match gross value, got %+v", detail)
	}
	if detail.IsProductPromotionSo == nil || !*detail.IsProductPromotionSo {
		t.Fatalf("expected reward detail sales flag to be true")
	}
	if detail.IsProductPromotionFinal == nil || !*detail.IsProductPromotionFinal {
		t.Fatalf("expected reward detail final flag to be true")
	}
	if promoBgTotal != 450 {
		t.Fatalf("expected promo background total 450, got %.2f", promoBgTotal)
	}
}

func TestDistributePromoToDetailRowsV2_PreservesDuplicateProductRowsByOrderDetailID(t *testing.T) {
	firstDetailID := 1001
	secondDetailID := 1002
	qtyA := 1.0
	qtyB := 3.0
	price := 100.0

	aggregate := map[int]promoAggregateRow{
		748: {
			Promo1:             40,
			PromoTotal:         40,
			Remarks:            []string{"PROMO-DUP"},
			IsProductPromotion: true,
		},
	}
	consultResp := []entity.ConsultPromoResp{{
		PromoID:          "PROMO-DUP",
		ProductsEligible: []int{748},
		RewardValue:      []entity.PromoRewardValue{{ProID: 748, Promo1: 40}},
		SlabPerScope:     string(model.PerScopeProduct),
	}}

	rows := distributePromoToDetailRowsV2(aggregate, []model.OrderDetailRead{
		{OrderDetailID: &firstDetailID, ProId: 748, Qty1: &qtyA, SellPrice1: &price},
		{OrderDetailID: &secondDetailID, ProId: 748, Qty1: &qtyB, SellPrice1: &price},
	}, promoSnapshotTabSalesOrder, consultResp)

	if len(rows) != 2 {
		t.Fatalf("expected 2 distributed rows, got %d", len(rows))
	}
	if rows[firstDetailID].PromoTotal != 10 {
		t.Fatalf("expected first row promo total 10, got %+v", rows[firstDetailID])
	}
	if rows[secondDetailID].PromoTotal != 30 {
		t.Fatalf("expected second row promo total 30, got %+v", rows[secondDetailID])
	}
	if len(rows[firstDetailID].Remarks) != 1 || rows[firstDetailID].Remarks[0] != "PROMO-DUP" {
		t.Fatalf("expected first row remarks to be preserved, got %+v", rows[firstDetailID].Remarks)
	}
	if !rows[secondDetailID].IsProductPromotion {
		t.Fatalf("expected product promotion flag to remain attached per distributed row")
	}
}

func TestCreateOrderDetailFromFinalOrder_PersistsIndependentFinalPromoFlag(t *testing.T) {
	finalFlag := true
	var storedDetail *model.OrderDetail

	orderRepo := &mockOrderRepository{
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			stored := *data
			storedDetail = &stored
			id := 9001
			data.OrderDetailID = &id
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo}
	stockUpdate, err := service.createOrderDetailFromFinalOrder(context.Background(), "SO2603170001", "C220010001", 301, time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC), entity.AddFinalOrderDetail{
		ProId:                   748,
		Qty1Final:               1,
		Qty2Final:               0,
		Qty3Final:               0,
		SellPriceSystem1:        100,
		SellPriceSystem2:        0,
		SellPriceSystem3:        0,
		SellPriceFinal1:         100,
		SellPriceFinal2:         0,
		SellPriceFinal3:         0,
		IsProductPromotionFinal: &finalFlag,
	})
	if err != nil {
		t.Fatalf("createOrderDetailFromFinalOrder returned error: %v", err)
	}
	if storedDetail == nil {
		t.Fatalf("expected final order detail to be stored")
	}
	if storedDetail.IsProductPromotionFinal == nil || !*storedDetail.IsProductPromotionFinal {
		t.Fatalf("expected final promo flag to be persisted independently")
	}
	if storedDetail.IsProductPromotionSo != nil {
		t.Fatalf("expected sales promo flag to remain untouched, got %+v", storedDetail.IsProductPromotionSo)
	}
	if stockUpdate == nil {
		t.Fatalf("expected final order stock update to be built")
	}
	if stockUpdate.RefDetId != 9001 {
		t.Fatalf("expected stock update ref_det_id=9001, got %+v", stockUpdate.RefDetId)
	}
	if stockUpdate.QtyOrder != 1 {
		t.Fatalf("expected stock update qty_order=1, got %+v", stockUpdate.QtyOrder)
	}
}

func intPtrForTest(v int64) *int {
	value := int(v)
	return &value
}

func TestPromoV2Only_NoLegacyConsultPathInSalesOrderFlow(t *testing.T) {
	type fileRule struct {
		name      string
		path      string
		forbidden []string
	}

	rules := []fileRule{
		{
			name: "controller create-update-updatefinal must not use legacy consult",
			path: "../controller/order_controller.go",
			forbidden: []string{
				"SetConsultPromotionRequest(",
				"ConsultPromotion(",
			},
		},
		{
			name: "service create-update-updatefinal-detail must not have legacy consult path",
			path: "order_service.go",
			forbidden: []string{
				"usePromoV2 :=",
				"ConsultPromotionBeforeStore(",
				"SetConsultPromotionRequest(",
				"legacyFallback :=",
			},
		},
	}

	for _, rule := range rules {
		rule := rule
		t.Run(rule.name, func(t *testing.T) {
			content, err := os.ReadFile(rule.path)
			if err != nil {
				t.Fatalf("failed reading %s: %v", rule.path, err)
			}
			text := string(content)
			for _, marker := range rule.forbidden {
				if strings.Contains(text, marker) {
					t.Fatalf("legacy marker %q still exists in %s", marker, rule.path)
				}
			}
		})
	}
}

func TestDetailV2_PreRolloutWithoutSnapshot_MustConsultV2ByTab(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603110002"
	outletID := int64(21)
	salesmanID := int64(11)
	whID := int64(301)
	roDate := time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC)
	convUnit2 := 10
	convUnit3 := 5
	qty := 52.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0
	consultCalled := 0

	orderRepo := &mockOrderRepositoryDetailV2{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, OutletID: &outletID, SalesmanId: &salesmanID, WhId: &whID}, nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{OrderDetailID: intPtr(2), RoNo: roNo, ProId: 748, ProCode: "PRO-748", ProName: "Product 748", ItemType: 1, Qty: &qty, Qty1: &qty1, Qty2: &qty2, Qty3: &qty3, QtyFinal: &qty, Qty1Final: &qty1, Qty2Final: &qty2, Qty3Final: &qty3, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, SellPriceFinal1: &sellPrice1, SellPriceFinal2: &sellPrice2, SellPriceFinal3: &sellPrice3, Vat: &vat, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3}}, nil
		},
		findRewardFn: func(inputRoNo string, inputCustID string) ([]model.OrderRewardRead, error) { return nil, nil },
		findWarehouseStockByWhIdAndProIds: func(inputCustID string, inputWhID int64, proIDs []int64) (map[int64]float64, error) {
			return map[int64]float64{748: 0}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ProCode: "PRO-748", ProName: "Product 748", SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
	}

	promotionRepo := &mockPromotionRepositoryDetailV2{findProductByIDAndCustIDFn: func(productID int64, custID string) (model.ProductRead, error) {
		return model.ProductRead{ProId: int(productID), ConvUnit2: 10, ConvUnit3: 5}, nil
	}}

	promotionV2Repo := &mockPromotionV2RepositoryDetailV2{
		findOutletByIDFn: func(outletID int64, custID string) (model.OutletPromo, error) {
			return model.OutletPromo{OutletID: int(outletID)}, nil
		},
		findSalesmanByIDFn: func(salesmanID int64, custID string) (model.SalesmanPromo, error) {
			return model.SalesmanPromo{WhId: int(whID)}, nil
		},
		findWarehouseByIDFn: func(warehouseID int64, custID string) (model.WarehousePromo, error) {
			return model.WarehousePromo{WhID: int(warehouseID)}, nil
		},
		findActivePromotionsByOutletFn: func(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) ([]model.PromotionV2, error) {
			consultCalled++
			return []model.PromotionV2{{PromoID: "PROMO-V2-PRE", PromoDesc: "Promo Pre Rollout", PromoType: model.PromotionTypeSlab}}, nil
		},
		findProductCriteriasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionProductCriteria, error) {
			return []model.PromotionProductCriteria{{PromoID: "PROMO-V2-PRE", ProID: 748}}, nil
		},
		findSlabsByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Slabs, error) {
			rewardValue := 10.0
			perScope := string(model.PerScopeProduct)
			return []model.PromotionV2Slabs{{PromoID: "PROMO-V2-PRE", RuleType: model.RuleTypeQuantity, RewardType: model.RewardTypeFixedValue, RewardValue: &rewardValue, PerScope: &perScope, RangeTo: 100}}, nil
		},
		findStratasByPromoIDsFn: func(promoIDs []string, custID string) ([]model.PromotionV2Strata, error) { return nil, nil },
		getAllRewardProductFromStockV2Fn: func(req entity.ConsultPromoV2Req, rewardCtx model.RewardContext) ([]model.PromotionRewardProduct, error) {
			return []model.PromotionRewardProduct{{ProID: 990, QtyStock: 100, ConvUnit2: 10, ConvUnit3: 5}}, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, PromotionRepository: promotionRepo, PromotionV2Repository: promotionV2Repo}

	response, err := service.DetailV2(roNo, custID, custID)
	if err != nil {
		t.Fatalf("DetailV2 returned error: %v", err)
	}
	if consultCalled == 0 {
		t.Fatalf("detail without snapshot must consult v2 on pre-rollout row")
	}
	if len(response.Details.Normal) != 1 {
		t.Fatalf("expected 1 normal detail, got %d", len(response.Details.Normal))
	}
	if response.Details.Normal[0].Promo1 == 0 {
		t.Fatalf("detail without snapshot must inject runtime promo, got %+v", response.Details.Normal[0])
	}
	if len(response.Details.FinalRemarks) == 0 {
		t.Fatalf("detail without snapshot must expose remarks from consult v2")
	}
}

func TestConsultDiscountBeforeStore_AppliesDiscountAfterPromo(t *testing.T) {
	promoValue := 200.0
	sellPrice1 := 1000.0
	qty1 := 1.0
	zero := 0.0
	vat := 0.0

	service := &orderServiceImpl{
		OrderRepository: &mockOrderRepositoryStore{
			findProductByIDFn: func(productID int) (model.ProductRead, error) {
				return model.ProductRead{ProId: productID, SellPrice1: sellPrice1}, nil
			},
			findOutletByIDFn: func(outletID int, custID string, parentCustID string) (model.OutletRead, error) {
				return model.OutletRead{OutletId: outletID}, nil
			},
			findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
				return model.DiscountRead{DiscountId: "DISC-10"}, nil
			},
			findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
				return model.DiscountCriteria{DiscountID: discountID, SlabRewardType: entity.DiscountRewardTypePercentage, SlabReward: 10}, nil
			},
		},
	}

	request := &entity.ConsultDiscountOrderBody{
		CustId:       "C220010001",
		ParentCustId: "C220010001",
		OutletID:     99,
		Details: entity.ConsultDiscountOrderDetWithGroup{Normal: []entity.ConsultDiscountOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &zero,
			Qty3:            &zero,
			SellPrice1:      &sellPrice1,
			SellPrice2:      &zero,
			SellPrice3:      &zero,
			PromoValue:      &promoValue,
			PromoValueFinal: &promoValue,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	err := service.ConsultDiscountBeforeStore(request)
	if err != nil {
		t.Fatalf("ConsultDiscountBeforeStore returned error: %v", err)
	}

	got := getValueOrDefault(request.Details.Normal[0].DiscValue, 0)
	want := 80.0 // (1000 - 200) * 10%
	if got != want {
		t.Fatalf("disc_value_item must follow formula (gross - promo_total_item) * disc/100: got %.2f want %.2f", got, want)
	}
}

func TestPromoV2Only_CSVAcceptanceCoverageForReferenceSOs(t *testing.T) {
	type expectedCSVResult struct {
		soNo            string
		promoRemarks    []string
		promoBarang     string
		promoUang       string
		regularDiscount string
	}

	expected := []expectedCSVResult{
		{soNo: "SO2603120003", promoRemarks: []string{"promostrata1"}, promoBarang: "2,000,000", promoUang: "", regularDiscount: ""},
		{soNo: "SO2603120004", promoRemarks: []string{"promostrata1"}, promoBarang: "2,000,000", promoUang: "0", regularDiscount: "0"},
		{soNo: "SO2603120005", promoRemarks: []string{"1112", "BUYBUYBUY"}, promoBarang: "4,000,000", promoUang: "1,840,000", regularDiscount: ""},
		{soNo: "SO2603120006", promoRemarks: []string{"1112", "BUYBUYBUY"}, promoBarang: "4,000,000", promoUang: "984,000", regularDiscount: "628,320"},
	}

	csvData, err := readCSVReferenceFixture("test promo integrasi sales order - Request mas angga.csv")
	if err != nil {
		t.Fatalf("failed reading csv reference: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(csvData)))
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed parsing csv reference: %v", err)
	}

	type sectionSummary struct {
		remarks         []string
		promoBarang     string
		promoUang       string
		regularDiscount string
	}

	summaries := map[string]*sectionSummary{}
	currentSO := ""
	for _, record := range records {
		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}
		if len(record) == 0 {
			continue
		}
		if strings.HasPrefix(record[0], "SO") {
			currentSO = record[0]
			if _, exists := summaries[currentSO]; !exists {
				summaries[currentSO] = &sectionSummary{}
			}
			continue
		}
		if currentSO == "" {
			continue
		}

		if strings.HasPrefix(record[0], "1.") || strings.HasPrefix(record[0], "2.") || strings.HasPrefix(record[0], "3.") {
			summary := summaries[currentSO]
			rowText := strings.ToLower(strings.Join(record, " "))
			for _, marker := range []string{"promostrata1", "BUYBUYBUY"} {
				if strings.Contains(rowText, strings.ToLower(marker)) && !containsString(summary.remarks, marker) {
					summary.remarks = append(summary.remarks, marker)
				}
			}
			if strings.Contains(rowText, "1112") || strings.Contains(rowText, "112") {
				if !containsString(summary.remarks, "1112") {
					summary.remarks = append(summary.remarks, "1112")
				}
			}
		}

		if len(record) > 13 {
			summary := summaries[currentSO]
			switch strings.ToLower(record[12]) {
			case "promo barang":
				summary.promoBarang = record[13]
			case "promo uang":
				summary.promoUang = record[13]
			case "regular diskon":
				summary.regularDiscount = record[13]
			}
		}
	}

	for _, tc := range expected {
		tc := tc
		t.Run(tc.soNo, func(t *testing.T) {
			summary, exists := summaries[tc.soNo]
			if !exists {
				t.Fatalf("SO %s not found in CSV reference", tc.soNo)
			}
			for _, remark := range tc.promoRemarks {
				if !containsString(summary.remarks, remark) {
					t.Fatalf("SO %s missing promo remark %q, got %+v", tc.soNo, remark, summary.remarks)
				}
			}
			if summary.promoBarang != tc.promoBarang {
				t.Fatalf("SO %s promo barang mismatch: got %q want %q", tc.soNo, summary.promoBarang, tc.promoBarang)
			}
			if summary.promoUang != tc.promoUang {
				t.Fatalf("SO %s promo uang mismatch: got %q want %q", tc.soNo, summary.promoUang, tc.promoUang)
			}
			if summary.regularDiscount != tc.regularDiscount {
				t.Fatalf("SO %s regular discount mismatch: got %q want %q", tc.soNo, summary.regularDiscount, tc.regularDiscount)
			}
		})
	}
}

func readCSVReferenceFixture(filename string) ([]byte, error) {
	wd, _ := os.Getwd()
	candidates := []string{
		filepath.Join("..", "docs", filename),
		filepath.Join("..", "..", "docs", filename),
		filepath.Join("..", "..", "..", "docs", filename),
		filepath.Join("..", "..", "..", "scylla-be", "docs", filename),
		filepath.Join("..", "..", "..", "..", "scylla-be", "docs", filename),
		filepath.Join("docs", filename),
		filepath.Join(wd, "..", "docs", filename),
		filepath.Join(wd, "..", "..", "docs", filename),
		filepath.Join(wd, "..", "..", "..", "docs", filename),
		filepath.Join(wd, "..", "..", "..", "scylla-be", "docs", filename),
		filepath.Join(wd, "..", "..", "..", "..", "scylla-be", "docs", filename),
		filepath.Join(wd, "docs", filename),
	}

	var lastErr error
	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}

	return nil, lastErr
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func stringPtrForTest(v string) *string {
	return &v
}

func TestUpdate_MobileNoChange_UnchangedPayload_HeaderOnlyNoStockMutation(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603260001"
	parentCustID := custID
	dataStatus := int64(entity.NEED_REVIEW)
	mobileSource := int64(2)
	roDate := time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)
	whID := int64(301)
	outletID := int64(21)
	salesmanID := int64(11)
	unit1 := "PCS"
	unit2 := "BOX"
	unit3 := "CRT"
	qtyPo1 := 1.0
	qtyPo2 := 2.0
	qtyPo3 := 3.0
	qtyFinal := 321.0
	qtyStored := 123.0
	sellPrice1 := 100.0
	sellPrice2 := 200.0
	sellPrice3 := 300.0
	vat := 11.0
	convUnit2 := 10
	convUnit3 := 12
	orderDetailID := int64(9101)
	orderDetailIDInt := int(orderDetailID)
	amount := 1400.0
	zero := 0.0

	updateHeaderCalled := 0
	findByNotInCalled := 0
	deleteNotInCalled := 0
	stockWriteCalled := 0
	storeDetailCalled := 0
	updateDetailCalled := 0

	orderRepo := &mockOrderRepository{
		findOutletByIDFn: func(outletID int, inputCustID string, inputParentCustID string) (model.OutletRead, error) {
			return model.OutletRead{OutletId: outletID, Address1: stringPtrForTest("Addr A")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: float32(convUnit2), ConvUnit3: float32(convUnit3), SellPrice1: sellPrice1, SellPrice2: sellPrice2, SellPrice3: sellPrice3}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, errors.New("no discount")
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, errors.New("no criteria")
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			updateHeaderCalled++
			return nil
		},
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			if inputRoNo != roNo || inputCustID != custID {
				t.Fatalf("unexpected FindByNo args: %s %s", inputRoNo, inputCustID)
			}
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, WhId: &whID, DataStatus: &dataStatus, DataSource: &mobileSource, OutletID: &outletID, SalesmanId: &salesmanID}, nil
		},
		findDetailByNotInDetailIDsFn: func(detailIDs []int64, inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			findByNotInCalled++
			return []model.OrderDetailRead{}, nil
		},
		deleteDetailNotInIDsFn: func(c context.Context, inputRoNo string, inputCustID string, ids []int64) error {
			deleteNotInCalled++
			return nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{{OrderDetailID: &orderDetailIDInt, ProId: 748, ItemType: 1, QtyFinal: &qtyFinal, QtyPo1: &qtyPo1, QtyPo2: &qtyPo2, QtyPo3: &qtyPo3, Qty: &qtyStored, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, Vat: &vat}}, nil
		},
		findDetailByDetailIDFn: func(detailID int64, inputRoNo string, inputCustID string) (model.OrderDetailRead, error) {
			if detailID != orderDetailID {
				t.Fatalf("unexpected detail id: %d", detailID)
			}
			return model.OrderDetailRead{OrderDetailID: &orderDetailIDInt, RoNo: roNo, ProId: 748, QtyFinal: &qtyFinal, QtyPo1: &qtyPo1, QtyPo2: &qtyPo2, QtyPo3: &qtyPo3, Qty: &qtyStored, SellPrice1: &sellPrice1, SellPrice2: &sellPrice2, SellPrice3: &sellPrice3, Vat: &vat, UnitId1: &unit1, UnitId2: &unit2, UnitId3: &unit3, MpConvUnit2: &convUnit2, MpConvUnit3: &convUnit3}, nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			storeDetailCalled++
			return nil
		},
		updateDetailFn: func(c context.Context, data *model.OrderDetail) error {
			updateDetailCalled++
			return nil
		},
		syncFinalOrderFieldsFn: func(c context.Context, orderDetailId int64) error { return nil },
		deletePromoDetailsFn:   func(c context.Context, inputRoNo string, inputCustID string) error { return nil },
		deleteRewardsFn:        func(c context.Context, inputRoNo string, inputCustID string) error { return nil },
		storeRewardFn:          func(c context.Context, data *model.OrderReward) error { return nil },
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled++
			return nil
		},
		getCurrentStockFn: func(c context.Context, inputCustID string, inputWhID int64, proID int64) (float64, error) {
			return 0, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.Update(roNo, entity.UpdateOrderBody{
		CustId:         custID,
		ParentCustId:   parentCustID,
		RoNo:           roNo,
		RoDate:         stringPtrForTest("2026-03-26"),
		ValDate:        stringPtrForTest("2026-03-26"),
		DueDate:        stringPtrForTest("2026-03-26"),
		WhId:           &whID,
		OutletID:       &outletID,
		SalesmanId:     &salesmanID,
		DataStatus:     intPtrForTest(dataStatus),
		DiscValue:      &zero,
		DiscValueFinal: &zero,
		VatValue:       &zero,
		VatValueFinal:  &zero,
		Details: entity.UpdateOrderDetWithGroup{
			Normal: []entity.UpdateOrderDetBody{{
				OrderDetId:      &orderDetailID,
				ProId:           748,
				Qty1:            &qtyPo1,
				Qty2:            &qtyPo2,
				Qty3:            &qtyPo3,
				QtyPo1:          &qtyPo1,
				QtyPo2:          &qtyPo2,
				QtyPo3:          &qtyPo3,
				SellPrice1:      &sellPrice1,
				SellPrice2:      &sellPrice2,
				SellPrice3:      &sellPrice3,
				UnitId1:         &unit1,
				UnitId2:         &unit2,
				UnitId3:         &unit3,
				ConvUnit2:       &convUnit2,
				ConvUnit3:       &convUnit3,
				PromoValue:      &zero,
				PromoValueFinal: &zero,
				DiscValue:       &zero,
				DiscValueFinal:  &zero,
				Vat:             &vat,
				VatValue:        &zero,
				VatValueFinal:   &zero,
				Amount:          &amount,
				AmountFinal:     &amount,
			}},
			Promo: []entity.CreateOrderDetBody{},
		},
		Rewards: []entity.CreateOrderRewardBody{},
	}, entity.ValidateResponse{})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if updateHeaderCalled != 1 {
		t.Fatalf("expected header update once, got %d", updateHeaderCalled)
	}
	if findByNotInCalled != 0 {
		t.Fatalf("expected FindDetailByNotInDetailIDs to be skipped for mobile no-change, got %d", findByNotInCalled)
	}
	if deleteNotInCalled != 0 {
		t.Fatalf("expected DeleteDetailNotInIDs to be skipped for mobile no-change, got %d", deleteNotInCalled)
	}
	if stockWriteCalled != 0 {
		t.Fatalf("expected stock writes to be skipped for mobile no-change, got %d", stockWriteCalled)
	}
	if storeDetailCalled != 0 {
		t.Fatalf("expected no new detail insert for unchanged payload, got %d", storeDetailCalled)
	}
	if updateDetailCalled != 0 {
		t.Fatalf("expected no detail update for mobile no-change, got %d", updateDetailCalled)
	}
}

func TestUpdate_MobileNoChange_EmptyDetailPayload_HeaderOnlyNoStockMutation(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603260002"
	parentCustID := custID
	dataStatus := int64(entity.PROCESSED)
	mobileSource := int64(2)
	roDate := time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)
	whID := int64(301)
	outletID := int64(21)
	salesmanID := int64(11)
	zero := 0.0

	updateHeaderCalled := 0
	findByNotInCalled := 0
	deleteNotInCalled := 0
	stockWriteCalled := 0

	orderRepo := &mockOrderRepository{
		findOutletByIDFn: func(outletID int, inputCustID string, inputParentCustID string) (model.OutletRead, error) {
			return model.OutletRead{OutletId: outletID, Address1: stringPtrForTest("Addr A")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, ConvUnit2: 10, ConvUnit3: 12, SellPrice1: 100, SellPrice2: 200, SellPrice3: 300}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, errors.New("no discount")
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, errors.New("no criteria")
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			updateHeaderCalled++
			return nil
		},
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{RoNo: roNo, CustID: custID, RoDate: &roDate, WhId: &whID, DataStatus: &dataStatus, DataSource: &mobileSource, OutletID: &outletID, SalesmanId: &salesmanID}, nil
		},
		findDetailByNotInDetailIDsFn: func(detailIDs []int64, inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			findByNotInCalled++
			return []model.OrderDetailRead{}, nil
		},
		deleteDetailNotInIDsFn: func(c context.Context, inputRoNo string, inputCustID string, ids []int64) error {
			deleteNotInCalled++
			return nil
		},
		findDetailFn: func(inputRoNo string, inputCustID string) ([]model.OrderDetailRead, error) {
			return []model.OrderDetailRead{}, nil
		},
		deletePromoDetailsFn: func(c context.Context, inputRoNo string, inputCustID string) error { return nil },
		deleteRewardsFn:      func(c context.Context, inputRoNo string, inputCustID string) error { return nil },
		storeRewardFn:        func(c context.Context, data *model.OrderReward) error { return nil },
		findDetailByDetailIDFn: func(detailID int64, inputRoNo string, inputCustID string) (model.OrderDetailRead, error) {
			return model.OrderDetailRead{}, nil
		},
		updateDetailFn:         func(c context.Context, data *model.OrderDetail) error { return nil },
		storeDetailFn:          func(c context.Context, data *model.OrderDetail) error { return nil },
		syncFinalOrderFieldsFn: func(c context.Context, orderDetailId int64) error { return nil },
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled++
			return nil
		},
		getCurrentStockFn: func(c context.Context, inputCustID string, inputWhID int64, proID int64) (float64, error) {
			return 0, nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.Update(roNo, entity.UpdateOrderBody{
		CustId:         custID,
		ParentCustId:   parentCustID,
		RoNo:           roNo,
		RoDate:         stringPtrForTest("2026-03-26"),
		ValDate:        stringPtrForTest("2026-03-26"),
		DueDate:        stringPtrForTest("2026-03-26"),
		WhId:           &whID,
		OutletID:       &outletID,
		SalesmanId:     &salesmanID,
		DataStatus:     intPtrForTest(dataStatus),
		DiscValue:      &zero,
		DiscValueFinal: &zero,
		VatValue:       &zero,
		VatValueFinal:  &zero,
		Details:        entity.UpdateOrderDetWithGroup{Normal: []entity.UpdateOrderDetBody{}, Promo: []entity.CreateOrderDetBody{}},
		Rewards:        []entity.CreateOrderRewardBody{},
	}, entity.ValidateResponse{})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if updateHeaderCalled != 1 {
		t.Fatalf("expected header update once, got %d", updateHeaderCalled)
	}
	if findByNotInCalled != 0 {
		t.Fatalf("expected FindDetailByNotInDetailIDs to be skipped for mobile process-only payload, got %d", findByNotInCalled)
	}
	if deleteNotInCalled != 0 {
		t.Fatalf("expected DeleteDetailNotInIDs to be skipped for mobile process-only payload, got %d", deleteNotInCalled)
	}
	if stockWriteCalled != 0 {
		t.Fatalf("expected stock writes to be skipped for mobile process-only payload, got %d", stockWriteCalled)
	}
}

func TestProcessEnhanceWithoutProductEdit_UpdatesProcessedStatusWithoutStockMutation(t *testing.T) {
	custID := "C220010001"
	roNo := "SO2603270001"
	updatedBy := int64(88)

	updateCalled := 0
	stockWriteCalled := 0

	orderRepo := &mockOrderRepository{
		findByNoFn: func(inputRoNo string, inputCustID string) (model.OrderList, error) {
			return model.OrderList{}, nil
		},
		updateFn: func(c context.Context, inputRoNo, inputCustID string, data model.Order) error {
			updateCalled++
			if inputRoNo != roNo {
				t.Fatalf("unexpected ro_no: %s", inputRoNo)
			}
			if inputCustID != custID {
				t.Fatalf("unexpected cust_id: %s", inputCustID)
			}
			if data.DataStatus == nil || *data.DataStatus != int64(entity.PROCESSED) {
				t.Fatalf("expected processed status update, got %+v", data.DataStatus)
			}
			if data.UpdatedBy == nil || *data.UpdatedBy != updatedBy {
				t.Fatalf("expected updated_by %d, got %+v", updatedBy, data.UpdatedBy)
			}
			return nil
		},
	}

	stockRepo := &mockStockRepository{
		salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled++
			return nil
		},
	}

	service := &orderServiceImpl{OrderRepository: orderRepo, StockRepository: stockRepo, Transaction: &mockDbtransaction{}}
	err := service.ProcessEnhanceWithoutProductEdit(context.Background(), roNo, custID, updatedBy)
	if err != nil {
		t.Fatalf("ProcessEnhanceWithoutProductEdit returned error: %v", err)
	}
	if updateCalled != 1 {
		t.Fatalf("expected one header update, got %d", updateCalled)
	}
	if stockWriteCalled != 0 {
		t.Fatalf("expected no stock mutation, got %d", stockWriteCalled)
	}
}

func TestUpdateEnhance_NormalizePurchaseDetailsAlias(t *testing.T) {
	request := entity.EditOrderEnhanceBody{
		PurchaseDetails:    []entity.EditPurchaseOrderDetail{{OrderDetailId: 99}},
		AddPurchaseDetails: []entity.AddPurchaseOrderDetail{{ProId: 77}},
	}

	err := normalizeEnhancePromoFlags(&request)
	if err != nil {
		t.Fatalf("normalizeEnhancePromoFlags returned error: %v", err)
	}
	if len(request.PurchaseOrder) != 1 {
		t.Fatalf("expected purchase_details alias to populate purchase_order, got %d entries", len(request.PurchaseOrder))
	}
	if request.PurchaseOrder[0].OrderDetailId != 99 {
		t.Fatalf("unexpected order detail id after normalization: %d", request.PurchaseOrder[0].OrderDetailId)
	}
	if len(request.AddPurchaseOrder) != 1 {
		t.Fatalf("expected add_purchase_details alias to populate add_purchase_order, got %d entries", len(request.AddPurchaseOrder))
	}
	if request.AddPurchaseOrder[0].ProId != 77 {
		t.Fatalf("unexpected pro_id after add alias normalization: %d", request.AddPurchaseOrder[0].ProId)
	}
}

func TestImportSecondarySales_ReplaceScope_ThreeScenarios(t *testing.T) {
	custID := "C220010001"
	makeOrder := func(date, no string) entity.CreateOrderBody {
		qty := 1.0
		return entity.CreateOrderBody{CustId: custID, RoNo: no, RoDate: &date, Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{ProId: 748, Qty: &qty, Qty1: &qty}}}}
	}
	tests := []struct {
		name   string
		parsed []entity.CreateOrderBody
		want   map[string]int
	}{
		{"one date", []entity.CreateOrderBody{makeOrder("2026-07-11", "A"), makeOrder("2026-07-11", "B")}, map[string]int{"2026-07-11": 2}},
		{"two dates", []entity.CreateOrderBody{makeOrder("2026-07-11", "A"), makeOrder("2026-07-13", "B")}, map[string]int{"2026-07-11": 1, "2026-07-13": 1}},
		{"three dates", []entity.CreateOrderBody{makeOrder("2026-07-11", "A"), makeOrder("2026-07-13", "B"), makeOrder("2026-07-14", "C")}, map[string]int{"2026-07-11": 1, "2026-07-13": 1, "2026-07-14": 1}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			counts := map[string]int{}
			repo := &mockOrderRepository{
				storeFn: func(ctx context.Context, order *model.Order) error {
					counts[order.RoDate.Format("2006-01-02")]++
					return nil
				},
				storeDetailFn: func(context.Context, *model.OrderDetail) error { return nil },
			}
			service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
			if err := service.importSecondarySales(context.Background(), custID, custID, 9, tc.parsed); err != nil {
				t.Fatal(err)
			}
			if fmt.Sprint(counts) != fmt.Sprint(tc.want) {
				t.Fatalf("counts=%v want=%v", counts, tc.want)
			}
		})
	}
}

func TestImportSecondarySales_LeavesOtherDatesIntact(t *testing.T) {
	var locked []time.Time
	repo := &mockOrderRepository{lockOrderByScopeFn: func(_ context.Context, _ string, dates []time.Time) error { locked = dates; return nil }}
	service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
	date := "2026-07-11"
	if err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{{CustId: "cust", RoDate: &date}}); err != nil {
		t.Fatal(err)
	}
	if len(locked) != 1 || locked[0].Format("2006-01-02") != date {
		t.Fatalf("scope=%v", locked)
	}
}

func TestImportSecondarySales_LeavesNonMappingIntact(t *testing.T) {
	calls := 0
	repo := &mockOrderRepository{deleteOrderByScopeFn: func(context.Context, string, []time.Time) (int64, error) { calls++; return 0, nil }}
	service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
	date := "2026-07-11"
	if err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{{CustId: "cust", RoDate: &date}}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("mapping delete calls=%d", calls)
	}
}

func TestImportSecondarySales_AllOrNothing_InsertFails_Rollback(t *testing.T) {
	stored := 0
	repo := &mockOrderRepository{
		storeFn: func(context.Context, *model.Order) error { stored++; return errors.New("insert failed") },
	}
	service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
	date := "2026-07-11"
	err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{{CustId: "cust", RoDate: &date}})
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("err=%v", err)
	}
	if stored != 1 {
		t.Fatalf("stored=%d", stored)
	}
}

func TestParseImportOrders_AppliesSevenDayRule_TooOld(t *testing.T) {
	today := time.Date(2026, 7, 15, 12, 0, 0, 0, jakartaLoc)
	if got := validateImportDate(today.AddDate(0, 0, -8), today); got != "Transaction Date cannot be more than 7 days before the current date." {
		t.Fatal(got)
	}
}

func TestParseImportOrders_AppliesSevenDayRule_Future(t *testing.T) {
	today := time.Date(2026, 7, 15, 12, 0, 0, 0, jakartaLoc)
	if got := validateImportDate(today.AddDate(0, 0, 1), today); got != "Transaction Date cannot be later than the current date." {
		t.Fatal(got)
	}
}

func TestParseImportOrders_SevenDayRule_AtBoundaryTodayMinus7_OK(t *testing.T) {
	today := time.Date(2026, 7, 15, 12, 0, 0, 0, jakartaLoc)
	if got := validateImportDate(today.AddDate(0, 0, -7), today); got != "" {
		t.Fatal(got)
	}
}

func TestImportSecondarySales_InvalidScopeDate(t *testing.T) {
	date := "not-a-date"
	service := &orderServiceImpl{OrderRepository: &mockOrderRepository{}, Transaction: &mockDbtransaction{}}
	if err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{{RoDate: &date}}); err == nil {
		t.Fatal("expected invalid scope date error")
	}
}

func TestImportSecondarySales_NilDateSkipsScope(t *testing.T) {
	repo := &mockOrderRepository{}
	service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
	if err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{{}}); err != nil {
		t.Fatal(err)
	}
}

func TestImportSecondarySales_OptionalDateParseErrors(t *testing.T) {
	date := "2026-07-11"
	badDates := []struct {
		name  string
		order entity.CreateOrderBody
	}{
		{"delivery", entity.CreateOrderBody{RoDate: &date, DeliveryDate: func() *string { v := "bad"; return &v }()}},
		{"due", entity.CreateOrderBody{RoDate: &date, DueDate: func() *string { v := "bad"; return &v }()}},
		{"invoice", entity.CreateOrderBody{RoDate: &date, InvoiceDate: func() *string { v := "bad"; return &v }()}},
	}
	for _, tc := range badDates {
		t.Run(tc.name, func(t *testing.T) {
			service := &orderServiceImpl{OrderRepository: &mockOrderRepository{}, Transaction: &mockDbtransaction{}}
			if err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{tc.order}); err == nil {
				t.Fatal("expected optional date parse error")
			}
		})
	}
}

func TestImportSecondarySales_StoreDetailError(t *testing.T) {
	date := "2026-07-11"
	qty := 1.0
	repo := &mockOrderRepository{storeDetailFn: func(context.Context, *model.OrderDetail) error { return errors.New("detail insert failed") }}
	service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
	order := entity.CreateOrderBody{CustId: "cust", RoNo: "RO-1", RoDate: &date, Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{ProId: 1, Qty: &qty}}}}
	if err := service.importSecondarySales(context.Background(), "cust", "parent", 1, []entity.CreateOrderBody{order}); err == nil || err.Error() != "detail insert failed" {
		t.Fatalf("err=%v", err)
	}
}

func TestImportSecondarySales_MapsOptionalDatesAndFields(t *testing.T) {
	roDate, deliveryDate, dueDate, invoiceDate := "2026-07-11", "2026-07-12", "2026-07-13", "2026-07-14"
	qty := 2.0
	var gotOrder *model.Order
	var gotDetail *model.OrderDetail
	repo := &mockOrderRepository{
		storeFn:       func(_ context.Context, order *model.Order) error { gotOrder = order; return nil },
		storeDetailFn: func(_ context.Context, detail *model.OrderDetail) error { gotDetail = detail; return nil },
	}
	service := &orderServiceImpl{OrderRepository: repo, Transaction: &mockDbtransaction{}}
	order := entity.CreateOrderBody{CustId: "cust", RoNo: "RO-1", RoDate: &roDate, DeliveryDate: &deliveryDate, DueDate: &dueDate, InvoiceDate: &invoiceDate, Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{SeqNo: 3, ProId: 9, Qty: &qty}}}}
	if err := service.importSecondarySales(context.Background(), "cust", "parent", 7, []entity.CreateOrderBody{order}); err != nil {
		t.Fatal(err)
	}
	if gotOrder == nil || gotOrder.RoDate == nil || gotOrder.DeliveryDate == nil || gotOrder.DueDate == nil || gotOrder.InvoiceDate == nil || gotDetail == nil {
		t.Fatalf("mapped order=%+v detail=%+v", gotOrder, gotDetail)
	}
	if gotDetail.Qty != qty || gotDetail.QtyFinal != qty {
		t.Fatalf("detail qty=%v/%v", gotDetail.Qty, gotDetail.QtyFinal)
	}
}
