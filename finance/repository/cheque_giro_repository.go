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
	RepositoryChequeGiroImpl struct {
		*gorm.DB
	}
)

type ChequeGiroRepository interface {
	Store(c context.Context, data *model.ChequeGiro) error
	FindByNo(ChequeGiroNo int, custId string) (whAdj model.ChequeGiroList, err error)
	FindAllByCustId(dataFilter entity.CheckGiroQueryFilter) ([]model.ChequeGiroList, int64, int, error)
	Update(c context.Context, ChequeGiroNo int, data model.ChequeGiro) error
	Delete(c context.Context, custId string, ChequeGiroNo int, deletedBy int64) error
	FindAllBankByCustId(dataFilter entity.CheckGiroQueryFilter) ([]model.BankLookup, int64, int, error)
	FindAllBankAccountByCustId(dataFilter entity.CheckGiroQueryFilter, bankID []int) ([]model.BankAccountLookup, int64, int, error)
}

func NewChequeGiroRepo(db *gorm.DB) *RepositoryChequeGiroImpl {
	return &RepositoryChequeGiroImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryChequeGiroImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryChequeGiroImpl) Store(c context.Context, data *model.ChequeGiro) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryChequeGiroImpl) FindByNo(ChequeGiroNo int, custId string) (whAdj model.ChequeGiroList, err error) {
	err = repository.Select(
		`acf.cheque_giro.*,us.user_fullname AS updated_by_name,sls.emp_id as salesman_id, sls.sales_name as sales_name,sp.sup_name,b.bank_name,o.outlet_name, SUM(appo.payment_amount) as used_amount, ab.payment as used_amount_outlet`).
		Joins("left join sys.m_user us on us.user_id = acf.cheque_giro.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.cheque_giro.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_supplier sp on sp.sup_id = acf.cheque_giro.sup_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_bank b on b.bank_id = acf.cheque_giro.bank_id AND b.cust_id = ?", custId).
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
		Joins("left join acf.account_payable_payment_options appo on appo.document_no = acf.cheque_giro.doc_no_cheque and appo.cust_id = ?", custId).
		Joins("left join mst.m_outlet_bank ob on ob.outlet_bank_id = acf.cheque_giro.outlet_bank_id AND o.cust_id = ?", custId).
		Where("acf.cheque_giro.cheque_giro_no = ? AND acf.cheque_giro.cust_id=?", ChequeGiroNo, custId).
		Where("acf.cheque_giro.is_del=false").
		Group("acf.cheque_giro.owner_id, acf.cheque_giro.cust_id, acf.cheque_giro.cheque_giro_no, acf.cheque_giro.doc_no_cheque, acf.cheque_giro.salesman_id, acf.cheque_giro.outlet_id, acf.cheque_giro.bank_id, acf.cheque_giro.bank_id_collecting, acf.cheque_giro.account_no, acf.cheque_giro.doc_date_cheque, acf.cheque_giro.due_date, acf.cheque_giro.amount, acf.cheque_giro.status_cheque, acf.cheque_giro.clearing_date, acf.cheque_giro.created_by, acf.cheque_giro.created_at, acf.cheque_giro.updated_by, acf.cheque_giro.updated_at, acf.cheque_giro.is_del, acf.cheque_giro.deleted_by, acf.cheque_giro.deleted_at, acf.cheque_giro.outlet_bank_id, acf.cheque_giro.sup_id, acf.cheque_giro.cash_no, acf.cheque_giro.cash_amount, acf.cheque_giro.transfer_no, acf.cheque_giro.transfer_amount, acf.cheque_giro.cheque_no, acf.cheque_giro.cheque_amount, acf.cheque_giro.reason, us.user_fullname, sls.emp_id, sls.sales_name, sp.sup_name, b.bank_name, o.outlet_name, ab.payment, cheque_giro.remaining_amount, cheque_giro.paid_amount").
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryChequeGiroImpl) FindAllByCustId(dataFilter entity.CheckGiroQueryFilter) ([]model.ChequeGiroList, int64, int, error) {
	var ChequeGiro []model.ChequeGiroList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("cheque_giro_no")
	query := repository.Select(
		`acf.cheque_giro.*,us.user_fullname AS updated_by_name,sls.emp_id as salesman_id, sls.sales_name as sales_name,sp.sup_name,b.bank_name,o.outlet_name`).
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

	if len(dataFilter.AccountNo) > 0 {
		query.Where("acf.cheque_giro.account_no in ?", dataFilter.AccountNo)
		queryCount.Where("acf.cheque_giro.account_no in ?", dataFilter.AccountNo)
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

func (repository *RepositoryChequeGiroImpl) Delete(c context.Context, custId string, ChequeGiroNo int, deletedBy int64) error {
	var data model.ChequeGiro
	result := repository.model(c).Model(&data).Where("cheque_giro_no=? AND cust_id = ? AND is_del= ? ", ChequeGiroNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryChequeGiroImpl) Update(c context.Context, ChequeGiroNo int, data model.ChequeGiro) error {

	result := repository.model(c).Model(&data).Where("cheque_giro_no=?", ChequeGiroNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}

func (repository *RepositoryChequeGiroImpl) FindAllBankByCustId(dataFilter entity.CheckGiroQueryFilter) ([]model.BankLookup, int64, int, error) {
	var ChequeGiro []model.BankLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.cheque_giro.bank_id))")
	query := repository.Select(
		`distinct(acf.cheque_giro.bank_id),mb.bank_name,mb.bank_code`).
		Joins("left join mst.m_bank mb on mb.bank_id = acf.cheque_giro.bank_id ")

	queryCount.Where("acf.cheque_giro.cust_id=?", dataFilter.CustId)
	query.Where("acf.cheque_giro.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.cheque_giro.is_del=false")
	query.Where("acf.cheque_giro.is_del=false")

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.cheque_giro.bank_name ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.cheque_giro.bank_name ILIKE ? )", "%"+dataFilter.Query+"%")
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
		query.Order("bank_name ASC")
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

func (repository *RepositoryChequeGiroImpl) FindAllBankAccountByCustId(dataFilter entity.CheckGiroQueryFilter, bankID []int) ([]model.BankAccountLookup, int64, int, error) {
	var ChequeGiro []model.BankAccountLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.cheque_giro.account_no))")
	query := repository.Select(
		`distinct(acf.cheque_giro.account_no)`)

	if len(dataFilter.BankID) > 0 {
		query.Where("acf.cheque_giro.bank_id in ?", dataFilter.BankID)
		queryCount.Where("acf.cheque_giro.bank_id in ?", dataFilter.BankID)
	}

	queryCount.Where("acf.cheque_giro.cust_id=?", dataFilter.CustId)
	query.Where("acf.cheque_giro.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.cheque_giro.is_del=false")
	query.Where("acf.cheque_giro.is_del=false")

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.cheque_giro.account_no ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.cheque_giro.account_no ILIKE ? )", "%"+dataFilter.Query+"%")
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
		query.Order("account_no ASC")
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
