package repository

import (
	"context"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/sql_helper"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ExpenseRepository interface {
	FindAllByCustId(ctx context.Context, dataFilter entity.ExpenseQueryFilter, custId string) ([]model.ExpenseListRead, int64, int, error)
	FindOneByExpenseId(ctx context.Context, expenseId int64, custId string) (model.ExpenseDetailRead, []model.ExpenseDetRead, []model.ExpenseFileRead, error)
	FindAllExpenseTypeLookup(ctx context.Context, dataFilter entity.ExpenseTypeQueryFilter) ([]model.ExpenseType, int64, int, error)
	FindAllOutletByPJP(ctx context.Context, salesmanId int64, custId string, params entity.OutletLookupQueryFilter) ([]model.OutletLookupPJP, int64, int, error)
	Store(ctx context.Context, data *model.Expense) error
	StoreDetail(ctx context.Context, data *model.ExpenseDet) error
	StoreFile(ctx context.Context, data *model.ExpenseFile) error
	Update(ctx context.Context, expenseId int64, custId string, data *model.Expense) error
	DeleteDetailsByExpenseId(ctx context.Context, expenseId int64, custId string) error
	Delete(ctx context.Context, expenseId int64, custId string, deletedBy int64) error
	FindFilesByExpenseId(ctx context.Context, expenseId int64, custId string) ([]model.ExpenseFileRead, error)
	FindFilesByExpenseFileIds(ctx context.Context, expenseId int64, expenseFileIds []int64, custId string) ([]model.ExpenseFileRead, error)
	DeleteFilesByExpenseFileIds(ctx context.Context, expenseFileIds []int64, custId string) error
	FindExpenseTypeById(ctx context.Context, expenseTypeId int) (model.ExpenseType, error)
	ValidateOutlets(ctx context.Context, outletIds []int, custId string) ([]int, error)
	FindAttendanceByUserAndDate(ctx context.Context, userId int64, custId string, date time.Time) (int, error)
	GetLastDocNo(ctx context.Context, custId string, date time.Time) (string, error)
}

type expenseRepositoryImpl struct {
	*gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
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

func (repo *expenseRepositoryImpl) FindAllByCustId(ctx context.Context, dataFilter entity.ExpenseQueryFilter, custId string) ([]model.ExpenseListRead, int64, int, error) {
	var expenses []model.ExpenseListRead
	var total int64

	query := repo.model(ctx).
		Table("acf.expense e").
		Select(`
			e.cust_id,
			e.expense_id,
			e.doc_no,
			e.expense_type_id,
			et.expense_type_name,
			e.date,
			e.amount,
			e.note,
			e.created_at,
			e.updated_at
		`).
		Joins("LEFT JOIN acf.expense_type et ON et.expense_type_id = e.expense_type_id").
		Where("e.cust_id = ? AND e.is_del = false", custId).
		Where("e.collector_id = ?", dataFilter.EmpID)

	// Date filter
	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		startDate := time.Unix(*dataFilter.StartDate, 0).Format(time.DateOnly)
		endDate := time.Unix(*dataFilter.EndDate, 0).Format(time.DateOnly)
		query = query.Where("e.date >= ? AND e.date <= ?", startDate, endDate)
	} else if dataFilter.StartDate != nil {
		startDate := time.Unix(*dataFilter.StartDate, 0).Format(time.DateOnly)
		query = query.Where("e.date >= ?", startDate)
	} else if dataFilter.EndDate != nil {
		endDate := time.Unix(*dataFilter.EndDate, 0).Format(time.DateOnly)
		query = query.Where("e.date <= ?", endDate)
	} else {
		// Default: 3 bulan terakhir jika tidak ada filter
		threeMonthsAgo := time.Now().AddDate(0, -3, 0)
		// Reset time to start of day
		threeMonthsAgo = time.Date(threeMonthsAgo.Year(), threeMonthsAgo.Month(), threeMonthsAgo.Day(), 0, 0, 0, 0, threeMonthsAgo.Location())
		query = query.Where("e.date >= ?", threeMonthsAgo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return expenses, 0, 0, err
	}

	// Sort
	sortBy := "e.created_at DESC"
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		var sortParts []string
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) == 2 {
				colName := colSort[0]
				direction := strings.ToUpper(colSort[1])
				if direction == "ASC" || direction == "DESC" {
					sortParts = append(sortParts, fmt.Sprintf("%s %s", colName, direction))
				}
			}
		}
		if len(sortParts) > 0 {
			sortBy = strings.Join(sortParts, ", ")
		}
	}
	query = query.Order(sortBy)

	// Pagination
	limit, _, offset, lastPage := repo.applyPagination(dataFilter.Limit, dataFilter.Page, total, 20)
	query = query.Limit(limit).Offset(offset)

	// Execute query
	if err := query.Scan(&expenses).Error; err != nil {
		return expenses, 0, 0, err
	}

	return expenses, total, lastPage, nil
}

