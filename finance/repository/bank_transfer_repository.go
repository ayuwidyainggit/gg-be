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
	RepositoryBankTransferImpl struct {
		*gorm.DB
	}
)

type BankTransferRepository interface {
	Store(c context.Context, data *model.BankTransfer) error
	FindByNo(BankTransferNo int, custId string, parentCustId string) (whAdj model.BankTransferList, err error)
	FindAllByCustId(dataFilter entity.BankTransferQueryFilter) ([]model.BankTransferList, int64, int, error)
	Update(c context.Context, BankTransferNo int, custId string, data model.BankTransfer) error
	Delete(c context.Context, custId string, BankTransferNo int, deletedBy int64) error
	FindAllBankByCustId(dataFilter entity.BankTransferQueryFilter) ([]model.BankLookupBankTransfer, int64, int, error)
	FindAllBankAccountByCustId(dataFilter entity.BankTransferQueryFilter, bankID []int) ([]model.BankAccountLookupBankTransfer, int64, int, error)
	GetLastDocNoBank(c context.Context, prefix string) (string, error)
	StoreFile(c context.Context, data *model.BankTransferFile) error
	DeleteFilesByDocNo(c context.Context, docNo string, custId string) error
	FindOne(c context.Context, bankTransferNo int, custId string) (model.BankTransfer, error)
	FindFilesByDocNo(c context.Context, docNo string, custId string) ([]model.BankTransferFile, error)
	FindDepositDataByDocumentNo(c context.Context, docNo string, custId string) ([]model.BankTransferDepositDataRow, error)
}

