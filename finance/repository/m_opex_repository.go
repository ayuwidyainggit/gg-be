package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryMOpexImpl struct {
		*gorm.DB
	}
)

type MOpexRepository interface {
	Store(c context.Context, data *model.MOpex) error
	FindById(OpexId int, custId string) (whAdj model.MOpexList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MOpexList, int64, int, error)
	Update(c context.Context, OpexId int, data model.MOpex) error
	Delete(c context.Context, custId string, OpexId int, deletedBy int64) error
	FindByCode(OpexCode string, custId string) (whAdj model.MOpexList, err error)
}

func NewMOpexRepo(db *gorm.DB) *RepositoryMOpexImpl {
	return &RepositoryMOpexImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMOpexImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMOpexImpl) Store(c context.Context, data *model.MOpex) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMOpexImpl) FindById(OpexId int, custId string) (whAdj model.MOpexList, err error) {
	err = repository.Select("acf.m_opex.*,us.user_fullname AS updated_by_name, c.coa_code, c.coa_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_opex.updated_by").
		Joins("left join acf.m_coa c on c.coa_id = acf.m_opex.coa_id AND c.cust_id = ?", custId).
		Where("acf.m_opex.opex_id = ? AND acf.m_opex.cust_id=?", OpexId, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryMOpexImpl) FindByCode(OpexCode string, custId string) (whAdj model.MOpexList, err error) {
	err = repository.Select("acf.m_opex.*,us.user_fullname AS updated_by_name, c.coa_code, c.coa_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_opex.updated_by").
		Joins("left join acf.m_coa c on c.coa_id = acf.m_opex.coa_id AND c.cust_id = ?", custId).
		Where("acf.m_opex.cust_id=? AND acf.m_opex.opex_code = ? ", custId, OpexCode).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryMOpexImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MOpexList, int64, int, error) {
	var mOpex []model.MOpexList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("opex_id")
	query := repository.Select("acf.m_opex.*,us.user_fullname AS updated_by_name, c.coa_code, c.coa_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_opex.updated_by").
		Joins("left join acf.m_coa c on c.coa_id = acf.m_opex.coa_id AND c.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.m_opex.cust_id=?", dataFilter.CustId)
	query.Where("acf.m_opex.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_opex.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.m_opex.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.m_opex.opex_name like ? or acf.m_opex.opex_code like ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query.Where("acf.m_opex.opex_name like ? or acf.m_opex.opex_code like ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")

	}

	if dataFilter.IsActive != nil {
		log.Println("is_active >>> ", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			queryCount.Where("acf.m_opex.is_active=?", true)
			query.Where("acf.m_opex.is_active=?", true)
		}
		if *dataFilter.IsActive == 2 {
			queryCount.Where("acf.m_opex.is_active=?", false)
			query.Where("acf.m_opex.is_active=?", false)
		}
	}

	// if dataFilter.Sort != "" {
	// 	sortBy, querySelect := "", ""
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	querySelect += sortBy
	// 	query.Order(querySelect)
	// } else {
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
		query.Order("opex_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&mOpex).Error
	if err != nil {
		return mOpex, total, 0, err
	}
	err = queryCount.Model(&mOpex).Count(&total).Error
	if err != nil {
		return mOpex, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return mOpex, total, lastPage, nil
}

func (repository *RepositoryMOpexImpl) Update(c context.Context, OpexId int, data model.MOpex) error {
	result := repository.model(c).Model(&data).Where("opex_id = ?", OpexId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryMOpexImpl) Delete(c context.Context, custId string, OpexId int, deletedBy int64) error {
	var data model.MOpex
	result := repository.model(c).Model(&data).Where("opex_id=? AND cust_id = ? AND is_del= ? ", OpexId, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
