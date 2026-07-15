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
	RepositoryMApDiscImpl struct {
		*gorm.DB
	}
)
type MApDiscRepository interface {
	Store(c context.Context, data *model.MApDisc) error
	FindByNo(apDiscID int64, custId string) (consg model.MApDiscList, err error)
	FindByProId(proId int64, custId string) (consg model.MApDiscList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MApDiscList, int64, int, error)
	Delete(c context.Context, custId string, apDiscID int64, deletedBy int64) error
	Update(c context.Context, apDiscID int64, data model.MApDisc) error
}

func NewMApDiscRepo(db *gorm.DB) *RepositoryMApDiscImpl {
	return &RepositoryMApDiscImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMApDiscImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryMApDiscImpl) Store(c context.Context, data *model.MApDisc) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMApDiscImpl) FindByNo(apDiscID int64, custId string) (mApDisc model.MApDiscList, err error) {
	err = repository.Select("acf.m_ap_disc.*,us.user_fullname AS updated_by_name,p.pro_code,p.pro_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.m_ap_disc.updated_by").
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = acf.m_ap_disc.pro_id").
		Where("acf.m_ap_disc.ap_disc_id = ? AND acf.m_ap_disc.cust_id = ?", apDiscID, custId).
		Take(&mApDisc).Error
	return mApDisc, err
}

func (repository *RepositoryMApDiscImpl) FindByProId(proId int64, custId string) (mApDisc model.MApDiscList, err error) {
	err = repository.Select("acf.m_ap_disc.*,us.user_fullname AS updated_by_name,p.pro_code,p.pro_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.m_ap_disc.updated_by").
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = acf.m_ap_disc.pro_id").
		Where("acf.m_ap_disc.pro_id = ? AND acf.m_ap_disc.cust_id = ?", proId, custId).
		Take(&mApDisc).Error
	return mApDisc, err
}

func (repository *RepositoryMApDiscImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MApDiscList, int64, int, error) {
	var apDisc []model.MApDiscList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ap_disc_id")
	query := repository.Select("acf.m_ap_disc.*,us.user_fullname AS updated_by_name,p.pro_code,p.pro_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_ap_disc.updated_by").
		Joins("left join mst.m_product p on p.pro_id = acf.m_ap_disc.pro_id AND p.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.m_ap_disc.cust_id=?", dataFilter.CustId)
	query.Where("acf.m_ap_disc.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_ap_disc.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.m_ap_disc.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {

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
		query.Order("ap_disc_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&apDisc).Error
	if err != nil {
		return apDisc, total, 0, err
	}
	err = queryCount.Model(&apDisc).Count(&total).Error
	if err != nil {
		return apDisc, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return apDisc, total, lastPage, nil
}
func (repository *RepositoryMApDiscImpl) Delete(c context.Context, custId string, apDiscID int64, deletedBy int64) error {
	var data model.MApDisc
	result := repository.model(c).Model(&data).Where("ap_disc_id=? AND cust_id = ? AND is_del= ? ", apDiscID, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryMApDiscImpl) Update(c context.Context, apDiscID int64, data model.MApDisc) error {
	result := repository.model(c).Model(&data).Where("ap_disc_id=?", apDiscID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
