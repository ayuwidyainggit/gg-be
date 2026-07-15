package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/str"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryOrderCanvasImpl struct {
		*gorm.DB
	}
)
type OrderCanvasRepository interface {
	Store(c context.Context, data *model.Order) error
	StoreNoOrder(c context.Context, data *model.NoOrder) error
	StoreDetail(c context.Context, data *model.OrderDetail) error
	FindByNo(RoNo string, custId string) (realOrder model.OrderList, err error)
	FindDetail(RoNo string, custId string) (details []model.OrderDetailRead, err error)
	FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error)
	FindAllNoOrderByCustId(dataFilter entity.NoOrderQueryFilter) ([]model.NoOrderList, int64, int, error)
	FindAllOutletBySalesmanId(dataFilter entity.OrderQueryFilter) ([]model.OutletBySalesman, int64, int, error)

	Update(c context.Context, RoNo string, data model.Order) error
	DeleteDetailNotInIDs(c context.Context, RoNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.OrderDetail) error
	Delete(c context.Context, custId string, RoNo string, deletedBy int64) error

	FindOneProductByProductIdAndCustId(productId int64, custId string, parentCustId string) (model.ProductConversion, error)
	FindProductByListID(productIDs []int64) (products []model.ProductOrder, err error)
	CountAllRoByCustId(custId string, roDate string) (int, error)

	SummaryTotalBySalesmanAndDate(salesmanID int, startDate string, endDate string, dataStatus []int, parentCustId string) (summaryOrder model.SummaryOrder, err error)
}

func NewOrderCanvasRepo(db *gorm.DB) *RepositoryOrderCanvasImpl {
	return &RepositoryOrderCanvasImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderCanvasImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderCanvasImpl) Store(c context.Context, data *model.Order) error {
	fmt.Println("Ro Number Repo: ", data.RoNo)
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderCanvasImpl) StoreNoOrder(c context.Context, data *model.NoOrder) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderCanvasImpl) StoreDetail(c context.Context, data *model.OrderDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderCanvasImpl) FindByNo(roNo string, custId string) (realOrder model.OrderList, err error) {
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

func (repository *RepositoryOrderCanvasImpl) FindDetail(roNo string, custId string) (details []model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("ro_no = ? AND sls.order_detail.cust_id=?", roNo, custId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderCanvasImpl) FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error) {
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

	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
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

func (repository *RepositoryOrderCanvasImpl) FindAllNoOrderByCustId(dataFilter entity.NoOrderQueryFilter) ([]model.NoOrderList, int64, int, error) {
	var ro []model.NoOrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("no_order_id")
	query := repository.Select(
		`sls.no_order.*, 
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, ot.outlet_name, tor.taking_order_name, tor.image_url,
			sls.sales_name`).
		Joins("left join sys.m_user us on us.user_id = sls.no_order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.no_order.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.no_order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_taking_order tor on tor.taking_order_id = sls.no_order.taking_order_id AND tor.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.no_order.cust_id=?", dataFilter.CustId)
	query.Where("sls.no_order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.no_order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.no_order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.no_order.ro_no=?", dataFilter.Query)
		query.Where("sls.no_order.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.no_order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.no_order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.no_order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.no_order.outlet_id in ?", dataFilter.OutletID)
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
		query.Order("no_order_id DESC")
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

func (repository *RepositoryOrderCanvasImpl) Update(c context.Context, RoNo string, data model.Order) error {
	result := repository.model(c).Model(&data).Where("ro_no=?", RoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryOrderCanvasImpl) DeleteDetailNotInIDs(c context.Context, RoNo string, IDs []int64) error {
	var Details model.OrderDetail
	err := repository.model(c).Where("ro_no=? AND order_detail_id not in (?) ", RoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryOrderCanvasImpl) UpdateDetail(c context.Context, Details *model.OrderDetail) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryOrderCanvasImpl) Delete(c context.Context, custId string, RoNo string, deletedBy int64) error {
	var data model.Order
	result := repository.model(c).Model(&data).Where("ro_no=? AND cust_id = ? AND is_del= ? ", RoNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderCanvasImpl) FindOneProductByProductIdAndCustId(productId int64, custId string, parentCustId string) (productConversion model.ProductConversion, err error) {
	err = repository.
		Select(`mst.m_product.cust_id,
			mst.m_product.pro_id,
			mst.m_product.conv_unit2, 
			mst.m_product.conv_unit3, 
			mst.m_product.conv_unit4, 
			mst.m_product.conv_unit5`).
		Where("mst.m_product.pro_id = ? AND mst.m_product.cust_id=?", productId, parentCustId).
		Take(&productConversion).Error
	return productConversion, err
}

func (repository *RepositoryOrderCanvasImpl) CountAllRoByCustId(custId string, roDate string) (int, error) {
	var ro []model.OrderList
	var total int64

	queryCount := repository.Select("ro_no")

	queryCount.Where("sls.order.cust_id = ?", custId)
	queryCount.Where("sls.order.ro_date = ?", roDate) // Menambahkan kondisi tanggal sekarang

	// queryCount.Where("sls.order.data_status=4")

	err := queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return int(total), nil
}

func (repository *RepositoryOrderCanvasImpl) FindAllOutletBySalesmanId(dataFilter entity.OrderQueryFilter) ([]model.OutletBySalesman, int64, int, error) {
	var ro []model.OutletBySalesman
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Table("sls.order").
		Select("distinct(sls.order.outlet_id)").
		Joins("left join mst.m_outlet mo on mo.outlet_id = sls.order.outlet_id AND mo.cust_id = ?", dataFilter.CustId).
		Where("sls.order.cust_id = ?", dataFilter.CustId)

	query := repository.Table("sls.order").
		Select("distinct(sls.order.outlet_id) as outlet_id, mo.outlet_code, mo.outlet_name").
		Joins("left join mst.m_outlet mo on mo.outlet_id = sls.order.outlet_id AND mo.cust_id = ?", dataFilter.CustId).
		Where("sls.order.cust_id = ?", dataFilter.CustId)

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(mo.outlet_code ILIKE ? OR mo.outlet_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query = query.Where("(mo.outlet_code ILIKE ? OR mo.outlet_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
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

	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
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
		query.Order("outlet_id DESC")
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

func (repository *RepositoryOrderCanvasImpl) FindProductByListID(productIDs []int64) (products []model.ProductOrder, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryOrderCanvasImpl) SummaryTotalBySalesmanAndDate(salesmanID int, startDate string, endDate string, dataStatus []int, parentCustId string) (summaryOrder model.SummaryOrder, err error) {

	err = repository.Table("sls.order").
		Select("SUM(total) as total_summary").
		Where("salesman_id = ? AND DATE(ro_date) BETWEEN ? AND ?", salesmanID, startDate, endDate).
		Where("data_status in ?", dataStatus).
		Where("is_del = false").
		Scan(&summaryOrder).Error
	return summaryOrder, err
}
