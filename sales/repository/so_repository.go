package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositorySoImpl struct {
		*gorm.DB
	}
)
type SoRepository interface {
	Store(c context.Context, data *model.So) error
	StoreDetail(c context.Context, data *model.SoDet) error
	FindByNo(SoNo string, custId string) (whAdj model.SoList, err error)
	FindDetail(SoNo string, custId string) (Details []model.SoDetRead, err error)
	FindAllByCustId(dataFilter entity.SoQueryFilter) ([]model.SoList, int64, int, error)
	Delete(c context.Context, custId string, SoNo string, deletedBy int64) error
	Update(c context.Context, SoNo string, data model.So) error
	DeleteDetailNotInIDs(c context.Context, SoNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.SoDet) error
	FindDownloadDataPo(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error)
	FindDownloadDataSo(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error)
	FindDownloadDataFinal(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error)
	FindDownloadQtySummary(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error)
}

func NewSoRepo(db *gorm.DB) *RepositorySoImpl {
	return &RepositorySoImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositorySoImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositorySoImpl) Store(c context.Context, data *model.So) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositorySoImpl) StoreDetail(c context.Context, data *model.SoDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositorySoImpl) FindByNo(SoNo string, custId string) (whAdj model.SoList, err error) {
	err = repository.
		Select("so.*, us.user_fullname AS updated_by_name,ot.outlet_code, ot.outlet_name").
		Joins("left join sys.m_user us on us.user_id = sls.so.updated_by").
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.so.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.so.cust_id=? AND sls.so.so_no = ?", custId, SoNo).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositorySoImpl) FindDetail(SoNo string, custId string) (Details []model.SoDetRead, err error) {
	err = repository.Select("sls.so_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = sls.so_det.pro_id").
		Where("sls.so_det.cust_id=? AND so_no = ?", custId, SoNo).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositorySoImpl) FindAllByCustId(dataFilter entity.SoQueryFilter) ([]model.SoList, int64, int, error) {
	var so []model.SoList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("so_no")
	query := repository.Select("so.*, us.user_fullname AS updated_by_name,ot.outlet_code, ot.outlet_name").
		Joins("left join sys.m_user us on us.user_id = sls.so.updated_by").
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.so.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.so.cust_id=?", dataFilter.CustId)
	query.Where("sls.so.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.so.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.so.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.so.so_no=?", dataFilter.Query)
		query.Where("sls.so.so_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.so.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.so.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.so.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.so.outlet_id in ?", dataFilter.OutletID)
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
		query.Order("so_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&so).Error
	if err != nil {
		return so, total, 0, err
	}
	err = queryCount.Model(&so).Count(&total).Error
	if err != nil {
		return so, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return so, total, lastPage, nil
}

func (repository *RepositorySoImpl) Delete(c context.Context, custId string, SoNo string, deletedBy int64) error {
	var data model.So
	result := repository.model(c).Model(&data).Where("cust_id = ? AND so_no=? AND is_del= ? ", custId, SoNo, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositorySoImpl) Update(c context.Context, SoNo string, data model.So) error {
	result := repository.model(c).Model(&data).Where("so_no=?", SoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositorySoImpl) DeleteDetailNotInIDs(c context.Context, SoNo string, IDs []int64) error {
	var Details model.SoDet
	err := repository.model(c).Where("so_no=? AND so_det_id not in (?) ", SoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositorySoImpl) UpdateDetail(c context.Context, Details *model.SoDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositorySoImpl) FindDownloadDataPo(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadPo, error) {
	var results []model.SoDownloadPo

	query := repository.Select(`
		sls.order.order_no,
		sls.order.po_no,
		sls.order.ro_no as so_no,
		sls.order.ro_date,
		sls.order.invoice_date,
		sls.order.invoice_no,
		mst.m_outlet.outlet_code,
		mst.m_outlet.outlet_name,
		mst.m_employee.emp_code as salesman_code,
		mst.m_salesman.sales_name as employee_name,
		COALESCE(product_cust.sup_code, product_cust_parent_supplier.sup_code, product_parent.sup_code, mst.m_supplier.sup_code, '') as supplier_code,
		COALESCE(product_cust.sup_name, product_cust_parent_supplier.sup_name, product_parent.sup_name, mst.m_supplier.sup_name, '') as supplier_name,
		COALESCE(product_cust.pro_code, product_parent.pro_code, mst.m_product.pro_code, '') as product_code,
		COALESCE(product_cust.pro_name, product_parent.pro_name, mst.m_product.pro_name, '') as product_name,
		sls.order_detail.unit_id3,
		sls.order_detail.unit_id2,
		sls.order_detail.unit_id1,
		sls.order_detail.sell_price_system3,
		sls.order_detail.sell_price_system2,
		sls.order_detail.sell_price_system1,
		sls.order_detail.sell_price_po3,
		sls.order_detail.sell_price_po2,
		sls.order_detail.sell_price_po1,
		sls.order_detail.qty_po3,
		sls.order_detail.qty_po2,
		sls.order_detail.qty_po1,
		sls.order_detail.vat_value_final,
		sls.order_detail.disc_value_final,
		sls.order_detail.vat
	`).
		Joins("INNER JOIN sls.order ON sls.order.ro_no = sls.order_detail.ro_no AND sls.order.cust_id = sls.order_detail.cust_id").
		Joins("LEFT JOIN mst.m_outlet ON mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_salesman ON mst.m_salesman.emp_id = sls.order.salesman_id AND mst.m_salesman.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_employee ON mst.m_employee.emp_id = mst.m_salesman.emp_id AND mst.m_employee.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_product ON mst.m_product.pro_id = sls.order_detail.pro_id AND mst.m_product.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN mst.m_supplier ON mst.m_supplier.sup_id = mst.m_product.sup_id AND mst.m_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, mp.sup_id, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.pro_id = sls.order_detail.pro_id AND mp.cust_id = ? LIMIT 1) AS product_cust ON TRUE", filter.CustId).
		Joins("LEFT JOIN mst.m_supplier AS product_cust_parent_supplier ON product_cust_parent_supplier.sup_id = product_cust.sup_id AND product_cust_parent_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.cust_id = ? AND product_cust.pro_code IS NOT NULL AND mp.pro_code = product_cust.pro_code LIMIT 1) AS product_parent ON TRUE", filter.ParentCustId).
		Where("sls.order_detail.cust_id = ?", filter.CustId).
		Where("sls.order_detail.item_type = 1")

	if filter.StartDate > 0 {
		startDate := str.UnixTimestampToUtcTime(filter.StartDate)
		query = query.Where("sls.order.ro_date >= ?", startDate)
	}

	if filter.EndDate > 0 {
		endDate := str.UnixTimestampToUtcTime(filter.EndDate)
		query = query.Where("sls.order.ro_date <= ?", endDate)
	}

	if len(filter.SalesmanId) > 0 {
		query = query.Where("sls.order.salesman_id IN ?", filter.SalesmanId)
	}

	err := query.Find(&results).Error
	return results, err
}

func (repository *RepositorySoImpl) FindDownloadDataSo(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadSo, error) {
	var results []model.SoDownloadSo

	query := repository.Select(`
		sls.order.order_no,
		sls.order.po_no,
		sls.order.ro_no as so_no,
		sls.order.ro_date,
		sls.order.invoice_date,
		sls.order.invoice_no,
		mst.m_outlet.outlet_code,
		mst.m_outlet.outlet_name,
	mst.m_employee.emp_code as salesman_code,
		mst.m_salesman.sales_name as employee_name,
		COALESCE(product_cust.sup_code, product_cust_parent_supplier.sup_code, product_parent.sup_code, mst.m_supplier.sup_code, '') as supplier_code,
		COALESCE(product_cust.sup_name, product_cust_parent_supplier.sup_name, product_parent.sup_name, mst.m_supplier.sup_name, '') as supplier_name,
		COALESCE(product_cust.pro_code, product_parent.pro_code, mst.m_product.pro_code, '') as product_code,
		COALESCE(product_cust.pro_name, product_parent.pro_name, mst.m_product.pro_name, '') as product_name,
		sls.order_detail.unit_id3,
		sls.order_detail.unit_id2,
		sls.order_detail.unit_id1,
		sls.order_detail.sell_price_system1,
		sls.order_detail.sell_price_system2,
		sls.order_detail.sell_price_system3,
		sls.order_detail.sell_price3,
		sls.order_detail.sell_price2,
		sls.order_detail.sell_price1,
		sls.order_detail.qty3,
		sls.order_detail.qty2,
		sls.order_detail.qty1,
		sls.order_detail.vat_value_final,
		sls.order_detail.disc_value_final,
		sls.order_detail.vat
	`).
		Joins("INNER JOIN sls.order ON sls.order.ro_no = sls.order_detail.ro_no AND sls.order.cust_id = sls.order_detail.cust_id").
		Joins("LEFT JOIN mst.m_outlet ON mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_salesman ON mst.m_salesman.emp_id = sls.order.salesman_id AND mst.m_salesman.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_employee ON mst.m_employee.emp_id = mst.m_salesman.emp_id AND mst.m_employee.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_product ON mst.m_product.pro_id = sls.order_detail.pro_id AND mst.m_product.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN mst.m_supplier ON mst.m_supplier.sup_id = mst.m_product.sup_id AND mst.m_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, mp.sup_id, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.pro_id = sls.order_detail.pro_id AND mp.cust_id = ? LIMIT 1) AS product_cust ON TRUE", filter.CustId).
		Joins("LEFT JOIN mst.m_supplier AS product_cust_parent_supplier ON product_cust_parent_supplier.sup_id = product_cust.sup_id AND product_cust_parent_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.cust_id = ? AND product_cust.pro_code IS NOT NULL AND mp.pro_code = product_cust.pro_code LIMIT 1) AS product_parent ON TRUE", filter.ParentCustId).
		Where("sls.order_detail.cust_id = ?", filter.CustId).
		Where("sls.order_detail.item_type = 1")

	if filter.StartDate > 0 {
		startDate := str.UnixTimestampToUtcTime(filter.StartDate)
		query = query.Where("sls.order.ro_date >= ?", startDate)
	}

	if filter.EndDate > 0 {
		endDate := str.UnixTimestampToUtcTime(filter.EndDate)
		query = query.Where("sls.order.ro_date <= ?", endDate)
	}

	if len(filter.SalesmanId) > 0 {
		query = query.Where("sls.order.salesman_id IN ?", filter.SalesmanId)
	}

	err := query.Find(&results).Error
	return results, err
}

func (repository *RepositorySoImpl) FindDownloadDataFinal(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadFinal, error) {
	var results []model.SoDownloadFinal

	query := repository.Select(`
		sls.order.order_no,
		sls.order.po_no,
		sls.order.ro_no as so_no,
		sls.order.ro_date,
		sls.order.invoice_date,
		sls.order.invoice_no,
		mst.m_outlet.outlet_code,
		mst.m_outlet.outlet_name,
		mst.m_employee.emp_code as salesman_code,
		mst.m_salesman.sales_name as employee_name,
		COALESCE(product_cust.sup_code, product_cust_parent_supplier.sup_code, product_parent.sup_code, mst.m_supplier.sup_code, '') as supplier_code,
		COALESCE(product_cust.sup_name, product_cust_parent_supplier.sup_name, product_parent.sup_name, mst.m_supplier.sup_name, '') as supplier_name,
		COALESCE(product_cust.pro_code, product_parent.pro_code, mst.m_product.pro_code, '') as product_code,
		COALESCE(product_cust.pro_name, product_parent.pro_name, mst.m_product.pro_name, '') as product_name,
		sls.order_detail.unit_id3,
		sls.order_detail.unit_id2,
		sls.order_detail.unit_id1,
		sls.order_detail.sell_price_system1,
		sls.order_detail.sell_price_system2,
		sls.order_detail.sell_price_system3,
		sls.order_detail.sell_price_final3,
		sls.order_detail.sell_price_final2,
		sls.order_detail.sell_price_final1,
		sls.order_detail.qty3_final,
		sls.order_detail.qty2_final,
		sls.order_detail.qty1_final,
		sls.order_detail.vat_value_final,
		sls.order_detail.disc_value_final,
		sls.order_detail.vat
	`).
		Joins("INNER JOIN sls.order ON sls.order.ro_no = sls.order_detail.ro_no AND sls.order.cust_id = sls.order_detail.cust_id").
		Joins("LEFT JOIN mst.m_outlet ON mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_salesman ON mst.m_salesman.emp_id = sls.order.salesman_id AND mst.m_salesman.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_employee ON mst.m_employee.emp_id = mst.m_salesman.emp_id AND mst.m_employee.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_product ON mst.m_product.pro_id = sls.order_detail.pro_id AND mst.m_product.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN mst.m_supplier ON mst.m_supplier.sup_id = mst.m_product.sup_id AND mst.m_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, mp.sup_id, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.pro_id = sls.order_detail.pro_id AND mp.cust_id = ? LIMIT 1) AS product_cust ON TRUE", filter.CustId).
		Joins("LEFT JOIN mst.m_supplier AS product_cust_parent_supplier ON product_cust_parent_supplier.sup_id = product_cust.sup_id AND product_cust_parent_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.cust_id = ? AND product_cust.pro_code IS NOT NULL AND mp.pro_code = product_cust.pro_code LIMIT 1) AS product_parent ON TRUE", filter.ParentCustId).
		Where("sls.order_detail.cust_id = ?", filter.CustId).
		Where("sls.order_detail.item_type = 1")

	if filter.StartDate > 0 {
		startDate := str.UnixTimestampToUtcTime(filter.StartDate)
		query = query.Where("sls.order.ro_date >= ?", startDate)
	}

	if filter.EndDate > 0 {
		endDate := str.UnixTimestampToUtcTime(filter.EndDate)
		query = query.Where("sls.order.ro_date <= ?", endDate)
	}

	if len(filter.SalesmanId) > 0 {
		query = query.Where("sls.order.salesman_id IN ?", filter.SalesmanId)
	}

	err := query.Find(&results).Error
	return results, err
}

func (repository *RepositorySoImpl) FindDownloadQtySummary(filter entity.SoDownloadQueryFilter) ([]model.SoDownloadQtySummary, error) {
	var results []model.SoDownloadQtySummary

	query := repository.Select(`
		sls.order.order_no,
		sls.order.po_no,
		sls.order.ro_no as so_no,
		sls.order.ro_date,
		sls.order.invoice_date,
		sls.order.invoice_no,
		mst.m_outlet.outlet_code,
		mst.m_outlet.outlet_name,
		mst.m_employee.emp_code as salesman_code,
		mst.m_salesman.sales_name as employee_name,
		COALESCE(product_cust.sup_code, product_cust_parent_supplier.sup_code, product_parent.sup_code, mst.m_supplier.sup_code, '') as supplier_code,
		COALESCE(product_cust.sup_name, product_cust_parent_supplier.sup_name, product_parent.sup_name, mst.m_supplier.sup_name, '') as supplier_name,
		COALESCE(product_cust.pro_code, product_parent.pro_code, mst.m_product.pro_code, '') as product_code,
		COALESCE(product_cust.pro_name, product_parent.pro_name, mst.m_product.pro_name, '') as product_name,
		sls.order_detail.unit_id3,
		sls.order_detail.unit_id2,
		sls.order_detail.unit_id1,
		sls.order_detail.qty_po3,
		sls.order_detail.qty_po2,
		sls.order_detail.qty_po1,
		sls.order_detail.qty3,
		sls.order_detail.qty2,
		sls.order_detail.qty1,
		sls.order_detail.qty3_final,
		sls.order_detail.qty2_final,
		sls.order_detail.qty1_final
	`).
		Joins("INNER JOIN sls.order ON sls.order.ro_no = sls.order_detail.ro_no AND sls.order.cust_id = sls.order_detail.cust_id").
		Joins("LEFT JOIN mst.m_outlet ON mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_salesman ON mst.m_salesman.emp_id = sls.order.salesman_id AND mst.m_salesman.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_employee ON mst.m_employee.emp_id = mst.m_salesman.emp_id AND mst.m_employee.cust_id = ?", filter.CustId).
		Joins("LEFT JOIN mst.m_product ON mst.m_product.pro_id = sls.order_detail.pro_id AND mst.m_product.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN mst.m_supplier ON mst.m_supplier.sup_id = mst.m_product.sup_id AND mst.m_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, mp.sup_id, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.pro_id = sls.order_detail.pro_id AND mp.cust_id = ? LIMIT 1) AS product_cust ON TRUE", filter.CustId).
		Joins("LEFT JOIN mst.m_supplier AS product_cust_parent_supplier ON product_cust_parent_supplier.sup_id = product_cust.sup_id AND product_cust_parent_supplier.cust_id = ?", filter.ParentCustId).
		Joins("LEFT JOIN LATERAL (SELECT mp.pro_code, mp.pro_name, ms.sup_code, ms.sup_name FROM mst.m_product mp LEFT JOIN mst.m_supplier ms ON ms.sup_id = mp.sup_id AND ms.cust_id = mp.cust_id WHERE mp.cust_id = ? AND product_cust.pro_code IS NOT NULL AND mp.pro_code = product_cust.pro_code LIMIT 1) AS product_parent ON TRUE", filter.ParentCustId).
		Where("sls.order_detail.cust_id = ?", filter.CustId).
		Where("sls.order_detail.item_type = 1")

	if filter.StartDate > 0 {
		startDate := str.UnixTimestampToUtcTime(filter.StartDate)
		query = query.Where("sls.order.ro_date >= ?", startDate)
	}

	if filter.EndDate > 0 {
		endDate := str.UnixTimestampToUtcTime(filter.EndDate)
		query = query.Where("sls.order.ro_date <= ?", endDate)
	}

	if len(filter.SalesmanId) > 0 {
		query = query.Where("sls.order.salesman_id IN ?", filter.SalesmanId)
	}

	err := query.Find(&results).Error
	return results, err
}
