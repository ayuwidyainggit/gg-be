package repository

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryApListImpl struct {
		*gorm.DB
	}
)
type ApListRepository interface {
	FindByNo(invNo string, ParentCustId string, custId string) (ap model.AccountPayableList, err error)
	FindDetail(invNo string, custId string) (Details []model.AccountPayableListDet, err error)
	FindAllByCustId(dataFilter entity.AccountPayableListQueryFilter) ([]model.AccountPayableList, int64, int, error)
}

func NewApListRepo(db *gorm.DB) *RepositoryApListImpl {
	return &RepositoryApListImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApListImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApListImpl) FindByNo(invNo string, ParentCustId string, custId string) (ap model.AccountPayableList, err error) {
	invNo, _ = url.QueryUnescape(invNo)

	err = repository.Select(`acf.account_payable.invoice_date as inv_date, acf.account_payable.due_date as inv_due_date, acf.account_payable.invoice_no as inv_no, acf.account_payable.total as inv_amount,
				case when sum(appd.total_payment) is null then 0 else sum(appd.total_payment) end as amount_paid, 
				case when appd.invoice_amount is null then acf.account_payable.total else (acf.account_payable.total - sum(appd.total_payment))  end as remaining_amount,
				ms.sup_id as supplier_id,ms.sup_code as supplier_code, ms.sup_name as supplier,
				CASE 
					WHEN LEFT(acf.account_payable.document_no, 3) = 'GRB' THEN gb.po_no
					ELSE gr.po_no
				END AS po_no_doc
				,
				sc.distributor_id as distributor_id,md.distributor_code as distributor_code, md.distributor_name as distributor,
				case when acf.account_payable.invoice_date < acf.account_payable.due_date then 'On Schedule' else 'Overdue' end as due_date_status, extract(day from now()) - extract(day from acf.account_payable.due_date ) as aging,
				us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name, acf.account_payable.created_at, acf.account_payable.updated_at`).
		Joins("left join acf.account_payable_payment_detail appd on appd.invoice_no = acf.account_payable.invoice_no AND appd.cust_id = acf.account_payable.cust_id").
		Joins("left join acf.account_payable_payment app on app.account_payable_payment_no = appd.account_payable_payment_no AND app.cust_id = acf.account_payable.cust_id").
		Joins("left join sys.m_user us on us.user_id = acf.account_payable.updated_by").
		Joins("left join smc.m_customer sc ON sc.cust_id = acf.account_payable.cust_id").
		Joins("left join inv.gr_branch gb ON gb.gr_branch_no = acf.account_payable.document_no").
		Joins("left join inv.gr gr ON gr.gr_no = acf.account_payable.document_no").
		Joins("left join mst.m_distributor md ON md.distributor_id = sc.distributor_id").
		Joins("left join mst.m_supplier ms on ms.sup_id = acf.account_payable.sup_id AND ms.cust_id = ?", ParentCustId).
		Where("acf.account_payable.invoice_no = ? AND acf.account_payable.cust_id=?", invNo, custId).
		Group(`acf.account_payable.invoice_date,
			app.account_payable_payment_date,
			acf.account_payable.due_date,
			acf.account_payable.invoice_no,
			appd.invoice_amount, ms.sup_id,
			po_no_doc,
			ms.sup_id,
			ms.sup_code,
			ms.sup_name,
			sc.distributor_id,
			md.distributor_code,
			md.distributor_name,
			us.user_fullname,
			acf.account_payable.created_at,
			acf.account_payable.updated_at,
			acf.account_payable.amount,
			acf.account_payable.total`).
		Take(&ap).Error
	return ap, err
}

