package repository

import (
	"context"
	"log"
	"mobile/model"
	"time"

	"gorm.io/gorm"
)

// PrintDailyReportRepository interface for print daily report repository operations
type PrintDailyReportRepository interface {
	FindPaymentDataByCustIdAndDate(ctx context.Context, custID string, date time.Time) ([]model.Cndn, error)
	FindSalesDataByCustIdAndDate(ctx context.Context, userID int64, custID string, date time.Time) ([]model.PaymentTrxSalesData, error)
	FindExpenseDataByCustIdAndDate(ctx context.Context, custID string, date time.Time, userID int64) ([]model.ExpenseListRead, error)
	FindAttendanceByUserAndDate(ctx context.Context, userId int64, custId string, date time.Time) (int, error)
}

type printDailyReportRepositoryImpl struct {
	*gorm.DB
}

// NewPrintDailyReportRepository creates a new print daily report repository
func NewPrintDailyReportRepository(db *gorm.DB) PrintDailyReportRepository {
	return &printDailyReportRepositoryImpl{db}
}

func (repo *printDailyReportRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

// FindPaymentDataByCustIdAndDate finds payment data from acf.cndn where cust_id and cndn_jenis = "credit"
func (repo *printDailyReportRepositoryImpl) FindPaymentDataByCustIdAndDate(ctx context.Context, custID string, date time.Time) ([]model.Cndn, error) {
	var cndns []model.Cndn

	query := repo.model(ctx).
		Table("acf.cndn").
		Where("cust_id = ?", custID).
		Where("cndn_jenis = ?", "credit").
		Where("DATE(cndn_date) = ?", date.Format("2006-01-02")).
		Where("is_del = false")

	err := query.Find(&cndns).Error
	if err != nil {
		log.Println("PrintDailyReportRepository, FindPaymentDataByCustIdAndDate, err:", err.Error())
		return nil, err
	}

	return cndns, nil
}

// FindSalesDataByCustIdAndDate finds selling data from acf.payment_trx and details
func (repo *printDailyReportRepositoryImpl) FindSalesDataByCustIdAndDate(ctx context.Context, userID int64, custID string, date time.Time) ([]model.PaymentTrxSalesData, error) {
	var salesData []model.PaymentTrxSalesData

	query := repo.model(ctx).
		Table("acf.payment_trx pt").
		Select("ptd.amount as payment_amount, pt2.payment_type_code").
		Joins("inner join acf.payment_trx_detail ptd on ptd.payment_trx_id = pt.payment_trx_id").
		Joins("inner join acf.payment_type pt2 on pt2.payment_type_id = ptd.pay_type").
		Where("pt.emp_id = ?", userID).
		Where("pt.cust_id = ?", custID).
		Where("pt.date = ?", date.Format("2006-01-02")).
		Where("pt.is_del = ?", false).
		Where("ptd.is_del = ?", false).
		Where("pt.trx_source != ?", "L")

	err := query.Scan(&salesData).Error
	if err != nil {
		log.Println("PrintDailyReportRepository, FindSalesDataByCustIdAndDate, err:", err.Error())
		return nil, err
	}

	return salesData, nil
}

// FindExpenseDataByCustIdAndDate finds expense data from acf.expense join acf.expense_type
func (repo *printDailyReportRepositoryImpl) FindExpenseDataByCustIdAndDate(ctx context.Context, custID string, date time.Time, userID int64) ([]model.ExpenseListRead, error) {
	var expenses []model.ExpenseListRead

	query := repo.model(ctx).
		Table("acf.expense e").
		Select(`
			e.cust_id,
			e.expense_id,
			e.expense_type_id,
			et.expense_type_name,
			e.date,
			e.amount,
			e.note,
			e.created_at,
			e.updated_at
		`).
		Joins("LEFT JOIN acf.expense_type et ON et.expense_type_id = e.expense_type_id").
		Where("e.cust_id = ?", custID).
		Where("e.collector_id = ?", userID).
		Where("DATE(e.date) = ?", date.Format("2006-01-02")).
		Where("e.is_del = false")

	err := query.Scan(&expenses).Error
	if err != nil {
		log.Println("PrintDailyReportRepository, FindExpenseDataByCustIdAndDate, err:", err.Error())
		return nil, err
	}

	return expenses, nil
}

func (repo *printDailyReportRepositoryImpl) FindAttendanceByUserAndDate(ctx context.Context, userId int64, custId string, date time.Time) (int, error) {
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
