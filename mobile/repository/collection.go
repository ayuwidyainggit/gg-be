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
	RepositoryCollectionImpl struct {
		*gorm.DB
	}
)
type CollectionRepository interface {
	// Store(c context.Context, data *model.Collection) error
	FindAllByCustId(dataFilter entity.CollectionQueryFilter) ([]model.CollectionList, model.CollectionTotal, int, error)
	FindAllByCustIdV2(dataFilter entity.CollectionQueryFilter) ([]model.CollectionList, model.CollectionTotal, int, error)
	//proses store
	CountAllByCustId(custId string, depositDate string) (int, error)
	CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error)
	StoreDetail(c context.Context, data *model.DepositDetail) (int, error)
	StorePayment(c context.Context, data *model.DepositPayment) (int, error)
	StoreDepositPaymentImage(c context.Context, data *model.DepositPaymentImage) (int, error)
	Store(c context.Context, data *model.Deposit) error
	StoreCollectionNoPayment(c context.Context, data *model.CollectionNoPayment) error

	// detail
	FindByNo(depositNo string, custId string) (whAdj model.DepositList, err error)
	FindDetailByNo(depositNo string, custId string) (whAdj []model.DepositDetailList, err error)
	FindDetailPaymentByNo(depositNo string, invoiceNo string, custId string) (whAdj []model.DepositPayment, err error)
	FindDetailPaymentInvoiceByNo(payType int, depositNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error)
	FindDetailPaymentByInvoice(payType int, depositNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error)
	FindPaymentImagesByNo(depositNo string, invoiceNo string) (whAdj []model.DepositPaymentImage, err error)

	//no payment
	FindMissedPaymentReasons(dataFilter entity.GeneralQueryFilter) ([]model.MissedPaymentReason, error)
	GetNewCollectionNo() (string, error)
	GetInvoiceNoByEmpId(custId string) ([]string, error)
	StoreCollection(c context.Context, data *model.CollectionModel) error
	StoreCollectionDetails(c context.Context, data []model.CollectionDetail) error
	GetInvoiceList(c context.Context, custID string, outletID int) ([]model.InvoiceList, error)
	GetCollectionList(c context.Context, filter entity.CollectionQueryFilter) (data []model.CollectionList, totals model.CollectionTotal, lastPage int, err error)
	CountCollectionList(c context.Context, filter entity.CollectionQueryFilter) (model.CollectionTotal, error)
}