func NewBankTransferRepo(db *gorm.DB) *RepositoryBankTransferImpl {
	return &RepositoryBankTransferImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryBankTransferImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryBankTransferImpl) Store(c context.Context, data *model.BankTransfer) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryBankTransferImpl) FindByNo(BankTransferNo int, custId string, parentCustId string) (whAdj model.BankTransferList, err error) {
	err = repository.Select(
		`acf.bank_transfer.*,us.user_fullname AS updated_by_name,sls.emp_id as salesman_id, sls.sales_name as sales_name,ep.emp_code as salesman_code,sp.sup_name,sp.sup_code,b.bank_name,o.outlet_name,o.outlet_code, appo.payment_amount as used_amount, ab.payment as used_amount_outlet`).
		Joins("left join sys.m_user us on us.user_id = acf.bank_transfer.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.bank_transfer.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_employee ep on ep.emp_id = sls.emp_id AND ep.cust_id = ?", custId).
		Joins("left join mst.m_supplier sp on sp.sup_id = acf.bank_transfer.sup_id AND sp.cust_id = ?", parentCustId).
		Joins("left join mst.m_bank b on b.bank_id = acf.bank_transfer.bank_id AND b.cust_id = ?", custId).
		Joins("left join mst.m_outlet o on o.outlet_id = acf.bank_transfer.outlet_id AND o.cust_id = ?", custId).
		Joins("left join acf.account_payable_payment_options appo on appo.document_no = acf.bank_transfer.doc_no_bank and appo.cust_id = ?", custId).
		Joins(`left join (
			select dp.document_no,
				coalesce(SUM(dp.payment_amount), 0) as payment
			from acf.deposit_payment dp
			inner join acf.deposit d on d.deposit_no = dp.deposit_no and d.cust_id = dp.cust_id
			where dp.pay_type = 3
				and dp.cust_id = ?
				and d.deposit_status in ?
			group by dp.document_no
		) ab on ab.document_no = acf.bank_transfer.doc_no_bank`, custId, []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}).
		Joins("left join mst.m_outlet_bank ob on ob.outlet_bank_id = acf.bank_transfer.outlet_bank_id AND o.cust_id = ?", custId).
		Where("acf.bank_transfer.bank_transfer_no = ? AND acf.bank_transfer.cust_id=?", BankTransferNo, custId).
		Where("acf.bank_transfer.is_del=false").
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryBankTransferImpl) FindAllByCustId(dataFilter entity.BankTransferQueryFilter) ([]model.BankTransferList, int64, int, error) {
	var BankTransfer []model.BankTransferList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("bank_transfer_no")
	query := repository.Select(
		`acf.bank_transfer.*,us.user_fullname AS updated_by_name,sls.emp_id as salesman_id, sls.sales_name as sales_name,ep.emp_code as salesman_code,sp.sup_name,b.bank_name,o.outlet_name,o.outlet_code`).
		Joins("left join sys.m_user us on us.user_id = acf.bank_transfer.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.bank_transfer.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_employee ep on ep.emp_id = sls.emp_id AND ep.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier sp on sp.sup_id = acf.bank_transfer.sup_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_bank b on b.bank_id = acf.bank_transfer.bank_id AND b.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet o on o.outlet_id = acf.bank_transfer.outlet_id AND o.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet_bank ob on ob.outlet_bank_id = acf.bank_transfer.outlet_bank_id AND o.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.bank_transfer.cust_id=?", dataFilter.CustId)
	query.Where("acf.bank_transfer.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.bank_transfer.is_del=false")
	query.Where("acf.bank_transfer.is_del=false")

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.bank_transfer.transfer_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.bank_transfer.transfer_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.bank_transfer.doc_no_bank ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.bank_transfer.doc_no_bank ILIKE ? )", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.BankID) > 0 {
		query.Where("acf.bank_transfer.bank_id in ?", dataFilter.BankID)
		queryCount.Where("acf.bank_transfer.bank_id in ?", dataFilter.BankID)
	}

	if len(dataFilter.AccountNo) > 0 {
		query.Where("acf.bank_transfer.account_no in ?", dataFilter.AccountNo)
		queryCount.Where("acf.bank_transfer.account_no in ?", dataFilter.AccountNo)
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
		query.Order("bank_transfer_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&BankTransfer).Error
	if err != nil {
		return BankTransfer, total, 0, err
	}
	err = queryCount.Model(&BankTransfer).Count(&total).Error
	if err != nil {
		return BankTransfer, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return BankTransfer, total, lastPage, nil
}

func (repository *RepositoryBankTransferImpl) Delete(c context.Context, custId string, BankTransferNo int, deletedBy int64) error {
	var data model.BankTransfer
	result := repository.model(c).Model(&data).Where("bank_transfer_no=? AND cust_id = ? AND is_del= ? ", BankTransferNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryBankTransferImpl) Update(c context.Context, BankTransferNo int, custId string, data model.BankTransfer) error {

	result := repository.model(c).Model(&data).Where("bank_transfer_no=? AND cust_id = ?", BankTransferNo, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}

func (repository *RepositoryBankTransferImpl) FindAllBankByCustId(dataFilter entity.BankTransferQueryFilter) ([]model.BankLookupBankTransfer, int64, int, error) {
	var BankTransfer []model.BankLookupBankTransfer
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.bank_transfer.bank_id))")
	query := repository.Select(
		`distinct(acf.bank_transfer.bank_id),mb.bank_name,mb.bank_code`).
		Joins("left join mst.m_bank mb on mb.bank_id = acf.bank_transfer.bank_id ")

	queryCount.Where("acf.bank_transfer.cust_id=?", dataFilter.CustId)
	query.Where("acf.bank_transfer.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.bank_transfer.is_del=false")
	query.Where("acf.bank_transfer.is_del=false")

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.bank_transfer.bank_name ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.bank_transfer.bank_name ILIKE ? )", "%"+dataFilter.Query+"%")
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

	err := query.Limit(limit).Offset(offset).Find(&BankTransfer).Error
	if err != nil {
		return BankTransfer, total, 0, err
	}
	err = queryCount.Model(&BankTransfer).Count(&total).Error
	if err != nil {
		return BankTransfer, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return BankTransfer, total, lastPage, nil
}

func (repository *RepositoryBankTransferImpl) FindAllBankAccountByCustId(dataFilter entity.BankTransferQueryFilter, bankID []int) ([]model.BankAccountLookupBankTransfer, int64, int, error) {
	var BankTransfer []model.BankAccountLookupBankTransfer
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("COUNT(distinct(acf.bank_transfer.account_no))")
	query := repository.Select(
		`distinct(acf.bank_transfer.account_no)`)

	if len(dataFilter.BankID) > 0 {
		query.Where("acf.bank_transfer.bank_id in ?", dataFilter.BankID)
		queryCount.Where("acf.bank_transfer.bank_id in ?", dataFilter.BankID)
	}

	queryCount.Where("acf.bank_transfer.cust_id=?", dataFilter.CustId)
	query.Where("acf.bank_transfer.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.bank_transfer.is_del=false")
	query.Where("acf.bank_transfer.is_del=false")

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(acf.bank_transfer.account_no ILIKE ? )", "%"+dataFilter.Query+"%")
		query = query.Where("(acf.bank_transfer.account_no ILIKE ? )", "%"+dataFilter.Query+"%")
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

	err := query.Limit(limit).Offset(offset).Find(&BankTransfer).Error
	if err != nil {
		return BankTransfer, total, 0, err
	}
	err = queryCount.Model(&BankTransfer).Count(&total).Error
	if err != nil {
		return BankTransfer, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return BankTransfer, total, lastPage, nil
}

func (repository *RepositoryBankTransferImpl) StoreFile(c context.Context, data *model.BankTransferFile) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryBankTransferImpl) GetLastDocNoBank(c context.Context, prefix string) (string, error) {
	var docNo string
	err := repository.model(c).Table("acf.bank_transfer").
		Select("doc_no_bank").
		Where("doc_no_bank LIKE ?", prefix+"%").
		Order("doc_no_bank DESC").
		Limit(1).
		Scan(&docNo).Error
	if err != nil {
		return "", err
	}
	return docNo, nil
}

func (repository *RepositoryBankTransferImpl) DeleteFilesByDocNo(c context.Context, docNo string, custId string) error {
	err := repository.model(c).Where("bank_transfer_no = ? AND cust_id = ?", docNo, custId).Delete(&model.BankTransferFile{}).Error
	return err
}

func (repository *RepositoryBankTransferImpl) FindOne(c context.Context, bankTransferNo int, custId string) (model.BankTransfer, error) {
	var data model.BankTransfer
	err := repository.model(c).Where("bank_transfer_no = ? AND cust_id = ?", bankTransferNo, custId).First(&data).Error
	return data, err
}

func (repository *RepositoryBankTransferImpl) FindFilesByDocNo(c context.Context, docNo string, custId string) ([]model.BankTransferFile, error) {
	var files []model.BankTransferFile
	err := repository.model(c).
		Table("acf.bank_transfer_files").
		Where("bank_transfer_no = ? AND cust_id = ?", docNo, custId).
		Order("bank_transfer_file_id ASC").
		Find(&files).Error
	return files, err
}

// FindDepositDataByDocumentNo returns deposit rows linked to this bank transfer (document_no, pay_type=3).
func (repository *RepositoryBankTransferImpl) FindDepositDataByDocumentNo(c context.Context, docNo string, custId string) ([]model.BankTransferDepositDataRow, error) {
	var rows []model.BankTransferDepositDataRow
	activeStatuses := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}
	err := repository.model(c).
		Table("acf.deposit d").
		Select("d.deposit_no, d.deposit_date, dd.invoice_no, dp.payment_amount AS used_amount").
		Joins("JOIN acf.deposit_detail dd ON dd.deposit_no = d.deposit_no AND dd.cust_id = d.cust_id").
		Joins("JOIN acf.deposit_payment dp ON dp.deposit_no = d.deposit_no AND dp.invoice_no = dd.invoice_no AND dp.document_no = ? AND dp.pay_type = 3 AND dp.cust_id = d.cust_id", docNo).
		Where("d.cust_id = ?", custId).
		Where("d.deposit_status IN ?", activeStatuses).
		Find(&rows).Error
	return rows, err
}
