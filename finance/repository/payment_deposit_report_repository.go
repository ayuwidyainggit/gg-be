package repository

import (
	"context"
	"finance/entity"
	"finance/model"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PaymentDepositReportRepository interface {
	FindAllPaymentDeposit(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error)
	FindPaymentDepositSummary(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error)
	FindAllPaymentDepositNoLimit(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, error)
	FindAllPaymentDepositDownload(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportDownloadRow, error)
	FindPaymentDepositRecapRows(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportRecapRow, error)
	InsertReportList(c context.Context, report model.ReportList) error
	UpdateReportList(c context.Context, reportID string, status int, fileBase64 string) error
	GetReportRunningNumber(custId string, date time.Time) (int, error)
}

type RepositoryPaymentDepositReportImpl struct {
	*gorm.DB
}

func NewPaymentDepositReportRepo(db *gorm.DB) *RepositoryPaymentDepositReportImpl {
	return &RepositoryPaymentDepositReportImpl{DB: db}
}

func (repo *RepositoryPaymentDepositReportImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repo *RepositoryPaymentDepositReportImpl) FindAllPaymentDeposit(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, int64, int, error) {
	countSQL, dataSQL, args := repo.buildCountAndDataQueries(dataFilter, custId)
	var total int64
	if err := repo.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 20
	}
	page := dataFilter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	dataSQL = fmt.Sprintf("%s LIMIT ? OFFSET ?", dataSQL)
	queryArgs := append(append([]interface{}{}, args...), limit, offset)

	var data []model.PaymentDepositReportRow
	if err := repo.Raw(dataSQL, queryArgs...).Scan(&data).Error; err != nil {
		return nil, 0, 0, err
	}

	return data, total, lastPage, nil
}

