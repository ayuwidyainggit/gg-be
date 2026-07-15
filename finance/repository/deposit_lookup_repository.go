package repository

import (
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"gorm.io/gorm"
)

type (
	RepositoryDepositLookupImpl struct {
		*gorm.DB
	}
)

type DepositLookupRepository interface {
	FindAllCollectionNoByCustId(dataFilter entity.GeneralQueryFilter) ([]model.CollectionNoLookup, int64, int, error)
	FindAllDepositNoByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositNoLookup, int64, int, error)
	FindAllDepositStatusByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositStatusLookup, int64, int, error)

	FindInvoiceByCollectionByCustId(dataFilter entity.GeneralQueryFilter) ([]model.InvoiceCollectionList, int64, int, error)

	FindAllBalancePaymentDepositByCustId(dataFilter entity.DepositLookupQueryFilter, mode string) ([]model.DepositPaymentLookup, int64, int, error)

	FindChequeGiroBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error)
	FindBankTransferBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error)
	FindReturnBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error)
	FindCNDNBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error)
}

func NewDepositLookupRepo(db *gorm.DB) *RepositoryDepositLookupImpl {
	return &RepositoryDepositLookupImpl{db}
}

// model returns query model with context with or without transaction extracted from context
// func (repo *RepositoryDepositLookupImpl) model(ctx context.Context) *gorm.DB {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx.WithContext(ctx)
// 	}
// 	return repo.WithContext(ctx)
// }