func (repo *expenseRepositoryImpl) FindOneByExpenseId(ctx context.Context, expenseId int64, custId string) (model.ExpenseDetailRead, []model.ExpenseDetRead, []model.ExpenseFileRead, error) {
	var expense model.ExpenseDetailRead
	var details []model.ExpenseDetRead
	var files []model.ExpenseFileRead

	// Get expense header
	err := repo.model(ctx).
		Table("acf.expense e").
		Select(`
			e.cust_id,
			e.expense_id,
			e.doc_no,
			e.expense_type_id,
			et.expense_type_code,
			et.expense_type_name,
			e.date,
			e.amount,
			e.note,
			e.created_at,
			e.updated_at
		`).
		Joins("LEFT JOIN acf.expense_type et ON et.expense_type_id = e.expense_type_id").
		Where("e.expense_id = ? AND e.cust_id = ? AND e.is_del = false", expenseId, custId).
		Scan(&expense).Error

	if err != nil {
		return expense, details, files, err
	}

	// Get expense details (outlets)
	err = repo.model(ctx).
		Table("acf.expense_det ed").
		Select(`
			ed.expense_det_id,
			ed.expense_id,
			ed.outlet_id,
			o.outlet_code,
			o.outlet_name,
			o.address1 as outlet_address1
		`).
		Joins("LEFT JOIN mst.m_outlet o ON o.outlet_id = ed.outlet_id AND o.cust_id = ?", custId).
		Where("ed.expense_id = ? AND ed.cust_id = ?", expenseId, custId).
		Scan(&details).Error

	if err != nil {
		return expense, details, files, err
	}

	// Get expense files
	err = repo.model(ctx).
		Table("acf.expense_file ef").
		Select(`
			ef.expense_file_id,
			ef.expense_id,
			ef.file_name,
			ef.file_url,
			ef.file_key,
			ef.media_category,
			ef.file_size,
			ef.created_at
		`).
		Where("ef.expense_id = ? AND ef.cust_id = ?", expenseId, custId).
		Order("ef.created_at ASC").
		Scan(&files).Error

	if err != nil {
		return expense, details, files, err
	}

	return expense, details, files, nil
}

func (repo *expenseRepositoryImpl) FindAllExpenseTypeLookup(ctx context.Context, dataFilter entity.ExpenseTypeQueryFilter) ([]model.ExpenseType, int64, int, error) {
	var expenseTypes []model.ExpenseType
	var total int64

	query := repo.model(ctx).
		Table("acf.expense_type").
		Where("is_del = false")

	// Filter is_active
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			query = query.Where("is_active = true")
		} else if *dataFilter.IsActive == 2 {
			query = query.Where("is_active = false")
		}
	} else {
		// Default: active only
		query = query.Where("is_active = true")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return expenseTypes, 0, 0, err
	}

	// Pagination
	limit, _, offset, lastPage := repo.applyPagination(dataFilter.Limit, dataFilter.Page, total, 100)

	// Execute query
	if err := query.
		Select("expense_type_id, expense_type_code, expense_type_name").
		Order("expense_type_id ASC").
		Limit(limit).
		Offset(offset).
		Find(&expenseTypes).Error; err != nil {
		return expenseTypes, 0, 0, err
	}

	return expenseTypes, total, lastPage, nil
}

