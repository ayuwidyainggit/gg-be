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
	RepositoryMChequeRejectImpl struct {
		*gorm.DB
	}
)

type MChequeRejectRepository interface {
	Store(c context.Context, data *model.MChequeReject) error
	FindById(ChqRejectId int, custId string) (whAdj model.MChequeRejectList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MChequeRejectList, int64, int, error)
	Update(c context.Context, ChqRejectId int, custId string, data model.MChequeReject) error
	Delete(c context.Context, custId string, ChqRejectId int, deletedBy int64) error
}

func NewMChequeRejectRepo(db *gorm.DB) *RepositoryMChequeRejectImpl {
	return &RepositoryMChequeRejectImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMChequeRejectImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMChequeRejectImpl) Store(c context.Context, data *model.MChequeReject) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMChequeRejectImpl) FindById(ChqRejectId int, custId string) (whAdj model.MChequeRejectList, err error) {
	err = repository.Select("acf.m_cheque_reject.*, us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_cheque_reject.updated_by").
		Where("acf.m_cheque_reject.chq_reject_id = ? AND acf.m_cheque_reject.cust_id=?", ChqRejectId, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryMChequeRejectImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MChequeRejectList, int64, int, error) {
	var MChequeReject []model.MChequeRejectList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("chq_reject_id")
	query := repository.Select("acf.m_cheque_reject.*, us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_cheque_reject.updated_by")

	queryCount.Where("acf.m_cheque_reject.cust_id=?", custId)
	query.Where("acf.m_cheque_reject.cust_id=?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_cheque_reject.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.m_cheque_reject.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("acf.m_cheque_reject.chq_reject_name like ?", "%"+dataFilter.Query+"%")
		queryCount.Where("acf.m_cheque_reject.chq_reject_name like ?", "%"+dataFilter.Query+"%")

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
		query.Order("chq_reject_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&MChequeReject).Error
	if err != nil {
		return MChequeReject, total, 0, err
	}
	err = queryCount.Model(&MChequeReject).Count(&total).Error
	if err != nil {
		return MChequeReject, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return MChequeReject, total, lastPage, nil
}

func (repository *RepositoryMChequeRejectImpl) Update(c context.Context, ChqRejectId int, custId string, data model.MChequeReject) error {
	result := repository.model(c).Model(&data).Where("chq_reject_id = ? AND cust_id = ?", ChqRejectId, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryMChequeRejectImpl) Delete(c context.Context, custId string, ChqRejectId int, deletedBy int64) error {
	var data model.MChequeReject
	result := repository.model(c).Model(&data).Where("chq_reject_id=? AND cust_id = ? AND is_del= ? ", ChqRejectId, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