func NewCollectionRepo(db *gorm.DB) *RepositoryCollectionImpl {
	return &RepositoryCollectionImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryCollectionImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryCollectionImpl) FindAllByCustId(dataFilter entity.CollectionQueryFilter) ([]model.CollectionList, model.CollectionTotal, int, error) {
	var ro []model.CollectionList
	var counting model.CollectionTotal
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	baseQuery := repository.Select(
		`acf.collection.*, 
		 us.user_fullname AS updated_by_name, 
		 ot.outlet_code, ot.outlet_name, 
		 sls.sales_name, 
		 cd.invoice_no`).
		Joins("left join sys.m_user us on us.user_id = acf.collection.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.collection.emp_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join acf.collection_det cd on cd.collection_no = acf.collection.collection_no").
		Joins("left join sls.order o on o.invoice_no = cd.invoice_no").
		Joins("join mst.m_outlet ot on o.outlet_id = ot.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	baseQueryCount := repository.Select(`COUNT(acf.collection.collection_no) as total, SUM(acf.collection.remaining_amount) as total_invoice`).
		Joins("left join sys.m_user us on us.user_id = acf.collection.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.collection.emp_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join acf.collection_det cd on cd.collection_no = acf.collection.collection_no").
		Joins("left join sls.order o on o.invoice_no = cd.invoice_no").
		Joins("join mst.m_outlet ot on o.outlet_id = ot.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	baseFilterQuery := "acf.collection.cust_id = ? AND acf.collection.is_del = false AND acf.collection.remaining_amount > 0"
	queryCount := baseQueryCount.Where(baseFilterQuery, dataFilter.CustId)
	query := baseQuery.Where(baseFilterQuery, dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.collection.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.collection.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	// if dataFilter.InvoiceFrom != nil && dataFilter.InvoiceTo != nil {
	// 	query.Where("acf.collection.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	// 	queryCount.Where("acf.collection.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	// }

	if dataFilter.Query != "" {
		queryCount.Where("acf.collection.collection_no=?", dataFilter.Query)
		query.Where("acf.collection.collection_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		if len(dataFilter.SalesmanId) > 1 {
			queryCount.Where("acf.collection.emp_id in (?)", dataFilter.SalesmanId)
			query.Where("acf.collection.emp_id in (?)", dataFilter.SalesmanId)
		} else {
			queryCount.Where("acf.collection.emp_id = ?", dataFilter.SalesmanId[0])
			query.Where("acf.collection.emp_id = ?", dataFilter.SalesmanId[0])
		}
	}

	if len(dataFilter.OutletID) > 0 {
		if len(dataFilter.OutletID) > 1 {
			queryCount.Where("ot.outlet_id in (?)", dataFilter.OutletID)
			query.Where("ot.outlet_id in (?)", dataFilter.OutletID)
		} else {
			queryCount.Where("ot.outlet_id = ?", dataFilter.OutletID[0])
			query.Where("ot.outlet_id = ?", dataFilter.OutletID[0])
		}
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
	// 	query.Collection(sortBy)
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
		query.Order("collection_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, model.CollectionTotal{}, 0, err
	}
	err = queryCount.Model(&ro).Scan(&counting).Error
	if err != nil {
		return ro, model.CollectionTotal{}, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(counting.Total) / float64(limit))))
	return ro, counting, lastPage, nil
}

func (repository *RepositoryCollectionImpl) CountAllByCustId(custId string, depositDate string) (int, error) {
	var deposit []model.DepositList
	var total int64

	queryCount := repository.Select("deposit_no")

	queryCount.Where("acf.deposit.cust_id = ?", custId)
	queryCount.Where("acf.deposit.deposit_date = ?", depositDate) // Menambahkan kondisi tanggal sekarang

	// queryCount.Where("acf.deposit.data_status=4")

	err := queryCount.Model(&deposit).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return int(total), nil
}

func (repository *RepositoryCollectionImpl) CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error) {
	var totalPayment float64

	err := repository.model(c).
		Table("acf.deposit_detail dd").
		Select("SUM(dd.total_payment) AS total_payment").
		Where("dd.invoice_no = ?", invoiceNo).
		Where("dd.cust_id = ?", custId).
		Group("dd.invoice_no").
		Scan(&totalPayment).Error

	if err != nil {
		// handle error
		return 0, err
	}

	return totalPayment, nil
}

func (repository *RepositoryCollectionImpl) Store(c context.Context, data *model.Deposit) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryCollectionImpl) StoreDetail(c context.Context, data *model.DepositDetail) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositDetailID, nil
}

func (repository *RepositoryCollectionImpl) StorePayment(c context.Context, data *model.DepositPayment) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositPaymentID, nil
}

func (repository *RepositoryCollectionImpl) StoreDepositPaymentImage(c context.Context, data *model.DepositPaymentImage) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositImageID, nil
}

func (repository *RepositoryCollectionImpl) StoreCollectionNoPayment(c context.Context, data *model.CollectionNoPayment) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryCollectionImpl) StoreCollection(c context.Context, data *model.CollectionModel) error {
	return repository.model(c).Create(data).Error
}

func (repository *RepositoryCollectionImpl) StoreCollectionDetails(c context.Context, data []model.CollectionDetail) error {
	if len(data) == 0 {
		return nil
	}
	return repository.model(c).Create(&data).Error
}

