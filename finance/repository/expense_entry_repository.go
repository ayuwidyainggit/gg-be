package repository

import (
	"context"
	"finance/entity"
	"finance/model"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ExpenseEntryRepository interface {
	FindAll(ctx context.Context, filter entity.ExpenseEntryQueryFilter) ([]model.ExpenseList, int64, int, error)
	FindById(ctx context.Context, custId string, expenseId int64) (model.Expense, error)
	FindDetailById(ctx context.Context, custId string, expenseId int64) (model.ExpenseList, error)
	FindDepositExpensesByExpenseId(ctx context.Context, custId string, expenseId int64) ([]model.DepositExpense, error)
	FindExpenseFilesByExpenseId(ctx context.Context, custId string, expenseId int64) ([]model.ExpenseFile, error)
	CountExpensesInCurrentMonth(ctx context.Context) (int, error)
	Store(ctx context.Context, data *model.Expense) error
	StoreExpenseFiles(ctx context.Context, custId string, expenseId int64, files []model.ExpenseFile) error
	Update(ctx context.Context, custId string, expenseId int64, data map[string]interface{}) error
	Delete(ctx context.Context, custId string, expenseId int64, deletedBy int64) error
}

const (
	expenseEntryCollectorSelect = "emp.emp_name AS collector_name"
	expenseEntryCollectorJoin   = "LEFT JOIN mst.m_employee emp ON emp.emp_id = acf.expense.collector_id AND emp.cust_id = ?"
)

type expenseEntryRepositoryImpl struct {
	*gorm.DB
}

func NewExpenseEntryRepo(db *gorm.DB) ExpenseEntryRepository {
	return &expenseEntryRepositoryImpl{db}
}

func (repo *expenseEntryRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repo *expenseEntryRepositoryImpl) FindAll(ctx context.Context, filter entity.ExpenseEntryQueryFilter) ([]model.ExpenseList, int64, int, error) {
	var expenses []model.ExpenseList
	var total int64

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 9999 {
		limit = 9999
	}

	queryCount := repo.buildExpenseEntryCountQuery(ctx, filter)
	query := repo.buildExpenseEntryListQuery(ctx, filter)

	if err := queryCount.Model(&model.Expense{}).Count(&total).Error; err != nil {
		return expenses, 0, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	if err := query.Limit(limit).Offset(offset).Find(&expenses).Error; err != nil {
		return expenses, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return expenses, total, lastPage, nil
}

func (repo *expenseEntryRepositoryImpl) buildExpenseEntryListQuery(ctx context.Context, filter entity.ExpenseEntryQueryFilter) *gorm.DB {
	query := repo.model(ctx).Table("acf.expense").
		Select(`acf.expense.*, et.expense_type_code, et.expense_type_name, ` + expenseEntryCollectorSelect)
	query = repo.applyExpenseEntryListFilters(query, filter)
	return query.Order(repo.buildSafeSort(filter.Sort))
}

func (repo *expenseEntryRepositoryImpl) buildExpenseEntryCountQuery(ctx context.Context, filter entity.ExpenseEntryQueryFilter) *gorm.DB {
	query := repo.model(ctx).Table("acf.expense").Select("expense_id")
	return repo.applyExpenseEntryListFilters(query, filter)
}

func (repo *expenseEntryRepositoryImpl) applyExpenseEntryListFilters(query *gorm.DB, filter entity.ExpenseEntryQueryFilter) *gorm.DB {
	query = query.Where("acf.expense.doc_no IS NOT NULL")
	query = query.Joins("LEFT JOIN acf.expense_type et ON et.expense_type_id = acf.expense.expense_type_id")
	query = query.Joins(expenseEntryCollectorJoin, filter.CustID)
	query = query.Where("acf.expense.cust_id = ?", filter.CustID)
	query = query.Where("acf.expense.is_del = false")

	if trimmedQuery := strings.TrimSpace(filter.Query); trimmedQuery != "" {
		query = query.Where("acf.expense.doc_no ILIKE ?", "%"+trimmedQuery+"%")
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		startDateParsed, err := time.Parse("2006-01-02", filter.StartDate)
		if err == nil {
			endDateParsed, err := time.Parse("2006-01-02", filter.EndDate)
			if err == nil {
				query = query.Where("acf.expense.date >= ? AND acf.expense.date <= ?", startDateParsed.Format("2006-01-02"), endDateParsed.Format("2006-01-02"))
			}
		}
	}

	if filter.MinBalance != nil {
		query = query.Where("acf.expense.balance >= ?", *filter.MinBalance)
	}

	if len(filter.CollectorIDs) > 0 {
		query = query.Where("acf.expense.collector_id IN ?", filter.CollectorIDs)
	}

	return query
}

func (repo *expenseEntryRepositoryImpl) buildSafeSort(sort string) string {
	sortMap := map[string]string{
		"created_date": "acf.expense.created_at",
		"date":         "acf.expense.date",
		"amount":       "acf.expense.amount",
		"balance":      "acf.expense.balance",
	}

	if strings.TrimSpace(sort) == "" {
		return "acf.expense.created_at DESC"
	}

	orders := make([]string, 0)
	for _, item := range strings.Split(sort, ",") {
		parts := strings.Split(strings.TrimSpace(item), ":")
		if len(parts) != 2 {
			continue
		}

		column, ok := sortMap[strings.TrimSpace(parts[0])]
		if !ok {
			continue
		}

		direction := strings.ToUpper(strings.TrimSpace(parts[1]))
		if direction != "ASC" && direction != "DESC" {
			continue
		}

		orders = append(orders, fmt.Sprintf("%s %s", column, direction))
	}

	if len(orders) == 0 {
		return "acf.expense.created_at DESC"
	}

	return strings.Join(orders, ", ")
}

func (repo *expenseEntryRepositoryImpl) FindById(ctx context.Context, custId string, expenseId int64) (model.Expense, error) {
	var expense model.Expense
	err := repo.model(ctx).
		Table("acf.expense").
		Where("acf.expense.cust_id = ? AND acf.expense.expense_id = ? AND acf.expense.is_del = false", custId, expenseId).
		Take(&expense).Error
	if err != nil {
		return expense, err
	}
	return expense, nil
}

func (repo *expenseEntryRepositoryImpl) FindDetailById(ctx context.Context, custId string, expenseId int64) (model.ExpenseList, error) {
	var row model.ExpenseList
	err := repo.model(ctx).Table("acf.expense").
		Select(`acf.expense.*, et.expense_type_code, et.expense_type_name, `+expenseEntryCollectorSelect+`,
			COALESCE(acf.expense.amount - COALESCE((
				SELECT COALESCE(SUM(payment_amount), 0)
				FROM acf.deposit_expense
				WHERE expense_id = acf.expense.expense_id
			), 0), acf.expense.amount) AS remaining_amount`).
		Joins("LEFT JOIN acf.expense_type et ON et.expense_type_id = acf.expense.expense_type_id").
		Joins(expenseEntryCollectorJoin, custId).
		Where("acf.expense.cust_id = ? AND acf.expense.expense_id = ? AND acf.expense.is_del = false", custId, expenseId).
		Take(&row).Error
	if err != nil {
		return row, err
	}
	return row, nil
}

func (repo *expenseEntryRepositoryImpl) FindDepositExpensesByExpenseId(ctx context.Context, custId string, expenseId int64) ([]model.DepositExpense, error) {
	var list []model.DepositExpense
	err := repo.model(ctx).Select(`
		acf.deposit_expense.*,
		e.doc_no,
		COALESCE(e.balance, 0) as balance,
		acf.deposit_expense.payment_amount
	`).Table("acf.deposit_expense").
		Joins("JOIN acf.deposit d ON d.deposit_no = acf.deposit_expense.deposit_no AND d.cust_id = ?", custId).
		Joins("JOIN acf.expense e ON e.expense_id = acf.deposit_expense.expense_id").
		Where("acf.deposit_expense.expense_id = ?", expenseId).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (repo *expenseEntryRepositoryImpl) FindExpenseFilesByExpenseId(ctx context.Context, custId string, expenseId int64) ([]model.ExpenseFile, error) {
	var list []model.ExpenseFile
	err := repo.model(ctx).Table("acf.expense_file").
		Where("cust_id = ? AND expense_id = ?", custId, expenseId).
		Order("expense_file_id ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (repo *expenseEntryRepositoryImpl) CountExpensesInCurrentMonth(ctx context.Context) (int, error) {
	var count int64
	err := repo.model(ctx).Table("acf.expense").
		Where("date_trunc('day', date) = date_trunc('day', CURRENT_DATE)").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *expenseEntryRepositoryImpl) Store(ctx context.Context, data *model.Expense) error {
	return repo.model(ctx).Table("acf.expense").Create(data).Error
}

func (repo *expenseEntryRepositoryImpl) StoreExpenseFiles(ctx context.Context, custId string, expenseId int64, files []model.ExpenseFile) error {
	if len(files) == 0 {
		return nil
	}
	return repo.model(ctx).Table("acf.expense_file").Create(&files).Error
}

func (repo *expenseEntryRepositoryImpl) Update(ctx context.Context, custId string, expenseId int64, data map[string]interface{}) error {
	result := repo.model(ctx).
		Table("acf.expense").
		Model(&model.Expense{}).
		Where("acf.expense.cust_id = ? AND acf.expense.expense_id = ? AND acf.expense.is_del = false", custId, expenseId).
		Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *expenseEntryRepositoryImpl) Delete(ctx context.Context, custId string, expenseId int64, deletedBy int64) error {
	now := time.Now()
	result := repo.model(ctx).
		Table("acf.expense").
		Model(&model.Expense{}).
		Where("acf.expense.cust_id = ? AND acf.expense.expense_id = ? AND acf.expense.is_del = false", custId, expenseId).
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
