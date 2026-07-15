package live_monitoring

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

// LiveMonitoringRepository defines the interface for live monitoring data access
type LiveMonitoringRepository interface {
	// Principal monitoring
	GetPrincipalEmployeeIDs(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) ([]int, error)
	GetPrincipalMonitoring(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string, limit, offset int) ([]model.LiveMonitoringPrincipalRow, error)
	GetPrincipalExtraCalls(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) ([]model.LiveMonitoringPrincipalRow, error)
	CountPrincipalMonitoring(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) (int64, error)

	// Distributor monitoring
	GetDistributorEmployeeIDs(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) ([]int, error)
	GetDistributorMonitoring(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string, limit, offset int) ([]model.LiveMonitoringDistributorRow, error)
	GetDistributorLatestVisitCoordinates(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empIDs []int) (map[string]model.LatestVisitCoordinateRow, error)
	GetDistributorEmployeeMeta(ctx context.Context, tx *gorm.DB, custIDs []string, empIDs []int) (map[int]model.DistributorEmployeeMetaRow, error)
	GetDistributorRouteMeta(ctx context.Context, tx *gorm.DB, custIDs []string, routeCodes []int64) (map[string]model.DistributorRouteMetaRow, error)
	GetDistributorOutletMeta(ctx context.Context, tx *gorm.DB, custIDs []string, outletIDs []int) (map[string]model.DistributorOutletMetaRow, error)
	GetDistributorAttendance(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) (map[int]model.AttendanceRow, error)
	GetDistributorCurrentCoordinates(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) (map[int]model.CurrentCoordinateRow, error)
	CountDistributorMonitoring(ctx context.Context, tx *gorm.DB, custIDs []string, date string, regionID, areaID, distributorID int, empIDs []int, statuses []string) (int64, error)

	// Detail queries
	GetVisitInformationPrincipal(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) (*model.VisitInformationRow, error)
	GetVisitInformationPrincipalFromHistory(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) (*model.VisitInformationRow, error)
	CountTotalVisitsPrincipal(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) (int64, error)
	GetVisitInformationDistributor(ctx context.Context, tx *gorm.DB, date string, empID, distributorID int) (*model.VisitInformationRow, error)
	CountDistributorPlannedVisits(ctx context.Context, tx *gorm.DB, date string, empID, distributorID int) (int64, error)
	CountDistributorExtraCalls(ctx context.Context, tx *gorm.DB, date string, empID, distributorID int) (int64, error)
	CountDistributorOnGoingVisits(ctx context.Context, tx *gorm.DB, date string, empID, distributorID int) (int64, error)
	CountDistributorVisitedVisits(ctx context.Context, tx *gorm.DB, date string, empID, distributorID int) (int64, error)
	CountDistributorSkippedVisits(ctx context.Context, tx *gorm.DB, date string, empID, distributorID int) (int64, error)
	GetSales(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.SalesRow, error)
	GetReturns(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.ReturnRow, error)
	GetCollections(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.CollectionRow, error)
	GetExpenses(ctx context.Context, tx *gorm.DB, custID string, empID int, date string) ([]model.ExpenseRow, error)
	GetShipments(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.ShipmentRow, error)
	GetSubmittedSurveyData(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.SurveyDataRow, error)
	GetActivityTime(ctx context.Context, tx *gorm.DB, date string, empID int) (*string, error)
	GetDistributorInfo(ctx context.Context, tx *gorm.DB, distributorID int) (*model.DistributorInfoRow, error)
	GetUserFullname(ctx context.Context, tx *gorm.DB, custID string) (*string, error)
	GetChildCustIDs(ctx context.Context, tx *gorm.DB, parentCustID string) ([]string, error)
	GetSalesmanCustID(ctx context.Context, tx *gorm.DB, empID int) (string, error)

	// Update locations
	GetEmployeeRole(ctx context.Context, tx *gorm.DB, empID int, jwtCust string) (string, error)
	GetUpdateLocations(ctx context.Context, tx *gorm.DB, empID int, date string, jwtCust string, branch string) ([]model.UpdateLocationRow, error)
}

type liveMonitoringRepository struct{}

// NewLiveMonitoringRepository creates a new instance of LiveMonitoringRepository
func NewLiveMonitoringRepository() LiveMonitoringRepository {
	return &liveMonitoringRepository{}
}