func (repository *RepositoryCollectionImpl) FindByNo(depositNo string, custId string) (whAdj model.DepositList, err error) {
	err = repository.Select(`
			acf.deposit.*,
			us.user_fullname AS updated_by_name,
			sls.sales_name as salesman_name,
			mpg.emp_grp_name,
			emp.emp_name,emp.emp_code,
			clc.collection_date,
			otgrp.ot_grp_id,otgrp.ot_grp_code,otgrp.ot_grp_name
		`).
		Joins("left join sys.m_user us on us.user_id = acf.deposit.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.deposit.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_emp_group mpg on mpg.emp_grp_id = acf.deposit.emp_grp_id AND mpg.cust_id = ?", custId).
		Joins("left join mst.m_employee emp on emp.emp_id = acf.deposit.emp_id AND emp.cust_id = ?", custId).
		Joins("left join acf.collection clc on clc.collection_no = acf.deposit.collection_no AND clc.cust_id = ?", custId).
		Joins("left join mst.m_outlet_group otgrp on otgrp.ot_grp_id = clc.ot_grp_id AND clc.cust_id = ?", custId).
		Where("acf.deposit.deposit_no = ? AND acf.deposit.cust_id=?", depositNo, custId).
		Where("acf.deposit.deleted_at IS NULL").
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryCollectionImpl) FindDetailByNo(depositNo string, custId string) (whAdj []model.DepositDetailList, err error) {
	err = repository.Select(`
			acf.deposit_detail.*,
			ro.invoice_date,ro.due_date,ro.order_no,ro.salesman_id,
			sls.sales_name as salesman_name,
			ot.outlet_id, ot.outlet_code, ot.outlet_name
		`).
		// Joins("left join sys.m_user us on us.user_id = acf.deposit_detail.updated_by").
		Joins("left join sls.order ro on ro.invoice_no = acf.deposit_detail.invoice_no AND ro.cust_id = ?", custId).
		Joins("left join mst.m_salesman sls on sls.emp_id = ro.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = ro.outlet_id AND ot.cust_id = ?", custId).
		Where("acf.deposit_detail.deposit_no = ? AND acf.deposit_detail.cust_id=?", depositNo, custId).
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryCollectionImpl) FindDetailPaymentByNo(depositNo string, invoiceNo string, custId string) (whAdj []model.DepositPayment, err error) {
	err = repository.Select(`
			acf.deposit_payment.*
		`).
		// Joins("left join sys.m_user us on us.user_id = acf.deposit_payment.updated_by").
		Where("acf.deposit_payment.deposit_no = ? AND acf.deposit_payment.invoice_no = ? AND acf.deposit_payment.cust_id=?", depositNo, invoiceNo, custId).
		// Where("acf.deposit_payment.is_del=false").
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryCollectionImpl) FindDetailPaymentInvoiceByNo(payType int, depositNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error) {
	err = repository.Select(`
			acf.deposit_payment.*,
			od.invoice_date,
			emp.emp_id as salesman_id,
			ot.outlet_id as outlet_id,
			emp.emp_code as salesman_code, 
			emp.emp_name as salesman_name,
			ot.outlet_code, ot.outlet_name
		`).
		Joins("LEFT JOIN acf.deposit_detail dd ON dd.deposit_no = acf.deposit_payment.deposit_no AND dd.invoice_no = acf.deposit_payment.invoice_no").
		Joins("LEFT JOIN sls.order od ON od.invoice_no = dd.invoice_no AND od.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = od.salesman_id AND emp.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = od.outlet_id AND ot.cust_id = ?", custId).
		Where("acf.deposit_payment.deposit_no = ? AND acf.deposit_payment.cust_id=?", depositNo, custId).
		Where("acf.deposit_payment.pay_type = ? AND acf.deposit_payment.cust_id=?", payType, custId).
		// Where("acf.deposit_payment.is_del=false").
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryCollectionImpl) FindDetailPaymentByInvoice(payType int, invoiceNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error) {
	err = repository.Select(`DISTINCT
			acf.deposit_payment.*,
			od.invoice_date,
			emp.emp_id as salesman_id,
			ot.outlet_id as outlet_id,
			emp.emp_code as salesman_code, 
			emp.emp_name as salesman_name,
			ot.outlet_code, ot.outlet_name
		`).
		Joins("LEFT JOIN acf.deposit_detail dd ON dd.invoice_no = acf.deposit_payment.invoice_no").
		Joins("LEFT JOIN sls.order od ON od.invoice_no = dd.invoice_no AND od.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = od.salesman_id AND emp.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = od.outlet_id AND ot.cust_id = ?", custId).
		Where("acf.deposit_payment.invoice_no = ? AND acf.deposit_payment.cust_id = ?", invoiceNo, custId).
		Where("acf.deposit_payment.pay_type = ? AND acf.deposit_payment.cust_id = ?", payType, custId).
		// Where("acf.deposit_payment.is_del=false").
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryCollectionImpl) FindPaymentImagesByNo(depositNo string, invoiceNo string) ([]model.DepositPaymentImage, error) {
	var images []model.DepositPaymentImage
	err := repository.model(context.Background()).
		Where("deposit_no = ? AND invoice_no = ? ", depositNo, invoiceNo).
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (repository *RepositoryCollectionImpl) FindMissedPaymentReasons(dataFilter entity.GeneralQueryFilter) ([]model.MissedPaymentReason, error) {
	var missedPayment []model.MissedPaymentReason

	query := repository.Select("*").Where("is_active = true AND cust_id = ?", dataFilter.CustId).Order("missed_payment_reasons_id DESC")

	err := query.Find(&missedPayment).Error
	if err != nil {
		return missedPayment, err
	}

	return missedPayment, nil

}

func (repository *RepositoryCollectionImpl) FindAllByCustIdV2(dataFilter entity.CollectionQueryFilter) ([]model.CollectionList, model.CollectionTotal, int, error) {
	var ro []model.CollectionList
	var counting model.CollectionTotal
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	baseQuery := repository.Select(
		`acf.collection.collection_no,
		acf.collection.emp_id,
		cd.invoice_no,
		cd.remaining_amount,
		cd.invoice_amount as total_amount,
		cd.paid_amount,
		o.ro_no,
		o.order_no,
		o.invoice_date as invoice_date_from,
		o.due_date as due_date_from`).
		Joins("LEFT JOIN acf.collection_det cd ON cd.collection_no = acf.collection.collection_no").
		Joins("LEFT JOIN sls.order o ON o.invoice_no = cd.invoice_no")

	baseQueryCount := repository.Select(`COUNT(cd.invoice_no) as total, SUM(cd.remaining_amount) as total_invoice`).
		Joins("LEFT JOIN acf.collection_det cd ON cd.collection_no = acf.collection.collection_no").
		Joins("LEFT JOIN sls.order o ON o.invoice_no = cd.invoice_no")

	baseFilterQuery := `acf.collection.cust_id = ? AND acf.collection.is_del = false AND cd.remaining_amount > 0 
					AND acf.collection.collection_no = (
					   SELECT MAX(c2.collection_no)
					   FROM acf.collection c2
					   JOIN acf.collection_det cd2 ON cd2.collection_no = c2.collection_no
					   WHERE cd2.invoice_no = cd.invoice_no
						 AND c2.cust_id = ?
						 AND c2.is_del = false
						 AND cd2.remaining_amount > 0
					 )
					AND NOT EXISTS (SELECT 1 FROM acf.deposit d WHERE d.collection_no = acf.collection.collection_no)`
	queryCount := baseQueryCount.Where(baseFilterQuery, dataFilter.CustId, dataFilter.CustId)
	query := baseQuery.Where(baseFilterQuery, dataFilter.CustId, dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.collection.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.collection.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.collection.collection_no=?", dataFilter.Query)
		query.Where("acf.collection.collection_no=?", dataFilter.Query)
	}

	// if len(dataFilter.SalesmanId) > 0 {
	// 	if len(dataFilter.SalesmanId) > 1 {
	// 		queryCount.Where("acf.collection.emp_id in (?)", dataFilter.SalesmanId)
	// 		query.Where("acf.collection.emp_id in (?)", dataFilter.SalesmanId)
	// 	} else {
	// 		queryCount.Where("acf.collection.emp_id = ?", dataFilter.SalesmanId[0])
	// 		query.Where("acf.collection.emp_id = ?", dataFilter.SalesmanId[0])
	// 	}
	// }

	if len(dataFilter.OutletID) > 0 {
		if len(dataFilter.OutletID) > 1 {
			queryCount.Where("o.outlet_id in (?)", dataFilter.OutletID)
			query.Where("o.outlet_id in (?)", dataFilter.OutletID)
		} else {
			queryCount.Where("o.outlet_id = ?", dataFilter.OutletID[0])
			query.Where("o.outlet_id = ?", dataFilter.OutletID[0])
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
		query.Order("collection_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, model.CollectionTotal{}, 0, err
	}
	err = queryCount.Model(&ro).Scan(&counting).Error
	if err != nil {
		return ro, model.CollectionTotal{}, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(counting.Total) / float64(limit))))
	return ro, counting, lastPage, nil
}

func (repository *RepositoryCollectionImpl) GetNewCollectionNo() (string, error) {
	var collectionNo string
	err := repository.
		Raw(`
		SELECT
			'CL' || TO_CHAR(CURRENT_DATE, 'YYMMDD') ||
			TO_CHAR(COALESCE(MAX(TO_NUMBER(SUBSTR(collection_no, 9, 4), '9999')), 0) + 1, 'FM0000') AS collection_no
		FROM
			acf.collection
		WHERE
			collection_no LIKE 'CL' || TO_CHAR(CURRENT_DATE, 'YYMMDD') || '%';
		`).Scan(&collectionNo).Error
	if err != nil {
		return "", err
	}
	return collectionNo, nil
}

func (repository *RepositoryCollectionImpl) GetInvoiceNoByEmpId(custId string) ([]string, error) {
	var invoiceNos []string
	err := repository.Table("acf.collection_det").Select("invoice_no").Where("cust_id = ?", custId).
		Where("acf.collection_det.created_at >= NOW() - INTERVAL '7 days'").
		Where("NOT EXISTS (SELECT 1 FROM acf.deposit d WHERE d.collection_no = acf.collection_det.collection_no)").
		Find(&invoiceNos).Error
	if err != nil {
		return nil, err
	}
	return invoiceNos, nil
}

func (repository *RepositoryCollectionImpl) GetInvoiceList(c context.Context, custID string, outletID int) ([]model.InvoiceList, error) {
	var data []model.InvoiceList

	err := repository.WithContext(c).
		Table("acf.collection c").Select("c.collection_no, cd.invoice_no").
		Joins("LEFT JOIN acf.collection_det cd ON cd.collection_no = c.collection_no").
		Joins("LEFT JOIN sls.order o ON o.invoice_no = cd.invoice_no").
		Where("c.is_del = false AND c.remaining_amount > 0 AND c.cust_id = ? AND o.outlet_id = ?", custID, outletID).
		Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repository *RepositoryCollectionImpl) GetCollectionList(c context.Context, filter entity.CollectionQueryFilter) (data []model.CollectionList, totals model.CollectionTotal, lastPage int, err error) {
	var limit int
	if filter.Limit <= 0 {
		limit = 10
	} else {
		limit = filter.Limit
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	baseFilter := `c.cust_id = ? AND c.is_del = false AND cd.remaining_amount > 0
		AND NOT EXISTS (
			SELECT 1 FROM acf.deposit_detail dd
			JOIN acf.deposit d ON d.deposit_no = dd.deposit_no
				AND dd.invoice_no = cd.invoice_no
				AND dd.cust_id = cd.cust_id
			WHERE d.collection_no = cd.collection_no
		)`

	baseArgs := []interface{}{filter.CustId}

	// Main query
	query := repository.WithContext(c).
		Table("acf.collection c").
		Select(`
			DISTINCT ON (cd.invoice_no)
			c.collection_no,
			c.emp_id,
			cd.invoice_no,
			cd.remaining_amount,
			cd.invoice_amount,
			cd.paid_amount,
			o.ro_no,
			o.order_no,
			o.invoice_date as invoice_date_from,
			o.due_date as due_date_from
		`).
		Joins("LEFT JOIN acf.collection_det cd ON cd.collection_no = c.collection_no").
		Joins("LEFT JOIN sls.order o ON o.invoice_no = cd.invoice_no").
		Where(baseFilter, baseArgs...)

	if len(filter.OutletID) > 0 {
		query = query.Where("o.outlet_id IN ?", filter.OutletID)
	}

	if err = query.
		Order("cd.invoice_no DESC").
		Limit(limit).
		Offset(offset).
		Scan(&data).Error; err != nil {
		return nil, totals, 0, err
	}

	totalData, err := repository.CountCollectionList(c, filter)
	if err != nil {
		return nil, totals, 0, err
	}

	lastPage = int(math.Ceil(float64(totalData.Total) / float64(limit)))
	return data, totalData, lastPage, nil
}

func (repository *RepositoryCollectionImpl) CountCollectionList(c context.Context, filter entity.CollectionQueryFilter) (model.CollectionTotal, error) {
	var result model.CollectionTotal

	q := `
	WITH total_data AS (
		SELECT DISTINCT ON (cd.invoice_no)
			c.collection_no,
			c.emp_id,
			cd.invoice_no,
			cd.remaining_amount,
			cd.invoice_amount,
			cd.paid_amount,
			o.ro_no,
			o.order_no,
			o.invoice_date,
			o.due_date
		FROM acf.collection c
		LEFT JOIN acf.collection_det cd ON cd.collection_no = c.collection_no
		LEFT JOIN sls.order o ON o.invoice_no = cd.invoice_no
		WHERE c.cust_id = ?
			AND c.is_del = false
			AND cd.remaining_amount > 0
			AND NOT EXISTS (
				SELECT 1 FROM acf.deposit_detail dd
				JOIN acf.deposit d ON d.deposit_no = dd.deposit_no
					AND dd.invoice_no = cd.invoice_no
					AND dd.cust_id = cd.cust_id
				WHERE d.collection_no = cd.collection_no
			)`

	args := []interface{}{filter.CustId}

	if len(filter.OutletID) > 0 {
		q += "\n			AND o.outlet_id IN ?"
		args = append(args, filter.OutletID)
	}

	q += `
	)
	SELECT
		COUNT(invoice_no) AS total,
		COALESCE(SUM(remaining_amount), 0) AS total_invoice
	FROM total_data`

	err := repository.WithContext(c).
		Raw(q, args...).
		Scan(&result).Error

	if err != nil {
		return model.CollectionTotal{}, err
	}

	return result, nil
}
