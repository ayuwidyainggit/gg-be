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
	RepositoryMemoJrImpl struct {
		*gorm.DB
	}
)

type MemoJrRepository interface {
	Store(c context.Context, data *model.MemoJr) error
	StoreDetail(c context.Context, data *model.MemoJrDet) error
	FindByNo(MjNo string, custId string) (memoJr model.MemoJrList, err error)
	FindDetail(MjNo string, custId string) (Details []model.MemoJrDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MemoJrList, int64, int, error)
	Delete(c context.Context, custId string, MjNo string, deletedBy int64) error
	Update(c context.Context, MjNo string, data model.MemoJr) error
	DeleteDetailNotInIDs(c context.Context, MjNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.MemoJrDet) error
}

func NewMemoJrRepo(db *gorm.DB) *RepositoryMemoJrImpl {
	return &RepositoryMemoJrImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMemoJrImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryMemoJrImpl) Store(c context.Context, data *model.MemoJr) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMemoJrImpl) StoreDetail(c context.Context, data *model.MemoJrDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryMemoJrImpl) FindByNo(MjNo string, custId string) (memoJr model.MemoJrList, err error) {
	err = repository.Select("acf.memo_jr.*, us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.memo_jr.updated_by").
		Where("mj_no = ? AND acf.memo_jr.cust_id=?", MjNo, custId).
		Take(&memoJr).Error
	return memoJr, err
}

func (repository *RepositoryMemoJrImpl) FindDetail(MjNo string, custId string) (Details []model.MemoJrDetRead, err error) {
	err = repository.Select("acf.memo_jr_det.*, c.coa_code, c.coa_name").
		Joins("left join acf.m_coa c on c.coa_id = acf.memo_jr_det.coa_id AND c.cust_id = ?", custId).
		Where("mj_no = ? AND acf.memo_jr_det.cust_id=?", MjNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryMemoJrImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MemoJrList, int64, int, error) {
	var memoJr []model.MemoJrList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("mj_no")
	query := repository.Select("acf.memo_jr.*, us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.memo_jr.updated_by")

	queryCount.Where("acf.memo_jr.cust_id=?", dataFilter.CustId)
	query.Where("acf.memo_jr.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.memo_jr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.memo_jr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.memo_jr.mj_no=?", dataFilter.Query)
		query.Where("acf.memo_jr.mj_no=?", dataFilter.Query)
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
		query.Order("mj_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&memoJr).Error
	if err != nil {
		return memoJr, total, 0, err
	}
	err = queryCount.Model(&memoJr).Count(&total).Error
	if err != nil {
		return memoJr, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return memoJr, total, lastPage, nil
}
func (repository *RepositoryMemoJrImpl) Delete(c context.Context, custId string, MjNo string, deletedBy int64) error {
	var data model.MemoJr
	result := repository.model(c).Model(&data).Where("mj_no=? AND cust_id = ? AND is_del= ? ", MjNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryMemoJrImpl) Update(c context.Context, MjNo string, data model.MemoJr) error {
	result := repository.model(c).Model(&data).Where("mj_no=?", MjNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected lala")
	// }
	return nil
}

func (repository *RepositoryMemoJrImpl) DeleteDetailNotInIDs(c context.Context, MjNo string, IDs []int64) error {
	var Details model.MemoJrDet
	err := repository.model(c).Where("mj_no=? AND memo_jr_det_id not in (?) ", MjNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryMemoJrImpl) UpdateDetail(c context.Context, Details *model.MemoJrDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