func (repository *RepositoryApListImpl) FindDetail(invNo string, custId string) (Details []model.AccountPayableListDet, err error) {
	err = repository.Model(&model.AccountPayablePaymentOptions{}).
		Select(`acf.account_payable_payment_options.pay_type as payment_method, appd.invoice_date as payment_date, appd.payment_balance as payment_balance, acf.account_payable_payment_options.document_no, acf.account_payable_payment_options.payment_amount as amount, app.updated_by, us.user_fullname as updated_by_name, app.updated_at`).
		Joins("inner join acf.account_payable_payment app on app.account_payable_payment_no = acf.account_payable_payment_options.account_payable_payment_no AND app.cust_id = acf.account_payable_payment_options.cust_id").
		Joins("left join acf.account_payable_payment_detail appd on appd.account_payable_payment_no = acf.account_payable_payment_options.account_payable_payment_no AND appd.invoice_no = acf.account_payable_payment_options.invoice_no AND appd.cust_id = acf.account_payable_payment_options.cust_id AND appd.total_payment = acf.account_payable_payment_options.payment_amount").
		Joins("left join sys.m_user us on us.user_id = app.updated_by").
		Where("acf.account_payable_payment_options.invoice_no = ? AND acf.account_payable_payment_options.cust_id = ?", invNo, custId).
		Where("app.deleted_at IS NULL").
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApListImpl) FindAllByCustId(dataFilter entity.AccountPayableListQueryFilter) ([]model.AccountPayableList, int64, int, error) {
	var ap []model.AccountPayableList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("invoice_no")
	query := repository.Select(`acf.account_payable.invoice_date as inv_date, acf.account_payable.due_date as inv_due_date, acf.account_payable.invoice_no as inv_no, acf.account_payable.total as inv_amount,
			case when sum(appd.total_payment) is null then 0 else sum(appd.total_payment) end as amount_paid, case when appd.invoice_amount is null then acf.account_payable.total else (acf.account_payable.total - sum(appd.total_payment)) end as remaining_amount,
			ms.sup_id as supplier_id,ms.sup_code as supplier_code, ms.sup_name as supplier,
			sc.distributor_id as distributor_id,md.distributor_code as distributor_code, md.distributor_name as distributor,
			case when acf.account_payable.invoice_date < acf.account_payable.due_date then 'On Schedule' else 'Overdue' end as due_date_status, extract(day from now()) - extract(day from acf.account_payable.due_date ) as aging,
			us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name, acf.account_payable.created_at, acf.account_payable.updated_at`).
		Joins("left join acf.account_payable_payment_detail appd on appd.invoice_no = acf.account_payable.invoice_no AND appd.cust_id = acf.account_payable.cust_id").
		Joins("left join acf.account_payable_payment app on app.account_payable_payment_no = appd.account_payable_payment_no AND app.cust_id = acf.account_payable.cust_id").
		Joins("left join sys.m_user us ON us.user_id = acf.account_payable.updated_by").
		Joins("left join smc.m_customer sc ON sc.cust_id = acf.account_payable.cust_id").
		Joins("left join mst.m_distributor md ON md.distributor_id = sc.distributor_id").
		Joins("left join mst.m_supplier ms ON ms.sup_id = acf.account_payable.sup_id AND ms.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.account_payable.cust_id=?", dataFilter.CustId)
	query.Where("acf.account_payable.cust_id=? ", dataFilter.CustId)
	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.account_payable.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.account_payable.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}
	if dataFilter.Query != "" {
		// queryCount.Where("acf.account_payable.invoice_no=?", dataFilter.Query)
		// query.Where("acf.account_payable.invoice_no=?", dataFilter.Query)

		query.Where("acf.account_payable.invoice_no LIKE ?", "%"+dataFilter.Query+"%")
		query.Or("acf.account_payable.invoice_no LIKE ?", "%"+dataFilter.Query+"%")
	}

	if dataFilter.InvoiceNo != "" {
		queryCount.Where(" acf.account_payable.invoice_no LIKE ?", "%"+dataFilter.InvoiceNo+"%")
		query.Where("acf.account_payable.invoice_no LIKE ?", "%"+dataFilter.InvoiceNo+"%")
	}

	if dataFilter.Supplier != 0 {
		queryCount.Where("acf.account_payable.sup_id=?", dataFilter.Supplier)
		query.Where("acf.account_payable.sup_id=?", dataFilter.Supplier)
	}

	query.Group(`acf.account_payable.invoice_date,
	app.account_payable_payment_date,
	acf.account_payable.due_date,
	acf.account_payable.invoice_no,
	appd.invoice_amount,
	ms.sup_id,
	ms.sup_code,
	ms.sup_name,
	sc.distributor_id,
	md.distributor_code,
	md.distributor_name,
	us.user_fullname,
	acf.account_payable.created_at,
	acf.account_payable.updated_at,
	acf.account_payable.amount,
	acf.account_payable.total`)

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
		query.Order("inv_date DESC")
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
