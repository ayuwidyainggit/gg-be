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

var (
	// ErrOrdersNotFound indicates that no orders were found for the given ro_no
	ErrOrdersNotFound = errors.New("orders not found")
)

type (
	RepositoryOrderImpl struct {
		*gorm.DB
	}
)
type OrderRepository interface {
	Store(c context.Context, data *model.Order) error
	StoreDetail(c context.Context, data *model.OrderDetail) error
	StoreReward(c context.Context, data *model.OrderReward) error
	FindByNo(RoNo string, custId string) (realOrder model.OrderList, err error)
	FindDetail(RoNo string, custId string) (details []model.OrderDetailRead, err error)
	FindReward(RoNo string, custId string) (details []model.OrderRewardRead, err error)
	FindDetailByDetailID(detailID int64, RoNo string, custId string) (details model.OrderDetailRead, err error)
	FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error)
	FindProformaInvoiceList(dataFilter entity.ProformaInvoiceQueryFilter) ([]model.OrderList, int64, int, error)
	FindAllOutletBySalesmanId(dataFilter entity.OrderQueryFilter) ([]model.OutletBySalesman, int64, int, error)
	FindOrdersByRoNos(ctx context.Context, roNos []string, custId string) ([]model.OrderList, error)
	FindOrderDetailsForProforma(ctx context.Context, roNos []string, custId string) ([]model.OrderDetailRead, error)
	UpdateProformaInvoiceFlags(ctx context.Context, roNos []string, custId string, userId int64) error
	FindDetailByNotInDetailIDs(detailIDs []int64, roNo string, custId string) (details []model.OrderDetailRead, err error)
	Update(c context.Context, RoNo, custID string, data model.Order) error
	DeleteDetailNotInIDs(c context.Context, RoNo, custID string, IDs []int64) error
	DeletePromoDetails(c context.Context, RoNo, custID string) error
	DeleteRewards(c context.Context, RoNo, custID string) error
	UpdateDetail(c context.Context, Details *model.OrderDetail) error
	Delete(c context.Context, custId string, RoNo string, deletedBy int64) error

	FindOneProductByProductIdAndCustId(productId int64, custId string, parentCustId string) (model.ProductConversion, error)

	CountAllRoByCustId(custId string, roDate string) (int, error)

	FindProductByListID(productIDs []int64) (products []model.Product, err error)
	FindSalesmanByCode(salesmanCode string, custId string) (detail model.SalesmanDetail, err error)
	FindWarehouseByCode(whCode string, custId string) (warehouse model.WarehouseLookup, err error)
	FindOutletByCode(outletCode string, custId string) (outlet model.OutletRead, err error)
	FindProductByCode(proCode string, custId string) (product model.ProductRead, err error)
	FindProductByName(proName string, custId string) (product model.ProductRead, err error)

	FindDiscountCriteria(proID int, outletID int, effectiveDate *int64, slabAmount float64) (criterias model.DiscountCriteria, err error)
	// Fungsi untuk consult Discount saat Create Order
	FindOutletByID(outletID int, custId string, parentCustId string) (outlet model.OutletRead, err error)
	FindProductByID(productID int) (product model.ProductRead, err error)
	FindDiscountByProductAndOutlet(product model.ProductRead, outlet model.OutletRead) (discount model.DiscountRead, err error)
	FindDiscountCriteriaBySubTotal(discountID string, subTotal int) (discountCriteria model.DiscountCriteria, err error)
	FindAllDiscountPrincipalsByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountPrincipalDetail, err error)

	FindByInvoiceNo(invoiceNo string, custId string) (order model.OrderList, err error)
	FindByNoNoCustID(roNo, custIDOrigin string) (realOrder model.OrderList, err error)
	FindDetailNoCustID(roNo, custIDOrigin string) (details []model.OrderDetailRead, err error)
	FindRewardNoCustID(roNo, custIDOrigin string) (rewards []model.OrderRewardRead, err error)
	FindFullPromoRewards(invoiceNo string, custID string) (rewards []model.FullPromoRewardRead, err error)
	FindOrderApprovalRequestDetailByRoAndEmp(roNo string, empID int64) (detail model.OrderApprovalRequestDetailRead, err error)
	FindSalesman(salesmanID int64, custId string) (detail model.SalesmanDetail, err error)
	FindMinimumPriceActiveByProID(proID int64, custId string) (detail model.ManageMinimumPrice, err error)
	FindOrderDetailByDetailID(detailID int64, custId string) (details model.OrderDetailRead, err error)
	FindOrderDetailsByIDs(detailIDs []int64, custId string) (details []model.OrderDetailRead, err error)
	UpdateDetailPartial(c context.Context, orderDetailId int64, custId string, updates map[string]interface{}) error
	FindWarehouseStockByWhIdAndProIds(custId string, whId int64, proIds []int64) (map[int64]float64, error)
	RefreshOrderDetailStock(c context.Context, orderDetailId int64, qty1Stok, qty2Stok, qty3Stok float64) error
	SyncFinalOrderFields(c context.Context, orderDetailId int64) error

	LockOrderByScope(ctx context.Context, custId string, roDates []time.Time) error
	DeleteOrderDetailByScope(ctx context.Context, custId string, roDates []time.Time) (int64, error)
	DeleteOrderByScope(ctx context.Context, custId string, roDates []time.Time) (int64, error)
}

