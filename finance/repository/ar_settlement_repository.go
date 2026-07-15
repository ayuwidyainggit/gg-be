package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryArSettlementImpl struct {
		*gorm.DB
	}
)

type ArSettlementRepository interface {
	FindOneByDepositNo(depositNo string, custId string) (arSettlement model.ArSettlementList, err error)
	FindDetail(depositNo string, custId string) (Details []model.ArSettlementPayment, err error)
	FindAllByCustId(dataFilter entity.ArSettlementQueryFilter) ([]model.ArSettlementList, int64, int, error)
	FindAllByCustIdNew(dataFilter entity.ArSettlementQueryFilter) ([]model.ArBranchSettlementList, int64, int, error)
	Approve(c context.Context, custId string, depositNo string, approvedBy int64) error
	Reject(c context.Context, custId string, depositNo string, rejectedBy int64) error
	ApproveBranch(c context.Context, custId string, depositNo string, approvedBy int64) error
	RejectBranch(c context.Context, custId string, depositNo string, rejectedBy int64) error

	FindAllCollectorByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (collector []model.SettlementCollectorFilter, total int64, lastPage int, err error)
	FindAllDepositStatusLookupMode(dataFilter entity.GeneralQueryFilter) (depositStatus []model.SettlementDepositStatusFilter, total int64, lastPage int, err error)
	FindDetailByInvoice(c context.Context, invoiceNo string, statuses []int, custId string) (Details []model.DepositDetailByInvoice, err error)
	SetInvoiceToPaidOff(c context.Context, invoiceNo []string, custId string) error

	FindBranchDetail(depositNo string, custId string) (Details []model.ArBranchSettlementPayment, err error)
	FindOneByBranchDepositNo(depositNo string, custId string) (arSettlement model.ArBranchSettlementList, err error)
	SumInvoiceRemainingForDepositSettlement(depositNo string, custId string) (float64, error)
	// FindBranchDetail(depositNo string, custId string) (Details []model.ArSettlementPayment, err error)
	FindBranchDetailByInvoice(c context.Context, invoiceNo string, statuses []int, custId string) (Details []model.DepositBranchDetailByInvoice, err error)
	SetBranchInvoiceToPaidOff(c context.Context, invoiceNo []string, custId string) error

	VerifyRejectData(c context.Context, depositNo string, custId string) (entity.RejectVerifyReport, error)
}

func NewArSettlementRepo(db *gorm.DB) *RepositoryArSettlementImpl {
	return &RepositoryArSettlementImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryArSettlementImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryArSettlementImpl) FindOneByDepositNo(depositNo string, custId string) (arSettlement model.ArSettlementList, err error) {
	err = repository.Select(`
		acf.deposit.deposit_date,
		acf.deposit.deposit_no,
		collection.collection_date,
		acf.deposit.collection_no,
		acf.deposit.deposit_status,
		acf.deposit.total_discount,
		acf.deposit.total_materai,
		acf.deposit.total_payment_balance,
		acf.deposit.total_payment,
		CASE
			WHEN acf.deposit.collection_no IS NOT NULL THEN COALESCE(collection.remaining_amount, acf.deposit.remaining_amount)
			ELSE acf.deposit.remaining_amount
		END as remaining_amount,
		acf.deposit.approved_by,
		cust.cust_id,
		cust.cust_name,
		ot_grp.ot_grp_id,
		ot_grp.ot_grp_code,
		ot_grp.ot_grp_name,
		employee.emp_id,
		employee.emp_code,
		employee.emp_name,
		approver_user.user_id,
		approver_user.user_fullname as approved_by_name
		`).
		Joins("left join smc.m_customer cust on cust.cust_id = acf.deposit.cust_id").
		Joins("left join mst.m_employee employee on employee.emp_id = acf.deposit.emp_id AND employee.cust_id = ?", custId).
		Joins("left join acf.collection collection on collection.collection_no = acf.deposit.collection_no AND collection.cust_id = ?", custId).
		Joins("left join mst.m_outlet_group ot_grp on ot_grp.ot_grp_id = collection.ot_grp_id").
		Joins("left join sys.m_user approver_user on approver_user.user_id = acf.deposit.approved_by").
		Where("acf.deposit.deposit_no = ? AND acf.deposit.cust_id=?", depositNo, custId).
		Take(&arSettlement).Error
	return arSettlement, err
}

