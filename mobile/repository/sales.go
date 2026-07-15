package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	RepositorySalesImpl struct {
		*gorm.DB
	}
)

type SalesRepository interface {
	GetLatestClockInDate(ctx context.Context, custID string, empID int64) (*time.Time, error)
	GetTotalOrder(ctx context.Context, custID string, salesmanID int64, startDate, endDate time.Time) (int64, error)
	GetTotalReturn(ctx context.Context, custID string, salesmanID int64, startDate, endDate time.Time) (int64, error)
	GetMonthlySalesTarget(ctx context.Context, custID string, salesmanID int64, month, year int) (int64, error)
}

// NewSalesRepository creates a new instance of SalesRepository
func NewSalesRepository(db *gorm.DB) *RepositorySalesImpl {
	return &RepositorySalesImpl{db}
}

// model returns a GORM DB instance with transaction support from context if available
func (repo *RepositorySalesImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

// GetLatestClockInDate retrieves the latest clock in date (type = 1) for a specific employee.
// It queries mobile.attendances table and joins with mst.m_employee to get emp_code from emp_id.
// Returns nil if no clock in record found.
func (repo *RepositorySalesImpl) GetLatestClockInDate(ctx context.Context, custID string, empID int64) (*time.Time, error) {
	var latestDate *time.Time
	err := repo.model(ctx).
		Table("mobile.attendances").
		Select("MAX(created_at)").
		Where("cust_id = ? AND emp_code = (SELECT emp_code FROM mst.m_employee WHERE emp_id = ? AND cust_id = ?) AND type = ?", custID, empID, custID, 1).
		Row().Scan(&latestDate)

	if err != nil {
		return nil, err
	}

	return latestDate, nil
}

// GetTotalOrder calculates the sum of total order amount for a specific salesman within a date range.
// Filters by cust_id, salesman_id, ro_date between startDate and endDate, and excludes deleted records (is_del = false).
// Returns 0 if no orders found or on error.
func (repo *RepositorySalesImpl) GetTotalOrder(ctx context.Context, custID string, salesmanID int64, startDate, endDate time.Time) (int64, error) {
	var total int64
	err := repo.model(ctx).
		Table("sls.order").
		Select("CAST(COALESCE(SUM(total), 0) AS BIGINT)").
		Where("cust_id = ? AND salesman_id = ? AND ro_date BETWEEN ? AND ? AND is_del = ? AND invoice_no IS NOT NULL", custID, salesmanID, startDate, endDate, false).
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetTotalReturn calculates the sum of total return amount for a specific salesman within a date range.
// Filters by cust_id, salesman_id, return_date between startDate and endDate, and excludes deleted records (is_del = false).
// Returns 0 if no returns found or on error.
func (repo *RepositorySalesImpl) GetTotalReturn(ctx context.Context, custID string, salesmanID int64, startDate, endDate time.Time) (int64, error) {
	var total int64
	// Filter returns based on invoice_no that exists in the orders within the date range
	subQuery := repo.model(ctx).
		Table("sls.order").
		Select("invoice_no").
		Where("cust_id = ? AND salesman_id = ? AND ro_date BETWEEN ? AND ? AND is_del = ? AND invoice_no IS NOT NULL", custID, salesmanID, startDate, endDate, false)

	err := repo.model(ctx).
		Table("sls.return").
		Select("CAST(COALESCE(SUM(total), 0) AS BIGINT)").
		Where("invoice_no IN (?)", subQuery).
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetMonthlySalesTarget retrieves the monthly sales target for a specific salesman.
// Joins mst.m_sales_target with mst.m_sales_allocated to get the allocated target amount.
// Returns the SUM of allocated amounts if multiple records exist.
// Filters by month, year, salesman_id, cust_id, and active/non-deleted records with status = 1.
func (repo *RepositorySalesImpl) GetMonthlySalesTarget(ctx context.Context, custID string, salesmanID int64, month, year int) (int64, error) {
	var total int64
	err := repo.model(ctx).
		Table("mst.m_sales_allocated msa").
		Select("CAST(COALESCE(SUM(msa.allocated), 0) AS BIGINT)").
		Joins("JOIN mst.m_sales_target mst ON msa.sales_team_id = mst.sales_target_id").
		Where("mst.month = ? AND mst.year = ? AND msa.salesman_id = ? AND msa.is_del = ? AND msa.is_active = ? AND mst.status = ? AND mst.cust_id = ?",
			month, year, salesmanID, false, true, 1, custID).
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}