func (repository *RepositoryDepositLookupImpl) FindAllCollectionNoByCustId(dataFilter entity.GeneralQueryFilter) ([]model.CollectionNoLookup, int64, int, error) {
	var DepositLookup []model.CollectionNoLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.deposit.collection_no))")
	query := repository.Select(
		`distinct(acf.deposit.collection_no)`)
	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)

	// Filter where 'deleted_at' is NULL or empty
	queryCount.Where("acf.deposit.deleted_at IS NULL")
	query.Where("acf.deposit.deleted_at IS NULL")

	if dataFilter.Query != "" {
		q := fmt.Sprintf("%%%s%%", strings.ToLower(dataFilter.Query))
		query.Where("LOWER(acf.deposit.collection_no) LIKE ?", q)
		queryCount.Where("LOWER(acf.deposit.collection_no) LIKE ?", q)
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
		query.Order("collection_no ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&DepositLookup).Error
	if err != nil {
		return DepositLookup, total, 0, err
	}
	err = queryCount.Model(&DepositLookup).Count(&total).Error
	if err != nil {
		return DepositLookup, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return DepositLookup, total, lastPage, nil
}

func (repository *RepositoryDepositLookupImpl) FindAllDepositNoByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositNoLookup, int64, int, error) {
	var DepositLookup []model.DepositNoLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.deposit.deposit_no))")
	query := repository.Select(
		`distinct(acf.deposit.deposit_no)`)
	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)

	// Filter where 'deleted_at' is NULL or empty
	queryCount.Where("acf.deposit.deleted_at IS NULL")
	query.Where("acf.deposit.deleted_at IS NULL")

	if dataFilter.Query != "" {
		q := fmt.Sprintf("%%%s%%", strings.ToLower(dataFilter.Query))
		query.Where("LOWER(acf.deposit.deposit_no) LIKE ?", q)
		queryCount.Where("LOWER(acf.deposit.deposit_no) LIKE ?", q)
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
		query.Order("deposit_no ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&DepositLookup).Error
	if err != nil {
		return DepositLookup, total, 0, err
	}
	err = queryCount.Model(&DepositLookup).Count(&total).Error
	if err != nil {
		return DepositLookup, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return DepositLookup, total, lastPage, nil
}

func (repository *RepositoryDepositLookupImpl) FindAllDepositStatusByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositStatusLookup, int64, int, error) {
	var DepositLookup []model.DepositStatusLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.deposit.deposit_status))")
	query := repository.Select(
		`distinct(acf.deposit.deposit_status)`)
	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)

	// Filter where 'deleted_at' is NULL or empty
	queryCount.Where("acf.deposit.deleted_at IS NULL")
	query.Where("acf.deposit.deleted_at IS NULL")

	if dataFilter.Query != "" {
		q := fmt.Sprintf("%%%s%%", strings.ToLower(dataFilter.Query))
		query.Where("LOWER(acf.deposit.deposit_status) LIKE ?", q)
		queryCount.Where("LOWER(acf.deposit.deposit_status) LIKE ?", q)
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
		query.Order("deposit_status ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&DepositLookup).Error
	if err != nil {
		return DepositLookup, total, 0, err
	}
	err = queryCount.Model(&DepositLookup).Count(&total).Error
	if err != nil {
		return DepositLookup, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return DepositLookup, total, lastPage, nil
}

func (repository *RepositoryDepositLookupImpl) FindInvoiceByCollectionByCustId(dataFilter entity.GeneralQueryFilter) ([]model.InvoiceCollectionList, int64, int, error) {
	var invoice []model.InvoiceCollectionList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ro_no").Joins("LEFT JOIN acf.collection_det cd ON cd.invoice_no = sls.order.invoice_no AND cd.cust_id = ?", dataFilter.CustId)
	query := repository.Select(
		`	cd.collection_no,
			sls.order.invoice_no, 
			sls.order.invoice_date, 
			sls.order.outlet_id, 
			sls.order.salesman_id, 
			sls.order.ro_no, 
			sls.order.total as invoice_amount, 
			(sls.order.total - coalesce(paid_invoices.paid_amount, 0)) as remaining_amount,
			ot.outlet_code, ot.outlet_name, 
			emp.emp_code as salesman_code, 
			emp.emp_name as salesman_name`).
		Joins("LEFT JOIN mst.m_salesman sales ON sales.emp_id = sls.order.salesman_id AND sales.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = sls.order.salesman_id AND emp.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN acf.collection_det cd ON cd.invoice_no = sls.order.invoice_no AND cd.cust_id = ?", dataFilter.CustId).
		Joins(`left join (
                        select acf.deposit_detail.invoice_no,
                        coalesce(sum(acf.deposit_detail.total_payment), 0) as paid_amount
                from acf.deposit_detail
                inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ? AND deposit.deposit_status IN (1, 2)
                where acf.deposit_detail.cust_id = ?
                group by acf.deposit_detail.invoice_no
        ) paid_invoices on paid_invoices.invoice_no = sls.order.invoice_no`, dataFilter.CustId, dataFilter.CustId)

	queryCount.Where("sls.order.cust_id = ?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		q := fmt.Sprintf("%%%s%%", strings.ToLower(dataFilter.Query))
		query.Where("LOWER(sls.order.ro_no) LIKE ?", q)
		queryCount.Where("LOWER(sls.order.ro_no) LIKE ?", q)
	}

	if dataFilter.CollectionNo != nil && *dataFilter.CollectionNo != "" {
		q := fmt.Sprintf("%%%s%%", strings.ToLower(*dataFilter.CollectionNo))
		query.Where("LOWER(cd.collection_no) LIKE ?", q)
		queryCount.Where("LOWER(cd.collection_no) LIKE ?", q)
	}

	queryCount.Where("sls.order.data_status=6")
	query.Where("sls.order.data_status=6")
	query.Where("remaining_amount > 0")

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

	err := query.Limit(limit).Offset(offset).Find(&invoice).Error
	if err != nil {
		return invoice, total, 0, err
	}
	err = queryCount.Model(&invoice).Count(&total).Error
	if err != nil {
		return invoice, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return invoice, total, lastPage, nil
}

func (repository *RepositoryDepositLookupImpl) FindAllBalancePaymentDepositByCustId(dataFilter entity.DepositLookupQueryFilter, mode string) ([]model.DepositPaymentLookup, int64, int, error) {
	var ChequeGiro []model.DepositPaymentLookup
	var total int64

	var docNo = "doc_no_cheque"
	var tableCheck = "acf.cheque_giro"
	var amountAs = "amount"
	var apType = 2

	if mode == "check" {
		docNo = "doc_no_cheque"
		tableCheck = "acf.cheque_giro"
		amountAs = "amount"
		apType = 2
	} else if mode == "transfer" {
		docNo = "doc_no_bank"
		tableCheck = "acf.bank_transfer"
		amountAs = "amount"
		apType = 3
	} else if mode == "cndn" {
		docNo = "cndn_no"
		tableCheck = "acf.cndn"
		amountAs = "amount"
		apType = 5
	} else if mode == "return" {
		docNo = "return_no"
		tableCheck = "sls.return"
		amountAs = "total"
		apType = 4
	}

	// Building the SQL queries
	selectCount := `COUNT(*) AS total`
	selectField := `cg.` + docNo + ` as doc_no, cg.` + amountAs + ` as amount, (cg.` + amountAs + ` - COALESCE(ab.payment,0)) as balance`
	qWhere := `WHERE cg.is_del = false AND cg.cust_id = '` + dataFilter.CustId + `'`

	if mode != "return" {
		qWhere += `AND cg.owner_id = 1`
	}

	// Adding search query filter
	if dataFilter.Query != "" {
		qWhere += ` AND cg.` + docNo + ` ILIKE '%' || '` + dataFilter.Query + `' || '%'`
	}

	if len(dataFilter.OutletId) > 0 {
		qWhere += ` AND cg.outlet_id IN (` + strings.Trim(strings.Join(strings.Split(fmt.Sprint(dataFilter.OutletId), " "), ","), "[]") + `)`
	}

	// // Add remaining_amount > 0 for cheque_giro, bank_transfer, cndn and return modes
	// if mode == "check" || mode == "transfer" || mode == "return" {
	// 	qWhere += ` AND cg.remaining_amount > 0`
	// 	// qWhere += ` AND (cg.` + amountAs + ` - COALESCE(ab.payment,0)) > 0`
	// } else if mode == "cndn" {
	// 	qWhere += ` AND cg.remaning_amount > 0`
	// }

	// Add remaining_amount > 0 for cheque_giro, bank_transfer, cndn and return modes
	if mode == "check" || mode == "transfer" || mode == "cndn" || mode == "return" {
		qWhere += ` AND (cg.` + amountAs + ` - COALESCE(ab.payment,0)) > 0`
	}

	qFrom := `
		FROM ` + tableCheck + ` cg
		LEFT JOIN (
        SELECT dp.document_no, COALESCE(SUM(dp.payment_amount), 0) AS payment
        FROM acf.deposit_payment dp
        WHERE dp.pay_type = ` + fmt.Sprintf("%d", apType) + ` AND dp.cust_id = '` + dataFilter.CustId + `'
        GROUP BY dp.document_no
    ) ab ON ab.document_no = cg.` + docNo

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// Executing the count query
	if err := repository.Raw(queryCount).Scan(&total).Error; err != nil {
		log.Error("depositLookupRepository, count total, err:", err)
		return ChequeGiro, 0, 0, err
	}

	// Handling sorting
	sortBy := `cg.` + docNo + ` DESC` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		sortArr := []string{}
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortArr = append(sortArr, fmt.Sprintf("%s %s", colSort[0], colSort[1]))
			}
		}
		if len(sortArr) > 0 {
			sortBy = strings.Join(sortArr, ", ")
		}
	}
	querySelect += ` ORDER BY ` + sortBy

	// Handling pagination
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	// Executing the select query
	if err := repository.Raw(querySelect).Scan(&ChequeGiro).Error; err != nil {
		log.Error("depositLookupRepository, FindAllBalancePaymentDepositByCustId, err:", err)
		return ChequeGiro, total, lastPage, err
	}

	return ChequeGiro, total, lastPage, nil
}

