package repository

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryChequeGiroClearingImpl struct {
		*gorm.DB
	}
)

type ChequeGiroClearingRepository interface {
	FindByNo(ChequeGiroNo int, custId string) (whAdj model.ChequeGiroClearingList, err error)
	FindAllByCustId(dataFilter entity.CheckGiroClearingQueryFilter) ([]model.ChequeGiroClearingList, int64, int, error)
	UpdateClearing(c context.Context, ChequeGiroNo int, custId string, data model.ChequeGiroClearingList) error
	FindDetailPaymentDepositByGiroNo(GiroNo string, custId string) (whAdj []model.DepositPayment, err error)
	StorePayment(c context.Context, data *model.DepositPayment) (int, error)
	DeleteAllDetailPaymentByDepositPaymentID(c context.Context, depositPaymentID []int) error
	FindDetailPaymentCashByDepositInvoiceNo(payType int, depositNo string, invoiceNo string, custId string) (whAdj []model.DepositPayment, err error)
	UpdateCashAmount(c context.Context, depositPaymentId int, data model.DepositPayment) error
}

func NewChequeGiroClearingRepo(db *gorm.DB) *RepositoryChequeGiroClearingImpl {
	return &RepositoryChequeGiroClearingImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryChequeGiroClearingImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryChequeGiroClearingImpl) FindByNo(ChequeGiroNo int, custId string) (whAdj model.ChequeGiroClearingList, err error) {
	err = repository.Select(
		`acf.cheque_giro.*,us.user_fullname AS updated_by_name,sls.emp_id as salesman_id, sls.sales_name as sales_name,sp.sup_name,b.bank_name,bc.bank_name as bank_id_collecting_name,o.outlet_name,o.outlet_code,ab.payment as used_amount`).
		Joins("left join sys.m_user us on us.user_id = acf.cheque_giro.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.cheque_giro.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_supplier sp on sp.sup_id = acf.cheque_giro.sup_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_bank b on b.bank_id = acf.cheque_giro.bank_id AND b.cust_id = ?", custId).
		Joins("left join mst.m_bank bc on bc.bank_id = acf.cheque_giro.bank_id_collecting AND bc.cust_id = ?", custId).
		Joins("left join mst.m_outlet o on o.outlet_id = acf.cheque_giro.outlet_id AND o.cust_id = ?", custId).
		Joins(`left join (
			select dp.document_no,
				coalesce(SUM(dp.payment_amount), 0) as payment
			from acf.deposit_payment dp
			inner join acf.deposit d on d.deposit_no = dp.deposit_no and d.cust_id = dp.cust_id
			where dp.pay_type = 2
				and dp.cust_id = ?
				and d.deposit_status in ?
			group by dp.document_no
		) ab on ab.document_no = acf.cheque_giro.doc_no_cheque`, custId, []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}).
		Joins("left join mst.m_outlet_bank ob on ob.outlet_bank_id = acf.cheque_giro.outlet_bank_id AND o.cust_id = ?", custId).
		Where("acf.cheque_giro.cheque_giro_no = ? AND acf.cheque_giro.cust_id=?", ChequeGiroNo, custId).
		Where("acf.cheque_giro.is_del=false").
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryChequeGiroClearingImpl) FindAllByCustId(dataFilter entity.CheckGiroClearingQueryFilter) ([]model.ChequeGiroClearingList, int64, int, error) {
	var ChequeGiro []model.ChequeGiroClearingList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("cheque_giro_no")
	query := repository.Select(
		`acf.cheque_giro.*,us.user_fullname AS updated_by_name,sls.emp_id as salesman_id, sls.sales_name as sales_name,sp.sup_name,b.bank_name,o.outlet_name,o.outlet_code`).
		Joins("left join sys.m_user us on us.user_id = acf.cheque_giro.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.cheque_giro.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier sp on sp.sup_id = acf.cheque_giro.sup_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_bank b on b.bank_id = acf.cheque_giro.bank_id AND b.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet o on o.outlet_id = acf.cheque_giro.outlet_id AND o.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.cheque_giro.cust_id=?", dataFilter.CustId)
	query.Where("acf.cheque_giro.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.cheque_giro.is_del=false")
	query.Where("acf.cheque_giro.is_del=false")

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.cheque_giro.doc_date_cheque between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.cheque_giro.doc_date_cheque between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.cheque_giro.doc_no_cheque ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.cheque_giro.doc_no_cheque ILIKE ? )", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.BankID) > 0 {
		query.Where("acf.cheque_giro.bank_id in ?", dataFilter.BankID)
		queryCount.Where("acf.cheque_giro.bank_id in ?", dataFilter.BankID)
	}

	if len(dataFilter.StatusID) > 0 {
		query.Where("acf.cheque_giro.status_cheque in ?", dataFilter.StatusID)
		queryCount.Where("acf.cheque_giro.status_cheque in ?", dataFilter.StatusID)
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
		query.Order("cheque_giro_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ChequeGiro).Error
	if err != nil {
		return ChequeGiro, total, 0, err
	}
	err = queryCount.Model(&ChequeGiro).Count(&total).Error
	if err != nil {
		return ChequeGiro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ChequeGiro, total, lastPage, nil
}

func (repository *RepositoryChequeGiroClearingImpl) UpdateClearing(c context.Context, ChequeGiroNo int, custId string, data model.ChequeGiroClearingList) error {

	result := repository.model(c).Model(&data).Where("cheque_giro_no=? AND cust_id = ?", ChequeGiroNo, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}

func (repository *RepositoryChequeGiroClearingImpl) UpdateCashAmount(c context.Context, depositPaymentId int, data model.DepositPayment) error {

	result := repository.model(c).Model(&data).Where("deposit_payment_id=?", depositPaymentId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}

func (repository *RepositoryChequeGiroClearingImpl) FindDetailPaymentDepositByGiroNo(GiroNo string, custId string) (whAdj []model.DepositPayment, err error) {
	err = repository.Select(`
			acf.deposit_payment.*,
			od.invoice_date
		`).
		Joins("LEFT JOIN acf.deposit_detail dd ON dd.deposit_no = acf.deposit_payment.deposit_no AND dd.invoice_no = acf.deposit_payment.invoice_no").
		Joins("LEFT JOIN sls.order od ON od.invoice_no = dd.invoice_no AND od.cust_id = ?", custId).
		Where("acf.deposit_payment.document_no = ? AND acf.deposit_payment.cust_id=?", GiroNo, custId).
		Where("acf.deposit_payment.pay_type = ? AND acf.deposit_payment.cust_id=?", 2, custId).
		// Where("acf.deposit_payment.is_del=false").
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryChequeGiroClearingImpl) FindDetailPaymentCashByDepositInvoiceNo(payType int, depositNo string, invoiceNo string, custId string) (whAdj []model.DepositPayment, err error) {
	err = repository.Select(`
			acf.deposit_payment.*,
			od.invoice_date
		`).
		Joins("LEFT JOIN acf.deposit_detail dd ON dd.deposit_no = acf.deposit_payment.deposit_no AND dd.invoice_no = acf.deposit_payment.invoice_no").
		Joins("LEFT JOIN sls.order od ON od.invoice_no = dd.invoice_no AND od.cust_id = ?", custId).
		Where("acf.deposit_payment.deposit_no = ? AND acf.deposit_payment.cust_id=?", depositNo, custId).
		Where("acf.deposit_payment.invoice_no = ? AND acf.deposit_payment.cust_id=?", invoiceNo, custId).
		Where("acf.deposit_payment.pay_type = ? AND acf.deposit_payment.cust_id=?", payType, custId).
		// Where("acf.deposit_payment.is_del=false").
		Find(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryChequeGiroClearingImpl) StorePayment(c context.Context, data *model.DepositPayment) (int, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.DepositPaymentID, nil
}

func (repository *RepositoryChequeGiroClearingImpl) DeleteAllDetailPaymentByDepositPaymentID(c context.Context, depositPaymentID []int) error {
	var Details model.DepositPayment
	err := repository.model(c).Where("deposit_payment_id IN (?)", depositPaymentID).Delete(&Details).Error

	return err
}
