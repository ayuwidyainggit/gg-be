package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/structs"
	"log"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	RepositoryStockImpl struct {
		*gorm.DB
	}
)

type StockRepository interface {
	FindAllByCustId(dataFilter entity.StockQueryFilter, custId, parentCustId string) ([]model.Stock, int64, int, error)
	Report(dataFilter entity.StockReportQueryFilter) ([]model.StockReport, int64, int, error)
	Store(c context.Context, data *model.Stock) error
	StoreBulk(c context.Context, data []*model.Stock) error
	StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error
	OpnameLookup(dataFilter entity.StockOpnameLookupQueryFilter) ([]model.StockOpnameLookup, int64, int, error)
}

func NewStockRepo(db *gorm.DB) *RepositoryStockImpl {
	return &RepositoryStockImpl{db}
}

func intSliceToString(slice []int) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(strSlice, ", ")
}

func productCustIDs(parentCustID, custID string) []string {
	custIDs := make([]string, 0, 2)

	if parentCustID != "" {
		custIDs = append(custIDs, parentCustID)
	}

	if custID != "" && custID != parentCustID {
		custIDs = append(custIDs, custID)
	}

	return custIDs
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryStockImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryStockImpl) Store(c context.Context, data *model.Stock) error {
	err := repository.model(c).Create(data).Error
	return err
}

func (repository *RepositoryStockImpl) StoreBulk(c context.Context, data []*model.Stock) error {
	err := repository.model(c).Create(data).Error
	return err
}

