package repository

import (
	"context"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryStockImpl struct {
		*gorm.DB
	}
)
type StockRepository interface {
	FindByEmpId(EmpId int64, custId string) (stockWarehouse model.StockWareHouseList, err error)
	FindByEmpIdCanvas(EmpId int64, custId string) (stockWarehouse model.StockWareHouseList, err error)
	FindDetailProductGudangUtama(dataFilter entity.StockQueryFilter, WhId int64) (details []model.DetilProductStock, err error)
	FindDetailProductGudangCanvas(dataFilter entity.StockQueryFilter, WhId int64) (details []model.DetilProductStock, err error)
}

func NewStockRepository(db *gorm.DB) *RepositoryStockImpl {
	return &RepositoryStockImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryStockImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryStockImpl) FindByEmpId(EmpId int64, custId string) (stockWarehouse model.StockWareHouseList, err error) {
	err = repository.
		Select(`mst.m_salesman.cust_id, mst.m_salesman.emp_id, mst.m_salesman.sales_name, mw.wh_id, mw.wh_code, mw.wh_name`).
		Joins("left join mst.m_warehouse mw on mw.wh_id = mst.m_salesman.wh_id and mw.cust_id = ?", custId).
		Where("mst.m_salesman.emp_id = ? and mst.m_salesman.is_active = true  AND mst.m_salesman.cust_id=?", EmpId, custId).
		Take(&stockWarehouse).Error
	return stockWarehouse, err
}

func (repository *RepositoryStockImpl) FindByEmpIdCanvas(EmpId int64, custId string) (stockWarehouse model.StockWareHouseList, err error) {
	err = repository.
		Select(`mst.m_salesman.cust_id, mst.m_salesman.emp_id, mst.m_salesman.sales_name, msc.wh_id, mw.wh_code, mw.wh_name`).
		Joins("left join mst.m_salesman_canvas msc on msc.emp_id = mst.m_salesman.emp_id").
		Joins("left join mst.m_warehouse mw on mw.wh_id = msc.wh_id and mw.cust_id = ?", custId).
		Where("mst.m_salesman.emp_id = ? and mst.m_salesman.is_active = true  AND mst.m_salesman.cust_id=?", EmpId, custId).
		Take(&stockWarehouse).Error
	return stockWarehouse, err
}

