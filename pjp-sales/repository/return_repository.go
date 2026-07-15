package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryReturnImpl struct {
		*gorm.DB
	}
)
type ReturnRepository interface {
	Store(c context.Context, data *model.Return) error
	StoreDetail(c context.Context, data *model.ReturnDetail) error
	FindOneByReturnNo(returnNo string, custId string, parentCustId string) (whAdj model.ReturnRead, err error)
	FindReturnDetail(returnNo string, custId string, parentCustId string) (Details []model.ReturnDetailRead, err error)
	CountReturnedProductQty(invoiceNo string, productId int64, custId string) (qtySummary model.ReturnedDetailRead, err error)
	FindAllByCustId(dataFilter entity.ReturnQueryFilter) ([]model.ReturnList, int64, int, error)
	// Delete(c context.Context, custId string, returnNo string, deletedBy int64) error
	Update(c context.Context, data *model.Return) error
	DeleteDetailNotInIDs(c context.Context, returnNo string, IDs []int64, custId string) error
	UpdateDetail(c context.Context, Details *model.ReturnDetail) error
	UpdateQuantity(c context.Context, Details *model.ReturnQuantity) error
	// UpdateReturnCost(c context.Context, Details *model.Return) error
	Print(c context.Context, custId string, returnNo string, printedBy int64) error

	FindAllSalesmanFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (salesmans []model.SalesmansFilter, total int64, lastPage int, err error)
	FindAllEmployeeByEmpGrpIdFilterByCustIdLookupMode(dataFilter entity.SalesmanQueryFilter) (salesmans []model.Employee, total int64, lastPage int, err error)
	FindAllRolesFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outlets []model.EmpGroupsFilter, total int64, lastPage int, err error)
	FindAllOutletFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outlets []model.OutletsFilter, total int64, lastPage int, err error)
	FindAllReturnStatusesLookupMode(dataFilter entity.GeneralQueryFilter) (returnStatus []model.ReturnStatusesFilter, total int64, lastPage int, err error)
	FindAllSalesmanFilterByCustIdLookupModeCreate(dataFilter entity.GeneralQueryFilter) (salesmans []model.SalesmansFilterCreate, total int64, lastPage int, err error)
	FindAllOutletFilterByCustIdLookupModeCreate(dataFilter entity.OutletCreateReturnQueryFilter) (outlets []model.OutletsFilterCreate, total int64, lastPage int, err error)
	FindAllProductByCustId(dataFilter entity.ProductListQueryFilter) ([]model.ProductList, int64, int, error)
	FindAllMasterReturnReasonLookupMode(dataFilter entity.GeneralQueryFilter) (outlets []model.ReturnReasonLookup, total int64, lastPage int, err error)
	FindAllMasterWarehouseLookupMode(dataFilter entity.WarehouseQueryFilter) (outlets []model.WarehouseLookup, total int64, lastPage int, err error)
	FindSalesmanById(salesmanId int, custId string, parentCustId string) (salesman []model.SalesmanRead, err error)
	FindAllMasterProductByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ProductsLookup, int64, int, error)
}

