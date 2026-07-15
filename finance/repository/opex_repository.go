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
	RepositoryOpexImpl struct {
		*gorm.DB
	}
)

type OpexRepository interface {
	Store(c context.Context, data *model.OpexTr) error
	StoreDetail(c context.Context, data *model.OpexTrDet) error
	FindByNo(opexTrNo string, custId string) (consg model.OpexTrList, err error)
	FindDetail(opexTrNo string, custId string) (Details []model.OpexTrDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.OpexTrList, int64, int, error)
	Delete(c context.Context, custId string, opexTrNo string, deletedBy int64) error
	Update(c context.Context, opexTrNo string, data model.OpexTr) error
	DeleteDetailNotInIDs(c context.Context, opexTrNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.OpexTrDet) error
}

func NewOpexRepo(db *gorm.DB) *RepositoryOpexImpl {
	return &RepositoryOpexImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOpexImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryOpexImpl) Store(c context.Context, data *model.OpexTr) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOpexImpl) StoreDetail(c context.Context, data *model.OpexTrDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryOpexImpl) FindByNo(opexTrNo string, custId string) (consg model.OpexTrList, err error) {
	err = repository.Select("acf.opex_tr.*,us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.opex_tr.updated_by").
		Where("acf.opex_tr.opex_tr_no = ? AND acf.opex_tr.cust_id=?", opexTrNo, custId).
		Take(&consg).Error
	return consg, err
}

func (repository *RepositoryOpexImpl) FindDetail(opexTrNo string, custId string) (Details []model.OpexTrDetRead, err error) {
	err = repository.Select("acf.opex_tr_det.*, o.opex_code,o.opex_name").
		Joins("left join acf.m_opex o on o.opex_id = acf.opex_tr_det.opex_id AND o.cust_id = ?", custId).
		Where("opex_tr_no = ? AND acf.opex_tr_det.cust_id=?", opexTrNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryOpexImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.OpexTrList, int64, int, error) {
	var opextr []model.OpexTrList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("opex_tr_no")
	query := repository.Select("acf.opex_tr.*,us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.opex_tr.updated_by")

	queryCount.Where("acf.opex_tr.cust_id=?", dataFilter.CustId)
	query.Where("acf.opex_tr.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.opex_tr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.opex_tr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.opex_tr.opex_tr_no=?", dataFilter.Query)
		query.Where("acf.opex_tr.opex_tr_no=?", dataFilter.Query)
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
		query.Order("opex_tr_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&opextr).Error
	if err != nil {
		return opextr, total, 0, err
	}
	err = queryCount.Model(&opextr).Count(&total).Error
	if err != nil {
		return opextr, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return opextr, total, lastPage, nil
}
func (repository *RepositoryOpexImpl) Delete(c context.Context, custId string, opexTrNo string, deletedBy int64) error {
	var data model.OpexTr
	result := repository.model(c).Model(&data).Where("opex_tr_no=? AND cust_id = ? AND is_del= ? ", opexTrNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryOpexImpl) Update(c context.Context, opexTrNo string, data model.OpexTr) error {
	result := repository.model(c).Model(&data).Where("opex_tr_no=?", opexTrNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryOpexImpl) DeleteDetailNotInIDs(c context.Context, opexTrNo string, IDs []int64) error {
	var Details model.OpexTrDet
	err := repository.model(c).Where("opex_tr_no=? AND opex_tr_det_id not in (?) ", opexTrNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryOpexImpl) UpdateDetail(c context.Context, Details *model.OpexTrDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
