package repository

import (
	"context"
	"mobile/model"
	"mobile/pkg/constant"

	"gorm.io/gorm"
)

type (
	RepositoryInvoicesImpl struct {
		*gorm.DB
	}
)
type InvoicesRepository interface {
	StorePaymentCndn(c context.Context, data *model.Cndn) error
	CountAllByCustId(custId string, depositDate string) (int, error)
	CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error)

	StoreDetail(c context.Context, data *model.DepositDetail) (int, error)
	StorePayment(c context.Context, data *model.DepositPayment) (int, error)
	StoreDepositPaymentImage(c context.Context, data *model.DepositPaymentImage) (int, error)
	Store(c context.Context, data *model.Deposit) error

	FindDetailPaymentByInvoice(payType int, depositNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error)
	FindPaymentImagesByNo(depositNo string, invoiceNo string) (whAdj []model.DepositPaymentImage, err error)
	GetRemainingOutstandingByOutletID(c context.Context, outletID int64) (float64, error)
	GetInvoiceListByDate(c context.Context, empID int64, date string) ([]model.InvoiceListItem, error)
	GetExpenseSummaryByDate(c context.Context, empID int64, date string) ([]model.ExpenseListItem, error)
	GetCollectionSummaryByDate(ctx context.Context, empID int64, date string) ([]model.CollectionListItem, []model.InvoicePayment, error)
	GetPaymentsByInvoiceNo(c context.Context, empID int64, invoiceNos []string) ([]model.InvoicePayment, error)
}

func NewInvoicesRepository(db *gorm.DB) *RepositoryInvoicesImpl {
	return &RepositoryInvoicesImpl{db}
}

func (repo *RepositoryInvoicesImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryInvoicesImpl) StorePaymentCndn(c context.Context, data *model.Cndn) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryInvoicesImpl) CountAllByCustId(custId string, depositDate string) (int, error) {
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

func (repository *RepositoryInvoicesImpl) CountRemainingAmountByInvoice(c context.Context, invoiceNo string, custId string) (float64, error) {
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

func (repository *RepositoryInvoicesImpl) StoreDetail(c context.Context, data *model.DepositDetail) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositDetailID, nil
}

func (repository *RepositoryInvoicesImpl) StorePayment(c context.Context, data *model.DepositPayment) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositPaymentID, nil
}

func (repository *RepositoryInvoicesImpl) StoreDepositPaymentImage(c context.Context, data *model.DepositPaymentImage) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositImageID, nil
}

