package repository

import (
	"context"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/str"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryOrderHistoryImpl struct {
		*gorm.DB
	}
)
type OrderHistoryRepository interface {
	FindDetail(RoNo string, custId string) (details []model.OrderDetailRead, err error)
	FindByNo(RoNo string, custId string) (realOrder model.OrderList, err error)
	FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error)
	FindDetailProductOrderHistory(roNo string) (details []model.OrderHistoryProduct, err error)
	FindDetailProductOrderHistoryPayment(roNo string) (details []model.OrderList, err error)
}

func NewOrderHistoryRepo(db *gorm.DB) *RepositoryOrderHistoryImpl {
	return &RepositoryOrderHistoryImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderHistoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderHistoryImpl) FindByNo(roNo string, custId string) (realOrder model.OrderList, err error) {
	err = repository.
		Select(`sls.order.*, 
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name,
			sls.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.order.ro_no = ? AND sls.order.cust_id=?", roNo, custId).
		Take(&realOrder).Error
	return realOrder, err
}

func (repository *RepositoryOrderHistoryImpl) FindDetail(roNo string, custId string) (details []model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("ro_no = ? AND sls.order_detail.cust_id=?", roNo, custId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderHistoryImpl) FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error) {
	var ro []model.OrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ro_no")
	query := repository.Select(
		`sls.order.*, 
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, ot.outlet_name, 
			sls.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.RoFrom != nil && dataFilter.RoTo != nil {
		query.Where("sls.order.ro_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.RoFrom), str.UnixTimestampToUtcTime(*dataFilter.RoTo))
		queryCount.Where("sls.order.ro_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.RoFrom), str.UnixTimestampToUtcTime(*dataFilter.RoTo))
	}

	if dataFilter.InvoiceFrom != nil && dataFilter.InvoiceTo != nil {
		query.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
		queryCount.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.order.ro_no=?", dataFilter.Query)
		query.Where("sls.order.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.order.data_status in ?", dataFilter.Status)
		query.Where("sls.order.data_status in ?", dataFilter.Status)
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
		query.Order("ro_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
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

func (repository *RepositoryOrderHistoryImpl) FindDetailProductOrderHistory(roNo string) (details []model.OrderHistoryProduct, err error) {
	err = repository.Select("sls.order_detail.ro_no, p.* ").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("sls.order_detail.ro_no = ?", roNo).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderHistoryImpl) FindDetailProductOrderHistoryPayment(roNo string) (details []model.OrderList, err error) {
	err = repository.Select("sls.order.invoice_no").
		Joins("left join acf.deposit_payment dp on dp.invoice_no = sls.order.invoice_no").
		Where("sls.order.ro_no = ?", roNo).
		Find(&details).Error
	return details, err
}
