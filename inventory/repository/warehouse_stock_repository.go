package repository

import (
	"context"
	"database/sql"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryWarehouseStockImpl struct {
		*gorm.DB
	}
)

type WarehouseStockRepository interface {
	FindAllByCustId(dataFilter entity.DistributorStockQueryFilter) ([]model.DistributorStockList, int64, int, error)
	FindAllWarehouse(dataFilter entity.WarehouseStockWhListQueryFilter) ([]model.WarehouseStockWhList, int64, int, error)
	Upsert(c context.Context, data *model.WarehouseStock) error
	UpdateQtyOnly(c context.Context, custID string, whID, proID int64, newQty float64) error
	UpsertWithExistingValue(c context.Context, data *model.WarehouseStock) error
	UpsertBulk(c context.Context, datas []*model.WarehouseStock) error
	ProductList(dataFilter entity.ProductWarehouseListQueryFilter) ([]model.ProductWarehouseList, int64, int, error)
}

func NewWarehouseStockRepo(db *gorm.DB) *RepositoryWarehouseStockImpl {
	return &RepositoryWarehouseStockImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryWarehouseStockImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryWarehouseStockImpl) Upsert(c context.Context, data *model.WarehouseStock) error {
	err := repository.Save(data).Error
	return err
}

func (repository *RepositoryWarehouseStockImpl) UpdateQtyOnly(c context.Context, custID string, whID, proID int64, newQty float64) error {
	now := time.Now().UTC().Unix()
	result := repository.model(c).
		Table("inv.warehouse_stock").
		Where("cust_id = ? AND wh_id = ? AND pro_id = ?", custID, whID, proID).
		Updates(map[string]interface{}{
			"qty":        newQty,
			"updated_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		newRow := &model.WarehouseStock{
			CustID:        custID,
			WhID:          whID,
			ProID:         proID,
			Qty:           newQty,
			QtyOnOrder:    0,
			QtyOnShipping: 0,
			QtyBs:         0,
			QtyExp:        0,
			UpdatedAt:     now,
		}
		return repository.model(c).Create(newRow).Error
	}
	return nil
}

func (repository *RepositoryWarehouseStockImpl) UpsertWithExistingValue(c context.Context, data *model.WarehouseStock) error {
	err := repository.model(c).Exec(
		`INSERT INTO inv.warehouse_stock (
			cust_id, wh_id, pro_id, qty, qty_on_order, qty_on_shipping, qty_bs, qty_exp
		) VALUES (
			@cust_id, @wh_id, @pro_id, @qty, @qty_on_order, @qty_on_shipping, @qty_bs, @qty_exp
		) ON CONFLICT ON CONSTRAINT warehouse_stock_pkey 
		DO UPDATE SET qty = inv.warehouse_stock.qty + EXCLUDED.qty;`,
		sql.Named("cust_id", data.CustID),
		sql.Named("wh_id", data.WhID),
		sql.Named("pro_id", data.ProID),
		sql.Named("qty", data.Qty),
		sql.Named("qty_on_order", data.QtyOnOrder),
		sql.Named("qty_on_shipping", data.QtyOnShipping),
		sql.Named("qty_bs", data.QtyBs),
		sql.Named("qty_exp", data.QtyExp)).Error
	if err != nil {
		log.Println("UpsertQty, error:", err.Error())
		return err
	}
	return nil
}

func (repository *RepositoryWarehouseStockImpl) UpsertBulk(c context.Context, data []*model.WarehouseStock) error {
	err := repository.Save(data).Error
	return err
}

