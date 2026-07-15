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

	"github.com/gofiber/fiber/v2/log"

	"gorm.io/gorm"
)

type (
	RepositoryApPaymentImpl struct {
		*gorm.DB
	}
)

type ApPaymentRepository interface {
	Store(c context.Context, data *model.AccountPayablePayment) error
	FindByNo(AccountPayablePaymentNo string, custId string, ParentCustId string) (whAdj model.AccountPayablePaymentList, err error)
	FindAllByCustId(dataFilter entity.ApPaymentQueryFilter) ([]model.AccountPayablePaymentList, int64, int, error)
	Update(c context.Context, AccountPayablePaymentNo string, custId string, data model.AccountPayablePayment) error
	Delete(c context.Context, custId string, AccountPayablePaymentNo string, deletedBy int64) error
	StoreApPaymentDetail(c context.Context, data *model.AccountPayablePaymentDetail) (int, error)
	StoreApPaymentOptions(c context.Context, data *model.AccountPayablePaymentOptions) (int, error)
	FindDetailByNo(AccountPayablePaymentNo string, custId string) (whAdj []model.AccountPayablePaymentDetailList, err error)
	FindDetailPaymentByNo(AccountPayablePaymentNo string, invoiceNo string, custId string, totalPayment float64) (whAdj []model.AccountPayablePaymentOptionsList, err error)
	DeleteAllDetailByApPayment(c context.Context, AccountPayablePaymentNo string, custId string) error
	DeleteAllDetailPaymentByApPayment(c context.Context, AccountPayablePaymentNo string, custId string) error
	FindDetailApPaymentOptionsByNo(payType int, AccountPayablePaymentNo string, custId string) (whAdj []model.AccountPayablePaymentOptionsList, err error)

	FindAllBalancePaymentDepositByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositPaymentLookup, int64, int, error)
	FindAllInvoiceNo(dataFilter entity.ApLookupSupplierInoviceReturnQueryFilter) ([]model.ApLookupSuppilerInvoiceReturnList, int64, int, error)
}

func NewApPaymentRepo(db *gorm.DB) *RepositoryApPaymentImpl {
	return &RepositoryApPaymentImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApPaymentImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApPaymentImpl) Store(c context.Context, data *model.AccountPayablePayment) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApPaymentImpl) StoreApPaymentDetail(c context.Context, data *model.AccountPayablePaymentDetail) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.AccountPayablePaymentDetailId, nil
}

func (repository *RepositoryApPaymentImpl) StoreApPaymentOptions(c context.Context, data *model.AccountPayablePaymentOptions) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.AccountPayablePaymentOptionsId, nil
}