// FindChequeGiroBalance retrieves cheque giro records with balance calculation
func (repository *RepositoryDepositLookupImpl) FindChequeGiroBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error) {
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}

	activeStatuses := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}

	selectField := `cg.doc_no_cheque AS doc_no, cg.amount, (cg.amount - COALESCE(cg_used.used, 0)) AS balance`
	qFrom := `
FROM acf.cheque_giro cg
LEFT JOIN (
	SELECT dp.document_no, COALESCE(SUM(dp.payment_amount), 0) AS used
	FROM acf.deposit_payment dp
	INNER JOIN acf.deposit d ON d.deposit_no = dp.deposit_no AND d.cust_id = dp.cust_id
	WHERE dp.pay_type = 2
		AND dp.cust_id = ?
		AND d.deposit_status IN ?
	GROUP BY dp.document_no
) cg_used ON cg_used.document_no = cg.doc_no_cheque`

	qWhere := `WHERE cg.cust_id = ? AND cg.salesman_id = ? AND (cg.amount - COALESCE(cg_used.used, 0)) > 0`
	args := []interface{}{dataFilter.CustId, activeStatuses, dataFilter.CustId, dataFilter.SalesmanId}

	if dataFilter.Query != "" {
		qWhere += ` AND cg.doc_no_cheque ILIKE ?`
		args = append(args, "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.OutletId) > 0 {
		parts := make([]string, len(dataFilter.OutletId))
		for i, id := range dataFilter.OutletId {
			parts[i] = strconv.Itoa(id)
		}
		qWhere += ` AND cg.outlet_id IN (` + strings.Join(parts, ",") + `)`
	}

	queryCount := `SELECT COUNT(*) AS total ` + qFrom + ` ` + qWhere
	var total int64
	if err := repository.Raw(queryCount, args...).Scan(&total).Error; err != nil {
		log.Error("FindChequeGiroBalance, count error:", err)
		return nil, 0, 0, err
	}

	sortBy := `cg.doc_no_cheque DESC`
	if dataFilter.Sort != "" {
		sortBy = dataFilter.Sort
	}
	offset := (dataFilter.Page - 1) * dataFilter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere + ` ORDER BY ` + sortBy +
		fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	var results []model.DepositPaymentLookup
	if err := repository.Raw(querySelect, args...).Scan(&results).Error; err != nil {
		log.Error("FindChequeGiroBalance, scan error:", err)
		return results, total, 0, err
	}

	return results, total, lastPage, nil
}

