package repository

import (
	"context"
	"mobile/model"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryAttendanceImpl struct {
		*gorm.DB
	}
)
type AttendanceRepository interface {
	Store(data *model.Attendance) error
	FindBetween(custID, empCode string, start, end time.Time) (attendances []model.AttendanceRead, err error)
	ExistsAttendanceInDateRange(custID, empCode string, start, end time.Time) (bool, error)
	CountPlanBySalesmanIDAndDate(salesmanID int64, date time.Time, isDistributor bool) (count int64, err error)
	GetWarehouseStockByEmpIDAndCustID(empID int64, custID string) (stock int64, err error)
	GetSalesmanInfoWithCanvas(empID int64) (salesmanInfo SalesmanInfoWithCanvas, err error)
}

// SalesmanInfoWithCanvas represents salesman information with canvas data
type SalesmanInfoWithCanvas struct {
	EmpID         int64   `gorm:"column:emp_id"`
	EmpCode       *string `gorm:"column:emp_code"`
	EmpName       *string `gorm:"column:emp_name"`
	OprType       *string `gorm:"column:opr_type"`
	OprTypeCanvas *string `gorm:"column:opr_type_canvas"`
	WhID          *int64  `gorm:"column:wh_id"`
	WhCode        *string `gorm:"column:wh_code"`
	WhNameCanvas  *string `gorm:"column:wh_name_canvas"`
}

func NewAttendanceRepository(db *gorm.DB) *RepositoryAttendanceImpl {
	return &RepositoryAttendanceImpl{db}
}

func (repo *RepositoryAttendanceImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryAttendanceImpl) Store(data *model.Attendance) error {
	err := repository.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryAttendanceImpl) FindBetween(custID, empCode string, start, end time.Time) (attendances []model.AttendanceRead, err error) {
	err = repository.
		Where("cust_id = ? AND emp_code = ? AND created_at between ? AND ?", custID, empCode, start, end).
		Find(&attendances).Error
	return attendances, err
}

func (repository *RepositoryAttendanceImpl) CountPlanBySalesmanIDAndDate(salesmanID int64, date time.Time, isDistributor bool) (count int64, err error) {
	if isDistributor {
		err = repository.Table("pjp.route_pop_permanent rpp").
			Select("COUNT(rpp.id)").
			Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = rpp.pjp_id").
			Where("pjp.salesman_id = ? AND rpp.day=  TO_CHAR(NOW(), 'Dy')", salesmanID).
			Where("pjp.approval_status IN ('Approved','approved') AND pjp.status = 'true'").
			Count(&count).Error
	} else {
		err = repository.Table("pjp_principles.route_pop_permanent rpp").
			Select("COUNT(rpp.id)").
			Joins("JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = rpp.pjp_id").
			Where("pjp.salesman_id = ? AND rpp.day= TO_CHAR(NOW(), 'Dy')", salesmanID).
			Where("pjp.approval_status IN ('Approved','approved') AND pjp.status = 'true'").
			Count(&count).Error
	}
	return count, err
}

func (repository *RepositoryAttendanceImpl) ExistsAttendanceInDateRange(custID, empCode string, start, end time.Time) (bool, error) {
	var attendanceID int64
	err := repository.
		Table("mobile.attendances a").
		Select("a.attendance_id").
		Where("a.emp_code = ? AND a.cust_id = ?", empCode, custID).
		Where("a.created_at::date >= ? AND a.created_at::date <= ?", start, end).
		Limit(1).
		Scan(&attendanceID).Error
	if err != nil {
		return false, err
	}
	return attendanceID > 0, nil
}

func (repository *RepositoryAttendanceImpl) GetWarehouseStockByEmpIDAndCustID(empID int64, custID string) (stock int64, err error) {

	type StockResult struct {
		WhID *int64   `gorm:"column:wh_id"`
		Qty  *float64 `gorm:"column:qty"`
	}
	var results []StockResult

	err = repository.Table("mst.m_salesman ms").
		Select("msc.wh_id, COALESCE(ws.qty, 0) AS qty").
		Joins("LEFT JOIN mst.m_salesman_canvas msc ON msc.emp_id = ms.emp_id AND msc.is_active = true").
		Joins("LEFT JOIN inv.warehouse_stock ws ON msc.wh_id = ws.wh_id AND ws.cust_id = ?", custID).
		Where("ms.emp_id = ?", empID).
		Scan(&results).Error

	if err != nil {
		return 0, err
	}

	var totalStock float64
	for _, result := range results {
		if result.WhID != nil && result.Qty != nil {
			totalStock += *result.Qty
		}
	}

	return int64(totalStock), nil
}

// GetSalesmanInfoWithCanvas retrieves salesman information with canvas data
// opr_type: "O" if is_taking_order = true, "" if false
// opr_type_canvas: "C" if is_active = true, "" if false or NULL
func (repository *RepositoryAttendanceImpl) GetSalesmanInfoWithCanvas(empID int64) (salesmanInfo SalesmanInfoWithCanvas, err error) {
	err = repository.Table("mst.m_salesman ms").
		Select(`
			ms.emp_id,
			me.emp_code,
			me.emp_name,
			CASE 
				WHEN ms.is_taking_order = true THEN 'O'
				ELSE ''
			END AS opr_type,
			CASE 
				WHEN msc.is_active = true THEN 'C'
				ELSE ''
			END AS opr_type_canvas,
			msc.wh_id,
			mw.wh_code,
			mw.wh_name AS wh_name_canvas
		`).
		Joins("LEFT JOIN mst.m_employee me ON me.emp_id = ms.emp_id").
		Joins("LEFT JOIN mst.m_salesman_canvas msc ON msc.emp_id = ms.emp_id").
		Joins("LEFT JOIN mst.m_warehouse mw ON mw.wh_id = msc.wh_id").
		Where("ms.emp_id = ?", empID).
		Scan(&salesmanInfo).Error

	return salesmanInfo, err
}
