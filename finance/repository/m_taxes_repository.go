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
	RepositoryMTaxesImpl struct {
		*gorm.DB
	}
)

type MTaxesRepository interface {
	Store(c context.Context, data *model.MTaxes) error
	GetByYearAndSequence(custId string, year int, sequence int) (tax model.MTaxes, err error)
	GetByID(custId string, id int64) (tax model.MTaxes, err error)
	FindAllByCustId(dataFilter entity.MTaxQueryFilter) ([]model.MTaxes, int64, int, error)
	Delete(c context.Context, custId string, id int64, deletedBy int64) error
	Update(c context.Context, mTaxID int64, data model.MTaxes) error
	GetByTransactionStatusAndSerial(custId string, year int, transactionStatusCode, serialCode string, from, to int) (tax model.MTaxes, err error)
	GetNewestSerialByStatus(custId string, status int, year int) (mtax model.MTaxes, err error)
	UpdateAfterGenerate(c context.Context, mTaxID int64, remainingQty int, status int, lastGenerated string) error
	GetSeriesFromNoRange(custId string, from, to, year int) (tax model.MTaxes, err error)
	GetSeriesToNoRange(custId string, from, to, year int) (tax model.MTaxes, err error)
	GetLastSequenceByYear(custId string, year int) (mtax model.MTaxes, err error)
}

func NewMTaxesRepo(db *gorm.DB) *RepositoryMTaxesImpl {
	return &RepositoryMTaxesImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMTaxesImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMTaxesImpl) Store(c context.Context, data *model.MTaxes) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMTaxesImpl) GetByYearAndSequence(custId string, year int, sequence int) (tax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND year = ? AND sequence = ?", custId, year, sequence).Take(&tax).
		Error

	return tax, err
}

func (repository *RepositoryMTaxesImpl) GetByTransactionStatusAndSerial(custId string, year int, transactionStatusCode, serialCode string, from, to int) (tax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND year = ? AND transaction_status_code = ? AND serial_code = ? AND acf.m_taxes.to between ? AND ?", custId, year, transactionStatusCode, serialCode, from, to).
		Order("m_tax_id DESC").Take(&tax).
		Error

	return tax, err
}

func (repository *RepositoryMTaxesImpl) GetByID(custId string, id int64) (tax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND m_tax_id = ?", custId, id).Take(&tax).
		Error

	return tax, err
}

func (repository *RepositoryMTaxesImpl) FindAllByCustId(dataFilter entity.MTaxQueryFilter) ([]model.MTaxes, int64, int, error) {
	var mTaxes []model.MTaxes
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("m_tax_id")
	query := repository.Select("*")

	queryCount.Where("acf.m_taxes.cust_id=?", dataFilter.CustId)
	query.Where("acf.m_taxes.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_taxes.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.m_taxes.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Status != 0 {
		queryCount.Where("acf.m_taxes.status=?", dataFilter.Status)
		query.Where("acf.m_taxes.status=?", dataFilter.Status)
	}
	if dataFilter.Year != 0 {
		queryCount.Where("acf.m_taxes.year=?", dataFilter.Year)
		query.Where("acf.m_taxes.year=?", dataFilter.Year)
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
		query.Order("m_tax_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&mTaxes).Error
	if err != nil {
		return mTaxes, total, 0, err
	}
	err = queryCount.Model(&mTaxes).Count(&total).Error
	if err != nil {
		return mTaxes, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return mTaxes, total, lastPage, nil
}

func (repository *RepositoryMTaxesImpl) Delete(c context.Context, custId string, id int64, deletedBy int64) error {
	var data model.MTaxes
	result := repository.model(c).Model(&data).Where("cust_id = ? AND m_tax_id = ? AND is_del= ? ", custId, id, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryMTaxesImpl) Update(c context.Context, mTaxID int64, data model.MTaxes) error {
	result := repository.model(c).Model(&data).Where("m_tax_id = ?", mTaxID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryMTaxesImpl) GetNewestSerialByStatus(custId string, status int, year int) (mtax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND year = ? AND status = ?", custId, year, status).
		Order("serial_from ASC").Take(&mtax).
		Error

	return mtax, err
}

func (repository *RepositoryMTaxesImpl) UpdateAfterGenerate(c context.Context, mTaxID int64, remainingQty int, status int, lastGenerated string) error {
	var data model.MTaxes
	result := repository.model(c).Model(&data).Where("m_tax_id = ?", mTaxID).Updates(map[string]interface{}{"remaining_qty": remainingQty, "status": status, "last_generated_tax": lastGenerated})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryMTaxesImpl) GetSeriesFromNoRange(custId string, from, to, year int) (tax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND serial_from between ? AND ?  AND year = ?", custId, from, to, year).Take(&tax).
		Error

	return tax, err
}

func (repository *RepositoryMTaxesImpl) GetSeriesToNoRange(custId string, from, to, year int) (tax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND serial_to between ? AND ?  AND year = ?", custId, from, to, year).Take(&tax).
		Error

	return tax, err
}

func (repository *RepositoryMTaxesImpl) GetLastSequenceByYear(custId string, year int) (mtax model.MTaxes, err error) {
	err = repository.
		Where("cust_id = ? AND year = ?", custId, year).
		Order("sequence ASC").Take(&mtax).
		Error

	return mtax, err
}