func (repository *RepositoryStockImpl) FindAllByCustId(dataFilter entity.StockQueryFilter, custId, parentCustId string) ([]model.Stock, int64, int, error) {
	var stocks []model.Stock
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("stock_id")
	queryCount.Where("stock.cust_id=?", custId)

	// var strFilterByWhId string
	// if dataFilter.WhId != nil {
	// 	if *dataFilter.WhId > 0 {
	// 		strFilterByWhId = `AND wh.wh_id = ` + strconv.Itoa(*dataFilter.WhId)
	// 	}
	// }

	query := repository.
		Select(`stock.*`)
		// Joins("LEFT JOIN mst.m_product p ON p.pro_id = m_product_dist.pro_id AND p.cust_id = ?", parentCustId).
		// Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = p.sup_id AND s.cust_id = ?", parentCustId).
		// Joins("LEFT JOIN mst.m_sub_brand1 sbr ON sbr.sbrand1_id = p.sbrand1_id AND sbr.cust_id = ?", parentCustId).
		// Joins("LEFT JOIN mst.m_brand br ON br.brand_id = sbr.brand_id AND br.cust_id = ?", parentCustId)

	query.Where("stock.cust_id = ?", custId)

	// if dataFilter.SupId != nil {
	// 	if *dataFilter.SupId > 0 {
	// 		query.Where("p.sup_id = ?", dataFilter.SupId)
	// 		queryCount.Where("p.sup_id = ?", dataFilter.SupId)
	// 	}
	// }

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("stock.pro_code ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&stocks).Error
	if err != nil {
		return stocks, total, 0, err
	}
	err = queryCount.Model(&stocks).Count(&total).Error
	if err != nil {
		return stocks, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return stocks, total, lastPage, nil

}

func (repository *RepositoryStockImpl) InvoiceSalesStockUpdates(c context.Context, stockUpdates []*entity.InvoiceSalesStockUpdate) error {
	var whStocks []model.WarehouseStock
	var stocks []*model.Stock
	mergedWhStocks := make(map[string]*model.WarehouseStock)
	for _, stockUpdate := range stockUpdates {
		key := fmt.Sprintf("%v-%v", stockUpdate.ProID, stockUpdate.WhID)
		qty := 0 - stockUpdate.QtyOrderBefore
		if existing, ok := mergedWhStocks[key]; ok {
			existing.QtyOnOrder += qty * -1
		} else {
			whStock := &model.WarehouseStock{
				CustID:     stockUpdate.CustID,
				WhID:       stockUpdate.WhID,
				ProID:      stockUpdate.ProID,
				QtyOnOrder: qty * -1,
			}
			mergedWhStocks[key] = whStock
		}

		// gudang
		stock := model.Stock{
			CustID:      stockUpdate.CustID,
			StockDate:   stockUpdate.StockDate,
			TrCode:      stockUpdate.TrCode,
			TrNo:        stockUpdate.TrNo,
			WhID:        stockUpdate.WhID,
			ProID:       stockUpdate.ProID,
			ItemCdn:     1,
			UnitPrice:   stockUpdate.UnitPrice,
			RefDetId:    stockUpdate.RefDetId,
			QtyOutOrder: qty,
		}

		stocks = append(stocks, &stock)
	}

	for _, whStock := range mergedWhStocks {
		whStocks = append(whStocks, *whStock)
	}
	log.Println("mergeWhStocks:", structs.StructToJson(whStocks))

	err := repository.UpsertWithExistingValueArr(c, whStocks)
	if err != nil {
		return err
	}

	err = repository.StoreBulk(c, stocks)
	if err != nil {
		return err
	}

	return nil
}

func (repository *RepositoryStockImpl) SalesStockUpdates(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
	var whStocks []model.WarehouseStock
	var stocks []*model.Stock
	mergedWhStocks := make(map[string]*model.WarehouseStock)

	for _, stockUpdate := range stockUpdates {
		key := fmt.Sprintf("%v-%v", stockUpdate.ProID, stockUpdate.WhID)
		qty := stockUpdate.QtyOrder
		qtyOnStock := stockUpdate.QtyOrder * -1

		if stockUpdate.QtyOrderBefore != nil {
			if *stockUpdate.QtyOrderBefore > stockUpdate.QtyOrder {
				qty = *stockUpdate.QtyOrderBefore - qty
				qtyOnStock = *stockUpdate.QtyOrderBefore - qty
			} else {
				qty = qty - *stockUpdate.QtyOrderBefore               // 10
				qtyOnStock = (qty - *stockUpdate.QtyOrderBefore) * -1 // -10
			}

		}

		if existing, ok := mergedWhStocks[key]; ok {
			existing.Qty += qty * -1
			existing.QtyOnOrder += qty
		} else {
			whStock := &model.WarehouseStock{
				CustID:     stockUpdate.CustID,
				WhID:       stockUpdate.WhID,
				ProID:      stockUpdate.ProID,
				QtyOnOrder: qty,
				Qty:        qtyOnStock,
			}
			mergedWhStocks[key] = whStock
		}

		// gudang
		stockWarehouse := model.Stock{
			CustID:      stockUpdate.CustID,
			StockDate:   stockUpdate.StockDate,
			TrCode:      stockUpdate.TrCode,
			TrNo:        stockUpdate.TrNo,
			WhID:        stockUpdate.WhID,
			ProID:       stockUpdate.ProID,
			ItemCdn:     1,
			UnitPrice:   stockUpdate.UnitPrice,
			RefDetId:    stockUpdate.RefDetId,
			QtyInOrder:  0,
			QtyOutOrder: qty,
		}

		stocks = append(stocks, &stockWarehouse)

		// gudang
		stock := model.Stock{
			CustID:      stockUpdate.CustID,
			StockDate:   stockUpdate.StockDate,
			TrCode:      "CO",
			TrNo:        stockUpdate.TrNo + "-" + "CO",
			WhID:        stockUpdate.WhID,
			ProID:       stockUpdate.ProID,
			ItemCdn:     1,
			UnitPrice:   stockUpdate.UnitPrice,
			RefDetId:    stockUpdate.RefDetId,
			QtyInOrder:  qty,
			QtyOutOrder: 0,
		}

		stocks = append(stocks, &stock)
	}

	for _, whStock := range mergedWhStocks {
		whStocks = append(whStocks, *whStock)
	}
	log.Println("mergeWhStocks:", structs.StructToJson(whStocks))

	err := repository.UpsertWithExistingValueArr(c, whStocks)
	if err != nil {
		return err
	}

	err = repository.StoreBulk(c, stocks)
	if err != nil {
		return err
	}

	return nil
}

func (repository *RepositoryStockImpl) StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error {
	var whStocks []model.WarehouseStock
	var stocks []*model.Stock

	mergedWhStocks := make(map[string]*model.WarehouseStock)

	for index, stockUpdate := range stockUpdates {
		key := fmt.Sprintf("%v-%v", stockUpdate.ProID, stockUpdate.WhID)
		var qty float64

		if stockUpdate.QtyOut > 0 && stockUpdate.QtyIn > 0 {
			return errors.New(fmt.Sprintf("qtyin and qtyout can't greater than 0 on same row for product ID %v at index %v", key, index))
		}

		if stockUpdate.QtyOut > 0 && stockUpdate.QtyIn == 0 {
			qty = stockUpdate.QtyOut * -1
		} else {
			qty = float64(stockUpdate.QtyIn)
		}

		if existing, ok := mergedWhStocks[key]; ok {
			existing.Qty += qty
		} else {
			whStock := &model.WarehouseStock{
				CustID: stockUpdate.CustID,
				WhID:   stockUpdate.WhID,
				ProID:  stockUpdate.ProID,
				Qty:    qty,
			}
			mergedWhStocks[key] = whStock
		}

		stock := model.Stock{
			CustID:    stockUpdate.CustID,
			StockDate: stockUpdate.StockDate,
			TrCode:    stockUpdate.TrCode,
			TrNo:      stockUpdate.TrNo,
			WhID:      stockUpdate.WhID,
			ProID:     stockUpdate.ProID,
			ItemCdn:   1,
			UnitPrice: stockUpdate.UnitPrice,
			RefDetId:  stockUpdate.RefDetId,
			QtyIn:     stockUpdate.QtyIn,
			QtyOut:    stockUpdate.QtyOut,
		}

		stocks = append(stocks, &stock)

	}

	for _, whStock := range mergedWhStocks {
		whStocks = append(whStocks, *whStock)
	}
	log.Println("mergeWhStocks:", structs.StructToJson(whStocks))

	err := repository.UpsertWithExistingValueArr(c, whStocks)
	if err != nil {
		return err
	}

	err = repository.StoreBulk(c, stocks)
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockImpl) UpsertWithExistingValueArr(c context.Context, data []model.WarehouseStock) error {
	err := repository.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "cust_id"}, {Name: "wh_id"}, {Name: "pro_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"qty":          gorm.Expr("inv.warehouse_stock.qty + EXCLUDED.qty"),
				"qty_on_order": gorm.Expr("inv.warehouse_stock.qty_on_order + EXCLUDED.qty_on_order"),
			}),
		},
	).Create(&data).Error
	if err != nil {
		log.Println("UpsertQty, error:", err.Error())
		return err
	}
	return nil
}

