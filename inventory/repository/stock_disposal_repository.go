package repository

import (
	"context"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"inventory/pkg/str"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	RepositoryStockDisposalImpl struct {
		*gorm.DB
	}
)

type StockDisposalRepository interface {
	Store(ctx context.Context, stockDisposal *model.StockDisposal) error
	CreateDetail(ctx context.Context, detail *model.StockDisposalDetail) (*model.StockDisposalDetail, error)
	FindByID(ctx context.Context, sdID int64, custID string) (*model.StockDisposal, error)
	FindByNumber(ctx context.Context, sdNumber string, custID, parentCustID, warehouseCustID string) (*model.StockDisposalList, error)
	FindDetail(ctx context.Context, sdID int64, custID string) ([]model.StockDisposalDetailList, error)
	FindAllByCustId(ctx context.Context, dataFilter entity.StockDisposalQueryFilter, custId, parentCustId string) ([]model.StockDisposalList, int64, int, error)
	FindProductByID(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error)
	FindWarehouseByID(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error)
	FindSupplierByID(ctx context.Context, supID int64, custID string) (*model.GrSupplier, error)
	GetAvailableStock(ctx context.Context, custID string, whID int64, proID int64) (float64, error)
	FindProductsForLookup(ctx context.Context, dataFilter entity.StockDisposalProductLookupQueryFilter, custId, parentCustId string) ([]model.StockDisposalProductLookup, int64, int, error)
}

