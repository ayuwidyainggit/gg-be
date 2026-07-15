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
	RepositoryCashImpl struct {
		*gorm.DB
	}
)

type CashRepository interface {
	Store(c context.Context, data *model.CashTr) error
	StoreDetail(c context.Context, data *model.CashTrDet) error
	FindByNo(CashTrNo string, custId string) (consg model.CashTrList, err error)
	FindDetail(CashTrNo string, custId string) (Details []model.CashTrDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.CashTrList, int64, int, error)
	Delete(c context.Context, custId string, CashTrNo string, deletedBy int64) error
	Update(c context.Context, CashTrNo string, data model.CashTr) error
	DeleteDetailNotInIDs(c context.Context, CashTrNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.CashTrDet) error
}

func NewCashRepo(db *gorm.DB) *RepositoryCashImpl {
	return &RepositoryCashImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryCashImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryCashImpl) Store(c context.Context, data *model.CashTr) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryCashImpl) StoreDetail(c context.Context, data *model.CashTrDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryCashImpl) FindByNo(cashTrNo string, custId string) (consg model.CashTrList, err error) {
	err = repository.Select("acf.cash_tr.*, us.user_fullname AS updated_by_name,c.coa_code as coa_code_to, c.coa_name as coa_name_to").
		Joins("left join sys.m_user us on us.user_id = acf.cash_tr.updated_by").
		Joins("left join acf.m_coa c on c.coa_id = acf.cash_tr.coa_id_to AND c.cust_id = ?", custId).
		Where("cash_tr_no = ? AND acf.cash_tr.cust_id=?", cashTrNo, custId).
		Take(&consg).Error
	return consg, err
}

func (repository *RepositoryCashImpl) FindDetail(cashTrNo string, custId string) (details []model.CashTrDetRead, err error) {
	err = repository.Select("acf.cash_tr_det.*, c.coa_code, c.coa_name").
		Joins("left join acf.m_coa c on c.coa_id = acf.cash_tr_det.coa_id AND c.cust_id = ?", custId).
		Where("cash_tr_no = ? AND acf.cash_tr_det.cust_id=?", cashTrNo, custId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryCashImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.CashTrList, int64, int, error) {
	var Cashtr []model.CashTrList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("cash_tr_no")
	query := repository.Select(
		`acf.cash_tr.*, 
			us.user_fullname AS updated_by_name,
			c.coa_code as coa_code_to, c.coa_name as coa_name_to`).
		Joins("left join sys.m_user us on us.user_id = acf.cash_tr.updated_by").
		Joins("left join acf.m_coa c on c.coa_id = acf.cash_tr.coa_id_to AND c.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.cash_tr.cust_id=?", dataFilter.CustId)
	query.Where("acf.cash_tr.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.cash_tr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.cash_tr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.cash_tr.cash_tr_no=?", dataFilter.Query)
		query.Where("acf.cash_tr.cash_tr_no=?", dataFilter.Query)
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
		query.Order("cash_tr_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&Cashtr).Error
	if err != nil {
		return Cashtr, total, 0, err
	}
	err = queryCount.Model(&Cashtr).Count(&total).Error
	if err != nil {
		return Cashtr, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return Cashtr, total, lastPage, nil
}

func (repository *RepositoryCashImpl) Delete(c context.Context, custId string, CashTrNo string, deletedBy int64) error {
	var data model.CashTr
	result := repository.model(c).Model(&data).Where("cash_tr_no=? AND cust_id = ? AND is_del= ? ", CashTrNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryCashImpl) Update(c context.Context, CashTrNo string, data model.CashTr) error {
	result := repository.model(c).Model(&data).Where("cash_tr_no=?", CashTrNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected lala")
	// }
	return nil
}

func (repository *RepositoryCashImpl) DeleteDetailNotInIDs(c context.Context, CashTrNo string, IDs []int64) error {
	var Details model.CashTrDet
	err := repository.model(c).Where("cash_tr_no=? AND cash_tr_det_id not in (?) ", CashTrNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryCashImpl) UpdateDetail(c context.Context, Details *model.CashTrDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