func (repository *RepositoryStockImpl) Report(dataFilter entity.StockReportQueryFilter) ([]model.StockReport, int64, int, error) {
	var stockReport []model.StockReport
	var total int64
	// productTenantCustIDs := productCustIDs(dataFilter.ParentCustID, dataFilter.CustID)
	productTenantCustIDs := []string{dataFilter.CustID}

	isShowPrice, _ := strconv.ParseBool(dataFilter.ShowPrice)
	qShowPrice := ""
	if isShowPrice {
		qShowPrice = `pro.purch_price1, pro.purch_price2, pro.purch_price3,
		pro.sell_price1, pro.sell_price2, pro.sell_price3, pro.vat, pro.vat_lg_purch, pro.vat_lg_sell,`
	}

	qDate := ""
	if dataFilter.Date != "" {
		qDate = "AND st.stock_date <= '" + dataFilter.Date + "'"
	}

	qOrderQty := "0 AS order_qty1, 0 AS order_qty2, 0 AS order_qty3"
	qOrderQtyJoin := ""
	orderDateFilter := dataFilter.OrderDate
	if orderDateFilter == "" {
		orderDateFilter = dataFilter.Date
	}
	if dataFilter.OutletID > 0 && orderDateFilter != "" {
		qOrderQty = `COALESCE(MAX(order_qty.order_qty1), 0) AS order_qty1,
		COALESCE(MAX(order_qty.order_qty2), 0) AS order_qty2,
		COALESCE(MAX(order_qty.order_qty3), 0) AS order_qty3`
		qOrderQtyJoin = fmt.Sprintf(`LEFT JOIN (
			SELECT 
				od.pro_id,
				COALESCE(SUM(od.qty1), 0) as order_qty1,
				COALESCE(SUM(od.qty2), 0) as order_qty2,
				COALESCE(SUM(od.qty3), 0) as order_qty3
			FROM sls.order_detail od
			INNER JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id
			WHERE o.cust_id = '%s'
			AND o.outlet_id = %d
			AND o.wh_id IN (%s)
			AND o.ro_date = '%s'
			AND o.data_status IN (1, 2, 3)
			AND od.item_type = 1
			GROUP BY od.pro_id
		) AS order_qty ON order_qty.pro_id = pro.pro_id`,
			dataFilter.CustID,
			dataFilter.OutletID,
			intSliceToString(dataFilter.WhID),
			orderDateFilter)
	}

	subQuery := repository.Select(`pro.pro_id, pro.pro_code, pro.pro_name, 
		pro.unit_id1, pro.unit_id2, pro.unit_id3, pro.conv_unit2, pro.conv_unit3, 
		`+qShowPrice+`
		pro.sup_id, pro.is_active, pro.deleted_at,
		COALESCE(SUM(st.qty_in), 0)-COALESCE(SUM(st.qty_out), 0) AS qty,
		COALESCE(SUM(st.qty_in_order), 0)-COALESCE(SUM(st.qty_out_order), 0) AS qty_order,
		`+qOrderQty).
		Joins("LEFT JOIN inv.stock st ON st.pro_id = pro.pro_id "+qDate+" AND st.cust_id = ? AND st.wh_id IN ?", dataFilter.CustID, dataFilter.WhID).
		Joins("LEFT JOIN mst.m_sub_brand1 msb on pro.sbrand1_id = msb.sbrand1_id AND msb.cust_id = ?", dataFilter.ParentCustID).
		Joins("LEFT JOIN mst.m_brand mb on msb.brand_id = mb.brand_id AND mb.cust_id = ?", dataFilter.ParentCustID)

	if qOrderQtyJoin != "" {
		subQuery = subQuery.Joins(qOrderQtyJoin)
	}

	subQuery = subQuery.Table("mst.m_product AS pro").Where("pro.cust_id IN ?", productTenantCustIDs)

	subQuery = subQuery.Where("pro.is_del IS FALSE")

	if dataFilter.Query != "" {
		subQuery = subQuery.Where("pro.pro_code ILIKE ? or pro.pro_name ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.BrandID) > 0 {
		subQuery = subQuery.Where("msb.brand_id in ?", dataFilter.BrandID)
	}

	if len(dataFilter.PCatID) > 0 {
		subQuery = subQuery.Where("pro.pcat_id in ?", dataFilter.PCatID)
	}

	if len(dataFilter.PLID) > 0 {
		subQuery = subQuery.Where("mb.pl_id in ?", dataFilter.PLID)
	}

	isActiveProductOnly, _ := strconv.ParseBool(dataFilter.ActiveProductOnly)
	if isActiveProductOnly {
		subQuery = subQuery.Where("pro.is_active = true")
	}

	if len(dataFilter.ProID) > 0 {
		subQuery = subQuery.Where("pro.pro_id IN ?", dataFilter.ProID)
	}
	if len(dataFilter.SupID) > 0 {
		subQuery = subQuery.Where("pro.sup_id IN ?", dataFilter.SupID)
	}
	// if len(dataFilter.WhID) > 0 {
	// 	subQuery.Where("st.wh_id IN ?", dataFilter.WhID)
	// }
	subQuery = subQuery.Group("pro.pro_id")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		subQuery = subQuery.Order("pro." + sortBy)
	} else {
		subQuery = subQuery.Order("pro.pro_id ASC")
	}

	// Using subquery in FROM clause
	query := repository.Table("(?) AS stock_report", subQuery)

	isIncludeZeroStock, _ := strconv.ParseBool(dataFilter.IncludeZeroStock)
	if !isIncludeZeroStock {
		query = query.Where("qty > 0 OR qty_order > 0")
	}

	query.Find(&stockReport)

	total = int64(len(stockReport))
	return stockReport, total, 1, nil

}

