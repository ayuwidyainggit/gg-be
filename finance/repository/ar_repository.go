package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryArImpl struct {
		*gorm.DB
	}
)

type ArRepository interface {
	// Store(c context.Context, data *model.Ar) error
	// StoreDetail(c context.Context, data *model.ArDet) error
	FindByInvoiceNo(invoiceNo string, custId string) (ar model.ArRead, err error)
	FindDetail(invoiceNo string, custId string, parentCustId string) (Details []model.ArPaymentRead, err error)
	FindAllByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.ArList, int64, int, error)
	CountInvoicePaidAmount(invoiceNo string, custId string) (invoice model.InvoicePaidAmount, err error)
	FindLastApprovedDeposit(invoiceNo string, custId string) (invoice model.LastApprovedDeposit, err error)
	// Delete(c context.Context, custId string, arNo string, deletedBy int64) error
	// Update(c context.Context, arNo string, data model.Ar) error
	// DeleteDetailNotInIDs(c context.Context, arNo string, IDs []int64) error
	// UpdateDetail(c context.Context, Details *model.ArDet) error
	// AR COLLECTION
	StoreCollection(c context.Context, data *model.Collection) error
	StoreCollectionDetail(c context.Context, data *model.CollectionDet) error
	FindCollectionByNo(collectionNo string, custId string, parentCustId string) (collection model.CollectionList, err error)
	FindCollectionDetail(collectionNo string, custId string) (Details []model.CollectionDetList, err error)
	FindAllCollectionByCustId(dataFilter entity.CollectionQueryFilter) ([]model.CollectionList, int64, int, error)
	DeleteCollection(c context.Context, custId string, collectionNo string, deletedBy int64) error
	UpdateCollection(c context.Context, collectionNo string, data model.Collection) error
	UpdateCollectionRemainingAmount(c context.Context, collectionNo string, custId string, remainingAmount float64) error
	DeleteCollectionDetailNotInIDs(c context.Context, collectionNo string, IDs []int64, custId string) error
	DeleteAllCollectionDetails(c context.Context, collectionNo string, custId string) error
	UpdateCollectionDetail(c context.Context, Details *model.CollectionDet) error
	PrintCollection(c context.Context, custId string, collectionNo string, printedBy int64) error

	FindAllEmployeeGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (employeeGroup []model.EmployeeGroup, total int64, lastPage int, err error)
	FindAllEmployeeByCustIdLookupMode(dataFilter entity.EmployeeListQueryFilter) (employee []model.Employee, total int64, lastPage int, err error)
	FindAllInvoiceByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, int64, int, error)
	FindAllCollectorByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (collector []model.Collector, total int64, lastPage int, err error)
	FindAllOutletGroupFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outletGroup []model.OutletGroupFilter, total int64, lastPage int, err error)
	FindAllOutletFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outlet []model.OutletFilter, total int64, lastPage int, err error)
	FindAllSalesmanFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (salesman []model.SalesmanFilter, total int64, lastPage int, err error)
}

