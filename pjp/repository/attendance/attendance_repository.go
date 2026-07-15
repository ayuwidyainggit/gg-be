package attendance

import (
	"context"
	"scyllax-pjp/constant"
	"time"

	"gorm.io/gorm"
)

// AttendanceRepository defines contract for attendance data access
type AttendanceRepository interface {
	GetDistributorPlanCount(ctx context.Context, tx *gorm.DB, salesmanID int, date time.Time) int
	GetPrincipalPlanCount(ctx context.Context, tx *gorm.DB, salesmanID int, date time.Time) int
	GetSalesmanWithCanvas(ctx context.Context, tx *gorm.DB, empID int) (*SalesmanCanvasInfo, error)
	GetWarehouseStock(ctx context.Context, tx *gorm.DB, empID int) (whID int, qty int, err error)
}

// SalesmanCanvasInfo holds salesman information with canvas details
type SalesmanCanvasInfo struct {
	EmpID          int
	EmpCode        string
	EmpName        string
	OprType        string
	OprTypeCanvas  *string
	IsActiveCanvas bool
	WhID           *int
	WhCode         *string
	WhName         *string
}

type attendanceRepository struct{}

// NewAttendanceRepository creates a new attendance repository instance
func NewAttendanceRepository() AttendanceRepository {
	return &attendanceRepository{}
}

// GetDistributorPlanCount returns the count of route plans for distributor salesman
// Query based on docs: pjp.route_pop_permanent + pjp.permanent_journey_plans
func (r *attendanceRepository) GetDistributorPlanCount(ctx context.Context, tx *gorm.DB, salesmanID int, date time.Time) int {
	var count int64

	tx.WithContext(ctx).
		Table("pjp.route_pop_permanent rpp").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = rpp.pjp_id").
		Where("pjp.salesman_id = ? AND rpp.date = ?", salesmanID, date.Format(constant.DateFormat)).
		Count(&count)

	return int(count)
}

// GetPrincipalPlanCount returns the count of route plans for principal salesman
// Query based on docs: pjp_principles.route_pop_permanent + pjp_principles.permanent_journey_plans
func (r *attendanceRepository) GetPrincipalPlanCount(ctx context.Context, tx *gorm.DB, salesmanID int, date time.Time) int {
	var count int64

	tx.WithContext(ctx).
		Table("pjp_principles.route_pop_permanent rpp").
		Joins("JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = rpp.pjp_id").
		Where("pjp.salesman_id = ? AND rpp.date = ?", salesmanID, date.Format(constant.DateFormat)).
		Count(&count)

	return int(count)
}

// GetSalesmanWithCanvas returns salesman info with canvas details
// Joins employee, salesman, salesman_canvas and warehouse tables
// Note: Citus requires join on distribution column (cust_id) for complex joins
func (r *attendanceRepository) GetSalesmanWithCanvas(ctx context.Context, tx *gorm.DB, empID int) (*SalesmanCanvasInfo, error) {
	var info SalesmanCanvasInfo

	query := `
		SELECT 
			me.emp_id,
			me.emp_code,
			me.emp_name,
			ms.opr_type,
			msc.opr_type as opr_type_canvas,
			COALESCE(msc.is_active, false) as is_active_canvas,
			msc.wh_id,
			mw.wh_code,
			mw.wh_name
		FROM mst.m_employee me
		LEFT JOIN mst.m_salesman ms ON ms.emp_id = me.emp_id AND ms.cust_id = me.cust_id
		LEFT JOIN mst.m_salesman_canvas msc ON msc.emp_id = me.emp_id AND msc.cust_id = me.cust_id
		LEFT JOIN mst.m_warehouse mw ON mw.wh_id = msc.wh_id AND mw.cust_id = me.cust_id
		WHERE me.emp_id = ?
		LIMIT 1
	`

	if err := tx.WithContext(ctx).Raw(query, empID).Scan(&info).Error; err != nil {
		return nil, err
	}

	return &info, nil
}

// GetWarehouseStock returns warehouse ID and stock quantity for canvas salesman
// Query based on docs: mst.m_salesman LEFT JOIN mst.m_salesman_canvas LEFT JOIN inv.warehouse_stock
// Note: Citus requires join on distribution column (cust_id) for complex joins
func (r *attendanceRepository) GetWarehouseStock(ctx context.Context, tx *gorm.DB, empID int) (int, int, error) {
	var result struct {
		WhID int
		Qty  float64
	}

	query := `
		SELECT 
			msc.wh_id, 
			COALESCE(SUM(ws.qty), 0) AS qty
		FROM mst.m_salesman ms
		LEFT JOIN mst.m_salesman_canvas msc ON msc.emp_id = ms.emp_id AND msc.cust_id = ms.cust_id
		LEFT JOIN inv.warehouse_stock ws ON msc.wh_id = ws.wh_id AND ws.cust_id = ms.cust_id
		WHERE ms.emp_id = ?
		GROUP BY msc.wh_id
		LIMIT 1
	`

	if err := tx.WithContext(ctx).Raw(query, empID).Scan(&result).Error; err != nil {
		return 0, 0, err
	}

	return result.WhID, int(result.Qty), nil
}
