package live_monitoring

import (
	"context"
	"fmt"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	livemonitoringrepo "scyllax-pjp/repository/live_monitoring"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var jakartaLocation = loadJakartaLocation()

// LiveMonitoringService defines the interface for live monitoring business logic
type LiveMonitoringService interface {
	// Principal monitoring
	GetPrincipalMonitoring(ctx context.Context, req request.LiveMonitoringRequest, custID string) ([]response.LiveMonitoringData, response.LiveMonitoringPaging, error)

	// Distributor monitoring
	GetDistributorMonitoring(ctx context.Context, req request.LiveMonitoringRequest, custID string) ([]response.LiveMonitoringData, response.LiveMonitoringPaging, error)

	// Detail monitoring
	GetMonitoringDetail(ctx context.Context, req request.LiveMonitoringDetailRequest, custID string, userID int64) (*response.LiveMonitoringDetailData, error)
	GetUpdateLocations(ctx context.Context, req request.UpdateLocationsRequest, custID string) (response.UpdateLocationsResponse, error)
}

type liveMonitoringService struct {
	repository livemonitoringrepo.LiveMonitoringRepository
	validate   *validator.Validate
	db         *gorm.DB
}

// NewLiveMonitoringService creates a new instance of LiveMonitoringService
func NewLiveMonitoringService(
	repo livemonitoringrepo.LiveMonitoringRepository,
	validate *validator.Validate,
	db *gorm.DB,
) LiveMonitoringService {
	return &liveMonitoringService{
		repository: repo,
		validate:   validate,
		db:         db,
	}
}

// epochToDateString converts epoch timestamp to YYYY-MM-DD format
func epochToDateString(epoch int64) string {
	t := time.Unix(epoch, 0).In(jakartaLocation)
	return t.Format("2006-01-02")
}

func buildJakartaBusinessDayEpochRange(date string) (int64, int64, error) {
	startAt, err := time.ParseInLocation("2006-01-02", date, jakartaLocation)
	if err != nil {
		return 0, 0, fmt.Errorf("parse business date: %w", err)
	}

	endAt := startAt.AddDate(0, 0, 1)

	return startAt.Unix(), endAt.Unix(), nil
}

func isBusinessDayEpoch(timestamp *int64, dayStartEpoch, dayEndEpoch int64) bool {
	if timestamp == nil {
		return false
	}

	normalizedTimestamp := normalizeBusinessEpoch(*timestamp)

	return normalizedTimestamp >= dayStartEpoch && normalizedTimestamp < dayEndEpoch
}

func normalizeBusinessEpoch(timestamp int64) int64 {
	if timestamp > 9999999999 || timestamp < -9999999999 {
		return timestamp / 1000
	}

	return timestamp
}

func loadJakartaLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.FixedZone("WIB", 7*60*60)
	}

	return loc
}

// calculatePaging calculates pagination information
func calculatePaging(totalRecord int64, page, limit int) response.LiveMonitoringPaging {
	totalPage := int(totalRecord) / limit
	if int(totalRecord)%limit > 0 {
		totalPage++
	}
	return response.LiveMonitoringPaging{
		TotalRecord: int(totalRecord),
		PageCurrent: page,
		PageLimit:   limit,
		PageTotal:   totalPage,
	}
}