func (repository *RepositoryApPaymentImpl) FindByNo(AccountPayablePaymentNo string, custId string, ParentCustId string) (whAdj model.AccountPayablePaymentList, err error) {
	err = repository.Select(`
			acf.account_payable_payment.*, acf.account_payable_payment.account_payable_payment_no as document_no,
			us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name,
			ms.sup_name,ms.sup_code,
			sc.distributor_id as distributor_id,md.distributor_code as distributor_code, md.distributor_name as distributor
		`).
		Joins("left join sys.m_user us on us.user_id = acf.account_payable_payment.updated_by").
		Joins("left join smc.m_customer sc ON sc.cust_id = acf.account_payable_payment.cust_id").
		Joins("left join mst.m_distributor md ON md.distributor_id = sc.distributor_id").
		Joins("left join mst.m_supplier ms on ms.sup_id = acf.account_payable_payment.sup_id AND ms.cust_id = ?", ParentCustId).
		Where("acf.account_payable_payment.account_payable_payment_no = ? AND acf.account_payable_payment.cust_id=?", AccountPayablePaymentNo, custId).
		Where("acf.account_payable_payment.deleted_at IS NULL").
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryApPaymentImpl) FindDetailByNo(AccountPayablePaymentNo string, custId string) (whAdj []model.AccountPayablePaymentDetailList, err error) {
	err = repository.Select(`
			acf.account_payable_payment_detail.*
		`).
		Where("acf.account_payable_payment_detail.account_payable_payment_no = ? AND acf.account_payable_payment_detail.cust_id=?", AccountPayablePaymentNo, custId).
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryApPaymentImpl) FindDetailPaymentByNo(AccountPayablePaymentNo string, invoiceNo string, custId string, totalPayment float64) (whAdj []model.AccountPayablePaymentOptionsList, err error) {
	query := repository.Model(&model.AccountPayablePaymentOptions{}).
		Select(`
			appd.invoice_date,
			acf.account_payable_payment_options.account_payable_payment_no,
			acf.account_payable_payment_options.invoice_no,
			acf.account_payable_payment_options.pay_type,
			acf.account_payable_payment_options.document_no,
			acf.account_payable_payment_options.balance,
			acf.account_payable_payment_options.payment_amount,
			GREATEST(COALESCE(appd.invoice_amount, 0) - COALESCE(paid.paid_amount, 0), 0) as remaining_amount
		`).
		Joins("left join acf.account_payable_payment_detail appd on appd.account_payable_payment_no = acf.account_payable_payment_options.account_payable_payment_no AND appd.invoice_no = acf.account_payable_payment_options.invoice_no AND appd.cust_id = acf.account_payable_payment_options.cust_id AND appd.total_payment = acf.account_payable_payment_options.payment_amount").
		Joins(`left join (
			select d.invoice_no,
				d.cust_id,
				coalesce(sum(d.total_payment), 0) as paid_amount
			from acf.account_payable_payment_detail d
			inner join acf.account_payable_payment p
				on p.account_payable_payment_no = d.account_payable_payment_no
				and p.cust_id = d.cust_id
			where p.deleted_at is null
			group by d.invoice_no, d.cust_id
		) paid on paid.invoice_no = acf.account_payable_payment_options.invoice_no and paid.cust_id = acf.account_payable_payment_options.cust_id`).
		Where("acf.account_payable_payment_options.account_payable_payment_no = ? AND acf.account_payable_payment_options.invoice_no = ? AND acf.account_payable_payment_options.cust_id=?", AccountPayablePaymentNo, invoiceNo, custId)

	// Pair options to the detail line when invoice_no is shared across multiple details
	if totalPayment > 0 {
		query = query.Where("acf.account_payable_payment_options.payment_amount = ?", totalPayment)
	}

	err = query.Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryApPaymentImpl) FindAllByCustId(dataFilter entity.ApPaymentQueryFilter) ([]model.AccountPayablePaymentList, int64, int, error) {
	var ApPayment []model.AccountPayablePaymentList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("account_payable_payment_no")
	query := repository.Select(`
			acf.account_payable_payment.*, acf.account_payable_payment.account_payable_payment_no as document_no,
			us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name,
			ms.sup_name,ms.sup_code, ms.sup_id,
			sc.distributor_id as distributor_id,md.distributor_code as distributor_code, md.distributor_name as distributor
		`).
		Joins("left join sys.m_user us on us.user_id = acf.account_payable_payment.updated_by").
		Joins("left join smc.m_customer sc ON sc.cust_id = acf.account_payable_payment.cust_id").
		Joins("left join mst.m_distributor md ON md.distributor_id = sc.distributor_id").
		Joins("left join mst.m_supplier ms on ms.sup_id = acf.account_payable_payment.sup_id AND ms.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.account_payable_payment.cust_id=?", dataFilter.CustId)
	query.Where("acf.account_payable_payment.cust_id=?", dataFilter.CustId)

	// Filter where 'deleted_at' is NULL or empty
	queryCount.Where("acf.account_payable_payment.deleted_at IS NULL")
	query.Where("acf.account_payable_payment.deleted_at IS NULL")

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.account_payable_payment.account_payable_payment_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.account_payable_payment.account_payable_payment_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("acf.account_payable_payment.account_payable_payment_no=?", dataFilter.Query)
		queryCount.Where("acf.account_payable_payment.account_payable_payment_no=?", dataFilter.Query)
	}

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.account_payable_payment.account_payable_payment_no ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.account_payable_payment.account_payable_payment_no ILIKE ? )", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.DocumentNo) > 0 {
		query.Where("acf.account_payable_payment.account_payable_payment_no in ?", dataFilter.DocumentNo)
		queryCount.Where("acf.account_payable_payment.account_payable_payment_no in ?", dataFilter.DocumentNo)
	}

	if dataFilter.SuppId != 0 {
		query.Where("acf.account_payable_payment.sup_id=?", dataFilter.SuppId)
		queryCount.Where("acf.account_payable_payment.sup_id=?", dataFilter.SuppId)
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
		query.Order("account_payable_payment_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ApPayment).Error
	if err != nil {
		return ApPayment, total, 0, err
	}
	err = queryCount.Model(&ApPayment).Count(&total).Error
	if err != nil {
		return ApPayment, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ApPayment, total, lastPage, nil
}

func (repository *RepositoryApPaymentImpl) Delete(c context.Context, custId string, AccountPayablePaymentNo string, deletedBy int64) error {
	var data model.AccountPayablePayment
	result := repository.model(c).Model(&data).Where("account_payable_payment_no=? AND cust_id = ?", AccountPayablePaymentNo, custId).
		Updates(map[string]interface{}{"deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApPaymentImpl) Update(c context.Context, AccountPayablePaymentNo string, custId string, data model.AccountPayablePayment) error {
	fmt.Println("mausuk reposi")
	result := repository.model(c).Model(&data).Where("account_payable_payment_no=? AND cust_id=?", AccountPayablePaymentNo, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository *RepositoryApPaymentImpl) DeleteAllDetailByApPayment(c context.Context, AccountPayablePaymentNo string, custId string) error {
	var Details model.AccountPayablePaymentDetail
	err := repository.model(c).Where("account_payable_payment_no=? AND cust_id=?", AccountPayablePaymentNo, custId).Delete(&Details).Error

	return err
}

func (repository *RepositoryApPaymentImpl) DeleteAllDetailPaymentByApPayment(c context.Context, AccountPayablePaymentNo string, custId string) error {
	var Details model.AccountPayablePaymentOptions
	err := repository.model(c).Where("account_payable_payment_no=? AND cust_id=?", AccountPayablePaymentNo, custId).Delete(&Details).Error

	return err
}

func (repository *RepositoryApPaymentImpl) FindDetailApPaymentOptionsByNo(payType int, AccountPayablePaymentNo string, custId string) (whAdj []model.AccountPayablePaymentOptionsList, err error) {
	subQuery := repository.
		Table("acf.account_payable_payment_options").
		Select(`appd.invoice_date,
            acf.account_payable_payment_options.account_payable_payment_no,
            acf.account_payable_payment_options.invoice_no,
            acf.account_payable_payment_options.pay_type,
            acf.account_payable_payment_options.document_no,
            acf.account_payable_payment_options.balance,
            acf.account_payable_payment_options.payment_amount,
            GREATEST(COALESCE(appd.invoice_amount, 0) - COALESCE(paid.paid_amount, 0), 0) AS remaining_amount`).
		Joins(`LEFT JOIN acf.account_payable_payment_detail appd ON appd.account_payable_payment_no = acf.account_payable_payment_options.account_payable_payment_no AND appd.invoice_no = acf.account_payable_payment_options.invoice_no AND appd.cust_id = acf.account_payable_payment_options.cust_id AND appd.total_payment = acf.account_payable_payment_options.payment_amount`).
		Joins(`LEFT JOIN (
			SELECT d.invoice_no,
				d.cust_id,
				COALESCE(SUM(d.total_payment), 0) AS paid_amount
			FROM acf.account_payable_payment_detail d
			INNER JOIN acf.account_payable_payment p
				ON p.account_payable_payment_no = d.account_payable_payment_no
				AND p.cust_id = d.cust_id
			WHERE p.deleted_at IS NULL
			GROUP BY d.invoice_no, d.cust_id
		) paid ON paid.invoice_no = acf.account_payable_payment_options.invoice_no AND paid.cust_id = acf.account_payable_payment_options.cust_id`).
		Where(`acf.account_payable_payment_options.account_payable_payment_no = ? AND acf.account_payable_payment_options.cust_id = ?`, AccountPayablePaymentNo, custId).
		Where(`acf.account_payable_payment_options.pay_type = ?`, payType)

	err = repository.
		Table("(?) as t", subQuery).
		Select("t.*, appd2.payment_balance").
		Joins(`
				LEFT JOIN LATERAL (
					SELECT payment_balance
					FROM acf.account_payable_payment_detail d
					WHERE d.account_payable_payment_no = t.account_payable_payment_no
					  AND d.invoice_no = t.invoice_no
					  AND d.cust_id = ?
					  AND d.total_payment = t.payment_amount
					LIMIT 1
				) appd2 ON true
			`, custId).
		Find(&whAdj).Error

	return whAdj, err
}

func (repository *RepositoryApPaymentImpl) FindAllBalancePaymentDepositByCustId(dataFilter entity.GeneralQueryFilter) ([]model.DepositPaymentLookup, int64, int, error) {
	var ChequeGiro []model.DepositPaymentLookup
	var total int64

	var docNo = "doc_no_cheque"
	var tableCheck = "acf.cheque_giro"
	var amountAs = "amount"
	// var apType = 2

	if dataFilter.Mode == "check" {
		docNo = "doc_no_cheque"
		tableCheck = "acf.cheque_giro"
		amountAs = "amount"
		// apType = 2
	} else if dataFilter.Mode == "transfer" {
		docNo = "doc_no_bank"
		tableCheck = "acf.bank_transfer"
		amountAs = "amount"
		// apType = 3
	} else if dataFilter.Mode == "cndn" {
		docNo = "cndn_no"
		tableCheck = "acf.cndn"
		amountAs = "amount"
		// apType = 4
	} else if dataFilter.Mode == "return" {
		docNo = "document_no"
		tableCheck = "acf.account_payable"
		amountAs = "total"
		// apType = 5
	}

	// Building the SQL queries
	selectCount := `COUNT(*) AS total`
	// selectField := `cg.` + docNo + ` as doc_no, cg.` + amountAs + ` as amount, (cg.` + amountAs + ` - SUM(appo.payment_amount)) as balance`
	selectField := `cg.` + docNo + ` as doc_no, cg.` + amountAs + ` as amount, case when sum(appo.payment_amount) is null then sum(cg.` + amountAs + `) else (cg.` + amountAs + ` - sum(appo.payment_amount))  end as balance`
	qWhere := `WHERE cg.is_del = false AND cg.cust_id = '` + dataFilter.CustId + `' and app.deleted_at is null  `

	if dataFilter.Mode != "return" {
		qWhere += `AND cg.owner_id = 2`
	}

	if dataFilter.Mode == "return" {
		qWhere += `AND cg.ap_type = 'R'`
	}

	// Adding search query filter
	if dataFilter.Query != "" {
		qWhere += ` AND cg.` + docNo + ` ILIKE '%' || '` + dataFilter.Query + `' || '%'`
	}

	qFrom := `
		FROM ` + tableCheck + ` cg LEFT JOIN acf.account_payable_payment_options appo on appo.document_no = cg.` + docNo + `
		left join acf.account_payable_payment app on app.account_payable_payment_no = appo.account_payable_payment_no`

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
	querySelect += `group by cg.` + docNo + `, cg.amount, cg.` + amountAs + ` `
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

func (repository *RepositoryApPaymentImpl) FindAllInvoiceNo(dataFilter entity.ApLookupSupplierInoviceReturnQueryFilter) ([]model.ApLookupSuppilerInvoiceReturnList, int64, int, error) {
	var ap []model.ApLookupSuppilerInvoiceReturnList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("invoice_no")
	query := repository.Select("acf.account_payable.invoice_no, acf.account_payable.ap_type, acf.account_payable.document_no, ROUND(acf.account_payable.total) as amount, acf.account_payable.sub_total, us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name,sup.sup_id, sup.sup_code,sup.sup_name, ROUND(coalesce(sum(CASE WHEN app_hdr.account_payable_payment_no IS NOT NULL THEN app.total_payment ELSE 0 END), 0)) as paid_amount, ROUND(case when (acf.account_payable.amount - coalesce(sum(CASE WHEN app_hdr.account_payable_payment_no IS NOT NULL THEN app.total_payment ELSE 0 END), 0)) != (acf.account_payable.amount) then (acf.account_payable.total - coalesce(sum(CASE WHEN app_hdr.account_payable_payment_no IS NOT NULL THEN app.total_payment ELSE 0 END), 0)) else acf.account_payable.total end) as remaining_amount  ").
		Joins("left join sys.m_user us on us.user_id = acf.account_payable.updated_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId).
		Joins("left join acf.account_payable_payment_detail app on app.invoice_no = acf.account_payable.invoice_no AND app.cust_id = acf.account_payable.cust_id").
		Joins("left join acf.account_payable_payment app_hdr on app_hdr.account_payable_payment_no = app.account_payable_payment_no AND app_hdr.cust_id = acf.account_payable.cust_id AND app_hdr.deleted_at IS NULL")

	queryCount.Where("acf.account_payable.cust_id=?", dataFilter.CustId)
	query.Where("acf.account_payable.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.account_payable.account_payable_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.account_payable.account_payable_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.account_payable.invoice_no=?", dataFilter.Query)
		query.Where("acf.account_payable.invoice_no=?", dataFilter.Query)
	}

	if dataFilter.DocumentNo != "" {
		queryCount.Where("acf.account_payable.document_no=?", dataFilter.DocumentNo)
		query.Where("acf.account_payable.document_no=?", dataFilter.DocumentNo)
	}

	if dataFilter.Type != "" {
		queryCount.Where("acf.account_payable.ap_type=?", dataFilter.Type)
		query.Where("acf.account_payable.ap_type=?", dataFilter.Type)
	}

	if dataFilter.SuppId != 0 {
		queryCount.Where("acf.account_payable.sup_id=?", dataFilter.SuppId)
		query.Where("acf.account_payable.sup_id=?", dataFilter.SuppId)
	}

	if dataFilter.InvoiceNo != "" {
		queryCount.Where("acf.account_payable.invoice_no=?", dataFilter.InvoiceNo)
		query.Where("acf.account_payable.invoice_no=?", dataFilter.InvoiceNo)
	}

	if dataFilter.ExcludeEmptyInvoice {
		query.Where("acf.account_payable.invoice_no <> '' AND acf.account_payable.is_can_return = true")
		queryCount.Where("acf.account_payable.invoice_no <> '' AND acf.account_payable.is_can_return = true")
	}

	query.Group("acf.account_payable.invoice_no, acf.account_payable.amount, acf.account_payable.sub_total, us.user_fullname, us.user_fullname, sup.sup_id,sup.sup_code,sup.sup_name,acf.account_payable.ap_type, acf.account_payable.document_no,acf.account_payable.total")

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

	err := query.Limit(limit).Offset(offset).Find(&ap).Error
	if err != nil {
		return ap, total, 0, err
	}
	err = queryCount.Model(&ap).Count(&total).Error
	if err != nil {
		return ap, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ap, total, lastPage, nil
}