func (repo *RepositoryPaymentDepositReportImpl) FindPaymentDepositSummary(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (model.PaymentDepositReportSummaryRow, error) {
	sql, args := repo.buildQuery(dataFilter, custId)
	summarySQL := fmt.Sprintf(`SELECT
		COALESCE(SUM(cash_amount), 0) AS total_cash,
		COALESCE(SUM(cheque_amount), 0) AS total_cheque,
		COALESCE(SUM(transfer_amount), 0) AS total_transfer,
		COALESCE(SUM(return_amount), 0) AS total_return,
		COALESCE(SUM(credit_debit_amount), 0) AS total_credit_debit,
		COALESCE(SUM(expense_amount), 0) AS total_expense
	FROM (%s) t`, sql)

	var summary model.PaymentDepositReportSummaryRow
	if err := repo.Raw(summarySQL, args...).Scan(&summary).Error; err != nil {
		return model.PaymentDepositReportSummaryRow{}, err
	}

	return summary, nil
}

func (repo *RepositoryPaymentDepositReportImpl) FindAllPaymentDepositNoLimit(dataFilter entity.PaymentDepositReportQueryFilter, custId string) ([]model.PaymentDepositReportRow, error) {
	_, dataSQL, args := repo.buildCountAndDataQueries(dataFilter, custId)

	var data []model.PaymentDepositReportRow
	if err := repo.Raw(dataSQL, args...).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *RepositoryPaymentDepositReportImpl) FindAllPaymentDepositDownload(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportDownloadRow, error) {
	query, args := repo.buildDownloadQuery(dataFilter, custId, parentCustId)
	var rows []model.PaymentDepositReportDownloadRow
	if err := repo.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *RepositoryPaymentDepositReportImpl) FindPaymentDepositRecapRows(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) ([]model.PaymentDepositReportRecapRow, error) {
	unionQuery, args := repo.buildDownloadUnionQuery(dataFilter, custId, parentCustId)
	recapQuery := fmt.Sprintf(`SELECT
		deposit_type,
		COALESCE(SUM(cash), 0) AS cash,
		COALESCE(SUM(cheque_giro), 0) AS cheque_giro,
		COALESCE(SUM(transfer), 0) AS transfer,
		COALESCE(SUM(return_amount), 0) AS return_amount,
		COALESCE(SUM(credit_debit), 0) AS credit_debit,
		COALESCE(SUM(discount), 0) AS discount,
		COALESCE(SUM(payment_balance), 0) AS payment_balance,
		COALESCE(SUM(expense), 0) AS expense
	FROM (%s) t
	GROUP BY deposit_type`, unionQuery)

	var rows []model.PaymentDepositReportRecapRow
	if err := repo.Raw(recapQuery, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *RepositoryPaymentDepositReportImpl) buildCountAndDataQueries(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (string, string, []interface{}) {
	sql, args := repo.buildQuery(dataFilter, custId)
	orderSQL := repo.buildSafeSort(dataFilter.Sort)
	countSQL := fmt.Sprintf("SELECT COUNT(1) FROM (%s) t", sql)
	dataSQL := fmt.Sprintf("SELECT * FROM (%s) t ORDER BY %s", sql, orderSQL)
	return countSQL, dataSQL, args
}

func (repo *RepositoryPaymentDepositReportImpl) InsertReportList(c context.Context, report model.ReportList) error {
	return repo.model(c).Create(&report).Error
}

func (repo *RepositoryPaymentDepositReportImpl) UpdateReportList(c context.Context, reportID string, status int, fileBase64 string) error {
	updates := map[string]interface{}{
		"file_status": status,
		"file_base64": fileBase64,
	}

	return repo.model(c).Model(&model.ReportList{}).Where("report_id = ?", reportID).Updates(updates).Error
}

func (repo *RepositoryPaymentDepositReportImpl) GetReportRunningNumber(custId string, date time.Time) (int, error) {
	var count int64
	prefix := fmt.Sprintf("%s-%s-%%", "DownloadDepositPayment", date.Format("020106"))

	err := repo.Model(&model.ReportList{}).
		Where("cust_id = ? AND report_name LIKE ?", custId, prefix).
		Count(&count).Error

	return int(count), err
}

func (repo *RepositoryPaymentDepositReportImpl) buildQuery(dataFilter entity.PaymentDepositReportQueryFilter, custId string) (string, []interface{}) {
	depositTypes := normalizeDepositTypeFilter(dataFilter.DepositType)
	depositNos := normalizeDepositNoFilter(dataFilter.DepositNo)
	empIDs := normalizeIntFilter(dataFilter.EmpID)
	startDate := dataFilter.StartDate
	endDate := dataFilter.EndDate
	search := strings.TrimSpace(dataFilter.Q)

	queries := make([]string, 0, len(depositTypes))
	args := make([]interface{}, 0)

	for _, depositType := range depositTypes {
		switch depositType {
		case "AR":
			query, queryArgs := repo.buildARPaymentDepositQuery(custId, startDate, endDate, empIDs, depositNos, search)
			queries = append(queries, query)
			args = append(args, queryArgs...)
		case "AP":
			query, queryArgs := repo.buildAPPaymentDepositQuery(custId, startDate, endDate, depositNos, search)
			queries = append(queries, query)
			args = append(args, queryArgs...)
		}
	}

	if len(queries) == 0 {
		return "SELECT NULL WHERE 1 = 0", nil
	}
	if len(queries) == 1 {
		return queries[0], args
	}

	return strings.Join(queries, " UNION ALL "), args
}

func (repo *RepositoryPaymentDepositReportImpl) buildARPaymentDepositQuery(custId, startDate, endDate string, empIDs []int, depositNos []string, search string) (string, []interface{}) {
	query := `SELECT
		d.deposit_date AS deposit_date,
		'AR' AS deposit_type,
		d.deposit_no AS deposit_no,
		d.emp_id AS collector_id,
		COALESCE(NULLIF(me.emp_code, ''), CAST(d.emp_id AS varchar)) AS collector_code,
		COALESCE(NULLIF(me.emp_name, ''), COALESCE(NULLIF(me.emp_code, ''), CAST(d.emp_id AS varchar))) AS collector_name,
		COALESCE(dp.cash_amount, 0) AS cash_amount,
		COALESCE(dp.cheque_bg_amount, 0) AS cheque_amount,
		COALESCE(dp.transfer_amount, 0) AS transfer_amount,
		COALESCE(dp.return_amount, 0) AS return_amount,
		COALESCE(dp.credit_debit_amount, 0) AS credit_debit_amount,
		COALESCE(de.expense_amount, 0) AS expense_amount,
		(COALESCE(dp.cash_amount, 0) + COALESCE(dp.cheque_bg_amount, 0) + COALESCE(dp.transfer_amount, 0) + COALESCE(dp.return_amount, 0) + COALESCE(dp.credit_debit_amount, 0) - COALESCE(de.expense_amount, 0)) AS total_payment
	FROM acf.deposit d
	JOIN (
		SELECT deposit_no, cust_id,
			SUM(CASE WHEN pay_type = 1 THEN payment_amount ELSE 0 END) AS cash_amount,
			SUM(CASE WHEN pay_type = 2 THEN payment_amount ELSE 0 END) AS cheque_bg_amount,
			SUM(CASE WHEN pay_type = 3 THEN payment_amount ELSE 0 END) AS transfer_amount,
			SUM(CASE WHEN pay_type = 4 THEN payment_amount ELSE 0 END) AS return_amount,
			SUM(CASE WHEN pay_type = 5 THEN payment_amount ELSE 0 END) AS credit_debit_amount
		FROM acf.deposit_payment
		WHERE cust_id = ?
		GROUP BY deposit_no, cust_id
	) dp ON dp.deposit_no = d.deposit_no AND dp.cust_id = d.cust_id
	LEFT JOIN (
		SELECT deposit_no, cust_id, SUM(payment_amount) AS expense_amount
		FROM acf.deposit_expense
		WHERE cust_id = ?
		GROUP BY deposit_no, cust_id
	) de ON de.deposit_no = d.deposit_no AND de.cust_id = d.cust_id
	LEFT JOIN mst.m_employee me ON me.emp_id = d.emp_id AND me.cust_id = d.cust_id
	WHERE d.cust_id = ? AND d.deleted_at IS NULL AND d.deposit_date BETWEEN ? AND ?`
	args := []interface{}{custId, custId, custId, startDate, endDate}

	if len(empIDs) > 0 {
		query += " AND d.emp_id IN ?"
		args = append(args, empIDs)
	}
	if len(depositNos) > 0 {
		query += " AND d.deposit_no IN ?"
		args = append(args, depositNos)
	}
	if search != "" {
		query += " AND (LOWER(d.deposit_no) LIKE ? OR LOWER(COALESCE(me.emp_name, '')) LIKE ?)"
		like := "%" + strings.ToLower(search) + "%"
		args = append(args, like, like)
	}

	return query, args
}

func (repo *RepositoryPaymentDepositReportImpl) buildAPPaymentDepositQuery(custId, startDate, endDate string, depositNos []string, search string) (string, []interface{}) {
	query := `SELECT
		app.account_payable_payment_date AS deposit_date,
		'AP' AS deposit_type,
		app.account_payable_payment_no AS deposit_no,
		NULL AS collector_id,
		NULL AS collector_code,
		NULL AS collector_name,
		COALESCE(appo.cash_amount, 0) AS cash_amount,
		COALESCE(appo.cheque_bg_amount, 0) AS cheque_amount,
		COALESCE(appo.transfer_amount, 0) AS transfer_amount,
		COALESCE(appo.return_amount, 0) AS return_amount,
		COALESCE(appo.credit_debit_amount, 0) AS credit_debit_amount,
		0 AS expense_amount,
		(COALESCE(appo.cash_amount, 0) + COALESCE(appo.cheque_bg_amount, 0) + COALESCE(appo.transfer_amount, 0) + COALESCE(appo.return_amount, 0) + COALESCE(appo.credit_debit_amount, 0)) AS total_payment
		FROM acf.account_payable_payment app
		LEFT JOIN (
			SELECT account_payable_payment_no, cust_id,
				SUM(CASE WHEN pay_type = 1 THEN payment_amount ELSE 0 END) AS cash_amount,
			SUM(CASE WHEN pay_type = 2 THEN payment_amount ELSE 0 END) AS cheque_bg_amount,
			SUM(CASE WHEN pay_type = 3 THEN payment_amount ELSE 0 END) AS transfer_amount,
			SUM(CASE WHEN pay_type = 4 THEN payment_amount ELSE 0 END) AS return_amount,
			SUM(CASE WHEN pay_type = 5 THEN payment_amount ELSE 0 END) AS credit_debit_amount
		FROM acf.account_payable_payment_options
		WHERE cust_id = ?
		GROUP BY account_payable_payment_no, cust_id
	) appo ON appo.account_payable_payment_no = app.account_payable_payment_no AND appo.cust_id = app.cust_id
	WHERE app.cust_id = ? AND app.deleted_by IS NULL AND app.account_payable_payment_date BETWEEN ? AND ?`
	args := []interface{}{custId, custId, startDate, endDate}

	if len(depositNos) > 0 {
		query += " AND app.account_payable_payment_no IN ?"
		args = append(args, depositNos)
	}
	if search != "" {
		query += " AND LOWER(app.account_payable_payment_no) LIKE ?"
		args = append(args, "%"+strings.ToLower(search)+"%")
	}

	return query, args
}

func (repo *RepositoryPaymentDepositReportImpl) buildDownloadUnionQuery(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) (string, []interface{}) {
	depositTypes := normalizeDepositTypeFilter(dataFilter.DepositType)
	depositNos := normalizeDepositNoFilter(dataFilter.DepositNo)
	empIDs := normalizeIntFilter(dataFilter.EmpID)
	queries := make([]string, 0, len(depositTypes))
	args := make([]interface{}, 0)

	for _, depositType := range depositTypes {
		switch depositType {
		case "AR":
			query, queryArgs := repo.buildDownloadARQuery(custId, parentCustId, dataFilter.StartDate, dataFilter.EndDate, empIDs, depositNos)
			queries = append(queries, query)
			args = append(args, queryArgs...)
		case "AP":
			query, queryArgs := repo.buildDownloadAPQuery(custId, parentCustId, dataFilter.StartDate, dataFilter.EndDate, depositNos)
			queries = append(queries, query)
			args = append(args, queryArgs...)
		}
	}

	if len(queries) == 0 {
		return "SELECT NULL::timestamp AS deposit_date WHERE 1 = 0", nil
	}
	if len(queries) == 1 {
		return queries[0], args
	}

	return strings.Join(queries, " UNION ALL "), args
}

func (repo *RepositoryPaymentDepositReportImpl) buildDownloadQuery(dataFilter entity.PaymentDepositReportQueryFilter, custId, parentCustId string) (string, []interface{}) {
	unionQuery, args := repo.buildDownloadUnionQuery(dataFilter, custId, parentCustId)
	return fmt.Sprintf("SELECT * FROM (%s) t ORDER BY t.deposit_date, t.deposit_no, t.document_date, t.document_no", unionQuery), args
}

func (repo *RepositoryPaymentDepositReportImpl) buildDownloadARQuery(custId, parentCustId, startDate, endDate string, empIDs []int, depositNos []string) (string, []interface{}) {
	query := `SELECT
		d.deposit_date AS deposit_date,
		'Account Receivable' AS deposit_type,
		d.deposit_no AS deposit_no,
		COALESCE(NULLIF(me.emp_name, ''), '') AS collector,
		dp.document_date AS document_date,
		COALESCE(NULLIF(dp.code, ''), '') AS code,
		COALESCE(NULLIF(dp.business_name, ''), '') AS business_name,
		COALESCE(NULLIF(dp.invoice_no, ''), '') AS document_no,
		COALESCE(dp.cash, 0) AS cash,
		COALESCE(dp.cheque_giro, 0) AS cheque_giro,
		COALESCE(dp.transfer, 0) AS transfer,
		COALESCE(dp.return_amount, 0) AS return_amount,
		COALESCE(dp.credit_debit, 0) AS credit_debit,
		COALESCE(dp.discount, 0) AS discount,
		COALESCE(dd.payment_balance, 0) AS payment_balance,
		0 AS expense,
		'' AS expense_name
	FROM acf.deposit d
	JOIN (
		SELECT
			dp2.deposit_no,
			dp2.cust_id,
			dp2.invoice_no,
			o.invoice_date AS document_date,
			mo.outlet_code AS code,
			mo.outlet_name AS business_name,
			SUM(CASE WHEN dp2.pay_type = 1 THEN COALESCE(dp2.payment_amount, 0) ELSE 0 END) AS cash,
			SUM(CASE WHEN dp2.pay_type = 2 THEN COALESCE(dp2.payment_amount, 0) ELSE 0 END) AS cheque_giro,
			SUM(CASE WHEN dp2.pay_type = 3 THEN COALESCE(dp2.payment_amount, 0) ELSE 0 END) AS transfer,
			SUM(CASE WHEN dp2.pay_type = 4 THEN COALESCE(dp2.payment_amount, 0) ELSE 0 END) AS return_amount,
			SUM(CASE WHEN dp2.pay_type = 5 THEN COALESCE(dp2.payment_amount, 0) ELSE 0 END) AS credit_debit,
			MAX(COALESCE(dp2.discount, 0)) AS discount
		FROM acf.deposit_payment dp2
		LEFT JOIN sls.order o ON o.invoice_no = dp2.invoice_no AND o.cust_id = dp2.cust_id
		LEFT JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND (mo.cust_id = dp2.cust_id OR mo.cust_id = ?)
		WHERE dp2.cust_id = ?
		GROUP BY dp2.deposit_no, dp2.cust_id, dp2.invoice_no, o.invoice_date, mo.outlet_code, mo.outlet_name
	) dp ON dp.deposit_no = d.deposit_no AND dp.cust_id = d.cust_id
	LEFT JOIN acf.deposit_detail dd ON dd.deposit_no = dp.deposit_no AND dd.invoice_no = dp.invoice_no AND dd.cust_id = d.cust_id
	LEFT JOIN mst.m_employee me ON me.emp_id = d.emp_id AND me.cust_id = d.cust_id
	WHERE d.cust_id = ? AND d.deleted_at IS NULL AND d.deposit_date BETWEEN ? AND ?`
	args := []interface{}{parentCustId, custId, custId, startDate, endDate}

	if len(empIDs) > 0 {
		query += " AND d.emp_id IN ?"
		args = append(args, empIDs)
	}
	if len(depositNos) > 0 {
		query += " AND d.deposit_no IN ?"
		args = append(args, depositNos)
	}
	query += `
	UNION ALL
	SELECT
		d.deposit_date AS deposit_date,
		'Account Receivable' AS deposit_type,
		d.deposit_no AS deposit_no,
		COALESCE(NULLIF(me.emp_name, ''), '') AS collector,
		ex.date AS document_date,
		'' AS code,
		'' AS business_name,
		COALESCE(NULLIF(ex.doc_no, ''), '') AS document_no,
		0 AS cash,
		0 AS cheque_giro,
		0 AS transfer,
		0 AS return_amount,
		0 AS credit_debit,
		0 AS discount,
		0 AS payment_balance,
		-ABS(COALESCE(SUM(de.payment_amount), 0)) AS expense,
		COALESCE(NULLIF(CONCAT_WS(' - ', NULLIF(etr.expense_type_code, ''), NULLIF(etr.expense_type_name, '')), ''), '') AS expense_name
	FROM acf.deposit d
	JOIN acf.deposit_expense de ON de.deposit_no = d.deposit_no AND de.cust_id = d.cust_id
	LEFT JOIN acf.expense ex ON ex.expense_id = de.expense_id AND ex.cust_id = d.cust_id AND ex.deleted_at IS NULL
	LEFT JOIN acf.expense_type etr ON etr.expense_type_id = ex.expense_type_id
	LEFT JOIN mst.m_employee me ON me.emp_id = d.emp_id AND me.cust_id = d.cust_id
	WHERE d.cust_id = ? AND d.deleted_at IS NULL AND d.deposit_date BETWEEN ? AND ?`
	args = append(args, custId, startDate, endDate)
	if len(empIDs) > 0 {
		query += " AND d.emp_id IN ?"
		args = append(args, empIDs)
	}
	if len(depositNos) > 0 {
		query += " AND d.deposit_no IN ?"
		args = append(args, depositNos)
	}
	query += " GROUP BY d.deposit_date, d.deposit_no, me.emp_name, ex.date, ex.doc_no, etr.expense_type_code, etr.expense_type_name HAVING COALESCE(SUM(de.payment_amount), 0) <> 0"

	return query, args
}

func (repo *RepositoryPaymentDepositReportImpl) buildDownloadAPQuery(custId, parentCustId, startDate, endDate string, depositNos []string) (string, []interface{}) {
	query := `SELECT
		app.account_payable_payment_date AS deposit_date,
		'Account Payable' AS deposit_type,
		app.account_payable_payment_no AS deposit_no,
		'' AS collector,
		appd.invoice_date AS document_date,
		COALESCE(NULLIF(sup.sup_code, ''), '') AS code,
		COALESCE(NULLIF(sup.sup_name, ''), '') AS business_name,
		COALESCE(NULLIF(appo.invoice_no, ''), COALESCE(NULLIF(appd.invoice_no, ''), '')) AS document_no,
		SUM(CASE WHEN appo.pay_type = 1 THEN COALESCE(appo.payment_amount, 0) ELSE 0 END) AS cash,
		SUM(CASE WHEN appo.pay_type = 2 THEN COALESCE(appo.payment_amount, 0) ELSE 0 END) AS cheque_giro,
		SUM(CASE WHEN appo.pay_type = 3 THEN COALESCE(appo.payment_amount, 0) ELSE 0 END) AS transfer,
		SUM(CASE WHEN appo.pay_type = 4 THEN COALESCE(appo.payment_amount, 0) ELSE 0 END) AS return_amount,
		SUM(CASE WHEN appo.pay_type = 5 THEN COALESCE(appo.payment_amount, 0) ELSE 0 END) AS credit_debit,
		COALESCE(appd.discount, 0) AS discount,
		COALESCE(appd.payment_balance, 0) AS payment_balance,
		0 AS expense,
		'' AS expense_name
	FROM acf.account_payable_payment app
	JOIN acf.account_payable_payment_options appo ON appo.account_payable_payment_no = app.account_payable_payment_no AND appo.cust_id = app.cust_id
	LEFT JOIN acf.account_payable_payment_detail appd ON appd.account_payable_payment_no = appo.account_payable_payment_no AND appd.invoice_no = appo.invoice_no AND appd.cust_id = appo.cust_id
	LEFT JOIN acf.account_payable ap ON ap.invoice_no = appo.invoice_no AND ap.cust_id = app.cust_id AND ap.deleted_at IS NULL
	LEFT JOIN mst.m_supplier sup ON sup.sup_id = ap.sup_id AND sup.cust_id = ?
	WHERE app.cust_id = ? AND app.deleted_by IS NULL AND app.account_payable_payment_date BETWEEN ? AND ?`
	args := []interface{}{parentCustId, custId, startDate, endDate}
	if len(depositNos) > 0 {
		query += " AND app.account_payable_payment_no IN ?"
		args = append(args, depositNos)
	}
	query += `
	GROUP BY
		app.account_payable_payment_date,
		app.account_payable_payment_no,
		appo.invoice_no,
		appd.invoice_no,
		appd.invoice_date,
		sup.sup_code,
		sup.sup_name,
		appd.discount,
		appd.payment_balance`
	return query, args
}

func (repo *RepositoryPaymentDepositReportImpl) buildSafeSort(sortValue string) string {
	defaultSort := "t.deposit_date DESC"
	if strings.TrimSpace(sortValue) == "" {
		return defaultSort
	}

	allowedFields := map[string]string{
		"created_date":   "t.deposit_date",
		"deposit_date":   "t.deposit_date",
		"deposit_no":     "t.deposit_no",
		"deposit_type":   "t.deposit_type",
		"collector_name": "t.collector_name",
		"total_payment":  "t.total_payment",
	}

	orders := make([]string, 0)
	for _, raw := range strings.Split(sortValue, ",") {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}

		sortPart := strings.SplitN(part, ":", 2)
		if len(sortPart) != 2 {
			continue
		}

		field := strings.TrimSpace(strings.ToLower(sortPart[0]))
		direction := strings.TrimSpace(strings.ToLower(sortPart[1]))

		column, ok := allowedFields[field]
		if !ok {
			continue
		}
		if direction != "asc" && direction != "desc" {
			continue
		}

		orders = append(orders, column+" "+strings.ToUpper(direction))
	}

	if len(orders) == 0 {
		return defaultSort
	}

	return strings.Join(orders, ", ")
}

func normalizeDepositTypeFilter(values []string) []string {
	set := map[string]struct{}{}
	for _, item := range values {
		for _, token := range strings.Split(item, ",") {
			value := strings.ToUpper(strings.TrimSpace(token))
			if value == "" {
				continue
			}
			set[value] = struct{}{}
		}
	}
	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func normalizeDepositNoFilter(depositNo []string) []string {
	return normalizeStringList(depositNo)
}

func normalizeIntFilter(values []string) []int {
	items := normalizeStringList(values)
	result := make([]int, 0, len(items))
	for _, item := range items {
		if v, err := strconv.Atoi(item); err == nil {
			result = append(result, v)
		}
	}
	return result
}

func normalizeStringList(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, item := range values {
		for _, token := range strings.Split(item, ",") {
			value := strings.TrimSpace(token)
			if value == "" {
				continue
			}
			normalized = append(normalized, value)
		}
	}
	return normalized
}