func (repository *RepositoryArSettlementImpl) FindDetail(depositNo string, custId string) (Details []model.ArSettlementPayment, err error) {
	// err = repository.Select(`
	// 	acf.deposit_detail.deposit_detail_id as deposit_detail_id,
	// 	acf.deposit_detail.pay_type,
	// 	acf.deposit_detail.document_no,
	// 	acf.deposit_detail.invoice_no,
	// 	acf.deposit_detail.invoice_amount,
	// 	acf.deposit_detail.discount,
	// 	acf.deposit_detail.materai,
	// 	acf.deposit_detail.payment_balance,
	// 	acf.deposit_detail.total_payment,
	// 	acf.deposit_detail.remaining_payment,
	// 	invoice.invoice_date,
	// 	employee.emp_id as salesman_id,
	// 	employee.emp_code as salesman_code,
	// 	employee.emp_name as salesman_name,
	// 	ot.outlet_id as outlet_id,
	// 	ot.outlet_code as outlet_code,
	// 	ot.outlet_name as outlet_name
	// 	`).
	// 	Joins("left join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ?", custId).
	// 	Joins("left join sls.order invoice on invoice.invoice_no = acf.deposit_detail.invoice_no AND invoice.cust_id = ?", custId).
	// 	Joins("left join mst.m_employee employee on employee.emp_id = invoice.salesman_id AND employee.cust_id = ?", custId).
	// 	Joins("left join mst.m_outlet ot on ot.outlet_id = invoice.outlet_id AND ot.cust_id = ?", custId).
	// 	Where("acf.deposit_detail.deposit_no = ?", depositNo).
	// 	Find(&Details).Error

	err = repository.Select(`
		cust.cust_id,
		cust.cust_name,
		acf.deposit_payment.deposit_payment_id,
		acf.deposit_payment.invoice_no,
		acf.deposit_payment.pay_type,
		acf.deposit_payment.document_no,
		acf.deposit_payment.balance,
		acf.deposit_payment.payment_amount,
		invoice.invoice_date,
		employee.emp_id as salesman_id,
		employee.emp_code as salesman_code,
		employee.emp_name as salesman_name,
		ot.outlet_id as outlet_id,
		ot.outlet_code as outlet_code,
		ot.outlet_name as outlet_name,
		deposit_detail.discount,
		deposit_detail.payment_balance,
		deposit_detail.materai,
		deposit_detail.total_payment,
		GREATEST(COALESCE(invoice.total, 0) - COALESCE(paid_invoices.paid_amount, 0), 0) as remaining_payment
		`).
		Joins("left join acf.deposit_detail deposit_detail on deposit_detail.deposit_no = acf.deposit_payment.deposit_no AND deposit_detail.invoice_no = acf.deposit_payment.invoice_no AND deposit_detail.cust_id = acf.deposit_payment.cust_id").
		Joins("left join smc.m_customer cust on cust.cust_id = acf.deposit_payment.cust_id").
		Joins("left join sls.order invoice on invoice.invoice_no = acf.deposit_payment.invoice_no AND invoice.cust_id = ?", custId).
		Joins(`left join (
			select dd.invoice_no,
				dd.cust_id,
				coalesce(sum(dd.total_payment), 0) as paid_amount
			from acf.deposit_detail dd
			inner join acf.deposit d
				on d.deposit_no = dd.deposit_no
				and d.cust_id = dd.cust_id
			where dd.cust_id = ?
				and d.deposit_status in ?
			group by dd.invoice_no, dd.cust_id
		) paid_invoices on paid_invoices.invoice_no = acf.deposit_payment.invoice_no AND paid_invoices.cust_id = acf.deposit_payment.cust_id`, custId, []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}).
		Joins("left join mst.m_employee employee on employee.emp_id = invoice.salesman_id AND employee.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = invoice.outlet_id AND ot.cust_id = ?", custId).
		Where("acf.deposit_payment.deposit_no = ?", depositNo).
		Where("acf.deposit_payment.cust_id = ?", custId).
		Find(&Details).Error

	return Details, err
}

func (repository *RepositoryArSettlementImpl) SumInvoiceRemainingForDepositSettlement(depositNo string, custId string) (float64, error) {
	statuses := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}
	var sum float64
	err := repository.Raw(`
		SELECT COALESCE(SUM(per_inv.remaining), 0)
		FROM (
			SELECT
				dp.invoice_no,
				GREATEST(COALESCE(MAX(inv.total), 0) - COALESCE(MAX(pi.paid_amount), 0), 0) AS remaining
			FROM acf.deposit_payment dp
			LEFT JOIN sls.order inv ON inv.invoice_no = dp.invoice_no AND inv.cust_id = ?
			LEFT JOIN (
				SELECT dd.invoice_no,
					dd.cust_id,
					COALESCE(SUM(dd.total_payment), 0) AS paid_amount
				FROM acf.deposit_detail dd
				INNER JOIN acf.deposit d ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id
				WHERE dd.cust_id = ?
					AND d.deposit_status IN ?
				GROUP BY dd.invoice_no, dd.cust_id
			) pi ON pi.invoice_no = dp.invoice_no AND pi.cust_id = dp.cust_id
			WHERE dp.deposit_no = ?
				AND dp.cust_id = ?
			GROUP BY dp.invoice_no, dp.cust_id
		) per_inv
	`, custId, custId, statuses, depositNo, custId).Scan(&sum).Error
	return sum, err
}