// FindBankTransferBalance retrieves bank transfer records with balance calculation
func (repository *RepositoryDepositLookupImpl) FindBankTransferBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error) {
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}

	activeStatuses := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}

	selectField := `bt.doc_no_bank AS doc_no, bt.amount, (bt.amount - COALESCE(bt_used.used, 0)) AS balance`
	qFrom := `
FROM acf.bank_transfer bt
LEFT JOIN (
	SELECT dp.document_no, COALESCE(SUM(dp.payment_amount), 0) AS used
	FROM acf.deposit_payment dp
	INNER JOIN acf.deposit d ON d.deposit_no = dp.deposit_no AND d.cust_id = dp.cust_id
	WHERE dp.pay_type = 3
		AND dp.cust_id = ?
		AND d.deposit_status IN ?
	GROUP BY dp.document_no
) bt_used ON bt_used.document_no = bt.doc_no_bank`

	qWhere := `WHERE bt.cust_id = ? AND bt.salesman_id = ? AND (bt.amount - COALESCE(bt_used.used, 0)) > 0`
	args := []interface{}{dataFilter.CustId, activeStatuses, dataFilter.CustId, dataFilter.SalesmanId}

	if dataFilter.Query != "" {
		qWhere += ` AND bt.doc_no_bank ILIKE ?`
		args = append(args, "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.OutletId) > 0 {
		parts := make([]string, len(dataFilter.OutletId))
		for i, id := range dataFilter.OutletId {
			parts[i] = strconv.Itoa(id)
		}
		qWhere += ` AND bt.outlet_id IN (` + strings.Join(parts, ",") + `)`
	}

	queryCount := `SELECT COUNT(*) AS total ` + qFrom + ` ` + qWhere
	var total int64
	if err := repository.Raw(queryCount, args...).Scan(&total).Error; err != nil {
		log.Error("FindBankTransferBalance, count error:", err)
		return nil, 0, 0, err
	}

	sortBy := `bt.doc_no_bank DESC`
	if dataFilter.Sort != "" {
		sortBy = dataFilter.Sort
	}
	offset := (dataFilter.Page - 1) * dataFilter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere + ` ORDER BY ` + sortBy +
		fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	var results []model.DepositPaymentLookup
	if err := repository.Raw(querySelect, args...).Scan(&results).Error; err != nil {
		log.Error("FindBankTransferBalance, scan error:", err)
		return results, total, 0, err
	}

	return results, total, lastPage, nil
}