func (repository *RepositoryStockImpl) FindDetailProductGudangUtama(dataFilter entity.StockQueryFilter, WhId int64) (details []model.DetilProductStock, err error) {

	query := repository.Select(`st.wh_id, pro.pro_id, pro.pro_code, pro.pro_name, 
    pro.unit_id1, pro.unit_id2, pro.unit_id3, pro.conv_unit2, pro.conv_unit3, 
    pro.purch_price1, pro.purch_price2, pro.purch_price3,
    pro.sell_price1, pro.sell_price2, pro.sell_price3, pro.vat, pro.vat_lg_purch, pro.vat_lg_sell,
    pro.sup_id, pro.is_active, pro.deleted_at,
    COALESCE(SUM(st.qty_in), 0)-COALESCE(SUM(st.qty_out), 0) AS qty,
    COALESCE(SUM(st.qty_in_order), 0)-COALESCE(SUM(st.qty_out_order), 0) AS qty_order `).
		Joins("LEFT JOIN inv.stock st ON st.pro_id = pro.pro_id ").
		Joins("LEFT JOIN mst.m_sub_brand1 msb on pro.sbrand1_id = msb.sbrand1_id").
		Joins("LEFT JOIN mst.m_brand mb on msb.brand_id = mb.brand_id").
		Where("st.wh_id = ?", WhId).
		Group("st.wh_id, pro.pro_id, pro.pro_code, pro.pro_name, pro.unit_id1, pro.unit_id2, pro.unit_id3, pro.conv_unit2, pro.conv_unit3, pro.purch_price1, pro.purch_price2, pro.purch_price3,pro.sell_price1, pro.sell_price2, pro.sell_price3, pro.vat, pro.vat_lg_purch, pro.vat_lg_sell, pro.sup_id, pro.is_active, pro.deleted_at").Table("mst.m_product AS pro")

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
		query.Order("pro.pro_name ASC")
	}

	if dataFilter.Query != "" {
		query.Where("pro.pro_code ILIKE ? OR pro.pro_name ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	err = query.Find(&details).Error
	if err != nil {
		return details, err
	}

	return details, err
}

func (repository *RepositoryStockImpl) FindDetailProductGudangCanvas(dataFilter entity.StockQueryFilter, WhId int64) (details []model.DetilProductStock, err error) {

	var products []model.DetilProductStock

	var queryFilter string
	if dataFilter.Query != "" {
		queryFilter = fmt.Sprintf("and LOWER(pro.pro_name) like '%%%s%%'", dataFilter.Query)
	}

	innerQuery := fmt.Sprintf(`
		SELECT 
			st.wh_id, pro.pro_id, pro.pro_code, pro.pro_name, 
			pro.unit_id1, pro.unit_id2, pro.unit_id3, pro.conv_unit2, pro.conv_unit3, 
			pro.purch_price1, pro.purch_price2, pro.purch_price3,
			pro.sell_price1, pro.sell_price2, pro.sell_price3, pro.vat, pro.vat_lg_purch, pro.vat_lg_sell,
			pro.sup_id, pro.is_active, pro.deleted_at,
			COALESCE (SUM(DISTINCT(st.qty_in)), 0) AS qty_stock,
			COALESCE(SUM(st.qty_in_order), 0)-COALESCE(SUM(st.qty_out_order), 0) AS qty_order,
			COALESCE(SUM(st.qty_in), 0)-COALESCE(SUM(st.qty_out), 0) AS qty 
		FROM mst.m_product AS pro 
		LEFT JOIN inv.stock st ON st.pro_id = pro.pro_id  
		WHERE st.wh_id = %d %s
		GROUP BY 
			st.wh_id, pro.pro_id, pro.pro_code, pro.pro_name, pro.unit_id1, pro.unit_id2, pro.unit_id3, pro.conv_unit2, pro.conv_unit3, pro.purch_price1, pro.purch_price2, 
			pro.purch_price3,pro.sell_price1, pro.sell_price2, pro.sell_price3, pro.vat, pro.vat_lg_purch, pro.vat_lg_sell, pro.sup_id, pro.is_active, pro.deleted_at 
		ORDER BY pro.pro_name ASC
	`, WhId, queryFilter)

	finalQuery := fmt.Sprintf(`SELECT 
		aa.*, sum(odd.qty1+odd.qty2+odd.qty3) as qty_order_old, SUM(odd.qty1_final) as qty_order1, SUM(odd.qty2_final) as qty_order2, SUM(odd.qty3_final) as qty_order3 
		FROM (%s) as aa  
		LEFT JOIN sls.order od on od.wh_id = aa.wh_id
		LEFT JOIN sls.order_detail odd on odd.ro_no = od.ro_no and odd.pro_id = aa.pro_id and (odd.pro_id = aa.pro_id and odd.pro_id = aa.pro_id) 
		GROUP BY aa.wh_id,aa.pro_id, aa.pro_code, aa.pro_name, aa.unit_id1, aa.unit_id2, aa.unit_id3, aa.conv_unit2, aa.conv_unit3, aa.purch_price1, aa.purch_price2, 
		aa.purch_price3, aa.sell_price1, aa.sell_price2, aa.sell_price3, aa.vat, aa.vat_lg_purch, aa.vat_lg_sell, aa.sup_id, aa.is_active, aa.deleted_at, aa.qty_stock, aa.qty_order, aa.qty`, innerQuery)

	err = repository.Raw(finalQuery).Scan(&products).Error
	if err != nil {
		return products, err
	}

	// Count total for pagination
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM (%s) as aa LEFT JOIN sls.order od on od.wh_id = aa.wh_id
		LEFT JOIN sls.order_detail odd on odd.ro_no = od.ro_no and odd.pro_id = aa.pro_id and (odd.pro_id = aa.pro_id and odd.pro_id = aa.pro_id) 
		GROUP BY aa.wh_id,aa.pro_id, aa.pro_code, aa.pro_name, aa.unit_id1, aa.unit_id2, aa.unit_id3, aa.conv_unit2, aa.conv_unit3, aa.purch_price1, aa.purch_price2, 
		aa.purch_price3, aa.sell_price1, aa.sell_price2, aa.sell_price3, aa.vat, aa.vat_lg_purch, aa.vat_lg_sell, aa.sup_id, aa.is_active, aa.deleted_at, aa.qty_stock, aa.qty`, innerQuery)
	var count int64
	err = repository.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return products, err
	}

	// lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return products, nil
}