func NewReturnRepo(db *gorm.DB) *RepositoryReturnImpl {
	return &RepositoryReturnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryReturnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryReturnImpl) Store(c context.Context, data *model.Return) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryReturnImpl) StoreDetail(c context.Context, data *model.ReturnDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryReturnImpl) FindOneByReturnNo(returnNo string, custId string, parentCustId string) (rtn model.ReturnRead, err error) {
	err = repository.Select(`
			sls.return.cust_id, 
			sls.return.refference_no, 
			sls.return.return_no, 
			sls.return.return_date, 
			sls.return.invoice_no, 
			sls.return.invoice_date, 
			sls.return.salesman_id, 
			salesman.emp_code as salesman_code, 
			salesman.emp_name as salesman_name, 
			sls.return.emp_id, 
			employee.emp_code as emp_code, 
			employee.emp_name as emp_name, 
			eg.emp_grp_code,
			eg.emp_grp_name,
			sls.return.outlet_id, 
			ot.outlet_code, 
			ot.outlet_name,
			ot.address1 as outlet_address,
			sls.return.tpr_cash_value, 
			sls.return.tpr_item_value, 
			sls.return.promo_value, 
			sls.return.promo_bg_value, 
			sls.return.discount, 
			sls.return.disc_value, 
			sls.return.vat, 
			sls.return.vat_value, 
			sls.return.sub_total, 
			sls.return.total, 
			sls.return.data_status,
			sls.return.is_printed,
			sls.return.printed_by,
			sls.return.printed_at,
			printer.user_fullname AS printed_by_name
		`).
		Joins("left join mst.m_employee salesman on salesman.emp_id = sls.return.salesman_id AND salesman.cust_id = ?", custId).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.return.emp_id AND employee.cust_id = ?", custId).
		Joins("left join mst.m_emp_group eg on eg.emp_grp_id = employee.emp_grp_id AND eg.cust_id = ?", parentCustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.return.outlet_id AND ot.cust_id = ?", custId).
		Joins("left join sys.m_user printer on printer.user_id = sls.return.printed_by").
		Where("sls.return.cust_id=? AND sls.return.return_no = ? ", custId, returnNo).
		Take(&rtn).Error

	return rtn, err
}

func (repository *RepositoryReturnImpl) FindReturnDetail(returnNo string, custId string, parentCustId string) (Details []model.ReturnDetailRead, err error) {
	queryReturnedProduct := `left join (
			select sls.return_det.order_detail_id,
				coalesce(sum(sls.return_det.qty1), 0) as returned_qty1, 
				coalesce(sum(sls.return_det.qty2), 0) as returned_qty2, 
				coalesce(sum(sls.return_det.qty3), 0) as returned_qty3
			from sls.return_det
			where sls.return_det.return_no = '` + returnNo + `' and sls.return_det.cust_id = '` + custId + `'
			group by sls.return_det.order_detail_id
		) returned_product on returned_product.order_detail_id = invoice_detail.order_detail_id`

	err = repository.Select(`
			sls.return_det.return_detail_id, 
			sls.return_det.order_detail_id, 
			sls.return_det.return_no, 
			sls.return_det.product_id, 
			product.pro_code as product_code, 
			product.pro_name as product_name, 
			sls.return_det.wh_id, 
			wh.wh_code, 
			wh.wh_name, 
			sls.return_det.item_type, 
			sls.return_det.item_cnd,
			sls.return_det.qty, 
			sls.return_det.qty1, 
			sls.return_det.qty2, 
			sls.return_det.qty3, 
			invoice_detail.qty1 as invoice_qty1, 
			invoice_detail.qty2 as invoice_qty2, 
			invoice_detail.qty3 as invoice_qty3,
			sls.return_det.sell_price1, 
			sls.return_det.sell_price2, 
			sls.return_det.sell_price3, 
			sls.return_det.unit_id1, 
			sls.return_det.unit_id2,
			sls.return_det.unit_id3,
			unit1.unit_name as unit_name1, 
			unit2.unit_name as unit_name2, 
			unit3.unit_name as unit_name3, 
			sls.return_det.conv_unit2, 
			sls.return_det.conv_unit3, 
			sls.return_det.vat, 
			sls.return_det.vat_value, 
			sls.return_det.disc_value, 
			sls.return_det.promo_value, 
			sls.return_det.sub_total, 
			sls.return_det.total, 
			sls.return_det.return_reason_id, 
			return_reason.return_reason_code, 
			return_reason.return_reason_name,
			product.volume,
			product.weight,
			product.volume1,product.volume2,product.volume3,product.weight1,product.weight2,product.weight3,
			coalesce(invoice_detail.qty1 - coalesce(returned_product.returned_qty1, 0), 0) as remaining_qty1,
			coalesce(invoice_detail.qty2 - coalesce(returned_product.returned_qty2, 0), 0) as remaining_qty2,
			coalesce(invoice_detail.qty3 - coalesce(returned_product.returned_qty3, 0), 0) as remaining_qty3
		`).
		Joins("left join sls.order_detail invoice_detail on invoice_detail.order_detail_id = sls.return_det.order_detail_id AND invoice_detail.cust_id = ?", custId).
		Joins("left join mst.m_product product on product.pro_id = sls.return_det.product_id AND product.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.return_det.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_return_reason return_reason on return_reason.return_reason_id = sls.return_det.return_reason_id AND return_reason.cust_id IN (?, ?)", custId, parentCustId).
		Joins("left join mst.m_unit unit1 on unit1.unit_id = sls.return_det.unit_id1 AND unit1.cust_id IN (?, ?)", custId, parentCustId).
		Joins("left join mst.m_unit unit2 on unit2.unit_id = sls.return_det.unit_id2 AND unit2.cust_id IN (?, ?)", custId, parentCustId).
		Joins("left join mst.m_unit unit3 on unit3.unit_id = sls.return_det.unit_id3 AND unit3.cust_id IN (?, ?)", custId, parentCustId).
		Joins(queryReturnedProduct).
		Where("sls.return_det.cust_id=? AND sls.return_det.return_no = ? ", custId, returnNo).
		Find(&Details).Error

	return Details, err
}

func (repository *RepositoryReturnImpl) CountReturnedProductQty(invoiceNo string, productId int64, custId string) (qtySummary model.ReturnedDetailRead, err error) {
	err = repository.Select(`
			coalesce(sum(sls.return_det.qty1), 0) as returned_qty1, 
			coalesce(sum(sls.return_det.qty2), 0) as returned_qty2, 
			coalesce(sum(sls.return_det.qty3), 0) as returned_qty3
		`).
		Where("sls.return_det.return_no in (select sls.return.return_no from sls.return where sls.return.invoice_no = ? and sls.return.cust_id = ? and sls.return.data_status IN (3, 4, 5, 6))", invoiceNo, custId).
		Where("sls.return_det.product_id = ?", productId).
		Where("sls.return_det.cust_id = ?", custId).
		Find(&qtySummary).Error

	return qtySummary, err
}

func (repository *RepositoryReturnImpl) FindAllByCustId(dataFilter entity.ReturnQueryFilter) ([]model.ReturnList, int64, int, error) {
	var rtn []model.ReturnList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("return_no")
	query := repository.Select(`
			sls.return.refference_no, 
			sls.return.return_no, 
			sls.return.return_date, 
			sls.return.invoice_no, 
			sls.return.invoice_date, 
			sls.return.salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name, 
			sls.return.outlet_id, 
			ot.outlet_code, 
			ot.outlet_name,
			ot.address1 as outlet_address,
			ot.latitude as outlet_latitude,
			ot.longitude as outlet_longitude,
			sls.return.data_status, 
			sls.return.created_by, 
			creator.user_fullname AS created_by_name,
			sls.return.created_at, 
			sls.return.reviewed_by, 
			reviewer.user_fullname AS reviewed_by_name,
			sls.return.reviewed_at
		`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.return.salesman_id AND employee.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.return.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user creator on creator.user_id = sls.return.created_by AND creator.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user reviewer on reviewer.user_id = sls.return.reviewed_by AND reviewer.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.return.cust_id=?", dataFilter.CustId)
	query.Where("sls.return.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.return.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.return.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.ReturnDateFrom != nil && dataFilter.ReturnDateTo != nil {
		query.Where("sls.return.return_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.ReturnDateFrom), str.UnixTimestampToUtcTime(*dataFilter.ReturnDateTo))
		queryCount.Where("sls.return.return_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.ReturnDateFrom), str.UnixTimestampToUtcTime(*dataFilter.ReturnDateTo))
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.return.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.return.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.return.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.return.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.return.data_status in ?", dataFilter.Status)
		query.Where("sls.return.data_status in ?", dataFilter.Status)
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.return.return_no LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("sls.return.return_no LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("return_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&rtn).Error
	if err != nil {
		return rtn, total, 0, err
	}
	err = queryCount.Model(&rtn).Count(&total).Error
	if err != nil {
		return rtn, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return rtn, total, lastPage, nil
}

/*
	func (repository *RepositoryReturnImpl) Delete(c context.Context, custId string, returnNo string, deletedBy int64) error {
		var data model.Return
		result := repository.model(c).Model(&data).Where("cust_id = ? AND return_no=? AND is_del= ? ", custId, returnNo, false).
			Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("no rows affected")
		}
		return nil
	}
*/
func (repository *RepositoryReturnImpl) Update(c context.Context, data *model.Return) error {
	// result := repository.model(c).Model(&data).Updates(data)
	custId := data.CustID
	data.CustID = ""
	returnNo := data.ReturnNo
	data.ReturnNo = ""
	// log.Println("Data Cust ID : ", data.CustID)
	result := repository.model(c).Model(&data).Where("return_no=?", returnNo).Where("cust_id = ?", custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryReturnImpl) DeleteDetailNotInIDs(c context.Context, returnNo string, IDs []int64, custId string) error {
	var Details model.ReturnDetail
	err := repository.model(c).Where("return_no=?", returnNo).Where("cust_id=?", custId).Where("return_detail_id not in (?) ", IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryReturnImpl) UpdateDetail(c context.Context, Details *model.ReturnDetail) error {
	// result := repository.model(c).Updates(&Details)
	// result := repository.model(c).Model(&Details).Updates(Details)
	result := repository.model(c).Select("Qty1", "Qty2", "Qty3", "PromoValue", "DiscValue", "Vat", "VatValue", "SubTotal", "Total", "ItemCnd", "WhId", "ReturnReasonID").Updates(Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryReturnImpl) UpdateQuantity(c context.Context, Details *model.ReturnQuantity) error {
	// result := repository.model(c).Updates(&Details)
	result := repository.model(c).Select("Qty1", "Qty2", "Qty3", "Vat", "VatValue", "SubTotal", "Total").Updates(&Details)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryReturnImpl) FindAllReturnStatusesLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ReturnStatusesFilter, int64, int, error) {

	var returnStatuses []model.ReturnStatusesFilter

	var total int64

	queryCount := repository.Select("sls.return.data_status")
	query := repository.Select(`sls.return.data_status as return_status`)

	queryCount.Where("sls.return.cust_id=?", dataFilter.CustId)
	query.Where("sls.return.cust_id=?", dataFilter.CustId)

	queryCount.Where("sls.return.data_status IS NOT NULL")
	query.Where("sls.return.data_status IS NOT NULL")

	queryCount.Group("sls.return.data_status")
	query.Group("sls.return.data_status")

	query.Order("sls.return.data_status ASC")

	err := query.Find(&returnStatuses).Error
	if err != nil {
		return returnStatuses, total, 0, err
	}

	total = int64(len(returnStatuses))
	lastPage := 1
	return returnStatuses, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllSalesmanFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.SalesmansFilter, int64, int, error) {

	var salesmans []model.SalesmansFilter

	var total int64

	queryCount := repository.Select("mst.m_employee.emp_id").
		Joins("inner join mst.m_employee on mst.m_employee.emp_id = sls.return.salesman_id AND mst.m_employee.cust_id = ?", dataFilter.CustId)
	query := repository.Select(`mst.m_employee.emp_id as salesman_id, mst.m_employee.emp_code as salesman_code, mst.m_employee.emp_name as salesman_name`).
		Joins("inner join mst.m_employee on mst.m_employee.emp_id = sls.return.salesman_id AND mst.m_employee.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.return.cust_id=?", dataFilter.CustId)
	query.Where("sls.return.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_employee.is_active=?", true)
	query.Where("mst.m_employee.is_active=?", true)

	queryCount.Where("mst.m_employee.is_del=?", false)
	query.Where("mst.m_employee.is_del=?", false)

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	queryCount.Group("mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name")
	query.Group("mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				if colSort[0] == "salesman_id" {
					colSort[0] = "emp_id"
				}
				if colSort[0] == "salesman_code" {
					colSort[0] = "emp_code"
				}
				if colSort[0] == "salesman_name" {
					colSort[0] = "emp_name"
				}
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_employee."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_employee.emp_id DESC")
	}

	err := query.Find(&salesmans).Error
	if err != nil {
		return salesmans, total, 0, err
	}

	total = int64(len(salesmans))
	lastPage := 1
	return salesmans, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllEmployeeByEmpGrpIdFilterByCustIdLookupMode(dataFilter entity.SalesmanQueryFilter) ([]model.Employee, int64, int, error) {

	var salesmans []model.Employee

	var total int64

	queryCount := repository.Select("mst.m_employee.emp_id")
	query := repository.Select(`mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name, mst.m_employee.emp_grp_id`)

	// queryCount.Where("sls.return.cust_id=?", dataFilter.CustId)
	// query.Where("sls.return.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_employee.is_active=?", true)
	query.Where("mst.m_employee.is_active=?", true)

	queryCount.Where("mst.m_employee.is_del=?", false)
	query.Where("mst.m_employee.is_del=?", false)

	queryCount.Where("mst.m_employee.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_employee.cust_id=?", dataFilter.CustId)

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.EmpGrpId) > 0 {
		queryCount.Where("mst.m_employee.emp_grp_id = ?", dataFilter.EmpGrpId)
		query.Where("mst.m_employee.emp_grp_id = ?", dataFilter.EmpGrpId)
	}

	queryCount.Group("mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name, mst.m_employee.emp_grp_id")
	query.Group("mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name, mst.m_employee.emp_grp_id")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_employee."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_employee.emp_id DESC")
	}

	err := query.Find(&salesmans).Error
	if err != nil {
		return salesmans, total, 0, err
	}

	total = int64(len(salesmans))
	lastPage := 1
	return salesmans, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllRolesFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.EmpGroupsFilter, int64, int, error) {

	var empGroups []model.EmpGroupsFilter

	var total int64

	queryCount := repository.Select("mst.m_emp_group.emp_grp_id")
	query := repository.Select("mst.m_emp_group.emp_grp_id, mst.m_emp_group.emp_grp_code, mst.m_emp_group.emp_grp_name")

	queryCount.Where("mst.m_emp_group.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_emp_group.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("(mst.m_emp_group.emp_grp_name ILIKE ? OR mst.m_emp_group.emp_grp_name ILIKE ?)", "driver", "salesman")
	query.Where("(mst.m_emp_group.emp_grp_name ILIKE ? OR mst.m_emp_group.emp_grp_name ILIKE ?)", "driver", "salesman")

	queryCount.Where("mst.m_emp_group.is_active=?", true)
	query.Where("mst.m_emp_group.is_active=?", true)

	queryCount.Where("mst.m_emp_group.is_del=?", false)
	query.Where("mst.m_emp_group.is_del=?", false)

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_emp_group.emp_grp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_emp_group.emp_grp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	queryCount.Group("mst.m_emp_group.emp_grp_id, mst.m_emp_group.emp_grp_code, mst.m_emp_group.emp_grp_name")
	query.Group("mst.m_emp_group.emp_grp_id, mst.m_emp_group.emp_grp_code, mst.m_emp_group.emp_grp_name")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_emp_group."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_emp_group.emp_grp_id DESC")
	}

	err := query.Find(&empGroups).Error
	if err != nil {
		return empGroups, total, 0, err
	}

	total = int64(len(empGroups))
	lastPage := 1
	return empGroups, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllOutletFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletsFilter, int64, int, error) {

	var outlets []model.OutletsFilter

	var total int64

	queryCount := repository.Select("mst.m_outlet.outlet_id").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.return.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId)
	query := repository.Select("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.return.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.return.cust_id=?", dataFilter.CustId)
	query.Where("sls.return.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_outlet.is_active=?", true)
	query.Where("mst.m_outlet.is_active=?", true)

	queryCount.Where("mst.m_outlet.is_del=?", false)
	query.Where("mst.m_outlet.is_del=?", false)

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	queryCount.Group("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name")
	query.Group("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_outlet."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_outlet.outlet_id DESC")
	}

	err := query.Find(&outlets).Error
	if err != nil {
		return outlets, total, 0, err
	}

	total = int64(len(outlets))
	lastPage := 1
	return outlets, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllSalesmanFilterByCustIdLookupModeCreate(dataFilter entity.GeneralQueryFilter) ([]model.SalesmansFilterCreate, int64, int, error) {

	var salesmans []model.SalesmansFilterCreate

	var total int64

	queryCount := repository.Select("mst.m_employee.emp_id").
		Joins("inner join mst.m_employee on mst.m_employee.emp_id = sls.order.salesman_id AND mst.m_employee.cust_id = ?", dataFilter.CustId)
	query := repository.Select(`mst.m_employee.emp_id as salesman_id, mst.m_employee.emp_code as salesman_code, mst.m_employee.emp_name as salesman_name`).
		Joins("inner join mst.m_employee on mst.m_employee.emp_id = sls.order.salesman_id AND mst.m_employee.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_employee.is_active=?", true)
	query.Where("mst.m_employee.is_active=?", true)

	queryCount.Where("mst.m_employee.is_del=?", false)
	query.Where("mst.m_employee.is_del=?", false)

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	queryCount.Group("mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name")
	query.Group("mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				if colSort[0] == "salesman_id" {
					colSort[0] = "emp_id"
				}
				if colSort[0] == "salesman_code" {
					colSort[0] = "emp_code"
				}
				if colSort[0] == "salesman_name" {
					colSort[0] = "emp_name"
				}
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_employee."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_employee.emp_id DESC")
	}

	err := query.Find(&salesmans).Error
	if err != nil {
		return salesmans, total, 0, err
	}

	total = int64(len(salesmans))
	lastPage := 1
	return salesmans, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllOutletFilterByCustIdLookupModeCreate(dataFilter entity.OutletCreateReturnQueryFilter) ([]model.OutletsFilterCreate, int64, int, error) {

	var outlets []model.OutletsFilterCreate

	var total int64

	queryCount := repository.Select("mst.m_outlet.outlet_id").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId)
	query := repository.Select("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_outlet.is_active=?", true)
	query.Where("mst.m_outlet.is_active=?", true)

	queryCount.Where("mst.m_outlet.is_del=?", false)
	query.Where("mst.m_outlet.is_del=?", false)

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("mst.m_outlet.outlet_id in ?", dataFilter.OutletID)
		query.Where("mst.m_outlet.outlet_id in ?", dataFilter.OutletID)
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	queryCount.Group("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name")
	query.Group("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_outlet."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_outlet.outlet_id DESC")
	}

	err := query.Find(&outlets).Error
	if err != nil {
		return outlets, total, 0, err
	}

	total = int64(len(outlets))
	lastPage := 1
	return outlets, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllProductByCustIdOld(dataFilter entity.ProductListQueryFilter) ([]model.ProductList, int64, int, error) {
	var products []model.ProductList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("sls.order_detail.order_detail_id").
		Joins("inner join sls.order invoice on invoice.ro_no = sls.order_detail.ro_no AND invoice.data_status = 4 AND invoice.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_product product on product.pro_id = sls.order_detail.pro_id")
	query := repository.Select(`
			sls.order_detail.order_detail_id, 
			invoice.invoice_no, 
			invoice.invoice_date, 
			sls.order_detail.pro_id as product_id, 
			product.pro_code as product_code, 
			product.pro_name as product_name, 
			sls.order_detail.qty1 as invoice_qty1, 
			sls.order_detail.qty2 as invoice_qty2, 
			sls.order_detail.qty3 as invoice_qty3,
			sls.order_detail.sell_price1 as sell_price1, 
			sls.order_detail.sell_price2 as sell_price2, 
			sls.order_detail.sell_price3 as sell_price3, 
			sls.order_detail.unit_id1, 
			sls.order_detail.unit_id2, 
			sls.order_detail.unit_id3, 
			unit1.unit_name as unit_name1, 
			unit2.unit_name as unit_name2,
			unit3.unit_name as unit_name3,
			sls.order_detail.conv_unit2,
			sls.order_detail.conv_unit3,
			sls.order_detail.vat
		`).
		Joins("inner join sls.order invoice on invoice.ro_no = sls.order_detail.ro_no AND invoice.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_product product on product.pro_id = sls.order_detail.pro_id AND product.cust_id = ?", dataFilter.ParentCustId).
		Joins("left join mst.m_unit unit1 on unit1.unit_id = sls.order_detail.unit_id1 AND unit1.cust_id = ?", dataFilter.ParentCustId).
		Joins("left join mst.m_unit unit2 on unit2.unit_id = sls.order_detail.unit_id2 AND unit2.cust_id = ?", dataFilter.ParentCustId).
		Joins("left join mst.m_unit unit3 on unit3.unit_id = sls.order_detail.unit_id3 AND unit3.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("sls.order_detail.cust_id=?", dataFilter.CustId)
	query.Where("sls.order_detail.cust_id=?", dataFilter.CustId)

	queryCount.Where("invoice.data_status = 4")
	query.Where("invoice.data_status = 4")

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("invoice.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("invoice.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("invoice.outlet_id in ?", dataFilter.OutletID)
		query.Where("invoice.outlet_id in ?", dataFilter.OutletID)
	}

	if dataFilter.Query != "" && dataFilter.SearchBy != "" {
		if dataFilter.SearchBy == "product" {
			queryCount.Where("product.pro_name ILIKE ? OR product.pro_code = ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
			query.Where("product.pro_name ILIKE ? OR product.pro_code = ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		} else {
			queryCount.Where("invoice.invoice_no = ?", "%"+dataFilter.Query+"%")
			query.Where("invoice.invoice_no = ?", "%"+dataFilter.Query+"%")
		}
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
		query.Order("invoice_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return products, total, 0, err
	}
	err = queryCount.Model(&products).Count(&total).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return products, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllProductByCustId(dataFilter entity.ProductListQueryFilter) ([]model.ProductList, int64, int, error) {
	var products []model.ProductList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryInvoice := "inner join sls.order invoice on invoice.ro_no = sls.order_detail.ro_no AND invoice.cust_id = '" + dataFilter.CustId + "' AND invoice.data_status IN (4, 5, 6, 7) AND invoice.invoice_no IS NOT NULL"

	if len(dataFilter.SalesmanId) > 0 {
		queryInvoice = queryInvoice + " AND invoice.salesman_id in (" + strings.Trim(strings.Join(strings.Split(fmt.Sprint(dataFilter.SalesmanId), " "), ","), "[]") + ")"
	}

	if len(dataFilter.OutletID) > 0 {
		queryInvoice = queryInvoice + " AND invoice.outlet_id in (" + strings.Trim(strings.Join(strings.Split(fmt.Sprint(dataFilter.OutletID), " "), ","), "[]") + ")"
	}

	queryProduct := "inner join mst.m_product product on product.pro_id = sls.order_detail.pro_id AND product.cust_id IN ('" + dataFilter.ParentCustId + "', '" + dataFilter.CustId + "')"

	if dataFilter.SearchBy == "product" {
		queryProduct = queryProduct + " AND (product.pro_name ILIKE '%" + dataFilter.Query + "%' OR product.pro_code = '" + dataFilter.Query + "')"
	}

	if dataFilter.SearchBy == "invoice" {
		queryInvoice = queryInvoice + " AND invoice.invoice_no = '" + dataFilter.Query + "'"
	}

	OneMonthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	queryReturnedProduct := `left join (
			select sls.return_det.order_detail_id,
				coalesce(sum(sls.return_det.qty1), 0) as returned_qty1, 
				coalesce(sum(sls.return_det.qty2), 0) as returned_qty2, 
				coalesce(sum(sls.return_det.qty3), 0) as returned_qty3
			from sls.return_det
			where sls.return_det.return_no in (
				select sls.return.return_no
				from sls.return
				where sls.return.cust_id = '` + dataFilter.CustId + `' 
					and sls.return.data_status IN (3, 4, 5, 6)
					and sls.return.invoice_date >= '` + OneMonthAgo + `'
			)
			and sls.return_det.order_detail_id is not null
			group by sls.return_det.order_detail_id
		) returned_product on returned_product.order_detail_id = sls.order_detail.order_detail_id`

	queryCount := repository.Select("sls.order_detail.order_detail_id").
		Joins(queryInvoice).
		Joins(queryProduct).
		Joins("left join mst.m_unit unit1 on unit1.unit_id = sls.order_detail.unit_id1 AND unit1.cust_id IN (?, ?)", dataFilter.ParentCustId, dataFilter.CustId).
		Joins("left join mst.m_unit unit2 on unit2.unit_id = sls.order_detail.unit_id2 AND unit2.cust_id IN (?, ?)", dataFilter.ParentCustId, dataFilter.CustId).
		Joins("left join mst.m_unit unit3 on unit3.unit_id = sls.order_detail.unit_id3 AND unit3.cust_id IN (?, ?)", dataFilter.ParentCustId, dataFilter.CustId).
		Joins(queryReturnedProduct)

	query := repository.Select(`
			sls.order_detail.order_detail_id, 
			invoice.invoice_no, 
			invoice.invoice_date, 
			invoice.wh_id, 
			warehouse.wh_code, 
			warehouse.wh_name, 
			sls.order_detail.pro_id as product_id, 
			product.pro_code as product_code, 
			product.pro_name as product_name, 
			sls.order_detail.qty1 as invoice_qty1, 
			sls.order_detail.qty2 as invoice_qty2, 
			sls.order_detail.qty3 as invoice_qty3,
			sls.order_detail.sell_price1, 
			sls.order_detail.sell_price2, 
			sls.order_detail.sell_price3, 
			sls.order_detail.unit_id1, 
			sls.order_detail.unit_id2, 
			sls.order_detail.unit_id3, 
			unit1.unit_name as unit_name1, 
			unit2.unit_name as unit_name2,
			unit3.unit_name as unit_name3,
			sls.order_detail.conv_unit2,
			sls.order_detail.conv_unit3,
			sls.order_detail.vat,
			coalesce(returned_product.returned_qty1, 0) as returned_qty1,
			coalesce(returned_product.returned_qty2, 0) as returned_qty2,
			coalesce(returned_product.returned_qty3, 0) as returned_qty3,
			(sls.order_detail.qty1 - coalesce(returned_product.returned_qty1, 0)) as remaining_qty1,
			(sls.order_detail.qty2 - coalesce(returned_product.returned_qty2, 0)) as remaining_qty2,
			(sls.order_detail.qty3 - coalesce(returned_product.returned_qty3, 0)) as remaining_qty3
		`).
		Joins(queryInvoice).
		Joins(queryProduct).
		Joins("left join mst.m_warehouse warehouse on warehouse.wh_id = invoice.wh_id AND warehouse.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_unit unit1 on unit1.unit_id = sls.order_detail.unit_id1 AND unit1.cust_id IN (?, ?)", dataFilter.ParentCustId, dataFilter.CustId).
		Joins("left join mst.m_unit unit2 on unit2.unit_id = sls.order_detail.unit_id2 AND unit2.cust_id IN (?, ?)", dataFilter.ParentCustId, dataFilter.CustId).
		Joins("left join mst.m_unit unit3 on unit3.unit_id = sls.order_detail.unit_id3 AND unit3.cust_id IN (?, ?)", dataFilter.ParentCustId, dataFilter.CustId).
		Joins(queryReturnedProduct)

	queryCount.Where("sls.order_detail.cust_id=?", dataFilter.CustId)
	query.Where("sls.order_detail.cust_id=?", dataFilter.CustId)

	queryCount.Where("sls.order_detail.item_type=1")
	query.Where("sls.order_detail.item_type=1")

	queryCount.Where("(sls.order_detail.qty1 - coalesce(returned_product.returned_qty1, 0)) > 0 OR (sls.order_detail.qty2 - coalesce(returned_product.returned_qty2, 0)) > 0 OR (sls.order_detail.qty3 - coalesce(returned_product.returned_qty3, 0)) > 0")
	query.Where("(sls.order_detail.qty1 - coalesce(returned_product.returned_qty1, 0)) > 0 OR (sls.order_detail.qty2 - coalesce(returned_product.returned_qty2, 0)) > 0 OR (sls.order_detail.qty3 - coalesce(returned_product.returned_qty3, 0)) > 0")

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

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return products, total, 0, err
	}
	err = queryCount.Model(&products).Count(&total).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return products, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllMasterReturnReasonLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ReturnReasonLookup, int64, int, error) {

	var return_reasons []model.ReturnReasonLookup
	var total int64

	query := repository.Select("mst.m_return_reason.return_reason_id, mst.m_return_reason.return_reason_name")

	query.Where("mst.m_return_reason.cust_id=?", dataFilter.ParentCustId)

	query.Where("mst.m_return_reason.is_active=?", true)

	query.Where("mst.m_return_reason.is_del=?", false)

	if dataFilter.Query != "" {
		query.Where("mst.m_return_reason.return_reason_name LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("return_reason_id DESC")
	}

	err := query.Find(&return_reasons).Error
	if err != nil {
		return return_reasons, total, 0, err
	}

	total = int64(len(return_reasons))
	lastPage := 1
	return return_reasons, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindAllMasterWarehouseLookupMode(dataFilter entity.WarehouseQueryFilter) ([]model.WarehouseLookup, int64, int, error) {

	var warehouses []model.WarehouseLookup
	var total int64

	query := repository.Select("mst.m_warehouse.wh_id, mst.m_warehouse.wh_code, mst.m_warehouse.wh_name")

	query.Where("mst.m_warehouse.cust_id=?", dataFilter.CustId)

	query.Where("mst.m_warehouse.is_active=?", true)

	query.Where("mst.m_warehouse.is_del=?", false)

	// if dataFilter.WhId != 0 {
	// 	query.Where("mst.m_warehouse.wh_id  ?", dataFilter.WhId)
	// }
	if len(dataFilter.WhId) > 0 {
		query.Where("mst.m_warehouse.wh_id in ?", dataFilter.WhId)
	}
	query.Where("mst.m_warehouse.wh_id <> 0")

	if dataFilter.StockType != "" {
		query.Where("mst.m_warehouse.stock_type = ?", dataFilter.StockType)
	}

	if dataFilter.Query != "" {
		query.Where("mst.m_warehouse.wh_name LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("wh_id DESC")
	}

	err := query.Find(&warehouses).Error
	if err != nil {
		return warehouses, total, 0, err
	}

	total = int64(len(warehouses))
	lastPage := 1
	return warehouses, total, lastPage, nil
}

func productsLookupSelectColumns() string {
	return `
			mst.m_product.pro_id as product_id, 
			mst.m_product.pro_code as product_code, 
			mst.m_product.pro_name as product_name, 
			mst.m_product.sell_price1, 
			mst.m_product.sell_price2, 
			mst.m_product.sell_price3, 
			mst.m_product.unit_id1, 
			mst.m_product.unit_id2, 
			mst.m_product.unit_id3, 
			mst.m_product.conv_unit2, 
			mst.m_product.conv_unit3, 
			mst.m_product.vat, 
			unit1.unit_name as unit_name1, 
			unit2.unit_name as unit_name2, 
			unit3.unit_name as unit_name3
		`
}

func (repository *RepositoryReturnImpl) productsLookupJoinUnits(query *gorm.DB, unitMasterCustID string) *gorm.DB {
	return query.
		Joins("left join mst.m_unit unit1 on unit1.unit_id = mst.m_product.unit_id1 AND unit1.cust_id = ?", unitMasterCustID).
		Joins("left join mst.m_unit unit2 on unit2.unit_id = mst.m_product.unit_id2 AND unit2.cust_id = ?", unitMasterCustID).
		Joins("left join mst.m_unit unit3 on unit3.unit_id = mst.m_product.unit_id3 AND unit3.cust_id = ?", unitMasterCustID)
}

func sortProductsLookupMerged(products []model.ProductsLookup, sortSpec string) {
	if len(products) == 0 {
		return
	}
	parseFirst := func() (field string, desc bool, ok bool) {
		if strings.TrimSpace(sortSpec) == "" {
			return "", true, false
		}
		parts := strings.Split(sortSpec, ",")
		if len(parts) == 0 {
			return "", true, false
		}
		segs := strings.Split(strings.TrimSpace(parts[0]), ":")
		if len(segs) < 2 {
			return "", true, false
		}
		field = strings.ToLower(strings.TrimSpace(strings.Replace(segs[0], "product", "pro", -1)))
		dir := strings.ToLower(strings.TrimSpace(segs[1]))
		return field, dir == "desc", true
	}
	field, desc, ok := parseFirst()
	if ok && field == "pro_id" {
		sort.Slice(products, func(i, j int) bool {
			a := int64(0)
			b := int64(0)
			if products[i].ProductId != nil {
				a = *products[i].ProductId
			}
			if products[j].ProductId != nil {
				b = *products[j].ProductId
			}
			if desc {
				return a > b
			}
			return a < b
		})
		return
	}
	sort.Slice(products, func(i, j int) bool {
		a := int64(0)
		b := int64(0)
		if products[i].ProductId != nil {
			a = *products[i].ProductId
		}
		if products[j].ProductId != nil {
			b = *products[j].ProductId
		}
		return a > b
	})
}

func shouldUseCustMasterOnly(sortSpec string) bool {
	if strings.TrimSpace(sortSpec) == "" {
		return false
	}
	parts := strings.Split(sortSpec, ",")
	if len(parts) == 0 {
		return false
	}
	segs := strings.Split(strings.TrimSpace(parts[0]), ":")
	if len(segs) < 2 {
		return false
	}
	field := strings.ToLower(strings.TrimSpace(strings.Replace(segs[0], "product", "pro", -1)))
	return field == "pro_id"
}

func (repository *RepositoryReturnImpl) FindAllMasterProductByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ProductsLookup, int64, int, error) {
	var total int64
	parentID := dataFilter.ParentCustId
	custID := dataFilter.CustId

	var invoicedProIDs []int64
	err := repository.Raw(`
		SELECT DISTINCT od.pro_id
		FROM sls.order_detail od
		INNER JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id AND o.cust_id = ?
		WHERE od.item_type = 1
			AND o.data_status IN (4, 5, 6, 7)
			AND o.invoice_no IS NOT NULL
	`, custID).Scan(&invoicedProIDs).Error
	if err != nil {
		return nil, 0, 0, err
	}

	byID := make(map[int64]model.ProductsLookup)

	appendFromQuery := func(q *gorm.DB) error {
		var batch []model.ProductsLookup
		if err := q.Find(&batch).Error; err != nil {
			return err
		}
		for _, row := range batch {
			if row.ProductId == nil {
				continue
			}
			byID[*row.ProductId] = row
		}
		return nil
	}

	useCustMasterOnly := shouldUseCustMasterOnly(dataFilter.Sort)

	if !useCustMasterOnly && len(invoicedProIDs) > 0 {
		q := repository.productsLookupJoinUnits(
			repository.Select(productsLookupSelectColumns()).Table("mst.m_product"),
			parentID,
		).
			Where("mst.m_product.cust_id = ?", parentID).
			Where("mst.m_product.is_active = ?", true).
			Where("mst.m_product.is_del = ?", false).
			Where("mst.m_product.pro_id IN ?", invoicedProIDs)
		if dataFilter.Query != "" {
			q = q.Where("mst.m_product.pro_name ILIKE ? OR mst.m_product.pro_code = ?", "%"+dataFilter.Query+"%", dataFilter.Query)
		}
		if err := appendFromQuery(q); err != nil {
			return nil, 0, 0, err
		}
	}

	if custID != "" {
		unitJoinCustID := parentID
		if unitJoinCustID == "" {
			unitJoinCustID = custID
		}
		q := repository.productsLookupJoinUnits(
			repository.Select(productsLookupSelectColumns()).Table("mst.m_product"),
			unitJoinCustID,
		).
			Where("mst.m_product.cust_id = ?", custID).
			Where("mst.m_product.is_active = ?", true).
			Where("mst.m_product.is_del = ?", false)
		if dataFilter.Query != "" {
			q = q.Where("mst.m_product.pro_name ILIKE ? OR mst.m_product.pro_code = ?", "%"+dataFilter.Query+"%", dataFilter.Query)
		}
		if err := appendFromQuery(q); err != nil {
			return nil, 0, 0, err
		}
	}

	products := make([]model.ProductsLookup, 0, len(byID))
	for _, row := range byID {
		products = append(products, row)
	}
	sortProductsLookupMerged(products, dataFilter.Sort)

	total = int64(len(products))
	limit := dataFilter.Limit
	if limit <= 0 {
		limit = len(products)
		if limit == 0 {
			limit = 10
		}
	}
	page := dataFilter.Page
	if page <= 0 {
		page = 1
	}
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}
	start := (page - 1) * limit
	if start >= len(products) {
		return []model.ProductsLookup{}, total, lastPage, nil
	}
	end := start + limit
	if end > len(products) {
		end = len(products)
	}
	return products[start:end], total, lastPage, nil
}

func (repository *RepositoryReturnImpl) FindSalesmanById(salesmanId int, custId string, parentCustId string) (rtn []model.SalesmanRead, err error) {
	// err = repository.Select(`
	// 		mst.m_salesman.emp_id as salesman_id,
	// 		mst.m_salesman.sales_name as salesman_name,
	// 		mst.m_salesman.wh_id
	// 	`).
	// 	Where("mst.m_salesman.emp_id=?", salesmanId).
	// 	Take(&rtn).Error

	// return rtn, err
	query := repository.Model(&model.SalesmanRead{}).Select(`
			mst.m_salesman.emp_id as salesman_id, 
			mst.m_salesman.sales_name as salesman_name, 
			mst.m_salesman.wh_id 
		`).
		Where("mst.m_salesman.emp_id=?", salesmanId).
		Where("mst.m_salesman.is_active=?", true)

	queryBranch := repository.Model(&model.SalesmanCanvasRead{}).Select(`
			mst.m_salesman_canvas.emp_id as salesman_id, 
			mst.m_salesman.sales_name as salesman_name, 
			mst.m_salesman_canvas.wh_id 
		`).
		Joins("inner join mst.m_salesman on mst.m_salesman.emp_id = mst.m_salesman_canvas.emp_id").
		Where("mst.m_salesman_canvas.emp_id=?", salesmanId).
		Where("mst.m_salesman_canvas.is_active=?", true)

	queryUnion := repository.Raw("SELECT * FROM (? UNION ALL ?) AS combined", query, queryBranch)
	if err := queryUnion.Scan(&rtn).Error; err != nil {
		return rtn, err
	}

	return rtn, nil
}

func (repository *RepositoryReturnImpl) Print(c context.Context, custId string, returnNo string, printedBy int64) error {
	var data model.Return
	result := repository.model(c).Model(&data).Where("return_no=? AND cust_id = ? AND is_printed= ? ", returnNo, custId, false).
		Updates(map[string]interface{}{"is_printed": true, "printed_by": printedBy, "printed_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
