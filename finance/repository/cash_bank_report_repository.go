package repository

import (
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryCashBankReportImpl struct {
		*gorm.DB
	}
)

type CashBankReportRepository interface {
	FindAllReportDeposit(dataFilter entity.DepositQueryFilter, custId string) ([]model.DepositReport, int64, int, error)

	FindAllDepositNoReportFilterByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositNoReportLookup, int64, int, error)
	FindAllDepositPaymentTypeReportFilterByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositPayTypeLookup, int64, int, error)
}

func NewCashBankReportRepo(db *gorm.DB) *RepositoryCashBankReportImpl {
	return &RepositoryCashBankReportImpl{db}
}

func (repository *RepositoryCashBankReportImpl) FindAllReportDeposit(dataFilter entity.DepositQueryFilter, custId string) ([]model.DepositReport, int64, int, error) {
	var Deposit []model.DepositReport
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("deposit_payment_id").
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT ON (deposit_no, invoice_no, cust_id)
					deposit_no, invoice_no, cust_id, notes
				FROM acf.deposit_detail
				WHERE cust_id = ?
				ORDER BY deposit_no, invoice_no, cust_id, notes
			) dd ON dd.deposit_no = acf.deposit_payment.deposit_no AND dd.invoice_no = acf.deposit_payment.invoice_no
		`, custId).
		Joins("LEFT JOIN acf.deposit dp ON dp.deposit_no = dd.deposit_no AND dp.cust_id = ?", custId)

	query := repository.Select(`
			acf.deposit_payment.*,
			od.invoice_date, od.total as invoice_amount, od.due_date, od.outlet_id,
			emps.emp_id as salesman_id,
			dp.deposit_date, dp.emp_grp_id, dp.emp_id, dp.salesman_id,
			ot.outlet_id as outlet_id,
			emps.emp_code as salesman_code, 
			emps.emp_name as salesman_name,
			empc.emp_name as emp_name,
			empc.emp_code as emp_code,
			mpg.emp_grp_name,
			ot.outlet_code, ot.outlet_name,
			CASE 
			WHEN acf.deposit_payment.pay_type = 1 THEN 'Cash'
			WHEN acf.deposit_payment.pay_type = 2 THEN 'Cheque/Bilyet Giro'
			WHEN acf.deposit_payment.pay_type = 3 THEN 'Transfer'
			WHEN acf.deposit_payment.pay_type = 4 THEN 'Return'
			WHEN acf.deposit_payment.pay_type = 5 THEN 'Credit Note'
			ELSE '-'
			END AS pay_type_name,
			CASE 
			WHEN acf.deposit_payment.pay_type = 2 THEN mbc.bank_name
			WHEN acf.deposit_payment.pay_type = 3 THEN mbt.bank_name
			ELSE '-'
			END AS bank_name,
			CASE 
			WHEN acf.deposit_payment.pay_type = 2 THEN cg.account_no
			WHEN acf.deposit_payment.pay_type = 3 THEN bt.account_no
			ELSE '-'
			END AS account_no,
			CASE 
			WHEN acf.deposit_payment.pay_type = 2 THEN cg.status_cheque
			ELSE NULL
			END AS clearing_status,
			CASE 
			WHEN acf.deposit_payment.pay_type = 2 THEN cg.clearing_date
			ELSE NULL
			END AS clearing_date,
			'Outlet' as owner_name,
			CASE 
			WHEN acf.deposit_payment.pay_type = 2 THEN cg.doc_date_cheque
			WHEN acf.deposit_payment.pay_type = 3 THEN bt.transfer_date
			WHEN acf.deposit_payment.pay_type = 4 THEN r.return_date
			WHEN acf.deposit_payment.pay_type = 5 THEN cn.cndn_date
			END AS document_date,
			dd.notes
		`).
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT ON (deposit_no, invoice_no, cust_id)
					deposit_no, invoice_no, cust_id, notes
				FROM acf.deposit_detail
				WHERE cust_id = ?
				ORDER BY deposit_no, invoice_no, cust_id, notes
			) dd ON dd.deposit_no = acf.deposit_payment.deposit_no AND dd.invoice_no = acf.deposit_payment.invoice_no
		`, custId).
		Joins("LEFT JOIN acf.deposit dp ON dp.deposit_no = dd.deposit_no AND dp.cust_id = ?", custId).
		Joins("LEFT JOIN sls.order od ON od.invoice_no = dd.invoice_no AND od.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_employee emps ON emps.emp_id = dp.salesman_id AND emps.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_employee empc ON empc.emp_id = dp.emp_id AND empc.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = od.outlet_id AND ot.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_emp_group mpg on mpg.emp_grp_id = dp.emp_grp_id AND mpg.cust_id = ?", dataFilter.ParentCustId).
		Joins("LEFT JOIN acf.cheque_giro cg on cg.doc_no_cheque = acf.deposit_payment.document_no AND cg.cust_id = ? AND cg.deleted_at ISNULL", custId).
		Joins("LEFT JOIN acf.bank_transfer bt on bt.doc_no_bank = acf.deposit_payment.document_no AND bt.cust_id = ? AND bt.deleted_at ISNULL", custId).
		Joins("LEFT JOIN acf.cndn cn on cn.cndn_no = acf.deposit_payment.document_no AND cn.cust_id = ? AND cn.deleted_at ISNULL", custId).
		Joins("LEFT JOIN sls.return r on r.return_no = acf.deposit_payment.document_no AND r.cust_id = ? AND r.deleted_at ISNULL", custId).
		Joins("LEFT JOIN mst.m_bank mbc on cg.bank_id = mbc.bank_id AND mbc.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_bank mbt on bt.bank_id = mbt.bank_id AND mbt.cust_id = ?", custId).
		Where("acf.deposit_payment.cust_id = ?", custId)
	// Where("acf.deposit_payment.is_del=false").

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("dp.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("dp.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.Type) > 0 {
		queryCount.Where("acf.deposit_payment.pay_type in ?", dataFilter.Type)
		query.Where("acf.deposit_payment.pay_type in ?", dataFilter.Type)
	}

	if len(dataFilter.DepositNo) > 0 {
		query.Where("acf.deposit_payment.deposit_no in ?", dataFilter.DepositNo)
		queryCount.Where("acf.deposit_payment.deposit_no in ?", dataFilter.DepositNo)
	}

	if len(dataFilter.DocumentNo) > 0 {
		query.Where("acf.deposit_payment.document_no in ?", dataFilter.DocumentNo)
		queryCount.Where("acf.deposit_payment.document_no in ?", dataFilter.DocumentNo)
	}

	if len(dataFilter.InvoiceNo) > 0 {
		query.Where("dd.invoice_no in ?", dataFilter.InvoiceNo)
		queryCount.Where("dd.invoice_no in ?", dataFilter.InvoiceNo)
	}

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.deposit_payment.invoice_no ILIKE ? OR acf.deposit_payment.document_no ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.deposit_payment.invoice_no ILIKE ? OR acf.deposit_payment.document_no ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
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
		query.Order("deposit_payment_id DESC")
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

func (repository *RepositoryCashBankReportImpl) FindAllDepositNoReportFilterByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositNoReportLookup, int64, int, error) {
	var DepositLookup []model.DepositNoReportLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.deposit_payment.deposit_no))")
	query := repository.Select(
		`distinct(acf.deposit_payment.deposit_no)`)
	queryCount.Where("acf.deposit_payment.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit_payment.cust_id=?", dataFilter.CustId)

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.deposit_payment.deposit_no ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.deposit_payment.deposit_no ILIKE ? )", "%"+dataFilter.Query+"%")
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

func (repository *RepositoryCashBankReportImpl) FindAllDepositPaymentTypeReportFilterByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositPayTypeLookup, int64, int, error) {
	var DepositLookup []model.DepositPayTypeLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.deposit_payment.pay_type))")
	query := repository.Select(
		`distinct(acf.deposit_payment.pay_type)`)
	queryCount.Where("acf.deposit_payment.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit_payment.cust_id=?", dataFilter.CustId)

	if dataFilter.Query != "" {
		query.Where("acf.deposit_payment.pay_type=?", dataFilter.Query)
		queryCount.Where("acf.deposit_payment.pay_type=?", dataFilter.Query)
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
		query.Order("pay_type ASC")
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