func (repository *RepositoryInvoicesImpl) Store(c context.Context, data *model.Deposit) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryInvoicesImpl) FindDetailPaymentByInvoice(payType int, invoiceNo string, custId string) (whAdj []model.DepositPaymentInvoice, err error) {
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

func (repository *RepositoryInvoicesImpl) FindPaymentImagesByNo(depositNo string, invoiceNo string) ([]model.DepositPaymentImage, error) {
	var images []model.DepositPaymentImage
	err := repository.model(context.Background()).
		Where("deposit_no = ? AND invoice_no = ? ", depositNo, invoiceNo).
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (repository *RepositoryInvoicesImpl) GetRemainingOutstandingByOutletID(c context.Context, outletID int64) (float64, error) {
	var total float64

	query := `
	COALESCE(
		SUM(
			CASE
				WHEN o.opr_type = 'C' and o.invoice_date::date = CURRENT_DATE 
				THEN COALESCE(pt.remaining_amount, o.total)
				ELSE (o.total - COALESCE(paid_invoices.paid_amount, 0))
			END
	), 0) as total_outstanding`

	err := repository.model(c).
		Table("sls.order o").
		Select(query).
		Joins("LEFT JOIN acf.payment_trx pt ON pt.po_number = o.order_no AND pt.trx_source = 'C' AND pt.outlet_id = o.outlet_id AND pt.cust_id = o.cust_id").
		Joins("LEFT JOIN (SELECT dd.invoice_no, dd.cust_id, COALESCE(SUM(dd.total_payment), 0) AS paid_amount FROM acf.deposit_detail dd INNER JOIN acf.deposit d ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id AND d.deposit_status IN (1, 2) GROUP BY dd.invoice_no, dd.cust_id) paid_invoices ON paid_invoices.invoice_no = o.invoice_no AND paid_invoices.cust_id = o.cust_id").
		Where("o.outlet_id = ? AND o.invoice_no IS NOT NULL", outletID).
		Where("o.data_status IN (?, ?)", constant.OrderStatusInvoicing, constant.OrderStatusCompleted).
		Where("CASE WHEN o.opr_type = 'C' and o.invoice_date::date = CURRENT_DATE THEN COALESCE(pt.remaining_amount, o.total) > 0 ELSE (o.total - COALESCE(paid_invoices.paid_amount, 0)) > 0 END").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (repository *RepositoryInvoicesImpl) GetInvoiceListByDate(c context.Context, empID int64, date string) ([]model.InvoiceListItem, error) {
	var result []model.InvoiceListItem

	err := repository.model(c).
		Table("sls.order o").
		Select(`
			o.invoice_no,
			o.invoice_date, 
			o.ro_no,
			o.order_no,
			o.due_date,
			o.outlet_id,
			mo.outlet_code,
			mo.outlet_name,
			o.salesman_id,
			me.emp_code AS salesman_code,
			ms.sales_name AS salesman_name,
			o.total_final AS invoice_amount,
			o.total_final AS remaining_amount,
			o.disc_value_final AS discount,
			o.notes
		`).
		Joins(`
			INNER JOIN mst.m_outlet mo 
			ON o.outlet_id = mo.outlet_id 
			AND o.cust_id = mo.cust_id
		`).
		Joins(`
			INNER JOIN mst.m_salesman ms 
			ON o.salesman_id = ms.emp_id 
			AND o.cust_id = ms.cust_id
		`).
		Joins(`
			INNER JOIN mst.m_employee me 
			ON o.salesman_id = me.emp_id 
			AND o.cust_id = me.cust_id
		`).
		Joins(`
		  INNER JOIN acf.payment_trx p 
				ON p.po_number = o.order_no 
				AND p.emp_id = ? AND p.trx_source != 'L'
				`, empID).
		Where("o.invoice_date::date = ?", date).
		Where("o.salesman_id = ?", empID).
		Where("o.data_status IN (?)", constant.OrderStatusInvoicing).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
func (repository *RepositoryInvoicesImpl) GetExpenseSummaryByDate(c context.Context, empId int64, date string) ([]model.ExpenseListItem, error) {
	var result []model.ExpenseListItem

	err := repository.model(c).
		Table("acf.expense").
		Select(`
			acf.expense.expense_id as expense_id,
			acf.expense.doc_no as doc_no,
			acf.expense_type.expense_type_name as expense_name,
			acf.expense.amount as amount
		`).
		Joins("INNER JOIN acf.expense_type ON acf.expense.expense_type_id = acf.expense_type.expense_type_id").
		Where("expense.collector_id = ?", empId).
		Where("expense.created_at::date = ?", date).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository *RepositoryInvoicesImpl) GetCollectionSummaryByDate(ctx context.Context, empID int64, date string) ([]model.CollectionListItem, []model.InvoicePayment, error) {
	// First, define a struct to hold the raw flat data from the join query
	type rawCollectionData struct {
		CollectionNo    string  `gorm:"column:collection_no"`
		InvoiceNumber   string  `gorm:"column:invoice_number"`
		InvoiceDate     string  `gorm:"column:invoice_date"`
		RONo            string  `gorm:"column:ro_no"`
		DueDate         string  `gorm:"column:due_date"`
		OutletID        int     `gorm:"column:outlet_id"`
		OutletCode      string  `gorm:"column:outlet_code"`
		OutletName      string  `gorm:"column:outlet_name"`
		SalesmanID      int     `gorm:"column:salesman_id"`
		SalesmanCode    string  `gorm:"column:salesman_code"`
		SalesmanName    string  `gorm:"column:salesman_name"`
		InvoiceAmount   float64 `gorm:"column:invoice_amount"`
		RemainingAmount float64 `gorm:"column:remaining_amount"`
		Discount        float64 `gorm:"column:discount"`
		Notes           string  `gorm:"column:notes"`
	}

	var rawData []rawCollectionData

	err := repository.model(ctx).
		Table("acf.collection c").
		Select(`
		c.collection_no,
		cd.invoice_amount,
		cd.invoice_no as invoice_number,
		o.invoice_date,
		o.ro_no,
		o.due_date,
		o.outlet_id,
		mo.outlet_code,
		mo.outlet_name,
		o.salesman_id,
		me.emp_code AS salesman_code,
		ms.sales_name AS salesman_name,
		cd.invoice_amount,
		o.total_final as remaining_amount,
		o.disc_value as discount,
		o.notes
	`).
		Joins(`
		INNER JOIN acf.collection_det cd 
		ON cd.collection_no = c.collection_no 
		AND cd.cust_id = c.cust_id
	`).
		Joins(`
		INNER JOIN sls."order" o 
		ON o.invoice_no = cd.invoice_no 
		AND o.cust_id = c.cust_id
	`).
		Joins(`
		INNER JOIN mst.m_outlet mo 
		ON o.outlet_id = mo.outlet_id 
		AND o.cust_id = mo.cust_id
	`).
		Joins(`
		INNER JOIN mst.m_employee me 
		ON c.emp_id = me.emp_id 
		AND c.cust_id = me.cust_id
	`).
		Joins(`
		INNER JOIN mst.m_salesman ms 
		ON o.salesman_id = ms.emp_id 
		AND o.cust_id = ms.cust_id
	`).
		Where("c.emp_id = ?", empID).
		Where("c.created_at::date = ?", date).
		Where("o.data_status IN (?)", constant.OrderStatusInvoicing).
		Where("").
		Scan(&rawData).Error

	if err != nil {
		return nil, nil, err
	}

	var (
		orderedKeys    []string
		invoiceNos     []string
		collectionsMap = make(map[string]*model.CollectionListItem)
	)

	for _, row := range rawData {
		if _, exists := collectionsMap[row.CollectionNo]; !exists {
			collectionsMap[row.CollectionNo] = &model.CollectionListItem{
				CollectionNo: row.CollectionNo,
				TotalAmount:  0,
				Details:      []model.CollectionSummaryDetail{},
			}
			orderedKeys = append(orderedKeys, row.CollectionNo)
		}

		detail := model.CollectionSummaryDetail{
			InvoiceNumber:   row.InvoiceNumber,
			InvoiceDate:     row.InvoiceDate,
			RONo:            row.RONo,
			DueDate:         row.DueDate,
			OutletID:        row.OutletID,
			OutletCode:      row.OutletCode,
			OutletName:      row.OutletName,
			SalesmanID:      row.SalesmanID,
			SalesmanCode:    row.SalesmanCode,
			SalesmanName:    row.SalesmanName,
			InvoiceAmount:   row.InvoiceAmount,
			RemainintAmount: row.RemainingAmount,
			Discount:        row.Discount,
			Notes:           row.Notes,
			IsCollection:    true,
		}

		collection := collectionsMap[row.CollectionNo]
		collection.Details = append(collection.Details, detail)
		invoiceNos = append(invoiceNos, row.InvoiceNumber)
	}

	// Convert the map back to a list (maintaining order)
	var result []model.CollectionListItem
	for _, key := range orderedKeys {
		result = append(result, *collectionsMap[key])
	}

	if len(invoiceNos) == 0 {
		return result, nil, nil
	}

	invoices, err := repository.GetPaymentsByInvoiceNo(ctx, empID, invoiceNos)
	if err != nil {
		repository.Logger.Error(ctx, "Failed to get payments by invoice no", err)
	}

	return result, invoices, nil
}

func (repository *RepositoryInvoicesImpl) GetPaymentsByInvoiceNo(c context.Context, empID int64, invoiceNos []string) ([]model.InvoicePayment, error) {
	var payments []model.InvoicePayment

	err := repository.model(c).
		Table("acf.payment_trx pt").
		Select(`
		pt.total_transaction as invoice_amount,
		pt.remaining_amount,
		pt.payment_amount as total_payment,
		pt.remaining_amount as remaining_payment,
		pt.po_number as invoice_no,
		ptd.pay_type,
		CASE WHEN ptd.pay_type = 3 THEN bt.doc_no_bank ELSE pt.document_no END AS document_no,
		ptd.amount as payment_amount
	`).
		Joins("INNER JOIN acf.payment_trx_detail ptd ON ptd.payment_trx_id = pt.payment_trx_id").
		Joins("LEFT JOIN acf.bank_transfer bt on ptd.bank_transfer_no = bt.bank_transfer_no").
		Where("pt.emp_id = ?", empID).
		Where("(pt.po_number IN ? AND pt.trx_source != 'L') OR (pt.date = CURRENT_DATE AND pt.trx_source != 'L')", invoiceNos). // []string
		Scan(&payments).Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}