func NewOrderRepo(db *gorm.DB) *RepositoryOrderImpl {
	return &RepositoryOrderImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderImpl) Store(c context.Context, data *model.Order) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) StoreDetail(c context.Context, data *model.OrderDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) StoreReward(c context.Context, data *model.OrderReward) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) FindByNo(roNo string, custId string) (realOrder model.OrderList, err error) {
	err = repository.
		Select(`sls.order.*,
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name, sls.order.address1, ot.address2, ot.inv_addr1, ot.inv_addr2,
			sls.sales_name,emp.emp_code as salesman_code,
			wh.wh_code, wh.wh_name,
			CASE 
				WHEN ot.obs_limit_action = 1 THEN 1
				WHEN ot.obs_limit_action = 2 THEN 2
				ELSE NULL
			END AS obs_limit_action,
			CASE 
				WHEN obs_limit_action = 1 THEN 'Warning'
				WHEN obs_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS obs_limit_action_name,
			CASE 
				WHEN ot.sales_inv_limit_action = 1 THEN 1
				WHEN ot.sales_inv_limit_action = 2 THEN 2
				ELSE NULL
			END AS sales_inv_limit_action,
			CASE 
				WHEN sales_inv_limit_action = 1 THEN 'Warning'
				WHEN sales_inv_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS sales_inv_limit_action_name,
			CASE 
				WHEN ot.credit_limit_action = 1 THEN 1
				WHEN ot.credit_limit_action = 2 THEN 2
				ELSE NULL
			END AS credit_limit_action,
			CASE 
				WHEN credit_limit_action = 1 THEN 'Warning'
				WHEN credit_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS credit_limit_action_name,
			CASE WHEN ot.credit_limit_type = 2 THEN 2 ELSE NULL END AS credit_limit_type,
			CASE WHEN ot.sales_inv_limit_type = 2 THEN 2 ELSE NULL END AS sales_inv_limit_type,
			CASE WHEN ot.obs_type = 2 THEN 2 ELSE NULL END AS obs_type,
			oar.order_approval_request_id`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_employee emp on emp.emp_id = sls.order.salesman_id AND emp.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", custId).
		Joins("left join sls.order_approval_requests oar on sls.order.ro_no = oar.ro_no AND oar.cust_id = ? AND oar.finished_at is null", custId).
		Where("sls.order.ro_no = ? AND sls.order.cust_id=?", roNo, custId).
		Take(&realOrder).Error
	return realOrder, err
}

func (repository *RepositoryOrderImpl) FindDetail(roNo string, custId string) (details []model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("ro_no = ? AND sls.order_detail.cust_id=?", roNo, custId).
		Order("sls.order_detail.item_type ASC").
		Order("sls.order_detail.order_detail_id ASC").
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderImpl) FindReward(roNo string, custId string) (rewards []model.OrderRewardRead, err error) {
	err = repository.Select(`
			sls.order_reward.*,
			CASE
				WHEN promo.promo_desc IS NOT NULL THEN promo.promo_desc
				WHEN discount.discount_desc IS NOT NULL THEN discount.discount_desc
			ELSE
				''
			END AS "reff_name"
			`).
		Joins("left join sls.promotions promo on promo.promo_id = sls.order_reward.reff_id").
		Joins("left join sls.discounts discount on discount.discount_id = sls.order_reward.reff_id").
		Where("ro_no = ? AND sls.order_reward.cust_id=?", roNo, custId).
		Find(&rewards).Error

	return rewards, err
}

func (repository *RepositoryOrderImpl) FindDetailByDetailID(detailID int64, roNo string, custId string) (detail model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("order_detail_id = ? AND ro_no = ? AND sls.order_detail.cust_id = ?", detailID, roNo, custId).
		Take(&detail).Error

	return detail, err
}

func (repository *RepositoryOrderImpl) FindDetailByNotInDetailIDs(detailIDs []int64, roNo string, custId string) (details []model.OrderDetailRead, err error) {

	if len(detailIDs) == 0 {
		detailIDs = append(detailIDs, 0)
	}

	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("order_detail_id not in ? AND ro_no = ? AND sls.order_detail.cust_id = ? AND sls.order_detail.item_type = 1", detailIDs, roNo, custId).
		Find(&details).Error

	return details, err
}

func (repository *RepositoryOrderImpl) FindDiscountCriteria(proID int, outletID int, effectiveDate *int64, slabAmount float64) (criterias model.DiscountCriteria, err error) {
	err = repository.Select("discount_criterias.*").
		Joins("JOIN sls.discounts d ON d.discount_id = discount_criterias.discount_id").
		Joins("JOIN sls.discount_principals dp ON dp.discount_id = d.discount_id").
		Joins("JOIN mst.m_product mp ON mp.principal_id = dp.principal_id").
		Joins("JOIN sls.discount_groups dg ON dg.discount_id = d.discount_id").
		Joins("JOIN mst.m_outlet mo ON mo.disc_grp_id = dg.disc_grp_id").
		Where("mp.pro_id = ? AND mo.outlet_id = ?", proID, outletID).
		Where("? BETWEEN d.effective_from AND d.effective_to", str.UnixTimestampToUtcTime(*effectiveDate)).
		Where("discount_criterias.slab_rule_from > ? AND discount_criterias.slab_rule_to < ?", slabAmount, slabAmount).
		Take(&criterias).Error
	return criterias, err
}

func (repository *RepositoryOrderImpl) FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error) {
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
			ot.outlet_code, ot.outlet_name, COALESCE(sls.order.address1, ot.address1) AS address1, ot.address2,
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

	if dataFilter.IsInvoice != nil && *dataFilter.IsInvoice {
		queryCount.Where("sls.order.invoice_no IS NOT NULL")
		query.Where("sls.order.invoice_no IS NOT NULL")
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

func (repository *RepositoryOrderImpl) FindProformaInvoiceList(dataFilter entity.ProformaInvoiceQueryFilter) ([]model.OrderList, int64, int, error) {
	var ro []model.OrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 5
	} else {
		limit = dataFilter.Limit
	}

	var orderModel model.OrderList
	queryCount := repository.Table("sls.order")
	query := repository.Select(
		`sls.order.*,
			ot.outlet_code, ot.outlet_name, COALESCE(sls.order.address1, ot.address1) AS address1, ot.address2,
			emp.emp_code AS salesman_code,
			sls.sales_name AS salesman_name`).
		Table("sls.order").
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_employee emp on emp.emp_id = sls.order.salesman_id AND emp.cust_id = ?", dataFilter.CustId)

	// Mandatory filter: data_status = 2 (PROCESSED)
	queryCount.Where("sls.order.cust_id=? AND sls.order.data_status=?", dataFilter.CustId, entity.PROCESSED)
	query.Where("sls.order.cust_id=? AND sls.order.data_status=?", dataFilter.CustId, entity.PROCESSED)

	// Date range filter (ro_date between start_date and end_date)
	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		startTime := str.UnixTimestampToUtcTime(*dataFilter.StartDate)
		endTime := str.UnixTimestampToUtcTime(*dataFilter.EndDate)
		// Convert to date for ro_date comparison
		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, time.UTC)
		queryCount.Where("sls.order.ro_date between ? AND ?", startDate, endDate)
		query.Where("sls.order.ro_date between ? AND ?", startDate, endDate)
	}

	// Salesman filter
	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	// Outlet filter (optional)
	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	// Sort
	sortBy := ""
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				// Map created_date to created_at for database column
				colName := colSort[0]
				if colName == "created_date" {
					colName = "sls.order.created_at"
				} else {
					colName = "sls.order." + colName
				}
				sortBy += fmt.Sprintf(`%s %s, `, colName, colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("sls.order.created_at DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&orderModel).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryOrderImpl) Update(c context.Context, RoNo, custID string, data model.Order) error {
	result := repository.model(c).Model(&data).Where("ro_no=? AND cust_id = ?", RoNo, custID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryOrderImpl) DeleteDetailNotInIDs(c context.Context, RoNo, custID string, IDs []int64) error {
	var Details model.OrderDetail

	if len(IDs) == 0 {
		IDs = append(IDs, 0)
	}

	err := repository.model(c).Where("ro_no=? AND cust_id = ? AND order_detail_id not in (?) AND item_type = 1", RoNo, custID, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryOrderImpl) DeletePromoDetails(c context.Context, RoNo, custID string) error {
	var Details model.OrderDetail
	err := repository.model(c).Where("ro_no=? AND cust_id =? AND item_type = 2", RoNo, custID).Delete(&Details).Error
	return err
}

func (repository *RepositoryOrderImpl) DeleteRewards(c context.Context, RoNo, custID string) error {
	var Reward model.OrderReward
	err := repository.model(c).Where("ro_no=? AND cust_id = ?", RoNo, custID).Delete(&Reward).Error
	return err
}

func (repository *RepositoryOrderImpl) UpdateDetail(c context.Context, Details *model.OrderDetail) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}

	// if Details.DiscountID == nil {
	// 	result := repository.model(c).Updates(map[string]interface{"discount_id": nil})
	// 	if result.Error != nil {
	// 		return result.Error
	// 	}
	// }
	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryOrderImpl) Delete(c context.Context, custId string, RoNo string, deletedBy int64) error {
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

func (repository *RepositoryOrderImpl) FindOneProductByProductIdAndCustId(productId int64, custId string, parentCustId string) (productConversion model.ProductConversion, err error) {
	err = repository.
		Select(`mst.m_product.cust_id,
			mst.m_product.pro_id,
			mst.m_product.conv_unit2, 
			mst.m_product.conv_unit3, 
			mst.m_product.conv_unit4, 
			mst.m_product.conv_unit5`).
		Where("mst.m_product.pro_id = ? AND mst.m_product.cust_id = ? AND mst.m_product.is_active = ? AND mst.m_product.is_del = ?", productId, parentCustId, true, false).
		Order("mst.m_product.updated_at DESC NULLS LAST, mst.m_product.created_at DESC NULLS LAST").
		Take(&productConversion).Error
	return productConversion, err
}

func (repository *RepositoryOrderImpl) CountAllRoByCustId(custId string, roDate string) (int, error) {
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

func (repository *RepositoryOrderImpl) FindAllOutletBySalesmanId(dataFilter entity.OrderQueryFilter) ([]model.OutletBySalesman, int64, int, error) {
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

func (repository *RepositoryOrderImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryOrderImpl) FindOutletByID(outletID int, custId string, parentCustId string) (outlet model.OutletRead, err error) {
	err = repository.Select(`
			outlet_id, 
			outlet_code, 
			outlet_name, 
			address1,
			address2,
			disc_grp_id,
			CASE 
				WHEN credit_limit_action = 1 THEN 1
				WHEN credit_limit_action = 2 THEN 2
				ELSE NULL
			END AS credit_limit_action,
			CASE 
				WHEN credit_limit_action = 1 THEN 'Warning'
				WHEN credit_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS credit_limit_action_name,
			CASE 
				WHEN sales_inv_limit_action = 1 THEN 1
				WHEN sales_inv_limit_action = 2 THEN 2
				ELSE NULL
			END AS sales_inv_limit_action,
			CASE 
				WHEN sales_inv_limit_action = 1 THEN 'Warning'
				WHEN sales_inv_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS sales_inv_limit_action_name,
			CASE 
				WHEN obs_limit_action = 1 THEN 1
				WHEN obs_limit_action = 2 THEN 2
				ELSE NULL
			END AS obs_limit_action,
			CASE 
				WHEN obs_limit_action = 1 THEN 'Warning'
				WHEN obs_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS obs_limit_action_name
		`).
		Where("mst.m_outlet.outlet_id=?", outletID).
		Where("mst.m_outlet.cust_id IN ?", []string{custId, parentCustId}).
		Order(gorm.Expr("CASE WHEN mst.m_outlet.cust_id = ? THEN 0 ELSE 1 END", custId)).
		Take(&outlet).Error

	return outlet, err
}

func (repository *RepositoryOrderImpl) FindProductByID(productID int) (product model.ProductRead, err error) {
	err = repository.Select(`
			mst.m_product.*
		`).
		Where("mst.m_product.pro_id=?", productID).
		Take(&product).Error

	return product, err
}

func (repository *RepositoryOrderImpl) FindDiscountByProductAndOutlet(product model.ProductRead, outlet model.OutletRead) (discount model.DiscountRead, err error) {
	err = repository.Select(`
			sls.discounts.discount_id, 
			sls.discounts.discount_desc,
			sls.discounts.created_at
		`).
		// Joins("inner join").
		Where("sls.discounts.discount_status_id = 2").
		Where("sls.discounts.publish_status_id = 2").
		Where("sls.discounts.discount_id IN (SELECT discount_id from sls.discount_principals WHERE principal_id = ?)", product.PrincipalId).
		Where("sls.discounts.discount_id IN (SELECT discount_id from sls.discount_groups WHERE disc_grp_id = ?)", outlet.DiscGrpId).
		Order("sls.discounts.created_at DESC").
		Take(&discount).Error

	return discount, err
}

func (repository *RepositoryOrderImpl) FindDiscountCriteriaBySubTotal(discountID string, subTotal int) (discountCriteria model.DiscountCriteria, err error) {
	err = repository.Select(`
			sls.discount_criterias.*
		`).
		Where("sls.discount_criterias.discount_id = ?", discountID).
		Where("sls.discount_criterias.slab_rule_from <= ?", subTotal).
		Where("sls.discount_criterias.slab_rule_to >= ?", subTotal).
		Take(&discountCriteria).Error

	return discountCriteria, err
}

func (repository *RepositoryOrderImpl) FindAllDiscountPrincipalsByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountPrincipalDetail, err error) {
	err = repository.
		Select(`discount_principals.*, princ.principal_code, princ.principal_name`).
		Joins("LEFT JOIN mst.m_principal princ ON princ.principal_id = discount_principals.principal_id AND princ.cust_id = ?", params.ParentCustId).
		Where("discount_principals.discount_id = ? AND discount_principals.cust_id = ?", params.DiscountID, params.CustID).
		Order("discount_principals.principal_id ASC").
		Find(&discountCriterias).Error
	return discountCriterias, err
}

func (repository *RepositoryOrderImpl) FindByNoNoCustID(roNo, custIDOrigin string) (realOrder model.OrderList, err error) {
	err = repository.
		Select(`sls.order.*,
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name, sls.order.address1, ot.address2,
			sls.sales_name,emp.emp_code as salesman_code,
			wh.wh_code, wh.wh_name,
			CASE 
				WHEN ot.obs_limit_action = 1 THEN 1
				WHEN ot.obs_limit_action = 2 THEN 2
				ELSE NULL
			END AS obs_limit_action,
			CASE 
				WHEN obs_limit_action = 1 THEN 'Warning'
				WHEN obs_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS obs_limit_action_name,
			CASE 
				WHEN ot.sales_inv_limit_action = 1 THEN 1
				WHEN ot.sales_inv_limit_action = 2 THEN 2
				ELSE NULL
			END AS sales_inv_limit_action,
			CASE 
				WHEN sales_inv_limit_action = 1 THEN 'Warning'
				WHEN sales_inv_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS sales_inv_limit_action_name,
			CASE 
				WHEN ot.credit_limit_action = 1 THEN 1
				WHEN ot.credit_limit_action = 2 THEN 2
				ELSE NULL
			END AS credit_limit_action,
			CASE 
				WHEN credit_limit_action = 1 THEN 'Warning'
				WHEN credit_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS credit_limit_action_name,
			CASE WHEN ot.credit_limit_type = 2 THEN 2 ELSE NULL END AS credit_limit_type,
			CASE WHEN ot.sales_inv_limit_type = 2 THEN 2 ELSE NULL END AS sales_inv_limit_type,
			CASE WHEN ot.obs_type = 2 THEN 2 ELSE NULL END AS obs_type,
			oar.order_approval_request_id`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id").
		Joins("left join mst.m_employee emp on emp.emp_id = sls.order.salesman_id").
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id").
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id").
		Joins("left join sls.order_approval_requests oar on sls.order.ro_no = oar.ro_no AND oar.finished_at is null").
		Where("sls.order.ro_no = ? AND sls.cust_id = ?", roNo, custIDOrigin).
		Take(&realOrder).Error
	return realOrder, err
}

func (repository *RepositoryOrderImpl) FindDetailNoCustID(roNo, custIDOrigin string) (details []model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("ro_no = ? AND sls.order_detail.cust_id = ?", roNo, custIDOrigin).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderImpl) FindRewardNoCustID(roNo, custIDOrigin string) (rewards []model.OrderRewardRead, err error) {
	err = repository.Select(`
			sls.order_reward.*,
			CASE
				WHEN promo.promo_desc IS NOT NULL THEN promo.promo_desc
				WHEN discount.discount_desc IS NOT NULL THEN discount.discount_desc
			ELSE
				''
			END AS "reff_name"
			`).
		Joins("left join sls.promotions promo on promo.promo_id = sls.order_reward.reff_id").
		Joins("left join sls.discounts discount on discount.discount_id = sls.order_reward.reff_id").
		Where("ro_no = ? AND sls.order_reward.cust_id = ?", roNo, custIDOrigin).
		Find(&rewards).Error

	return rewards, err
}

func (repository *RepositoryOrderImpl) FindOrderApprovalRequestDetailByRoAndEmp(roNo string, empID int64) (detail model.OrderApprovalRequestDetailRead, err error) {
	err = repository.Select(`sls.order_approval_requests_details.*, emp.emp_name, emp.emp_code`).
		Joins("LEFT JOIN mst.m_employee emp ON sls.order_approval_requests_details.emp_id = emp.emp_id").
		Joins("LEFT JOIN sls.order_approval_requests oar on sls.order_approval_requests_details.order_approval_request_id = oar.order_approval_request_id").
		Where("oar.ro_no = ? AND sls.order_approval_requests_details.emp_id=? AND oar.finished_at is null", roNo, empID).
		Order("level ASC, seq ASC").
		Take(&detail).Error
	return detail, err
}

func (repository *RepositoryOrderImpl) FindSalesman(salesmanID int64, custId string) (detail model.SalesmanDetail, err error) {
	err = repository.Select(`
			mst.m_salesman.emp_id as salesman_id, 
			mst.m_salesman.sales_name as salesman_name, 
			mst.m_salesman.sales_team_id,
			mst.m_salesman.wh_id,
			mst.m_salesman.allow_input_price
		`).
		Where("mst.m_salesman.emp_id=?", salesmanID).
		Take(&detail).Error

	return detail, err
}

func (repository *RepositoryOrderImpl) FindMinimumPriceActiveByProID(proID int64, custId string) (detail model.ManageMinimumPrice, err error) {
	err = repository.Select(`
			*
		`).
		Where("mst.manage_minimum_price.cust_id=? AND mst.manage_minimum_price.pro_id=? AND mst.manage_minimum_price.status_manage_minimum_price=?", custId, proID, model.STATUS_MANAGE_PRICE_ACTIVE).
		Take(&detail).Error

	return detail, err
}

func (repository *RepositoryOrderImpl) FindByInvoiceNo(invoiceNo string, custId string) (order model.OrderList, err error) {
	err = repository.
		Select(`sls.order.*,
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name, sls.order.address1, ot.address2, ot.inv_addr1, ot.inv_addr2,
			sls.sales_name,emp.emp_code as salesman_code,
			wh.wh_code, wh.wh_name,
			CASE 
				WHEN ot.obs_limit_action = 1 THEN 1
				WHEN ot.obs_limit_action = 2 THEN 2
				ELSE NULL
			END AS obs_limit_action,
			CASE 
				WHEN obs_limit_action = 1 THEN 'Warning'
				WHEN obs_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS obs_limit_action_name,
			CASE 
				WHEN ot.sales_inv_limit_action = 1 THEN 1
				WHEN ot.sales_inv_limit_action = 2 THEN 2
				ELSE NULL
			END AS sales_inv_limit_action,
			CASE 
				WHEN sales_inv_limit_action = 1 THEN 'Warning'
				WHEN sales_inv_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS sales_inv_limit_action_name,
			CASE 
				WHEN ot.credit_limit_action = 1 THEN 1
				WHEN ot.credit_limit_action = 2 THEN 2
				ELSE NULL
			END AS credit_limit_action,
			CASE 
				WHEN credit_limit_action = 1 THEN 'Warning'
				WHEN credit_limit_action = 2 THEN 'Restricted'
				ELSE ''
			END AS credit_limit_action_name,
			CASE WHEN ot.credit_limit_type = 2 THEN 2 ELSE NULL END AS credit_limit_type,
			CASE WHEN ot.sales_inv_limit_type = 2 THEN 2 ELSE NULL END AS sales_inv_limit_type,
			CASE WHEN ot.obs_type = 2 THEN 2 ELSE NULL END AS obs_type,
			oar.order_approval_request_id`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_employee emp on emp.emp_id = sls.order.salesman_id AND emp.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", custId).
		Joins("left join sls.order_approval_requests oar on sls.order.ro_no = oar.ro_no AND oar.cust_id = ? AND oar.finished_at is null", custId).
		Where("sls.order.invoice_no = ? AND sls.order.cust_id=?", invoiceNo, custId).
		Take(&order).Error
	return order, err
}

func (repository *RepositoryOrderImpl) FindFullPromoRewards(invoiceNo string, custId string) (rewards []model.FullPromoRewardRead, err error) {
	err = repository.Select(`
			sls.order_reward.*,
			promo.promo_id,
			promo.promo_desc,
			promo.is_multiplied,
			criterias.slab_rule_type,
			criterias.slab_rule_from,
			criterias.slab_rule_to,
			criterias.slab_rule_uom,
			criterias.slab_reward_type,
			criterias.slab_reward,
			criterias.slab_reward_uom
			`).
		Joins("inner join sls.promotions promo on promo.promo_id = sls.order_reward.reff_id and promo.cust_id = ?", custId).
		Joins("inner join sls.order invoice on invoice.ro_no = sls.order_reward.ro_no and invoice.cust_id = ?", custId).
		Joins("inner join sls.promo_criterias criterias on criterias.slab_id = sls.order_reward.slab_id and criterias.cust_id = ?", custId).
		// Where("ro_no = ?", invoiceNo).
		Where("invoice.invoice_no=?", invoiceNo).
		Where("sls.order_reward.cust_id=?", custId).
		Where("sls.order_reward.reward_type_id=1").
		Find(&rewards).Error

	return rewards, err
}

func (repository *RepositoryOrderImpl) FindOrderDetailByDetailID(detailID int64, custId string) (detail model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("order_detail_id = ?", detailID).
		Where("sls.order_detail.cust_id = ?", custId).
		Take(&detail).Error

	return detail, err
}

func (repository *RepositoryOrderImpl) FindOrderDetailsByIDs(detailIDs []int64, custId string) (details []model.OrderDetailRead, err error) {
	if len(detailIDs) == 0 {
		return []model.OrderDetailRead{}, nil
	}

	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("sls.order_detail.order_detail_id IN ?", detailIDs).
		Where("sls.order_detail.cust_id = ?", custId).
		Find(&details).Error

	return details, err
}

func (repository *RepositoryOrderImpl) FindOrdersByRoNos(ctx context.Context, roNos []string, custId string) ([]model.OrderList, error) {
	var orders []model.OrderList
	query := repository.model(ctx).
		Select("sls.order.*, salesman.sales_name, ot.outlet_code, ot.outlet_name, sls.order.address1, ot.zip_code").
		Table("sls.order").
		Joins("left join mst.m_salesman salesman on salesman.emp_id = sls.order.salesman_id AND salesman.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.order.ro_no IN ? AND sls.order.cust_id = ?", roNos, custId)

	err := query.Find(&orders).Error
	return orders, err
}

func (repository *RepositoryOrderImpl) FindOrderDetailsForProforma(ctx context.Context, roNos []string, custId string) ([]model.OrderDetailRead, error) {
	var details []model.OrderDetailRead
	err := repository.model(ctx).
		Select("sls.order_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2, p.conv_unit3 as mconv_unit3").
		Table("sls.order_detail").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("sls.order_detail.ro_no IN ? AND sls.order_detail.cust_id = ?", roNos, custId).
		Order("sls.order_detail.item_type ASC").
		Order("sls.order_detail.order_detail_id ASC").
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderImpl) UpdateProformaInvoiceFlags(ctx context.Context, roNos []string, custId string, userId int64) error {
	// Validate roNos is not empty to prevent WHERE IN with empty slice
	if len(roNos) == 0 {
		return fmt.Errorf("roNos cannot be empty")
	}

	now := time.Now()

	// Only update if first_issue_date IS NULL
	// Update 3 fields: first_issue_date, is_proforma_inv, generate_by
	result := repository.model(ctx).
		Model(&model.Order{}).
		Where("ro_no IN ? AND cust_id = ? AND first_issue_date IS NULL", roNos, custId).
		Updates(map[string]interface{}{
			"first_issue_date": now,
			"is_proforma_inv":  true,
			"generate_by":      userId,
		})
	if result.Error != nil {
		return result.Error
	}

	// Check if at least one order was found (even if not updated)
	var count int64
	countResult := repository.model(ctx).
		Model(&model.Order{}).
		Where("ro_no IN ? AND cust_id = ?", roNos, custId).
		Count(&count)
	if countResult.Error != nil {
		return countResult.Error
	}
	if count == 0 {
		return fmt.Errorf("no orders found for ro_no: %v: %w", roNos, ErrOrdersNotFound)
	}

	// If first_issue_date already exists (NOT NULL), no update needed - this is expected behavior
	return nil
}

// UpdateDetailPartial updates specific fields of order_detail by order_detail_id
func (repository *RepositoryOrderImpl) UpdateDetailPartial(c context.Context, orderDetailId int64, custId string, updates map[string]interface{}) error {
	var detail model.OrderDetail
	result := repository.model(c).
		Model(&detail).
		Where("order_detail_id = ? AND cust_id = ?", orderDetailId, custId).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no order detail found with given id")
	}
	return nil
}

// FindWarehouseStockByWhIdAndProIds fetches warehouse stock qty for given warehouse and product IDs
// Returns map[pro_id] -> qty
// Source-of-truth: inv.stock ledger (same basis as Inventory stock report), not inv.warehouse_stock snapshot.
// Computes qty = SUM(qty_in) - SUM(qty_out) from the latest cumulative ledger available for the warehouse/product.
func (repository *RepositoryOrderImpl) FindWarehouseStockByWhIdAndProIds(custId string, whId int64, proIds []int64) (map[int64]float64, error) {
	result := make(map[int64]float64)

	if len(proIds) == 0 {
		return result, nil
	}

	type ledgerRow struct {
		ProID int64
		Qty   float64
	}

	var rows []ledgerRow
	err := repository.
		Table("inv.stock AS st").
		Select("st.pro_id AS pro_id, COALESCE(SUM(st.qty_in), 0) - COALESCE(SUM(st.qty_out), 0) AS qty").
		Where("st.cust_id = ? AND st.wh_id = ? AND st.pro_id IN ?", custId, whId, proIds).
		Group("st.pro_id").
		Scan(&rows).Error

	if err != nil {
		return result, err
	}

	for _, stock := range rows {
		result[stock.ProID] = stock.Qty
	}

	return result, nil
}

func (repository *RepositoryOrderImpl) RefreshOrderDetailStock(c context.Context, orderDetailId int64, qty1Stok, qty2Stok, qty3Stok float64) error {
	return repository.model(c).
		Table("sls.order_detail").
		Where("order_detail_id = ?", orderDetailId).
		Updates(map[string]interface{}{
			"qty1_stok": qty1Stok,
			"qty2_stok": qty2Stok,
			"qty3_stok": qty3Stok,
		}).Error
}

func (repository *RepositoryOrderImpl) FindSalesmanByCode(salesmanCode string, custId string) (detail model.SalesmanDetail, err error) {
	salesmanCode = strings.TrimSpace(salesmanCode)
	err = repository.Select(`mst.m_salesman.cust_id            AS cust_id,
	                          mst.m_salesman.emp_id             AS salesman_id,
	                          mst.m_employee.emp_code           AS salesman_code,
	                          mst.m_salesman.sales_name         AS salesman_name,
	                          COALESCE(mst.m_salesman.allow_input_price, false) AS allow_input_price,
	                          mst.m_salesman.wh_id              AS wh_id`).
		Joins("LEFT JOIN mst.m_employee ON mst.m_employee.emp_id = mst.m_salesman.emp_id AND mst.m_employee.cust_id = mst.m_salesman.cust_id").
		Where("mst.m_employee.emp_code = ? AND mst.m_salesman.cust_id = ?", salesmanCode, custId).
		Take(&detail).Error
	return detail, err
}

func (repository *RepositoryOrderImpl) FindWarehouseByCode(whCode string, custId string) (warehouse model.WarehouseLookup, err error) {
	err = repository.Select("mst.m_warehouse.*").
		Where("wh_code = ? AND cust_id = ? AND is_del = false", whCode, custId).
		Take(&warehouse).Error
	return warehouse, err
}

func (repository *RepositoryOrderImpl) FindOutletByCode(outletCode string, custId string) (outlet model.OutletRead, err error) {
	err = repository.Select("mst.m_outlet.*").
		Where("outlet_code = ? AND cust_id = ? AND is_del = false", outletCode, custId).
		Take(&outlet).Error
	return outlet, err
}

func (repository *RepositoryOrderImpl) FindProductByCode(proCode string, custId string) (product model.ProductRead, err error) {
	err = repository.Select("mst.m_product.*").
		Where("pro_code = ? AND cust_id = ? AND is_del = false", proCode, custId).
		Take(&product).Error
	return product, err
}

func (repository *RepositoryOrderImpl) FindProductByName(proName string, custId string) (product model.ProductRead, err error) {
	err = repository.Select("mst.m_product.*").
		Where("LOWER(pro_name) = LOWER(?) AND cust_id = ? AND is_del = false", proName, custId).
		Order("pro_id ASC").
		Take(&product).Error
	return product, err
}

func (repository *RepositoryOrderImpl) LockOrderByScope(ctx context.Context, custId string, roDates []time.Time) error {
	for _, d := range roDates {
		dateStr := d.Format("2006-01-02")
		key := custId + ":" + dateStr
		if err := repository.model(ctx).Exec("SELECT pg_advisory_xact_lock(hashtextextended($1, 0))", key).Error; err != nil {
			return err
		}
	}
	dateStrs := make([]string, len(roDates))
	for i, d := range roDates {
		dateStrs[i] = d.Format("2006-01-02")
	}
	var count int64
	if err := repository.model(ctx).Raw("SELECT 1 FROM sls.order WHERE cust_id=$1 AND is_sales_mapping=true AND ro_date = ANY($2::date[]) FOR UPDATE", custId, dateStrs).Scan(&count).Error; err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) DeleteOrderDetailByScope(ctx context.Context, custId string, roDates []time.Time) (int64, error) {
	dateStrs := make([]string, len(roDates))
	for i, d := range roDates {
		dateStrs[i] = d.Format("2006-01-02")
	}
	res := repository.model(ctx).Exec("DELETE FROM sls.order_detail d USING sls.order o WHERE o.ro_no=d.ro_no AND o.cust_id=$1 AND o.is_sales_mapping=true AND o.ro_date = ANY($2::date[])", custId, dateStrs)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (repository *RepositoryOrderImpl) DeleteOrderByScope(ctx context.Context, custId string, roDates []time.Time) (int64, error) {
	dateStrs := make([]string, len(roDates))
	for i, d := range roDates {
		dateStrs[i] = d.Format("2006-01-02")
	}
	res := repository.model(ctx).Exec("DELETE FROM sls.order WHERE cust_id=$1 AND is_sales_mapping=true AND ro_date = ANY($2::date[])", custId, dateStrs)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (repository *RepositoryOrderImpl) SyncFinalOrderFields(c context.Context, orderDetailId int64) error {
	// Syncs regular fields to final fields when order has no proforma invoice yet.
	return repository.model(c).
		Table("sls.order_detail as od").
		Joins("JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id").
		Where("od.order_detail_id = ?", orderDetailId).
		Where("o.is_proforma_inv IS NOT TRUE").
		Updates(map[string]interface{}{
			"qty_final":                  gorm.Expr("od.qty"),
			"qty1_final":                 gorm.Expr("od.qty1"),
			"qty2_final":                 gorm.Expr("od.qty2"),
			"qty3_final":                 gorm.Expr("od.qty3"),
			"disc_value_final":           gorm.Expr("od.disc_value"),
			"vat_value_final":            gorm.Expr("od.vat_value"),
			"amount_final":               gorm.Expr("od.amount"),
			"promo_value_final":          gorm.Expr("od.promo_value"),
			"promo_final1":               gorm.Expr("od.promo_so1"),
			"promo_final2":               gorm.Expr("od.promo_so2"),
			"promo_final3":               gorm.Expr("od.promo_so3"),
			"promo_final4":               gorm.Expr("od.promo_so4"),
			"promo_final5":               gorm.Expr("od.promo_so5"),
			"promo_remarks_final":        gorm.Expr("od.promo_remarks_so"),
			"is_product_promotion_final": gorm.Expr("od.is_product_promotion_so"),
		}).Error
}
