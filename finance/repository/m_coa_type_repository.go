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
	RepositoryMCoaTypeImpl struct {
		*gorm.DB
	}
)

type MCoaTypeRepository interface {
	Store(c context.Context, data *model.MCoaType) error
	FindByID(CoaTypeID int64) (ApCndn model.MCoaTypeList, err error)
	FindAll(dataFilter entity.GeneralQueryFilter) ([]model.MCoaTypeList, int64, int, error)
	FindAllLookup(dataFilter entity.GeneralQueryFilter) ([]model.MCoaTypeList, int64, int, error)
	Delete(c context.Context, CoaTypeID int64, deletedBy int64) error
	Update(c context.Context, CoaTypeID int64, data model.MCoaType) error
}

func NewMCoaTypeRepo(db *gorm.DB) *RepositoryMCoaTypeImpl {
	return &RepositoryMCoaTypeImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryMCoaTypeImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryMCoaTypeImpl) Store(c context.Context, data *model.MCoaType) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryMCoaTypeImpl) FindByID(CoaTypeID int64) (coaType model.MCoaTypeList, err error) {
	err = repository.Select("acf.m_coa_type.*, us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_coa_type.updated_by").
		Where("coa_type_id = ?", CoaTypeID).
		Take(&coaType).Error
	return coaType, err
}

func (repository *RepositoryMCoaTypeImpl) FindAll(dataFilter entity.GeneralQueryFilter) ([]model.MCoaTypeList, int64, int, error) {
	var coaTypes []model.MCoaTypeList
	var total int64
	var limit, page, offset int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("coa_type_id")
	query := repository.Select("acf.m_coa_type.*, us.user_fullname AS updated_by_name").
		Joins("left join sys.m_user us on us.user_id = acf.m_coa_type.updated_by")

	// if dataFilter.ParentCustId != "" {
	// 	query.Where("acf.m_coa_type.cust_id = ?", dataFilter.ParentCustId)
	// 	queryCount.Where("acf.m_coa_type.cust_id = ?", dataFilter.ParentCustId)
	// }

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_coa_type.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.m_coa_type.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.m_coa_type.coa_type_name like ?", "%"+dataFilter.Query+"%")
		query.Where("acf.m_coa_type.coa_type_name like ?", "%"+dataFilter.Query+"%")

	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			queryCount.Where("acf.m_coa_type.is_active = ?", true)
			query.Where("acf.m_coa_type.is_active = ?", true)
		}
		if *dataFilter.IsActive == 2 {
			queryCount.Where("acf.m_coa_type.is_active = ?", false)
			query.Where("acf.m_coa_type.is_active = ?", false)
		}
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
		query.Order("coa_type_id DESC")
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

	err := query.Find(&coaTypes).Error
	if err != nil {
		return coaTypes, total, 0, err
	}
	err = queryCount.Model(&coaTypes).Count(&total).Error
	if err != nil {
		return coaTypes, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return coaTypes, total, lastPage, nil
}

func (repository *RepositoryMCoaTypeImpl) FindAllLookup(dataFilter entity.GeneralQueryFilter) ([]model.MCoaTypeList, int64, int, error) {
	var coaTypes []model.MCoaTypeList

	query := repository.Select(`
			acf.m_coa_type.coa_type_id, 
			acf.m_coa_type.coa_type_name, 
			acf.m_coa_type.coa_group, 
			acf.m_coa_type.def_blc, 
			acf.m_coa_type.sort_index, 
			acf.m_coa_type.coa_kind`)

	query.Where(`acf.m_coa_type.cust_id = ? 
				AND acf.m_coa_type.is_active = ?
				AND acf.m_coa_type.is_del = ?`,
		dataFilter.CustId, true, false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.m_coa_type.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("acf.m_coa_type.coa_type_name LIKE ?", "%"+dataFilter.Query+"%")

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
		query.Order("coa_type_id DESC")
	}

	err := query.Find(&coaTypes).Error
	if err != nil {
		return coaTypes, 0, 0, err
	}

	return coaTypes, 0, 0, nil
}

func (repository *RepositoryMCoaTypeImpl) Delete(c context.Context, CoaTypeID int64, deletedBy int64) error {
	var data model.MCoaType
	result := repository.model(c).Model(&data).Where("coa_type_id=? AND is_del= ? ", CoaTypeID, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryMCoaTypeImpl) Update(c context.Context, CoaTypeID int64, data model.MCoaType) error {
	result := repository.model(c).Model(&data).Where("coa_type_id=?", CoaTypeID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
