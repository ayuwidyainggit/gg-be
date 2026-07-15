package destination

import (
	"context"
	"fmt"
	"scyllax-pjp/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DestinationRepository interface {
	FindAllOutletsByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) []model.Destination
	FindAllOutletsByPjpCode(ctx context.Context, tx *gorm.DB, pjpCode int, custId string) []model.Destination
	CreateBulk(ctx context.Context, tx *gorm.DB, outlets []model.Destination)
	GetById(ctx context.Context, tx *gorm.DB, id int, custId string) model.Destination
	UpdatePivot(ctx context.Context, tx *gorm.DB, route model.Destination)
	GetAll(ctx context.Context, tx *gorm.DB, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]model.Destination, int)
	GetAllEnhance(ctx context.Context, tx *gorm.DB, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]model.Pjp, int)
}

type destinationRepository struct{}

func NewDestinationRepository() DestinationRepository {
	return &destinationRepository{}
}

func applyFilter(query *gorm.DB, field string, value interface{}) *gorm.DB {
	switch field {
	case "salesman_name":
		if v, ok := value.(string); ok && v != "" {
			return query.Joins("LEFT JOIN pjp_principles.permanent_journey_plans ON pjp_principles.destinations.pjp_id = pjp_principles.permanent_journey_plans.id").
				Where("pjp_principles.permanent_journey_plans.salesman_name = ?", v)
		}
	case "salesman_code":
		if v, ok := value.(string); ok && v != "" {
			codes := parseCommaSeparated(v)
			return query.Joins("LEFT JOIN pjp_principles.permanent_journey_plans ON pjp_principles.destinations.pjp_id = pjp_principles.permanent_journey_plans.id").
				Where("pjp_principles.permanent_journey_plans.salesman_code IN ?", codes)
		}
	case "status":
		if v, ok := value.(string); ok && v != "" {
			return query.Where("status IN ?", parseCommaSeparated(v))
		}
	case "start_date":
		if v, ok := value.(time.Time); ok && !v.IsZero() {
			return query.Where("created_at >= ?", v)
		}
	case "end_date":
		if v, ok := value.(time.Time); ok && !v.IsZero() {
			return query.Where("created_at <= ?", v.AddDate(0, 0, 1))
		}
	default:
		switch v := value.(type) {
		case string:
			if v != "" {
				return query.Where(fmt.Sprintf("%s IN ?", field), parseCommaSeparated(v))
			}
		case int:
			if v != 0 {
				return query.Where(fmt.Sprintf("%s = ?", field), v)
			}
		}
	}
	return query
}

func parseCommaSeparated(val string) []string {
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
