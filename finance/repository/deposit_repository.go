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
	RepositoryDepositImpl struct {
		*gorm.DB
	}
)

type DepositRepository interface {
	Store(c context.Context, data *model.Deposit) error
	FindByNo(depositNo string, custId string) (whAdj model.DepositList, err error)
	FindAllByCustId(dataFilter entity.DepositQueryFilter) ([]model.DepositList, int64, int, error)
	FindDepositNumberListByCustId(dataFilter entity.DepositNumberListQueryFilter) ([]model.DepositNumberList, int64, int, error)
	Update(c context.Context, depositNo string, custId string, data model.Deposit) error
	Delete(c context.Context, custId string, depositNo string, deletedBy int64) error
	CountAllByCustId(custId string, depositDate string) (int, error)
	StoreDetail(c context.Context, data *model.DepositDetail) (int, error)
	StorePayment(c context.Context, data *model.DepositPayment) (int, error)
	FindDetailByNo(depositNo string, custId string) (whAdj []model.DepositDetailList, err error)
	FindDetailPaymentByNo(depositNo string, invoiceNo string, custId string) (whAdj []model.DepositPayment, err error)
	CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error)
	DeleteAllDetailByDeposit(c context.Context, depositNo string) error
	DeleteAllDetailPaymentByDeposit(c context.Context, depositNo string) error
	DeleteAllExpenseByDeposit(c context.Context, depositNo string, custId string) error
	// RestoreExpensesByDeposit adds back deposit_expense.payment_amount to acf.expense.balance for given deposit
	RestoreExpensesByDeposit(c context.Context, depositNo string, custID string) error
	FindDetailPaymentInvoiceByNo(payType int, depositNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error)
	StoreExpense(c context.Context, data *model.DepositExpense) (int, error)
	DeductExpense(c context.Context, data *model.DepositExpense) error
	FindExpenseByDepositNo(depositNo string, custId string) ([]model.DepositExpense, error)
	CalcCollectionPaidByInvoice(c context.Context, data *model.DepositDetail) error
	UpdateAmountProgressionCheque(c context.Context, data *model.DepositPayment, auditor int64) error
	UpdateAmountProgressionTransfer(c context.Context, data *model.DepositPayment, auditor int64) error
	UpdateAmountProgressionReturn(c context.Context, data *model.DepositPayment, auditor int64) error
	UpdateAmountProgressionCNDN(c context.Context, data *model.DepositPayment, auditor int64) error
	FindProofOfPayment(depositNo string, q string, typ string, custId string) (items []map[string]interface{}, err error)
}

func NewDepositRepo(db *gorm.DB) *RepositoryDepositImpl {
	return &RepositoryDepositImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryDepositImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryDepositImpl) Store(c context.Context, data *model.Deposit) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryDepositImpl) StoreDetail(c context.Context, data *model.DepositDetail) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositDetailID, nil
}

func (repository *RepositoryDepositImpl) StorePayment(c context.Context, data *model.DepositPayment) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositPaymentID, nil
}

