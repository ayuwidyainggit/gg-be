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
	RepositoryMCoaImpl struct {
		*gorm.DB
	}
)

type MCoaRepository interface {
	Store(c context.Context, data *model.MCoa) error
	FindByID(CoaID int64, custId string) (ApCndn model.MCoaList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MCoaList, int64, int, error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter) ([]model.MCoaList, int64, int, error)
	Delete(c context.Context, custId string, CoaID int64, deletedBy int64) error
	Update(c context.Context, coaID int64, data model.MCoa, custId string) error
}

func NewMCoaRepo(db *gorm.DB) *RepositoryMCoaImpl {
	return &RepositoryMCoaImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMCoaImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMCoaImpl) Store(c context.Context, data *model.MCoa) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMCoaImpl) FindByID(CoaID int64, custId string) (mcoa model.MCoaList, err error) {
	err = repository.Select("acf.m_coa.*, us.user_fullname AS updated_by_name,c.coa_type_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_coa.updated_by").
		Joins("left join acf.m_coa_type c on c.coa_type_id = acf.m_coa.coa_type_id AND c.cust_id = ?", custId).
		Where("acf.m_coa.coa_id = ? AND acf.m_coa.cust_id=?", CoaID, custId).
		Take(&mcoa).Error
	return mcoa, err
}
func (repository *RepositoryMCoaImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MCoaList, int64, int, error) {
	var mcoa []model.MCoaList
	var total int64
	var limit, page, offset int

	queryCount := repository.Select("coa_id")
	query := repository.Select("acf.m_coa.*, us.user_fullname AS updated_by_name,c.coa_type_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_coa.updated_by").
		Joins("left join acf.m_coa_type c on c.coa_type_id = acf.m_coa.coa_type_id AND c.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.m_coa.cust_id=?", dataFilter.CustId)
	query.Where("acf.m_coa.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_coa.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.m_coa.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.m_coa.coa_code like ?", "%"+dataFilter.Query+"%")
		query.Where("acf.m_coa.coa_code like ?", "%"+dataFilter.Query+"%")

	}

	// if dataFilter.Sort != "" {
	// 	// query.Order("coa_id DESC")
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			// sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 			query.Order(colSort[0] + " " + colSort[1])
	// 		}
	// 	}
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
		query.Order("coa_id DESC")
	}

	if dataFilter.Mode != "lookup" {
		if dataFilter.Limit == 0 {
			limit = 10
		} else {
			limit = dataFilter.Limit
		}
		page = dataFilter.Page
		if page-1 < 1 {
			page = 1
		}
		offset = (page - 1) * dataFilter.Limit

		query.Limit(limit).Offset(offset)
	}

	err := query.Find(&mcoa).Error
	if err != nil {
		return mcoa, total, 0, err
	}
	err = queryCount.Model(&mcoa).Count(&total).Error
	if err != nil {
		return mcoa, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return mcoa, total, lastPage, nil
}

func (repository *RepositoryMCoaImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter) ([]model.MCoaList, int64, int, error) {
	var mcoa []model.MCoaList

	query := repository.Select("acf.m_coa.*, c.coa_type_name").
		Joins("LEFT JOIN acf.m_coa_type c ON c.coa_type_id = acf.m_coa.coa_type_id AND c.cust_id = ?", dataFilter.CustId)

	query.Where("acf.m_coa.cust_id = ?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_coa.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("acf.m_coa.coa_code LIKE ?", "%"+dataFilter.Query+"%")
	}

	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			query.Order(colSort[0] + " " + colSort[1])
	// 		}
	// 	}
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
		query.Order("coa_id DESC")
	}

	err := query.Find(&mcoa).Error
	if err != nil {
		return mcoa, 0, 0, err
	}
	return mcoa, 0, 0, nil
}

func (repository *RepositoryMCoaImpl) Delete(c context.Context, custId string, CoaID int64, deletedBy int64) error {
	var data model.MCoa
	result := repository.model(c).Model(&data).Where("coa_id=? AND cust_id = ? AND is_del= ? ", CoaID, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryMCoaImpl) Update(c context.Context, coaID int64, data model.MCoa, custId string) error {
	result := repository.model(c).Model(&data).Where("cust_id = ? AND coa_id = ?", custId, coaID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}