func NewArRepo(db *gorm.DB) *RepositoryArImpl {
	return &RepositoryArImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryArImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

/*
	func (repository *RepositoryArImpl) Store(c context.Context, data *model.Ar) error {
		err := repository.model(c).Create(data).Error
		if err != nil {
			return err
		}
		return nil
	}

	func (repository *RepositoryArImpl) StoreDetail(c context.Context, data *model.ArDet) error {
		err := repository.model(c).Create(data).Error
		if err != nil {
			return err
		}
		return nil
	}
*/
func (repository *RepositoryArImpl) FindByInvoiceNo(invoiceNo string, custId string) (ar model.ArRead, err error) {
	err = repository.Select(`
		sls.order.*, 
		sls.order.total as invoice_amount, 
		mst.m_employee.emp_id as salesman_id,
		mst.m_employee.emp_code as salesman_code,
		mst.m_employee.emp_name as salesman_name,
		mst.m_outlet.outlet_id as outlet_id,
		mst.m_outlet.outlet_code as outlet_code,
		mst.m_outlet.outlet_name as outlet_name
		`).
		Joins("left join mst.m_employee on mst.m_employee.emp_id = sls.order.salesman_id AND mst.m_employee.cust_id = ?", custId).
		Joins("left join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", custId).
		Where("sls.order.invoice_no = ? AND sls.order.cust_id=?", invoiceNo, custId).
		Take(&ar).Error
	return ar, err
}

func (repository *RepositoryArImpl) FindDetail(invoiceNo string, custId string, parentCustId string) (Details []model.ArPaymentRead, err error) {
	err = repository.Select(`
		acf.deposit_payment.deposit_payment_id, 
		deposit.deposit_date as visit_date, 
		deposit.deposit_date, 
		deposit.deposit_no, 
		employee.emp_id,
		employee.emp_code,
		employee.emp_name,
		emp_group.emp_grp_id,
		emp_group.emp_grp_code,
		emp_group.emp_grp_name,
		deposit.deposit_status as verification_status, 
		deposit_detail.total_payment, 
		deposit_detail.remaining_payment, 
		acf.deposit_payment.pay_type as payment_method, 
		acf.deposit_payment.payment_amount as amount,
		deposit.approved_by as verified_by, 
		approver.user_fullname as verified_by_name, 
		deposit.approved_at as verified_date, 
		acf.deposit_payment.invoice_no as reason,
		acf.deposit_payment.document_no as additional_info
		`).
		Joins("left join acf.deposit_detail deposit_detail on deposit_detail.deposit_no = acf.deposit_payment.deposit_no AND deposit_detail.invoice_no = acf.deposit_payment.invoice_no AND deposit_detail.cust_id = ?", custId).
		Joins("left join acf.deposit deposit on deposit.deposit_no = acf.deposit_payment.deposit_no AND deposit.cust_id = ?", custId).
		Joins("left join mst.m_employee employee on employee.emp_id = deposit.emp_id AND employee.cust_id = ?", custId).
		Joins("left join mst.m_emp_group emp_group on emp_group.emp_grp_id = employee.emp_grp_id AND emp_group.cust_id = ?", parentCustId).
		Joins("left join sys.m_user approver on approver.user_id = deposit.approved_by").
		Where("acf.deposit_payment.invoice_no = ? AND acf.deposit_payment.cust_id = ?", invoiceNo, custId).
		Order("acf.deposit_payment.deposit_payment_id DESC").
		Find(&Details).Error

	return Details, err
}
func (repository *RepositoryArImpl) FindAllByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.ArList, int64, int, error) {
	var ro []model.ArList
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryPaidInvoices := `left join (
			select acf.deposit_detail.invoice_no,
			coalesce(sum(acf.deposit_detail.total_payment), 0) as paid_amount
		from acf.deposit_detail
		inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = '` + dataFilter.CustId + `' AND deposit.deposit_status IN (1, 2) 
		where acf.deposit_detail.cust_id = '` + dataFilter.CustId + `'
		group by acf.deposit_detail.invoice_no
	) paid_invoices on paid_invoices.invoice_no = sls.order.invoice_no`

	queryCount := repository.Select("sls.order.invoice_no")
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
		Joins("left join mst.m_employee employee on employee.emp_id = sls.order.salesman_id AND employee.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins(queryPaidInvoices)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId).Joins(queryPaidInvoices)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.InvoiceFrom != nil && dataFilter.InvoiceTo != nil {
		query.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
		queryCount.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.order.invoice_no=?", dataFilter.Query)
		query.Where("sls.order.invoice_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if dataFilter.InvoiceStatus != nil {
		switch entity.InvoiceStatus(*dataFilter.InvoiceStatus) {
		case entity.InvoiceStatusPaid:
			queryCount.Where("(sls.order.total - COALESCE(paid_invoices.paid_amount, 0)) <= 0")
			query.Where("(sls.order.total - COALESCE(paid_invoices.paid_amount, 0)) <= 0")
		case entity.InvoiceStatusOutstanding:
			queryCount.Where("(sls.order.total - COALESCE(paid_invoices.paid_amount, 0)) > 0")
			query.Where("(sls.order.total - COALESCE(paid_invoices.paid_amount, 0)) > 0")
		default:
			queryCount.Where("1=0")
			query.Where("1=0")
		}
	}

	// queryCount.Where("sls.order.data_status = 6")
	// query.Where("sls.order.data_status = 6")

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

	if dataFilter.Mode != "print" {
		page := dataFilter.Page
		if page-1 < 1 {
			page = 1
		}
		offset := (page - 1) * dataFilter.Limit

		query.Limit(limit).Offset(offset)
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

func (repository *RepositoryArImpl) CountInvoicePaidAmount(invoiceNo string, custId string) (invoice model.InvoicePaidAmount, err error) {
	err = repository.Select(`
			coalesce(sum(acf.deposit_detail.total_payment), 0) as paid_amount
		`).
		Joins("inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ? AND deposit.deposit_status IN (1, 2)", custId).
		Where("acf.deposit_detail.invoice_no = ?", invoiceNo).
		Where("acf.deposit_detail.cust_id = ?", custId).
		Take(&invoice).Error

	return invoice, err
}

func (repository *RepositoryArImpl) FindLastApprovedDeposit(invoiceNo string, custId string) (deposit model.LastApprovedDeposit, err error) {
	err = repository.Select(`
			acf.deposit_detail.invoice_no,
			deposit.deposit_no,
			deposit.approved_at
		`).
		Joins("inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ? AND deposit.deposit_status = 2", custId).
		Where("acf.deposit_detail.invoice_no = ?", invoiceNo).
		Where("acf.deposit_detail.cust_id = ?", custId).
		Order("deposit.approved_at DESC").
		Take(&deposit).Error

	return deposit, err
}

/*
	func (repository *RepositoryArImpl) Delete(c context.Context, custId string, arNo string, deletedBy int64) error {
		var data model.Ar
		result := repository.model(c).Model(&data).Where("ar_no=? AND cust_id = ? AND is_del= ? ", arNo, custId, false).
			Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("no rows affected")
		}
		return nil
	}

func (repository *RepositoryArImpl) Update(c context.Context, arNo string, data model.Ar) error {

		result := repository.model(c).Model(&data).Where("ar_no=?", arNo).Updates(data)
		if result.Error != nil {
			return result.Error
		}

		return nil
	}

	func (repository *RepositoryArImpl) DeleteDetailNotInIDs(c context.Context, arNo string, IDs []int64) error {
		var Details model.ArDet
		err := repository.model(c).Where("ar_no=? AND ar_det_id not in (?) ", arNo, IDs).Delete(&Details).Error
		return err
	}

	func (repository *RepositoryArImpl) UpdateDetail(c context.Context, Details *model.ArDet) error {
		result := repository.model(c).Updates(&Details)
		if result.Error != nil {
			return result.Error
		}
		// if result.RowsAffected == 0 {
		// 	return errors.New("no rows affected")
		// }
		return nil
	}
*/
func (repository *RepositoryArImpl) StoreCollection(c context.Context, data *model.Collection) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryArImpl) StoreCollectionDetail(c context.Context, data *model.CollectionDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryArImpl) FindCollectionByNo(collectionNo string, custId string, parentCustId string) (collection model.CollectionList, err error) {
	err = repository.Select("acf.collection.*, us1.user_fullname AS created_by_name, us2.user_fullname AS updated_by_name, us3.user_fullname AS deleted_by_name, us4.user_fullname AS printed_by_name, emp.emp_id, emp.emp_code, emp.emp_name, ot_grp.ot_grp_id, ot_grp.ot_grp_code, ot_grp.ot_grp_name").
		Joins("left join sys.m_user us1 on us1.user_id = acf.collection.created_by").
		Joins("left join sys.m_user us2 on us2.user_id = acf.collection.updated_by").
		Joins("left join sys.m_user us3 on us3.user_id = acf.collection.deleted_by").
		Joins("left join sys.m_user us4 on us4.user_id = acf.collection.printed_by").
		Joins("left join mst.m_employee emp on emp.emp_id = acf.collection.emp_id AND emp.cust_id = ?", custId).
		Joins("left join mst.m_outlet_group ot_grp on ot_grp.ot_grp_id = acf.collection.ot_grp_id AND ot_grp.cust_id = ?", parentCustId).
		Where("acf.collection.collection_no = ? AND acf.collection.cust_id=?", collectionNo, custId).
		Take(&collection).Error
	return collection, err
}

func (repository *RepositoryArImpl) FindCollectionDetail(collectionNo string, custId string) (Details []model.CollectionDetList, err error) {
	err = repository.Select(`
			acf.collection_det.*, 
			acf.collection_det.paid_by_invoice AS total_invoice_amount,
			us1.user_fullname AS created_by_name,
			us1.user_fullname AS created_by_name,
			invoice.ro_no as sales_order,
			invoice.invoice_date,
			invoice.due_date,
			salesman.emp_id as salesman_id,
			salesman.emp_code as salesman_code,
			salesman.emp_name as salesman_name,
			outlet.outlet_id,
			outlet.outlet_code,
			outlet.outlet_name
		`).
		Joins("left join sys.m_user us1 on us1.user_id = acf.collection_det.created_by").
		Joins("left join sls.order invoice on invoice.invoice_no = acf.collection_det.invoice_no and invoice.cust_id=?", custId).
		Joins("left join mst.m_employee salesman on salesman.emp_id = invoice.salesman_id").
		Joins("left join mst.m_outlet outlet on outlet.outlet_id = invoice.outlet_id").
		Where("acf.collection_det.collection_no = ? AND acf.collection_det.cust_id=?", collectionNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryArImpl) FindAllCollectionByCustId(dataFilter entity.CollectionQueryFilter) ([]model.CollectionList, int64, int, error) {
	var collection []model.CollectionList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("collection_no")
	query := repository.Select("acf.collection.*, us.user_fullname AS updated_by_name,emp.emp_code,emp.emp_name, emp_grp.emp_grp_id, emp_grp.emp_grp_code, emp_grp.emp_grp_name").
		Joins("left join sys.m_user us on us.user_id = acf.collection.updated_by").
		Joins("left join mst.m_employee emp on emp.emp_id = acf.collection.emp_id AND emp.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_emp_group emp_grp on emp_grp.emp_grp_id = emp.emp_grp_id AND emp_grp.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.collection.cust_id=?", dataFilter.CustId)
	query.Where("acf.collection.cust_id=?", dataFilter.CustId)

	if dataFilter.Mode == "deposit" {
		queryCount.Where("NOT EXISTS (SELECT 1 FROM acf.deposit d WHERE d.cust_id = ? AND acf.collection.collection_no = d.collection_no)", dataFilter.CustId)
		query.Where("NOT EXISTS (SELECT 1 FROM acf.deposit d WHERE d.cust_id = ? AND acf.collection.collection_no = d.collection_no)", dataFilter.CustId)
	}

	if dataFilter.CollectionDateFrom != nil && dataFilter.CollectionDateTo != nil {
		query.Where("acf.collection.collection_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.CollectionDateFrom), str.UnixTimestampToUtcTime(*dataFilter.CollectionDateTo))
		queryCount.Where("acf.collection.collection_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.CollectionDateFrom), str.UnixTimestampToUtcTime(*dataFilter.CollectionDateTo))
	}

	if dataFilter.Query != "" {
		q := fmt.Sprintf("%%%s%%", strings.ToLower(dataFilter.Query))
		query.Where("LOWER(acf.collection.collection_no) LIKE ?", q)
		queryCount.Where("LOWER(acf.collection.collection_no) LIKE ?", q)
	}

	if len(dataFilter.EmpId) > 0 {
		queryCount.Where("acf.collection.emp_id in ?", dataFilter.EmpId)
		query.Where("acf.collection.emp_id in ?", dataFilter.EmpId)
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
		query.Order("acf.collection.collection_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&collection).Error
	if err != nil {
		return collection, total, 0, err
	}
	err = queryCount.Model(&collection).Count(&total).Error
	if err != nil {
		return collection, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return collection, total, lastPage, nil
}

func (repository *RepositoryArImpl) DeleteCollection(c context.Context, custId string, collectionNo string, deletedBy int64) error {
	var data model.Collection
	result := repository.model(c).Model(&data).Where("collection_no=? AND cust_id = ? AND is_del= ? ", collectionNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryArImpl) UpdateCollection(c context.Context, collectionNo string, data model.Collection) error {

	custId := data.CustID
	data.CustID = ""
	data.CollectionNo = ""
	result := repository.model(c).Model(&data).Where("collection_no=?", collectionNo).Where("cust_id=?", custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryArImpl) UpdateCollectionRemainingAmount(c context.Context, collectionNo string, custId string, remainingAmount float64) error {
	result := repository.model(c).Model(&model.Collection{}).
		Where("collection_no=?", collectionNo).
		Where("cust_id=?", custId).
		Update("remaining_amount", remainingAmount)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (repository *RepositoryArImpl) DeleteCollectionDetailNotInIDs(c context.Context, collectionNo string, IDs []int64, custId string) error {
	var Details model.CollectionDet
	err := repository.model(c).Where("collection_no=?", collectionNo).Where("collection_det_id not in (?)", IDs).Where("cust_id=?", custId).Delete(&Details).Error
	return err
}

func (repository *RepositoryArImpl) DeleteAllCollectionDetails(c context.Context, collectionNo string, custId string) error {
	var Details model.CollectionDet
	err := repository.model(c).Where("collection_no=?", collectionNo).Where("cust_id=?", custId).Delete(&Details).Error
	return err
}

func (repository *RepositoryArImpl) UpdateCollectionDetail(c context.Context, Details *model.CollectionDet) error {
	if Details.CollectionDetID == nil {
		return errors.New("collection detail id is required")
	}

	result := repository.model(c).
		Model(&model.CollectionDet{}).
		Where("collection_det_id = ?", *Details.CollectionDetID).
		Where("cust_id = ?", Details.CustID).
		Updates(map[string]interface{}{
			"invoice_no":       Details.InvoiceNo,
			"salesman_id":      Details.SalesmanID,
			"invoice_amount":   Details.InvoiceAmount,
			"remaining_amount": Details.RemainingAmount,
			"paid_amount":      Details.PaidAmount,
		})
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryArImpl) PrintCollection(c context.Context, custId string, collectionNo string, printedBy int64) error {
	var data model.Collection
	result := repository.model(c).Model(&data).Where("collection_no=? AND cust_id = ? AND is_printed= ? ", collectionNo, custId, false).
		Updates(map[string]interface{}{"is_printed": true, "printed_by": printedBy, "printed_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArImpl) FindAllEmployeeGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.EmployeeGroup, int64, int, error) {

	var empGroups []model.EmployeeGroup

	var total int64

	queryCount := repository.Select("emp_grp_id")
	query := repository.Select(`mst.m_emp_group.emp_grp_id, mst.m_emp_group.emp_grp_code, mst.m_emp_group.emp_grp_name`)

	queryCount.Where("mst.m_emp_group.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_emp_group.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_emp_group.is_active=?", true)
	query.Where("mst.m_emp_group.is_active=?", true)

	queryCount.Where("mst.m_emp_group.is_del=?", false)
	query.Where("mst.m_emp_group.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_emp_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_emp_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_emp_group.emp_grp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_emp_group.emp_grp_name LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("emp_grp_id DESC")
	}

	err := query.Find(&empGroups).Error
	if err != nil {
		return empGroups, total, 0, err
	}
	err = queryCount.Model(&empGroups).Count(&total).Error
	if err != nil {
		return empGroups, total, 0, err
	}

	lastPage := 1
	return empGroups, total, lastPage, nil
}

/*
func (repository *RepositoryArImpl) FindAllOutletGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletGroup, int64, int, error) {

		var outletGroups []model.OutletGroup

		var total int64

		queryCount := repository.Select("ot_grp_id")
		query := repository.Select(`mst.m_outlet_group.ot_grp_id, mst.m_outlet_group.ot_grp_code, mst.m_outlet_group.ot_grp_name`)

		queryCount.Where("mst.m_outlet_group.cust_id=?", dataFilter.CustId)
		query.Where("mst.m_outlet_group.cust_id=?", dataFilter.CustId)

		queryCount.Where("mst.m_outlet_group.is_active=?", true)
		query.Where("mst.m_outlet_group.is_active=?", true)

		queryCount.Where("mst.m_outlet_group.is_del=?", false)
		query.Where("mst.m_outlet_group.is_del=?", false)

		if dataFilter.From != nil && dataFilter.To != nil {
			query.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
			queryCount.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		}

		if dataFilter.Query != "" {
			queryCount.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
			query.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
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
			query.Order("ot_grp_id DESC")
		}

		err := query.Find(&outletGroups).Error
		if err != nil {
			return outletGroups, total, 0, err
		}
		err = queryCount.Model(&outletGroups).Count(&total).Error
		if err != nil {
			return outletGroups, total, 0, err
		}

		lastPage := 1
		return outletGroups, total, lastPage, nil
	}
*/
func (repository *RepositoryArImpl) FindAllEmployeeByCustIdLookupMode(dataFilter entity.EmployeeListQueryFilter) ([]model.Employee, int64, int, error) {

	var employees []model.Employee

	var total int64

	queryCount := repository.Select("emp_id")
	query := repository.Select(`mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name`)

	queryCount.Where("mst.m_employee.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_employee.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_employee.is_active=?", true)
	query.Where("mst.m_employee.is_active=?", true)

	queryCount.Where("mst.m_employee.is_del=?", false)
	query.Where("mst.m_employee.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_employee.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_employee.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.EmpGrpID) > 0 {
		queryCount.Where("mst.m_employee.emp_grp_id in ?", dataFilter.EmpGrpID)
		query.Where("mst.m_employee.emp_grp_id in ?", dataFilter.EmpGrpID)
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_employee.emp_name ILIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_employee.emp_name ILIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("emp_id DESC")
	}

	err := query.Find(&employees).Error
	if err != nil {
		return employees, total, 0, err
	}
	err = queryCount.Model(&employees).Count(&total).Error
	if err != nil {
		return employees, total, 0, err
	}

	lastPage := 1
	return employees, total, lastPage, nil
}

/*
func (repository *RepositoryArImpl) FindAllSalesmanByCustIdLookupMode(dataFilter entity.SalesmanListQueryFilter) ([]model.Salesman, int64, int, error) {

		var salesmans []model.Salesman

		var total int64

		queryCount := repository.Select("mst.m_salesman.emp_id")
		query := repository.Select(
			`mst.m_salesman.emp_id as salesman_id,
			employee.emp_code as salesman_code,
			employee.emp_name as salesman_name`).
			Joins("left join mst.m_employee employee on employee.emp_id = mst.m_salesman.emp_id AND employee.cust_id = ?", dataFilter.CustId)

		queryCount.Where("mst.m_salesman.cust_id=?", dataFilter.CustId)
		query.Where("mst.m_salesman.cust_id=?", dataFilter.CustId)

		queryCount.Where("mst.m_salesman.is_active=?", true)
		query.Where("mst.m_salesman.is_active=?", true)

		queryCount.Where("mst.m_salesman.is_del=?", false)
		query.Where("mst.m_salesman.is_del=?", false)

		if dataFilter.From != nil && dataFilter.To != nil {
			query.Where("mst.m_salesman.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
			queryCount.Where("mst.m_salesman.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		}

		if dataFilter.Query != "" {
			queryCount.Where("mst.m_salesman.sales_name LIKE ?", "%"+dataFilter.Query+"%")
			query.Where("mst.m_salesman.sales_name LIKE ?", "%"+dataFilter.Query+"%")
		}

		sortBy := ``
		if dataFilter.Sort != "" {
			mSortBy := strings.Split(dataFilter.Sort, ",")
			for _, row := range mSortBy {
				colSort := strings.Split(row, ":")
				if len(colSort) > 1 {
					if colSort[0] == "salesman_id" {
						colSort[0] = "mst.m_salesman.emp_id"
					}
					if colSort[0] == "salesman_code" {
						colSort[0] = "employee.emp_code"
					}
					if colSort[0] == "salesman_name" {
						colSort[0] = "mst.m_salesman.sales_name"
					}
					sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
				}
			}
			sortBy = strings.TrimSuffix(sortBy, ", ")
			query.Order(sortBy)
		} else {
			query.Order("mst.m_salesman.emp_id DESC")
		}

		err := query.Find(&salesmans).Error
		if err != nil {
			return salesmans, total, 0, err
		}
		err = queryCount.Model(&salesmans).Count(&total).Error
		if err != nil {
			return salesmans, total, 0, err
		}

		lastPage := 1
		return salesmans, total, lastPage, nil
	}
*/
func (repository *RepositoryArImpl) FindAllInvoiceByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, int64, int, error) {
	var ro []model.InvoiceList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryPaidInvoices := `left join (
			select acf.deposit_detail.invoice_no,
			coalesce(sum(acf.deposit_detail.total_payment), 0) as paid_amount
		from acf.deposit_detail
		inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = '` + dataFilter.CustId + `' AND deposit.deposit_status IN (1, 2) 
		where acf.deposit_detail.cust_id = '` + dataFilter.CustId + `'
		group by acf.deposit_detail.invoice_no
	) paid_invoices on paid_invoices.invoice_no = sls.order.invoice_no`

	queryCount := repository.Select("sls.order.invoice_no").
		Joins("left join mst.m_employee employee on employee.emp_id = sls.order.salesman_id AND employee.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet_group ot_grp on ot_grp.ot_grp_id = ot.ot_grp_id AND ot_grp.cust_id = ?", dataFilter.ParentCustId).
		Joins(queryPaidInvoices)

	query := repository.Select(
		`sls.order.invoice_no, 
			sls.order.invoice_date, 
			sls.order.due_date, 
			sls.order.outlet_id, 
			sls.order.salesman_id, 
			sls.order.ro_no,
			sls.order.total as invoice_amount,
			ot.outlet_code, ot.outlet_name, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name,
			coalesce(paid_invoices.paid_amount, 0) as paid_amount,
			(sls.order.total - coalesce(paid_invoices.paid_amount, 0)) as remaining_amount
		`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.order.salesman_id AND employee.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet_group ot_grp on ot_grp.ot_grp_id = ot.ot_grp_id AND ot_grp.cust_id = ?", dataFilter.ParentCustId).
		Joins(queryPaidInvoices)

	queryCount.Where("sls.order.invoice_no IS NOT NULL")
	query.Where("sls.order.invoice_no IS NOT NULL")

	queryCount.Where("(sls.order.total - coalesce(paid_invoices.paid_amount, 0)) > 0")
	query.Where("(sls.order.total - coalesce(paid_invoices.paid_amount, 0)) > 0")

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.InvoiceFrom != nil && dataFilter.InvoiceTo != nil {
		query.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
		queryCount.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	}

	if dataFilter.DueFrom != nil && dataFilter.DueTo != nil {
		query.Where("sls.order.due_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DueFrom), str.UnixTimestampToUtcTime(*dataFilter.DueTo))
		queryCount.Where("sls.order.due_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DueFrom), str.UnixTimestampToUtcTime(*dataFilter.DueTo))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.order.invoice_no=?", dataFilter.Query)
		query.Where("sls.order.invoice_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletGroupID) > 0 {
		queryCount.Where("ot_grp.ot_grp_id in ?", dataFilter.OutletGroupID)
		query.Where("ot_grp.ot_grp_id in ?", dataFilter.OutletGroupID)
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
		query.Order("sls.order.ro_no DESC")
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

/*
func (repository *RepositoryArImpl) FindAllOutletByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.Outlet, int64, int, error) {

		var outlets []model.Outlet

		var total int64

		queryCount := repository.Select("outlet_id")
		query := repository.Select(`mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name`)

		queryCount.Where("mst.m_outlet.cust_id=?", dataFilter.CustId)
		query.Where("mst.m_outlet.cust_id=?", dataFilter.CustId)

		queryCount.Where("mst.m_outlet.is_active=?", true)
		query.Where("mst.m_outlet.is_active=?", true)

		queryCount.Where("mst.m_outlet.is_del=?", false)
		query.Where("mst.m_outlet.is_del=?", false)

		if dataFilter.From != nil && dataFilter.To != nil {
			query.Where("mst.m_outlet.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
			queryCount.Where("mst.m_outlet.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		}

		if dataFilter.Query != "" {
			queryCount.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
			query.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
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
			query.Order("outlet_id DESC")
		}

		err := query.Find(&outlets).Error
		if err != nil {
			return outlets, total, 0, err
		}
		err = queryCount.Model(&outlets).Count(&total).Error
		if err != nil {
			return outlets, total, 0, err
		}

		lastPage := 1
		return outlets, total, lastPage, nil
	}
*/
func (repository *RepositoryArImpl) FindAllCollectorByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.Collector, int64, int, error) {
	var collectors []model.Collector
	var total int64
	query, queryCount, limit, page := repository.buildCollectorLookupQueries(dataFilter)
	offset := (page - 1) * limit

	err := query.Order(repository.buildCollectorLookupSort(dataFilter.Sort)).Limit(limit).Offset(offset).Scan(&collectors).Error
	if err != nil {
		return collectors, total, 0, err
	}

	err = queryCount.Count(&total).Error
	if err != nil {
		return collectors, total, 0, err
	}

	lastPage := 1
	if total > 0 {
		lastPage = int(math.Ceil(float64(total) / float64(limit)))
	}
	return collectors, total, lastPage, nil
}

func (repository *RepositoryArImpl) buildCollectorLookupQueries(dataFilter entity.GeneralQueryFilter) (*gorm.DB, *gorm.DB, int, int) {
	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 9999 {
		limit = 9999
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}

	selectQuery := `
		d.emp_id AS emp_id,
		COALESCE(NULLIF(emp.emp_code, ''), CAST(d.emp_id AS varchar)) AS emp_code,
		COALESCE(NULLIF(emp.emp_name, ''), COALESCE(NULLIF(emp.emp_code, ''), CAST(d.emp_id AS varchar))) AS emp_name
	`

	query := repository.Table("acf.deposit d").
		Select(selectQuery).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = d.emp_id AND emp.cust_id = ?", dataFilter.CustId).
		Where("d.cust_id = ?", dataFilter.CustId).
		Where("d.deleted_at IS NULL").
		Where("d.emp_id IS NOT NULL")

	queryCount := repository.Table("acf.deposit d").
		Distinct("d.emp_id").
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = d.emp_id AND emp.cust_id = ?", dataFilter.CustId).
		Where("d.cust_id = ?", dataFilter.CustId).
		Where("d.deleted_at IS NULL").
		Where("d.emp_id IS NOT NULL")

	if dataFilter.Query != "" {
		searchQuery := "%" + dataFilter.Query + "%"
		query = query.Where("(CAST(d.emp_id AS varchar) ILIKE ? OR COALESCE(emp.emp_code, '') ILIKE ? OR COALESCE(emp.emp_name, '') ILIKE ?)", searchQuery, searchQuery, searchQuery)
		queryCount = queryCount.Where("(CAST(d.emp_id AS varchar) ILIKE ? OR COALESCE(emp.emp_code, '') ILIKE ? OR COALESCE(emp.emp_name, '') ILIKE ?)", searchQuery, searchQuery, searchQuery)
	}

	query = query.Group("d.emp_id, emp.emp_code, emp.emp_name")

	return query, queryCount, limit, page
}

func (repository *RepositoryArImpl) buildCollectorLookupSort(sort string) string {
	defaultSort := "d.emp_id DESC"
	if strings.TrimSpace(sort) == "" {
		return defaultSort
	}

	allowedFields := map[string]string{
		"emp_id":   "d.emp_id",
		"emp_code": "emp.emp_code",
		"emp_name": "emp.emp_name",
	}

	orderClauses := make([]string, 0)
	for _, rawSort := range strings.Split(sort, ",") {
		sortToken := strings.TrimSpace(rawSort)
		if sortToken == "" {
			continue
		}

		sortParts := strings.SplitN(sortToken, ":", 2)
		if len(sortParts) != 2 {
			continue
		}

		column, ok := allowedFields[strings.TrimSpace(sortParts[0])]
		if !ok {
			continue
		}

		direction := strings.ToUpper(strings.TrimSpace(sortParts[1]))
		if direction != "ASC" && direction != "DESC" {
			continue
		}

		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", column, direction))
	}

	if len(orderClauses) == 0 {
		return defaultSort
	}

	return strings.Join(orderClauses, ", ")
}

func (repository *RepositoryArImpl) FindAllOutletGroupFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletGroupFilter, int64, int, error) {

	var outletGroups []model.OutletGroupFilter

	var total int64

	// queryCount := repository.Select("ot_grp_id")
	// query := repository.Select(`mst.m_outlet_group.ot_grp_id, mst.m_outlet_group.ot_grp_code, mst.m_outlet_group.ot_grp_name`)

	queryCount := repository.Select("mst.m_outlet_group.ot_grp_id").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId).
		Joins("inner join mst.m_outlet_group on mst.m_outlet_group.ot_grp_id = mst.m_outlet.ot_grp_id AND mst.m_outlet_group.cust_id = ?", dataFilter.ParentCustId)
	query := repository.Select("mst.m_outlet_group.ot_grp_id, mst.m_outlet_group.ot_grp_code, mst.m_outlet_group.ot_grp_name").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId).
		Joins("inner join mst.m_outlet_group on mst.m_outlet_group.ot_grp_id = mst.m_outlet.ot_grp_id AND mst.m_outlet_group.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	queryCount.Where("sls.order.data_status = 6")
	query.Where("sls.order.data_status = 6")

	queryCount.Where("mst.m_outlet_group.is_active=?", true)
	query.Where("mst.m_outlet_group.is_active=?", true)

	queryCount.Where("mst.m_outlet_group.is_del=?", false)
	query.Where("mst.m_outlet_group.is_del=?", false)

	// if dataFilter.From != nil && dataFilter.To != nil {
	// 	query.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	// 	queryCount.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	// }

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	queryCount.Group("mst.m_outlet_group.ot_grp_id, mst.m_outlet_group.ot_grp_code, mst.m_outlet_group.ot_grp_name")
	query.Group("mst.m_outlet_group.ot_grp_id, mst.m_outlet_group.ot_grp_code, mst.m_outlet_group.ot_grp_name")

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_outlet_group."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_outlet_group.ot_grp_id DESC")
	}

	err := query.Find(&outletGroups).Error
	if err != nil {
		return outletGroups, total, 0, err
	}

	total = int64(len(outletGroups))
	lastPage := 1
	return outletGroups, total, lastPage, nil
}

func (repository *RepositoryArImpl) FindAllSalesmanFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.SalesmanFilter, int64, int, error) {

	var salesmans []model.SalesmanFilter

	var total int64

	// queryCount := repository.Select("mst.m_salesman.emp_id")
	// query := repository.Select(
	// 	`mst.m_salesman.emp_id as salesman_id,
	// 	employee.emp_code as salesman_code,
	// 	employee.emp_name as salesman_name`).
	// 	Joins("left join mst.m_employee employee on employee.emp_id = mst.m_salesman.emp_id AND employee.cust_id = ?", dataFilter.CustId)

	queryCount := repository.Select("mst.m_employee.emp_id").
		Joins("inner join mst.m_employee on mst.m_employee.emp_id = sls.order.salesman_id AND mst.m_employee.cust_id = ?", dataFilter.CustId)
	query := repository.Select(`mst.m_employee.emp_id as salesman_id, mst.m_employee.emp_code as salesman_code, mst.m_employee.emp_name as salesman_name`).
		Joins("inner join mst.m_employee on mst.m_employee.emp_id = sls.order.salesman_id AND mst.m_employee.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	queryCount.Where("sls.order.data_status = 6")
	query.Where("sls.order.data_status = 6")

	queryCount.Where("mst.m_employee.is_active=?", true)
	query.Where("mst.m_employee.is_active=?", true)

	queryCount.Where("mst.m_employee.is_del=?", false)
	query.Where("mst.m_employee.is_del=?", false)

	// if dataFilter.From != nil && dataFilter.To != nil {
	// 	query.Where("mst.m_salesman.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	// 	queryCount.Where("mst.m_salesman.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	// }

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

func (repository *RepositoryArImpl) FindAllOutletFilterByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletFilter, int64, int, error) {

	var outlets []model.OutletFilter

	var total int64

	queryCount := repository.Select("mst.m_outlet.outlet_id").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId)
	query := repository.Select("mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name").
		Joins("inner join mst.m_outlet on mst.m_outlet.outlet_id = sls.order.outlet_id AND mst.m_outlet.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	queryCount.Where("sls.order.data_status = 6")
	query.Where("sls.order.data_status = 6")

	queryCount.Where("mst.m_outlet.is_active=?", true)
	query.Where("mst.m_outlet.is_active=?", true)

	queryCount.Where("mst.m_outlet.is_del=?", false)
	query.Where("mst.m_outlet.is_del=?", false)

	// if dataFilter.From != nil && dataFilter.To != nil {
	// 	query.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	// 	queryCount.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	// }

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
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
