package repository

import (
	"context"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryValidateOrderImpl struct {
		*gorm.DB
	}
)
type ValidateOrderRepository interface {
	StockReport(dataFilter entity.ValidateOrderBody) ([]model.StockReport, int64, int, error)
	FindProductByListID(productIDs []int64) (products []model.ProductOrder, err error)
	FindAllArByCustId(dataFilter entity.ValidateOrderBody) ([]model.ArList, int64, int, error)
	FindAllArDetailByCustId(dataFilter entity.ValidateOrderDetailBody) ([]model.ArList, int64, int, error)
	CountInvoicePaidAmount(invoiceNo string, custId string) (invoice model.InvoicePaidAmount, err error)
	FindDetailOutletID(OutletID int64) (outlets model.OutletValidate, err error)
}

func NewValidateOrderRepo(db *gorm.DB) *RepositoryValidateOrderImpl {
	return &RepositoryValidateOrderImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryValidateOrderImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryValidateOrderImpl) StockReport(dataFilter entity.ValidateOrderBody) ([]model.StockReport, int64, int, error) {
	var stockReport []model.StockReport
	var total int64

	qDate := ""
	if dataFilter.Date != "" {
		qDate = "AND st.stock_date <= '" + dataFilter.Date + "'"
	}

	subQuery := repository.Select(`pro.pro_id, pro.pro_code, pro.pro_name, 
		pro.unit_id1, pro.unit_id2, pro.unit_id3, pro.conv_unit2, pro.conv_unit3, 
		pro.sup_id, pro.is_active, pro.deleted_at,
		COALESCE(SUM(st.qty_in), 0)-COALESCE(SUM(st.qty_out), 0) AS qty`).
		Joins("LEFT JOIN inv.stock st ON st.pro_id = pro.pro_id "+qDate+" AND st.cust_id = ? AND st.wh_id = ?", dataFilter.CustID, dataFilter.WhID).
		Where("pro.cust_id = ?", dataFilter.ParentCustID).Table("mst.m_product AS pro")

	isActiveProductOnly, _ := strconv.ParseBool(dataFilter.ActiveProductOnly)
	if isActiveProductOnly {
		subQuery.Where("pro.is_active = true")
	}

	if len(dataFilter.ProID) > 0 {
		subQuery.Where("pro.pro_id IN ?", dataFilter.ProID)
	}
	// if len(dataFilter.WhID) > 0 {
	// 	subQuery.Where("st.wh_id IN ?", dataFilter.WhID)
	// }
	subQuery.Group("pro.pro_id")

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
		subQuery.Order("pro." + sortBy)
	} else {
		subQuery.Order("pro.pro_id ASC")
	}

	// Using subquery in FROM clause
	query := repository.Table("(?) AS stock_report", subQuery)

	query.Find(&stockReport)

	total = int64(len(stockReport))
	return stockReport, total, 1, nil
}

func (repository *RepositoryValidateOrderImpl) FindProductByListID(productIDs []int64) (products []model.ProductOrder, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryValidateOrderImpl) FindAllArByCustId(dataFilter entity.ValidateOrderBody) ([]model.ArList, int64, int, error) {
	var ro []model.ArList
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("invoice_no")
	query := repository.Select(
		`sls.order.invoice_no, 
			sls.order.invoice_date, 
			sls.order.due_date, 
			sls.order.outlet_id, 
			sls.order.salesman_id, 
			sls.order.total as invoice_amount,
			ot.outlet_code, ot.outlet_name, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.order.salesman_id AND employee.cust_id = ?", dataFilter.CustID).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustID)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustID)
	query.Where("sls.order.cust_id=?", dataFilter.CustID)

	queryCount.Where("sls.order.outlet_id = ?", dataFilter.OutletID)
	query.Where("sls.order.outlet_id = ?", dataFilter.OutletID)

	queryCount.Where("sls.order.invoice_no IS NOT NULL")
	query.Where("sls.order.invoice_no IS NOT NULL")

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
		query.Order("invoice_no DESC")
	}

	err := query.Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryValidateOrderImpl) FindAllArDetailByCustId(dataFilter entity.ValidateOrderDetailBody) ([]model.ArList, int64, int, error) {
	var ro []model.ArList
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("invoice_no")
	query := repository.Select(
		`sls.order.invoice_no, 
			sls.order.invoice_date, 
			sls.order.due_date, 
			sls.order.outlet_id, 
			sls.order.salesman_id, 
			sls.order.total as invoice_amount,
			ot.outlet_code, ot.outlet_name, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.order.salesman_id AND employee.cust_id = ?", dataFilter.CustID).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustID)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustID)
	query.Where("sls.order.cust_id=?", dataFilter.CustID)

	queryCount.Where("sls.order.outlet_id = ?", dataFilter.OutletID)
	query.Where("sls.order.outlet_id = ?", dataFilter.OutletID)

	queryCount.Where("sls.order.invoice_no IS NOT NULL")
	query.Where("sls.order.invoice_no IS NOT NULL")

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
		query.Order("invoice_no DESC")
	}

	err := query.Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryValidateOrderImpl) CountInvoicePaidAmount(invoiceNo string, custId string) (invoice model.InvoicePaidAmount, err error) {
	err = repository.Select(`
			coalesce(sum(acf.deposit_detail.total_payment), 0) as paid_amount
		`).
		Joins("inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ? AND deposit.deposit_status IN (1, 2)", custId).
		Where("acf.deposit_detail.invoice_no = ?", invoiceNo).
		Take(&invoice).Error

	return invoice, err
}

func (repository *RepositoryValidateOrderImpl) FindDetailOutletID(OutletID int64) (outlets model.OutletValidate, err error) {
	err = repository.
		Select("*").
		Where("outlet_id = ?", OutletID).
		Take(&outlets).Error

	return outlets, err
}