func (repository *RepositoryDepositImpl) FindByNo(depositNo string, custId string) (whAdj model.DepositList, err error) {
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

func (repository *RepositoryDepositImpl) FindDetailByNo(depositNo string, custId string) (whAdj []model.DepositDetailList, err error) {
	err = repository.Select(`
			acf.deposit_detail.*,
			ro.invoice_date,ro.due_date,ro.ro_no,ro.salesman_id,
			sls.sales_name as salesman_name,
			ot.outlet_id, ot.outlet_code, ot.outlet_name,
			(ro.total - coalesce(paid_invoices.paid_amount, 0)) as remaining_amount
		`).
		// Joins("left join sys.m_user us on us.user_id = acf.deposit_detail.updated_by").
		Joins("left join sls.order ro on ro.invoice_no = acf.deposit_detail.invoice_no AND ro.cust_id = ?", custId).
		Joins("left join mst.m_salesman sls on sls.emp_id = ro.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = ro.outlet_id AND ot.cust_id = ?", custId).
		Joins(`left join (
                        select acf.deposit_detail.invoice_no,
                        coalesce(sum(acf.deposit_detail.total_payment), 0) as paid_amount
                from acf.deposit_detail
                inner join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ? AND deposit.deposit_status IN (1, 2)
                where acf.deposit_detail.cust_id = ?
                group by acf.deposit_detail.invoice_no
        ) paid_invoices on paid_invoices.invoice_no = ro.invoice_no`, custId, custId).
		Where("acf.deposit_detail.deposit_no = ? AND acf.deposit_detail.cust_id=?", depositNo, custId).
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryDepositImpl) FindDetailPaymentByNo(depositNo string, invoiceNo string, custId string) (whAdj []model.DepositPayment, err error) {
	err = repository.Select(`
			acf.deposit_payment.*
		`).
		// Joins("left join sys.m_user us on us.user_id = acf.deposit_payment.updated_by").
		Where("acf.deposit_payment.deposit_no = ? AND acf.deposit_payment.invoice_no = ? AND acf.deposit_payment.cust_id=?", depositNo, invoiceNo, custId).
		// Where("acf.deposit_payment.is_del=false").
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryDepositImpl) FindDetailPaymentInvoiceByNo(payType int, depositNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error) {
	err = repository.Select(`
		DISTINCT ON (acf.deposit_payment.deposit_no, acf.deposit_payment.invoice_no) 
		acf.deposit_payment.*,
		od.invoice_date,
		emp.emp_id AS salesman_id,
		ot.outlet_id AS outlet_id,
		emp.emp_code AS salesman_code,
		emp.emp_name AS salesman_name,
		ot.outlet_code,
		ot.outlet_name,
		dd.notes
	`).
		Joins("LEFT JOIN acf.deposit_detail dd ON dd.deposit_no = acf.deposit_payment.deposit_no AND dd.invoice_no = acf.deposit_payment.invoice_no AND dd.cust_id = ?", custId).
		Joins("LEFT JOIN sls.order od ON od.invoice_no = dd.invoice_no AND od.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = od.salesman_id AND emp.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = od.outlet_id AND ot.cust_id = ?", custId).
		Where("acf.deposit_payment.deposit_no = ? AND acf.deposit_payment.cust_id = ?", depositNo, custId).
		Where("acf.deposit_payment.pay_type = ? AND acf.deposit_payment.cust_id = ?", payType, custId).
		Find(&whAdj).Error

	return whAdj, err
}

func (repository *RepositoryDepositImpl) FindProofOfPayment(depositNo string, q string, typ string, custId string) (items []map[string]interface{}, err error) {
	type fileRow struct {
		DocNo         string      `gorm:"column:doc_no"`
		FileID        int         `gorm:"column:file_id"`
		DetID         interface{} `gorm:"column:det_id"`
		FileName      string      `gorm:"column:file_name"`
		FileURL       string      `gorm:"column:file_url"`
		FileKey       string      `gorm:"column:file_key"`
		MediaCategory string      `gorm:"column:media_category"`
		FileSize      int64       `gorm:"column:file_size"`
	}

	var rows []fileRow

	like := "%" + q + "%"

	// choose query based on type
	if strings.ToUpper(typ) == "EXPENSE" {
		// join deposit_expense to filter by deposit_no
		sql := `SELECT e.doc_no AS doc_no, ef.expense_file_id AS file_id, ef.expense_id AS det_id, ef.file_name, ef.file_url, ef.file_key, ef.media_category, ef.file_size
			FROM acf.deposit_expense de
			JOIN acf.expense e ON e.expense_id = de.expense_id
			JOIN acf.expense_file ef ON e.expense_id = ef.expense_id
			JOIN acf.deposit d ON d.deposit_no = de.deposit_no AND d.cust_id = ?
			WHERE de.deposit_no = ?`
		// optionally filter by doc_no if q provided
		if q != "" {
			sql += " AND e.doc_no ILIKE ?"
			if err := repository.Raw(sql, custId, depositNo, like).Scan(&rows).Error; err != nil {
				return nil, err
			}
		} else {
			if err := repository.Raw(sql, custId, depositNo).Scan(&rows).Error; err != nil {
				return nil, err
			}
		}
	} else if strings.ToUpper(typ) == "TRANSFER" {
		// join deposit_payment to filter by deposit_no (deposit_payment.document_no -> bank_transfer.doc_no_bank)
		sql := `SELECT bt.doc_no_bank AS doc_no, btf.bank_transfer_file_id AS file_id, bt.bank_transfer_no AS det_id, btf.file_name, btf.file_url, btf.file_key, btf.media_category, btf.file_size
			FROM acf.bank_transfer bt
			JOIN acf.bank_transfer_files btf ON bt.doc_no_bank = btf.bank_transfer_no AND bt.cust_id = btf.cust_id
			JOIN acf.deposit_payment dp ON dp.document_no = bt.doc_no_bank AND dp.deposit_no = ? AND dp.cust_id = ?
			WHERE dp.deposit_no = ?`
		if q != "" {
			sql += " AND bt.doc_no_bank ILIKE ?"
			if err := repository.Raw(sql, depositNo, custId, depositNo, like).Scan(&rows).Error; err != nil {
				return nil, err
			}
		} else {
			if err := repository.Raw(sql, depositNo, custId, depositNo).Scan(&rows).Error; err != nil {
				return nil, err
			}
		}
	} else {
		// both: union expense and transfer
		// expense part
		expenseSQL := `SELECT e.doc_no AS doc_no, ef.expense_file_id AS file_id, ef.expense_id AS det_id, ef.file_name, ef.file_url, ef.file_key::text AS file_key, ef.media_category, ef.file_size
			FROM acf.deposit_expense de
			JOIN acf.expense e ON e.expense_id = de.expense_id
			JOIN acf.expense_file ef ON e.expense_id = ef.expense_id
			JOIN acf.deposit d ON d.deposit_no = de.deposit_no AND d.cust_id = ?
			WHERE de.deposit_no = ?`
		// transfer part
		transferSQL := `SELECT bt.doc_no_bank AS doc_no, btf.bank_transfer_file_id AS file_id, bt.bank_transfer_no AS det_id, btf.file_name, btf.file_url, btf.file_key, btf.media_category, btf.file_size
			FROM acf.bank_transfer bt
			JOIN acf.bank_transfer_files btf ON bt.doc_no_bank = btf.bank_transfer_no AND bt.cust_id = btf.cust_id
			JOIN acf.deposit_payment dp ON dp.document_no = bt.doc_no_bank AND dp.deposit_no = ? AND dp.cust_id = ?
			WHERE dp.deposit_no = ?`

		if q != "" {
			expenseSQL += " AND e.doc_no ILIKE ?"
			transferSQL += " AND bt.doc_no_bank ILIKE ?"
			unionSQL := expenseSQL + " UNION ALL " + transferSQL
			if err := repository.Raw(unionSQL, custId, depositNo, depositNo, custId, depositNo, like, like).Scan(&rows).Error; err != nil {
				return nil, err
			}
		} else {
			unionSQL := expenseSQL + " UNION ALL " + transferSQL
			if err := repository.Raw(unionSQL, custId, depositNo, depositNo, custId, depositNo).Scan(&rows).Error; err != nil {
				return nil, err
			}
		}
	}

	// group by doc_no
	grouped := make(map[string][]map[string]interface{})
	for _, r := range rows {
		file := map[string]interface{}{
			"file_id":        r.FileID,
			"det_id":         r.DetID,
			"file_name":      r.FileName,
			"file_url":       r.FileURL,
			"file_key":       r.FileKey,
			"media_category": r.MediaCategory,
			"file_size":      r.FileSize,
		}
		grouped[r.DocNo] = append(grouped[r.DocNo], file)
	}

	items = make([]map[string]interface{}, 0, len(grouped))
	for doc, files := range grouped {
		items = append(items, map[string]interface{}{"doc_no": doc, "files": files})
	}

	return items, nil
}

func (repository *RepositoryDepositImpl) FindAllByCustId(dataFilter entity.DepositQueryFilter) ([]model.DepositList, int64, int, error) {
	var Deposit []model.DepositList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("deposit_no")
	query := repository.Select(`
			acf.deposit.*,
			us.user_fullname AS updated_by_name,
			sls.sales_name as salesman_name,
			mpg.emp_grp_name,
			emp.emp_name,
			clc.collection_date
		`).
		Joins("left join sys.m_user us on us.user_id = acf.deposit.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.deposit.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_emp_group mpg on mpg.emp_grp_id = acf.deposit.emp_grp_id AND mpg.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_employee emp on emp.emp_id = acf.deposit.emp_id AND emp.cust_id = ?", dataFilter.CustId).
		Joins("left join acf.collection clc on clc.collection_no = acf.deposit.collection_no AND clc.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)

	// Filter where 'deleted_at' is NULL or empty
	queryCount.Where("acf.deposit.deleted_at IS NULL")
	query.Where("acf.deposit.deleted_at IS NULL")

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.deposit.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.deposit.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("acf.deposit.deposit_no=?", dataFilter.Query)
		queryCount.Where("acf.deposit.deposit_no=?", dataFilter.Query)
	}

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.deposit.deposit_no ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.deposit.deposit_no ILIKE ? )", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.Status) > 0 {
		query.Where("acf.deposit.deposit_status in ?", dataFilter.Status)
		queryCount.Where("acf.deposit.deposit_status in ?", dataFilter.Status)
	}

	if len(dataFilter.CollectionNo) > 0 {
		query.Where("acf.deposit.collection_no in ?", dataFilter.CollectionNo)
		queryCount.Where("acf.deposit.collection_no in ?", dataFilter.CollectionNo)
	}

	cond := repository.DB.Session(&gorm.Session{})

	if len(dataFilter.DepositNo) > 0 {
		cond = cond.Where("acf.deposit.deposit_no IN ?", dataFilter.DepositNo)
	}

	if len(dataFilter.DocumentNo) > 0 {
		if len(dataFilter.DepositNo) > 0 {
			cond = cond.Or("acf.deposit.deposit_no IN ?", dataFilter.DocumentNo)
		} else {
			cond = cond.Where("acf.deposit.deposit_no IN ?", dataFilter.DocumentNo)
		}
	}

	query = query.Where(cond)
	queryCount = queryCount.Where(cond)

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
		query.Order("deposit_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&Deposit).Error
	if err != nil {
		return Deposit, total, 0, err
	}
	err = queryCount.Model(&Deposit).Count(&total).Error
	if err != nil {
		return Deposit, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return Deposit, total, lastPage, nil
}

func (repository *RepositoryDepositImpl) FindDepositNumberListByCustId(dataFilter entity.DepositNumberListQueryFilter) ([]model.DepositNumberList, int64, int, error) {
	var deposits []model.DepositNumberList
	var total int64

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

	query := repository.Table("acf.deposit d").
		Select("d.deposit_no, d.emp_id as collector_id, d.deposit_date").
		Where("d.cust_id = ?", dataFilter.CustId).
		Where("d.deleted_at IS NULL")

	queryCount := repository.Table("acf.deposit d").
		Where("d.cust_id = ?", dataFilter.CustId).
		Where("d.deleted_at IS NULL")

	if dataFilter.Query != "" {
		query = query.Where("d.deposit_no ILIKE ?", "%"+dataFilter.Query+"%")
		queryCount = queryCount.Where("d.deposit_no ILIKE ?", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.CollectorIDs) > 0 {
		query = query.Where("d.emp_id IN ?", dataFilter.CollectorIDs)
		queryCount = queryCount.Where("d.emp_id IN ?", dataFilter.CollectorIDs)
	}

	sortBy := repository.depositNumberListSort(dataFilter.Sort)
	offset := (page - 1) * limit

	err := query.Order(sortBy).Limit(limit).Offset(offset).Find(&deposits).Error
	if err != nil {
		return deposits, total, 0, err
	}

	err = queryCount.Count(&total).Error
	if err != nil {
		return deposits, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return deposits, total, lastPage, nil
}

func (repository *RepositoryDepositImpl) depositNumberListSort(sort string) string {
	allowedColumns := map[string]string{
		"created_date": "d.created_at",
		"deposit_date": "d.deposit_date",
		"deposit_no":   "d.deposit_no",
	}

	defaultSort := "d.created_at DESC"
	if sort == "" {
		return defaultSort
	}

	sortToken := strings.TrimSpace(strings.Split(sort, ",")[0])
	sortParts := strings.Split(sortToken, ":")
	if len(sortParts) != 2 {
		return defaultSort
	}

	column, ok := allowedColumns[strings.TrimSpace(sortParts[0])]
	if !ok {
		return defaultSort
	}

	direction := strings.ToUpper(strings.TrimSpace(sortParts[1]))
	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func (repository *RepositoryDepositImpl) Delete(c context.Context, custId string, depositNo string, deletedBy int64) error {
	var data model.Deposit
	result := repository.model(c).Model(&data).Where("deposit_no=? AND cust_id = ?", depositNo, custId).
		Updates(map[string]interface{}{"deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryDepositImpl) Update(c context.Context, depositNo string, custId string, data model.Deposit) error {
	result := repository.model(c).Model(&data).Where("deposit_no=? AND cust_id = ?", depositNo, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	if data.RemainingAmount < 1 {
		result := repository.model(c).Model(&data).Where("deposit_no = ? AND cust_id = ?", depositNo, custId).Updates(map[string]interface{}{
			"remaining_amount": 0,
		})
		if result.Error != nil {
			return result.Error
		}
	}

	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}

func (repository *RepositoryDepositImpl) CountAllByCustId(custId string, depositDate string) (int, error) {
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

func (repository *RepositoryDepositImpl) CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error) {
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

func (repository *RepositoryDepositImpl) DeleteAllDetailByDeposit(c context.Context, depositNo string) error {
	var Details model.DepositDetail
	err := repository.model(c).Where("deposit_no=?", depositNo).Delete(&Details).Error

	return err
}

func (repository *RepositoryDepositImpl) DeleteAllDetailPaymentByDeposit(c context.Context, depositNo string) error {
	var Details model.DepositPayment
	err := repository.model(c).Where("deposit_no=?", depositNo).Delete(&Details).Error

	return err
}

func (repository *RepositoryDepositImpl) DeleteAllExpenseByDeposit(c context.Context, depositNo string, custId string) error {
	var Expenses model.DepositExpense
	err := repository.model(c).Where("deposit_no=? AND cust_id=?", depositNo, custId).Delete(&Expenses).Error

	return err
}

func (repository *RepositoryDepositImpl) RestoreExpensesByDeposit(c context.Context, depositNo string, custID string) error {
	// Add back payment_amount from deposit_expense to corresponding expense.balance
	// Use a single SQL update joining the subquery of deposit_expense
	return repository.model(c).Exec(`
		UPDATE acf.expense e
		SET balance = COALESCE(e.balance, 0) + sub.payment_amount,
			updated_at = NOW()
		FROM (
			SELECT expense_id, COALESCE(payment_amount, 0) AS payment_amount
			FROM acf.deposit_expense
			WHERE deposit_no = ?
			AND cust_id = ?
		) sub
		WHERE e.expense_id = sub.expense_id
	`, depositNo, custID).Error
}

func (repository *RepositoryDepositImpl) StoreExpense(c context.Context, data *model.DepositExpense) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return int(data.DepositExpenseID), nil
}

func (repository *RepositoryDepositImpl) DeductExpense(c context.Context, data *model.DepositExpense) error {
	return repository.model(c).Exec("UPDATE acf.expense SET balance = balance - ?, updated_at = NOW(), updated_by = ? WHERE expense_id = ?", data.PaymentAmount, data.CreatedBy, data.ExpenseID).Error
}

func (repository *RepositoryDepositImpl) FindExpenseByDepositNo(depositNo string, custId string) ([]model.DepositExpense, error) {
	var expenses []model.DepositExpense
	err := repository.Select(`
		de.deposit_expense_id,
		de.cust_id,
		e.expense_id,
		e.doc_no,
		COALESCE(e.balance, 0) as balance,
		de.payment_amount
	`).
		Table("acf.deposit_expense de").
		Joins("JOIN acf.expense e ON e.expense_id = de.expense_id").
		Joins("JOIN acf.deposit d ON d.deposit_no = de.deposit_no AND d.cust_id = ?", custId).
		Where("de.deposit_no = ?", depositNo).
		Where("de.cust_id = ?", custId).
		Where("(de.deleted_at IS NULL AND (de.is_del = false OR de.is_del IS NULL))").
		Scan(&expenses).Error
	return expenses, err
}

func (repository *RepositoryDepositImpl) CalcCollectionPaidByInvoice(c context.Context, data *model.DepositDetail) error {
	// recalculate invoice and collection payment snapshots from deposit detail source of truth
	err := repository.model(c).Exec(`
		UPDATE acf.collection_det cd
		SET paid_by_invoice = COALESCE((
			SELECT SUM(dd.total_payment)
			FROM acf.deposit_detail dd
			INNER JOIN acf.deposit d ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id
			WHERE dd.invoice_no = cd.invoice_no
			AND dd.cust_id = cd.cust_id
			AND d.deposit_status IN (1, 2)
		), 0),
		paid_amount = COALESCE((
			SELECT SUM(dd.total_payment)
			FROM acf.deposit_detail dd
			INNER JOIN acf.deposit d ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id
			WHERE dd.invoice_no = cd.invoice_no
			AND dd.cust_id = cd.cust_id
			AND d.collection_no = cd.collection_no
			AND d.deposit_status IN (1, 2)
		), 0)
		,
		remaining_amount = GREATEST(
			COALESCE(cd.invoice_amount, 0) - COALESCE((
				SELECT SUM(dd.total_payment)
				FROM acf.deposit_detail dd
				INNER JOIN acf.deposit d ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id
				WHERE dd.invoice_no = cd.invoice_no
				AND dd.cust_id = cd.cust_id
				AND d.deposit_status IN (1, 2)
			), 0),
			0
		)
		WHERE cd.invoice_no = ?
		AND cd.cust_id = ?
	`, data.InvoiceNo, data.CustID).Error
	if err != nil {
		return err
	}

	return repository.model(c).Exec(`
		WITH affected_collections AS (
			SELECT DISTINCT cd.collection_no, cd.cust_id
			FROM acf.collection_det cd
			WHERE cd.invoice_no = ?
			AND cd.cust_id = ?
		)
		UPDATE acf.collection c
		SET remaining_amount = sub.total_remaining_amount
		FROM (
			SELECT cd.collection_no, cd.cust_id, COALESCE(SUM(cd.remaining_amount), 0) AS total_remaining_amount
			FROM acf.collection_det cd
			INNER JOIN affected_collections ac ON ac.collection_no = cd.collection_no AND ac.cust_id = cd.cust_id
			GROUP BY cd.collection_no, cd.cust_id
		) sub
		WHERE c.collection_no = sub.collection_no
		AND c.cust_id = sub.cust_id
	`, data.InvoiceNo, data.CustID).Error
}

func (repository *RepositoryDepositImpl) UpdateAmountProgressionCheque(c context.Context, data *model.DepositPayment, auditor int64) error {
	return repository.model(c).Exec("UPDATE acf.cheque_giro SET paid_amount = ?, remaining_amount = amount - paid_amount, updated_at = NOW(), updated_by = ? WHERE cust_id = ? AND doc_no_cheque = ?", data.PaymentAmount, auditor, data.CustID, data.DocumentNo).Error
}

func (repository *RepositoryDepositImpl) UpdateAmountProgressionTransfer(c context.Context, data *model.DepositPayment, auditor int64) error {
	return repository.model(c).Exec("UPDATE acf.bank_transfer SET paid_amount = ?, remaining_amount = amount - paid_amount, updated_at = NOW(), updated_by = ? WHERE cust_id = ? AND doc_no_bank = ?", data.PaymentAmount, auditor, data.CustID, data.DocumentNo).Error
}

func (repository *RepositoryDepositImpl) UpdateAmountProgressionReturn(c context.Context, data *model.DepositPayment, auditor int64) error {
	return repository.model(c).Exec("UPDATE sls.return SET paid_amount = ?, remaining_amount = total - paid_amount, updated_at = NOW(), updated_by = ? WHERE cust_id = ? AND return_no = ?", data.PaymentAmount, auditor, data.CustID, data.DocumentNo).Error
}

func (repository *RepositoryDepositImpl) UpdateAmountProgressionCNDN(c context.Context, data *model.DepositPayment, auditor int64) error {
	return repository.model(c).Exec("UPDATE acf.cndn SET used_amount = ?, remaning_amount = amount - used_amount, updated_at = NOW(), updated_by = ? WHERE cust_id = ? AND cndn_no = ?", data.PaymentAmount, auditor, data.CustID, data.DocumentNo).Error
}