func (repository *RepositoryArSettlementImpl) FindAllByCustId(dataFilter entity.ArSettlementQueryFilter) ([]model.ArSettlementList, int64, int, error) {
	var settlements []model.ArSettlementList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("deposit_no")
	query := repository.Select(`	
			acf.deposit.cust_id, 
			acf.deposit.deposit_no, 
			acf.deposit.deposit_date, 
			acf.deposit.total_payment,
			CASE
				WHEN acf.deposit.collection_no IS NOT NULL THEN COALESCE(collection.remaining_amount, acf.deposit.remaining_amount)
				ELSE acf.deposit.remaining_amount
			END as remaining_amount, 
			acf.deposit.deposit_status, 
			collector.emp_id, 
			collector.emp_code, 
			collector.emp_name,
			approver_user.user_id as approved_by,  
			approver_user.user_fullname as approved_by_name
		`).
		Joins("left join mst.m_employee collector on collector.emp_id = acf.deposit.salesman_id AND collector.cust_id = ?", dataFilter.CustId).
		Joins("left join acf.collection collection on collection.collection_no = acf.deposit.collection_no AND collection.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user approver_user on approver_user.user_id = acf.deposit.approved_by AND approver_user.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)

	if dataFilter.DepositDateFrom != nil && dataFilter.DepositDateTo != nil {
		query.Where("acf.deposit.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
		queryCount.Where("acf.deposit.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.deposit.deposit_no=?", dataFilter.Query)
		query.Where("acf.deposit.deposit_no=?", dataFilter.Query)
	}

	if len(dataFilter.EmpId) > 0 {
		queryCount.Where("acf.deposit.emp_id in ?", dataFilter.EmpId)
		query.Where("acf.deposit.emp_id in ?", dataFilter.EmpId)
	}

	if len(dataFilter.DepositStatus) > 0 {
		queryCount.Where("acf.deposit.deposit_status in ?", dataFilter.DepositStatus)
		query.Where("acf.deposit.deposit_status in ?", dataFilter.DepositStatus)
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
		query.Order("acf.deposit.deposit_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&settlements).Error
	if err != nil {
		return settlements, total, 0, err
	}
	err = queryCount.Model(&settlements).Count(&total).Error
	if err != nil {
		return settlements, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return settlements, total, lastPage, nil
}

func (repository *RepositoryArSettlementImpl) FindAllByCustIdUnion(dataFilter entity.ArSettlementQueryFilter) ([]model.ArSettlementList, int64, int, error) {
	var settlements []model.ArSettlementList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Model(&model.ArSettlement{}).Select("deposit_no")
	query := repository.Model(&model.ArSettlement{}).Select(`
			acf.deposit.cust_id, 
			acf.deposit.deposit_no, 
			acf.deposit.deposit_date, 
			acf.deposit.total_payment,
			CASE
				WHEN acf.deposit.collection_no IS NOT NULL THEN COALESCE(collection.remaining_amount, acf.deposit.remaining_amount)
				ELSE acf.deposit.remaining_amount
			END as remaining_amount, 
			acf.deposit.deposit_status, 
			collector.emp_id, 
			collector.emp_code, 
			collector.emp_name,
			approver_user.user_id as approved_by,  
			approver_user.user_fullname as approved_by_name
		`).
		Joins("left join mst.m_employee collector on collector.emp_id = acf.deposit.salesman_id AND collector.cust_id = ?", dataFilter.CustId).
		Joins("left join acf.collection collection on collection.collection_no = acf.deposit.collection_no AND collection.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user approver_user on approver_user.user_id = acf.deposit.approved_by AND approver_user.cust_id = ?", dataFilter.CustId)

	queryBranchCount := repository.Model(&model.ArBranchSettlement{}).Select("deposit_no")
	queryBranch := repository.Model(&model.ArBranchSettlement{}).Select(`
			inv.gr_branch_payment.cust_id, 
			inv.gr_branch_payment.deposit_no, 
			inv.gr_branch_payment.deposit_date, 
			inv.gr_branch_payment.total_payment,
			CASE
				WHEN inv.gr_branch_payment.verification_status = `+strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED)+` THEN 0
				ELSE inv.gr_branch_payment.total_payment
			END 
			as remaining_amount,
			--inv.gr_branch_payment.total_payment as remaining_amount,
			--(gr.total - coalesce(paid_invoices.paid_amount, 0)) as remaining_amount,
			inv.gr_branch_payment.verification_status as deposit_status, 
			collector.emp_id, 
			collector.emp_code, 
			collector.emp_name,
			approver_user.user_id as approved_by,  
			approver_user.user_fullname as approved_by_name
		`).
		Joins("left join mst.m_employee collector on collector.emp_id = 0 AND collector.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user approver_user on approver_user.user_id = inv.gr_branch_payment.verified_by AND approver_user.cust_id = ?", dataFilter.CustId).
		Joins("left join inv.gr_branch gr on gr.invoice_no_branch = inv.gr_branch_payment.invoice_no_branch AND gr.cust_id = inv.gr_branch_payment.cust_id")
		// Joins(`
		// left join (
		// 	select
		// 		inv.gr_branch_payment.invoice_no_branch,
		// 		inv.gr_branch_payment.cust_id,
		// 		(coalesce(sum(inv.gr_branch_payment.payment_amount), 0) + coalesce(sum(inv.gr_branch_payment.discount), 0)) as paid_amount
		// 	from inv.gr_branch_payment
		// 	left join inv.gr_branch on inv.gr_branch_payment.invoice_no_branch = inv.gr_branch.invoice_no_branch AND gr_branch.cust_id = inv.gr_branch_payment.cust_id
		// 	where inv.gr_branch_payment.verification_status in (` + strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED) + `)
		// 	group by inv.gr_branch_payment.invoice_no_branch, inv.gr_branch_payment.cust_id
		// ) paid_invoices on paid_invoices.invoice_no_branch = gr.invoice_no_branch AND gr.cust_id = paid_invoices.cust_id
		// `)

		// querySelect := fmt.Sprintf(`SELECT * FROM (? UNION ALL ?) AS combined`, query, queryBranch)
		// querySelectCount := fmt.Sprintf(`SELECT COUNT(*) as total FROM (? UNION ALL ?) AS combined`, queryCount, queryBranchCount)

		// queryUnionCount := repository.Raw("? UNION ALL ?", queryCount, queryBranchCount)
		// queryUnion := repository.Raw("? UNION ALL ?", query, queryBranch)

		// if err := .Scan(&settlements).Error; err != nil {
		// 	return settlements, total, 0, err
		// }

		// if err := repository.Raw("? UNION ALL ?", queryCount, queryBranchCount).Count(&total).Error; err != nil {
		// 	return settlements, total, 0, err
		// }

		// querySelect = querySelect + fmt.Sprintf(` WHERE`, query, queryBranch)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	queryBranch.Where("inv.gr_branch_payment.cust_id=?", dataFilter.CustId)
	queryBranchCount.Where("inv.gr_branch_payment.cust_id=?", dataFilter.CustId)

	if dataFilter.DepositDateFrom != nil && dataFilter.DepositDateTo != nil {
		query.Where("acf.deposit.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
		queryCount.Where("acf.deposit.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
		queryBranch.Where("inv.gr_branch_payment.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
		queryBranchCount.Where("inv.gr_branch_payment.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
	}

	if dataFilter.Query != "" {
		query.Where("acf.deposit.deposit_no=?", dataFilter.Query)
		queryCount.Where("acf.deposit.deposit_no=?", dataFilter.Query)
		queryBranch.Where("inv.gr_branch_payment.deposit_no=?", dataFilter.Query)
		queryBranchCount.Where("inv.gr_branch_payment.deposit_no=?", dataFilter.Query)
	}

	if len(dataFilter.EmpId) > 0 {
		query.Where("acf.deposit.emp_id in ?", dataFilter.EmpId)
		queryCount.Where("acf.deposit.emp_id in ?", dataFilter.EmpId)
		// queryBranch.Where("inv.gr_branch_payment.emp_id in ?", dataFilter.EmpId)
		// queryBranchCount.Where("inv.gr_branch_payment.emp_id in ?", dataFilter.EmpId)
	}

	if len(dataFilter.DepositStatus) > 0 {
		query.Where("acf.deposit.deposit_status in ?", dataFilter.DepositStatus)
		queryCount.Where("acf.deposit.deposit_status in ?", dataFilter.DepositStatus)
		queryBranch.Where("inv.gr_branch_payment.verification_status in ?", dataFilter.DepositStatus)
		queryBranchCount.Where("inv.gr_branch_payment.verification_status in ?", dataFilter.DepositStatus)
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
		// queryUnion.Order(sortBy)
	} else {
		sortBy = "deposit_no DESC"
		// queryUnion.Order("deposit_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	// err := query.Limit(limit).Offset(offset).Find(&settlements).Error
	// if err != nil {
	// 	return settlements, total, 0, err
	// }
	// err = queryCount.Model(&settlements).Count(&total).Error
	// if err != nil {
	// 	return settlements, total, 0, err
	// }

	queryUnionCount := repository.Raw("SELECT COUNT(*) as total FROM (? UNION ALL ?) AS combined", queryCount, queryBranchCount)
	queryUnion := repository.Raw("SELECT * FROM (? UNION ALL ?) AS combined ORDER BY ? LIMIT ? OFFSET ?", query, queryBranch, sortBy, limit, offset)

	if err := queryUnion.Scan(&settlements).Error; err != nil {
		return settlements, total, 0, err
	}
	fmt.Println("OKe cust")
	if err := queryUnionCount.Scan(&total).Error; err != nil {
		return settlements, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return settlements, total, lastPage, nil
}

func (repository *RepositoryArSettlementImpl) FindAllByCustIdNew(dataFilter entity.ArSettlementQueryFilter) ([]model.ArBranchSettlementList, int64, int, error) {
	var settlements []model.ArBranchSettlementList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryPaidInvoices := `left join (
			select inv.gr_branch_payment.invoice_no_branch, 
			inv.gr_branch_payment.cust_id,
			(coalesce(sum(inv.gr_branch_payment.payment_amount), 0) + coalesce(sum(inv.gr_branch_payment.discount), 0) + coalesce(sum(inv.gr_branch_payment.payment_balance), 0)) as paid_amount
		from inv.gr_branch_payment
		left join inv.gr_branch on inv.gr_branch_payment.invoice_no_branch = inv.gr_branch.invoice_no_branch AND inv.gr_branch.cust_id = inv.gr_branch_payment.cust_id
		where inv.gr_branch_payment.verification_status = ` + strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED) + `
		group by inv.gr_branch_payment.invoice_no_branch, inv.gr_branch_payment.cust_id
	) paid_invoices on paid_invoices.invoice_no_branch = gr.invoice_no_branch AND gr.cust_id = paid_invoices.cust_id`

	queryCount := repository.Select("deposit_no")
	query := repository.Select(`
			inv.gr_branch_payment.cust_id, 
			inv.gr_branch_payment.deposit_no, 
			inv.gr_branch_payment.deposit_date, 
			inv.gr_branch_payment.total_payment,
			CASE
				WHEN (gr.total - coalesce(paid_invoices.paid_amount, 0)) < 0 THEN 0
				ELSE (gr.total - coalesce(paid_invoices.paid_amount, 0))
			END 
			as remaining_amount,
			--inv.gr_branch_payment.total_payment as remaining_amount,
			--(gr.total - coalesce(paid_invoices.paid_amount, 0)) as remaining_amount,
			inv.gr_branch_payment.verification_status as deposit_status, 
			collector.emp_id, 
			collector.emp_code, 
			collector.emp_name,
			approver_user.user_id as approved_by,  
			approver_user.user_fullname as approved_by_name
		`).
		Joins("left join mst.m_employee collector on collector.emp_id = 0 AND collector.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user approver_user on approver_user.user_id = inv.gr_branch_payment.verified_by").
		Joins("left join inv.gr_branch gr on gr.invoice_no_branch = inv.gr_branch_payment.invoice_no_branch AND gr.cust_id = inv.gr_branch_payment.cust_id").
		Joins("inner join smc.m_customer cust on cust.cust_id = inv.gr_branch_payment.cust_id AND cust.parent_cust_id = ?", dataFilter.ParentCustId).
		Joins(queryPaidInvoices)

	// queryCount.Where("inv.gr_branch_payment.cust_id=?", dataFilter.CustId)
	// query.Where("inv.gr_branch_payment.cust_id=?", dataFilter.CustId)

	if dataFilter.DepositDateFrom != nil && dataFilter.DepositDateTo != nil {
		query.Where("inv.gr_branch_payment.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
		queryCount.Where("inv.gr_branch_payment.deposit_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.DepositDateFrom), str.UnixTimestampToUtcTime(*dataFilter.DepositDateTo))
	}

	if dataFilter.Query != "" {
		queryCount.Where("inv.gr_branch_payment.deposit_no=?", dataFilter.Query)
		query.Where("inv.gr_branch_payment.deposit_no=?", dataFilter.Query)
	}

	// if len(dataFilter.EmpId) > 0 {
	// 	queryCount.Where("acf.deposit.emp_id in ?", dataFilter.EmpId)
	// 	query.Where("acf.deposit.emp_id in ?", dataFilter.EmpId)
	// }

	if len(dataFilter.DepositStatus) > 0 {
		queryCount.Where("inv.gr_branch_payment.verification_status in ?", dataFilter.DepositStatus)
		query.Where("inv.gr_branch_payment.verification_status in ?", dataFilter.DepositStatus)
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
		query.Order("inv.gr_branch_payment.deposit_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&settlements).Error
	if err != nil {
		return settlements, total, 0, err
	}
	err = queryCount.Model(&settlements).Count(&total).Error
	if err != nil {
		return settlements, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return settlements, total, lastPage, nil
}

func (repository *RepositoryArSettlementImpl) FindAllCollectorByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.SettlementCollectorFilter, int64, int, error) {
	var collectors []model.SettlementCollectorFilter
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
		Select(`
			d.emp_id AS emp_id,
			COALESCE(NULLIF(emp.emp_code, ''), CAST(d.emp_id AS varchar)) AS emp_code,
			COALESCE(NULLIF(emp.emp_name, ''), COALESCE(NULLIF(emp.emp_code, ''), CAST(d.emp_id AS varchar))) AS emp_name
		`).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = d.emp_id AND emp.cust_id = ?", dataFilter.CustId).
		Where("d.cust_id = ?", dataFilter.CustId).
		Where("d.deleted_at IS NULL").
		Where("d.emp_id IS NOT NULL").
		Group("d.emp_id, emp.emp_code, emp.emp_name")

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

	sortBy := "d.emp_id DESC"
	if strings.TrimSpace(dataFilter.Sort) != "" {
		allowedFields := map[string]string{
			"emp_id":   "d.emp_id",
			"emp_code": "emp.emp_code",
			"emp_name": "emp.emp_name",
		}

		orderClauses := make([]string, 0)
		for _, rawSort := range strings.Split(dataFilter.Sort, ",") {
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

		if len(orderClauses) > 0 {
			sortBy = strings.Join(orderClauses, ", ")
		}
	}

	offset := (page - 1) * limit
	err := query.Order(sortBy).Limit(limit).Offset(offset).Scan(&collectors).Error
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

func (repository *RepositoryArSettlementImpl) FindAllDepositStatusLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.SettlementDepositStatusFilter, int64, int, error) {

	var depositStatuses []model.SettlementDepositStatusFilter

	var total int64

	queryCount := repository.Select("acf.deposit.deposit_status")
	query := repository.Select(`acf.deposit.deposit_status`)

	queryCount.Where("acf.deposit.cust_id=?", dataFilter.CustId)
	query.Where("acf.deposit.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.deposit.deposit_status IS NOT NULL")
	query.Where("acf.deposit.deposit_status IS NOT NULL")

	queryCount.Group("acf.deposit.deposit_status")
	query.Group("acf.deposit.deposit_status")

	query.Order("acf.deposit.deposit_status ASC")

	err := query.Find(&depositStatuses).Error
	if err != nil {
		return depositStatuses, total, 0, err
	}

	total = int64(len(depositStatuses))
	lastPage := 1
	return depositStatuses, total, lastPage, nil
}

func (repository *RepositoryArSettlementImpl) Approve(c context.Context, custId string, depositNo string, approvedBy int64) error {
	var data model.ArSettlement
	result := repository.model(c).Model(&data).Where("deposit_no=?", depositNo).Where("cust_id = ?", custId).
		Updates(map[string]interface{}{"deposit_status": 2, "is_approved": true, "approved_by": approvedBy, "approved_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArSettlementImpl) Reject(c context.Context, custId string, depositNo string, rejectedBy int64) error {
	db := repository.model(c)
	// 1. Revert expense: acf.expense.balance += deposit_expense.payment_amount
	if err := repository.revertExpenseBalanceForDeposit(c, depositNo, rejectedBy); err != nil {
		return err
	}
	// 2. Revert cheque/giro (pay_type=2), bank_transfer (3), return (4), cndn (5)
	if err := repository.revertChequeGiroForDeposit(c, depositNo, custId, rejectedBy); err != nil {
		return err
	}
	if err := repository.revertBankTransferForDeposit(c, depositNo, custId, rejectedBy); err != nil {
		return err
	}
	if err := repository.revertReturnForDeposit(c, depositNo, custId, rejectedBy); err != nil {
		return err
	}
	if err := repository.revertCndnForDeposit(c, depositNo, custId, rejectedBy); err != nil {
		return err
	}
	var data model.ArSettlement
	result := db.Model(&data).Where("deposit_no=?", depositNo).Where("cust_id = ?", custId).
		Updates(map[string]interface{}{"deposit_status": 3, "is_approved": false, "approved_by": nil, "approved_at": nil})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArSettlementImpl) revertExpenseBalanceForDeposit(c context.Context, depositNo string, updatedBy int64) error {
	return repository.model(c).Exec(`
		UPDATE acf.expense e
		SET balance = e.balance + de.payment_amount,
			updated_at = NOW(),
			updated_by = ?
		FROM acf.deposit_expense de
		WHERE de.deposit_no = ? AND de.expense_id = e.expense_id
	`, updatedBy, depositNo).Error
}

func (repository *RepositoryArSettlementImpl) revertChequeGiroForDeposit(c context.Context, depositNo string, custId string, updatedBy int64) error {
	return repository.model(c).Exec(`
		UPDATE acf.cheque_giro cg
		SET paid_amount = cg.paid_amount - sub.sum_amt,
			remaining_amount = cg.remaining_amount + sub.sum_amt,
			updated_at = NOW(),
			updated_by = ?
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 2
			GROUP BY document_no
		) sub
		WHERE cg.doc_no_cheque = sub.document_no AND cg.cust_id = ?
	`, updatedBy, depositNo, custId, custId).Error
}

func (repository *RepositoryArSettlementImpl) revertBankTransferForDeposit(c context.Context, depositNo string, custId string, updatedBy int64) error {
	return repository.model(c).Exec(`
		UPDATE acf.bank_transfer bt
		SET paid_amount = bt.paid_amount - sub.sum_amt,
			remaining_amount = bt.remaining_amount + sub.sum_amt,
			updated_at = NOW(),
			updated_by = ?
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 3
			GROUP BY document_no
		) sub
		WHERE bt.doc_no_bank = sub.document_no AND bt.cust_id = ?
	`, updatedBy, depositNo, custId, custId).Error
}

func (repository *RepositoryArSettlementImpl) revertReturnForDeposit(c context.Context, depositNo string, custId string, updatedBy int64) error {
	return repository.model(c).Exec(`
		UPDATE sls.return r
		SET paid_amount = r.paid_amount - sub.sum_amt,
			remaining_amount = r.remaining_amount + sub.sum_amt,
			updated_at = NOW(),
			updated_by = ?
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 4
			GROUP BY document_no
		) sub
		WHERE r.return_no = sub.document_no AND r.cust_id = ?
	`, updatedBy, depositNo, custId, custId).Error
}

func (repository *RepositoryArSettlementImpl) revertCndnForDeposit(c context.Context, depositNo string, custId string, updatedBy int64) error {
	return repository.model(c).Exec(`
		UPDATE acf.cndn cn
		SET used_amount = COALESCE(cn.used_amount, 0) - sub.sum_amt,
			remaning_amount = COALESCE(cn.remaning_amount, 0) + sub.sum_amt,
			updated_at = NOW(),
			updated_by = ?
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 5
			GROUP BY document_no
		) sub
		WHERE cn.cndn_no = sub.document_no AND cn.cust_id = ?
	`, updatedBy, depositNo, custId, custId).Error
}

func (repository *RepositoryArSettlementImpl) ApproveBranch(c context.Context, custId string, depositNo string, approvedBy int64) error {
	var data model.ArBranchSettlement
	result := repository.model(c).Model(&data).Where("deposit_no=?", depositNo).Where("cust_id = ?", custId).
		Updates(map[string]interface{}{"verification_status": 2, "verified_by": approvedBy, "verified_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArSettlementImpl) RejectBranch(c context.Context, custId string, depositNo string, rejectedBy int64) error {
	var data model.ArBranchSettlement
	result := repository.model(c).Model(&data).Where("deposit_no=?", depositNo).Where("cust_id = ?", custId).
		Updates(map[string]interface{}{"verification_status": 3, "verified_by": rejectedBy, "verified_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArSettlementImpl) VerifyRejectData(c context.Context, depositNo string, custId string) (entity.RejectVerifyReport, error) {
	out := entity.RejectVerifyReport{DepositNo: depositNo, CustId: custId}

	var expenseTotal float64
	err := repository.model(c).Raw(`
		SELECT COALESCE(SUM(de.payment_amount), 0)
		FROM acf.deposit_expense de
		WHERE de.deposit_no = ?
	`, depositNo).Scan(&expenseTotal).Error
	if err != nil {
		return out, err
	}
	out.Expense.TotalPaymentAmount = expenseTotal

	var expenseRows []struct {
		ExpenseId      int64   `gorm:"column:expense_id"`
		PaymentAmount  float64 `gorm:"column:payment_amount"`
		CurrentBalance float64 `gorm:"column:current_balance"`
		ExpenseExists  bool    `gorm:"column:expense_exists"`
	}
	err = repository.model(c).Raw(`
		SELECT de.expense_id,
		       de.payment_amount,
		       COALESCE(e.balance, 0) AS current_balance,
		       (e.expense_id IS NOT NULL) AS expense_exists
		FROM acf.deposit_expense de
		LEFT JOIN acf.expense e ON e.expense_id = de.expense_id AND e.cust_id = ?
		WHERE de.deposit_no = ?
	`, custId, depositNo).Scan(&expenseRows).Error
	if err != nil {
		return out, err
	}
	for _, r := range expenseRows {
		out.Expense.Items = append(out.Expense.Items, entity.RejectVerifyExpenseItem{
			ExpenseId:      r.ExpenseId,
			PaymentAmount:  r.PaymentAmount,
			CurrentBalance: r.CurrentBalance,
			ExpenseExists:  r.ExpenseExists,
		})
	}

	var cgRows []struct {
		DocumentNo             string  `gorm:"column:document_no"`
		SumAmount              float64 `gorm:"column:sum_amt"`
		CurrentPaidAmount      float64 `gorm:"column:paid_amount"`
		CurrentRemainingAmount float64 `gorm:"column:remaining_amount"`
		RowExists              bool    `gorm:"column:row_exists"`
	}
	err = repository.model(c).Raw(`
		SELECT sub.document_no, sub.sum_amt,
		       COALESCE(cg.paid_amount, 0) AS paid_amount,
		       COALESCE(cg.remaining_amount, 0) AS remaining_amount,
		       (cg.doc_no_cheque IS NOT NULL) AS row_exists
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 2
			GROUP BY document_no
		) sub
		LEFT JOIN acf.cheque_giro cg ON cg.doc_no_cheque = sub.document_no AND cg.cust_id = ?
	`, depositNo, custId, custId).Scan(&cgRows).Error
	if err != nil {
		return out, err
	}
	for _, r := range cgRows {
		out.ChequeGiro = append(out.ChequeGiro, entity.RejectVerifyDocAmount{
			DocumentNo:             r.DocumentNo,
			SumAmount:              r.SumAmount,
			CurrentPaidAmount:      r.CurrentPaidAmount,
			CurrentRemainingAmount: r.CurrentRemainingAmount,
			RowExists:              r.RowExists,
		})
	}

	var btRows []struct {
		DocumentNo             string  `gorm:"column:document_no"`
		SumAmount              float64 `gorm:"column:sum_amt"`
		CurrentPaidAmount      float64 `gorm:"column:paid_amount"`
		CurrentRemainingAmount float64 `gorm:"column:remaining_amount"`
		RowExists              bool    `gorm:"column:row_exists"`
	}
	err = repository.model(c).Raw(`
		SELECT sub.document_no, sub.sum_amt,
		       COALESCE(bt.paid_amount, 0) AS paid_amount,
		       COALESCE(bt.remaining_amount, 0) AS remaining_amount,
		       (bt.doc_no_bank IS NOT NULL) AS row_exists
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 3
			GROUP BY document_no
		) sub
		LEFT JOIN acf.bank_transfer bt ON bt.doc_no_bank = sub.document_no AND bt.cust_id = ?
	`, depositNo, custId, custId).Scan(&btRows).Error
	if err != nil {
		return out, err
	}
	for _, r := range btRows {
		out.BankTransfer = append(out.BankTransfer, entity.RejectVerifyDocAmount{
			DocumentNo:             r.DocumentNo,
			SumAmount:              r.SumAmount,
			CurrentPaidAmount:      r.CurrentPaidAmount,
			CurrentRemainingAmount: r.CurrentRemainingAmount,
			RowExists:              r.RowExists,
		})
	}

	var retRows []struct {
		DocumentNo             string  `gorm:"column:document_no"`
		SumAmount              float64 `gorm:"column:sum_amt"`
		CurrentPaidAmount      float64 `gorm:"column:paid_amount"`
		CurrentRemainingAmount float64 `gorm:"column:remaining_amount"`
		RowExists              bool    `gorm:"column:row_exists"`
	}
	err = repository.model(c).Raw(`
		SELECT sub.document_no, sub.sum_amt,
		       COALESCE(r.paid_amount, 0) AS paid_amount,
		       COALESCE(r.remaining_amount, 0) AS remaining_amount,
		       (r.return_no IS NOT NULL) AS row_exists
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 4
			GROUP BY document_no
		) sub
		LEFT JOIN sls.return r ON r.return_no = sub.document_no AND r.cust_id = ?
	`, depositNo, custId, custId).Scan(&retRows).Error
	if err != nil {
		return out, err
	}
	for _, r := range retRows {
		out.Return = append(out.Return, entity.RejectVerifyDocAmount{
			DocumentNo:             r.DocumentNo,
			SumAmount:              r.SumAmount,
			CurrentPaidAmount:      r.CurrentPaidAmount,
			CurrentRemainingAmount: r.CurrentRemainingAmount,
			RowExists:              r.RowExists,
		})
	}

	var cndnRows []struct {
		DocumentNo            string  `gorm:"column:document_no"`
		SumAmount             float64 `gorm:"column:sum_amt"`
		CurrentUsedAmount     float64 `gorm:"column:used_amount"`
		CurrentRemaningAmount float64 `gorm:"column:remaning_amount"`
		RowExists             bool    `gorm:"column:row_exists"`
	}
	err = repository.model(c).Raw(`
		SELECT sub.document_no, sub.sum_amt,
		       COALESCE(cn.used_amount, 0) AS used_amount,
		       COALESCE(cn.remaning_amount, 0) AS remaning_amount,
		       (cn.cndn_no IS NOT NULL) AS row_exists
		FROM (
			SELECT document_no, COALESCE(SUM(payment_amount), 0) AS sum_amt
			FROM acf.deposit_payment
			WHERE deposit_no = ? AND cust_id = ? AND pay_type = 5
			GROUP BY document_no
		) sub
		LEFT JOIN acf.cndn cn ON cn.cndn_no = sub.document_no AND cn.cust_id = ?
	`, depositNo, custId, custId).Scan(&cndnRows).Error
	if err != nil {
		return out, err
	}
	for _, r := range cndnRows {
		out.Cndn = append(out.Cndn, entity.RejectVerifyCndnAmount{
			DocumentNo:            r.DocumentNo,
			SumAmount:             r.SumAmount,
			CurrentUsedAmount:     r.CurrentUsedAmount,
			CurrentRemaningAmount: r.CurrentRemaningAmount,
			RowExists:             r.RowExists,
		})
	}

	payTypeNames := map[int]string{1: "Cash", 2: "Cheque", 3: "Transfer", 4: "Return", 5: "Debit/Credit"}
	var dpRows []struct {
		DocumentNo    string  `gorm:"column:document_no"`
		PayType       int     `gorm:"column:pay_type"`
		PaymentAmount float64 `gorm:"column:payment_amount"`
	}
	err = repository.model(c).Raw(`
		SELECT document_no, pay_type, COALESCE(SUM(payment_amount), 0) AS payment_amount
		FROM acf.deposit_payment
		WHERE deposit_no = ? AND cust_id = ? AND pay_type IN (2, 3, 4, 5)
		GROUP BY document_no, pay_type
	`, depositNo, custId).Scan(&dpRows).Error
	if err != nil {
		return out, err
	}
	for _, r := range dpRows {
		item := entity.RejectVerifyDepositPaymentItem{
			DocumentNo:    r.DocumentNo,
			PayType:       r.PayType,
			PayTypeName:   payTypeNames[r.PayType],
			PaymentAmount: r.PaymentAmount,
			PayTypeOK:     r.PayType >= 2 && r.PayType <= 5,
		}
		switch r.PayType {
		case 2:
			var n int64
			repository.model(c).Raw(`SELECT COUNT(*) FROM acf.cheque_giro WHERE doc_no_cheque = ? AND cust_id = ?`, r.DocumentNo, custId).Scan(&n)
			item.RowExists = n > 0
		case 3:
			var n int64
			repository.model(c).Raw(`SELECT COUNT(*) FROM acf.bank_transfer WHERE doc_no_bank = ? AND cust_id = ?`, r.DocumentNo, custId).Scan(&n)
			item.RowExists = n > 0
		case 4:
			var n int64
			repository.model(c).Raw(`SELECT COUNT(*) FROM sls.return WHERE return_no = ? AND cust_id = ?`, r.DocumentNo, custId).Scan(&n)
			item.RowExists = n > 0
		case 5:
			var n int64
			repository.model(c).Raw(`SELECT COUNT(*) FROM acf.cndn WHERE cndn_no = ? AND cust_id = ?`, r.DocumentNo, custId).Scan(&n)
			item.RowExists = n > 0
		default:
			item.RowExists = false
		}
		out.DepositPaymentValidation = append(out.DepositPaymentValidation, item)
	}

	return out, nil
}

func (repository *RepositoryArSettlementImpl) FindDetailByInvoice(c context.Context, invoiceNo string, statuses []int, custId string) (Details []model.DepositDetailByInvoice, err error) {

	err = repository.model(c).Select(`
		acf.deposit_detail.*
		`).
		Joins("left join acf.deposit d on d.deposit_no = acf.deposit_detail.deposit_no and d.cust_id = ?", custId).
		Where("acf.deposit_detail.invoice_no = ?", invoiceNo).
		Where("acf.deposit_detail.cust_id = ?", custId).
		Where("d.deposit_status in ?", statuses).
		Find(&Details).Error

	return Details, err
}

func (repository *RepositoryArSettlementImpl) SetInvoiceToPaidOff(c context.Context, invoiceNo []string, custId string) error {
	result := repository.model(c).Table("sls.order").Where("invoice_no in ?", invoiceNo).Where("cust_id = ?", custId).
		Updates(map[string]interface{}{"is_paid_off": true})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArSettlementImpl) FindBranchDetailByInvoice(c context.Context, invoiceNo string, statuses []int, custId string) (Details []model.DepositBranchDetailByInvoice, err error) {

	err = repository.model(c).Select(`
		inv.gr_branch_payment.*
		`).
		Joins("left join inv.gr_branch gr on gr.invoice_no_branch = inv.gr_branch_payment.invoice_no_branch").
		Where("inv.gr_branch_payment.invoice_no_branch = ?", invoiceNo).
		Where("inv.gr_branch_payment.cust_id = ?", custId).
		Where("inv.gr_branch_payment.verification_status in ?", statuses).
		Find(&Details).Error

	return Details, err
}

func (repository *RepositoryArSettlementImpl) SetBranchInvoiceToPaidOff(c context.Context, invoiceNo []string, custId string) error {
	result := repository.model(c).Table("inv.gr_branch").Where("invoice_no in ?", invoiceNo).Where("cust_id = ?", custId).
		Updates(map[string]interface{}{"is_paid_off": true})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArSettlementImpl) FindBranchDetail(depositNo string, custId string) (Details []model.ArBranchSettlementPayment, err error) {

	queryPaidInvoices := `left join (
			select inv.gr_branch_payment.invoice_no_branch, 
			inv.gr_branch_payment.cust_id,
			(coalesce(sum(inv.gr_branch_payment.payment_amount), 0) + coalesce(sum(inv.gr_branch_payment.discount), 0)) as paid_amount
		from inv.gr_branch_payment
		left join inv.gr_branch on inv.gr_branch_payment.invoice_no_branch = inv.gr_branch.invoice_no_branch AND inv.gr_branch.cust_id = inv.gr_branch_payment.cust_id
		where inv.gr_branch_payment.verification_status = ` + strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED) + `
		group by inv.gr_branch_payment.invoice_no_branch, inv.gr_branch_payment.cust_id
	) paid_invoices on paid_invoices.invoice_no_branch = gr.invoice_no_branch AND gr.cust_id = paid_invoices.cust_id`

	err = repository.Select(`
				inv.gr_branch_payment.cust_id,
				inv.gr_branch_payment.gr_branch_payment_id as deposit_payment_id,
				inv.gr_branch_payment.invoice_no_branch as invoice_no,
				gr.invoice_date_branch as invoice_date,
				inv.gr_branch_payment.payment_type as pay_type,
				--inv.gr_branch_payment.document_no,
				--inv.gr_branch_payment.balance,
				inv.gr_branch_payment.payment_amount,
				employee.emp_id as salesman_id,
				employee.emp_code as salesman_code,
				employee.emp_name as salesman_name,
				ot.outlet_id as outlet_id,
				ot.outlet_code as outlet_code,
				ot.outlet_name as outlet_name,
				inv.gr_branch_payment.discount,
				inv.gr_branch_payment.payment_balance,
				--inv.gr_branch_payment.materai,
				inv.gr_branch_payment.total_payment,
				CASE
					WHEN (gr.total - coalesce(paid_invoices.paid_amount, 0)) < 0 THEN 0
					ELSE (gr.total - coalesce(paid_invoices.paid_amount, 0))
				END 
				as remaining_payment
				--inv.gr_branch_payment.remaining_payment
			`).
		// Joins("left join acf.deposit_detail deposit_detail on deposit_detail.deposit_no = acf.deposit_payment.deposit_no AND deposit_detail.invoice_no = acf.deposit_payment.invoice_no").
		// Joins("left join acf.deposit deposit on deposit.deposit_no = acf.deposit_detail.deposit_no AND deposit.cust_id = ?", custId).
		// Joins("left join sls.order invoice on invoice.invoice_no = acf.deposit_payment.invoice_no AND invoice.cust_id = ?", custId).
		Joins("left join inv.gr_branch gr on gr.invoice_no_branch = inv.gr_branch_payment.invoice_no_branch AND gr.cust_id = ?", custId).
		Joins("left join mst.m_employee employee on employee.emp_id = 0 AND employee.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = 0 AND ot.cust_id = ?", custId).
		Joins(queryPaidInvoices).
		Where("inv.gr_branch_payment.deposit_no = ? and inv.gr_branch_payment.cust_id = ?", depositNo, custId).
		Find(&Details).Error

	return Details, err
}

func (repository *RepositoryArSettlementImpl) FindOneByBranchDepositNo(depositNo string, custId string) (arSettlement model.ArBranchSettlementList, err error) {
	queryPaidInvoices := `left join (
			select inv.gr_branch_payment.invoice_no_branch, 
			inv.gr_branch_payment.cust_id,
			(coalesce(sum(inv.gr_branch_payment.payment_amount), 0) + coalesce(sum(inv.gr_branch_payment.discount), 0) + coalesce(sum(inv.gr_branch_payment.payment_balance), 0)) as paid_amount
		from inv.gr_branch_payment
		left join inv.gr_branch on inv.gr_branch_payment.invoice_no_branch = inv.gr_branch.invoice_no_branch AND inv.gr_branch.cust_id = inv.gr_branch_payment.cust_id
		where inv.gr_branch_payment.verification_status = ` + strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED) + `
		group by inv.gr_branch_payment.invoice_no_branch, inv.gr_branch_payment.cust_id
	) paid_invoices on paid_invoices.invoice_no_branch = gr.invoice_no_branch AND gr.cust_id = paid_invoices.cust_id`

	err = repository.Select(`
			--inv.gr_branch_payment.gr_branch_payment_id as deposit_payment_id,
			--inv.gr_branch_payment.payment_type as pay_type,
			--gr.invoice_no_branch,
			--gr.invoice_date_branch,
			inv.gr_branch_payment.deposit_date,
			inv.gr_branch_payment.deposit_no,
			inv.gr_branch_payment.verification_status as deposit_status,
			inv.gr_branch_payment.payment_balance as total_payment_balance,
			inv.gr_branch_payment.total_payment,
			inv.gr_branch_payment.payment_balance as total_payment_balance,
			inv.gr_branch_payment.discount as total_discount,
			CASE
				WHEN (gr.total - coalesce(paid_invoices.paid_amount, 0)) < 0 THEN 0
				ELSE (gr.total - coalesce(paid_invoices.paid_amount, 0))
			END 
			as remaining_amount,
			inv.gr_branch_payment.verified_by as approved_by,
			collection.collection_date,
			collection.collection_no,
			cust.cust_id,
			cust.cust_name,
			ot_grp.ot_grp_id,
			ot_grp.ot_grp_code,
			ot_grp.ot_grp_name,
			employee.emp_id,
			employee.emp_code,
			employee.emp_name,
			approver_user.user_id,
			approver_user.user_fullname as approved_by_name
		`).
		Joins("left join inv.gr_branch gr on gr.invoice_no_branch = inv.gr_branch_payment.invoice_no_branch AND gr.cust_id = ?", custId).
		Joins("left join smc.m_customer cust on cust.cust_id = inv.gr_branch_payment.cust_id").
		Joins("left join mst.m_employee employee on employee.emp_id = 0 AND employee.cust_id = ?", custId).
		Joins("left join acf.collection collection on collection.collection_no = inv.gr_branch_payment.deposit_no AND collection.cust_id = ?", custId).
		Joins("left join mst.m_outlet_group ot_grp on ot_grp.ot_grp_id = 0").
		Joins("left join sys.m_user approver_user on approver_user.user_id = inv.gr_branch_payment.verified_by").
		Joins(queryPaidInvoices).
		Where("inv.gr_branch_payment.deposit_no = ? AND inv.gr_branch_payment.cust_id=?", depositNo, custId).
		Take(&arSettlement).Error
	return arSettlement, err
}