func NewStockDisposalRepo(db *gorm.DB) *RepositoryStockDisposalImpl {
	return &RepositoryStockDisposalImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryStockDisposalImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repo *RepositoryStockDisposalImpl) Store(ctx context.Context, stockDisposal *model.StockDisposal) error {
	err := repo.model(ctx).Create(stockDisposal).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RepositoryStockDisposalImpl) CreateDetail(ctx context.Context, detail *model.StockDisposalDetail) (*model.StockDisposalDetail, error) {
	err := repo.model(ctx).Create(detail).Error
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (repo *RepositoryStockDisposalImpl) FindByID(ctx context.Context, sdID int64, custID string) (*model.StockDisposal, error) {
	var stockDisposal model.StockDisposal
	query := repo.model(ctx).
		Where("sd_id = ? AND is_del = false", sdID)

	if len(custID) < 10 {
		query = query.Where("cust_id LIKE ?", custID+"%")
	} else {
		query = query.Where("cust_id = ?", custID)
	}

	err := query.
		First(&stockDisposal).Error
	if err != nil {
		return nil, err
	}
	return &stockDisposal, nil
}

func (repo *RepositoryStockDisposalImpl) FindByNumber(ctx context.Context, sdNumber string, custID, parentCustID, warehouseCustID string) (*model.StockDisposalList, error) {
	var stockDisposal model.StockDisposalList
	warehouseJoinCustID := warehouseCustID
	if warehouseJoinCustID == "" {
		warehouseJoinCustID = custID
	}
	err := repo.model(ctx).
		Select("stock_disposal.*, us.user_fullname AS created_by_name, us2.user_fullname AS updated_by_name, sup.sup_code, sup.sup_name, wh.wh_code, wh.wh_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = stock_disposal.created_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = stock_disposal.updated_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = stock_disposal.sup_id AND sup.cust_id = ?", parentCustID).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = stock_disposal.wh_id AND wh.cust_id = ?", warehouseJoinCustID).
		Where("stock_disposal.sd_number = ? AND stock_disposal.cust_id = ? AND stock_disposal.is_del = false", sdNumber, custID).
		Take(&stockDisposal).Error
	if err != nil {
		return nil, err
	}
	return &stockDisposal, nil
}

func (repo *RepositoryStockDisposalImpl) FindDetail(ctx context.Context, sdID int64, custID string) ([]model.StockDisposalDetailList, error) {
	var details []model.StockDisposalDetailList
	err := repo.model(ctx).
		Select("stock_disposal_detail.*, pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3").
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = stock_disposal_detail.pro_id").
		Where("stock_disposal_detail.sd_id = ? AND stock_disposal_detail.cust_id = ? AND stock_disposal_detail.is_del = false", sdID, custID).
		Order("stock_disposal_detail.sd_detail_id ASC").
		Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (repo *RepositoryStockDisposalImpl) buildListQueryBase(ctx context.Context, dataFilter entity.StockDisposalQueryFilter, custId, parentCustId string, forCount bool) *gorm.DB {
	var query *gorm.DB
	if forCount {
		query = repo.model(ctx).Select("stock_disposal.sd_id")
	} else {
		// Calculate subtotal from sum of gross_price in details (without VAT)
		query = repo.model(ctx).
			Select(`stock_disposal.*, 
				us.user_fullname AS created_by_name, 
				us2.user_fullname AS updated_by_name, 
				sup.sup_code, 
				sup.sup_name, 
				wh.wh_code, 
				wh.wh_name,
				COALESCE((
					SELECT SUM(gross_price) 
					FROM inv.stock_disposal_detail 
					WHERE stock_disposal_detail.sd_id = stock_disposal.sd_id 
					AND stock_disposal_detail.cust_id = stock_disposal.cust_id 
					AND stock_disposal_detail.is_del = false
				), 0) AS calculated_subtotal`)
	}

	query = query.
		Joins("LEFT JOIN sys.m_user us ON us.user_id = stock_disposal.created_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = stock_disposal.updated_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = stock_disposal.sup_id AND sup.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = stock_disposal.wh_id AND wh.cust_id = stock_disposal.cust_id").
		Where("stock_disposal.is_del = false")

	if len(custId) < 10 {
		query = query.Where("stock_disposal.cust_id LIKE ?", custId+"%")
	} else {
		query = query.Where("stock_disposal.cust_id = ?", custId)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query = query.Where("stock_disposal.disposal_date BETWEEN ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.WhID) > 0 {
		query = query.Where("stock_disposal.wh_id IN ?", dataFilter.WhID)
	}

	if len(dataFilter.SupID) > 0 {
		query = query.Where("stock_disposal.sup_id IN ?", dataFilter.SupID)
	}

	if dataFilter.StockType != "" {
		query = query.Where("stock_disposal.stock_type = ?", dataFilter.StockType)
	}

	if dataFilter.Query != "" {
		query = query.Where(`(
			stock_disposal.sd_number ILIKE ? OR
			sup.sup_name ILIKE ? OR
			wh.wh_name ILIKE ?
		)`, "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	return query
}

func (repo *RepositoryStockDisposalImpl) FindAllByCustId(ctx context.Context, dataFilter entity.StockDisposalQueryFilter, custId, parentCustId string) ([]model.StockDisposalList, int64, int, error) {
	var stockDisposals []model.StockDisposalList
	var total int64
	var limit int
	if dataFilter.Limit <= 0 || dataFilter.Limit > 9999 {
		limit = constant.DEFAULT_PAGE_LIMIT
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repo.buildListQueryBase(ctx, dataFilter, custId, parentCustId, true)
	query := repo.buildListQueryBase(ctx, dataFilter, custId, parentCustId, false)

	if dataFilter.Sort != "" {
		var sortByBuilder strings.Builder
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for i, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				if i > 0 {
					sortByBuilder.WriteString(", ")
				}
				sortByBuilder.WriteString(fmt.Sprintf(`stock_disposal.%s %s`, colSort[0], colSort[1]))
			}
		}
		if sortByBuilder.Len() > 0 {
			query.Order(sortByBuilder.String())
		} else {
			query.Order("stock_disposal.created_at DESC")
		}
	} else {
		query.Order("stock_disposal.created_at DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&stockDisposals).Error
	if err != nil {
		return stockDisposals, total, 0, err
	}

	err = queryCount.Model(&stockDisposals).Count(&total).Error
	if err != nil {
		return stockDisposals, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return stockDisposals, total, lastPage, nil
}

func (repo *RepositoryStockDisposalImpl) FindProductByID(ctx context.Context, proID int64, custID, parentCustID string) (*model.Product, error) {
	var product model.Product
	query := repo.model(ctx).
		Where("pro_id = ? AND is_del = false AND is_active = true", proID)

	if custID == parentCustID {
		query = query.Where("cust_id LIKE ?", parentCustID+"%")
	} else {
		tenantCustIDs := productCustIDs(parentCustID, custID)
		if len(tenantCustIDs) > 0 {
			query = query.Where("cust_id IN ?", tenantCustIDs)
		}
	}

	err := query.
		Order(clause.Expr{SQL: "CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END", Vars: []interface{}{parentCustID, custID}}).
		Take(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *RepositoryStockDisposalImpl) FindWarehouseByID(ctx context.Context, whID int64, custID string) (*model.WarehouseStockWhList, error) {
	var warehouse model.WarehouseStockWhList
	query := repo.model(ctx).
		Table("mst.m_warehouse").
		Select("cust_id, wh_id, wh_code, wh_name, stock_type").
		Where("wh_id = ? AND is_active = true AND is_del = false", whID)

	if len(custID) < 10 {
		query = query.Where("cust_id LIKE ?", custID+"%")
	} else {
		query = query.Where("cust_id = ?", custID)
	}

	err := query.
		Take(&warehouse).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (repo *RepositoryStockDisposalImpl) FindSupplierByID(ctx context.Context, supID int64, custID string) (*model.GrSupplier, error) {
	var supplier model.GrSupplier
	err := repo.model(ctx).
		Table("mst.m_supplier").
		Select("sup_id, sup_code, sup_name").
		Where("sup_id = ? AND cust_id = ? AND is_del = false", supID, custID).
		Take(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// GetAvailableStock returns available stock quantity for a product in a warehouse
func (repo *RepositoryStockDisposalImpl) GetAvailableStock(ctx context.Context, custID string, whID int64, proID int64) (float64, error) {
	var availableQty float64
	query := repo.model(ctx).
		Table("inv.warehouse_stock").
		Select("COALESCE(qty, 0)").
		Where("wh_id = ? AND pro_id = ?", whID, proID)

	if len(custID) < 10 {
		query = query.Where("cust_id LIKE ?", custID+"%")
	} else {
		query = query.Where("cust_id = ?", custID)
	}

	err := query.
		Scan(&availableQty).Error
	if err != nil {
		return 0, err
	}
	return availableQty, nil
}

// buildProductLookupQuery builds base query for product lookup with all filters applied
// Uses same logic as FindApprovalProducts for in_transit_stock and qty1/2/3 calculation
func (repo *RepositoryStockDisposalImpl) buildProductLookupQuery(ctx context.Context, dataFilter entity.StockDisposalProductLookupQueryFilter, custId, parentCustId string, forCount bool) *gorm.DB {
	// Determine if user is principal (custId == parentCustId)
	isPrincipal := custId == parentCustId

	// Build subqueries using helper functions
	whsTotalQtySubquery := BuildWarehouseStockTotalSubquery(isPrincipal, custId, dataFilter.WhID)
	inTransitStock1Subquery, inTransitStock2Subquery, inTransitStock3Subquery := BuildInTransitStockSubqueries(isPrincipal, custId)

	// Build qty1, qty2, qty3 calculation expressions using helper function
	qty1Expression, qty2Expression, qty3Expression := BuildQtyCalculationExpressions(whsTotalQtySubquery, "p")

	var query *gorm.DB

	if forCount {
		// For count query, use subquery (same as main query)
		query = repo.model(ctx).
			Table("mst.m_product p").
			Select("COUNT(DISTINCT p.pro_id)")
	} else {
		query = repo.model(ctx).
			Table("mst.m_product p").
			Select(`p.pro_id, 
				p.pro_code, 
				p.pro_name, 
				p.vat, 
				p.conv_unit2, 
				p.conv_unit3, 
				p.unit_id1, 
				p.unit_id2, 
				p.unit_id3, 
				p.purch_price1, 
				p.purch_price2, 
				p.purch_price3,
				COALESCE(p.min_stock_qty, 0) AS min_stock_qty,
				COALESCE(p.saf_stock_qty, 0) AS saf_stock_qty,
				` + inTransitStock1Subquery + ` AS in_transit_stock1,
				` + inTransitStock2Subquery + ` AS in_transit_stock2,
				` + inTransitStock3Subquery + ` AS in_transit_stock3,
				` + whsTotalQtySubquery + ` AS total_qty,
				` + qty3Expression + `,
				` + qty2Expression + `,
				` + qty1Expression)
	}

	// Apply WHERE conditions
	_ = isPrincipal
	query = query.
		Joins("JOIN smc.m_customer mc ON mc.cust_id = p.cust_id").
		Where("p.is_del = false AND p.is_active = true")

	if isPrincipal {
		query = query.Where("p.cust_id LIKE ?", parentCustId+"%")
	} else {
		productTenantCustIDs := productCustIDs(parentCustId, custId)
		if len(productTenantCustIDs) > 0 {
			query = query.Where("p.cust_id IN ?", productTenantCustIDs)
		}
	}

	// Filter by sup_id if provided
	if dataFilter.SupID != nil {
		query = query.Where("p.sup_id = ?", *dataFilter.SupID)
	}

	if dataFilter.DistributorID != nil {
		query = query.Where("mc.distributor_id = ?", *dataFilter.DistributorID)
	}

	// Search query
	if dataFilter.Query != "" {
		query = query.Where("p.pro_code ILIKE ? OR p.pro_name ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if dataFilter.ZeroStock != nil {
		if *dataFilter.ZeroStock {
			query = query.Where(whsTotalQtySubquery + " = 0")
		} else {
			query = query.Where(whsTotalQtySubquery + " > 0")
		}
	}

	return query
}

func (repo *RepositoryStockDisposalImpl) FindProductsForLookup(ctx context.Context, dataFilter entity.StockDisposalProductLookupQueryFilter, custId, parentCustId string) ([]model.StockDisposalProductLookup, int64, int, error) {
	var products []model.StockDisposalProductLookup
	var total int64
	var limit int
	if dataFilter.Limit <= 0 || dataFilter.Limit > 9999 {
		limit = constant.DEFAULT_PAGE_LIMIT
	} else {
		limit = dataFilter.Limit
	}

	// Build queries using helper
	query := repo.buildProductLookupQuery(ctx, dataFilter, custId, parentCustId, false)
	queryCount := repo.buildProductLookupQuery(ctx, dataFilter, custId, parentCustId, true)

	// Pagination
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Execute query
	err := query.Order("p.pro_name ASC").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return products, total, 0, err
	}

	// Count total
	err = queryCount.Scan(&total).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return products, total, lastPage, nil
}