func (repository *RepositoryWarehouseStockImpl) FindAllByCustId(dataFilter entity.DistributorStockQueryFilter) ([]model.DistributorStockList, int64, int, error) {
	var whStocks []model.DistributorStockList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	showPrice, _ := strconv.ParseBool(dataFilter.ShowPrice)
	selectPrice := ``
	if showPrice {
		selectPrice = ` m_product.purch_price1, m_product.purch_price2, m_product.purch_price3, 
						m_product.sell_price1, m_product.sell_price2, m_product.sell_price3, `
	}
	queryCount := repository.Select("m_product.pro_id").
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = m_product.pro_id AND whs.cust_id = ? AND whs.wh_id = ?", dataFilter.CustID, dataFilter.WhID).
		Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = m_product.sup_id AND s.cust_id = m_product.cust_id")
	queryCount.Where("m_product.cust_id=?", dataFilter.CustID)

	query := repository.
		Select(`m_product.pro_id, m_product.pro_code, m_product.pro_name, 
				m_product.unit_id1, m_product.unit_id2, m_product.unit_id3,
				m_product.conv_unit2, m_product.conv_unit3,
				`+selectPrice+`
				s.sup_id, s.sup_code, s.sup_name,
				whs.qty, whs.qty_on_order, whs.qty_on_shipping, whs.qty_bs, whs.qty_exp, whs.updated_at`).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = m_product.pro_id AND whs.cust_id = ? AND whs.wh_id = ?", dataFilter.CustID, dataFilter.WhID).
		Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = m_product.sup_id AND s.cust_id = m_product.cust_id")

	query.Where("m_product.cust_id = ?", dataFilter.CustID)

	if len(dataFilter.ProID) > 0 {
		queryCount.Where("m_product.pro_id IN ?", dataFilter.ProID)
		query.Where("m_product.pro_id IN ?", dataFilter.ProID)
	}

	if len(dataFilter.SupID) > 0 {
		queryCount.Where("s.sup_id IN ?", dataFilter.SupID)
		query.Where("s.sup_id IN ?", dataFilter.SupID)
	}

	zeroStock, _ := strconv.ParseBool(dataFilter.ZeroStock)
	if !zeroStock {
		queryCount.Where("whs.qty > 0")
		query.Where("whs.qty > 0")
	}

	activeProduct, _ := strconv.ParseBool(dataFilter.ActiveProduct)
	if activeProduct {
		queryCount.Where("m_product.is_active = true")
		query.Where("m_product.is_active = true")
	}

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
		query.Order("m_product.pro_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&whStocks).Error
	if err != nil {
		return whStocks, total, 0, err
	}
	err = queryCount.Model(&whStocks).Count(&total).Error
	if err != nil {
		return whStocks, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return whStocks, total, lastPage, nil

}

func (repository *RepositoryWarehouseStockImpl) FindAllWarehouse(dataFilter entity.WarehouseStockWhListQueryFilter) ([]model.WarehouseStockWhList, int64, int, error) {
	var warehouses []model.WarehouseStockWhList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("wh_id").
		Joins("LEFT JOIN mst.m_warehouse w ON w.wh_id = warehouse_stock.wh_id AND w.cust_id = ?", dataFilter.CustID).
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = warehouse_stock.pro_id AND p.cust_id = ?", dataFilter.ParentCustID).
		Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = p.sup_id AND s.cust_id = ?", dataFilter.ParentCustID)
	queryCount.Where("warehouse_stock.cust_id=?", dataFilter.CustID)

	query := repository.
		Select(`warehouse_stock.wh_id, w.wh_code, w.wh_name, w.stock_type`).
		Joins("LEFT JOIN mst.m_warehouse w ON w.wh_id = warehouse_stock.wh_id AND w.cust_id = ?", dataFilter.CustID).
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = warehouse_stock.pro_id AND p.cust_id = ?", dataFilter.ParentCustID).
		Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = p.sup_id AND s.cust_id = ?", dataFilter.ParentCustID)

	query.Where("warehouse_stock.cust_id = ?", dataFilter.CustID)

	if len(dataFilter.ProID) > 0 {
		queryCount.Where("warehouse_stock.pro_id IN ?", dataFilter.ProID)
		query.Where("warehouse_stock.pro_id IN ?", dataFilter.ProID)
	}

	if len(dataFilter.SupID) > 0 {
		queryCount.Where("s.sup_id IN ?", dataFilter.SupID)
		query.Where("s.sup_id IN ?", dataFilter.SupID)
	}

	if len(dataFilter.WhID) > 0 {
		queryCount.Where("warehouse_stock.wh_id IN ?", dataFilter.WhID)
		query.Where("warehouse_stock.wh_id IN ?", dataFilter.WhID)
	}

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
		query.Order("warehouse_stock.wh_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Group("warehouse_stock.wh_id, w.wh_code, w.wh_name, w.stock_type").Find(&warehouses).Error
	if err != nil {
		return warehouses, total, 0, err
	}
	err = queryCount.Model(&warehouses).Distinct("warehouse_stock.wh_id").Count(&total).Error
	if err != nil {
		return warehouses, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return warehouses, total, lastPage, nil

}

func (repository *RepositoryWarehouseStockImpl) ProductList(dataFilter entity.ProductWarehouseListQueryFilter) ([]model.ProductWarehouseList, int64, int, error) {
	var whStocks []model.ProductWarehouseList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	timeNowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
	transDate := currentDate.Format("2006-01-02")

	queryCount := repository.Select("p.pro_id").
		Joins("left join mst.m_product p on inv.warehouse_stock.pro_id = p.pro_id")

	queryCount.Where("inv.warehouse_stock.cust_id =?", dataFilter.CustID)

	query := repository.Select(`p.pro_id, p.pro_code, p.pro_name, inv.warehouse_stock.qty, p.conv_unit2, p.conv_unit3, p.unit_id1, p.unit_id2, p.unit_id3, p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell,
		CASE WHEN (mtp.sell_price1=0 or mtp.sell_price1 is null) THEN (CASE WHEN (mtp_mg_pr_sell.sell_price1=0 or mtp_mg_pr_sell.sell_price1 is null) THEN p.sell_price1 ELSE mtp_mg_pr_sell.sell_price1 END) ELSE mtp.sell_price1 END AS sell_price1,
		CASE WHEN (mtp.sell_price2=0 or mtp.sell_price2 is null) THEN (CASE WHEN (mtp_mg_pr_sell.sell_price2=0 or mtp_mg_pr_sell.sell_price2 is null) THEN p.sell_price2 ELSE mtp_mg_pr_sell.sell_price2 END) ELSE mtp.sell_price2 END AS sell_price2,
		CASE WHEN (mtp.sell_price3=0 or mtp.sell_price3 is null) THEN (CASE WHEN (mtp_mg_pr_sell.sell_price3=0 or mtp_mg_pr_sell.sell_price3 is null) THEN p.sell_price3 ELSE mtp_mg_pr_sell.sell_price3 END) ELSE mtp.sell_price3 END AS sell_price3,
		CASE WHEN mtp_mg_pr.purch_price1 is NULL THEN p.purch_price1 ELSE mtp_mg_pr.purch_price1 END AS purch_price1,
		CASE WHEN mtp_mg_pr.purch_price2 is NULL THEN p.purch_price2 ELSE mtp_mg_pr.purch_price2 END AS purch_price2,
		CASE WHEN mtp_mg_pr.purch_price3 is NULL THEN p.purch_price3 ELSE mtp_mg_pr.purch_price3 END AS purch_price3`).
		Joins("LEFT JOIN mst.m_distributor dist ON dist.cust_id = ? AND dist.distributor_id = ?", dataFilter.ParentCustID, dataFilter.DistributorID).
		Joins("left join mst.m_product p on inv.warehouse_stock.pro_id = p.pro_id").
		Joins("LEFT JOIN smc.m_customer cus ON cus.cust_id = ?", dataFilter.CustID).
		Joins("LEFT JOIN mst.m_transaction_price mtp ON mtp.pro_id = p.pro_id AND mtp.cust_id = ? AND mtp.outlet_id = ? AND (? BETWEEN mtp.start_date AND mtp.end_date) ", dataFilter.CustID, 0, transDate).
		Joins("LEFT JOIN LATERAL (SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3 FROM mst.m_transaction_price mtp_mg_pr WHERE mtp_mg_pr.cust_id = ? AND mtp_mg_pr.pro_id = p.pro_id AND mtp_mg_pr.start_date <= ? AND (mtp_mg_pr.distributor_id = (CASE WHEN mtp_mg_pr.coverage='N' THEN 0 ELSE dist.distributor_id END) OR mtp_mg_pr.price_group_reff = dist.dist_price_grp_id)ORDER BY mtp_mg_pr.start_date DESC LIMIT 1) mtp_mg_pr ON true", dataFilter.ParentCustID, transDate).
		Joins("LEFT JOIN LATERAL (SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3 FROM mst.m_transaction_price mtp_mg_pr_sell WHERE mtp_mg_pr_sell.cust_id = ? AND mtp_mg_pr_sell.pro_id = p.pro_id and mtp_mg_pr_sell.source = 10 AND mtp_mg_pr_sell.start_date <= ? AND (mtp_mg_pr_sell.distributor_id = (CASE WHEN mtp_mg_pr_sell.coverage = 'N' THEN 0 ELSE dist.distributor_id END) OR mtp_mg_pr_sell.price_group_reff = dist.dist_price_grp_id) ORDER BY mtp_mg_pr_sell.start_date DESC LIMIT 1) mtp_mg_pr_sell ON true", dataFilter.ParentCustID, transDate)

	if dataFilter.WhID != 0 {
		query.Where(" inv.warehouse_stock.wh_id = ?", dataFilter.WhID)
		queryCount.Where(" inv.warehouse_stock.wh_id = ?", dataFilter.WhID)

	}

	if dataFilter.ProID != 0 {
		query.Where(" inv.warehouse_stock.pro_id = ? ", dataFilter.ProID)
		queryCount.Where(" inv.warehouse_stock.pro_id = ? ", dataFilter.ProID)
	}

	if dataFilter.Query != "" {
		query.Where(" p.pro_name ILIKE ? OR p.pro_code ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		queryCount.Where(" p.pro_name ILIKE ? OR p.pro_code ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Order("p.pro_name ASC").Limit(limit).Offset(offset).Find(&whStocks).Error
	if err != nil {
		return whStocks, total, 0, err
	}

	err = queryCount.Model(&whStocks).Count(&total).Error
	if err != nil {
		return whStocks, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))

	return whStocks, total, lastPage, nil

}
