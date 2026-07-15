package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ExpenseRepository interface {
	FindAll(ctx context.Context, dataFilter entity.ExpenseQueryFilter) ([]model.ExpenseTypeList, int64, int, error)
	FindById(ctx context.Context, custId string, expenseTypeId int) (model.ExpenseTypeList, error)
	FindByCodeAndName(ctx context.Context, custId, code, name string) (model.ExpenseTypeList, error)
	Store(ctx context.Context, data *model.ExpenseType) error
	Update(ctx context.Context, custId string, expenseTypeId int, data map[string]interface{}) error
	Delete(ctx context.Context, custId string, expenseTypeId int, deletedBy int64) error
}

type expenseRepositoryImpl struct {
	*gorm.DB
}

func NewExpenseRepo(db *gorm.DB) ExpenseRepository {
	return &expenseRepositoryImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *expenseRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repo *expenseRepositoryImpl) FindAll(ctx context.Context, dataFilter entity.ExpenseQueryFilter) ([]model.ExpenseTypeList, int64, int, error) {
	var expenseTypes []model.ExpenseTypeList
	var total int64

	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 5
	}

	queryCount := repo.buildExpenseTypeListBaseQuery(ctx, dataFilter).
		Select("expense_type_id")
	query := repo.buildExpenseTypeListBaseQuery(ctx, dataFilter).
		Select("acf.expense_type.*, us.user_fullname AS updated_by_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.expense_type.updated_by")

	// Sort
	sortBy := "acf.expense_type.created_at DESC"
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		var sortParts []string
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) == 2 {
				colName := colSort[0]
				direction := strings.ToUpper(colSort[1])
				if direction == "ASC" || direction == "DESC" {
					// Map created_date to created_at
					if colName == "created_date" {
						colName = "acf.expense_type.created_at"
					} else {
						colName = fmt.Sprintf("acf.expense_type.%s", colName)
					}
					sortParts = append(sortParts, fmt.Sprintf("%s %s", colName, direction))
				}
			}
		}
		if len(sortParts) > 0 {
			sortBy = strings.Join(sortParts, ", ")
		}
	}
	query = query.Order(sortBy)

	// Count total
	if err := queryCount.Model(&model.ExpenseType{}).Count(&total).Error; err != nil {
		return expenseTypes, 0, 0, err
	}

	// Pagination
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Execute query
	if err := query.Limit(limit).Offset(offset).Find(&expenseTypes).Error; err != nil {
		return expenseTypes, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return expenseTypes, total, lastPage, nil
}

func (repo *expenseRepositoryImpl) buildExpenseTypeListBaseQuery(ctx context.Context, dataFilter entity.ExpenseQueryFilter) *gorm.DB {
	query := repo.model(ctx).
		Table("acf.expense_type").
		Where("acf.expense_type.is_del = false AND acf.expense_type.cust_id = ?", dataFilter.ParentCustId)

	if dataFilter.Q != "" {
		searchPattern := "%" + dataFilter.Q + "%"
		query = query.Where("(acf.expense_type.expense_type_code LIKE ? OR acf.expense_type.expense_type_name LIKE ?)", searchPattern, searchPattern)
	}

	return query
}

func (repo *expenseRepositoryImpl) FindById(ctx context.Context, custId string, expenseTypeId int) (model.ExpenseTypeList, error) {
	var expenseType model.ExpenseTypeList
	err := repo.model(ctx).
		Table("acf.expense_type").
		Select("acf.expense_type.*, us.user_fullname AS updated_by_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.expense_type.updated_by").
		Where("acf.expense_type.expense_type_id = ? AND acf.expense_type.cust_id = ? AND acf.expense_type.is_del = false", expenseTypeId, custId).
		Take(&expenseType).Error
	if err != nil {
		return expenseType, err
	}
	return expenseType, nil
}

func (repo *expenseRepositoryImpl) FindByCodeAndName(ctx context.Context, custId, code, name string) (model.ExpenseTypeList, error) {
	var expenseType model.ExpenseTypeList
	err := repo.model(ctx).
		Table("acf.expense_type").
		Select("acf.expense_type.*, us.user_fullname AS updated_by_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.expense_type.updated_by").
		Where("acf.expense_type.expense_type_code = ? AND acf.expense_type.expense_type_name = ? AND acf.expense_type.cust_id = ? AND acf.expense_type.is_del = false", code, name, custId).
		Take(&expenseType).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return expenseType, gorm.ErrRecordNotFound
	}
	if err != nil {
		return expenseType, err
	}
	return expenseType, nil
}

func (repo *expenseRepositoryImpl) Store(ctx context.Context, data *model.ExpenseType) error {
	// Use GORM Create with map to explicitly set is_active even when false
	// This ensures PostgreSQL will use the provided value instead of DEFAULT true
	return repo.model(ctx).
		Table("acf.expense_type").
		Create(map[string]interface{}{
			"cust_id":           data.CustID,
			"expense_type_code": data.ExpenseTypeCode,
			"expense_type_name": data.ExpenseTypeName,
			"is_active":         data.IsActive,
			"created_by":        data.CreatedBy,
			"is_del":            data.IsDel,
		}).Error
}

func (repo *expenseRepositoryImpl) Update(ctx context.Context, custId string, expenseTypeId int, data map[string]interface{}) error {
	// Use map to update fields, which allows updating boolean false values correctly
	result := repo.model(ctx).
		Table("acf.expense_type").
		Model(&model.ExpenseType{}).
		Where("acf.expense_type.expense_type_id = ? AND acf.expense_type.cust_id = ? AND acf.expense_type.is_del = false", expenseTypeId, custId).
		Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *expenseRepositoryImpl) Delete(ctx context.Context, custId string, expenseTypeId int, deletedBy int64) error {
	now := time.Now()
	result := repo.model(ctx).
		Model(&model.ExpenseType{}).
		Table("acf.expense_type").
		Where("acf.expense_type.expense_type_id = ? AND acf.expense_type.cust_id = ? AND acf.expense_type.is_del = false", expenseTypeId, custId).
		Updates(map[string]interface{}{
			"is_del":     true,
			"deleted_by": deletedBy,
			"deleted_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