func (repo *expenseRepositoryImpl) FindAllOutletByPJP(ctx context.Context, salesmanId int64, custId string, params entity.OutletLookupQueryFilter) ([]model.OutletLookupPJP, int64, int, error) {
	var outlets []model.OutletLookupPJP
	var total int64

	if params.Page <= 0 || params.Page >= 999 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit >= 999 {
		params.Limit = 5
	}

	query := repo.model(ctx).
		Table("pjp.outlet_visit_list ovl").
		Select("DISTINCT o.outlet_id, o.outlet_code, o.outlet_name").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.pjp_code = ovl.pjp_code").
		Joins("JOIN mst.m_outlet o ON o.outlet_id = ovl.outlet_id AND o.cust_id = ?", custId).
		Where("pjp.salesman_id = ? AND o.is_del = false", salesmanId)

	// Filter is_active
	if params.IsActive != nil {
		if *params.IsActive == 1 {
			query = query.Where("o.is_active = true")
		} else if *params.IsActive == 2 {
			query = query.Where("o.is_active = false")
		}
	} else {
		// Default: active only
		query = query.Where("o.is_active = true")
	}

	if params.Search != "" {
		q := "%" + params.Search + "%"
		query = query.Where("o.outlet_code ILIKE ? OR o.outlet_name ILIKE ?", q, q)
	}

	sort := "o.outlet_code ASC"
	sortAvailable := map[string]string{
		"outlet_code": "o.outlet_code",
		"outlet_name": "o.outlet_name",
	}

	if params.Sort != "" {
		s, sn := sql_helper.ParseSort(params.Sort)
		if sq, ok := sortAvailable[s]; ok {
			sort = fmt.Sprintf("%s %s", sq, sn)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return outlets, 0, 0, err
	}

	// Pagination
	limit, _, offset, lastPage := repo.applyPagination(100, 1, total, 100)

	// Execute query
	if err := query.
		Order(sort).
		Limit(limit).
		Offset(offset).
		Find(&outlets).Error; err != nil {
		return outlets, 0, 0, err
	}

	return outlets, total, lastPage, nil
}

func (repo *expenseRepositoryImpl) Store(ctx context.Context, data *model.Expense) error {
	return repo.model(ctx).Create(data).Error
}

func (repo *expenseRepositoryImpl) StoreDetail(ctx context.Context, data *model.ExpenseDet) error {
	return repo.model(ctx).Create(data).Error
}

func (repo *expenseRepositoryImpl) StoreFile(ctx context.Context, data *model.ExpenseFile) error {
	return repo.model(ctx).Create(data).Error
}

func (repo *expenseRepositoryImpl) Update(ctx context.Context, expenseId int64, custId string, data *model.Expense) error {
	return repo.model(ctx).
		Where("expense_id = ? AND cust_id = ? AND is_del = false", expenseId, custId).
		Updates(data).Error
}

func (repo *expenseRepositoryImpl) DeleteDetailsByExpenseId(ctx context.Context, expenseId int64, custId string) error {
	return repo.model(ctx).
		Where("expense_id = ? AND cust_id = ?", expenseId, custId).
		Delete(&model.ExpenseDet{}).Error
}

func (repo *expenseRepositoryImpl) Delete(ctx context.Context, expenseId int64, custId string, deletedBy int64) error {
	now := time.Now()
	return repo.model(ctx).
		Model(&model.Expense{}).
		Where("expense_id = ? AND cust_id = ? AND is_del = false", expenseId, custId).
		Updates(map[string]interface{}{
			"is_del":     true,
			"deleted_by": deletedBy,
			"deleted_at": now,
		}).Error
}

func (repo *expenseRepositoryImpl) FindFilesByExpenseId(ctx context.Context, expenseId int64, custId string) ([]model.ExpenseFileRead, error) {
	var files []model.ExpenseFileRead
	err := repo.model(ctx).
		Table("acf.expense_file").
		Select("expense_file_id, expense_id, file_name, file_url, file_key, media_category, file_size, created_at").
		Where("expense_id = ? AND cust_id = ?", expenseId, custId).
		Scan(&files).Error
	return files, err
}

func (repo *expenseRepositoryImpl) FindFilesByExpenseFileIds(ctx context.Context, expenseId int64, expenseFileIds []int64, custId string) ([]model.ExpenseFileRead, error) {
	var files []model.ExpenseFileRead
	err := repo.model(ctx).
		Table("acf.expense_file").
		Select("expense_file_id, expense_id, file_name, file_url, file_key, media_category, file_size, created_at").
		Where("expense_file_id IN ? AND expense_id = ? AND cust_id = ?", expenseFileIds, expenseId, custId).
		Scan(&files).Error
	return files, err
}

func (repo *expenseRepositoryImpl) DeleteFilesByExpenseFileIds(ctx context.Context, expenseFileIds []int64, custId string) error {
	return repo.model(ctx).
		Where("expense_file_id IN ? AND cust_id = ?", expenseFileIds, custId).
		Delete(&model.ExpenseFile{}).Error
}

func (repo *expenseRepositoryImpl) FindExpenseTypeById(ctx context.Context, expenseTypeId int) (model.ExpenseType, error) {
	var expenseType model.ExpenseType
	err := repo.model(ctx).
		Table("acf.expense_type").
		Where("expense_type_id = ? AND is_del = false AND is_active = true", expenseTypeId).
		First(&expenseType).Error
	return expenseType, err
}

func (repo *expenseRepositoryImpl) ValidateOutlets(ctx context.Context, outletIds []int, custId string) ([]int, error) {
	var validOutletIds []int
	err := repo.model(ctx).
		Table("mst.m_outlet").
		Select("outlet_id").
		Where("outlet_id IN ? AND cust_id = ? AND is_del = false AND is_active = true", outletIds, custId).
		Pluck("outlet_id", &validOutletIds).Error
	return validOutletIds, err
}

func (repo *expenseRepositoryImpl) FindAttendanceByUserAndDate(ctx context.Context, userId int64, custId string, date time.Time) (int, error) {
	var attendanceType int
	// Start of day
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	// End of day
	endDate := startDate.Add(24 * time.Hour)

	err := repo.model(ctx).
		Table("mobile.attendances a").
		Select("a.type").
		Joins("JOIN mst.m_employee e ON e.emp_code = a.emp_code AND e.cust_id = a.cust_id").
		Joins("JOIN sys.m_user u ON u.emp_id = e.emp_id AND u.cust_id = e.cust_id").
		Where("u.user_id = ? AND u.cust_id = ?", userId, custId).
		Where("a.created_at >= ? AND a.created_at < ?", startDate, endDate).
		Order("a.created_at DESC").
		Limit(1).
		Scan(&attendanceType).Error

	if err != nil {
		return 0, err
	}
	return attendanceType, nil
}

func (repo *expenseRepositoryImpl) GetLastDocNo(ctx context.Context, custId string, date time.Time) (string, error) {
	var docNo string
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	err := repo.model(ctx).
		Table("acf.expense").
		Select("doc_no").
		Where("cust_id = ? AND date = ?", custId, dateOnly).
		Order("doc_no DESC").
		Limit(1).
		Scan(&docNo).Error

	return docNo, err
}

// applyPagination calculates pagination parameters (limit, page, offset, lastPage)
func (repo *expenseRepositoryImpl) applyPagination(limit, page int, total int64, defaultLimit int) (int, int, int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return limit, page, offset, lastPage
}