// FindCNDNBalance retrieves CNDN records with balance calculation
func (repository *RepositoryDepositLookupImpl) FindCNDNBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error) {
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}

	activeStatuses := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}

	selectField := `c.cndn_no AS doc_no, COALESCE(c.amount, 0) AS amount, (COALESCE(c.amount, 0) - COALESCE(c_used.used, 0)) AS balance`
	qFrom := `
FROM acf.cndn c
LEFT JOIN (
	SELECT dp.document_no, COALESCE(SUM(dp.payment_amount), 0) AS used
	FROM acf.deposit_payment dp
	INNER JOIN acf.deposit d ON d.deposit_no = dp.deposit_no AND d.cust_id = dp.cust_id
	WHERE dp.pay_type = 5
		AND dp.cust_id = ?
		AND d.deposit_status IN ?
	GROUP BY dp.document_no
) c_used ON c_used.document_no = c.cndn_no`

	qWhere := `WHERE c.cust_id = ? AND c.created_by = ? AND (COALESCE(c.amount, 0) - COALESCE(c_used.used, 0)) > 0`
	args := []interface{}{dataFilter.CustId, activeStatuses, dataFilter.CustId, dataFilter.SalesmanId}

	if dataFilter.Query != "" {
		qWhere += ` AND c.cndn_no ILIKE ?`
		args = append(args, "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.OutletId) > 0 {
		parts := make([]string, len(dataFilter.OutletId))
		for i, id := range dataFilter.OutletId {
			parts[i] = strconv.Itoa(id)
		}
		qWhere += ` AND c.outlet_id IN (` + strings.Join(parts, ",") + `)`
	}

	queryCount := `SELECT COUNT(*) AS total ` + qFrom + ` ` + qWhere
	var total int64
	if err := repository.Raw(queryCount, args...).Scan(&total).Error; err != nil {
		log.Error("FindCNDNBalance, count error:", err)
		return nil, 0, 0, err
	}

	sortBy := `c.cndn_no DESC`
	if dataFilter.Sort != "" {
		sortBy = dataFilter.Sort
	}
	offset := (dataFilter.Page - 1) * dataFilter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere + ` ORDER BY ` + sortBy +
		fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	var results []model.DepositPaymentLookup
	if err := repository.Raw(querySelect, args...).Scan(&results).Error; err != nil {
		log.Error("FindCNDNBalance, scan error:", err)
		return results, total, 0, err
	}

	return results, total, lastPage, nil
}

// FindReturnBalance retrieves return records with balance calculation
func (repository *RepositoryDepositLookupImpl) FindReturnBalance(dataFilter entity.DepositLookupQueryFilter) ([]model.DepositPaymentLookup, int64, int, error) {
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}

	activeStatuses := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}

	selectField := `r.return_no AS doc_no, r.total AS amount, (r.total - COALESCE(r_used.used, 0)) AS balance`
	qFrom := `
FROM sls."return" r
LEFT JOIN (
	SELECT dp.document_no, COALESCE(SUM(dp.payment_amount), 0) AS used
	FROM acf.deposit_payment dp
	INNER JOIN acf.deposit d ON d.deposit_no = dp.deposit_no AND d.cust_id = dp.cust_id
	WHERE dp.pay_type = 4
		AND dp.cust_id = ?
		AND d.deposit_status IN ?
	GROUP BY dp.document_no
) r_used ON r_used.document_no = r.return_no`

	qWhere := `WHERE r.cust_id = ? AND r.salesman_id = ? AND (r.total - COALESCE(r_used.used, 0)) > 0`
	args := []interface{}{dataFilter.CustId, activeStatuses, dataFilter.CustId, dataFilter.SalesmanId}

	if dataFilter.Query != "" {
		qWhere += ` AND r.return_no ILIKE ?`
		args = append(args, "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.OutletId) > 0 {
		parts := make([]string, len(dataFilter.OutletId))
		for i, id := range dataFilter.OutletId {
			parts[i] = strconv.Itoa(id)
		}
		qWhere += ` AND r.outlet_id IN (` + strings.Join(parts, ",") + `)`
	}

	queryCount := `SELECT COUNT(*) AS total ` + qFrom + ` ` + qWhere
	var total int64
	if err := repository.Raw(queryCount, args...).Scan(&total).Error; err != nil {
		log.Error("FindReturnBalance, count error:", err)
		return nil, 0, 0, err
	}

	sortBy := `r.return_no DESC`
	if dataFilter.Sort != "" {
		sortBy = dataFilter.Sort
	}
	offset := (dataFilter.Page - 1) * dataFilter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere + ` ORDER BY ` + sortBy +
		fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	var results []model.DepositPaymentLookup
	if err := repository.Raw(querySelect, args...).Scan(&results).Error; err != nil {
		log.Error("FindReturnBalance, scan error:", err)
		return results, total, 0, err
	}

	return results, total, lastPage, nil
}