func (repository *RepositoryStockImpl) OpnameLookup(dataFilter entity.StockOpnameLookupQueryFilter) ([]model.StockOpnameLookup, int64, int, error) {
	var stockOpnameLookup []model.StockOpnameLookup
	var total int64
	productTenantCustIDs := productCustIDs(dataFilter.ParentCustID, dataFilter.CustID)

	qDate := ""
	if dataFilter.Date != "" {
		qDate = "AND st.stock_date <= '" + dataFilter.Date + "'"
	}

	subQuery := repository.Select(`pro.pro_id, pro.pro_code, pro.pro_name,
	COALESCE(SUM(st.qty_in), 0)-COALESCE(SUM(st.qty_out), 0) AS qty`).
		Table("mst.m_product AS pro").
		Joins("LEFT JOIN inv.stock st ON st.pro_id = pro.pro_id "+qDate+" AND st.cust_id = ? AND st.wh_id IN ?", dataFilter.CustID, dataFilter.WhID).
		Where("pro.cust_id IN ?", productTenantCustIDs).
		Group("pro.pro_id")

	// Using subquery in FROM clause
	query := repository.Table("(?) AS stock", subQuery)
	isIncludeZeroStock, _ := strconv.ParseBool(dataFilter.IncludeZeroStock)
	if !isIncludeZeroStock {
		query = query.Where("qty > 0")
	}

	subSubQuery := repository.Table("(?) AS stock_opname_lookup", query)
	if dataFilter.ProductHierarchy == 1 {
		subSubQuery = subSubQuery.Select(`pro_id AS id, pro_code AS code, pro_name AS name`).
			Order("pro_code ASC")
		subSubQuery.Find(&stockOpnameLookup)
	} else {
		subSubQuery = subSubQuery.Joins("LEFT JOIN mst.m_product p ON p.pro_id = stock_opname_lookup.pro_id AND p.cust_id IN ?", productTenantCustIDs)
		if dataFilter.ProductHierarchy == 2 {
			subSubQuery = subSubQuery.Select(`mpc.pcat_id AS id, mpc.pcat_code AS code, mpc.pcat_name AS name`).
				Joins("LEFT JOIN mst.m_product_cat mpc ON mpc.pcat_id = p.pcat_id AND mpc.cust_id = ?", dataFilter.ParentCustID).
				Group("mpc.pcat_id, mpc.pcat_code, mpc.pcat_name").
				Order("mpc.pcat_code ASC")
		} else {
			subSubQuery = subSubQuery.Joins("LEFT JOIN mst.m_sub_brand1 msb ON p.sbrand1_id = msb.sbrand1_id AND msb.cust_id = ?", dataFilter.ParentCustID).
				Joins("LEFT JOIN mst.m_brand mb ON mb.brand_id = msb.brand_id AND mb.cust_id = ?", dataFilter.ParentCustID)
			if dataFilter.ProductHierarchy == 3 {
				subSubQuery = subSubQuery.Select(`mb.brand_id AS id, mb.brand_code AS code, mb.brand_name AS name`).
					Group("mb.brand_id, mb.brand_code, mb.brand_name").
					Order("mb.brand_code ASC")
			} else if dataFilter.ProductHierarchy == 4 {
				subSubQuery = subSubQuery.Select(`mpl.pl_id AS id, mpl.pl_code AS code, mpl.pl_name AS name`).
					Joins("LEFT JOIN mst.m_product_line mpl on mpl.pl_id = mb.pl_id AND mpl.cust_id = ?", dataFilter.ParentCustID).
					Group("mpl.pl_id, mpl.pl_code, mpl.pl_name").
					Order("mpl.pl_code ASC")
			}
		}
	}

	subSubQuery.Find(&stockOpnameLookup)
	total = int64(len(stockOpnameLookup))
	return stockOpnameLookup, total, 1, nil

}
